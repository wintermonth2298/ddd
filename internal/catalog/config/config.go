package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wintermonth2298/library-ddd/internal/pkg/env"
)

type Config struct {
	PSQL PSQL
}

type PSQL struct {
	Port     string
	User     string
	Password string
	DB       string
	Host     string
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Panicf("No .env file found (fallback to OS environment)")
	}

	requiredEnvVars := []string{
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
	}

	if err := env.CheckEnvVars(requiredEnvVars); err != nil {
		log.Panicf("check env vars: %v", err)
	}

	return &Config{
		PSQL: PSQL{
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DB:       os.Getenv("POSTGRES_DB"),
			Host:     os.Getenv("POSTGRES_HOST"),
		},
	}
}
