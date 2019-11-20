package config

import (
	"github.com/go-ini/ini"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type ServerConfig struct {
	// 服务端端口
	Port int
	// 代理端口模式
	PortMode int
}

// 代理端口模式
const (
	portModeRandom    = 0 // 使用随机端口(60000+)
	portModeIdentical = 1 // 使用同名端口
	portModeOffset    = 2 // 偏移端口(比如被代理端口是 3000，代理端口就是 4000)

	portOffset = 1000
)

var serverConfig ServerConfig

// 检查 PortMode 合法性
func _checkPortMode(portMode int) bool {
	return portMode >= 0 && portMode <= 2
}

// 从参数中解析配置
func _parseServerConfig(args []string) {
	if len(args) == 0 {
		log.Fatalln("More args in need")
	}
	// port, portMode
	port, err := strconv.Atoi(strings.TrimSpace(args[0]))
	portMode := portModeRandom
	if len(args) >= 2 {
		portMode, err = strconv.Atoi(strings.TrimSpace(args[1]))
	}
	if err != nil || !checkPort(port) || !_checkPortMode(portMode) {
		log.Fatalln("Fail to parse args.", args)
	}
	serverConfig = ServerConfig{Port: port, PortMode: portMode}
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
	portMode, err := server("port-mode").Int()
	if err != nil {
		log.Fatalln("Fail to parse config.ini", err.Error())
	}

	serverConfig = ServerConfig{Port: port, PortMode: portMode}
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

// 生成代理端口
func NewProxyPort(portMode, port int) int {
	switch portMode {
	case portModeIdentical:
		// 直接返回被代理端口
		return port
	case portModeOffset:
		// 基于被代理端口，加上偏移量
		return port + portOffset
	case portModeRandom:
		// 随机生成 60000+ 的端口
		return 60000 + rand.Intn(5535)
	}
	log.Printf("Make proxy port failed, use original port %d replace\n", port)
	return port
}

// 检查端口是否可被代理
func CheckProxyPort(portMode, port int) bool {
	switch portMode {
	case portModeIdentical:
		// 同名端口模式下，不能代理保留端口
		return port > 1024
	case portModeOffset:
		// 偏移模式下，被代理端口不能超过 64535
		return port < 65535-portOffset
	default:
		return true
	}
}
