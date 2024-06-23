package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

var (
	TokenExpiry             = getEnvAsDuration("TOKEN_EXPIRY", time.Hour*24)
	JWTSecretKey            = []byte(getEnv("JWT_SECRET_KEY", "default-secret"))
	MongoDBURI              = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	DatabaseName            = getEnv("DATABASE_NAME", "drivefluency")
	KEYCLOAK_URL            = getEnv("KEYCLOAK_URL", "http://keycloak:8080/")
	KEYCLOAK_REALM          = getEnv("KEYCLOAK_REALM", "drivefluency")
	KEYCLOAK_CLIENT_ID      = getEnv("KEYCLOAK_CLIENT_ID", "df-backend-go")
	KEYCLOAK_CLIENT_SECRET  = getEnv("KEYCLOAK_CLIENT_SECRET", "JWGNgb3pasLTIAPdXSznyMeoVBM0qWdh")
	KEYCLOAK_ADMIN          = getEnv("KEYCLOAK_ADMIN", "admin")
	KEYCLOAK_ADMIN_PASSWORD = getEnv("KEYCLOAK_ADMIN_PASSWORD", "admin")
)

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Invalid duration for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return value
}
