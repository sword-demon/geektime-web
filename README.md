# go web 框架学习

## Request Query

```go
func queryParams(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fmt.Fprintf(w, "query is %v\n", values)
}
```
- 除了`body`，我们可以传参的地方是`Query`
- 所有的值都被解释为字符串，所以需要自己解析为数字等


## Request Header
- Header 大体是两类，一种是 HTTP 预定义的，一类是自己定义的
- Go 会自动降 Header 名字转为标准名字 -- 其实就是大小写调整，每个单词首字母大写
- 一般用`X`开头来表明是自己定义的，比如：`X-mycompany-your=header`

## Form 表单
- Form 和 ParseForm
- 要先调用 ParseForm，不调用拿到的 Form 则是`nil`
- 建议加上`Content-Type: application/x-www-form-urlencoded`

## Beego Context 设计

- Input：对输入的封装
- Output：对输出的封装
- Response：对响应的封装

```go
type Context struct {
	Input *BeegoInput
	Output *BeegoOutput
	Request *http.Request
	ResponseWriter *Response
	_xsrfToken string
}
```

```go
type BeegoInput struct {
	Context *Context
	CruSession session.Store
	pnames []string
	pvalues []string
	data map[interface{}]interface{}
	dataLock sync.RWMutext
	RequestBody []byte
	RunMethod string
	RunController refect.Type
}
```

- 反向应用了 Context
- 直接耦合了 Session，Beego 直接内置了 Session 的支持
- 维持了不同部门的输入

```go
type BeegoOutput struct {
    Context *Context
	Status int
	EnableGzip bool
}
```

```go
type Response struct {
	http.ResponseWriter
	Started bool
	Status int
	Elapsed time.Duration
}
```

- 反向引用了`Context`
- 维持住了`Status`，即 HTTP 状态响应码

### Beego 处理输入的方法
- Bind 一族方法，用于将 Body 转化为具体的结构体
- Input 中的 Bind 方法，尝试将各个部分的输入都绑定到一个结构体里
- Input 尝试从各个部位获取输入的方法
- Input 中还有各种判断的方法，判断是否是`Get`、`Post`请求的方法等

### Beego 处理输出的方法
- Resp 一族方法，用于将输入序列化后输出
- Render 方法尝试渲染模板，输出响应
- Output 中定义的输出各种格式数据的方法