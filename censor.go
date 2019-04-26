package main

import "fmt"

// Censor ds
type Censor struct {
	dictionary *Dictionary
}

//NewCensor defines new censor instance
func NewCensor(dictionary *Dictionary) (*Censor, error) {
	if err := dictionary.parseFile("dictionary.txt"); err != nil {
		return nil, err
	}
	censor := &Censor{dictionary: dictionary}
	return censor, nil
}

func (censor *Censor) run(comment string) bool {
	words := censor.dictionary.prepareString(comment)
	for _, word := range words {
		if censor.dictionary.words[word] {
			return true
		}
	}
	return false
}

func (censor *Censor) update(word string) bool {
	fmt.Println(word)
	return true
}

func (censor *Censor) reload() bool {
	return true
}
