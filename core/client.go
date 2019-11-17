package core

import (
	"log"
	"net"
	"time"
)

func dial(targetAddr string, maxRedialTimes int) net.Conn {
	redialTimes := 0
	for {
		log.Println("Dial to", targetAddr)
		conn, err := net.Dial("tcp", targetAddr)
		if err == nil {
			log.Println("Dial success ->", targetAddr)
			return conn
		}
		switch {
		case maxRedialTimes < 0:
			// 无限重连模式，每1分钟一次
			time.Sleep(1 * time.Minute)
		case redialTimes < maxRedialTimes:
			redialTimes++
			// 有限重连模式，每5秒一次
			log.Printf("Dial failed, start the %dth reconnection. error: %s", redialTimes, err.Error())
			time.Sleep(redialIntervalTime * time.Second)
		default:
			log.Println("Dial failed ->", err.Error())
			return nil
		}
	}
}

func handleClientConn(localAddr, serverAddr string) {
	for {
		conn := dial(localAddr, 0)
		if conn == nil {
			return
		}
		serverConn := dial(serverAddr, 0)
		if serverConn == nil {
			return
		}
		_, err := serverConn.Write([]byte(localAddr))
		if err != nil {
			log.Println("Fail to send address", err.Error())
			return
		}
		forward(conn, serverConn)
	}
}

func Client(config ClientConfig) {
	for _, addr := range config.GetLocalAddr() {
		go handleClientConn(addr, config.ServerAddr)
	}
	select {}
}
