package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// MySQL ...
type MySQL struct {
	Master *sql.DB
	Slave  *sql.DB
	config *MySQLConfig
}

// NewMySQL вернет полностью готовый объект для работы с БД.
func NewMySQL(config *MySQLConfig) *MySQL {
	return &MySQL{
		config: config,
	}
}

// OpenConnections ...
func (mysql *MySQL) OpenConnections() error {
	var err error
	mysql.Master, err = openMySQLDatabaseConnection(mysql.config.Master)
	if err != nil {
		return errors.Wrap(err, "open connection to master")
	}
	if mysql.config.slaveDefined {
		mysql.Slave, err = openMySQLDatabaseConnection(mysql.config.Slave)
		if err != nil {
			return errors.Wrap(err, "open connection to slave")
		}
	} else {
		mysql.Slave = mysql.Master
	}
	return nil
}

// Ping ...
func (mysql *MySQL) Ping() error {
	var err error
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()
	err = mysql.Master.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, "ping mysql master")
	}
	if mysql.config.slaveDefined {
		err = mysql.Slave.PingContext(ctx)
		if err != nil {
			return errors.Wrap(err, "ping mysql slave")
		}
	}
	return nil
}

// Close закрывает все открытые соединения с БД
func (mysql *MySQL) Close() {
	err := mysql.Master.Close()
	if err != nil {
		log.Errorln(errors.Wrap(err, "close mysql master"))
	}
	err = mysql.Slave.Close()
	if err != nil {
		log.Errorln(errors.Wrap(err, "close mysql slave"))
	}
}

//
func openMySQLDatabaseConnection(config *MySQLConnectionConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.getDSN())
	if err != nil {
		return nil, errors.Wrapf(err, "open mysql connection"+
			" for dsn \"%s\"", config.getDSN())
	}
	// установка параметров соединения
	db.SetMaxIdleConns(config.Parameters.MaxIdleConns)
	db.SetMaxOpenConns(config.Parameters.MaxOpenConns)
	db.SetConnMaxLifetime(config.Parameters.ConnMaxLifetime)
	return db, nil
}
