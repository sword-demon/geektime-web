//go:build v2

package v2

import "net/http"

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	Stat(addr string) error
	addRoute(method string, path string, handler HandleFunc)
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

func (h *HTTPServer) Stat(addr string) error {
	return http.ListenAndServe(addr, h)
}

func (h *HTTPServer) addRoute(method string, path string, handler HandleFunc) {
	//TODO implement me
	panic("implement me")
}

func (h *HTTPServer) Post(path string, handler HandleFunc) {
	h.addRoute(http.MethodPost, path, handler)
}

func (h *HTTPServer) Get(path string, handler HandleFunc) {
	h.addRoute(http.MethodGet, path, handler)
}

func (h *HTTPServer) serve(ctx *Context) {

}
