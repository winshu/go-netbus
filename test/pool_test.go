package test

import (
	"../pool"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

// 测试连接复用

var (
	addr = "127.0.0.1:6666"
)

func server() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on ", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		//go handleRequest(conn)
	}
}

func BenchmarkPool(b *testing.B) {
	go server()
	time.Sleep(1e9)

	//factory 创建连接的方法
	//close 关闭连接的方法
	//创建一个连接池： 初始化5，最大连接30
	poolConfig := &pool.Config{
		InitialCap: 2,
		MaxCap:     10,
		Factory: func() (interface{}, error) {
			log.Println("new dial to", addr)
			return net.Dial("tcp", addr)
		},
		Close: func(v interface{}) error { return v.(net.Conn).Close() },
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}

	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		log.Fatalln("err=", err)
	}

	for i := 0; i < 10; i++ {
		go func() {
			v, err := p.Get()
			if err != nil {
				log.Println("Get error", err.Error())
				return
			}
			conn := v.(net.Conn)
			log.Println("conn", &conn)
			time.Sleep(100 * time.Millisecond)

			err = p.Put(v)
			if err != nil {
				log.Println("Put error", err.Error())
				return
			}
			fmt.Println("len=", p.Len())
		}()
	}
	time.Sleep(5e9)
}
