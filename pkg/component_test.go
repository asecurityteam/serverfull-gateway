package serverfullgw

import (
	"bytes"
	"context"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestComponentDoesNotAllowInvalidTemplates(t *testing.T) {
	ctx := context.Background()
	// The range of acceptable template strings is fairly broad which
	// makes it hard to validate all possible inputs. At the very
	// least we can ensure that any valid go template string is
	// accepted and invalid strings are not.
	text := `#! no closing characters`
	c := &Component{}
	conf := c.Settings()

	conf.Request = text
	conf.Success = `{"status": 200, "body": {"v2":"#!.Response.Body.v!#"}}`
	conf.Error = `{"status": 500, "bodyPassthrough": true}`

	_, err := c.New(ctx, conf)
	if err == nil {
		t.Error("did not fail on bad request template")
	}

	conf.Request = `{}`
	conf.Success = text
	conf.Error = `{"status": 500, "bodyPassthrough": true}`
	_, err = c.New(ctx, conf)
	if err == nil {
		t.Error("did not fail on bad success template")
	}

	conf.Request = `{}`
	conf.Success = `{"status": 200, "body": {"v2":"#!.Response.Body.v!#"}}`
	conf.Error = text
	_, err = c.New(ctx, conf)
	if err == nil {
		t.Error("did not fail on bad error template")
	}
}

func TestComponentTemplateFunctions(t *testing.T) {
	jsonTemplate := `#! json . !#`
	reduceTemplate := `#! reduce . !#`
	data := map[string]string{"r2": "sample", "r1": "app", "r3": "name"}

	jt, err := template.New("jt").Funcs(fns).Delims("#!", "!#").Parse(jsonTemplate)
	if err != nil {
		t.Error("Problem parsing json test template")
	}
	var jsonOutput bytes.Buffer
	err = jt.Execute(&jsonOutput, data)
	if err != nil {
		t.Error("Problem running json custom template function")
	}
	//Json.Marshal orders the keys for us
	expectedJSON := `{"r1":"app","r2":"sample","r3":"name"}`
	assert.Equal(t, expectedJSON, jsonOutput.String())

	var reduceOutput bytes.Buffer
	rt, err := template.New("rt").Funcs(fns).Delims("#!", "!#").Parse(reduceTemplate)
	if err != nil {
		t.Error("Problem parsing reduce test template")
	}
	err = rt.Execute(&reduceOutput, data)
	if err != nil {
		t.Error("Problem running reduce custom template function")
	}
	expectedReduce := "app/sample/name"
	assert.Equal(t, expectedReduce, reduceOutput.String())
}
