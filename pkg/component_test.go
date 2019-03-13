package serverfullgw

import (
	"context"
	"testing"
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
