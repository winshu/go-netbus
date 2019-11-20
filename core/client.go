package core

import (
	"../config"
	"log"
	"net"
	"time"
)

// 拨号
func _dial(targetAddr config.NetAddress /*目标地址*/, maxRedialTimes int /*最大重拨次数*/) net.Conn {
	redialTimes := 0
	for {
		log.Println("Dial to", targetAddr)
		conn, err := net.Dial("tcp", targetAddr.String())
		if err == nil {
			log.Println("Dial success ->", targetAddr)
			return conn
		}

		redialTimes++
		if maxRedialTimes < 0 || redialTimes < maxRedialTimes {
			// 重连模式，每5秒一次
			log.Printf("Dial failed, retry(%d) after %dnd", redialTimes, retryIntervalTime)
			time.Sleep(retryIntervalTime * time.Second)
		} else {
			log.Println("Dial failed ->", err.Error())
			return nil
		}
	}
}

func _requestHeader(serverConn net.Conn, localAddr config.NetAddress) (config.NetAddress, bool) {
	if !sendHeader(serverConn, localAddr) {
		return config.NetAddress{}, false
	}
	header, ok := receiveHeader(serverConn)
	if !ok {
		log.Println("Send header error")
		return config.NetAddress{}, false
	}
	return header, true
}

// 处理客户端连接
func _handleClientConn(localAddr, serverAddr config.NetAddress, maxRedialTimes int) {
	for {
		// 本地服务拨号
		conn := _dial(localAddr, maxRedialTimes)
		if conn == nil {
			return
		}
		// 代理服务拨号
		serverConn := _dial(serverAddr, maxRedialTimes)
		if serverConn == nil {
			return
		}
		// 请求头
		if header, ok := _requestHeader(serverConn, localAddr); ok {
			log.Println("Access address", header)
			forward(conn, serverConn)
		}
	}
}

func Client(cfg config.ClientConfig) {
	for _, addr := range cfg.LocalAddr {
		go _handleClientConn(addr, cfg.ServerAddr, cfg.MaxRedialTimes)
	}
	select {}
}
