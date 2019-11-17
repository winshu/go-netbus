package test

import (
	"../core"
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	config := core.ClientConfig{
		ServerAddr: "127.0.0.1:6666",
		LocalAddr:  "127.0.0.1:7456",
	}
	core.Client(config)
}
