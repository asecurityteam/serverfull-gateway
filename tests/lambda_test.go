// +build integration

package tests

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	serverfullgw "github.com/asecurityteam/serverfull-gateway/pkg"
	transportd "github.com/asecurityteam/transportd/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLambda(t *testing.T) {
	f, err := os.Open("specs/complete.yaml")
	assert.Nil(t, err)
	defer f.Close()
	spec, _ := ioutil.ReadAll(f)
	reqs := make(chan *http.Request, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqs <- r
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"v":"R"}`))
	}))
	defer srv.Close()
	os.Setenv("TEST_HOST", srv.URL)

	done := make(chan error)
	rt, err := transportd.New(context.Background(), spec, serverfullgw.Lambda)
	require.Nil(t, err)
	rt.Exit = done
	go func() { _ = rt.Run() }()

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:9090", http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()

	respB, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte(`{"v2":"R"}`), respB)
	done <- nil

	req = <-reqs
	assert.Equal(t, "/2015-03-31/functions/test/invocations", req.URL.Path)
}
