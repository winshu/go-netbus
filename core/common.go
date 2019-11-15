package core

import (
	"log"
	"net"
	"sync"
)

func forward(conn1 net.Conn, conn2 net.Conn) {
	log.Printf("[+] start transmit. [%s],[%s] <-> [%s],[%s] \n", conn1.LocalAddr(), conn1.RemoteAddr(), conn2.LocalAddr(), conn2.RemoteAddr())
	var wg sync.WaitGroup
	// wait tow goroutines
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	//blocking when the wg is locked
	wg.Wait()
}

// 受理请求
func accept(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("[x] accept connect [%s] failed. %s", conn.RemoteAddr(), err.Error())
		return nil
	}
	log.Printf("[√] accept a new client. remote address:[%s], local address:[%s]", conn.RemoteAddr(), conn.LocalAddr())
	return conn
}
