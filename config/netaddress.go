package config

import (
	"fmt"
	"log"
	"regexp"
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

// 检查地址是否合法，支持多个
func checkNetAddress(addresses string) bool {
	for _, addr := range strings.Split(addresses, ",") {
		if _, ok := ParseNetAddress(addr); !ok {
			return false
		}
	}
	return true
}

// 检查端口是否合法
func checkPort(port int) bool {
	return port > 0 && port <= 65535
}
