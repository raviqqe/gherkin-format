package main_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/raviqqe/gherkin-format"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	f, err := os.CreateTemp("", "")
	assert.Nil(t, err)

	_, err = f.WriteString("Feature: Foo")
	assert.Nil(t, err)

	assert.Nil(t, main.Run([]string{f.Name()}))

	assert.Nil(t, os.Remove(f.Name()))
}

func TestCommandWithNonExistentFile(t *testing.T) {
	assert.NotNil(t, main.Run([]string{"non-existent.feature"}))
}

func TestCommandWithDirectory(t *testing.T) {
	d, err := os.MkdirTemp("", "")
	assert.Nil(t, err)

	f := filepath.Join(d, "foo.feature")
	err = os.WriteFile(f, []byte("Feature:  Foo"), 0600)
	assert.Nil(t, err)

	assert.Nil(t, main.Run([]string{d}))

	bs, err := os.ReadFile(f)
	assert.Nil(t, err)
	assert.Equal(t, "Feature: Foo\n", string(bs))

	assert.Nil(t, os.RemoveAll(d))
}
