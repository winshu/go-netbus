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

// 鉴权
func _requestAuth(token string, cfg config.ClientConfig) (resp Protocol, ok bool) {
	serverConn := _dial(cfg.ServerAddr, cfg.MaxRedialTimes)
	defer closeConn(serverConn)
	if serverConn == nil {
		return
	}

	// 验证身份
	// 如果没有配置固定端口
	ports := cfg.AccessPort
	if len(ports) == 0 {
		ports = config.ExtractPorts(cfg.LocalAddr)
	}

	req := Protocol{
		Result: protocolResultSuccess,
		Type:   protocolTypeAuth, // 鉴权
		Ports:  ports,
		Token:  token,
	}
	if !sendProtocol(serverConn, req) {
		return
	}
	return receiveProtocol(serverConn)
}

// 请求连接
func _requestConn(serverConn net.Conn, token string, accessPort int) bool {
	tokenX := fmt.Sprintf("%s%d", token, accessPort)
	reqProtocol := Protocol{
		Result: protocolResultSuccess,
		Type:   protocolTypeConn,
		Ports:  []int{accessPort},
		Token:  tokenX,
	}
	if !sendProtocol(serverConn, reqProtocol) {
		return false
	}
	_, ok := receiveProtocol(serverConn)
	return ok
}

// 处理客户端连接
func _handleClientConn(token string, local config.NetAddress, server config.NetAddress, accessPort int, maxRedialTimes int) {
	var conn, serverConn net.Conn
	for {
		// 代理服务拨号，失败则关闭客户端
		if serverConn = _dial(server, maxRedialTimes); serverConn == nil {
			break
		}
		// 发送建立连接请求
		if !_requestConn(serverConn, token, accessPort) {
			continue
		}
		// 接收到请求，则拨号连接内网服务
		// 如果内网服务不通，尝试重连后放弃
		if conn = _dial(local, maxRedialTimes); conn == nil {
			closeConn(serverConn)
			break
		}
		log.Printf("Proxy address [%s] --> [%s:%d]\n", local, server.IP, accessPort)
		forward(conn, serverConn)
	}
}

// 入口
func Client(cfg config.ClientConfig) {
	// token 随机生成
	token := util.RandToken(cfg.Key, protocolTokenLength)

	//鉴权
	protocol, ok := _requestAuth(token, cfg)
	if !ok {
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
