package config

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Port       int
	RandomPort bool
}

var serverConfig ServerConfig

// ------------------------------------------------------------------------------
func parseServerConfig(args []string) {
	if len(args) == 0 {
		log.Fatalln("More args in need")
	}
	port, err := strconv.Atoi(strings.TrimSpace(args[0]))
	if err != nil {
		log.Fatalln("Parse args failed")
	}

	serverConfig = ServerConfig{Port: port, RandomPort: false}
	log.Println("Init server config from args finished", serverConfig)
}

func loadServerConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to load config.ini", err.Error())
	}

	server := func(key string) *ini.Key {
		return cfg.Section("server").Key(key)
	}
	port, _ := server("port").Int()
	randomPort, _ := server("random-port").Bool()
	serverConfig = ServerConfig{Port: port, RandomPort: randomPort}

	log.Println("Init server config from config.ini finished", serverConfig)
}

func InitServerConfig(args []string) ServerConfig {
	if len(args) == 0 {
		loadServerConfig()
	} else {
		parseServerConfig(args)
	}
	return serverConfig
}

// ------------------------------------------------------------------------------
