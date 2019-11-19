package main

import (
	"./config"
	"./core"
	"fmt"
	"log"
	"os"
	"time"
)

func printWelcome(args []string) {
	fmt.Println("+----------------------------------------------------------------+")
	fmt.Println("| Welcome to use NetBus version 1.0.0 .                          |")
	fmt.Println("| Code by winshu at 2019-10-19                                   |")
	fmt.Println("| If you have some problem when you use the tool,                |")
	fmt.Println("| Please submit issue at : https://gitee.com/winshu/go-netbus .  |")
	fmt.Println("+----------------------------------------------------------------+")
	fmt.Println()
	// sleep one second because the fmt is not thread-safety.
	// if not to do this, fmt.Print will print after the log.Print.
	time.Sleep(time.Second)
}

func printHelp() {
	fmt.Println(`method A: "-server" load "config.ini" and start as server`)
	fmt.Println(`          "-client " load "config.ini" and start as client`)
	fmt.Println(`method B: "-server <port>" start as server, and listening at port x', e.g. -server 6666`)
	fmt.Println(`          "-client <server:port> <local:port>" start as client, e.g. -client 123.54.23.67:6666 127.0.0.1:3306`)
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	args := os.Args
	argc := len(os.Args)
	printWelcome(args)

	if argc < 2 {
		printHelp()
		os.Exit(0)
	}

	// 获取其余参数
	argsConfig := args[2:]

	switch args[1] {
	case "-server":
		serverConfig := config.InitServerConfig(argsConfig)
		core.Server(serverConfig)
	case "-client":
		clientConfig := config.InitClientConfig(argsConfig)
		core.Client(clientConfig)
	default:
		printHelp()
	}
}
