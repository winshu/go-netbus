package config

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
)

type ServerConfig struct {
	// 服务端端口
	Port          int
	CustomPortKey string
	RandomPortKey string
}

var serverConfig ServerConfig

// 从参数中解析配置
func _parseServerConfig(args []string) {
	if len(args) == 0 {
		log.Fatalln("More args in need")
	}
	// port, portMode
	port, err := strconv.Atoi(strings.TrimSpace(args[0]))
	if err != nil || !checkPort(port) {
		log.Fatalln("Fail to parse args.", args)
	}
	serverConfig = ServerConfig{Port: port}
	log.Println("Init server config from args finished", serverConfig)
}

// 从配置文件中加载配置
func _loadServerConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to load config.ini", err.Error())
	}

	server := func(key string) *ini.Key {
		return cfg.Section("server").Key(key)
	}
	port, err := server("port").Int()

	serverConfig = ServerConfig{Port: port}
	log.Println("Init server config from config.ini finished", serverConfig)
}

// 初始化服务端配置，支持从参数中读取或者从配置文件中读取
func InitServerConfig(args []string) ServerConfig {
	if len(args) == 0 {
		_loadServerConfig()
	} else {
		_parseServerConfig(args)
	}
	return serverConfig
}
