//go:build v1

package v1

import (
	"testing"
)

func TestServer(t *testing.T) {
	var s Server
	//http.ListenAndServe(":8080", s)

	s.Start(":8081")
}
