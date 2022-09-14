package sidecar

import (
	"context"
	pitaya "github.com/topfreegames/pitaya/v2/pkg"
	"github.com/topfreegames/pitaya/v2/pkg/cluster"
	"github.com/topfreegames/pitaya/v2/pkg/config"
	"github.com/topfreegames/pitaya/v2/pkg/constants"
	"github.com/topfreegames/pitaya/v2/pkg/errors"
	"google.golang.org/protobuf/proto"

	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"sync/atomic"

	"github.com/opentracing/opentracing-go"
	"github.com/topfreegames/pitaya/v2/pkg/logger"
	"github.com/topfreegames/pitaya/v2/pkg/logger/logrus"
	"github.com/topfreegames/pitaya/v2/pkg/protos"
	"github.com/topfreegames/pitaya/v2/pkg/tracing"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	clusterprotos "github.com/topfreegames/pitaya/v2/examples/demo/protos"
)

// TODO: implement jaeger into this, test everything, if connection dies this
// will go to hell, reconnection doesnt work, what can we do?  Bench results:
// 40k req/sec bidirectional stream, around 13k req/sec (RPC)

// TODO investigate why I will get drops in send rate every now and then during
// the benchmark test. I imagine maybe it's due to garbage collection?

// TODO I can utilize reutilizable objects, such as with a poll and reduce
// allocations here

// TODO separate sidecar and sidecar server into isolated files

// Sidecar main struct to keep state
type Sidecar struct {
	config        *config.BuilderConfig
	sidecarServer protos.SidecarServer
	stopChan      chan bool
	callChan      chan (*Call)
	sdChan        chan (*protos.SDEvent)
	shouldRun     bool
	listener      net.Listener
	debug bool
}

// Call struct represents an incoming RPC call from other servers
type Call struct {
	ctx   context.Context
	req   *protos.Request
	done  chan (bool)
	res   *protos.Response
	err   *protos.Error
	reqId uint64
}

// SidecarServer is the implementation of the GRPC server used to communicate
// with the sidecar client
type SidecarServer struct {
	protos.UnimplementedPitayaServer
	pitayaApp   pitaya.Pitaya
}

var (
	sidecar = &Sidecar{
		stopChan:  make(chan bool),
		callChan:  make(chan *Call),
		sdChan:    make(chan *protos.SDEvent),
		shouldRun: true,
	}
	reqId    uint64
	reqMutex sync.RWMutex
	reqMap   = make(map[uint64]*Call)
	empty    = &emptypb.Empty{}
)

// AddServer is called by the ServiceDiscovery when a new  pitaya server is
// added. We have it here so that we stream add and removed servers to sidecar
// client.
func (s *Sidecar) AddServer(server *cluster.Server) {
	s.sdChan <- &protos.SDEvent{
		Server: &protos.Server{
			Id:       server.ID,
			Frontend: server.Frontend,
			Type:     server.Type,
			Metadata: server.Metadata,
			Hostname: server.Hostname,
		},
		Event: protos.SDEvent_ADD,
	}
}

// RemoveServer is called by the ServiceDiscovery when a pitaya server is
// removed from the cluster.  We have it here so that we stream add and removed
// servers to sidecar client.
func (s *Sidecar) RemoveServer(server *cluster.Server) {
	s.sdChan <- &protos.SDEvent{
		Server: &protos.Server{
			Id:       server.ID,
			Frontend: server.Frontend,
			Type:     server.Type,
			Metadata: server.Metadata,
			Hostname: server.Hostname,
		},
		Event: protos.SDEvent_REMOVE,
	}
}

// Call receives an RPC request from other pitaya servers and forward it to the
// sidecar client so that it processes them, afterwards it gets the client
// response and send it back to the callee
func (s *Sidecar) Call(ctx context.Context, req *protos.Request) (*protos.Response, error) {
	call := &Call{
		ctx:   ctx,
		req:   req,
		done:  make(chan (bool), 1),
		reqId: atomic.AddUint64(&reqId, 1),
	}

	reqMutex.Lock()
	reqMap[call.reqId] = call
	reqMutex.Unlock()

	logger.Log.Info("Message on route:", req.Msg.Route)

	var msg clusterprotos.RPCMsg
	proto.Unmarshal(req.Msg.GetData(), &msg)

	logger.Log.Info("Message data:", msg.Msg)

	s.callChan <- call

	defer func() {
		reqMutex.Lock()
		delete(reqMap, call.reqId)
		reqMutex.Unlock()
	}()

	select {
	case <-call.done:
		return call.res, nil
	case <-time.After(sidecar.config.Pitaya.Sidecar.CallTimeout):
		close(call.done)
		return &protos.Response{}, constants.ErrSidecarCallTimeout
	}
}

// finishRPC is called when the sidecar client returns the answer to an RPC
// call, after this method happens, the Call method above returns
func (s *SidecarServer) finishRPC(ctx context.Context, res *protos.RPCResponse) {
	reqMutex.RLock()
	defer reqMutex.RUnlock()
	call, ok := reqMap[res.ReqId]
	if ok {
		logger.Log.Info("Response on route:", call.req.Msg.Route)

		var msg clusterprotos.RPCRes
		proto.Unmarshal(res.Res.Data, &msg)

		logger.Log.Info("Response data:", msg.Msg)
		call.res = res.Res
		call.err = res.Err
		close(call.done)
	}
}

// SessionBindRemote is meant to frontend servers so its not implemented here
func (s *Sidecar) SessionBindRemote(ctx context.Context, msg *protos.BindMsg) (*protos.Response, error) {
	return nil, constants.ErrNotImplemented
}

// PushToUser is meant to frontend servers so its not implemented here
func (s *Sidecar) PushToUser(ctx context.Context, push *protos.Push) (*protos.Response, error) {
	return nil, constants.ErrNotImplemented
}

// KickUser is meant to frontend servers so its not implemented here
func (s *Sidecar) KickUser(ctx context.Context, kick *protos.KickMsg) (*protos.KickAnswer, error) {
	return nil, constants.ErrNotImplemented
}

// GetServer is called by the sidecar client to get the information from
// a pitaya server by passing its ID
func (s *SidecarServer) GetServer(ctx context.Context, in *protos.Server) (*protos.Server, error) {
	server, err := pitaya.GetServerByID(in.Id)
	if err != nil {
		return nil, err
	}
	res := &protos.Server{
		Id:       server.ID,
		Frontend: server.Frontend,
		Type:     server.Type,
		Metadata: server.Metadata,
		Hostname: server.Hostname,
	}
	return res, nil
}

func (s *SidecarServer) Heartbeat(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	logger.Log.Debug("Received heartbeat from the sidecar client")
	return empty, nil
}

// ListenSD keeps a stream open between the sidecar client and server, it sends
// service discovery events to the sidecar client, such as add or removal of
// servers in the cluster
func (s *SidecarServer) ListenSD(empty *emptypb.Empty, stream protos.Sidecar_ListenSDServer) error {
	for sidecar.shouldRun {
		select {
		case evt := <-sidecar.sdChan:
			err := stream.Send(evt)
			if err != nil {
				logger.Log.Warnf("error sending sd event to sidecar client: %s", err.Error)
			}
		case <-sidecar.stopChan:
			sidecar.shouldRun = false
		}
	}
	logger.Log.Info("exiting sidecar ListenSD routine because stopChan was closed")
	return nil
}

// ListenRPC keeps a bidirectional stream open between the sidecar client and
// server, it sends incoming RPC from other pitaya servers to the client and
// also listens for incoming answers from the client. This method is the most
// important one here and is where is defined our async model.
func (s *SidecarServer) ListenRPC(stream protos.Sidecar_ListenRPCServer) error {
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				select {
				case <-sidecar.stopChan:
					return
				default:
					close(sidecar.stopChan)
					return
				}
			}
			// TODO fix context to fix tracing
			s.finishRPC(context.Background(), res)
		}
	}()
	for sidecar.shouldRun {
		select {
		case call := <-sidecar.callChan:
			err := stream.Send(&protos.SidecarRequest{ReqId: call.reqId, Req: call.req})
			if err != nil {
				call.err = &protos.Error{Msg: err.Error(), Code: errors.ErrBadRequestCode}
				close(call.done)
			}
		case <-sidecar.stopChan:
			sidecar.shouldRun = false
		}
	}
	logger.Log.Info("exiting sidecar ListenRPC routine because stopChan was closed")
	return nil
}

func getCtxWithParentSpan(ctx context.Context, op string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		carrier := opentracing.HTTPHeadersCarrier(md)
		spanContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			logger.Log.Debugf("tracing: could not extract span from context!")
		} else {
			return tracing.StartSpan(ctx, op, opentracing.Tags{}, spanContext)
		}
	}
	return ctx
}

// SendRPC is called by the sidecar client when it wants to send RPC requests to
// other pitaya servers
func (s *SidecarServer) SendRPC(ctx context.Context, in *protos.RequestTo) (*protos.Response, error) {
	pCtx := getCtxWithParentSpan(ctx, in.Msg.Route)
	ret, err := s.pitayaApp.RawRPC(pCtx, in.ServerID, in.Msg.Route, in.Msg.Data)
	defer tracing.FinishSpan(pCtx, err)
	return ret, err
}

// SendPush is called by the sidecar client when it wants to send a push to an
// user through a frontend server
func (s *SidecarServer) SendPush(ctx context.Context, in *protos.PushRequest) (*protos.PushResponse, error) {
	push := in.GetPush()
	failedUids, err := pitaya.SendPushToUsers(push.Route, push.GetData(), []string{push.Uid}, in.FrontendType)
	res := &protos.PushResponse{
		FailedUids: failedUids,
	}
	if err != nil {
		res.HasFailed = true
	} else {
		res.HasFailed = false
	}
	return res, nil // can't send the error here because if we do, it will throw an exception in csharp side
}

// SendKick is called by the sidecar client when it wants to send a kick to an
// user through a frontend server
func (s *SidecarServer) SendKick(ctx context.Context, in *protos.KickRequest) (*protos.PushResponse, error) {
	failedUids, err := pitaya.SendKickToUsers([]string{in.GetKick().GetUserId()}, in.FrontendType)
	res := &protos.PushResponse{
		FailedUids: failedUids,
	}
	if err != nil {
		res.HasFailed = true
	} else {
		res.HasFailed = false
	}
	return res, nil // can't send the error here because if we do, it will throw an exception in csharp side
}

// StartPitaya instantiates a pitaya server and starts it. It must be called
// during the initialization of the sidecar client, all other methods will only
// work when this one was already called
func (s *SidecarServer) StartPitaya(ctx context.Context, req *protos.StartPitayaRequest) (*protos.Error, error) {
	pitayaconfig := req.GetConfig()

	builder := pitaya.NewDefaultBuilder(
		pitayaconfig.GetFrontend(),
		pitayaconfig.GetType(),
		pitaya.Cluster,
		pitayaconfig.GetMetadata(),
		*sidecar.config,
		)

	builder.ServiceDiscovery.AddListener(sidecar)

	// register the sidecar as the pitaya server so that calls will be delivered
	// here and we can forward to the remote process
	builder.RPCServer.SetPitayaServer(sidecar)

	s.pitayaApp = builder.Build()

	// Start our own logger
	log := logrus.New()

	s.pitayaApp.SetDebug(req.GetDebugLog())

	pitaya.SetLogger(log.WithField("source", "sidecar"))

	// TODO support frontend servers
	if pitayaconfig.GetFrontend() {
		//t := acceptor.NewTCPAcceptor(":3250") pitaya.AddAcceptor(t)
		logger.Log.Fatal("Frontend servers not supported yet")
	}

	// TODO maybe we should return error in pitaya start. maybe recover from fatal
	// TODO make this method return error so that I can catch it
	go func() {
		s.pitayaApp.Start()
	}()
	return &protos.Error{}, nil
}

// StopPitaya stops the instantiated pitaya server and must always be called
// when the client is dying so that we can correctly gracefully shutdown pitaya
func (s *SidecarServer) StopPitaya(ctx context.Context, req *emptypb.Empty) (*protos.Error, error) {
	logger.Log.Info("received stop request, will stop pitaya server")
	select {
	case <-sidecar.stopChan:
		break
	default:
		close(sidecar.stopChan)
	}
	return &protos.Error{}, nil
}

func checkError(err error) {
	if err != nil {
		logger.Log.Fatalf("failed to start sidecar: %s", err)
	}
}

// TODO need to replace pitayas own jaeger config with something that actually makes sense
func configureJaeger(debug bool) {
	cfg, err := jaegercfg.FromEnv()
	if debug {
		cfg.ServiceName = "pitaya-sidecar"
		cfg.Sampler.Type = "const"
		cfg.Sampler.Param = 1
	}
	if cfg.ServiceName == "" {
		logger.Log.Error("Could not init jaeger tracer without ServiceName, either set environment JAEGER_SERVICE_NAME or cfg.ServiceName = \"my-api\"")
		return
	}
	if err != nil {
		logger.Log.Error("Could not parse Jaeger env vars: %s", err.Error())
		return
	}
	tracer, _, err := cfg.NewTracer()
	if err != nil {
		logger.Log.Error("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	opentracing.SetGlobalTracer(tracer)
	logger.Log.Infof("Tracer configured for %s", cfg.Reporter.LocalAgentHostPort)
}

// StartSidecar starts the sidecar server, it instantiates the GRPC server and
// listens for incoming client connections. This is the very first method that
// is called when the sidecar is starting.
func StartSidecar(cfg *config.BuilderConfig, debug bool, bindAddr, bindProto string) {
	sidecar.config = cfg
	sidecar.sidecarServer = &SidecarServer{}
	sidecar.debug = debug
	if bindProto != "unix" && bindProto != "tcp" {
		logger.Log.Fatal("only supported schemes are unix and tcp, review your bindaddr config")
	}

	var err error
	sidecar.listener, err = net.Listen(bindProto, bindAddr)
	checkError(err)
	defer sidecar.listener.Close()
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	protos.RegisterSidecarServer(grpcServer, sidecar.sidecarServer)
	go func() {
		err = grpcServer.Serve(sidecar.listener)
		if err != nil {
			logger.Log.Errorf("error serving GRPC: %s", err)
			select {
			case <-sidecar.stopChan:
				break
			default:
				close(sidecar.stopChan)
			}
		}
	}()

	// TODO: what to do if received sigint/term without receiving stop request from client?
	logger.Log.Infof("sidecar listening at %s", sidecar.listener.Addr())

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	// TODO make jaeger optional and configure with configs
	configureJaeger(true)

	// stop server
	select {
	case s := <-sg:
		logger.Log.Warn("got signal: ", s, ", shutting down...")
		close(sidecar.stopChan)
		break
	case <-sidecar.stopChan:
		logger.Log.Warn("the app will shutdown in a few seconds")
	}
	pitaya.Shutdown()
}
