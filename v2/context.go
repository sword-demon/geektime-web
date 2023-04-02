//go:build v2

package v2

import "net/http"

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
}
