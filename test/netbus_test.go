package test

import (
	"../core"
	"fmt"
	"log"
	"testing"
	"time"
)

// -----------------------------------------------------------------------------

// netbus
func TestServer(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	config := core.ServerConfig{
		Port:       6666,
		RandomPort: true,
	}
	core.Server(config)
}

func TestClient(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	config := core.ClientConfig{
		ServerAddr: "127.0.0.1:6666",
		LocalAddr:  "127.0.0.1:7456",
	}
	core.Client(config)
}

// -----------------------------------------------------------------------------

// server
func TestPort2Port(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	core.Port2Port("6666", "8456")
}

// client
func TestHost2Host(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	core.Host2Host("127.0.0.1:7456", "127.0.0.1:6666")
}

func TestPort2Host(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	core.Port2Host("6666", "127.0.0.1:7456")
}

func TestStruct2Byte(t *testing.T) {
	header := &core.Header{
		Address:   core.NetAddress{IP: "127.0.0.1", Port: 6666},
		Timestamp: time.Now().Unix(),
	}
	bytes := core.Serialize(header)
	fmt.Println("ptestStruct.data is : ", *core.Deserialize(&bytes))
}
