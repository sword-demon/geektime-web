//go:build v2

package v2

import "net/http"

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	Start(addr string) error
}

var _ Server = &HTTPServer{}

type HTTPServer struct {
	router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{router: newRouter()}
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	h.serve(ctx)
}

func (h *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, h)
}

func (h *HTTPServer) Post(path string, handler HandleFunc) {
	h.addRoute(http.MethodPost, path, handler)
}

func (h *HTTPServer) Get(path string, handler HandleFunc) {
	h.addRoute(http.MethodGet, path, handler)
}

func (h *HTTPServer) serve(ctx *Context) {
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil {
		// 路由没有命中 404
		ctx.Resp.WriteHeader(404)
		_, _ = ctx.Resp.Write([]byte("not found"))
		return
	}
	ctx.PathParams = info.pathParams
	info.n.handler(ctx)
}
