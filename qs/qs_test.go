package qs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	assert.Equal(t, "1", MustParse("https://example.com/a?a=1&b=2")["a"])
	assert.Equal(t, "2", MustParse("https://example.com/a?a=1&b=2")["b"])
	assert.Equal(t, "1", MustParse("https://example.com/a/?a=1&b=2")["a"])
	assert.Equal(t, "2", MustParse("https://example.com/a/?a=1&b=2")["b"])
	assert.Equal(t, "1", MustParse("/?a=1&b=2")["a"])
	assert.Equal(t, "2", MustParse("/?a=1&b=2")["b"])
	assert.Equal(t, "1", MustParse("?a=1&b=2")["a"])
	assert.Equal(t, "2", MustParse("?a=1&b=2")["b"])
	assert.Equal(t, "1", MustParse("a=1&b=2")["a"])
	assert.Equal(t, "2", MustParse("a=1&b=2")["b"])
}

func TestStringify(t *testing.T) {
	assert.Equal(t, "a=1&b=2", MustStringify(map[string]string{"a": "1", "b": "2"}))
	assert.Equal(t, "a=1&b=2", MustStringify(map[string]interface{}{"a": 1, "b": 2}))
	assert.Equal(t, "a=1&b=2", MustStringify(struct {
		A string `json:"a"`
		B string `json:"b"`
	}{A: "1", B: "2"}))
}
