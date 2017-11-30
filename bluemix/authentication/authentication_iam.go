package authentication

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
)

const (
	iamClientID     = "bx"
	iamClientSecret = "bx"
)

type IAMError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	ErrorDetails string `json:"errorDetails"`
}

func (e IAMError) detailsOrMessage() string {
	if e.ErrorDetails != "" {
		return e.ErrorDetails
	}
	return e.ErrorMessage
}

type IAMAuthRepository interface {
	AuthenticatePassword(username string, password string) (iamToken, uaaToken Token, err error)
	AuthenticateSSO(passcode string) (iamToken, uaaToken Token, err error)
	AuthenticateAPIKey(apiKey string) (iamToken, uaaToken Token, err error)
	RefreshToken(refreshToken string) (iamToken, uaaToken Token, err error)
	LinkAccounts(refreshToken string, accounts core_config.AccountsInfo) (iamToken, uaaToken Token, err error)
}

type IAMConfig struct {
	// the token endpoint. for example: https://iam.example.com/indentity/token
	TokenEndpoint string
	// client ID and secret may be configurable in future
	// ClientID      string
	// ClientSecret  string
}

type iamAuthRepository struct {
	config IAMConfig
	client *rest.Client
}

func NewIAMAuthRepository(config IAMConfig, client *rest.Client) IAMAuthRepository {
	return &iamAuthRepository{
		config: config,
		client: client,
	}
}

func (auth *iamAuthRepository) AuthenticatePassword(username string, password string) (Token, Token, error) {
	data := map[string]string{
		"grant_type": "password",
		"username":   username,
		"password":   password,
	}
	return auth.getToken(data)
}

func (auth *iamAuthRepository) AuthenticateAPIKey(apiKey string) (Token, Token, error) {
	data := map[string]string{
		"grant_type": "urn:ibm:params:oauth:grant-type:apikey",
		"apikey":     apiKey,
	}
	return auth.getToken(data)
}

func (auth *iamAuthRepository) AuthenticateSSO(passcode string) (Token, Token, error) {
	data := map[string]string{
		"grant_type": "urn:ibm:params:oauth:grant-type:passcode",
		"passcode":   passcode,
	}
	return auth.getToken(data)
}

func (auth *iamAuthRepository) RefreshToken(refreshToken string) (Token, Token, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	return auth.getToken(data)
}

func (auth *iamAuthRepository) LinkAccounts(refreshToken string, accounts core_config.AccountsInfo) (Token, Token, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"bss_account":   accounts.AccountID,
		"ims_account":   accounts.IMSAccountID,
	}
	return auth.getToken(data)
}

func (auth *iamAuthRepository) getToken(data map[string]string) (Token, Token, error) {
	tokenRequest := rest.PostRequest(auth.config.TokenEndpoint).
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", iamClientID, iamClientSecret)))).
		Field("response_type", "cloud_iam,uaa").
		Field("uaa_client_id", "cf").
		Field("uaa_client_secret", "")

	for k, v := range data {
		tokenRequest.Field(k, v)
	}

	var tokenResponse = struct {
		AccessToken     string `json:"access_token"`
		RefreshToken    string `json:"refresh_token"`
		UAAAccessToken  string `json:"uaa_token"`
		UAARefreshToken string `json:"uaa_refresh_token"`
		TokenType       string `json:"token_type"`
	}{}

	if err := auth.sendRequest(tokenRequest, &tokenResponse); err != nil {
		return Token{}, Token{}, err
	}

	iamToken := Token{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		TokenType:    tokenResponse.TokenType,
	}
	uaaToken := Token{
		AccessToken:  tokenResponse.UAAAccessToken,
		RefreshToken: tokenResponse.UAARefreshToken,
		TokenType:    tokenResponse.TokenType,
	}
	return iamToken, uaaToken, nil
}

func (auth *iamAuthRepository) sendRequest(req *rest.Request, respV interface{}) error {
	_, err := auth.client.Do(req, respV, nil)
	switch err := err.(type) {
	case *rest.ErrorResponse:
		var apiErr IAMError
		if e := json.Unmarshal([]byte(err.Message), &apiErr); e == nil {
			switch apiErr.ErrorCode {
			case "":
			case "BXNIM0407E":
				return NewInvalidTokenError(apiErr.detailsOrMessage())
			default:
				return NewServerError(err.StatusCode, apiErr.ErrorCode, apiErr.detailsOrMessage())
			}
		}
	}
	return err
}
