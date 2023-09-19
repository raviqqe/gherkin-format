package main_test

import (
	"testing"

	"github.com/raviqqe/gherkin-format"
	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, c := range []struct {
		parameters []string
		arguments  main.Arguments
	}{
		{[]string{"path"}, main.Arguments{Path: "path"}},
	} {
		args, err := main.GetArguments(c.parameters)

		assert.Nil(t, err)
		assert.Equal(t, c.arguments, args)
	}
}
