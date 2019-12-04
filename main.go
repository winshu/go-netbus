package main

import (
	"./config"
	"./core"
	"./nb"
	"fmt"
	"log"
	"os"
	"time"
)

func printWelcome(args []string) {
	fmt.Println("+----------------------------------------------------------------+")
	fmt.Println("| Welcome to use NetBus version 1.0.2                            |")
	fmt.Println("| Code by winshu at 2019/12/04                                   |")
	fmt.Println("| If you have some problem when you use the tool,                |")
	fmt.Println("| Please submit issue at : https://gitee.com/winshu/go-netbus    |")
	fmt.Println("+----------------------------------------------------------------+")
	fmt.Println()
	// sleep one second because the fmt is not thread-safety.
	// if not to do this, fmt.Print will print after the log.Print.
	time.Sleep(500 * time.Millisecond)
}

func printHelp() {
	fmt.Println(`A: "-server" load "config.ini" and start as server`)
	fmt.Println(`   "-client " load "config.ini" and start as client`)
	fmt.Println(`B: "-server <port>" start as server, and listening at port x', e.g. -server 6666`)
	fmt.Println(`   "-client <server:port> <local:port>" start as client, e.g. -client 123.54.23.67:6666 127.0.0.1:3306`)
	fmt.Println(`more details please read "README.md"`)
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
		// 外网
		clientConfig := config.InitClientConfig(argsConfig)
		core.Client(clientConfig)
	case "-nbt":
		// 隐藏彩蛋，支持端口转发
		if len(argsConfig) == 2 {
			nb.Port2Host(argsConfig[0], argsConfig[1])
		}
	case "-nbs":
		// 隐藏彩蛋，单端口服务端
		if len(argsConfig) == 2 {
			nb.Port2Port(argsConfig[0], argsConfig[1])
		}
	case "-nbc":
		// 隐藏彩蛋，单端口服务端
		if len(argsConfig) == 2 {
			// 隐藏彩蛋，支持端口转发
			nb.Host2Host(argsConfig[0], argsConfig[1])
		}
	default:
		printHelp()
	}
}
