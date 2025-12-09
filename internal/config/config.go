package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	StorageType      string
	B2AccountID      string
	B2ApplicationKey string
	B2BucketName     string
	BunnyZoneName    string
	BunnyAccessKey   string
	BunnyReadOnlyKey string
	BunnyEndpoint    string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		Port:             getEnv("PORT", "8080"),
		StorageType:      getEnv("STORAGE_TYPE", "b2"),
		B2AccountID:      getEnv("B2_ACCOUNT_ID", ""),
		B2ApplicationKey: getEnv("B2_APPLICATION_KEY", ""),
		B2BucketName:     getEnv("B2_BUCKET_NAME", ""),
		BunnyZoneName:    getEnv("BUNNY_ZONE_NAME", ""),
		BunnyAccessKey:   getEnv("BUNNY_ACCESS_KEY", ""),
		BunnyReadOnlyKey: getEnv("BUNNY_READ_ONLY_KEY", ""),
		BunnyEndpoint:    getEnv("BUNNY_ENDPOINT", "de"), // Default to Falkenstein (de)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(value)
	}
	return fallback
}
