package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	f, err := os.CreateTemp("", "")
	assert.Nil(t, err)

	_, err = f.WriteString("Feature: Foo")
	assert.Nil(t, err)

	assert.Nil(t, command([]string{f.Name()}))

	os.Remove(f.Name())
}

func TestCommandWithNonExistentFile(t *testing.T) {
	assert.NotNil(t, command([]string{"non-existent.feature"}))
}

func TestCommandWithDirectory(t *testing.T) {
	d, err := os.MkdirTemp("", "")
	assert.Nil(t, err)

	f := filepath.Join(d, "foo.feature")
	err = os.WriteFile(f, []byte("Feature:  Foo"), 0600)
	assert.Nil(t, err)

	assert.Nil(t, command([]string{d}))

	bs, err := os.ReadFile(f)
	assert.Nil(t, err)
	assert.Equal(t, "Feature: Foo\n", string(bs))

	os.RemoveAll(d)
}
