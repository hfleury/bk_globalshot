package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DbDsn    string
	CfgToken ConfigToken
}

type ConfigToken struct {
	TokenKey    string
	TokenExpiry time.Duration
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfgToken := ConfigToken{
		TokenKey:    getEnv("TOKEN_KEY", "your-secret-key-here"),
		TokenExpiry: time.Hour * 72,
	}

	return Config{
		DbDsn:    getEnv("DB_DSN", "user=globalshotuser password=globalshotsecret dbname=globalshotdb sslmode=disable"),
		CfgToken: cfgToken,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
