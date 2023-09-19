package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFile(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	_, err = f.Write([]byte("Feature: Foo"))
	assert.Nil(t, err)

	assert.Nil(t, formatFile(f.Name(), ioutil.Discard))
}

func TestFormatFileError(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	_, err = f.Write([]byte("Feature"))
	assert.Nil(t, err)

	assert.NotNil(t, formatFile(f.Name(), ioutil.Discard))
}

func TestFormatFilesWithNonReadableSourceDir(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(d)

	assert.NotNil(t, formatFiles("foo", d))
}