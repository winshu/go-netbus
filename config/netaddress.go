package config

import (
	"fmt"
	"strconv"
	"strings"
)

type NetAddress struct {
	IP   string
	Port int
}

func (t NetAddress) String() string {
	return fmt.Sprintf("%s:%d", t.IP, t.Port)
}

// 解析地址
func ParseNetAddress(host string) NetAddress {
	arr := strings.Split(host, ":")
	if len(arr) != 2 {
		return NetAddress{}
	}
	port, err := strconv.Atoi(strings.TrimSpace(arr[1]))
	if err != nil {
		port = 0
	}
	return NetAddress{strings.TrimSpace(arr[0]), port}
}
