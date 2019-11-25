package config

import (
	"github.com/go-ini/ini"
	"log"
	"strings"
)

type ClientConfig struct {
	ServerAddr     NetAddress
	LocalAddr      []NetAddress
	AccessPort     []int
	MaxRedialTimes int
}

var clientConfig ClientConfig

// 从参数中解析配置
func _parseClientConfig(args []string) {
	if len(args) < 2 {
		log.Fatalln("More args in need.", args)
	}
	// 解析地址
	serverAddr, ok1 := ParseNetAddress(strings.TrimSpace(args[0]))
	localAddr, ok2 := ParseNetAddresses(strings.TrimSpace(args[1]))
	if !ok1 || !ok2 {
		log.Fatalln("Fail to parse address, the format is 'ip:port', such as '127.0.0.1:1024'")
	}

	clientConfig = ClientConfig{
		ServerAddr:     serverAddr,
		LocalAddr:      localAddr,
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
	//_accessPort := client("access-port").String()
	maxRedialTimes, err := client("max-redial-times").Int()
	if err != nil {
		maxRedialTimes = 20
	}

	// 解析地址
	serverAddr, ok1 := ParseNetAddress(_serverAddr)
	localAddr, ok2 := ParseNetAddresses(_localAddr)
	if !ok1 || !ok2 {
		log.Fatalln("Fail to parse address, the format is 'ip:port', such as '127.0.0.1:1024'")
	}

	clientConfig = ClientConfig{
		ServerAddr:     serverAddr,
		LocalAddr:      localAddr,
		AccessPort:     []int{7001},
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
