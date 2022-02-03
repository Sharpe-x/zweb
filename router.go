package zweb

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	// // roots key eg, roots['GET'] roots['POST']
	roots map[string]*node // 存储每种请求方式的Trie 树根节点。
	// // handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
	handlers map[string]HandlerFunc //  存储每种请求方式的 HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// // Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)

	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRouter /assets/program/go/base
func (r *router) getRouter(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // assets program go base
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {

	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key]) // 将从路由匹配得到的 Handler 添加到 c.handlers列表中
		// 中间件A B
		/*func A(c *Context) {
			part1
			c.Next()
			part2
		}
		func B(c *Context) {
			part3
			c.Next()
			part4
		}*/
		// 最终执行顺序 part1 -> part3 -> r.handlers[key] -> part 4 -> part2

	} else {
		c.handlers = append(c.handlers, func(context *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
