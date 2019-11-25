package core

import (
	"../config"
	"fmt"
	"log"
	"math/rand"
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
// key:   NetAddress.String()
// value: net.Listener
var listeners sync.Map

// 根据地址获取已建立的监听，如果不存在则创建监听
func _loadListener(key string) net.Listener {
	listener, exists := listeners.Load(key)
	if exists {
		return listener.(net.Listener)
	}
	log.Println("Listener key not exists", key)
	return nil
}

// 创建访问端口
func _makeAccessPort(header Header, cfg config.ServerConfig) ([]int, bool) {
	if strings.HasPrefix(header.Token, cfg.CustomPortKey) {
		// 允许自定义端口、同名端口
		return header.Ports, true
	} else if strings.HasPrefix(header.Token, cfg.RandomPortKey) {
		// 允许随机端口
		randPorts := make([]int, len(header.Ports))
		for i := 0; i < len(header.Ports); i++ {
			randPorts[i] = 60000 + rand.Intn(5535)
		}
		return randPorts, true
	} else {
		// 无权限访问
		return nil, false
	}
}

// 处理连接
func _handleServerConn(conn net.Conn, cfg config.ServerConfig) {
	var listener net.Listener

	header := receiveHeader(conn)
	switch header.Type {
	case 0:
		// 正常通讯
		// 取监听
		listener = _loadListener(header.Token)
	case 1:
		// 鉴权
		// 创建监听
		accessPort, ok := _makeAccessPort(header, cfg)
		// 失败
		if !ok {
			sendHeader(conn, Header{
				Type:  1,
				Ports: header.Ports,
				Token: header.Token,
			})
			return
		}

		for _, port := range accessPort {
			key := fmt.Sprintf("%s%d", header.Token, port)
			listener = _listen(port)
			listeners.Store(key, listener)
		}

		// 返回消息
		sendHeader(conn, Header{
			Result: 1,
			Type:   1,
			Ports:  accessPort,
			Token:  header.Token,
		})
		// 正常关闭
		closeConn(conn)
		return
	default:
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

	for {
		conn := _accept(serverListener)
		if conn != nil {
			go _handleServerConn(conn, cfg)
		}
	}
}
