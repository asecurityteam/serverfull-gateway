package serverfullgw

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path"
	"text/template"
)

type httpError struct {
	// Code is the HTTP status code.
	Code int `json:"code"`
	// Status is the HTTP status string.
	Status string `json:"status"`
	// Reason is the debug data.
	Reason string `json:"reason"`
}

func newError(code int, reason string) *http.Response {
	b, _ := json.Marshal(httpError{
		Code:   code,
		Status: http.StatusText(code),
		Reason: reason,
	})
	return &http.Response{
		Status:     http.StatusText(code),
		StatusCode: code,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
}

// TransportResponse is used to unmarshal the rendered response template.
type TransportResponse struct {
	Status      int                 `json:"status"`
	Header      map[string][]string `json:"header"`
	Passthrough bool                `json:"bodyPassthrough"`
	Body        json.RawMessage     `json:"body"`
}

// LambdaTransport is a decorator that rewrites a request into a call
// to the AWS Lambda Invoke API.
type LambdaTransport struct {
	Wrapped                 http.RoundTripper
	URLParamFn              func(context.Context) map[string]string
	Name                    string
	Async                   bool
	RequestTemplate         *template.Template
	ResponseSuccessTemplate *template.Template
	ResponseErrorTemplate   *template.Template
	Signer                  Signer
}

// RoundTrip converts an incoming request to a Lambda Invoke API call.
func (r *LambdaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tReq, err := NewRequest(r.URLParamFn, req)
	if err != nil {
		return newError(http.StatusBadRequest, err.Error()), nil
	}

	var reqBody bytes.Buffer
	err = r.RequestTemplate.Execute(&reqBody, TemplateContext{Request: tReq})
	if err != nil {
		return newError(http.StatusInternalServerError, err.Error()), nil
	}

	req.Method = http.MethodPost
	req.URL.Path = path.Join("/", "2015-03-31", "functions", r.Name, "invocations")
	req.Header = http.Header{}
	req.ContentLength = int64(reqBody.Len())
	req.Body = io.NopCloser(&reqBody)
	if r.Async {
		req.Header.Set("X-Amz-Invocation-Type", "Event")
	}

	if err = r.Signer.Sign(req, bytes.NewReader(reqBody.Bytes())); err != nil {
		return newError(http.StatusInternalServerError, err.Error()), nil
	}
	resp, err := r.Wrapped.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tResp, err := NewResponse(resp)
	if err != nil {
		return newError(http.StatusBadGateway, err.Error()), nil
	}

	tp := r.ResponseSuccessTemplate
	if tResp.Status < 200 || tResp.Status >= 300 {
		tp = r.ResponseErrorTemplate
	}
	var respBody bytes.Buffer
	err = tp.Execute(&respBody, TemplateContext{Request: tReq, Response: tResp})
	if err != nil {
		return newError(http.StatusInternalServerError, err.Error()), nil
	}

	var tr TransportResponse
	if err = json.Unmarshal(respBody.Bytes(), &tr); err != nil {
		return newError(http.StatusInternalServerError, err.Error()), nil
	}
	var finalBody io.ReadCloser = http.NoBody
	if len(tr.Body) > 0 {
		finalBody = io.NopCloser(bytes.NewReader([]byte(tr.Body)))
	}
	if tr.Passthrough && len(tResp.Body) > 0 {
		finalB, _ := json.Marshal(tResp.Body)
		finalBody = io.NopCloser(bytes.NewReader(finalB))
	}
	if tr.Header == nil {
		tr.Header = make(map[string][]string)
	}
	if len(tr.Header["Content-Type"]) < 1 {
		tr.Header["Content-Type"] = append(tr.Header["Content-Type"], "application/json")
	}
	return &http.Response{
		Body:       finalBody,
		Header:     http.Header(tr.Header),
		StatusCode: tr.Status,
		Status:     http.StatusText(tr.Status),
	}, nil
}
