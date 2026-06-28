package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
}

type AppConfig struct {
	Port int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port, err := getEnvAsInt("APP_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("get app port: %w", err)
	}
	return &Config{
		App: AppConfig{
			Port: port,
		},
	}, nil

}

func getEnvAsInt(key string, defaultValue int) (int, error) {
	port := os.Getenv(key)
	if port == "" {
		return defaultValue, nil
	}

	intValue, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("env %s must be valid integer: %w", key, err)
	}
	return intValue, nil
}
