package config

import (
	"os"
	"strconv"
	"time"
)

type JWT struct {
	SecretKey string        `mapstructure:"SECRET_KEY"`
	Expire    time.Duration `mapstructure:"EXPIRE"`
}

func LoadJWT() JWT {
	return JWT{
		SecretKey: getENV("SECRET_KEY", ""),
		Expire:    time.Minute * time.Duration(parseInt(getENV("EXPIRE", "15"))),
	}
}

func getENV(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
