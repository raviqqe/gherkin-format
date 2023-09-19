package main_test

import (
	"testing"

	"github.com/raviqqe/gherkin-format"
	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	args, err := main.GetArguments([]string{"path"})

	assert.Nil(t, err)
	assert.Equal(t, main.Arguments{Path: "path"}, args)
}
