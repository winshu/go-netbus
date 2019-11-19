package core

import (
	"../config"
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

func Server(cfg config.ServerConfig) {
	serverListener := listen(cfg.Port)
	if serverListener == nil {
		os.Exit(1)
	}
	listeners := make(map[string]net.Listener, 10)
	for {
		conn := accept(serverListener)
		if conn == nil {
			continue
		}

		buffer := make([]byte, headerLengthInByte)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Println("Received header failed", err.Error())
			_ = conn.Close()
			continue
		}

		address := config.ParseNetAddress(string(buffer))
		log.Println("Received header", address)

		listener, exists := listeners[address.String()]
		if !exists {
			listener = listen(address.Port + 1000)
			listeners[address.String()] = listener
		}

		proxyConn := accept(listener)
		if conn == nil || proxyConn == nil {
			log.Println("Accept client failed, retry after", retryIntervalTime, "seconds")
			time.Sleep(retryIntervalTime * time.Second)
			continue
		}
		forward(conn, proxyConn)
	}
}
