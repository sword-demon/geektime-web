//go:build v2

package v2

import (
	"fmt"
	"strings"
)

type router struct {
	// trees 按照 HTTP 方法来阻止
	trees map[string]*node
}

func newRouter() router {
	return router{trees: map[string]*node{}}
}

// addRoute 添加路由
// 加一些限制:
// path 必须以 "/" 开头 不能以 "/" 结尾，中间也不能有连续的 "//"
func (r *router) addRoute(method string, path string, handlerFunc HandleFunc) {

	if path == "" {
		panic("web: 路径不能为空字符串")
	}

	// 开头不能没有/
	if path[0] != '/' {
		panic("web: 路径必须以 / 开头")
	}

	// 结尾
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路径不能以 / 结尾")
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

	// 如果是根节点
	// 根节点特殊处理
	if path == "/" {
		// 根节点重复注册
		if root.handler != nil {
			panic("web: 路由冲突，重复注册[/]")
		}
		root.handler = handlerFunc
		return
	}

	// 去除 path 的第一个空字符串防止切割有空字符串
	//path = path[1:]
	// 切割 path
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		// 中间连续 ///
		if seg == "" {
			panic("web: 不能有连续的 / ")
		}
		// 递归下去找准位置
		// 如果中途有节点不存在，就要创建出来
		children := root.childOrCreate(seg)
		root = children
	}

	// 覆盖之前检测是否已注册
	if root.handler != nil {
		// 这里需要原始的 path 变量， 所以上面分割直接使用 path[1]
		panic(fmt.Sprintf("web: 路由冲突，重复注册[%s]", path))
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
