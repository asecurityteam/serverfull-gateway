package serverfullgw

import (
	"testing"

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
