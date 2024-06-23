package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lucasgarciaf/df-backend-go/config"
)

var keycloakToken string

// GetClientToken retrieves the Keycloak client token
func GetClientToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", config.KEYCLOAK_CLIENT_ID)
	data.Set("client_secret", config.KEYCLOAK_CLIENT_SECRET)

	req, err := http.NewRequest("POST", fmt.Sprintf("%srealms/%s/protocol/openid-connect/token", config.KEYCLOAK_URL, config.KEYCLOAK_REALM), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to obtain client token: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("failed to parse access token")
	}

	keycloakToken = token
	return keycloakToken, nil
}
