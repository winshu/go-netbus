package main

import (
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
	fmt.Println(`usage: "-server" start as server`)
	fmt.Println(`       "-client" start as client`)
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

	switch args[1] {
	case "-server":
		serverConfig := core.InitServerConfig()
		core.Server(serverConfig)
	case "-client":
		clientConfig := core.InitClientConfig()
		core.Client(clientConfig)
	default:
		printHelp()
	}
}
