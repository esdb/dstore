package test

import (
	"testing"
	"github.com/v2pro/plz/counselor"
	"github.com/esdb/dstore/endpoint/http"
)

func Test_http(t *testing.T) {
	counselor.SetObject("dstore", "http",
		"DstoreHttpConfig", []byte(`{
		"ListenAddr": "127.0.0.1:9776"
	}`))
	http.StartHttpEndpoints()
}
