package core

import (
	"../config"
	"fmt"
	"log"
	"net"
	"os"
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
		log.Println("Accept connect failed ->", err.Error())
		return nil
	}
	//log.Println("Accept a new client ->", conn.RemoteAddr())
	return conn
}

// 所有访问监听器
// key:   accessPort
// value: net.Listener
var listeners sync.Map

// 创建访问端口
func _auth(req Protocol, cfg config.ServerConfig) (port int, ok bool) {
	if len(req.Key) < protocolKeyMinLength || len(req.Key) > protocolKeyMaxLength {
		return
	}
	if req.Key == cfg.Key {
		return req.Port, true
	}
	return
}

func _fetchListener(accessPort int, cfg config.ServerConfig) net.Listener {
	// 获取监听
	listener, exists := listeners.Load(accessPort)
	if exists {
		log.Printf("Port [%d] is already listening\n", accessPort)
		return listener.(net.Listener)
	}

	// 若不存在，则创建监听
	listener = _listen(accessPort)
	if listener == nil {
		return nil
	}
	listeners.Store(accessPort, listener)
	return listener.(net.Listener)
}

// 处理连接
func _handleServerConn(bridgeConn net.Conn, cfg config.ServerConfig) {
	// 接收协议消息
	req, ok := receiveProtocol(bridgeConn)
	if !ok {
		sendProtocol(bridgeConn, req.NewResult(protocolResultFailToReceive))
		closeConn(bridgeConn)
		return
	}
	// 检查权限
	accessPort, ok := _auth(req, cfg)
	if !ok {
		log.Println("Unauthorized access", req.String())
		sendProtocol(bridgeConn, req.NewResult(protocolResultFailToAuth))
		closeConn(bridgeConn)
		return
	}

	// 建立连接
	serverListener := _fetchListener(accessPort, cfg)
	if serverListener == nil {
		log.Println("Fail to fail server listener", req.String())
		sendProtocol(bridgeConn, req.NewResult(protocolResultFailToListen))
		closeConn(bridgeConn)
		return
	}

	serverConn := _accept(serverListener)
	if serverConn == nil {
		closeConn(bridgeConn)
		closeConn(serverConn)
		return
	}
	log.Println("Accept a new server ->", req.String())
	// 通知客户端，开始通讯
	if sendProtocol(bridgeConn, Protocol{Port: accessPort, Key: req.Key}) {
		forward(bridgeConn, serverConn)
	}
}

// 入口
func Server(cfg config.ServerConfig) {
	log.Println("Load config", cfg)

	// 监听桥接端口
	bridgeListener := _listen(cfg.Port)
	if bridgeListener == nil {
		os.Exit(1)
	}

	for {
		// 受理来自客户端的请求
		bridgeConn := _accept(bridgeListener)
		if bridgeConn != nil {
			log.Println("Accept a new client ->", bridgeConn.RemoteAddr(), bridgeConn.LocalAddr())
			go _handleServerConn(bridgeConn, cfg)
		}
	}
}
