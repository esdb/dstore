package test

import (
	"testing"
	"github.com/v2pro/plz/counselor"
	"github.com/esdb/dstore/endpoint"
	"github.com/v2pro/plz.service/http"
	"time"
	"github.com/v2pro/plz/countlog"
	"context"
	"github.com/esdb/lstore"
	"github.com/stretchr/testify/require"
	"os"
)

var ctx = countlog.Ctx(context.Background())

func testStore(cfg *lstore.Config) *lstore.Store {
	cfg.Directory = "/run/store"
	os.RemoveAll(cfg.Directory)
	store, err := lstore.New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return store
}

func Test_http(t *testing.T) {
	should := require.New(t)
	counselor.SetObject("dstore", "http",
		"DstoreHttpConfig", []byte(`{
		"ListenAddr": "127.0.0.1:9776"
	}`))
	store := testStore(&lstore.Config{})
	go endpoint.StartHttpServer(store)
	time.Sleep(time.Millisecond * 100)
	var appendEntry endpoint.AppendEntryFunc
	client := http.NewClient()
	client.Handle("POST", "http://127.0.0.1:9776/append", &appendEntry)
	resp, err := appendEntry(ctx, &endpoint.AppendEntryRequest{
		Entry: &lstore.Entry{
			IntValues: []int64{1},
		},
	})
	should.NoError(err)
	should.Equal(lstore.Offset(1), resp.Offset)
}
