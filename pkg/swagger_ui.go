package serverfullgw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	transportd "github.com/asecurityteam/transportd/pkg"
	"github.com/getkin/kin-openapi/openapi3"
	packr "github.com/gobuffalo/packr/v2"
)

type swaggerUITransport struct {
	Wrapped         http.RoundTripper
	Spec            *openapi3.Swagger
	SwaggerPathRoot string
}

func (r *swaggerUITransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Find the root path from which swagger should be served.
	// For example, if the spec doc has "swaggerui" enabled under /api/{somevar*}
	// we want to serve swaggerui under /api/*.  Remember transportd will only
	// call through this RoundTrip on an already-best-matched path, so we only
	// need to find how the current request path best matches the spec paths to
	// find the real configured root.
	configuredRootPath := findConfiguredRootPath(r.Spec, req.URL.Path)

	bytesofspec, _ := json.Marshal(r.Spec) // error is safe to ignore here

	if strings.EqualFold(req.URL.Path, configuredRootPath+"/spec-doc") {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer(bytesofspec)), StatusCode: http.StatusOK}, nil
	} else {

		files := packr.New("swaggerui", "../swaggerui")
		filePath := req.URL.Path[len(configuredRootPath)+1 : len(req.URL.Path)]
		fileContents, e := files.Find(filePath)
		if e != nil {
			return newError(http.StatusNotFound, fmt.Sprintf("Path was not found: %v", e)), nil
		}

		response := http.Response{Body: ioutil.NopCloser(bytes.NewBuffer(fileContents)), StatusCode: http.StatusOK}
		contentTypeHeader := "text/html"
		if strings.HasSuffix(filePath, ".js") {
			contentTypeHeader = "application/javascript"
		} else if strings.HasSuffix(filePath, ".json") {
			contentTypeHeader = "application/json"
		}

		response.Header.Add("Content-Type", contentTypeHeader)
		return &response, nil
	}
}

// SwaggerUIConfig is not really used
type SwaggerUIConfig struct{}

// Name of the config root
func (*SwaggerUIConfig) Name() string {
	return "swaggerui"
}

// SwaggerUIConfigComponent is a plugin
type SwaggerUIConfigComponent struct{}

// SwaggerUI satisfies the NewComponent signature
func SwaggerUI(_ context.Context, _ string, _ string, _ string) (interface{}, error) {
	return &SwaggerUIConfigComponent{}, nil
}

// Settings generates a config populated with defaults.
func (*SwaggerUIConfigComponent) Settings() *SwaggerUIConfig {
	return &SwaggerUIConfig{}
}

// New generates the middleware.
func (*SwaggerUIConfigComponent) New(ctx context.Context, conf *SwaggerUIConfig) (func(tripper http.RoundTripper) http.RoundTripper, error) {
	originalSpec := ctx.Value(transportd.ContextKeyOpenAPISpec).(*openapi3.Swagger)
	// we need a deep copy with the extensions removed; we don't want to reveal our configuration extensions
	bytesofspec, e := json.Marshal(originalSpec)
	if e != nil {
		return nil, e
	}
	spec, e := newSpecification(bytesofspec)
	if e != nil {
		return nil, e
	}

	spec.ExtensionProps.Extensions = nil
	for k := range spec.Paths {
		spec.Paths[k].ExtensionProps.Extensions = nil
		if spec.Paths[k].Get != nil {
			spec.Paths[k].Get.ExtensionProps.Extensions = nil
		}
		if spec.Paths[k].Post != nil {
			spec.Paths[k].Post.ExtensionProps.Extensions = nil
		}
	}

	return func(wrapped http.RoundTripper) http.RoundTripper {
		return &swaggerUITransport{
			Wrapped: wrapped,
			Spec:    spec,
		}
	}, nil
}

func newSpecification(source []byte) (*openapi3.Swagger, error) {
	envProcessor := transportd.NewEnvProcessor()
	source, err := envProcessor.Process(source)
	if err != nil {
		return nil, err
	}

	loader := openapi3.NewSwaggerLoader()
	swagger, errYaml := loader.LoadSwaggerFromData(source)
	var errJSON error
	if errYaml != nil {
		swagger, errJSON = loader.LoadSwaggerFromData(source)
	}
	if errYaml != nil && errJSON != nil {
		return nil, errJSON
	}
	return swagger, nil
}

func findConfiguredRootPath(spec *openapi3.Swagger, requestPath string) string {
	highScore := 0
	winner := ""
	for k := range spec.Paths {
		currentScore := 0
		pathToCheck := k[0:strings.LastIndex(k, "/")]
		if pathToCheck == "" {
			pathToCheck = k
		}
		for i := 0; i < min(len(pathToCheck), len(requestPath)); i++ {
			if pathToCheck[i] == requestPath[i] {
				currentScore = currentScore + 1
			} else {
				break
			}
		}
		if currentScore > highScore {
			highScore = currentScore
			winner = pathToCheck
		}

	}
	return winner
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
