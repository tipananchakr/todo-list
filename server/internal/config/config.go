package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	MongoURI       string
	DatabaseName   string
	TodoCollection string
	UserCollection string
	JWTSecret      string
}

func Load() (Config, error) {
	_ = godotenv.Load(".env")

	cfg := Config{
		Port:           getEnv("PORT", "3000"),
		MongoURI:       os.Getenv("MONGODB_URI"),
		DatabaseName:   getEnv("MONGODB_DATABASE", "golang_db"),
		TodoCollection: getEnv("MONGODB_TODO_COLLECTION", "todos"),
		UserCollection: getEnv("MONGODB_USER_COLLECTION", "users"),
		JWTSecret:      getEnv("JWT_SECRET", "local-dev-secret"),
	}

	if cfg.MongoURI == "" {
		return Config{}, fmt.Errorf("MONGODB_URI is required")
	}

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
