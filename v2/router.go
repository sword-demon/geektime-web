//go:build v2

package v2

import "strings"

type router struct {
	// trees 按照 HTTP 方法来阻止
	trees map[string]*node
}

func newRouter() router {
	return router{trees: map[string]*node{}}
}

func (r *router) addRoute(method string, path string, handlerFunc HandleFunc) {
	// 判断路径是否合法
	if path == "" {
		panic("web: 路由是空字符串")
	}
	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 理由不能以 / 结尾")
	}

	// 首先找到树
	root, ok := r.trees[method]
	if !ok {
		// 说明还没有根节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	// 去除 path 的第一个空字符串防止切割有空字符串
	path = path[1:]
	// 切割 path
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		// 递归下去找准位置
		// 如果中途有节点不存在，就要创建出来
		children := root.childOrCreate(seg)
		root = children
	}

	root.handler = handlerFunc
}

func (n *node) childOrCreate(seg string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		// 要新建一个
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

type node struct {
	path string
	// children 子节点
	// 子节点的 path => node
	children map[string]*node
	// handler 命中路由之后执行的逻辑
	handler HandleFunc
}
