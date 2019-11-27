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
		log.Println("Accept connect failed ->", conn.RemoteAddr(), err.Error())
		return nil
	}
	log.Println("Accept a new client ->", conn.RemoteAddr())
	return conn
}

// 所有访问监听器
// key:   token
// value: net.Listener
var listeners sync.Map

// 从 listeners 中加载监听
func _loadListener(conn net.Conn, protocol Protocol) net.Listener {
	listener, exists := listeners.Load(protocol.Ports[0])
	result := protocolResultFail
	if exists {
		result = protocolResultSuccess
	}
	// 发送处理结果
	sendProtocol(conn, Protocol{
		Result: result,
		Type:   protocolTypeAuth,
		Ports:  protocol.Ports,
		Token:  protocol.Token,
	})
	if listener == nil {
		return nil
	}
	// listener 为空时，不能强转，go 语法
	return listener.(net.Listener)
}

// 创建监听
func _buildListener(conn net.Conn, protocol Protocol, cfg config.ServerConfig) bool {
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
		// 如果端口已经在监听，则重复利用
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
func _handleServerConn(conn net.Conn, cfg config.ServerConfig) {
	var listener net.Listener
	protocol := receiveProtocol(conn)

	switch protocol.Type {
	case protocolTypeNormal:
		if listener = _loadListener(conn, protocol); listener == nil {
			return
		}
	case protocolTypeAuth:
		_buildListener(conn, protocol, cfg)
		// 注意不能关闭连接
		return
	default:
		// 非法类型
		log.Println("Forbidden type", protocol)
		closeConn(conn)
		return
	}
	// 代理连接
	proxyConn := _accept(listener)
	if conn != nil && proxyConn != nil {
		forward(conn, proxyConn)
	}
}

// 入口
func Server(cfg config.ServerConfig) {
	serverListener := _listen(cfg.Port)
	if serverListener == nil {
		os.Exit(1)
	}

	for {
		conn := _accept(serverListener)
		if conn != nil {
			go _handleServerConn(conn, cfg)
		}
	}
}
