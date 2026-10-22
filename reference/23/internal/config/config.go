package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Address     string
	DatabaseURL string
	Workers     int
}

func Load() (Config, error) {
	c := Config{Address: os.Getenv("HTTP_ADDR"), DatabaseURL: os.Getenv("DATABASE_URL"), Workers: 2}
	if c.Address == "" {
		c.Address = "127.0.0.1:8080"
	}
	if c.DatabaseURL == "" {
		return c, errors.New("нужна DATABASE_URL")
	}
	if v := os.Getenv("WORKER_COUNT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 32 {
			return c, errors.New("WORKER_COUNT: целое число от 1 до 32")
		}
		c.Workers = n
	}
	return c, nil
}
