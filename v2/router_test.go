//go:build v2

package v2

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"

	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	// 第一个步骤是构造路由树
	// 第二个步骤是验证路由树
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail/:id",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},

		// 对以下内容加校验，不能支持用户这么去写
		// login///
		// /login///
		// //login///
		// /login/a/a//a
	}

	var mockHandler HandleFunc = func(ctx *Context) {

	}
	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 在这里断言路由树和你预期的一模一样
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": {
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
								paramChild: &node{
									path:    ":id",
									handler: mockHandler,
								},
							},
						},
						startChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
				},
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"create": {
								path:    "create",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}

	// 断言两者相等， 不能使用 assert.Equal
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	r = newRouter()
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	}, "web: 路径必须以 / 开头")

	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	}, "web: 路径不能以 / 结尾")

	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/a////", mockHandler)
	}, "web: 不能有连续的 / ")

	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "web: 路由冲突，重复注册[/]")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "web: 路由冲突，重复注册[/a/b/c]")

	// 可用的几个 http method 校验
	// 将 AddRoute 改为私有， 直接使用自己写的 Get 等方法去默认提供给用户使用，所以不需要校验
	// mockHandler 为 nil 的情况 需不需要校验 不需要校验，相当于不会注册

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/*", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	}, "web: 不允许同时注册路径参数和通配符匹配, 已有通配符匹配")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/*", mockHandler)
	}, "web: 不允许同时注册路径参数和通配符匹配, 已有路径参数匹配")
}

func (r *router) equal(y router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的 http method"), false
		}
		// v dst 要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	// 比较长度
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不匹配"), false
	}

	// 比较 startChild
	if n.startChild != nil {
		// 严格判断还需要判断 y.startChild 是否是 nil
		msg, ok := n.startChild.equal(y.startChild)
		if !ok {
			return msg, ok
		}
	}

	// 比较 paramChild
	if n.paramChild != nil {
		// 严格判断还需要判断 y.startChild 是否是 nil
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
	}

	// 比较 handler
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("handler 不相等"), false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	testRoute := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodDelete,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/login/:username",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		//{
		//	method: http.MethodPost,
		//	path:   "/*",
		//},
		//{
		//	method: http.MethodPost,
		//	path:   "/*/*",
		//},
		//{
		//	method: http.MethodPost,
		//	path:   "/*/abc",
		//},
	}
	var mockHandler HandleFunc = func(ctx *Context) {

	}
	r := newRouter()
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		wantNode  *node
	}{
		{
			// 方法都不存在
			name:   "method not found",
			method: http.MethodOptions,
			path:   "/order/detail",
		},
		{
			// 完全命中
			name:      "order detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantFound: true,
			wantNode: &node{
				handler: mockHandler,
				path:    "detail",
			},
		},
		{
			// 完全命中
			name:      "order start",
			method:    http.MethodGet,
			path:      "/order/abc",
			wantFound: true,
			wantNode: &node{
				handler: mockHandler,
				path:    "*",
			},
		},
		{
			// 命中了但是没有 handler
			name:      "order",
			method:    http.MethodGet,
			path:      "/order",
			wantFound: true,
			wantNode: &node{
				path: "order",
				children: map[string]*node{
					"detail": {
						path:    "detail",
						handler: mockHandler,
					},
				},
			},
		},
		{
			// 根节点
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
			},
		},
		{
			// path nod found
			name:   "path not found",
			method: http.MethodGet,
			path:   "/aaaabbbccc",
		},
		{
			// :username 路径参数匹配
			name:      "login username",
			method:    http.MethodPost,
			path:      "login/wujie",
			wantFound: true,
			wantNode: &node{
				path:    ":username",
				handler: mockHandler,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			msg, ok := tc.wantNode.equal(n)
			assert.True(t, ok, msg)
		})
	}
}
