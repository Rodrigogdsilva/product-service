package config

import "os"

type Config struct {
	ListenAddr     string
	JWTSecret      string
	InternalAPIKey string
	DatabaseURL    string
}

func Load() *Config {
	return &Config{
		ListenAddr:     getEnv("LISTEN_ADDR", ":8083"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		InternalAPIKey: getEnv("INTERNAL_API_KEY", ""),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
