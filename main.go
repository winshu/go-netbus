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
	fmt.Println(`usage: "-listen port1 port2" example: "nb -listen 1997 2017" `)
	fmt.Println(`       "-tran port1 ip:port2" example: "nb -tran 1997 192.168.1.2:3389" `)
	fmt.Println(`       "-slave ip1:port1 ip2:port2" example: "nb -slave 127.0.0.1:3389 8.8.8.8:1997" `)
	fmt.Println(`============================================================`)
	fmt.Println(`optional argument: "-log logpath" . example: "nb -listen 1997 2017 -log d:/nb" `)
	fmt.Println(`log filename format: Y_m_d_H_i_s-agrs1-args2-args3.log`)
	fmt.Println(`============================================================`)
	fmt.Println(`if you want more help, please read "README.md". `)
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
