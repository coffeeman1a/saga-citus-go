package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const (
	StatusReserved = "pending"
)

var (
	Port        string
	DatabaseURL string
	LogLevel    string
	LogFormat   string
)

func Init() {
	_ = godotenv.Load()

	Port = getEnv("PORT", "8080")
	DatabaseURL = getEnv("DATABASE_URL", "postgres://username:password@localhost:5432/database_name")
	LogFormat = getEnv("LOG_FROMAT", "text")
	LogLevel = getEnv("LOG_LEVEL", "debug")

	initLogger()
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func initLogger() {
	if LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	if LogLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}
}
