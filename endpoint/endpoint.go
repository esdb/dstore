package endpoint

import (
	"github.com/v2pro/plz/counselor"
	"github.com/json-iterator/go"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz.service/http"
	"github.com/esdb/lstore"
	"github.com/v2pro/plz"
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
	server.Handle("/append", AppendFunc(endpoint.append))
	server.Handle("/search", SearchFunc(endpoint.search))
	server.Start(cfg.ListenAddr)
}

type dstoreEndpoint struct {
	store *lstore.Store
}

type AppendRequest struct {
	*lstore.Entry
}

type AppendResponse struct {
	Offset lstore.Offset
}

type AppendFunc func(ctx *countlog.Context, req *AppendRequest) (*AppendResponse, error)

func (endpoint *dstoreEndpoint) append(ctx *countlog.Context, req *AppendRequest) (*AppendResponse, error) {
	offset, err := endpoint.store.Append(ctx, req.Entry)
	if err != nil {
		return nil, err
	}
	return &AppendResponse{offset}, nil
}

type SearchRequest struct {
	Query string
}

type SearchResponse struct {
	Rows []lstore.Row
}

type SearchFunc func(ctx *countlog.Context, req *SearchRequest) (*SearchResponse, error)

func (endpoint *dstoreEndpoint) search(ctx *countlog.Context, req *SearchRequest) (*SearchResponse, error) {
	reader, err := endpoint.store.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer plz.Close(reader)
	reader.SearchForward(ctx, &lstore.SearchRequest{

	})
	return nil, nil
}
