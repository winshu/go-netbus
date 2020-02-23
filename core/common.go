package core

import (
	"io"
	"net"
	"sync"
)

const (
	// 重连间隔时间
	retryIntervalTime = 5
)

// 连接数据复制
func connCopy(source, target net.Conn, wg *sync.WaitGroup) {
	//if _, err := io.Copy(source, target); err != nil {
	if _, err := io.Copy(target, source); err != nil {
		//log.Println("Connection interrupted", err)
	}
	_ = source.Close()
	wg.Done()
}

// 连接转发
func forward(conn1, conn2 net.Conn) {
	//log.Printf("Forward channel [%s/%s] <-> [%s/%s]\n",
	//	conn1.RemoteAddr(), conn1.LocalAddr(), conn2.RemoteAddr(), conn2.LocalAddr())

	var wg sync.WaitGroup
	// wait tow goroutines
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	//blocking when the wg is locked
	wg.Wait()
}

// 关闭连接
func closeConn(connections ...net.Conn) {
	for _, conn := range connections {
		if conn != nil {
			_ = conn.Close()
		}
	}
}
