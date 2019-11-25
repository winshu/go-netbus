package core

import (
	"../util"
	"bytes"
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
	Result int    // 结果：0 失败，1 成功
	Type   int    // 消息类型：0 正常，1 鉴权
	Ports  []int  // 端口列表
	Token  string // 令牌
}

func (h *Header) String() string {
	ports := strings.Replace(strings.Trim(fmt.Sprint(h.Ports), "[]"), " ", ",", -1)
	return fmt.Sprintf("%d|%d|%s|%s", h.Result, h.Type, ports, h.Token)
}

func ParseHeader(body string) Header {
	arr := strings.Split(body, "|")
	if len(arr) != 4 {
		return Header{}
	}
	params, err := util.Atoi(arr[0:2])
	if err != nil {
		return Header{}
	}
	var ports []int
	if len(arr[2]) > 0 {
		portsArr := strings.Split(arr[2], ",")
		ports, err = util.Atoi(portsArr)
		if err != nil {
			return Header{}
		}
	} else {
		ports = []int{}
	}

	return Header{
		Result: params[0],
		Type:   params[1],
		Ports:  ports,
		Token:  arr[3],
	}
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

func closeConn(conn net.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}

// 消息头：长度|代理模式|端口列表
// 举例:  6|0|7001,7002,7003
// 发送消息头，包含了地址信息
func sendHeader(conn net.Conn, header Header) bool {
	buffer := bytes.NewBuffer([]byte{})

	length := byte(len(header.String()))
	buffer.WriteByte(length)
	buffer.WriteString(header.String())
	log.Println("Send header", header.String())

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		log.Printf("Send header failed. [%s] %s\n", header.String(), err.Error())
		_ = conn.Close()
		return false
	}
	return true
}

// 接收消息头，包含了地址信息
func receiveHeader(conn net.Conn) Header {
	// 读取消息长度
	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Println("Receive header failed.", err.Error())
		return Header{}
	}
	// 读取消息体
	length := buffer[0]
	buffer = make([]byte, length)
	_, err = conn.Read(buffer)
	if err != nil {
		log.Println("Receive header failed.", err.Error())
		return Header{}
	}
	// 解析消息
	body := strings.TrimSpace(string(buffer))
	log.Println("-----------------------------------------> Receive header", body)
	return ParseHeader(body)
}
