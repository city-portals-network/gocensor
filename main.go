package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	docopt "github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const usage = `gocensor

Usage:
  gocensor --config <path> [options]

Options:
  --config <path>              Configuration file in YAML format.
                                 CLI args overwrite parameters in config file.
  --host <host>                HTTP Server host.
  --port <port>                HTTP Server port.
  -h --help                    Show this screen.
`

const welcomeMessage = `Gocensor is running!

PID: %+v

Log level: %s

HTTP server configuration: %v

`

var (
	version = "dev"
	log     *logrus.Logger
)

func main() {

	cfg := createConfig(parseCLIArgs())
	createPidFile(cfg.Pidfile)

	renderWelcomeScreen(cfg)
	initializeLogger()
	setVerbosityLevel(cfg)

	censor, err := NewCensor()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	server := NewServer(censor, cfg)
	server.Routes()
	err = server.Run()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func setVerbosityLevel(config *Config) {
	log.SetLevel(config.LogLevel)
}

//
func parseCLIArgs() cliArgs {
	args, err := docopt.Parse(usage, nil, true, version, true)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "parse cli args"))
	}
	return args
}

func createPidFile(pidfile string) error {
	pid := []byte(strconv.Itoa(os.Getpid()))
	err := ioutil.WriteFile(pidfile, pid, 0664)
	if err != nil {
		log.Fatalln(errors.Wrapf(err, "create pid file \"%s\"", pidfile))
	}
	return nil
}

func initializeLogger() {
	log = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.TextFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.TraceLevel,
	}
}

func renderWelcomeScreen(config *Config) {
	fmt.Printf(
		welcomeMessage,
		os.Getpid(),
		config.LogLevel.String(),
		config.Server,
	)
}
