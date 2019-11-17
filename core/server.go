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

func handleServerConn(conn net.Conn) {
	// 建立连接后，接收发来的被代理者信息
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Fail to read local addresses", err.Error())
		return
	}
	address := string(buffer[:n])
	_, port := ParseHost(address)
	log.Println("proxy address", address)

	proxy := listen(port + 1000)
	if proxy == nil {
		return
	}
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
	if listener == nil {
		os.Exit(1)
	}
	for {
		conn := accept(listener)
		if conn != nil {
			go handleServerConn(conn)
		}
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
