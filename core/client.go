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
			log.Printf("Dial failed, retry(%d) after %d seconeds.", redialTimes, retryIntervalTime)
			time.Sleep(retryIntervalTime * time.Second)
		} else {
			log.Println("Dial failed ->", err.Error())
			return nil
		}
	}
}

func _requestHeader(serverConn net.Conn) (Header, bool) {
	reqHeader := Header{Token: token}
	if !sendHeader(serverConn, reqHeader) {
		return Header{}, false
	}
	respHeader, ok := receiveHeader(serverConn)
	if !ok {
		log.Println("Send header error")
		return Header{}, false
	}
	return respHeader, true
}

// 处理客户端连接
func _handleClientConn(index int, cfg config.ClientConfig) {
	for {
		// 本地服务拨号
		conn := _dial(cfg.LocalAddr[index], cfg.MaxRedialTimes)
		if conn == nil {
			return
		}
		// 代理服务拨号
		serverConn := _dial(cfg.ServerAddr, cfg.MaxRedialTimes)
		if serverConn == nil {
			return
		}
		// 请求头
		if _, ok := _requestHeader(serverConn); ok {
			forward(conn, serverConn)
		} else {
			// 关闭连接
			closeConn(conn)
			closeConn(serverConn)
		}
	}
}

func _auth(cfg config.ClientConfig) bool {
	serverConn := _dial(cfg.ServerAddr, cfg.MaxRedialTimes)
	if serverConn == nil {
		return false
	}

	// 验证身份
	header := Header{
		Type:  1,
		Ports: cfg.AccessPort,
		Token: token,
	}
	if !sendHeader(serverConn, header) {
		return false
	}
	_, ok := receiveHeader(serverConn)
	return ok
}

// 客户端令牌，每次启动时生成
var token string

func Client(cfg config.ClientConfig) {
	token = "123456"

	// 身份验证
	if !_auth(cfg) {
		log.Fatalln("Fail to auth")
	}

	for i, _ := range cfg.LocalAddr {
		go _handleClientConn(i, cfg)
	}
	select {}
}
