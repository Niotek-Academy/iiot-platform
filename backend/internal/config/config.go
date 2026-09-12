package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	ServerAddr     string
	JWTSecret      string
	JWTExpiryHours int

	// Phase 4 — telemetry ingestion
	FactoryIOMode    string // "simulator" or "opcua"
	OPCUAEndpoint    string
	BufferCapacity   int
	SnapshotInterval time.Duration

	// Phase 6 - Ai
	AIServiceMode        string // "mock" or "http"
	AIServiceURL         string
	AIEvalIntervalSeconds int
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

		FactoryIOMode:    getEnv("FACTORYIO_MODE", "simulator"),
		OPCUAEndpoint:    getEnv("OPCUA_ENDPOINT", "opc.tcp://localhost:4840"),
		BufferCapacity:   getEnvInt("BUFFER_CAPACITY", 30),
		SnapshotInterval: time.Duration(getEnvInt("SNAPSHOT_INTERVAL_SECONDS", 1)) * time.Second,

		AIServiceMode:         getEnv("AI_SERVICE_MODE", "mock"),
		AIServiceURL:          getEnv("AI_SERVICE_URL", "http://localhost:9000"),
		AIEvalIntervalSeconds: getEnvInt("AI_EVAL_INTERVAL_SECONDS", 5),
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