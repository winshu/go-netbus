package test

import (
	"../config"
	"../core"
	"log"
	"testing"
)

// -----------------------------------------------------------------------------

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

// netbus
func TestServer(t *testing.T) {
	cfg := config.ServerConfig{
		Port:     6666,
		PortMode: 2,
	}
	core.Server(cfg)
}

func TestClient(t *testing.T) {
	cfg := config.ClientConfig{
		ServerAddr:     "127.0.0.1:6666",
		LocalAddr:      "127.0.0.1:7456",
		MaxRedialTimes: 10,
	}
	core.Client(cfg)
}
