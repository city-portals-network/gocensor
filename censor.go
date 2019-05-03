package main

// Censor ds
type Censor struct {
	dictionary *Dictionary
}

//NewCensor defines new censor instance
func NewCensor(dictionary *Dictionary) (*Censor, error) {
	if err := dictionary.parse(); err != nil {
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
	return true
}

func (censor *Censor) reload() bool {
	return true
}
