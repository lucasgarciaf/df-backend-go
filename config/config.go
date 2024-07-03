package config

import (
	"log"
	"os"
	"time"
)

var (
	MongoDBURI           = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	DatabaseName         = getEnv("DATABASE_NAME", "drivefluency")
	KeycloakURL          = getEnv("KEYCLOAK_URL", "http://localhost:8080")
	KeycloakRealm        = getEnv("KEYCLOAK_REALM", "myrealm")
	KeycloakClientID     = getEnv("KEYCLOAK_CLIENT_ID", "myclient")
	KeycloakClientSecret = getEnv("KEYCLOAK_CLIENT_SECRET", "mysecret")
	JWTSecretKey         = []byte(getEnv("JWT_SECRET_KEY", "myjwtsecret"))
	TokenExpiry          = time.Hour * 24 // Token expiry duration set to 24 hours
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if defaultValue == "" {
		log.Fatalf("Environment variable %s is not set and no default value provided", key)
	}
	return defaultValue
}
