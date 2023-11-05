//go:build v2

package v2

import (
	"fmt"
	"net/http"
	"sync"
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

	h.Post("/form", func(ctx *Context) {
		ctx.Req.ParseForm()
	})

	h.Get("/values/:id", func(ctx *Context) {
		// 使用 StringValue 返回值可以进行链式调用来解析数据
		id, err := ctx.PathValueV1("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("id 输入不对"))
			return
		}

		ctx.Resp.Write([]byte(fmt.Sprintf("hello id: %d", id)))
	})

	type User struct {
		Name string `json:"name"`
	}

	h.Get("/user/:id", func(ctx *Context) {
		ctx.RespJSON(200, User{Name: "张三"})
	})

	err := h.Start(":8087")
	if err != nil {
		return
	}
}

type SafeContext struct {
	Context
	mutex sync.RWMutex
}

func (c *SafeContext) RespJSON1(val any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Context.RespJSONOK(val)
}
