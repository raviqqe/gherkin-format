package main_test

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
		args, err := getArguments(c.parameters)

		assert.Nil(t, err)
		assert.Equal(t, c.arguments, args)
	}
}

func TestParseArgumentsWithoutUsage(t *testing.T) {
	err := parseArguments("", []string{"path"}, &arguments{})

	assert.Error(t, err)
}

func TestParseArgumentsPanic(t *testing.T) {
	assert.Panics(t, func() {
		parseArguments(usage, []string{"path"}, arguments{})
	})
}
