package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	// 协议-结果
	protocolResultFail            = 0 // 失败，默认值
	protocolResultSuccess         = 1 // 成功
	protocolResultFailToSend      = 2 // 发送失败
	protocolResultFailToReceive   = 3 // 接收失败
	protocolResultFailToAuth      = 4 // 鉴权失败
	protocolResultFailToListen    = 5 // 监听失败
	protocolResultVersionMismatch = 6 // 版本不匹配

	// 协议发送超时时间
	protocolSendTimeout = 3

	// 版本号(单调递增)
	protocolVersion = 1
)

// 协议格式
// 结果|版本号|原端口|访问端口|Key
// 1|0|1|3306|13306|winshu

// 协议
type Protocol struct {
	Result     byte   // 结果：0 失败，1 成功
	Version    uint32 // 版本号，单调递增
	AccessPort uint32 // 访问端口
	Port       uint32 // 原端口
	Key        string // 身份验证
}

// 转字符串
func (p *Protocol) String() string {
	return fmt.Sprintf("%d|%d|%d|%d|%s", p.Result, p.Version, p.Port, p.AccessPort, p.Key)
}

// 返回一个新结果
func (p *Protocol) NewResult(newResult byte) Protocol {
	return Protocol{
		Result:     newResult,
		Version:    p.Version,
		Port:       p.Port,
		AccessPort: p.AccessPort,
		Key:        p.Key,
	}
}

func (p *Protocol) Bytes() []byte {
	buffer := bytes.NewBuffer([]byte{})

	buffer.WriteByte(p.Result)
	_ = binary.Write(buffer, binary.BigEndian, p.Version)
	_ = binary.Write(buffer, binary.BigEndian, p.Port)
	_ = binary.Write(buffer, binary.BigEndian, p.AccessPort)
	buffer.WriteString(p.Key)
	return buffer.Bytes()
}

func (p *Protocol) Len() byte {
	return byte(len(p.Bytes()))
}

// 是否成功
func (p *Protocol) Success() bool {
	return p.Result == protocolResultSuccess
}

// 解析协议
func _parseProtocol(body []byte) Protocol {
	return Protocol{
		Result:     body[0],
		Version:    binary.BigEndian.Uint32(body[1:5]),
		Port:       binary.BigEndian.Uint32(body[5:9]),
		AccessPort: binary.BigEndian.Uint32(body[9:13]),
		Key:        string(body[13:]),
	}
	//log.Println("Parse Protocol", protocol.String())
}

// 发送协议
// 第一个字节为协议长度
// 协议长度只支持到255
func sendProtocol(conn net.Conn, req Protocol) bool {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte(req.Len())
	buffer.Write(req.Bytes())

	// 设置写超时时间，避免连接断开的问题
	if err := conn.SetWriteDeadline(time.Now().Add(protocolSendTimeout * time.Second)); err != nil {
		log.Println("Fail to set write deadline.", err.Error())
		return false
	}
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		log.Printf("Send protocol failed. [%s] %s\n", req.String(), err.Error())
		return false
	}
	// 清空写超时设置
	if err := conn.SetWriteDeadline(time.Time{}); err != nil {
		log.Println("Fail to clear write deadline.", err.Error())
		return false
	}
	//log.Println("Send protocol", req.String())
	return true
}

// 接收协议
// 第一个字节为协议长度
func receiveProtocol(conn net.Conn) Protocol {
	var err error
	var length byte

	if err = binary.Read(conn, binary.BigEndian, &length); err != nil {
		log.Println("Parse protocol length failed.", err.Error())
		return Protocol{Result: protocolResultFailToReceive}
	}
	// 读取协议内容
	body := make([]byte, length)
	if err = binary.Read(conn, binary.BigEndian, &body); err != nil {
		log.Println("Parse protocol body failed.", err.Error())
		return Protocol{Result: protocolResultFailToReceive}
	}
	return _parseProtocol(body)
}
