package test

import (
	"../core"
	"fmt"
	"log"
	"net/http"
	"testing"
)

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

func TestToByte(t *testing.T) {
	arr := []string{"127.0.0.1:3306", "255.255.255.255:65535", "我是中国人:333"}
	for _, v := range arr {
		n := core.FormatHeader(v)
		fmt.Println([]byte(n), len([]byte(n)))
	}
}
