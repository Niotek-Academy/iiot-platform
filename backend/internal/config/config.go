package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds everything the backend needs to boot.
type Config struct {
	DatabaseURL string
	ServerAddr  string
	JWTSecret      string
	JWTExpiryHours int
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on process environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required but not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required but not set")
	}

	return Config{
		DatabaseURL:    dbURL,
		ServerAddr:     getEnv("SERVER_ADDR", ":8080"),
		JWTSecret:      jwtSecret,
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 12),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}