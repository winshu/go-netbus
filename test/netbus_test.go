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
		Key:  "winshu",
		Port: 6666,
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
		AccessPort:     []uint32{13306},
		MaxRedialTimes: 10,
	}
	core.Client(cfg)
}

func TestProtocol(t *testing.T) {
	seed := "winshu"
	key, _ := config.NewKey(seed, "2019-12-31")
	fmt.Println(key)
	fmt.Println(config.CheckKey(seed, key))
}
