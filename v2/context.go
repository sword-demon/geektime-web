//go:build v2

package v2

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	// cacheQueryValues url.Values 引入URL 查询参数缓存
	cacheQueryValues url.Values
}

func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("web: 输入为 nil")
	}

	if c.Req.Body == nil {
		return errors.New("web: Body 为 nil")
	}

	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	values, ok := c.Req.Form[key]
	if !ok {
		return "", errors.New("web: key不存在")
	}

	// 注意这里的 values 的类型是 []string
	return values[0], nil
}

// QueryValue 获取 url 中的 query 参数解析
func (c *Context) QueryValue(key string) (string, error) {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	values, ok := c.cacheQueryValues[key]
	if !ok {
		return "", errors.New("web: key不存在")
	}
	// 用户区别不出来真的有值，但是值恰好是空字符串还是没有值
	// 每次都 ParseForm 都要重新解析，所以这里直接使用 Get
	// 和表单比起来，它是没有缓存的，所以每次都要解析
	// 避免多次解析, 稍微缓存一下
	return values[0], nil
}

// PathValue 路径参数解析
func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web: key不存在")
	}
	return val, nil
}

func (c *Context) PathValueV1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key不存在"),
		}
	}
	return StringValue{val: val}
}

type StringValue struct {
	val string
	err error
}

// AsInt64 扩展性函数 将字符串转为 int64
func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
