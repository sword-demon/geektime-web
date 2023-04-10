//go:build v1

package v1

import "net/http"

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr: 监听的地址， 如果只指定端口: 可以使用 ":8081"
	// 或者 "localhost:8081"
	Start(addr string) error
	// addRoute 注册一个路由
	// method 是 HTTP 方法
	// path 路径 必须以 / 开头
	AddRoute(method string, path string, handler HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
}

func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:  r,
		Resp: w,
	}
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {

}

func (h *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, h)
}

func (h *HTTPServer) AddRoute(method string, path string, handler HandleFunc) {
	//TODO implement me
	panic("implement me")
}

// Get 实现 get 请求方法
func (h *HTTPServer) Get(path string, handler HandleFunc) {
	h.AddRoute(http.MethodGet, path, handler)
}
