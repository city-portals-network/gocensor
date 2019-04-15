package main

//NewCensor ds
func NewCensor() *Censor {
	words := make(map[string]bool)
	dictionary := &dictionary{words: words}
	dictionary.parseFile("dictionary.txt")
	censor := &Censor{dictionary: dictionary}
	return censor
}

// Censor ds
type Censor struct {
	dictionary *dictionary
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
