// +build integration

package tests

import (
	"context"
	"testing"

	"github.com/asecurityteam/serverfull-gateway/pkg"
	"github.com/asecurityteam/transportd/pkg"
	"github.com/asecurityteam/transportd/pkg/components"
	"github.com/stretchr/testify/assert"
)

func TestHelpNoErrors(t *testing.T) {
	h, err := transportd.Help(
		context.Background(),
		append(components.Defaults, serverfullgw.Lambda)...,
	)
	// Basic sanity check that the help output works with
	// the native components and is not empty.
	assert.Nil(t, err)
	assert.NotEmpty(t, h)
}
