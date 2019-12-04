package config

import (
	"../util"
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
)

const (
	// 默认最大重连次数
	defaultMaxRedialTimes = 20
)

// 客户端配置
type ClientConfig struct {
	Key            string       // 参考服务端配置 custom-port-key random-port-key
	ServerAddr     NetAddress   // 服务端地址
	LocalAddr      []NetAddress // 内网服务地址
	AccessPort     []int        // 代理访问端口(不能为空)
	MaxRedialTimes int          // 最大重连次数
}

var clientConfig ClientConfig

// TODO 如果 AccessPort 不为空，需要校验长度是否与 LocalAddr 一致

// 从参数中解析配置
func _parseClientConfig(args []string) ClientConfig {
	if len(args) < 3 {
		log.Fatalln("More args in need.", args)
	}

	config := ClientConfig{MaxRedialTimes: defaultMaxRedialTimes}
	var ok bool

	// 1 Key
	config.Key = strings.TrimSpace(args[0])
	// 2 ServerAddr
	if config.ServerAddr, ok = ParseNetAddress(strings.TrimSpace(args[1])); !ok {
		log.Fatalln("Fail to parse ServerAddr")
	}
	// 3 LocalAddr
	if config.LocalAddr, ok = ParseNetAddresses(strings.TrimSpace(args[2])); !ok {
		log.Fatalln("Fail to parse LocalAddr")
	}
	// 4 AccessPort
	var err error

	if len(args) >= 4 {
		if config.AccessPort, err = util.AtoInt2(args[3]); err != nil {
			log.Fatalln("Fail to parse AccessPort")
		}
	}
	// 如果未配置访问端口，则访问端口与内网服务端口相同
	if len(config.AccessPort) == 0 {
		config.AccessPort = ExtractPorts(config.LocalAddr)
	}
	// 如果访问端口与内网服务地址不一样，则配置检查不通过
	if len(config.AccessPort) != len(config.LocalAddr) {
		log.Fatalln("len(AccessPort) must equals len(LocalAddr)")
	}
	// 5 MaxRedialTimes
	if len(args) >= 5 {
		if config.MaxRedialTimes, err = strconv.Atoi(args[4]); err != nil {
			log.Fatalln("Fail to parse MaxRedialTimes")
		}
	}
	return config
}

// 从配置文件中加载配置
func _loadClientConfig() ClientConfig {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to load config.ini", err.Error())
	}

	client := func(key string) *ini.Key {
		return cfg.Section("client").Key(key)
	}
	args := make([]string, 5)

	args[0] = client("key").String()
	args[1] = client("server-host").String()
	args[2] = client("local-host").String()
	args[3] = client("access-port").String()
	args[4] = client("max-redial-times").String()

	return _parseClientConfig(args)
}

// 初始化客户端配置，支持从参数中读取或者从配置文件中读取
func InitClientConfig(args []string) ClientConfig {
	if len(args) == 0 {
		clientConfig = _loadClientConfig()
	} else {
		clientConfig = _parseClientConfig(args)
	}
	return clientConfig
}
