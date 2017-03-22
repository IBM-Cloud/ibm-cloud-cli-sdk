package authentication

import (
	"encoding/base64"
	"fmt"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
)

type UAAError struct {
	ErrorCode   string `json:"error"`
	Description string `json:"error_description"`
}

type UAATokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type UAARepository struct {
	config core_config.ReadWriter
	client *rest.Client
}

func NewUAARepository(config core_config.ReadWriter, client *rest.Client) *UAARepository {
	return &UAARepository{
		config: config,
		client: client,
	}
}

func (auth *UAARepository) AuthenticatePassword(username string, password string) error {
	return auth.getToken(map[string]string{
		"grant_type": "password",
		"username":   username,
		"password":   password,
	})
}

func (auth *UAARepository) AuthenticateSSO(passcode string) error {
	return auth.getToken(map[string]string{
		"grant_type": "password",
		"passcode":   passcode,
	})
}

func (auth *UAARepository) AuthenticateAPIKey(apiKey string) error {
	return auth.AuthenticatePassword("apikey", apiKey)
}

func (auth *UAARepository) RefreshToken() (string, error) {
	err := auth.getToken(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": auth.config.UAARefreshToken(),
	})
	if err != nil {
		return "", err
	}

	return auth.config.UAAToken(), nil
}

func (auth *UAARepository) getToken(data map[string]string) error {
	request := rest.PostRequest(auth.config.AuthenticationEndpoint()+"/oauth/token").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:"))).
		Field("scope", "")

	for k, v := range data {
		request.Field(k, v)
	}

	var tokens UAATokenResponse
	var apiErr UAAError

	resp, err := auth.client.Do(request, &tokens, &apiErr)
	if err != nil {
		return err
	}
	if apiErr.ErrorCode != "" {
		if apiErr.ErrorCode == "invalid-token" {
			return NewInvalidTokenError(apiErr.Description)
		}
		return NewServerError(resp.StatusCode, apiErr.ErrorCode, apiErr.Description)
	}

	uaaToken := fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken)
	uaaRefreshToken := tokens.RefreshToken

	auth.config.SetUAAToken(uaaToken)
	auth.config.SetUAARefreshToken(uaaRefreshToken)

	return nil
}
