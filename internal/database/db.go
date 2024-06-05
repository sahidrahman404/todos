package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	Bun *bun.DB
}

func New(dsn string) (*DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxIdleTime(5 * time.Minute)
	sqldb.SetConnMaxLifetime(2 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := sqldb.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	return &DB{Bun: db}, nil
}
