package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const (
	StatusPending = "pending"
)

var (
	Port            string
	RedisURL        string
	LogLevel        string
	LogFormat       string
	UsersService    string
	PaymentsService string
	OrdersService   string
)

func Init() {
	_ = godotenv.Load()

	Port = getEnv("PORT", "8080")
	RedisURL = getEnv("REDIS_URL", "redis://default:password@redis/0")
	UsersService = getEnv("USERS_URL", "http://saga_users-service:8081")
	OrdersService = getEnv("ORDERS_URL", "http://saga_orders-service:8082")
	PaymentsService = getEnv("PAYMENTS_URL", "http://saga_payments-service:8083")
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
