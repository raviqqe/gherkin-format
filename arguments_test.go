package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, c := range []struct {
		parameters []string
		arguments
	}{
		{[]string{"path"}, arguments{Path: "path"}},
	} {
		assert.Equal(t, c.arguments, getArguments(c.parameters))
	}
}

func TestParseArgumentsPanic(t *testing.T) {
	assert.Panics(t, func() {
		parseArguments("", []string{"path"}, &arguments{})
	})

	assert.Panics(t, func() {
		parseArguments(usage, []string{"path"}, arguments{})
	})
}
