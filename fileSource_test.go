package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileSource(t *testing.T) {
	assert := assert.New(t)

	source := NewFileSource("test.txt")

	assert.Equal(
		&FileSource{FileSource: "test.txt"},
		source,
	)
}

func TestFileSourceAppend(t *testing.T) {
	assert := assert.New(t)

	source := NewFileSource("test.txt")
	err := source.append("test1")
	assert.Nil(err)
	source.append("test2")
	source.append("test3")

	assert.Equal(
		&FileSource{FileSource: "test.txt"},
		source,
	)
}
