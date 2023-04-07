package gee

import (
	"log"
	"net/http"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc
	parent      *RouterGroup
	engine      *Engine
}

// 具体的哈希表，通过路径来映射
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

// 有点像抽象函数，这个由 main 定义具体的逻辑，通过 GET 和 POST 映射进路由表
type HandleFunc func(*Context)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 创建新的 group
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 小写，就私有了
func (group *RouterGroup) addRoute(method, pattern1 string, handler HandleFunc) {
	pattern := group.prefix + pattern1
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRoute("POST", pattern, handler)
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
