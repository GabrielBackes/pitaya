package sidecar

import (
	pitaya "github.com/topfreegames/pitaya/v2/pkg"
	config2 "github.com/topfreegames/pitaya/v2/pkg/config"
	"github.com/topfreegames/pitaya/v2/pkg/protos"
)
type SidecarBuilder struct {
	 pitaya.Builder
	 server  protos.PitayaServer
}

func NewSidecarBuilder(isFrontend bool, serverType string, serverMode pitaya.ServerMode, serverMetadata map[string]string, builderConfig config2.BuilderConfig) *SidecarBuilder {
	builder := pitaya.NewDefaultBuilder(isFrontend, serverType, serverMode, serverMetadata, builderConfig)
	return &SidecarBuilder{Builder : *builder}
}

func(s *SidecarBuilder) SetPitayaServer(server protos.PitayaServer){
	s.server = server
}

func(s *SidecarBuilder) Build() pitaya.Pitaya{
	app := s.Builder.Build()
	s.RPCServer.SetPitayaServer(s.server)
	return app
}


