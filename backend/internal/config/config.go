package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds everything the backend needs to boot.
type Config struct {
	DatabaseURL string
	ServerAddr  string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on process environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required but not set")
	}

	return Config{
		DatabaseURL: dbURL,
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}