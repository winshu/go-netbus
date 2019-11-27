package config

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// 网络地址
type NetAddress struct {
	IP   string
	Port int
}

// 转字符串
func (t NetAddress) String() string {
	return fmt.Sprintf("%s:%d", t.IP, t.Port)
}

// 解析多个地址
func ParseNetAddresses(addresses string) ([]NetAddress, bool) {
	arr := strings.Split(addresses, ",")
	result := make([]NetAddress, len(arr))

	var ok bool
	for i, addr := range arr {
		result[i], ok = ParseNetAddress(addr)
		if !ok {
			return nil, false
		}
	}
	return result, true
}

// 解析单个地址
func ParseNetAddress(address string) (NetAddress, bool) {
	arr := strings.Split(strings.TrimSpace(address), ":")
	if len(arr) != 2 {
		log.Println("Fail to parse address")
		return NetAddress{}, false
	}

	ip := strings.TrimSpace(arr[0])
	ipPattern := `^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$`
	if ok, err := regexp.MatchString(ipPattern, ip); !ok || err != nil {
		log.Println("Fail to parse address ip")
		return NetAddress{}, false
	}

	port, err := strconv.Atoi(strings.TrimSpace(arr[1]))
	if err != nil || !checkPort(port) {
		log.Println("Fail to parse address port")
		return NetAddress{}, false
	}
	return NetAddress{ip, port}, true
}

// 检查端口是否合法
func checkPort(port int) bool {
	return port > 0 && port <= 65535
}

// 从地址中抽取出端口
func extractPorts(address []NetAddress) []int {
	accessPort := make([]int, len(address))
	for i, addr := range address {
		accessPort[i] = addr.Port
	}
	return accessPort
}
