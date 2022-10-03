package sidecar

import (
	"context"
	"github.com/topfreegames/pitaya/v2/pkg/cluster"
	"github.com/topfreegames/pitaya/v2/pkg/config"
	"github.com/topfreegames/pitaya/v2/pkg/constants"
	"time"

	"sync/atomic"

	"github.com/topfreegames/pitaya/v2/pkg/protos"
)

// TODO: implement jaeger into this, test everything, if connection dies this
// will go to hell, reconnection doesnt work, what can we do?  Bench results:
// 40k req/sec bidirectional stream, around 13k req/sec (RPC)

// TODO investigate why I will get drops in send rate every now and then during
// the benchmark test. I imagine maybe it's due to garbage collection?

// TODO I can utilize reutilizable objects, such as with a poll and reduce
// allocations here

// Remote main struct to keep state
type Remote struct {
	config config.SidecarConfig
}

func NewRemote(config config.SidecarConfig) *Remote {
	return &Remote{
		config : config,
	}
}

// AddServer is called by the ServiceDiscovery when a new  pitaya server is
// added. We have it here so that we stream add and removed servers to sidecar
// client.
func (s *Remote) AddServer(server *cluster.Server) {
	sdChan <- &protos.SDEvent{
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
func (s *Remote) RemoveServer(server *cluster.Server) {
	sdChan <- &protos.SDEvent{
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
func (s *Remote) Call(ctx context.Context, req *protos.Request) (*protos.Response, error) {
	call := &Call{
		ctx:   ctx,
		req:   req,
		done:  make(chan (bool), 1),
		reqId: atomic.AddUint64(&redId, 1),
	}

	reqMutex.Lock()
	reqMap[call.reqId] = call
	reqMutex.Unlock()

	callChan <- call

	defer func() {
		reqMutex.Lock()
		delete(reqMap, call.reqId)
		reqMutex.Unlock()
	}()

	select {
	case <-call.done:
		return call.res, nil
	case <-time.After(s.config.CallTimeout):
		close(call.done)
		return &protos.Response{}, constants.ErrSidecarCallTimeout
	}
}

// SessionBindRemote is meant to frontend servers so its not implemented here
func (s *Remote) SessionBindRemote(ctx context.Context, msg *protos.BindMsg) (*protos.Response, error) {
	return nil, constants.ErrNotImplemented
}

// PushToUser is meant to frontend servers so its not implemented here
func (s *Remote) PushToUser(ctx context.Context, push *protos.Push) (*protos.Response, error) {
	return nil, constants.ErrNotImplemented
}

// KickUser is meant to frontend servers so its not implemented here
func (s *Remote) KickUser(ctx context.Context, kick *protos.KickMsg) (*protos.KickAnswer, error) {
	return nil, constants.ErrNotImplemented
}