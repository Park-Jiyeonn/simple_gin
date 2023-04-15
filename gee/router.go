package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}
func parsePattern(pattern string) []string {
	sp := strings.Split(pattern, "/")
	ret := make([]string, 0)
	for _, val := range sp {
		if val != "" {
			ret = append(ret, val)
			if val[0] == '*' {
				break
			}
		}
	}
	return ret
}
func (r *router) addRoute(method string, pattern string, handler HandleFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 用户的链接可能是 /hello/Jiyeon，注册的路由是 /hello/:name，这时候要将name映射到Jiyeon
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	httpParts := parsePattern(path)
	f := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	x := root.search(httpParts, 0)

	if x != nil {
		register_parts := parsePattern(x.pattern)
		for i, part := range register_parts {
			if part[0] == ':' {
				f[part[1:]] = httpParts[i] // 实现f[name] = Jiyeon
			}
			if part[0] == '*' {
				// 本来是分离的字符串，现在通过/连接他们
				f[part[1:]] = strings.Join(httpParts[i:], "/")
			}
		}
		return x, f
	}
	return nil, nil
}

// 中间件先走，因为中间件先 append 进来了，中间件同意 继续 Next了，才可以走到路由对应的 function
func (r *router) handle(c *Context) {
	x, f := r.getRoute(c.Method, c.Path)
	if x != nil {
		c.f = f // f里面是参数对应的值
		key := c.Method + "-" + x.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND %s\n", c.Path)
		})
	}
	c.Next()
}
