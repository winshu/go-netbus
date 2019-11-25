package test

import (
	"../config"
	"../core"
	"fmt"
	"log"
	"testing"
)

// -----------------------------------------------------------------------------

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
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
		ServerAddr: config.NetAddress{IP: "10.3.8.119", Port: 6666},
		LocalAddr: []config.NetAddress{
			{"127.0.0.1", 7456},
		},
		MaxRedialTimes: 10,
	}
	core.Client(cfg)
}

func TestHeader(t *testing.T) {
	for i, v := range []int{3, 4, 5, 6} {
		fmt.Println(i, v)
	}
}
