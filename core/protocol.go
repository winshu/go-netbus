package core

import (
	"../util"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	// 协议长度=字段数量
	protocolLength = 4

	// 协议-结果
	protocolResultFail    = 0 // 失败，默认值
	protocolResultSuccess = 1 // 成功

	// 协议-消息类型
	protocolTypeNormal = 0 // 正常，默认值
	protocolTypeAuth   = 1 // 鉴权

	// Token 长度
	protocolTokenLength = 16 // Token 长度
)

// 协议格式
// 结果|消息类型|端口列表|令牌
// 0|0|7001,7002|customabcd8000

// 协议
type Protocol struct {
	Result int    // 结果：0 失败，1 成功
	Type   int    // 消息类型：0 正常，1 鉴权
	Ports  []int  // 端口列表
	Token  string // 令牌
}

// 转字符串
func (h *Protocol) String() string {
	ports := strings.Replace(strings.Trim(fmt.Sprint(h.Ports), "[]"), " ", ",", -1)
	return fmt.Sprintf("%d|%d|%s|%s", h.Result, h.Type, ports, h.Token)
}

// 解析协议
func _parseProtocol(body string) Protocol {
	// 拆解字符
	arr := strings.Split(body, "|")
	if len(arr) != protocolLength {
		return Protocol{}
	}
	// 前两个字段
	params, err := util.AtoInt(arr[0:2])
	if err != nil {
		return Protocol{}
	}
	// 第三个字段
	var ports []int
	if len(arr[2]) > 0 {
		if ports, err = util.AtoInt2(arr[2]); err != nil {
			return Protocol{}
		}
	} else {
		return Protocol{}
	}

	return Protocol{
		Result: params[0],
		Type:   params[1],
		Ports:  ports,
		Token:  arr[3],
	}
}

// 发送协议
// 第一个字节为协议长度
// 协议长度只支持到255
func sendProtocol(conn net.Conn, protocol Protocol) bool {
	buffer := bytes.NewBuffer([]byte{})
	length := byte(len(protocol.String()))

	// 协议长度
	buffer.WriteByte(length)
	// 协议内容
	buffer.WriteString(protocol.String())

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		log.Printf("Send protocol failed. [%s] %s\n", protocol.String(), err.Error())
		_ = conn.Close()
		return false
	}
	//log.Println("Send protocol", protocol.String())
	return true
}

// 接收协议
// 第一个字节为协议长度
func receiveProtocol(conn net.Conn) Protocol {
	var err error

	// 读取协议长度
	buffer := make([]byte, 1)
	if _, err := conn.Read(buffer); err != nil {
		log.Println("Parse protocol length failed.", err.Error())
		return Protocol{}
	}
	// 读取协议内容
	buffer = make([]byte, buffer[0])
	if _, err = conn.Read(buffer); err != nil {
		log.Println("Parse protocol body failed.", err.Error())
		return Protocol{}
	}
	// 解析消息
	body := strings.TrimSpace(string(buffer))
	//log.Println("----------> Receive protocol", body)
	return _parseProtocol(body)
}
