package core

import (
	"../config"
	"../util"
	"fmt"
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
			//log.Println("Dial success ->", targetAddr)
			return conn
		}

		redialTimes++
		if maxRedialTimes < 0 || redialTimes < maxRedialTimes {
			// 重连模式，每5秒一次
			log.Printf("Dial failed, retry(%d) after %d seconeds.", redialTimes, retryIntervalTime)
			time.Sleep(retryIntervalTime * time.Second)
		} else {
			log.Println("Dial failed ->", err.Error())
			return nil
		}
	}
}

// 处理客户端连接
func _handleClientConn(token string, local config.NetAddress, server config.NetAddress, accessPort int, maxRedialTimes int) {
	var conn, serverConn net.Conn
	for {
		// 本地服务拨号
		if conn = _dial(local, maxRedialTimes); conn == nil {
			return
		}
		// 代理服务拨号
		if serverConn = _dial(server, maxRedialTimes); serverConn == nil {
			return
		}

		tokenX := fmt.Sprintf("%s%d", token, accessPort)
		reqProtocol := Protocol{Token: tokenX, Ports: []int{accessPort}}
		if sendProtocol(serverConn, reqProtocol) {
			log.Printf("Proxy address [%s] --> [%s:%d]\n", local, server.IP, accessPort)
			forward(conn, serverConn)
		}
	}
}

// 鉴权
func _requestAuth(token string, cfg config.ClientConfig) Protocol {
	serverConn := _dial(cfg.ServerAddr, cfg.MaxRedialTimes)
	defer closeConn(serverConn)
	if serverConn == nil {
		return Protocol{}
	}

	// 验证身份
	// 如果没有配置固定端口
	ports := cfg.AccessPort
	if len(ports) == 0 {
		ports = config.ExtractPorts(cfg.LocalAddr)
	}

	header := Protocol{
		Type:  protocolTypeAuth, // 鉴权
		Ports: ports,
		Token: token,
	}
	if !sendProtocol(serverConn, header) {
		return Protocol{}
	}
	return receiveProtocol(serverConn)
}

// 入口
func Client(cfg config.ClientConfig) {
	// token 随机生成
	token := util.RandToken(cfg.Key, protocolTokenLength)
	var protocol Protocol

	//鉴权
	if protocol = _requestAuth(token, cfg); protocol.Result != protocolResultSuccess {
		log.Fatalln("Fail to auth")
	}
	log.Println("Auth success", protocol)

	// 连接
	for i, local := range cfg.LocalAddr {
		accessPort := protocol.Ports[i]
		go _handleClientConn(token, local, cfg.ServerAddr, accessPort, cfg.MaxRedialTimes)
	}
	select {}
}
