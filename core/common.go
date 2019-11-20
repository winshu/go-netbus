package core

import (
	"../config"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	// 重连间隔时间
	retryIntervalTime = 5
	// 固定报文头长度
	_headerLengthInByte = 32
)

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

// 固定报文头长度
func _formatHeader(header string) string {
	if len(header) > _headerLengthInByte {
		return header[:_headerLengthInByte]
	}
	return header + strings.Repeat(" ", _headerLengthInByte-len(header))
}

// 发送消息头，包含了地址信息
func sendHeader(conn net.Conn, address config.NetAddress) bool {
	_header := _formatHeader(address.String())
	if _, err := conn.Write([]byte(_header)); err != nil {
		log.Printf("Send header failed. [%s] %s\n", _header, err.Error())
		_ = conn.Close()
		return false
	}
	return true
}

// 接收消息头，包含了地址信息
func receiveHeader(conn net.Conn) (config.NetAddress, bool) {
	buffer := make([]byte, _headerLengthInByte)
	_, err := conn.Read(buffer)
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
