package cache

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/testkit"
	"github.com/novakit/view"
)

type Model struct {
	Key string
	Val string
}

func TestCache_All(t *testing.T) {
	n := nova.New()
	n.Use(Handler(Options{}))
	n.Use(view.Handler())
	router.Route(n).Get("/set/:key/:value").Use(func(c *nova.Context) error {
		cch := Extract(c)
		vals := router.PathParams(c)
		v := view.Extract(c)
		model := Model{Key: vals.Get("key"), Val: vals.Get("value")}
		err := cch.SetJSON(vals.Get("key"), model, 1)
		if err != nil {
			t.Error("failed")
		}
		v.Text("OK")
		return nil
	})
	router.Route(n).Get("/get/:key").Use(func(c *nova.Context) error {
		vals := router.PathParams(c)
		cch := Extract(c)
		v := view.Extract(c)
		s, err := cch.Get(vals.Get("key"))
		if err != nil {
			t.Error("failed")
		}
		v.Text(s)
		return nil
	})
	router.Route(n).Get("/get/:key/ne").Use(func(c *nova.Context) error {
		vals := router.PathParams(c)
		cch := Extract(c)
		v := view.Extract(c)
		s, err := cch.Get(vals.Get("key"))
		if err != ErrKeyNotFound {
			t.Error("failed")
		}
		v.Text(s)
		return nil
	})
	var req *http.Request
	var res *testkit.DummyResponse

	req, _ = http.NewRequest(http.MethodGet, "/set/key1/value1", nil)
	res = testkit.NewDummyResponse()

	n.ServeHTTP(res, req)

	req, _ = http.NewRequest(http.MethodGet, "/get/key1", nil)
	res = testkit.NewDummyResponse()

	n.ServeHTTP(res, req)

	m := &Model{}
	json.Unmarshal([]byte(res.String()), m)

	if m.Key != "key1" || m.Val != "value1" {
		t.Fatal("failed")
	}

	time.Sleep(time.Second * 2)

	req, _ = http.NewRequest(http.MethodGet, "/get/key1/ne", nil)
	res = testkit.NewDummyResponse()

	n.ServeHTTP(res, req)
}
