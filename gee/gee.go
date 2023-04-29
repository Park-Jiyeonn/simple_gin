package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

// 具体的哈希表，通过路径来映射
type Engine struct {
	*RouterGroup  // 为了能够访问 Group 的方法，但是 Group 不可以访问 Engine 的方法。
	router        *router
	groups        []*RouterGroup     // 由 Engine 创建所有的 Group
	htmlTemplates *template.Template // 将所有的模板加载进内存
	funcMap       template.FuncMap   // 所有的自定义模板渲染函数
}

type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc
	parent      *RouterGroup
	engine      *Engine //所有的 Group 都使用这个 Engine
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
	//log.Printf("Route %4s - %s", method, pattern)
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

// 接收任意数量的 midllewares
func (group *RouterGroup) Use(midllewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, midllewares...)
}

// 得实现一下 http 接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var midllewares []HandleFunc
	for _, group := range engine.groups {
		// 如果请求中有这个分组的前缀，那么这个分组的中间件就要加进 Context 的中间件
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			midllewares = append(midllewares, group.middlewares...)
			//fmt.Println("It should be here!!!!!!")
		}
	}
	c := NewContext(w, req)
	c.handlers = midllewares
	c.engine = engine
	engine.router.handle(c)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	// FileServer 会创建一个 http.Handler 的值，用 fs 这个文件系统来提供服务
	// StripPrefix 将请求中的 url 前缀去掉，从而能够得到文件的路径
	// 看了看源码，StripPrefix 就是对请求连接进行了一个修改，然后返回 FileServer 的Handler
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		// 将文件的内容作为 http 相应发送回客户端
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Template
// ---------------------------------------------------------------------------------------------------------

// 把 relativePath 映射到 root
func (group *RouterGroup) Static(relativePath string, root string) {
	// dir 函数将路径暴露给 web服务器，客户端可以用http协议来访问指定目录下的文件和子目录，返回一个文件系统
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPath := path.Join(relativePath, "/*filepath")
	group.GET(urlPath, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// -------------------------------------------------------------------------------------------------------

func Default() *Engine {
	enging := New()
	enging.Use(Logger(), Recovery())
	return enging
}
