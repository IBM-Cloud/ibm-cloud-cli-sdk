package authentication

import (
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
)

type IAMError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	ErrorDetails string `json:"errorDetails"`
}

func (e IAMError) Description() string {
	if e.ErrorDetails != "" {
		return e.ErrorDetails
	}
	return e.ErrorMessage
}

type IAMTokenResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	UAAAccessToken  string `json:"uaa_token"`
	UAARefreshToken string `json:"uaa_refresh_token"`
	TokenType       string `json:"token_type"`
}

type IAMAuthRepository interface {
	AuthenticatePassword(username string, password string) error
	AuthenticateAPIKey(apiKey string) error
	RefreshToken() (string, error)
	RefreshTokenWithUpdatedAccounts(core_config.AccountsInfo) (string, error)
}

type iamAuthRepository struct {
	config core_config.ReadWriter
	client *rest.Client
}

func NewIAMAuthRepository(config core_config.ReadWriter, client *rest.Client) IAMAuthRepository {
	return &iamAuthRepository{
		config: config,
		client: client,
	}
}

func (auth *iamAuthRepository) AuthenticatePassword(username string, password string) error {
	return auth.getToken(map[string]string{
		"grant_type": "password",
		"username":   username,
		"password":   password,
	})
}

func (auth *iamAuthRepository) AuthenticateSSO(passcode string) error {
	return auth.getToken(map[string]string{
		"grant_type": "password",
		"passcode":   passcode,
	})
}

func (auth *iamAuthRepository) AuthenticateAPIKey(apiKey string) error {
	return auth.getToken(map[string]string{
		"grant_type": "urn:ibm:params:oauth:grant-type:apikey",
		"apikey":     apiKey,
	})
}

func (auth *iamAuthRepository) RefreshToken() (string, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": auth.config.IAMRefreshToken(),
	}

	err := auth.getToken(data)
	if err != nil {
		return "", err
	}

	return auth.config.IAMToken(), nil
}

func (auth *iamAuthRepository) RefreshTokenWithUpdatedAccounts(updatedAccounts core_config.AccountsInfo) (string, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": auth.config.IAMRefreshToken(),
		"bss_account":   updatedAccounts.AccountID,
		"ims_account":   updatedAccounts.IMSAccountID,
	}

	err := auth.getToken(data)
	if err != nil {
		return "", err
	}

	return auth.config.IAMToken(), nil
}

func (auth *iamAuthRepository) getToken(data map[string]string) error {
	request := rest.PostRequest(IAMTokenEndpoint(auth.config.APIEndpoint())+"/oidc/token").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("bx:bx"))).
		Field("response_type", "cloud_iam,uaa").
		Field("uaa_client_id", "cf").
		Field("uaa_client_secret", "")

	for k, v := range data {
		request.Field(k, v)
	}

	var tokens IAMTokenResponse
	var apiErr IAMError

	resp, err := auth.client.Do(request, &tokens, &apiErr)
	if err != nil {
		return err
	}

	if apiErr.ErrorCode != "" {
		if apiErr.ErrorCode == "BXNIM0407E" {
			return NewInvalidTokenError(apiErr.Description())
		}
		return NewServerError(resp.StatusCode, apiErr.ErrorCode, apiErr.Description())
	}

	auth.config.SetIAMToken(fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	auth.config.SetIAMRefreshToken(tokens.RefreshToken)

	auth.config.SetUAAToken(fmt.Sprintf("%s %s", tokens.TokenType, tokens.UAAAccessToken))
	auth.config.SetUAARefreshToken(tokens.UAARefreshToken)

	return nil
}

func IAMTokenEndpoint(apiEndpoint string) string {
	return regexp.MustCompile(`(^https?://)[^\.]+(\..+)+`).ReplaceAllString(apiEndpoint, "${1}iam${2}")
}
