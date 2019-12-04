package core

import (
	"../config"
	"../util"
	"fmt"
	"log"
	"net"
	"os"
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
		log.Println("Accept connect failed ->", err.Error())
		return nil
	}
	//log.Println("Accept a new client ->", conn.RemoteAddr())
	return conn
}

// 所有访问监听器
// key:   token
// value: net.Listener
var listeners sync.Map

// 从 listeners 中加载监听
func _loadListener(protocol Protocol) net.Listener {
	listener, exists := listeners.Load(protocol.Ports[0])
	if !exists {
		return nil
	}
	// listener 为空时，不能强转，go 语法
	return listener.(net.Listener)
}

// 创建监听
func _initAccessListener(conn net.Conn, protocol Protocol, cfg config.ServerConfig) bool {
	// 创建访问端口
	accessPort, ok := _makeAccessPort(protocol, cfg)
	if !ok { // 失败，发送失败结果到客户端
		sendProtocol(conn, Protocol{
			Result: protocolResultFail,
			Type:   protocolTypeAuth,
			Ports:  protocol.Ports,
			Token:  protocol.Token,
		})
		return false
	}

	// 创建监听失败
	var listener net.Listener
	for _, port := range accessPort {
		// 如果端口已经在监听，则重复使用
		if _, exists := listeners.Load(port); exists {
			log.Printf("Port %d is already listening\n", port)
			continue
		}

		listener = _listen(port)
		if listener == nil {
			// 监听失败
			sendProtocol(conn, Protocol{
				Result: protocolResultFail,
				Type:   protocolTypeAuth,
				Ports:  accessPort,
				Token:  protocol.Token,
			})
			return false
		}
		// 存放监听映射信息
		listeners.Store(port, listener)
	}

	// 发送鉴权结果到客户端
	return sendProtocol(conn, Protocol{
		Result: protocolResultSuccess,
		Type:   protocolTypeAuth,
		Ports:  accessPort,
		Token:  protocol.Token,
	})
}

// 创建访问端口
func _makeAccessPort(protocol Protocol, cfg config.ServerConfig) ([]int, bool) {
	if len(protocol.Ports) == 0 {
		return nil, false
	}
	if strings.HasPrefix(protocol.Token, cfg.CustomPortKey) {
		// 允许自定义端口、同名端口
		return protocol.Ports, true
	} else if strings.HasPrefix(protocol.Token, cfg.RandomPortKey) {
		// 允许随机端口
		randPorts := util.RandPorts(len(protocol.Ports))
		return randPorts, true
	}
	// 无权限访问
	return nil, false
}

// 处理连接
func _handleServerConn(bridgeConn net.Conn, cfg config.ServerConfig) {
	var serverListener net.Listener
	protocol, ok := receiveProtocol(bridgeConn)
	if !ok {
		return
	}
	log.Println("----------> Receive protocol", protocol)
	switch protocol.Type {
	case protocolTypeConn:
		if serverListener = _loadListener(protocol); serverListener == nil {
			// 获取监听失败
			return
		}
	case protocolTypeAuth:
		if !_initAccessListener(bridgeConn, protocol, cfg) {
			log.Println("Fail to auth", protocol)
		}
		// 鉴权成功需要返回
		closeConn(bridgeConn)
		return
	default:
		// 非法类型
		log.Println("Forbidden protocol type", protocol)
		closeConn(bridgeConn)
		return
	}

	serverConn := _accept(serverListener)
	if serverConn == nil {
		closeConn(bridgeConn)
		closeConn(serverConn)
		return
	}
	log.Println("Accept a new server ->", serverConn.RemoteAddr(), serverConn.LocalAddr())
	// 通知客户端，开始通讯
	if sendProtocol(bridgeConn, Protocol{
		Result: protocolResultSuccess,
		Type:   protocolTypeConn,
		Ports:  protocol.Ports,
		Token:  protocol.Token,
	}) {
		forward(bridgeConn, serverConn)
	}
}

// 入口
func Server(cfg config.ServerConfig) {
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
