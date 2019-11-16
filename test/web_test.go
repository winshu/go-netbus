package test

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func handleWebRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("path", r.URL.Path)
	_, _ = fmt.Fprintf(w, "Hello Wrold!") //这个写入到w的是输出到客户端的
}
func TestWeb(t *testing.T) {
	log.Println("Listen to port 7001")
	http.HandleFunc("/", handleWebRequest)   //设置访问的路由
	err := http.ListenAndServe(":7001", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
