package test

import (
	"../core"
	"log"
	"testing"
)

func TestServer(t *testing.T) {
	//config := core.ServerConfig{
	//	Port:       6666,
	//	RandomPort: true,
	//}
	//core.SingleServer(config)
	core.Port2Port("6666", "8456")
}

func TestClient(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//config := core.ClientConfig{
	//	ServerAddr: "127.0.0.1:6666",
	//	LocalAddr:  "127.0.0.1:7456",
	//}
	//core.Client(config)
	core.Host2Host("127.0.0.1:7456", "127.0.0.1:6666")
}
