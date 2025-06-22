package config

import (
	"log"
	"os"
	"time"
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
	log.Println("Loading config from environment")

	cfgToken := ConfigToken{
		TokenKey:    getEnv("TOKEN_PRIVATE_KEY", "your-secret-key-here"),
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
