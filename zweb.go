package zweb

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc defines the handler process request used by zweb
//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

// 分组的意义
// 分组控制(Group Control)是 Web 框架应提供的基础功能之一。所谓分组，是指路由的分组。如果没有路由分组，我们需要针对每一个路由进行控制。但是真实的业务场景中，往往某一组路由需要相似的处理。
// 以/post开头的路由匿名可访问。 以/admin开头的路由需要鉴权 以/api开头的路由是 RESTFUL 接口，可以对接第三方平台，需要三方平台鉴权。
// /post是一个分组，/post/a和/post/b可以是该分组下的子分组。作用在/post分组上的中间件(middleware)，也都会作用在子分组，子分组还可以应用自己特有的中间件。
// 中间件可以给框架提供无限的扩展能力，应用在分组上，可以使得分组控制的收益更为明显，而不是共享相同的路由前缀这么简单。例如/admin的分组，可以应用鉴权中间件；
///分组应用日志中间件，/是默认的最顶层的分组，也就意味着给所有的路由，即整个框架增加了记录日志的能力。

// 一个 Group 对象需要具备哪些属性呢？首先是前缀(prefix)，比如/，或者/api；要支持分组嵌套，那么需要知道当前分组的父亲(parent)是谁；当然了，按照我们一开始的分析，中间件是应用在分组上的，那还需要存储应用在该分组上的中间件(middlewares)
// 如果Group对象需要直接映射路由规则 Group对象，还需要有访问Router的能力，为了方便，我们可以在Group中，保存一个指针，指向Engine，整个框架的所有资源都是由Engine统一协调的，那么就可以通过Engine间接地访问各种接口了。

// RouterGroup 分组控制
type RouterGroup struct {
	prefix      string        // 前缀(prefix)，比如/，或者/api
	middlewares []HandlerFunc // 支持的中间间
	//parent      *RouterGroup  // 要支持分组嵌套
	engine *Engine // all groups share a Engine instance
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		//parent: rg,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 由于Engine从某种意义上继承了RouterGroup的所有属性和方法，因为 (*Engine).engine 是指向自己的
//。这样实现，我们既可以通过Engine添加路由，也可以通过分组添加路由。
func (rg *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := rg.prefix + comp
	//log.Printf("Route %4s - %s", method, pattern)
	rg.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (rg *RouterGroup) GET(pattern string, handler HandlerFunc) {
	rg.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (rg *RouterGroup) POST(pattern string, handler HandlerFunc) {
	rg.addRoute("POST", pattern, handler)
}

// Engine implements the interface of ServerHTTP
type Engine struct {
	// 进一步地抽象，将Engine作为最顶层的分组，也就是说Engine拥有RouterGroup所有的能力。
	*RouterGroup
	//router map[string]HandlerFunc
	router        *router
	groups        []*RouterGroup     //store all groups
	htmlTemplates *template.Template // HTML 模板渲染 将所有的模板加载进内存
	funcMap       template.FuncMap   // HTML 模板渲染 定义模板渲染函数
}

// New is the constructor of zweb.Engine
func New() *Engine {

	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
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

func (rg *RouterGroup) Use(middleware ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middleware...)
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(writer, request)
	c.handlers = middlewares
	c.engine = e
	e.router.handle(c)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

// 服务端渲染
func (rg *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(rg.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		// 		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

func (rg *RouterGroup) Static(relativePath, root string) {
	handler := rg.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	rg.GET(urlPattern, handler)
}
