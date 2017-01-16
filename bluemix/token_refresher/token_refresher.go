package token_refresher

import (
	"encoding/base64"
	"fmt"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
)

type InvalidTokenError struct {
	Description string
}

func NewInvalidTokenError(description string) *InvalidTokenError {
	return &InvalidTokenError{Description: description}
}

func (e *InvalidTokenError) Error() string {
	return "Invalid auth token: " + e.Description
}

type TokenRefresher interface {
	RefreshAuthToken() (newToken string, refreshToken string, err error)
}

type uaaTokenRefresher struct {
	config core_config.ReadWriter
	client *rest.Client
}

func NewTokenRefresher(config core_config.ReadWriter) *uaaTokenRefresher {
	return &uaaTokenRefresher{
		config: config,
		client: rest.NewClient(),
	}
}

func (t *uaaTokenRefresher) RefreshAuthToken() (string, string, error) {
	req := rest.PostRequest(fmt.Sprintf("%s/oauth/token", t.config.UaaEndpoint())).
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:"))).
		Field("refresh_token", t.config.RefreshToken()).
		Field("grant_type", "refresh_token").
		Field("scope", "")

	tokens := new(uaaTokenResponse)
	apiErr := new(uaaErrorResponse)
	resp, err := t.client.Do(req, tokens, apiErr)

	if err != nil {
		return "", "", err
	}

	if apiErr.Code != "" {
		if apiErr.Code == "invalid-token" {
			return "", "", NewInvalidTokenError(apiErr.Description)
		} else {
			return "", "", fmt.Errorf("Error response from server. StatusCode: %d; description: %s", resp.StatusCode, apiErr.Description)
		}
	}

	accessToken := fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken)
	refreshToken := tokens.RefreshToken

	t.config.SetAccessToken(accessToken)
	t.config.SetRefreshToken(refreshToken)

	return accessToken, refreshToken, nil
}

type uaaErrorResponse struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

type uaaTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}
