package zweb

import (
	"net/http"
)

// HandlerFunc defines the handler process request used by zweb
//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

// Engine implements the interface of ServerHTTP
type Engine struct {
	//router map[string]HandlerFunc
	router *router
}

// New is the constructor of zweb.Engine
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

// addRoute 添加路由
func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
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
	c := newContext(writer, request)
	e.router.handle(c)
}
