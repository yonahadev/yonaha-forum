// config/config.go

package config

import (
	"fmt"

	"os"

	"github.com/joho/godotenv"
)

func Init() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbSSLMode)

	return connStr, nil
}
