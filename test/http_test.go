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

func Test_append(t *testing.T) {
	should := require.New(t)
	counselor.SetObject("dstore", "http",
		"DstoreHttpConfig", []byte(`{
		"ListenAddr": "127.0.0.1:9776"
	}`))
	store := testStore(&lstore.Config{})
	go endpoint.StartHttpServer(store)
	time.Sleep(time.Millisecond * 100)
	var appendEntry endpoint.AppendFunc
	client := http.NewClient()
	client.Handle("POST", "http://127.0.0.1:9776/append", &appendEntry)
	resp, err := appendEntry(ctx, &endpoint.AppendRequest{
		Entry: &lstore.Entry{
			IntValues: []int64{1},
		},
	})
	should.NoError(err)
	should.Equal(lstore.Offset(1), resp.Offset)
}


func Test_search(t *testing.T) {
	should := require.New(t)
	counselor.SetObject("dstore", "http",
		"DstoreHttpConfig", []byte(`{
		"ListenAddr": "127.0.0.1:9776"
	}`))
	store := testStore(&lstore.Config{})
	go endpoint.StartHttpServer(store)
	time.Sleep(time.Millisecond * 100)
	var appendEntry endpoint.AppendFunc
	client := http.NewClient()
	client.Handle("POST", "http://127.0.0.1:9776/append", &appendEntry)
	appendEntry(ctx, &endpoint.AppendRequest{
		Entry: &lstore.Entry{
			IntValues: []int64{1},
		},
	})
	var search endpoint.SearchFunc
	client.Handle("POST", "http://127.0.0.1:9776/search", &search)
	resp, err := search(ctx, &endpoint.SearchRequest{
		Query: "SELECT * FROM store WHERE field=1",
	})
	should.NoError(err)
	should.Equal(1, len(resp.Rows))
}