package core

import (
	"../config"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// 监听端口
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

// 受理请求
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

// 处理连接
func handleServerConn(conn net.Conn, cfg config.ServerConfig) {
	// 接收消息头，并检查端口是否可代理
	address, ok := receiveHeader(conn)
	if !ok || !config.CheckProxyPort(cfg.PortMode, address.Port) {
		return
	}
	// 取出监听
	listener := loadOrStore(address, cfg.PortMode)
	// 代理连接
	proxyConn := accept(listener)
	if conn == nil || proxyConn == nil {
		log.Println("Accept client failed, retry after", retryIntervalTime, "seconds")
		time.Sleep(retryIntervalTime * time.Second)
		return
	}
	forward(conn, proxyConn)
}

// 所有访问监听器
// key:   NetAddress.String()
// value: net.Listener
var listeners sync.Map

// 根据地址获取已建立的监听，如果不存在则创建监听
func loadOrStore(address config.NetAddress, portMode int) net.Listener {
	listener, _ := listeners.LoadOrStore(address.String(), func() net.Listener {
		// 不存在，创建监听，并放入监听池
		proxyPort := config.NewProxyPort(portMode, address.Port)
		return listen(proxyPort)
	}())
	return listener.(net.Listener)
}

func Server(cfg config.ServerConfig) {
	serverListener := listen(cfg.Port)
	if serverListener == nil {
		os.Exit(1)
	}
	// TODO 需要记录已使用的端口
	for {
		conn := accept(serverListener)
		if conn != nil {
			go handleServerConn(conn, cfg)
		}
	}
}
