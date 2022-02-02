package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine is the uni handler for all requests
type Engine struct{}

// 基于net/http标准库实现Web框架的入口。
/*type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}*/

// 拦截了所有的HTTP请求，
// 在这里我们可以自由定义路由映射的规则，
//也可以统一添加一些处理逻辑，例如日志、异常处理等。

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		_, _ = fmt.Fprintf(w, "URL.Path: %s\n", r.URL.Path)
	case "/hello":
		for k, v := range r.Header {
			_, _ = fmt.Fprintf(w, "Header: %s = %s\n", k, v)
		}
	default:
		_, _ = fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":10001", engine))
}
