//go:build v2

package v2

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string
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
