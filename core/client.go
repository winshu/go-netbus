package core

import (
	"../config"
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

func handleClientConn(localAddr, serverAddr string, maxRedialTimes int) {
	for {
		conn := dial(localAddr, maxRedialTimes)
		if conn == nil {
			return
		}
		serverConn := dial(serverAddr, maxRedialTimes)
		if serverConn == nil {
			return
		}
		header := FormatHeader(localAddr)
		_, err := serverConn.Write([]byte(header))
		if err != nil {
			log.Println("Send header error", err.Error())
			return
		}
		forward(conn, serverConn)
	}
}

func Client(cfg config.ClientConfig) {
	for _, addr := range cfg.GetLocalAddr() {
		go handleClientConn(addr, cfg.ServerAddr, cfg.MaxRedialTimes)
	}
	select {}
}
