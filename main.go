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
	fmt.Println("| Welcome to use NetBus version 1.0.5                            |")
	fmt.Println("| Code by winshu at 2019/12/15                                   |")
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
	fmt.Println(`B: "-server <key> <port>" start as server, and listening at port x', e.g. -server 6666`)
	fmt.Println(`   "-client <key> <server:port> <local:port> [access-port] [max-redial-times]" start as client,`)
	fmt.Println(`   "e.g. -client winshu 123.54.23.67:6666 127.0.0.1:3306`)
	fmt.Println(`Generate trial key: `)
	fmt.Println(`   "-generate <key> [expired-time]" make a trial client key, e.g. -generate winshu 2019-12-31`)
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
		// 外网服务
		serverConfig := config.InitServerConfig(argsConfig)
		core.Server(serverConfig)
	case "-client":
		// 内网服务
		clientConfig := config.InitClientConfig(argsConfig)
		core.Client(clientConfig)
	case "-generate":
		// 生成短期 key
		var seed, expired string
		if len(argsConfig) > 0 {
			seed = argsConfig[0]
		}
		if len(argsConfig) > 1 {
			expired = argsConfig[1]
		}
		if len(argsConfig) > 0 {
			trialKey, _ := config.NewKey(seed, expired)
			fmt.Println("You got a new key ->    ", trialKey)
		}
	case "-check":
		if len(argsConfig) == 2 {
			fmt.Println(config.CheckKey(argsConfig[0], argsConfig[1]))
		}
	default:
		printHelp()
	}
}
