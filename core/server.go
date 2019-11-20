package core

import (
	"../config"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

// 监听端口
func _listen(port int) net.Listener {
	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Listen failed, the port may be used or closed", port)
		return nil
	}
	log.Println("Listening at address", address)
	return listener
}

// 受理请求
func _accept(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		log.Println("Accept connect failed ->", conn.RemoteAddr(), err.Error())
		return nil
	}
	log.Println("Accept a new client ->", conn.RemoteAddr())
	return conn
}

// 所有访问监听器
// key:   NetAddress.String()
// value: net.Listener
var listeners sync.Map

// 根据地址获取已建立的监听，如果不存在则创建监听
func _loadOrStore(originalAddr config.NetAddress, portMode int) net.Listener {
	listener, exists := listeners.Load(originalAddr.String())
	if !exists {
		proxyPort := config.NewProxyPort(portMode, originalAddr.Port)
		listener = _listen(proxyPort)
		listeners.Store(originalAddr.String(), listener)
	}
	return listener.(net.Listener)
}

func _buildAccessAddress(conn net.Conn, listener net.Listener) config.NetAddress {
	serverAddress := conn.RemoteAddr().String()
	ipIndex := strings.LastIndex(serverAddress, ":")

	listenerAddress := listener.Addr().String()
	portIndex := strings.LastIndex(listenerAddress, ":")
	port, _ := strconv.Atoi(listenerAddress[portIndex+1:])

	return config.NetAddress{IP: serverAddress[:ipIndex], Port: port}
}

// 处理连接
func _handleServerConn(conn net.Conn, cfg config.ServerConfig) {
	// 接收消息头，并检查端口是否可代理
	address, ok := receiveHeader(conn)
	if !ok || !config.CheckProxyPort(cfg.PortMode, address.Port) {
		closeConn(conn)
		return
	}
	// 取出监听
	// TODO FIX 存在监听无法被 close 的问题
	listener := _loadOrStore(address, cfg.PortMode)
	// 回写访问地址
	header := _buildAccessAddress(conn, listener)
	if !sendHeader(conn, header) {
		closeConn(conn)
		return
	}

	// 代理连接
	proxyConn := _accept(listener)
	if conn != nil && proxyConn != nil {
		forward(conn, proxyConn)
	} else {
		closeConn(conn)
		closeConn(proxyConn)
	}
}

func Server(cfg config.ServerConfig) {
	serverListener := _listen(cfg.Port)
	if serverListener == nil {
		os.Exit(1)
	}
	// TODO 需要记录已使用的端口
	for {
		conn := _accept(serverListener)
		if conn != nil {
			go _handleServerConn(conn, cfg)
		}
	}
}
