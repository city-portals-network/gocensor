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

func (censor *Censor) check(comment string) *CensorResponse {
	response := NewCensorResponse(comment)
	words := censor.dictionary.prepareString(comment)
	for _, word := range words {
		if censor.dictionary.words[word] {
			return response.setResult(true)
		}
	}
	return response.setResult(false)
}

func (censor *Censor) append(word string) error {
	if err := censor.dictionary.append(word); err != nil {
		return err
	}
	return nil
}

func (censor *Censor) delete(word string) error {
	if err := censor.dictionary.delete(word); err != nil {
		return err
	}
	return nil
}

func (censor *Censor) reload() error {
	if err := censor.dictionary.parse(); err != nil {
		return err
	}
	return nil
}
