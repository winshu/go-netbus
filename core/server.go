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
func _listen(port uint32) net.Listener {
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
var (
	listeners  sync.Map
	listenerMu sync.Mutex
)

// 获取监听
func _fetchListener(accessPort uint32) net.Listener {
	listener, exists := listeners.Load(accessPort)
	if exists {
		return listener.(net.Listener)
	}

	// 若不存在，则创建监听
	listenerMu.Lock()
	defer listenerMu.Unlock()
	// 双重检查
	listener, exists = listeners.Load(accessPort)
	if !exists {
		listener = _listen(accessPort)
		if listener != nil {
			listeners.Store(accessPort, listener)
		}
	}
	return listener.(net.Listener)
}

// 处理连接
func _handleBridgeConn(bridgeConn net.Conn, cfg config.ServerConfig) {
	// 接收协议消息
	req := receiveProtocol(bridgeConn)
	if !req.Success() {
		log.Println("Fail to receive protocol", req.String())
		sendProtocol(bridgeConn, req.NewResult(protocolResultFailToReceive))
		closeConn(bridgeConn)
		return
	}
	// 检查版本号
	if req.Version != protocolVersion {
		log.Println("Version mismatch", req.String())
		sendProtocol(bridgeConn, req.NewResult(protocolResultVersionMismatch))
		closeConn(bridgeConn)
		return
	}

	// 检查权限
	if _, ok := config.CheckKey(cfg.Key, req.Key); !ok {
		log.Println("Unauthorized access", req.String())
		sendProtocol(bridgeConn, req.NewResult(protocolResultFailToAuth))
		closeConn(bridgeConn)
		return
	}

	// 建立连接
	serverListener := _fetchListener(req.AccessPort)
	if serverListener == nil {
		log.Println("Fail to fetch server listener", req.String())
		sendProtocol(bridgeConn, req.NewResult(protocolResultFailToListen))
		closeConn(bridgeConn)
		return
	}

	serverConn := _accept(serverListener)
	if serverConn == nil {
		closeConn(bridgeConn, serverConn)
		return
	}
	log.Println("Tunnel connected ->", req.String())

	// 通知客户端，开始通讯
	if sendProtocol(bridgeConn, req.NewResult(protocolResultSuccess)) {
		forward(bridgeConn, serverConn)
	} else {
		log.Println("Tunnel interrupted")
		closeConn(bridgeConn, serverConn)
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
			//log.Println("New bridge ->", bridgeConn.RemoteAddr())
			go _handleBridgeConn(bridgeConn, cfg)
		}
	}
}
