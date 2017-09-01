package authentication

import (
	"encoding/base64"
	"fmt"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
)

const (
	uaaClientID     = "cf"
	uaaClientSecret = ""
)

type UAAError struct {
	ErrorCode   string `json:"error"`
	Description string `json:"error_description"`
}

type UAARepository struct {
	config core_config.Reader
	client *rest.Client
}

func NewUAARepository(config core_config.Reader, client *rest.Client) *UAARepository {
	return &UAARepository{
		config: config,
		client: client,
	}
}

func (auth *UAARepository) AuthenticatePassword(username string, password string) (Token, error) {
	data := map[string]string{
		"grant_type": "password",
		"username":   username,
		"password":   password,
	}
	return auth.getToken(data)
}

func (auth *UAARepository) AuthenticateSSO(passcode string) (Token, error) {
	data := map[string]string{
		"grant_type": "password",
		"passcode":   passcode,
	}
	return auth.getToken(data)
}

func (auth *UAARepository) AuthenticateAPIKey(apiKey string) (Token, error) {
	return auth.AuthenticatePassword("apikey", apiKey)
}

func (auth *UAARepository) RefreshToken(refreshToken string) (Token, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	return auth.getToken(data)
}

func (auth *UAARepository) getToken(data map[string]string) (Token, error) {
	tokenRequest := rest.PostRequest(auth.config.AuthenticationEndpoint()+"/oauth/token").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", uaaClientID, uaaClientSecret)))).
		Field("scope", "")

	for k, v := range data {
		tokenRequest.Field(k, v)
	}

	var tokenResponse = struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
	}{}

	var apiErr UAAError
	resp, err := auth.client.Do(tokenRequest, &tokenResponse, &apiErr)
	if err != nil {
		return Token{}, err
	}
	if apiErr.ErrorCode != "" {
		if apiErr.ErrorCode == "invalid-token" {
			return Token{}, NewInvalidTokenError(apiErr.Description)
		}
		return Token{}, NewServerError(resp.StatusCode, apiErr.ErrorCode, apiErr.Description)
	}

	token := Token{
		AccessToken:  fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.AccessToken),
		RefreshToken: tokenResponse.RefreshToken,
	}

	return token, nil
}
