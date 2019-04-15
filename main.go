package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	version  = "dev"
	hostname string
	log      *logrus.Logger
)

func main() {
	var err error
	initializeLogger()
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatalln(errors.Wrap(err, "get hostname failed"))
	}

	censor := NewCensor()
	if censor.run("Фильм \"Дурак\" напоминает. Наверняка были предпосылки, но чьё-то авось, преступная халатность или жадность довели до трагедии ни в ") {
		log.Infof("True")
	} else {
		log.Infof("FALSE")
	}
}

func initializeLogger() {
	log = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.TraceLevel,
	}
}
