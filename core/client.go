package core

import (
	"../config"
	"log"
	"net"
	"runtime"
	"time"
)

// 拨号
func _dial(targetAddr config.NetAddress /*目标地址*/, maxRedialTimes int /*最大重拨次数*/) net.Conn {
	redialTimes := 0
	for {
		conn, err := net.Dial("tcp", targetAddr.String())
		if err == nil {
			//log.Printf("Dial to [%s] success.\n", targetAddr)
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
func _requestConn(serverConn net.Conn, key string, port int, accessPort int) (Protocol, bool) {
	reqProtocol := Protocol{
		AccessPort: accessPort,
		Port:       port,
		Key:        key,
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

	connChan := make(chan net.Conn)
	flagChan := make(chan bool)

	// 拨号
	go func(connCh chan net.Conn, flagCh chan bool) {
		for {
			select {
			case <-flagCh:
				go func(ch chan net.Conn) {
					conn := _dial(server, cfg.MaxRedialTimes)
					if conn == nil {
						runtime.Goexit()
					}
					log.Printf("Proxy service [%s] -> [%s:%d]\n", local.String(), server.IP, accessPort)
					resp, ok := _requestConn(conn, cfg.Key, local.Port, accessPort)
					if !ok || resp.Result != protocolResultSuccess {
						log.Println("Fail to request conn.", resp.String())
						closeConn(conn)
						return
					}
					ch <- conn
				}(connCh)
			default:
				// default
			}
		}
	}(connChan, flagChan)

	// 连接
	go func(connCh chan net.Conn, flagCh chan bool) {
		for {
			select {
			case cn := <-connCh:
				go func(conn net.Conn) {
					localConn := _dial(local, cfg.MaxRedialTimes)
					if localConn == nil {
						closeConn(conn)
						flagCh <- true // 通知创建连接
						return
					}
					flagCh <- true // 通知创建连接
					forward(localConn, conn)
				}(cn)
			default:
				// default
			}
		}
	}(connChan, flagChan)

	// 初始化连接
	flagChan <- true
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
