package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/Nerzal/gocloak/v13"
	"github.com/lucasgarciaf/df-backend-go/config"
)

func GetKeycloakToken(email, password string) (string, error) {
	client := gocloak.NewClient(config.KeycloakURL)
	ctx := context.Background()

	token, err := client.Login(ctx, config.KeycloakClientID, config.KeycloakClientSecret, config.KeycloakRealm, email, password)
	if err != nil {
		log.Printf("Login to Keycloak failed: %v", err)
		return "", fmt.Errorf("login to Keycloak failed: %w", err)
	}

	return token.AccessToken, nil
}

func GetClientToken() (string, error) {
	client := gocloak.NewClient(config.KeycloakURL)
	ctx := context.Background()

	token, err := client.LoginClient(ctx, config.KeycloakClientID, config.KeycloakClientSecret, config.KeycloakRealm)
	if err != nil {
		log.Printf("Login to Keycloak failed: %v", err)
		return "", fmt.Errorf("login to Keycloak failed: %w", err)
	}

	return token.AccessToken, nil
}
