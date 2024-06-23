package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	kconfig "github.com/lucasgarciaf/df-backend-go/config"
	"golang.org/x/oauth2"
)

var (
	provider *oidc.Provider
	config   oauth2.Config
	verifier *oidc.IDTokenVerifier
)

func InitKeycloak() error {
	var err error
	ctx := context.Background()
	provider, err = oidc.NewProvider(ctx, kconfig.KEYCLOAK_URL+"realms/"+kconfig.KEYCLOAK_REALM)
	if err != nil {
		return err
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: kconfig.KEYCLOAK_CLIENT_ID})
	return nil
}

func AuthMiddleware(c *gin.Context) {
	rawToken := c.GetHeader("Authorization")
	if rawToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	rawToken = strings.TrimPrefix(rawToken, "Bearer ")

	ctx := context.Background()
	idToken, err := verifier.Verify(ctx, rawToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token claims"})
		return
	}

	c.Set("user", claims)
	c.Next()
}
