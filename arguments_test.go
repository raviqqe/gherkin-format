package main_test

import (
	"testing"

	"github.com/raviqqe/gherkin-format"
	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	args, err := main.GetArguments([]string{})

	assert.Nil(t, err)
	assert.Equal(t, main.Arguments{Paths: []string{}}, args)
}

func TestGetArgumentsPath(t *testing.T) {
	args, err := main.GetArguments([]string{"path"})

	assert.Nil(t, err)
	assert.Equal(t, main.Arguments{Paths: []string{"path"}}, args)
}

func TestGetArgumentsMultiplePaths(t *testing.T) {
	args, err := main.GetArguments([]string{"path1", "path2"})

	assert.Nil(t, err)
	assert.Equal(t, main.Arguments{Paths: []string{"path1", "path2"}}, args)
}
