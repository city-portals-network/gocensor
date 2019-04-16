package main

// Censor ds
type Censor struct {
	dictionary *dictionary
}

//NewCensor defines new censor instance
func NewCensor() (*Censor, error) {
	words := make(map[string]bool)
	dictionary := &dictionary{words: words}
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
