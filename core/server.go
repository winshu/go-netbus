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

// 接收消息头，包含了地址信息
func receiveHeader(conn net.Conn) (config.NetAddress, bool) {
	buffer := make([]byte, headerLengthInByte)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Println("Receive header failed", err.Error())
		_ = conn.Close()
		return config.NetAddress{}, false
	}
	address, _ := config.ParseNetAddress(string(buffer))
	log.Println("Receive header", address)
	return address, true
}

func Server(cfg config.ServerConfig) {
	serverListener := listen(cfg.Port)
	if serverListener == nil {
		os.Exit(1)
	}
	// TODO 需要记录已使用的端口
	listeners := make(map[string]net.Listener, 10)
	for {
		conn := accept(serverListener)
		if conn == nil {
			continue
		}

		// 接收消息头，并检查端口是否可代理
		address, ok := receiveHeader(conn)
		if !ok || !config.CheckProxyPort(cfg.PortMode, address.Port) {
			continue
		}

		// 取出监听
		listener, exists := listeners[address.String()]
		if !exists {
			proxyPort := config.NewProxyPort(cfg.PortMode, address.Port)
			listener = listen(proxyPort)
			listeners[address.String()] = listener
		}

		// 代理连接
		proxyConn := accept(listener)
		if conn == nil || proxyConn == nil {
			log.Println("Accept client failed, retry after", retryIntervalTime, "seconds")
			time.Sleep(retryIntervalTime * time.Second)
			continue
		}
		forward(conn, proxyConn)
	}
}
