// +build integration

package tests

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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
	reqs := make(chan *http.Request, 100)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqs <- r
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"v":"R"}`))
	}))
	defer srv.Close()
	os.Setenv("TEST_HOST", srv.URL)
	port, err := getPort()
	require.Nil(t, err)
	os.Setenv("TEST_SERVER_PORT", port)

	done := make(chan error, 1)
	rt, err := transportd.New(context.Background(), spec, serverfullgw.Lambda)
	require.Nil(t, err)
	rt.Exit = done
	go func() {
		if runErr := rt.Run(); runErr != nil {
			t.Log(runErr)
		}
	}()
	stop := time.Now().Add(time.Second)
	var resp *http.Response
	var req *http.Request
	for time.Now().Before(stop) {
		req, _ = http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port, http.NoBody)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Log(err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		defer resp.Body.Close()
		break
	}
	require.Nil(t, err)

	respB, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, []byte(`{"v2":"R"}`), respB)
	done <- nil

	req = <-reqs
	assert.Equal(t, "/2015-03-31/functions/test/invocations", req.URL.Path)
}
