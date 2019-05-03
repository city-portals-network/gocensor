package main

import (
	"bufio"
	"os"

	"github.com/pkg/errors"
)

type FileSource struct {
	FileSource string
	File       *os.File
}

func NewFileSource(filename string) (*FileSource, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, errors.Wrap(err, "file  "+filename)
	}
	return &FileSource{FileSource: filename}, nil
}

func (s *FileSource) append(word string) error {
	file, err := os.OpenFile(s.FileSource, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "cant open file "+s.FileSource)
	}

	defer file.Close()

	if _, err = file.WriteString(word + "\r\n"); err != nil {
		return errors.Wrap(err, "cant append to file "+s.FileSource)
	}
	return nil
}

func (s *FileSource) delete(word string) error {
	words, err := s.parse()
	if err != nil {
		return err
	}
	delete(words, word)
	tmpFilename := s.FileSource + "-tmp"
	file, err := os.Create(tmpFilename)
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "cant append to file "+s.FileSource)
	}

	for word := range words {
		if _, err = file.WriteString(word); err != nil {
			return errors.Wrap(err, "cant append to file "+s.FileSource)
		}
	}
	err = os.Rename(tmpFilename, s.FileSource)
	if err != nil {
		return errors.Wrap(err, "cant append to file "+s.FileSource)
	}
	return nil
}

func (s *FileSource) parse() (map[string]bool, error) {
	file, err := os.Open(s.FileSource)
	if err != nil {
		return nil, errors.Wrap(err, "cant parse file "+s.FileSource)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wordsRow := make(map[string]bool)
	for scanner.Scan() {
		wordsRow[scanner.Text()] = true
	}

	return wordsRow, nil
}
