package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/kljensen/snowball"
)

// dictionary ds
type dictionary struct {
	sync.Mutex
	words map[string]bool
}

func (dictionary *dictionary) prepareString(text string) []string {
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

func (dictionary *dictionary) parseFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
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
}
