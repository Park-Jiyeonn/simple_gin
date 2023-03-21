package gee

import (
	"fmt"
	"net/http"
)

// 有点像抽象函数，这个由 main 定义具体的逻辑，通过 GET 和 POST 映射进路由表
type HandleFunc func(w http.ResponseWriter, r *http.Request)

// 具体的哈希表，通过路径来映射
type Engine struct {
	router map[string]HandleFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandleFunc)}
}

// 小写，就私有了
func (engine *Engine) addRoute(method, pattern string, handler HandleFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandleFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandleFunc) {
	engine.addRoute("POST", pattern, handler)
}

// main 要调用， 具体的，xx.Run("8081")
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// 得实现一下 http 接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL.Path)
	}
}
