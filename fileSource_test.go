package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createFile(filename string, words string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if words != "" {
		f.WriteString(words)
	}
	defer f.Close()
	return err
}
func TestNewFileSource(t *testing.T) {
	assert := assert.New(t)
	filename := "test.txt"

	source, err := NewFileSource(filename)
	assert.EqualError(err, "file  "+filename+": stat "+filename+": no such file or directory")
	err = createFile(filename, "")
	assert.NoError(err)
	source, err = NewFileSource(filename)
	assert.NoError(err)
	assert.Equal(
		&FileSource{FileSource: filename},
		source,
	)
	os.Remove(filename)
}

func TestFileSourceAppend(t *testing.T) {
	assert := assert.New(t)
	filename := "test.txt"

	err := createFile(filename, "")
	assert.NoError(err)
	source, _ := NewFileSource(filename)
	err = source.append("test1")
	assert.Nil(err)
	err = source.append("test2")
	assert.Nil(err)
	os.Remove(filename)
}

func TestFileSourceParse(t *testing.T) {
	assert := assert.New(t)
	filename := "test.txt"

	err := createFile(filename, "test1\r\ntest2")
	assert.NoError(err)

	source, err := NewFileSource(filename)
	assert.NoError(err)
	words, err := source.parse()
	assert.NoError(err)
	wordsResult := map[string]bool{"test1": true, "test2": true}
	assert.Equal(words, wordsResult)
	os.Remove(filename)
}

func TestFileSourceDelete(t *testing.T) {
	assert := assert.New(t)
	filename := "test.txt"

	err := createFile(filename, "test1\r\ntest2")
	assert.NoError(err)

	source, err := NewFileSource(filename)
	assert.NoError(err)
	err = source.delete("test1")
	assert.NoError(err)

	wordsResult := map[string]bool{"test2": true}
	words, err := source.parse()
	assert.NoError(err)
	assert.Equal(words, wordsResult)
	os.Remove(filename)
}
