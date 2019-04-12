package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	version  = "dev"
	hostname string
	log      *logrus.Logger
)

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
type dictionary struct {
	sync.Mutex
	words map[string]bool
}

func (censor *Censor) run(comment string) bool {
	reg, _ := regexp.Compile("/[a-zA-Z]+/")
	comment = reg.ReplaceAllLiteralString(comment, " ")
	comment = strings.ToLower(comment)
	words := strings.Split(comment, " ")
	for _, word := range words {
		if censor.dictionary.words[word] {
			return true
		}
	}
	return false
}

func main() {

	var err error
	initializeLogger()
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatalln(errors.Wrap(err, "get hostname failed"))
	}

	censor := NewCensor()
	if censor.run("Фильм \"Дурак\" хуй напоминает. Наверняка были предпосылки, но чьё-то авось, преступная халатность или жадность довели до трагедии ни в ") {
		log.Infof("True")
	} else {
		log.Infof("FALSE")
	}
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
		dictionary.words[scanner.Text()] = true
	}
	log.Infof("appended words")
	dictionary.Unlock()
}

func initializeLogger() {
	log = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.TraceLevel,
	}
}
