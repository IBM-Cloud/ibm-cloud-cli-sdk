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

type UAARepository interface {
	AuthenticatePassword(username string, password string) (Token, error)
	AuthenticatePasswordAndConnectIAM(username string, password string, iamAccessToken string) (Token, error)
	AuthenticateSSO(passcode string) (Token, error)
	AuthenticateSSOAndConnectIAM(passcode string, iamAccessToken string) (Token, error)
	AuthenticateAPIKey(apiKey string) (Token, error)
	AuthenticateWithIAMToken(iamAccessToken string) (Token, error)
	DisconnectIAM(uaaToken string) error
	RefreshToken(uaaRefreshToken string) (Token, error)
}

type uaaRepository struct {
	config core_config.Reader
	client *rest.Client
}

func NewUAARepository(config core_config.Reader, client *rest.Client) UAARepository {
	return &uaaRepository{
		config: config,
		client: client,
	}
}

func (auth *uaaRepository) AuthenticatePassword(username string, password string) (Token, error) {
	return auth.getToken(passwordTokenParams(username, password))
}

func (auth *uaaRepository) AuthenticatePasswordAndConnectIAM(username string, password string, iamAccessToken string) (Token, error) {
	return auth.getTokenAndConnectIAM(passwordTokenParams(username, password), iamAccessToken)
}

func passwordTokenParams(username, password string) map[string]string {
	return map[string]string{
		"grant_type": "password",
		"username":   username,
		"password":   password,
	}
}

func (auth *uaaRepository) AuthenticateSSO(passcode string) (Token, error) {
	return auth.getToken(ssoTokenParams(passcode))
}

func (auth *uaaRepository) AuthenticateSSOAndConnectIAM(passcode string, iamAccessToken string) (Token, error) {
	return auth.getTokenAndConnectIAM(ssoTokenParams(passcode), iamAccessToken)
}

func ssoTokenParams(passcode string) map[string]string {
	return map[string]string{
		"grant_type": "password",
		"passcode":   passcode,
	}
}

func (auth *uaaRepository) AuthenticateAPIKey(apiKey string) (Token, error) {
	return auth.getToken(apiKeyTokenParams(apiKey))
}

func apiKeyTokenParams(apiKey string) map[string]string {
	return map[string]string{
		"grant_type": "password",
		"username":   "apikey",
		"password":   apiKey,
	}
}

func (auth *uaaRepository) AuthenticateWithIAMToken(iamAccessToken string) (Token, error) {
	return auth.getToken(iamGrantTokenParams(iamAccessToken))
}

func iamGrantTokenParams(iamAccessToken string) map[string]string {
	return map[string]string{
		"grant_type": "iam_token",
		"iam_token":  iamAccessToken,
	}
}

func (auth *uaaRepository) RefreshToken(refreshToken string) (Token, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	return auth.getToken(data)
}

func (auth *uaaRepository) getTokenAndConnectIAM(v map[string]string, iamToken string) (Token, error) {
	v["connect_to_iam_token"] = iamToken
	return auth.getToken(v)
}

func (auth *uaaRepository) getToken(data map[string]string) (Token, error) {
	req := rest.PostRequest(auth.config.AuthenticationEndpoint()+"/oauth/token").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", uaaClientID, uaaClientSecret)))).
		Field("scope", "")

	for k, v := range data {
		req.Field(k, v)
	}

	var tokens = struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := auth.sendRequest(req, &tokens); err != nil {
		return Token{}, err
	}

	return Token{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
	}, nil
}

func (auth *uaaRepository) DisconnectIAM(uaaToken string) error {
	req := rest.PostRequest(auth.config.AuthenticationEndpoint()+"/oauth/token/disconnect").
		Set("Authorization", uaaToken)
	return auth.sendRequest(req, nil)
}

func (auth *uaaRepository) sendRequest(req *rest.Request, respV interface{}) error {
	var apiErr UAAError
	resp, err := auth.client.Do(req, respV, &apiErr)
	if err != nil {
		return err
	}

	switch apiErr.ErrorCode {
	case "":
		return nil
	case "invalid-token":
		return NewInvalidTokenError(apiErr.Description)
	default:
		return NewServerError(resp.StatusCode, apiErr.ErrorCode, apiErr.Description)
	}
}
