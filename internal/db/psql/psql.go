package psql

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ttrueno/rl2-final/config"
	"github.com/ttrueno/rl2-final/internal/lib/e"
)

var (
	ErrNoRecords = errors.New("no records found")
)

type DB struct {
	Conn *pgx.Conn
}

func Connect(ctx context.Context, cfg config.DbConnConfig) (_ *DB, err error) {
	var errmsg = "psql.Connect"

	defer func() { err = e.WrapIfErr(errmsg, err) }()
	poolConfig, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		return nil, err
	}

	maxConnIdleTime, err := time.ParseDuration(cfg.ConnPoolConfig.MaxConnIdleTime)
	if err != nil {
		return nil, err
	}
	maxConnLifetime, err := time.ParseDuration(cfg.ConnPoolConfig.MaxConnLifeTime)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.ConnPoolConfig.MaxConns)
	poolConfig.MaxConnIdleTime = maxConnIdleTime
	poolConfig.MaxConnLifetime = maxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return &DB{
		Conn: conn.Conn(),
	}, nil
}

func (db *DB) Close(ctx context.Context) error {
	return db.Conn.Close(ctx)
}
