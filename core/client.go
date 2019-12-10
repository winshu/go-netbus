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
			log.Printf("Dial to [%s] failed, redial(%d) after %d seconeds.", targetAddr.String(), redialTimes, retryIntervalTime)
			time.Sleep(retryIntervalTime * time.Second)
		} else {
			log.Printf("Dial to [%s] failed. %s\n", targetAddr.String(), err.Error())
			return nil
		}
	}
}

// 请求连接
func _requestConn(serverConn net.Conn, key string, port uint32, accessPort uint32) Protocol {
	reqProtocol := Protocol{
		Result:     protocolResultSuccess,
		AccessPort: accessPort,
		Port:       port,
		Key:        key,
	}
	if !sendProtocol(serverConn, reqProtocol) {
		return reqProtocol.NewResult(protocolResultFailToSend)
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
					resp := _requestConn(conn, cfg.Key, local.Port, accessPort)
					if resp.Result == protocolResultSuccess {
						ch <- conn
						return
					}
					// 连接中断，重新连接
					log.Printf("bridge connection interrupted, try to redial. [%d] [%s]\n", resp.Result, local.String())
					closeConn(conn)
					flagCh <- true
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
					// 本地连接，不需要重新拨号
					if localConn := _dial(local, 0); localConn != nil {
						// 通知创建新桥
						flagCh <- true
						forward(localConn, conn)
					} else {
						// 放弃连接，重新建桥
						closeConn(conn)
						flagCh <- true
					}
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
