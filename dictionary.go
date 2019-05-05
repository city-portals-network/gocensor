package main

import (
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

func (dictionary *Dictionary) parse() error {
	wordsPrepared := make(map[string]bool)
	words, err := dictionary.Source.parse()
	if err != nil {
		return err
	}

	for word, _ := range words {
		stemWord := dictionary.stemWord(word)
		wordsPrepared[stemWord] = true
	}
	dictionary.Lock()
	dictionary.words = wordsPrepared
	dictionary.Unlock()
	return nil
}

func (dictionary *Dictionary) append(word string) error {
	err := dictionary.Source.append(word)
	if err != nil {
		return err
	}
	dictionary.Lock()
	dictionary.words[word] = true
	dictionary.Unlock()
	return nil
}

func (dictionary *Dictionary) delete(word string) error {
	err := dictionary.Source.delete(word)
	if err != nil {
		return err
	}
	dictionary.Lock()
	delete(dictionary.words, word)
	dictionary.Unlock()
	return nil
}

func (dictionary *Dictionary) reload() error {
	return dictionary.parse()
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
		sWord := dictionary.stemWord(word)
		stemmedWords = append(stemmedWords, sWord)
	}
	return stemmedWords
}

func (dictionary *Dictionary) stemWord(word string) string {
	stemmed, err := snowball.Stem(word, "russian", false)
	if err == nil {
		return stemmed
	}
	return word
}
