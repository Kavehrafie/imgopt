package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	B2AccountID      string
	B2ApplicationKey string
	B2BucketName     string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		Port:             getEnv("PORT", "8080"),
		B2AccountID:      getEnv("B2_ACCOUNT_ID", ""),
		B2ApplicationKey: getEnv("B2_APPLICATION_KEY", ""),
		B2BucketName:     getEnv("B2_BUCKET_NAME", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
