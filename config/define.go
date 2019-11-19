package config

import (
	"github.com/go-ini/ini"
	"log"
)

func loadConfig() *ini.File {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to load config", err.Error())
	}
	return cfg
}
