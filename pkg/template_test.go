package serverfullgw

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestMultiMap_Get(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		m    MultiMap
		args args
		want string
	}{
		{
			name: "nil",
			m:    nil,
			args: args{name: "a"},
			want: "",
		},
		{
			name: "missing",
			m:    MultiMap{},
			args: args{name: "a"},
			want: "",
		},
		{
			name: "multiple values",
			m:    MultiMap{"a": []string{"b", "c"}},
			args: args{name: "a"},
			want: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Get(tt.args.name); got != tt.want {
				t.Errorf("MultiMap.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newRequest(headers MultiMap, query MultiMap) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost/", http.NoBody)
	req.Header = http.Header(headers)
	req.URL.RawQuery = url.Values(query).Encode()
	return req
}
func newRequestBody(headers MultiMap, query MultiMap, body string) *http.Request {
	req := newRequest(headers, query)
	req.Body = io.NopCloser(bytes.NewBufferString(body))
	req.ContentLength = int64(len(body))
	return req
}

func TestNewRequest(t *testing.T) {
	params := make(map[string]string)
	paramsFn := func(context.Context) map[string]string { return params }
	headers := MultiMap{"a": []string{"b", "c"}}
	query := MultiMap{"q": []string{"w", "e"}}
	type args struct {
		urlParamFn func(context.Context) map[string]string
		r          *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    Request
		wantErr bool
	}{
		{
			name: "empty body",
			args: args{
				urlParamFn: paramsFn,
				r:          newRequest(headers, query),
			},
			want: Request{
				Body:   map[string]interface{}{},
				Header: headers,
				URL:    params,
				Query:  query,
			},
			wantErr: false,
		},
		{
			name: "body not json",
			args: args{
				urlParamFn: paramsFn,
				r:          newRequestBody(headers, query, `notjson`),
			},
			want:    Request{},
			wantErr: true,
		},
		{
			name: "body json",
			args: args{
				urlParamFn: paramsFn,
				r:          newRequestBody(headers, query, `{"a": "b"}`),
			},
			want: Request{
				Body: map[string]interface{}{
					"a": "b",
				},
				Header: headers,
				URL:    params,
				Query:  query,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRequest(tt.args.urlParamFn, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResponse(t *testing.T) {
	headers := MultiMap{"a": []string{"b", "c"}}
	type args struct {
		r *http.Response
	}
	tests := []struct {
		name    string
		args    args
		want    Response
		wantErr bool
	}{
		{
			name: "empty body",
			args: args{r: &http.Response{
				ContentLength: 0,
				Header:        http.Header(headers),
				StatusCode:    200,
				Body:          http.NoBody,
			}},
			want: Response{
				Status: 200,
				Header: headers,
				Body:   make(map[string]interface{}),
			},
			wantErr: false,
		},
		{
			name: "body not json",
			args: args{r: &http.Response{
				ContentLength: int64(len(`notjson`)),
				Header:        http.Header(headers),
				StatusCode:    200,
				Body:          io.NopCloser(bytes.NewBufferString(`notjson`)),
			}},
			want:    Response{},
			wantErr: true,
		},
		{
			name: "body json",
			args: args{r: &http.Response{
				ContentLength: int64(len(`{"a": "b"}`)),
				Header:        http.Header(headers),
				StatusCode:    200,
				Body:          io.NopCloser(bytes.NewBufferString(`{"a": "b"}`)),
			}},
			want: Response{
				Status: 200,
				Header: headers,
				Body: map[string]interface{}{
					"a": "b",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewResponse(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
