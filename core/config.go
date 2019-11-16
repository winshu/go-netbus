package core

import (
	"github.com/go-ini/ini"
	"log"
	"strings"
)

type ServerConfig struct {
	Port       int
	RandomPort bool
}

type ClientConfig struct {
	ServerAddr string
	LocalAddr  string
}

func (t ClientConfig) GetLocalAddr() []string {
	str := strings.ReplaceAll(t.LocalAddr, " ", "")
	return strings.Split(str, ",")
}

var (
	serverConfig ServerConfig
	clientConfig ClientConfig
)

func loadConfig() *ini.File {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to load config", err.Error())
	}
	return cfg
}

func InitServerConfig() ServerConfig {
	cfg := loadConfig()
	server := func(key string) *ini.Key {
		return cfg.Section("server").Key(key)
	}
	port, _ := server("port").Int()
	randomPort, _ := server("random-port").Bool()
	serverConfig = ServerConfig{Port: port, RandomPort: randomPort}

	log.Println("Init server config finished", serverConfig)
	return serverConfig
}

func InitClientConfig() ClientConfig {
	cfg := loadConfig()
	client := func(key string) *ini.Key {
		return cfg.Section("client").Key(key)
	}
	serverAddr := client("server-host").String()
	localAddr := client("local-host").String()
	clientConfig = ClientConfig{ServerAddr: serverAddr, LocalAddr: localAddr}
	log.Println("Init client config finished", clientConfig)
	return clientConfig
}
