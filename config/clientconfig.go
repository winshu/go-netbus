package config

import (
	"github.com/go-ini/ini"
	"log"
	"strings"
)

type ClientConfig struct {
	ServerAddr     string
	LocalAddr      string
	MaxRedialTimes int
}

func (t ClientConfig) GetLocalAddr() []string {
	str := strings.ReplaceAll(t.LocalAddr, " ", "")
	return strings.Split(str, ",")
}

var clientConfig ClientConfig

// 从参数中解析配置
func parseClientConfig(args []string) {
	if len(args) < 2 {
		log.Fatalln("More args in need")
	}
	serverAddr := strings.TrimSpace(args[0])
	localAddr := strings.TrimSpace(args[1])

	clientConfig = ClientConfig{
		ServerAddr:     serverAddr,
		LocalAddr:      localAddr,
		MaxRedialTimes: 20,
	}
	log.Println("Init client config from args finished", clientConfig)
}

// 从配置文件中加载配置
func loadClientConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to load config.ini", err.Error())
	}

	client := func(key string) *ini.Key {
		return cfg.Section("client").Key(key)
	}
	serverAddr := client("server-host").String()
	localAddr := client("local-host").String()
	maxRedialTimes, err := client("max-redial-times").Int()
	if err != nil {
		maxRedialTimes = 20
	}
	clientConfig = ClientConfig{
		ServerAddr:     serverAddr,
		LocalAddr:      localAddr,
		MaxRedialTimes: maxRedialTimes,
	}
	log.Println("Init client config from config.ini finished", clientConfig)
}

// 初始化客户端配置，支持从参数中读取或者从配置文件中读取
func InitClientConfig(args []string) ClientConfig {
	if len(args) == 0 {
		loadClientConfig()
	} else {
		parseClientConfig(args)
	}
	return clientConfig
}
