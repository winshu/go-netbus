package test

import (
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
	config := core.ServerConfig{
		Port:       6666,
		RandomPort: true,
	}
	core.Server(config)
}

func TestClient(t *testing.T) {
	config := core.ClientConfig{
		ServerAddr:     "127.0.0.1:6666",
		LocalAddr:      "127.0.0.1:7456",
		MaxRedialTimes: 10,
	}
	core.Client(config)
}
