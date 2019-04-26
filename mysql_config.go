package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type MySQLConfig struct {
	Master       *MySQLConnectionConfig `yaml:"master"`
	Slave        *MySQLConnectionConfig `yaml:"slave"`
	slaveDefined bool                   `yaml:"-"`
}

type MySQLConnectionConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Base       string `yaml:"base"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Parameters struct {
		MaxIdleConns            int           `yaml:"max_idle_conns"`
		MaxOpenConns            int           `yaml:"max_open_conns"`
		ConnMaxLifetimeAsString string        `yaml:"conn_max_lifetime"`
		ConnMaxLifetime         time.Duration `yaml:"-"`
	} `yaml:"parameters"`
}

// Parse ...
func (config *MySQLConfig) Parse() error {
	if config.Master == nil {
		return errors.New("mysql master is not defined")
	}
	err := config.Master.Parse()
	if err != nil {
		return errors.Wrap(err, "parse master config")
	}
	return nil
}

// Parse ...
func (config *MySQLConnectionConfig) Parse() error {
	var err error
	if config.Parameters.ConnMaxLifetimeAsString != "" {
		config.Parameters.ConnMaxLifetime, err = time.ParseDuration(
			config.Parameters.ConnMaxLifetimeAsString,
		)
		if err != nil {
			return errors.Wrap(err, "invalid conn_max_lifetime")
		}
	}
	return nil
}

//
func (config *MySQLConnectionConfig) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Base,
	)
}
