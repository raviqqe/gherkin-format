package main_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/raviqqe/gherkin-format"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	assert.Nil(t, main.Format(bytes.NewBufferString("Feature: Foo"), bytes.NewBufferString("")))
}

func TestFormatPaths(t *testing.T) {
	f, err := os.CreateTemp("", "*.feature")
	assert.Nil(t, err)
	defer func() { assert.Nil(t, os.Remove(f.Name())) }()

	_, err = f.Write([]byte("Feature:   Foo"))
	assert.Nil(t, err)

	assert.Nil(t, main.FormatPaths([]string{f.Name()}))
	bs, err := os.ReadFile(f.Name())
	assert.Nil(t, err)
	assert.Equal(t, "Feature: Foo\n", string(bs))
}

func TestCheckPaths(t *testing.T) {
	f, err := os.CreateTemp("", "*.feature")
	assert.Nil(t, err)
	defer func() { assert.Nil(t, os.Remove(f.Name())) }()

	_, err = f.Write([]byte("Feature: Foo\n"))
	assert.Nil(t, err)

	assert.Nil(t, main.CheckPaths([]string{f.Name()}))
}

func TestFormatPathsError(t *testing.T) {
	f, err := os.CreateTemp("", "*.feature")
	assert.Nil(t, err)
	defer func() { assert.Nil(t, os.Remove(f.Name())) }()

	_, err = f.Write([]byte("Feature"))
	assert.Nil(t, err)

	assert.NotNil(t, main.FormatPaths([]string{f.Name()}))
}

func TestCheckPathsError(t *testing.T) {
	f, err := os.CreateTemp("", "*.feature")
	assert.Nil(t, err)
	defer func() { assert.Nil(t, os.Remove(f.Name())) }()

	_, err = f.Write([]byte("Feature:   Foo"))
	assert.Nil(t, err)

	assert.NotNil(t, main.CheckPaths([]string{f.Name()}))
}

func TestFormatPathsWithNonReadableDirectory(t *testing.T) {
	assert.NotNil(t, main.FormatPaths([]string{"foo"}))
}

func TestCheckPathsWithNonReadableDirectory(t *testing.T) {
	assert.NotNil(t, main.CheckPaths([]string{"foo"}))
}
