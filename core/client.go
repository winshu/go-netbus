package core

import (
	"log"
	"net"
	"time"
)

func dial(targetAddr string, maxRedialTimes int) net.Conn {
	redialTimes := 0
	for {
		conn, err := net.Dial("tcp", targetAddr)
		if err == nil {
			log.Println("Connection success ->", targetAddr)
			return conn
		}
		switch {
		case maxRedialTimes < 0:
			// 无限重连模式，每1分钟一次
			time.Sleep(1 * time.Minute)
		case redialTimes < maxRedialTimes:
			redialTimes++
			// 有限重连模式，每5秒一次
			log.Printf("Connection failed, start the %dth reconnection. error: %s", redialTimes, err.Error())
			time.Sleep(redialIntervalTime * time.Second)
		default:
			log.Println("Connection failed ->", err.Error())
			return nil
		}
	}
}

func Client(config ClientConfig) {
	serverConn := dial(config.ServerAddr, 0)
	//for {
	conn := dial(config.LocalAddr, 0)
	forward(conn, serverConn)
	//}
}
