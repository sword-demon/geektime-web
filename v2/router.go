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

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	// 沿着树深度遍历查找下去
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	// 如果是根节点直接返回
	if path == "/" {
		return &matchInfo{n: root}, true
	}
	// 把前置和后置的 / 去掉
	path = strings.Trim(path, "/")
	// 按照 / 切割
	segs := strings.Split(path, "/")
	// 构造 pathParams
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		// 命中了路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			// path 是 :id 这种形式
			pathParams[child.path[1:]] = seg
		}
		root = child
	}

	// 代表我确实有这个节点
	// 但是节点是不是用户注册的业务逻辑 有 handler 的 就不一定了
	//return root, root.handler != nil
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
}

func (n *node) childOrCreate(seg string) *node {
	// 如果用户注册路由的第一段是冒号
	if seg[0] == ':' {
		if n.startChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配, 已有通配符匹配")
		}
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}
	if seg == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配, 已有路径参数匹配")
		}
		n.startChild = &node{
			path: seg,
		}
		return n.startChild
	}
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

// childOf 优先考虑静态匹配 匹配不上再考虑 通配符匹配
// 第一个返回值是子节点
// 第二个是标志是否是路径参数
// 第三个标志命中了没有
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.startChild, false, n.startChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.startChild, false, n.startChild != nil
	}
	return child, false, ok
}

type node struct {
	path string

	// 静态节点
	// children 子节点
	// 子节点的 path => node
	children map[string]*node
	// handler 命中路由之后执行的逻辑
	handler HandleFunc

	// 加一个通配符匹配
	startChild *node

	// 路径参数匹配
	paramChild *node
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
