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
		conn, err := net.Dial("tcp", targetAddr.String())
		if err == nil {
			log.Printf("Dial to [%s] success.\n", targetAddr)
			return conn
		}

		redialTimes++
		if maxRedialTimes < 0 || redialTimes < maxRedialTimes {
			// 重连模式，每5秒一次
			log.Printf("Dial to [%s] failed, retry(%d) after %d seconeds.", targetAddr.String(), redialTimes, retryIntervalTime)
			time.Sleep(retryIntervalTime * time.Second)
		} else {
			log.Printf("Dial to [%s] failed. %s\n", targetAddr.String(), err.Error())
			return nil
		}
	}
}

// 请求连接
func _requestConn(serverConn net.Conn, key string, accessPort int) (Protocol, bool) {
	reqProtocol := Protocol{
		Port: accessPort,
		Key:  key,
	}
	if !sendProtocol(serverConn, reqProtocol) {
		return Protocol{Result: protocolResultFailToSend}, false
	}
	return receiveProtocol(serverConn)
}

// 处理客户端连接
func _handleClientConn(cfg config.ClientConfig, index int) {
	server := cfg.ServerAddr
	local := cfg.LocalAddr[index]
	accessPort := cfg.AccessPort[index]

	var conn, serverConn net.Conn
	for {
		// 代理服务拨号，失败则关闭客户端
		if serverConn = _dial(server, cfg.MaxRedialTimes); serverConn == nil {
			continue
		}
		// 发送连接请求
		resp, ok := _requestConn(serverConn, cfg.Key, accessPort)
		if !ok || resp.Result != protocolResultSuccess {
			log.Println("Fail to request conn.", resp.String())
			closeConn(serverConn)
			break
		}
		log.Printf("Proxy address [%s] --> [%s:%d]\n", local, server.IP, resp.Port)
		// 如果内网服务不通，尝试重连后放弃
		if conn = _dial(local, cfg.MaxRedialTimes); conn == nil {
			closeConn(serverConn)
			continue
		}
		forward(conn, serverConn)
	}
}

// 入口
func Client(cfg config.ClientConfig) {
	log.Println("Load config", cfg)

	// 遍历所有端口
	for index := range cfg.LocalAddr {
		go _handleClientConn(cfg, index)
	}
	select {}
}
