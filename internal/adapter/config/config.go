package config

import (
	ut "bam/internal/core/utils"
	"os"

	"github.com/joho/godotenv"
)

type (
	DB struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
	}
    // HTTP struct {
	// 	Env            string
	// 	URL            string
	// 	Port           string
	// 	AllowedOrigins string
	// 	Prefix         string
	// }
)

type Container struct {
	DB *DB
}

func New() (*Container, error) {
	if ut.GetEnv("APP_ENV", "development") != "production" {
		if _, err := os.Stat(".env"); err == nil {
			err := godotenv.Load()
			if err != nil {
				return nil, err
			}
		}
	}
	db := &DB{
		Host:     ut.GetEnv("DB_HOST", ""),
		Port:     ut.GetEnv("DB_PORT", ""),
		Username: ut.GetEnv("DB_USERNAME", ""),
		Password: ut.GetEnv("DB_PASSWORD", ""),
		DBName:   ut.GetEnv("DB_NAME", ""),
	}

    // http := &HTTP{
	// 	Env:            ut.GetEnv("APP_ENV", ""),
	// 	URL:            ut.GetEnv("HTTP_URL", ""),
	// 	Port:           ut.GetEnv("HTTP_PORT", ""),
	// 	AllowedOrigins: ut.GetEnv("HTTP_ALLOWED_ORIGINS", ""),
	// }

	return &Container{
		DB: db,
		// HTTP:   http,

	}, nil
}
