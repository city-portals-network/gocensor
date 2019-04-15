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

//TODO add config.yml
type Config struct {
	Server           *ServerConfig `yaml:"server"`
	Compress         bool          `env:"COMPRESS"`
	Port             string        `yaml:"server.port"`
	Pidfile          string        `yaml:"pidfile"`
	LogLevel         logrus.Level  `yaml:"-"`
	LogLevelAsString string        `yaml:"log_level"`
	Debug            bool          `yaml:"debug"`
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	ReadTimeoutString  string `yaml:"read_timeout"`
	WriteTimeoutString string `yaml:"write_timeout"`
	// подготовленные данные:
	// host + port
	ListenAddr string
	// значения таймаутов после парсинга
	// строковых значений из yml
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Parse парсит значения конфига и
// проверяет их валидность
func (config *ServerConfig) Parse() error {
	config.ListenAddr =
		config.Host + ":" +
			config.Port

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
	return nil
}

//
func createConfig(args cliArgs) *Config {
	config, err := NewConfigFromYamlFile(args["--config"].(string))
	if err != nil {
		log.Fatalln(errors.Wrapf(err, "new config from yaml file"))
	}
	// config.OverwriteWithCLIArgs(args)
	err = config.Parse()
	if err != nil {
		log.Fatalln(errors.Wrap(err, "parse config"))
	}
	return config
}

// Parse проверяет валидность всех параметров конфигурации
func (config *Config) Parse() error {
	var err error
	err = config.Server.Parse()
	if err != nil {
		return errors.Wrap(err, "parse rest server config")
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

// NewConfigFromYamlFile вернет Config с данными из yaml файла.
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
	return config, nil
}
