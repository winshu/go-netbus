package core

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func listen(port int) net.Listener {
	host := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Println("Listen failed, the port may be used or closed", host)
		return nil
	}
	log.Println("Listening at address", host)
	return listener
}

func accept(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		log.Println("Accept connect failed ->", conn.RemoteAddr(), err.Error())
		return nil
	}
	log.Println("Accept a new client ->", conn.RemoteAddr())
	return conn
}

func Server(config ServerConfig) {
	serverListener := listen(config.Port)
	if serverListener == nil {
		os.Exit(1)
	}
	listeners := make(map[string]net.Listener, 10)
	for {
		conn := accept(serverListener)
		if conn == nil {
			continue
		}
		address := readHeader(conn)
		if address.invalid {
			log.Println("Read original address failed")
			continue
		}

		listener, ok := listeners[address.String()]
		if !ok {
			listener = listen(address.Port + 1000)
			listeners[address.String()] = listener
		}

		proxyConn := accept(listener)
		if conn == nil || proxyConn == nil {
			log.Println("Accept client failed, retry at", 5, "seconds")
			time.Sleep(5 * time.Second)
			continue
		}
		forward(conn, proxyConn)
	}
}

func SingleServer(config ServerConfig) {
	listener1 := listen(config.Port)
	listener2 := listen(8456)

	for {
		conn1 := accept(listener1)
		buffer := make([]byte, 1024)
		n, err := conn1.Read(buffer)
		if err != nil {
			log.Println("Fail to read local addresses", err.Error())
			return
		}
		address := string(buffer[:n])
		log.Println("proxy address", address)

		conn2 := accept(listener2)
		if conn1 == nil || conn2 == nil {
			log.Println("Accept client failed, retry at", 5, "seconds")
			time.Sleep(5 * time.Second)
			continue
		}
		forward(conn1, conn2)
	}
}
