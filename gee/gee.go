package gee

import (
	"net/http"
)

// 有点像抽象函数，这个由 main 定义具体的逻辑，通过 GET 和 POST 映射进路由表
type HandleFunc func(*Context)

// 具体的哈希表，通过路径来映射
type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

// 小写，就私有了
func (engine *Engine) addRoute(method, pattern string, handler HandleFunc) {
	engine.router.addRoute(method, pattern, handler)
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
	c := NewContext(w, req)
	engine.router.handle(c)
}
