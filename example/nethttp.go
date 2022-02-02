package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	// :10000表示在 10000 端口监听。
	// 第二个参数则代表处理所有的HTTP请求的实例，nil 代表使用标准库中的实例处理。
	// 第二个参数，也是是基于net/http标准库实现Web框架的入口。
	log.Fatal(http.ListenAndServe("127.0.0.1:10000", nil))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprintf(w, "URL.Path: %s\n", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		_, _ = fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
