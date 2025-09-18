package config

import "os"

type Config struct {
	ListenAddr     string
	InternalAPIKey string
	DatabaseURL    string
	AuthServiceURL string
}

func Load() *Config {
	return &Config{
		ListenAddr:     getEnv("LISTEN_ADDR", ":8083"),
		InternalAPIKey: getEnv("INTERNAL_API_KEY", ""),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
