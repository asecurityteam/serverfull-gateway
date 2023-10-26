package serverfullgw

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/asecurityteam/runhttp"
)

var (
	emptyRequest  = Request{}
	emptyResponse = Response{}
)

// MultiMap adds a helper for extracting the first element
// of a slice of values. This is a stand-in for the url.Value
// and http.Query types which are map[string][]string.
type MultiMap map[string][]string

// Get the first matching value or an empty string.
func (m MultiMap) Get(name string) string {
	if m == nil {
		return ""
	}
	vs := m[name]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Request is a container for all available HTTP request values
// in a template.
type Request struct {
	Query  MultiMap
	Header MultiMap
	URL    map[string]string
	Body   map[string]interface{}
}

// Response is a container for all available response values
// in a template.
type Response struct {
	Header MultiMap
	Body   map[string]interface{}
	Status int
}

// TemplateContext is the value given as the root object context
// when rendering a template.
type TemplateContext struct {
	Request Request
	// Response is only populated when rendering a response template
	// and contains the unmarshalled JSON from the Lambda response.
	Response Response
}

// NewResponse converts an http.Response into a template Response.
func NewResponse(ctx context.Context, r *http.Response) (Response, error) {
	logger := runhttp.LoggerFromContext(ctx)
	d := make(map[string]interface{})
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(struct {
			Message string `logevent:"message,default=response-body-exception"`
			Reason  string `logevent:"reason"`
		}{
			Reason: err.Error(),
		})
		return emptyResponse, err
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, &d)
		if err != nil {
			reason := fmt.Sprintf("response body is not valid JSON: %s, error: %s", string(b), err.Error())
			logger.Error(struct {
				Message string `logevent:"message,default=response-json-unmarshal-exception"`
				Reason  string `logevent:"reason"`
			}{
				Reason: reason,
			})
			return emptyResponse, err
		}
	}
	return Response{
		Body:   d,
		Header: MultiMap(r.Header),
		Status: r.StatusCode,
	}, nil
}

// NewRequest converts an http.Request into a template Request.
func NewRequest(urlParamFn func(context.Context) map[string]string, r *http.Request) (Request, error) {
	logger := runhttp.LoggerFromContext(r.Context())
	d := make(map[string]interface{})
	var b []byte
	var err error
	if r.Body != nil {
		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error(struct {
				Message string `logevent:"message,default=request-body-exception"`
				Reason  string `logevent:"reason"`
			}{
				Reason: err.Error(),
			})
			return emptyRequest, err
		}
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, &d)
		if err != nil {
			logger.Error(struct {
				Message string `logevent:"message,default=request-json-unmarshal-exception"`
				Reason  string `logevent:"reason"`
			}{
				Reason: err.Error(),
			})
			return emptyRequest, err
		}
	}
	return Request{
		Query:  MultiMap(r.URL.Query()),
		Header: MultiMap(r.Header),
		URL:    urlParamFn(r.Context()),
		Body:   d,
	}, nil
}
