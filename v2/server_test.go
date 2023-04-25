//go:build v2

package v2

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h := NewHTTPServer()

	h.addRoute(http.MethodGet, "/user", func(ctx *Context) {
		fmt.Println("处理第一件事")
		fmt.Println("处理第二件事")
	})

	handle1 := func(ctx *Context) {
		fmt.Println("处理第三件事")
	}

	handle2 := func(ctx *Context) {
		fmt.Println("处理第四件事")
	}

	h.addRoute(http.MethodGet, "/handle1", handle1)
	h.addRoute(http.MethodGet, "/handle2", handle2)

	h.addRoute(http.MethodGet, "/order/detail", func(ctx *Context) {
		_, _ = ctx.Resp.Write([]byte("hello order detail"))
	})

	// 通配符匹配
	h.Get("/order/*", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
	})

	err := h.Start(":8087")
	if err != nil {
		return
	}
}
