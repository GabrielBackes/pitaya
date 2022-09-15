package main

import (
	"flag"
	"fmt"
	"github.com/topfreegames/pitaya/v2/pkg/config"
	"github.com/topfreegames/pitaya/v2/sidecar"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	tmpDir := os.TempDir()
	debug := flag.Bool("debug", false, "turn debug on")
	bind := flag.String("bind", filepath.FromSlash(fmt.Sprintf("%s/pitaya.sock", strings.TrimSuffix(tmpDir, "/"))), "bind address of the sidecar")
	bindProtocol := flag.String( "bindProtocol", "unix", "bind address of the sidecar")

	flag.Parse()

	cfg := config.NewBuilderConfig(config.NewConfig())
	sidecar := sidecar.NewSidecar(*cfg, *debug)
	sidecar.StartSidecar(*bind, *bindProtocol)
}
