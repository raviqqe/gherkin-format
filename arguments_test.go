package main_test

import (
	"testing"

	"github.com/raviqqe/gherkin-format"
	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	args, err := main.GetArguments([]string{})

	assert.Nil(t, err)
	assert.Equal(t, main.Arguments{Paths: nil}, args)
}

func TestGetArgumentsPath(t *testing.T) {
	args, err := main.GetArguments([]string{"path"})

	assert.Nil(t, err)
	assert.Equal(t, main.Arguments{Paths: []string{"path"}}, args)
}

func TestGetArgumentsPaths(t *testing.T) {
	args, err := main.GetArguments([]string{"path1", "path2"})

	assert.Nil(t, err)
	assert.Equal(t, main.Arguments{Paths: []string{"path"}}, args)
}

func TestGetArgumentsCheck(t *testing.T) {
	args, err := main.GetArguments([]string{"--check"})

	assert.Nil(t, err)
	assert.Equal(t, main.Arguments{Check: true, Paths: []string{}}, args)
}
