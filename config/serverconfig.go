package config

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
)

// 服务端配置
type ServerConfig struct {
	Port          int    // 服务端口
	CustomPortKey string // 自定义端口的 Key
	RandomPortKey string // 随机端口的 Key
}

var serverConfig ServerConfig

// 从参数中解析配置
func _parseServerConfig(args []string) ServerConfig {
	if len(args) < 3 {
		log.Fatalln("More args in need")
	}
	// 1 port
	port, err := strconv.Atoi(strings.TrimSpace(args[0]))
	if err != nil || !checkPort(port) {
		log.Fatalln("Fail to parse args.", args)
	}
	// 2 custom-port-key
	customPortKey := strings.TrimSpace(args[1])
	// 3 random-port-key
	randomPortKey := strings.TrimSpace(args[2])

	return ServerConfig{
		Port:          port,
		CustomPortKey: customPortKey,
		RandomPortKey: randomPortKey,
	}
}

// 从配置文件中加载配置
func _loadServerConfig() ServerConfig {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to load config.ini", err.Error())
	}
	server := func(key string) *ini.Key {
		return cfg.Section("server").Key(key)
	}

	args := make([]string, 3)
	args[0] = server("port").String()
	args[1] = server("custom-port-key").String()
	args[2] = server("random-port-key").String()

	return _parseServerConfig(args)
}

// 初始化服务端配置，支持从参数中读取或者从配置文件中读取
func InitServerConfig(args []string) ServerConfig {
	if len(args) == 0 {
		serverConfig = _loadServerConfig()
	} else {
		serverConfig = _parseServerConfig(args)
	}
	log.Println("Init server config from config.ini finished", serverConfig)
	return serverConfig
}
