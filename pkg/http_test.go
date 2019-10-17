package serverfullgw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"testing"
	"text/template"

	"github.com/golang/mock/gomock"
)

func newTemplate(pattern string) *template.Template {
	return template.Must(template.New("").Delims("#!", "!#").Parse(pattern))
}

func TestLambdaTransport_RoundTrip(t *testing.T) {
	var lambdaResponse = `{"v": "R"}` // The response body for all lambda calls in this test.
	type fields struct {
		Name                    string
		Async                   bool
		RequestTemplate         *template.Template
		ResponseSuccessTemplate *template.Template
		ResponseErrorTemplate   *template.Template
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		want              *http.Response
		wantErr           bool
		wantCallWrapped   bool
		wantErrorResponse bool
		wantSign          bool
		wantSignError     bool
	}{
		{
			name: "success",
			fields: fields{
				Name:            "success",
				Async:           false,
				RequestTemplate: newTemplate(`{"v": "#!.Request.Header.T!#"}`),
				ResponseSuccessTemplate: newTemplate(
					`{"status": 200, "header": {}, "bodyPassthrough": false, "body": {"v2": "#!.Response.Body.v!#"}}`,
				),
				ResponseErrorTemplate: newTemplate(
					`{"status": 500, "header": {}, "bodyPassthrough": false, "body": {}}`,
				),
			},
			args: args{
				req: newRequestBody(MultiMap{"T": []string{"V"}}, MultiMap{}, `{}`),
			},
			want: &http.Response{
				StatusCode: 200,
				Status:     http.StatusText(200),
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"v2": "R"}`)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			},
			wantErr:           false,
			wantCallWrapped:   true,
			wantErrorResponse: false,
			wantSign:          true,
			wantSignError:     false,
		},
		{
			name: "success async",
			fields: fields{
				Name:            "success",
				Async:           true,
				RequestTemplate: newTemplate(`{"v": "#!.Request.Header.T!#"}`),
				ResponseSuccessTemplate: newTemplate(
					`{"status": 202, "header": {}, "bodyPassthrough": false}`,
				),
				ResponseErrorTemplate: newTemplate(
					`{"status": 500, "header": {}, "bodyPassthrough": false, "body": {}}`,
				),
			},
			args: args{
				req: newRequestBody(MultiMap{"T": []string{"V"}}, MultiMap{}, `{}`),
			},
			want: &http.Response{
				StatusCode: 202,
				Status:     http.StatusText(202),
				Body:       http.NoBody,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			},
			wantErr:           false,
			wantCallWrapped:   true,
			wantErrorResponse: false,
			wantSign:          true,
			wantSignError:     false,
		},
		{
			name: "success passthrough",
			fields: fields{
				Name:            "success",
				Async:           false,
				RequestTemplate: newTemplate(`{"v": "#!.Request.Header.T!#"}`),
				ResponseSuccessTemplate: newTemplate(
					`{"status": 200, "header": {}, "bodyPassthrough": true}`,
				),
				ResponseErrorTemplate: newTemplate(
					`{"status": 500, "header": {}, "bodyPassthrough": false, "body": {}}`,
				),
			},
			args: args{
				req: newRequestBody(MultiMap{"T": []string{"V"}}, MultiMap{}, `{}`),
			},
			want: &http.Response{
				StatusCode: 200,
				Status:     http.StatusText(200),
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"v":"R"}`)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			},
			wantErr:           false,
			wantCallWrapped:   true,
			wantErrorResponse: false,
			wantSign:          true,
			wantSignError:     false,
		},
		{
			name: "body not json",
			fields: fields{
				Name:            "success",
				Async:           false,
				RequestTemplate: newTemplate(`{"v": "#!.Request.Header.T!#"}`),
				ResponseSuccessTemplate: newTemplate(
					`{"status": 200, "header": {}, "bodyPassthrough": false, "body": {}}`,
				),
				ResponseErrorTemplate: newTemplate(
					`{"status": 500, "header": {}, "bodyPassthrough": false, "body": {}}`,
				),
			},
			args: args{
				req: newRequestBody(MultiMap{"T": []string{"V"}}, MultiMap{}, `o`),
			},
			want: &http.Response{
				StatusCode: 400,
				Status:     http.StatusText(400),
				Body: ioutil.NopCloser(bytes.NewBufferString(
					`{"code":400,"status":"Bad Request"}`,
				)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			wantErr:           false,
			wantCallWrapped:   false,
			wantErrorResponse: true,
			wantSign:          false,
			wantSignError:     false,
		},
		{
			name: "bad request template",
			fields: fields{
				Name:            "success",
				Async:           false,
				RequestTemplate: newTemplate(`{"v": "#!.Request.sHeader.T!#"}`),
				ResponseSuccessTemplate: newTemplate(
					`{"status": 200, "header": {}, "bodyPassthrough": false, "body": {}}`,
				),
				ResponseErrorTemplate: newTemplate(
					`{"status": 500, "header": {}, "bodyPassthrough": false, "body": {}}`,
				),
			},
			args: args{
				req: newRequestBody(MultiMap{"T": []string{"V"}}, MultiMap{}, `{}`),
			},
			want: &http.Response{
				StatusCode: 500,
				Status:     http.StatusText(500),
				Body: ioutil.NopCloser(bytes.NewBufferString(
					`{"code":500,"status":"Internal Server Error"}`,
				)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			wantErr:           false,
			wantCallWrapped:   false,
			wantErrorResponse: true,
			wantSign:          false,
			wantSignError:     false,
		},
		{
			name: "bad success template",
			fields: fields{
				Name:            "success",
				Async:           false,
				RequestTemplate: newTemplate(`{"v": "#!.Request.Header.T!#"}`),
				ResponseSuccessTemplate: newTemplate(
					`{"status": 200, "header": {}, "bodyPassthrough": false, "body": {"v": #!.Response.sBody.v!#}}`,
				),
				ResponseErrorTemplate: newTemplate(
					`{"status": 500, "header": {}, "bodyPassthrough": false, "body": {}}`,
				),
			},
			args: args{
				req: newRequestBody(MultiMap{"T": []string{"V"}}, MultiMap{}, `{}`),
			},
			want: &http.Response{
				StatusCode: 500,
				Status:     http.StatusText(500),
				Body: ioutil.NopCloser(bytes.NewBufferString(
					`{"code":500,"status":"Internal Server Error"}`,
				)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			wantErr:           false,
			wantCallWrapped:   true,
			wantErrorResponse: true,
			wantSign:          true,
			wantSignError:     false,
		},
		{
			name: "bad credentials",
			fields: fields{
				Name:            "success",
				Async:           false,
				RequestTemplate: newTemplate(`{"v": "#!.Request.Header.T!#"}`),
				ResponseSuccessTemplate: newTemplate(
					`{"status": 200, "header": {}, "bodyPassthrough": false, "body": {"v2": "#!.Response.Body.v!#"}}`,
				),
				ResponseErrorTemplate: newTemplate(
					`{"status": 500, "header": {}, "bodyPassthrough": false, "body": {}}`,
				),
			},
			args: args{
				req: newRequestBody(MultiMap{"T": []string{"V"}}, MultiMap{}, `{}`),
			},
			want: &http.Response{
				StatusCode: 500,
				Status:     http.StatusText(500),
				Body: ioutil.NopCloser(bytes.NewBufferString(
					`{"code":500,"status":"Internal Server Error"}`,
				)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			wantErr:           false,
			wantCallWrapped:   false,
			wantErrorResponse: true,
			wantSign:          true,
			wantSignError:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wrapped := NewMockRoundTripper(ctrl)
			signer := NewMockSigner(ctrl)
			r := &LambdaTransport{
				Wrapped:                 wrapped,
				URLParamFn:              func(context.Context) map[string]string { return make(map[string]string) },
				Name:                    tt.fields.Name,
				Async:                   tt.fields.Async,
				RequestTemplate:         tt.fields.RequestTemplate,
				ResponseSuccessTemplate: tt.fields.ResponseSuccessTemplate,
				ResponseErrorTemplate:   tt.fields.ResponseErrorTemplate,
				Signer:                  signer,
			}
			if tt.wantSign {
				if tt.wantSignError {
					signer.EXPECT().Sign(gomock.Any(), gomock.Any()).Return(fmt.Errorf("sign error"))
				} else {
					signer.EXPECT().Sign(gomock.Any(), gomock.Any()).Return(nil)
				}
			}
			if tt.wantCallWrapped {
				wrapped.EXPECT().RoundTrip(gomock.Any()).Do(
					func(req *http.Request) {
						wantPath := path.Join("/", "2015-03-31", "functions", tt.fields.Name, "invocations")
						if req.URL.Path != wantPath {
							t.Errorf("wrong lambda path %s, want %s", req.URL.Path, wantPath)
						}
						if tt.fields.Async {
							if req.Header.Get("X-Amz-Invocation-Type") != "Event" {
								t.Errorf("did not send correct invocation type header %s", req.Header.Get("X-Amz-Invocation-Type"))
							}
						}
					},
				).Return(&http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(lambdaResponse)),
					// Setting ContentLength to zero in order to test for cases where the body
					// is non-zero but the length is not reported. A Go HTTP server usually
					// writes the content length for responses but only if the body is under
					// a certain size and there are no calls to Flush(). Larger payloads, then,
					// result in a missing content-length value. Previous versions of the project
					// checked the ContentLength attributes in a switch which failed when the
					// body was populated but mis-reported.
					ContentLength: 0,
				}, nil)
			}

			got, err := r.RoundTrip(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("LambdaTransport.RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.StatusCode != tt.want.StatusCode {
					t.Errorf("LambdaTransport.RoundTrip().Status = %v, want %v", got.StatusCode, tt.want.StatusCode)
				}
				if !reflect.DeepEqual(got.Header, tt.want.Header) {
					t.Errorf("LambdaTransport.RoundTrip().Status = %v, want %v", got.Header, tt.want.Header)
				}
				gotB, _ := ioutil.ReadAll(got.Body)
				wantB, _ := ioutil.ReadAll(tt.want.Body)
				if !tt.wantErrorResponse {
					if !bytes.Equal(gotB, wantB) {
						t.Errorf("LambdaTransport.RoundTrip().Body = %v, want %v", string(gotB), string(wantB))
					}
				}
				if tt.wantErrorResponse {
					var gotE httpError
					var wantE httpError
					_ = json.Unmarshal(gotB, &gotE)
					_ = json.Unmarshal(wantB, &wantE)
					if gotE.Status != wantE.Status {
						t.Errorf("LambdaTransport.RoundTrip().Body.Status = %v, want %v", gotE.Status, wantE.Status)
					}
					if gotE.Code != wantE.Code {
						t.Errorf("LambdaTransport.RoundTrip().Body.Code = %v, want %v", gotE.Code, wantE.Code)
					}
					// Skipping check on .Reason because it is depending on various internal error messages.
				}
			}
		})
	}
}
