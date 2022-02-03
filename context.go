package zweb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
提供了访问Query和PostForm参数的方法。
提供了快速构造String/Data/JSON/HTML响应的方法.
*/

// H map[string]interface{} H
type H map[string]interface{}

// Context  封装*http.Request和http.ResponseWriter的方法，简化相关接口的调用
//  扩展性和复杂性留在了内部 对外简化了接口。
type Context struct {
	// origin info
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int

	//middlewares
	handlers []HandlerFunc
	index    int // index是记录当前执行到第几个中间件
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

func (ctx *Context) Next() {
	ctx.index++
	s := len(ctx.handlers)
	for ; ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

func (ctx *Context) Fail(code int, err string) {
	ctx.index = len(ctx.handlers)
	ctx.JSON(code, H{"message": err})
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}

func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	_, _ = ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "applications/json;charset=utf-8")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		// http.Error(c.Writer, err.Error(), 500)这里不起作用
		// 前面已经执行了 ctx.Status(code) 也就是调用了 WriteHeader(code)  w.Header().Set()不起作用 ,返回码将不会再更改 所以直接panic
		// WriteHeader必须在Write之前调用

		// 正确的顺序
		//  w.Header().Set("xx")
		//  w.WriteHeader(code)
		//  w.Write([]byte("hello world\n"))

		// https://www.zhangshengrong.com/p/zD1yDr25Xr/
		panic("Encode failed")
	}
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	_, _ = ctx.Writer.Write(data)
}

func (ctx *Context) HTML(code int, html string) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	_, _ = ctx.Writer.Write([]byte(html))
}

func (ctx *Context) Param(key string) string {
	value, _ := ctx.Params[key]
	return value
}
