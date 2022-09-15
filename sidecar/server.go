package sidecar

import (
	"context"
	pitaya "github.com/topfreegames/pitaya/v2/pkg"
	"github.com/topfreegames/pitaya/v2/pkg/errors"
	"github.com/topfreegames/pitaya/v2/pkg/logger"
	"github.com/topfreegames/pitaya/v2/pkg/logger/logrus"
	"github.com/topfreegames/pitaya/v2/pkg/protos"
	"github.com/topfreegames/pitaya/v2/pkg/tracing"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server is the implementation of the GRPC server used to communicate
// with the sidecar client
type Server struct {
	pitayaApp   pitaya.Pitaya
	sidecar *Sidecar
	protos.UnimplementedPitayaServer
}


// GetServer is called by the sidecar client to get the information from
// a pitaya server by passing its ID
func (s *Server) GetServer(ctx context.Context, in *protos.Server) (*protos.Server, error) {
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

func (s *Server) Heartbeat(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	logger.Log.Debug("Received heartbeat from the sidecar client")
	return empty, nil
}

// ListenSD keeps a stream open between the sidecar client and server, it sends
// service discovery events to the sidecar client, such as add or removal of
// servers in the cluster
func (s *Server) ListenSD(empty *emptypb.Empty, stream protos.Sidecar_ListenSDServer) error {
	for shouldRun {
		select {
		case evt := <-sdChan:
			err := stream.Send(evt)
			if err != nil {
				logger.Log.Warnf("error sending sd event to sidecar client: %s", err.Error)
			}
		case <-stopChan:
			shouldRun = false
		}
	}
	logger.Log.Info("exiting sidecar ListenSD routine because stopChan was closed")
	return nil
}

// ListenRPC keeps a bidirectional stream open between the sidecar client and
// server, it sends incoming RPC from other pitaya servers to the client and
// also listens for incoming answers from the client. This method is the most
// important one here and is where is defined our async model.
func (s *Server) ListenRPC(stream protos.Sidecar_ListenRPCServer) error {
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				select {
				case <-stopChan:
					return
				default:
					close(stopChan)
					return
				}
			}
			// TODO fix context to fix tracing
			s.finishRPC(context.Background(), res)
		}
	}()
	for shouldRun {
		select {
		case call := <- callChan:
			err := stream.Send(&protos.SidecarRequest{ReqId: call.reqId, Req: call.req})
			if err != nil {
				call.err = &protos.Error{Msg: err.Error(), Code: errors.ErrBadRequestCode}
				close(call.done)
			}
		case <-stopChan:
			shouldRun = false
		}
	}
	logger.Log.Info("exiting sidecar ListenRPC routine because stopChan was closed")
	return nil
}

// finishRPC is called when the sidecar client returns the answer to an RPC
// call, after this method happens, the Call method above returns
func (s *Server) finishRPC(ctx context.Context, res *protos.RPCResponse) {
	reqMutex.RLock()
	defer reqMutex.RUnlock()
	call, ok := reqMap[res.ReqId]
	if ok {
		call.res = res.Res
		call.err = res.Err
		close(call.done)
	}
}


// SendRPC is called by the sidecar client when it wants to send RPC requests to
// other pitaya servers
func (s *Server) SendRPC(ctx context.Context, in *protos.RequestTo) (*protos.Response, error) {
	pCtx := getCtxWithParentSpan(ctx, in.Msg.Route)
	ret, err := s.pitayaApp.RawRPC(pCtx, in.ServerID, in.Msg.Route, in.Msg.Data)
	defer tracing.FinishSpan(pCtx, err)
	return ret, err
}

// SendPush is called by the sidecar client when it wants to send a push to an
// user through a frontend server
func (s *Server) SendPush(ctx context.Context, in *protos.PushRequest) (*protos.PushResponse, error) {
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
func (s *Server) SendKick(ctx context.Context, in *protos.KickRequest) (*protos.PushResponse, error) {
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
func (s *Server) StartPitaya(ctx context.Context, req *protos.StartPitayaRequest) (*protos.Error, error) {
	pitayaconfig := req.GetConfig()

	builder := NewSidecarBuilder(
		pitayaconfig.GetFrontend(),
		pitayaconfig.GetType(),
		pitaya.Cluster,
		pitayaconfig.GetMetadata(),
		s.sidecar.config,
	)

	builder.ServiceDiscovery.AddListener(s.sidecar)
	// register the sidecar as the pitaya server so that calls will be delivered
	// here and we can forward to the remote process
	builder.SetPitayaServer(s.sidecar)

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
func (s *Server) StopPitaya(ctx context.Context, req *emptypb.Empty) (*protos.Error, error) {
	logger.Log.Info("received stop request, will stop pitaya server")
	select {
	case <-stopChan:
		break
	default:
		close(stopChan)
	}
	return &protos.Error{}, nil
}