package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	Master
	Slave

	master *sqlx.DB
	slave  *sqlx.DB

	driver string
}

func New(driver, master, slave string) (*DB, error) {
	masterDB, err := sqlx.Open(driver, master)
	if err != nil {
		return nil, err
	}

	slaveDB, err := sqlx.Open(driver, slave)
	if err != nil {
		return nil, err
	}

	return &DB{
		Master: masterDB,
		Slave:  slaveDB,

		master: masterDB,
		slave:  slaveDB,
		driver: driver,
	}, nil
}

type Master interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)

	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Slave interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func (db *DB) GetMaster() *sqlx.DB {
	return db.master
}

func (db *DB) GetSlave() *sqlx.DB {
	return db.slave
}

func (db *DB) Rebind(query string) string {
	return sqlx.Rebind(sqlx.BindType(db.driver), query)
}
