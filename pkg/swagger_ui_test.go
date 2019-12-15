package serverfullgw

import (
	"context"
	"fmt"
	"testing"

	transportd "github.com/asecurityteam/transportd/pkg"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMaxMin(t *testing.T) {
	// max
	result := max(0, 1)
	assert.Equal(t, 1, result)
	result = max(1, 1)
	assert.Equal(t, 1, result)
	result = max(1, 0)
	assert.Equal(t, 1, result)

	// min
	result = min(0, 1)
	assert.Equal(t, 0, result)
	result = min(0, 0)
	assert.Equal(t, 0, result)
	result = min(1, 0)
	assert.Equal(t, 0, result)

}

func TestContextSpecDoc(t *testing.T) {

	type Something struct{}
	extensions := make(map[string]interface{})
	extensions["blarg"] = Something{}
	extensionProps := openapi3.ExtensionProps{Extensions: extensions}
	paths := make(map[string]*openapi3.PathItem)
	operation := openapi3.NewOperation()
	operation.ExtensionProps = extensionProps
	paths["one"] = &openapi3.PathItem{
		ExtensionProps: extensionProps,
		Get:            operation,
	}

	spec := openapi3.Swagger{
		ExtensionProps: extensionProps,
		Paths:          paths,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, transportd.ContextKeyOpenAPISpec, &spec)

	c := &SwaggerUIConfigComponent{}
	conf := c.Settings()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, err := c.New(ctx, conf)

	assert.Nil(t, err, fmt.Sprintf("Expected no errors but got %v", err))
}
