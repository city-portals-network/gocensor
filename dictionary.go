package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/kljensen/snowball"
)

// Dictionary ds
type Dictionary struct {
	Source Source
	sync.Mutex
	words map[string]bool
}

//NewDictionary return dict
func NewDictionary(source Source) *Dictionary {
	words := make(map[string]bool)
	return &Dictionary{Source: source, words: words}
}

func (dictionary *Dictionary) parse() {
	words, _ := dictionary.Source.parse()
	dictionary.words = words
}

func (dictionary *Dictionary) prepareString(text string) []string {
	reg, _ := regexp.Compile("[^a-zа-я]+")
	text = reg.ReplaceAllLiteralString(text, " ")
	text = strings.ToLower(text)
	words := strings.Split(text, " ")
	var stemmedWords []string
	for _, word := range words {
		if len([]rune(word)) < 2 {
			continue
		}
		stemmed, err := snowball.Stem(word, "russian", false)
		if err == nil {
			stemmedWords = append(stemmedWords, stemmed)
		}
	}
	return stemmedWords
}

func (dictionary *Dictionary) parseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	dictionary.Lock()
	for scanner.Scan() {
		words := dictionary.prepareString(scanner.Text())
		for _, word := range words {
			dictionary.words[word] = true
		}
	}
	log.Infof("appended words")
	dictionary.Unlock()
	return nil
}
