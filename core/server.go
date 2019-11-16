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

func Server(config ServerConfig) {
	listener := listen(config.Port)
	accessListener := listen(8001)

	for {
		conn := accept(listener)
		accessConn := accept(accessListener)
		if conn == nil || accessConn == nil {
			log.Println("Accept client failed, retry at", 5, "seconds")
			time.Sleep(5 * time.Second)
			continue
		}
		forward(conn, accessConn)
	}
}
