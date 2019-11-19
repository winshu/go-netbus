package core

import (
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
	headerLengthInByte = 32
)

func connCopy(source, target net.Conn, wg *sync.WaitGroup) {
	_, err := io.Copy(source, target)
	if err != nil {
		log.Println("Connection interrupted")
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
func FormatHeader(header string) string {
	if len(header) > headerLengthInByte {
		return header[:headerLengthInByte]
	}
	return header + strings.Repeat(" ", headerLengthInByte-len(header))
}
