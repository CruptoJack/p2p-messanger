package database

import (
	"context"
	"fmt"
	"os"
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

func ConfigFromENV() Config {
	return Config{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLmode:  "disable",
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLmode)
}

func NewPool(ctx context.Context, cnf Config) (*pgxpool.Pool, error) {
	pcfg, err := pgxpool.ParseConfig(cnf.DSN())
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
