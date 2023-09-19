package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFile(t *testing.T) {
	f, err := os.CreateTemp("", "")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	_, err = f.Write([]byte("Feature: Foo"))
	assert.Nil(t, err)

	assert.Nil(t, formatFile(f.Name()))
}

func TestFormatFileError(t *testing.T) {
	f, err := os.CreateTemp("", "")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	_, err = f.Write([]byte("Feature"))
	assert.Nil(t, err)

	assert.NotNil(t, formatFile(f.Name()))
}

func TestFormatFilesWithNonReadableSourceDir(t *testing.T) {
	assert.NotNil(t, formatFiles("foo"))
}
