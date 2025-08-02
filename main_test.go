package main_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/raviqqe/gherkin-format"
	"github.com/stretchr/testify/assert"
)

func createTempFeature(t *testing.T) *os.File {
	f, err := os.CreateTemp("", "*.feature")
	assert.Nil(t, err)
	return f
}

func TestCommand(t *testing.T) {
	f := createTempFeature(t)

	_, err := f.WriteString("Feature: Foo")
	assert.Nil(t, err)

	assert.Nil(t, main.Run([]string{f.Name()}, io.Discard))

	assert.Nil(t, os.Remove(f.Name()))
}

func TestCommandWithNonExistentFile(t *testing.T) {
	assert.NotNil(t, main.Run([]string{"non-existent.feature"}, io.Discard))
}

func TestCommandWithDirectory(t *testing.T) {
	d, err := os.MkdirTemp("", "")
	assert.Nil(t, err)

	f := filepath.Join(d, "foo.feature")
	err = os.WriteFile(f, []byte("Feature:  Foo"), 0600)
	assert.Nil(t, err)

	assert.Nil(t, main.Run([]string{d}, io.Discard))

	bs, err := os.ReadFile(f)
	assert.Nil(t, err)
	assert.Equal(t, "Feature: Foo\n", string(bs))

	assert.Nil(t, os.RemoveAll(d))
}

func TestCommandCheckSuccess(t *testing.T) {
	f := createTempFeature(t)

	_, err := f.WriteString("Feature: Foo")
	assert.Nil(t, err)

	assert.Nil(t, main.Run([]string{"-check", f.Name()}, io.Discard))

	assert.Nil(t, os.Remove(f.Name()))
}

func TestCommandCheckFailure(t *testing.T) {
	f := createTempFeature(t)

	_, err := f.WriteString("Feature:  Foo")
	assert.Nil(t, err)

	assert.Error(t, main.Run([]string{"-check", f.Name()}, io.Discard))

	assert.Nil(t, os.Remove(f.Name()))
}

func TestCommandVersion(t *testing.T) {
	b := bytes.NewBuffer(nil)

	assert.Nil(t, main.Run([]string{"-version"}, b))
	assert.Equal(t, "0.1.0\n", b.String())
}
