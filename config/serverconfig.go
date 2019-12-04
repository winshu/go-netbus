package config

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
)

// 服务端配置
type ServerConfig struct {
	Port int    // 服务端口
	Key  string // 6-16 个字符，用于身份校验
}

var serverConfig ServerConfig

// 从参数中解析配置
func _parseServerConfig(args []string) ServerConfig {
	if len(args) < 2 {
		log.Fatalln("More args in need")
	}
	// 0 key
	key := strings.TrimSpace(args[0])

	// 1 port
	port, err := strconv.Atoi(strings.TrimSpace(args[1]))
	if err != nil || !checkPort(port) {
		log.Fatalln("Fail to parse args.", args)
	}

	return ServerConfig{
		Port: port,
		Key:  key,
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

	args := make([]string, 2)
	args[0] = server("key").String()
	args[1] = server("port").String()

	return _parseServerConfig(args)
}

// 初始化服务端配置，支持从参数中读取或者从配置文件中读取
func InitServerConfig(args []string) ServerConfig {
	if len(args) == 0 {
		serverConfig = _loadServerConfig()
	} else {
		serverConfig = _parseServerConfig(args)
	}
	return serverConfig
}
