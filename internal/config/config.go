package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr   string
	MaxRequests int
	BanDuration time.Duration
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	maxRequests, err := strconv.Atoi(os.Getenv("MAX_REQUESTS"))
	if err != nil {
		return nil, err
	}

	banDuration, err := time.ParseDuration(os.Getenv("BAN_DURATION"))
	if err != nil {
		return nil, err
	}

	return &Config{
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		MaxRequests: maxRequests,
		BanDuration: banDuration,
	}, nil
}
