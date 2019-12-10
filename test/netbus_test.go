package test

import (
	"../config"
	"../core"
	"encoding/binary"
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
		AccessPort:     []uint32{13306},
		MaxRedialTimes: 10,
	}
	core.Client(cfg)
}

type pro struct {
	r byte
	p uint32
	a uint32
	k string
}

func TestProtocol(t *testing.T) {
	buffer := make([]byte, 8)

	binary.BigEndian.PutUint32(buffer[1:8], 22)
	fmt.Println(buffer)
}
