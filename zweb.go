package zweb

import (
	"fmt"
	"net/http"
)

// HandlerFunc defines the handler process request used by zweb
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implements the interface of ServerHTTP
type Engine struct {
	router map[string]HandlerFunc
}

// New is the constructor of zweb.Engine
func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

// addRoute 添加路由
func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handler
}

// GET defines the method to add GET request
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// POST defines the method to add GET request
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := request.Method + "-" + request.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(writer, request)
	} else {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(writer, "404 Not Found: %s\n", request.URL)
	}
}
