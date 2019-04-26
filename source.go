package main

import (
	"bufio"
	"os"
)

//Source retur
type Source interface {
	parse() (map[string]bool, error)
}

type FileSource struct {
	FileSource string
}

func NewFileSource(filename string) *FileSource {
	return &FileSource{FileSource: filename}
}

func (s *FileSource) parse() (map[string]bool, error) {
	file, err := os.Open(s.FileSource)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wordsRow := make(map[string]bool)
	for scanner.Scan() {
		wordsRow[scanner.Text()] = true
	}

	return wordsRow, nil
}

type MysqlSource struct {
	MySQL *MySQL
}

func (s *MysqlSource) parse() (map[string]bool, error) {
	return nil, nil
}
