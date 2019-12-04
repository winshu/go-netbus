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
		Port: 6666,
		Key:  "winshu",
	}
	core.Server(cfg)
}

func TestClient(t *testing.T) {
	cfg := config.ClientConfig{
		Key: "winshu",
		ServerAddr: config.NetAddress{
			IP: "127.0.0.1", Port: 6666,
		},
		LocalAddr: []config.NetAddress{
			{"127.0.0.1", 3306},
		},
		AccessPort:     []int{13306},
		MaxRedialTimes: 10,
	}
	core.Client(cfg)
}

func TestHeader(t *testing.T) {
	arr := []string{"aa", "bb", "cc"}
	for i := range arr {
		fmt.Println(i)
	}
}
