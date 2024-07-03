package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/lucasgarciaf/df-backend-go/config"
)

var client gocloak.GoCloak

func InitKeycloak() error {
	client = *gocloak.NewClient(config.KeycloakURL)
	token, err := client.LoginClient(context.TODO(), config.KeycloakClientID, config.KeycloakClientSecret, config.KeycloakRealm)
	if err != nil {
		return fmt.Errorf("login to keycloak failed: %w", err)
	}
	fmt.Printf("Successfully logged into Keycloak with token: %v\n", token.AccessToken)
	return nil
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		rptResult, err := client.RetrospectToken(context.TODO(), tokenStr, config.KeycloakClientID, config.KeycloakClientSecret, config.KeycloakRealm)
		if err != nil || !*rptResult.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Next()
	}
}
