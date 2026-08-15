package config

import (
	"os"
	"strconv"
)

type Config struct {
	Workers   int
	BatchSize int
	Retries   int
}

func Load() *Config {
	return &Config{
		Workers:   getInt("MAILER_WORKERS", 2),
		BatchSize: getInt("MAILER_BATCH_SIZE", 2),
		Retries:   getInt("MAILER_RETRIES", 1),
	}
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
