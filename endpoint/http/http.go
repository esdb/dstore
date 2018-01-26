package http

import (
	"github.com/v2pro/plz/counselor"
	"net/http"
	"github.com/json-iterator/go"
	"github.com/v2pro/plz/countlog"
)

type Config struct {
	ListenAddr string
}

func init() {
	counselor.RegisterParserByFunc("DstoreHttpConfig", func(data []byte) (interface{}, error) {
		var cfg Config
		err := jsoniter.Unmarshal(data, &cfg)
		return cfg, err
	})
}

func StartHttpEndpoints() {
	cfg, found := counselor.GetObject("dstore", "http").(Config)
	if !found {
		countlog.Error("event!missing dstore http config")
		return
	}
	http.ListenAndServe(cfg.ListenAddr, nil)
}
