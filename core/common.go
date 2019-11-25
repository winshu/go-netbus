package core

import (
	"../config"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	// 重连间隔时间
	retryIntervalTime = 5
)

type Header struct {
	Result byte
	Type   byte
	Mode   byte
	Ports  []int
}

func (h Header) String() string {
	ports := strings.Replace(strings.Trim(fmt.Sprint(h.Ports), "[]"), " ", ",", -1)
	return fmt.Sprintf("%d|%d|%s", h.Result, h.Mode, ports)
}

func connCopy(source, target net.Conn, wg *sync.WaitGroup) {
	_, err := io.Copy(source, target)
	if err != nil {
		//log.Println("Connection interrupted")
	}
	_ = source.Close()
	wg.Done()
}

func forward(conn1, conn2 net.Conn) {
	log.Printf("Forward channel [%s/%s] <-> [%s/%s]\n",
		conn1.RemoteAddr(), conn1.LocalAddr(), conn2.RemoteAddr(), conn2.LocalAddr())

	var wg sync.WaitGroup
	// wait tow goroutines
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	//blocking when the wg is locked
	wg.Wait()
}

// 消息头：长度|代理模式|端口列表
// 举例:  6|0|7001,7002,7003
// 发送消息头，包含了地址信息
func sendHeader(conn net.Conn, header Header) bool {
	if _, err := conn.Write([]byte(header.String())); err != nil {
		log.Printf("Send header failed. [%s] %s\n", header.String(), err.Error())
		_ = conn.Close()
		return false
	}
	return true
}

// 接收消息头，包含了地址信息
func receiveHeader(conn net.Conn) (config.NetAddress, bool) {
	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Println("Receive header failed.", err.Error())
		_ = conn.Close()
		return config.NetAddress{}, false
	}
	length := buffer[0]
	buffer = make([]byte, length)
	_, err = conn.Read(buffer)
	if err != nil {
		log.Println("Receive header failed.", err.Error())
		_ = conn.Close()
		return config.NetAddress{}, false
	}

	header := strings.TrimSpace(string(buffer))
	log.Println("-----------------------------------------> Receive header", header)
	address, _ := config.ParseNetAddress(header)
	return address, true
}

func closeConn(conn net.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}
