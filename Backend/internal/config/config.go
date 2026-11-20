package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBURL      string
	JWT_SECRET string
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		Port:       get("PORT", "8080"),
		DBURL:      get("DB_URL", "postgres://postgres:postgres@localhost:5432/login_with_go_vue?sslmode=disable"),
		JWT_SECRET: get("JWT_SECRET", "dev-secret"),
	}
}

func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
