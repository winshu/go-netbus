package core

import (
	"fmt"
	"log"
	"net"
	"time"
)

func listen(port int) net.Listener {
	host := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalln("Listen failed", host)
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

func handleServerConn(port int, conn net.Conn) {
	proxy := listen(port + 1000)
	for {
		proxyConn := accept(proxy)
		if conn == nil || proxyConn == nil {
			log.Println("Accept client failed, retry at", 5, "seconds")
			time.Sleep(5 * time.Second)
			continue
		}
		forward(conn, proxyConn)
	}
}

func Server(config ServerConfig) {
	listener := listen(config.Port)

	for {
		conn := accept(listener)

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Fail to read local addresses", err.Error())
			continue
		}
		address := string(buffer[:n])
		_, port := ParseHost(address)
		log.Println("proxy address", address)

		go handleServerConn(port, conn)
	}
}
