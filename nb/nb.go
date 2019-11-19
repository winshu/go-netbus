package nb

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// 网上的一个案例

const timeout = 5

func Port2Port(port1 string, port2 string) {
	listen1 := listen("0.0.0.0:" + port1)
	listen2 := listen("0.0.0.0:" + port2)
	log.Println("[√]", "listen port:", port1, "and", port2, "success. waiting for client...")
	for {
		conn1 := accept(listen1)
		conn2 := accept(listen2)
		if conn1 == nil || conn2 == nil {
			log.Println("[x]", "accept client failed. retry in ", timeout, " seconds. ")
			time.Sleep(timeout * time.Second)
			continue
		}
		forward(conn1, conn2)
	}
}

func Port2Host(allowPort string, targetAddress string) {
	server := listen("0.0.0.0:" + allowPort)
	for {
		conn := accept(server)
		if conn == nil {
			continue
		}
		go func(targetAddress string) {
			log.Println("[+]", "start connect host:["+targetAddress+"]")
			target, err := net.Dial("tcp", targetAddress)
			if err != nil {
				// temporarily unavailable, don't use fatal.
				log.Println("[x]", "connect target address ["+targetAddress+"] failed. retry in ", timeout, "seconds. ")
				conn.Close()
				log.Println("[←]", "close the connect at local:["+conn.LocalAddr().String()+"] and remote:["+conn.RemoteAddr().String()+"]")
				time.Sleep(timeout * time.Second)
				return
			}
			log.Println("[→]", "connect target address ["+targetAddress+"] success.")
			forward(target, conn)
		}(targetAddress)
	}
}

func Host2Host(address1, address2 string) {
	for {
		log.Println("[+]", "try to connect host:["+address1+"] and ["+address2+"]")
		var host1, host2 net.Conn
		var err error
		for {
			host1, err = net.Dial("tcp", address1)
			if err == nil {
				log.Println("[→]", "connect ["+address1+"] success.")
				break
			} else {
				log.Println("[x]", "connect target address ["+address1+"] failed. retry in ", timeout, " seconds. ")
				time.Sleep(timeout * time.Second)
			}
		}
		for {
			host2, err = net.Dial("tcp", address2)
			if err == nil {
				log.Println("[→]", "connect ["+address2+"] success.")
				break
			} else {
				log.Println("[x]", "connect target address ["+address2+"] failed. retry in ", timeout, " seconds. ")
				time.Sleep(timeout * time.Second)
			}
		}
		forward(host1, host2)
	}
}

func listen(address string) net.Listener {
	log.Println("[+]", "try to start server on:["+address+"]")
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("[x]", "listen address ["+address+"] faild.")
	}
	log.Println("[√]", "start listen at address:["+address+"]")
	return server
}

func accept(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		log.Println("[x]", "accept connect ["+conn.RemoteAddr().String()+"] failed.", err.Error())
		return nil
	}
	log.Println("[√]", "accept a new client. remote address:["+conn.RemoteAddr().String()+"], local address:["+conn.LocalAddr().String()+"]")
	return conn
}

func forward(conn1 net.Conn, conn2 net.Conn) {
	log.Printf("[+] start transmit. [%s],[%s] <-> [%s],[%s] \n", conn1.LocalAddr().String(), conn1.RemoteAddr().String(), conn2.LocalAddr().String(), conn2.RemoteAddr().String())
	var wg sync.WaitGroup
	// wait tow goroutines
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	//blocking when the wg is locked
	wg.Wait()
}

func connCopy(conn1 net.Conn, conn2 net.Conn, wg *sync.WaitGroup) {
	// 读数据，读完了不要关闭连接
	_, err := io.Copy(conn1, conn2)
	if err != nil {
		log.Println("Copy error", err.Error())
	}
	conn1.Close()
	log.Println("[←]", "close the connect at local:["+conn1.LocalAddr().String()+"] and remote:["+conn1.RemoteAddr().String()+"]")
	wg.Done()
}
