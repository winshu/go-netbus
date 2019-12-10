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
	Port uint32
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

	port, err := parsePort(arr[1])
	if err != nil || !checkPort(port) {
		log.Println("Fail to parse address port")
		return NetAddress{}, false
	}
	return NetAddress{ip, uint32(port)}, true
}

// 从地址中抽取出端口
func ExtractPorts(address []NetAddress) []uint32 {
	accessPort := make([]uint32, len(address))
	for i, addr := range address {
		accessPort[i] = addr.Port
	}
	return accessPort
}

// 解析多个端口
func parsePorts(str string) ([]uint32, error) {
	str = strings.ReplaceAll(str, " ", "")
	arr := strings.Split(str, ",")

	var err error
	result := make([]uint32, len(arr))

	for i, v := range arr {
		if result[i], err = parsePort(v); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// 解析单个端口
func parsePort(str string) (uint32, error) {
	var port int
	var err error
	str = strings.TrimSpace(str)
	if port, err = strconv.Atoi(str); err == nil {
		return uint32(port), nil
	}
	return 0, err
}

// 检查端口是否合法
func checkPort(port uint32) bool {
	return port > 0 && port <= 65535
}
