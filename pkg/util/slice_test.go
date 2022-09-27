package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInSlice(t *testing.T) {
	assert.True(t, InSlice[string]("a", []string{"a", "b", "c"}))
	assert.False(t, InSlice[string]("d", []string{"a", "b", "c"}))
	assert.True(t, InSlice[int](1, []int{1, 2, 3}))
	assert.False(t, InSlice[int](4, []int{1, 2, 3}))
}
