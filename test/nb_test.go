package test

import (
	"../nb"
	"testing"
)

// -----------------------------------------------------------------------------

// server
func TestPort2Port(t *testing.T) {
	nb.Port2Port("6666", "8456")
}

// client
func TestHost2Host(t *testing.T) {
	nb.Host2Host("127.0.0.1:7456", "127.0.0.1:6666")
}

func TestPort2Host(t *testing.T) {
	nb.Port2Host("6666", "127.0.0.1:7456")
}
