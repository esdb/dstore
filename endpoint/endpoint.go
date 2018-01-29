package endpoint

import (
	"github.com/v2pro/plz/counselor"
	"github.com/json-iterator/go"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz.service/http"
	"github.com/esdb/lstore"
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

func StartHttpServer(store *lstore.Store) {
	cfg, found := counselor.GetObject("dstore", "http").(Config)
	if !found {
		countlog.Error("event!missing dstore http config")
		return
	}
	endpoint := &dstoreEndpoint{store}
	server := http.NewServer()
	server.Handle("/append", AppendEntryFunc(endpoint.appendEntry))
	server.Start(cfg.ListenAddr)
}

type dstoreEndpoint struct {
	store *lstore.Store
}

type AppendEntryRequest struct {
	*lstore.Entry
}

type AppendEntryResponse struct {
	Offset lstore.Offset
}

type AppendEntryFunc func(ctx *countlog.Context, req *AppendEntryRequest) (*AppendEntryResponse, error)

func (endpoint *dstoreEndpoint) appendEntry(ctx *countlog.Context, req *AppendEntryRequest) (*AppendEntryResponse, error) {
	offset, err := endpoint.store.Append(ctx, req.Entry)
	if err != nil {
		return nil, err
	}
	return &AppendEntryResponse{offset}, nil
}
