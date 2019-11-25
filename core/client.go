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
		log.Println("Dial to", targetAddr)
		conn, err := net.Dial("tcp", targetAddr.String())
		if err == nil {
			log.Println("Dial success ->", targetAddr)
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

func _requestHeader(serverConn net.Conn, token string, accessPort int) Header {
	tokenX := fmt.Sprintf("%s%d", token, accessPort)
	reqHeader := Header{Token: tokenX, Ports: []int{accessPort}}
	if !sendHeader(serverConn, reqHeader) {
		return Header{Result: 0}
	}
	return receiveHeader(serverConn)
}

// 处理客户端连接
func _handleClientConn(token string, local config.NetAddress, server config.NetAddress, accessPort int, maxRedialTimes int) {
	for {
		// 本地服务拨号
		conn := _dial(local, maxRedialTimes)
		if conn == nil {
			return
		}
		// 代理服务拨号
		serverConn := _dial(server, maxRedialTimes)
		if serverConn == nil {
			return
		}
		// 请求头
		if header := _requestHeader(serverConn, token, accessPort); header.Result == 1 {
			forward(conn, serverConn)
		} else {
			// 关闭连接
			closeConn(conn)
			closeConn(serverConn)
		}
	}
}

func _auth(token string, cfg config.ClientConfig) Header {
	serverConn := _dial(cfg.ServerAddr, cfg.MaxRedialTimes)
	if serverConn == nil {
		return Header{}
	}
	// 验证身份
	// 如果没有配置固定端口
	ports := cfg.AccessPort
	if len(ports) == 0 {
		ports = make([]int, len(cfg.LocalAddr))
		for i, addr := range cfg.LocalAddr {
			ports[i] = addr.Port
		}
	}

	header := Header{
		Type:  1, // 鉴权
		Ports: ports,
		Token: token,
	}
	if !sendHeader(serverConn, header) {
		return Header{}
	}
	return receiveHeader(serverConn)
}

func Client(cfg config.ClientConfig) {
	// token 随机生成
	token := util.RandToken(cfg.Key, 16)
	var header Header

	// 身份验证
	if header = _auth(token, cfg); header.Result != 1 {
		log.Fatalln("Fail to auth")
	}
	for i, local := range cfg.LocalAddr {
		accessPort := header.Ports[i]
		go _handleClientConn(token, local, cfg.ServerAddr, accessPort, cfg.MaxRedialTimes)
	}
	select {}
}
