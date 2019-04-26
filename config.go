package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type cliArgs map[string]interface{}

//Config defines app config from
type Config struct {
	Server           *ServerConfig `yaml:"server"`
	MySQL            *MySQLConfig  `yaml:"mysql"`
	Port             string        `yaml:"server.port"`
	Pidfile          string        `yaml:"pidfile"`
	Filename         string        `yaml:"filename"`
	Source           string        `yaml:"source"`
	LogLevel         logrus.Level  `yaml:"-"`
	LogLevelAsString string        `yaml:"log_level"`
	Debug            bool          `yaml:"debug"`
}

// ServerConfig defines server configuration
type ServerConfig struct {
	Host                  string `yaml:"host"`
	Port                  string `yaml:"port"`
	ReadTimeoutString     string `yaml:"read_timeout"`
	WriteTimeoutString    string `yaml:"write_timeout"`
	GracefulTimeoutString string `yaml:"graceful_timeout"`
	MaxKeepaliveString    string `yaml:"max_keepalive"`
	ListenAddr            string
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	GracefulTimeout       time.Duration
	MaxKeepalive          time.Duration
}

// Parse defines server config and check validation
func (config *ServerConfig) Parse() error {
	config.ListenAddr = config.Host + ":" + config.Port
	var err error
	config.ReadTimeout, err = time.ParseDuration(
		config.ReadTimeoutString,
	)
	if err != nil {
		return errors.Wrap(err, "invalid read timeout")
	}
	config.WriteTimeout, err = time.ParseDuration(
		config.WriteTimeoutString,
	)

	if err != nil {
		return errors.Wrap(err, "invalid write timeout")
	}

	config.GracefulTimeout, err = time.ParseDuration(
		config.GracefulTimeoutString,
	)
	if err != nil {
		return errors.Wrap(err, "invalid graceful timeout")
	}
	config.MaxKeepalive, err = time.ParseDuration(
		config.MaxKeepaliveString,
	)
	if err != nil {
		return errors.Wrap(err, "invalid max keepalive")
	}

	return nil
}

//
func createConfig(args cliArgs) (*Config, error) {
	config, err := NewConfigFromYamlFile(args["--config"].(string))
	if err != nil {
		log.Fatalln(errors.Wrapf(err, "new config from yaml file"))
		return nil, err
	}

	// config.OverwriteWithCLIArgs(args)
	err = config.Parse()
	if err != nil {
		return nil, err

	}
	return config, nil
}

// Parse check valid yaml file
func (config *Config) Parse() error {
	var err error
	err = config.Server.Parse()
	if err != nil {
		return errors.Wrap(err, "parse rest server config")
	}

	if config.Source == "mysql" {
		if config.MySQL == nil {
			return errors.New("mysql config is not defined")
		}
		err = config.MySQL.Parse()
		if err != nil {
			return errors.Wrap(err, "parse mysql config")
		}
	} else if config.Source == "file" {
		log.Infoln("Use " + config.Filename)
		if config.Filename == "" {
			return errors.New("undefined source dict File")
		}
	} else {
		return errors.New("undefined source dict")
	}

	if config.LogLevelAsString == "" {
		config.LogLevelAsString = logrus.WarnLevel.String()
	}
	config.LogLevel, err = logrus.ParseLevel(config.LogLevelAsString)
	if err != nil {
		return errors.Wrap(err, "parse log level")
	}

	return nil
}

// NewConfigFromYamlFile defines Config with parsed yaml file
func NewConfigFromYamlFile(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "get info about file \"%s\" failed", path)
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file \"%s\" failed", path)
	}

	config := &Config{}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal failed")
	}
	log.Infoln(config.Source)
	return config, nil
}
