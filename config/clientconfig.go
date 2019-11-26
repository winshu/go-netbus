package config

import (
	"../util"
	"github.com/go-ini/ini"
	"log"
	"strings"
)

type ClientConfig struct {
	Key            string
	ServerAddr     NetAddress
	LocalAddr      []NetAddress
	AccessPort     []int
	MaxRedialTimes int
}

var clientConfig ClientConfig

// TODO 如果 AccessPort 不为空，需要校验长度是否与 LocalAddr 一致

func _copyPort(localAddr []NetAddress) []int {
	accessPort := make([]int, len(localAddr))
	for i, addr := range localAddr {
		accessPort[i] = addr.Port
	}
	return accessPort
}

// 从参数中解析配置
func _parseClientConfig(args []string) {
	if len(args) < 3 {
		log.Fatalln("More args in need.", args)
	}
	key := strings.TrimSpace(args[0])
	// 解析地址
	serverAddr, ok1 := ParseNetAddress(strings.TrimSpace(args[1]))
	localAddr, ok2 := ParseNetAddresses(strings.TrimSpace(args[2]))
	if !ok1 || !ok2 {
		log.Fatalln("Fail to parse address, the format is 'ip:port', such as '127.0.0.1:1024'")
	}

	var accessPort []int
	var err error
	if len(args) >= 4 {
		accessPort, err = util.AtoInt2(args[3])
		if err != nil {
			log.Fatalln("Fail to parse AccessPort")
		}
	}
	if len(accessPort) == 0 {
		accessPort = _copyPort(localAddr)
	}

	clientConfig = ClientConfig{
		Key:            key,
		ServerAddr:     serverAddr,
		LocalAddr:      localAddr,
		AccessPort:     accessPort,
		MaxRedialTimes: 20,
	}
	log.Println("Init client config from args finished", clientConfig)
}

// 从配置文件中加载配置
func _loadClientConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to load config.ini", err.Error())
	}

	client := func(key string) *ini.Key {
		return cfg.Section("client").Key(key)
	}
	_serverAddr := client("server-host").String()
	_localAddr := client("local-host").String()

	// 解析地址
	serverAddr, ok1 := ParseNetAddress(_serverAddr)
	localAddr, ok2 := ParseNetAddresses(_localAddr)
	if !ok1 || !ok2 {
		log.Fatalln("Fail to parse address, the format is 'ip:port', such as '127.0.0.1:1024'")
	}

	_accessPort := client("access-port").String()
	var accessPort []int
	if len(_accessPort) == 0 {
		accessPort = _copyPort(localAddr)
	}

	maxRedialTimes, err := client("max-redial-times").Int()
	if err != nil {
		maxRedialTimes = 20
	}

	clientConfig = ClientConfig{
		ServerAddr:     serverAddr,
		LocalAddr:      localAddr,
		AccessPort:     accessPort,
		MaxRedialTimes: maxRedialTimes,
	}
	log.Println("Init client config from config.ini finished", clientConfig)
}

// 初始化客户端配置，支持从参数中读取或者从配置文件中读取
func InitClientConfig(args []string) ClientConfig {
	if len(args) == 0 {
		_loadClientConfig()
	} else {
		_parseClientConfig(args)
	}
	return clientConfig
}
