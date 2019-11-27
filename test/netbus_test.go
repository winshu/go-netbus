package test

import (
	"../config"
	"../core"
	"../util"
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
		Port:          6666,
		CustomPortKey: "custom",
		RandomPortKey: "random",
	}
	core.Server(cfg)
}

func TestClient(t *testing.T) {
	cfg := config.ClientConfig{
		Key: "custom",
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
	fmt.Println(util.RandToken("abcde", 10))
}
