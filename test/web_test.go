package test

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

// 本地模拟多个 web 服务

func handleWebRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("path", r.URL.Path)
	_, _ = fmt.Fprintf(w, r.Host)
}

func listenOnPort(port int) {
	log.Println("Listen to port", port)
	http.HandleFunc("/", handleWebRequest)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func TestWeb7001(t *testing.T) {
	listenOnPort(7001)
}

func TestWeb7002(t *testing.T) {
	listenOnPort(7002)
}
