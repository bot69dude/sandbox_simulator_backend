package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	ServerPort  string
	DatabaseURL string
	Environment string
}

func LoadConfig() (*ServerConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &ServerConfig{
		ServerPort:  GetEnv("SERVER_PORT", "8080"),
		DatabaseURL: GetEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
		Environment: GetEnv("ENVIRONMENT", "development"),
	}, nil
}

func GetEnv(key, defaultValue string) string {
	if values, exists := os.LookupEnv(key); exists {
		return values
	}
	return defaultValue
}