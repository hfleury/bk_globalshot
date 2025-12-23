package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	DbDsn     string
	ServerPort string
	CfgToken  ConfigToken
}

type ConfigToken struct {
	TokenKey    string
	TokenExpiry time.Duration
}

func LoadConfig() Config {
	log.Println("Loading config from environment")

	cfgToken := ConfigToken{
		TokenKey:    getEnv("TOKEN_PRIVATE_KEY", "8a23b8605a2b0a753cc84e3e8154833d3d82039b97bc124d2f4ca17d1590df88e881b06f9cc476dbbb3ba97337dd6e4626d53b6c36b2178da1824ea4ee61e6d8"),
		TokenExpiry: time.Hour * 72,
	}

	return Config{
		DbDsn:      getEnv("DB_DSN", "user=globalshotuser password=globalshotsecret dbname=globalshotdb sslmode=disable host=127.0.0.1 port=5432"),
		ServerPort: getEnv("PORT", "8080"),
		CfgToken:   cfgToken,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
