package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLmode  string
}

func (c Config) DNS() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLmode)
}

func NewPool(ctx context.Context, cnf Config) (*pgxpool.Pool, error) {
	pcfg, err := pgxpool.ParseConfig(cnf.DNS())
	if err != nil {
		return nil, err
	}
	pcfg.MaxConns = 25
	pcfg.MinConns = 5
	pcfg.MaxConnLifetime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
