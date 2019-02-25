package serverfullgw

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	transportd "github.com/asecurityteam/transportd/pkg"
)

var (
	fns = template.FuncMap{
		"json": func(v interface{}) (string, error) {
			b, err := json.Marshal(v)
			return string(b), err
		},
	}
)

// Config contains all settings for the lambda plugin.
type Config struct {
	ARN     string `description:"Name/ARN of the lambda to invoke."`
	Async   bool   `description:"Fire and forget rather than wait for a response."`
	Request string `description:"Template string to transform incoming requests to Lambda requests."`
	Success string `description:"Template string to transform a success response into a proxy response."`
	Error   string `description:"Template string to transform a Lambda error response into a proxy response."`
}

// Name of the config root.
func (*Config) Name() string {
	return "lambda"
}

// Component implements the settings.Component interface.
type Component struct{}

// Lambda satisfies the transportd.NeComponent signature.
func Lambda(_ context.Context, _ string, _ string, _ string) (interface{}, error) {
	return &Component{}, nil
}

// Settings returns a config with defaults set.
func (*Component) Settings() *Config {
	return &Config{}
}

// New creates the middleware.
func (*Component) New(_ context.Context, conf *Config) (func(http.RoundTripper) http.RoundTripper, error) {
	rT, err := template.New("request").Funcs(fns).Delims("#!", "!#").Parse(conf.Request)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request template: %s", err.Error())
	}
	sT, err := template.New("success").Funcs(fns).Delims("#!", "!#").Parse(conf.Success)
	if err != nil {
		return nil, fmt.Errorf("failed to parse success template: %s", err.Error())
	}
	eT, err := template.New("error").Funcs(fns).Delims("#!", "!#").Parse(conf.Error)
	if err != nil {
		return nil, fmt.Errorf("failed to parse error template: %s", err.Error())
	}
	return func(wrapped http.RoundTripper) http.RoundTripper {
		return &LambdaTransport{
			Wrapped:                 wrapped,
			URLParamFn:              transportd.PathParamsFromContext,
			Name:                    conf.ARN,
			Async:                   conf.Async,
			RequestTemplate:         rT,
			ResponseSuccessTemplate: sT,
			ResponseErrorTemplate:   eT,
		}
	}, nil
}
