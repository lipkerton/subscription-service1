package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
}

type AppConfig struct {
	Port int
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	appPort, err := getEnvAsInt("APP_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("get app port: %w", err)
	}

	postgresPort, err := getEnvAsInt("POSTGRES_PORT", 5432)
	if err != nil {
		return nil, fmt.Errorf("get postgres port: %w", err)
	}
	return &Config{
		App: AppConfig{
			Port: appPort,
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     postgresPort,
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DB", "subscriptions_db"),
			SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
		},
	}, nil
}

func (c PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
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
