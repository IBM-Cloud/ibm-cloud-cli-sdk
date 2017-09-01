package authentication

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"

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

type IAMAuthRepository struct {
	config core_config.Reader
	client *rest.Client
}

func NewIAMAuthRepository(config core_config.Reader, client *rest.Client) *IAMAuthRepository {
	return &IAMAuthRepository{
		config: config,
		client: client,
	}
}

func (auth *IAMAuthRepository) AuthenticatePassword(username string, password string) (Token, Token, error) {
	data := map[string]string{
		"grant_type": "password",
		"username":   username,
		"password":   password,
	}
	return auth.getToken(data)
}

func (auth *IAMAuthRepository) AuthenticateAPIKey(apiKey string) (Token, Token, error) {
	data := map[string]string{
		"grant_type": "urn:ibm:params:oauth:grant-type:apikey",
		"apikey":     apiKey,
	}
	return auth.getToken(data)
}

func (auth *IAMAuthRepository) AuthenticateSSO(passcode string) (Token, Token, error) {
	data := map[string]string{
		"grant_type": "urn:ibm:params:oauth:grant-type:passcode",
		"passcode":   passcode,
	}
	return auth.getToken(data)
}

func (auth *IAMAuthRepository) RefreshToken(refreshToken string) (Token, Token, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	return auth.getToken(data)
}

func (auth *IAMAuthRepository) LinkAccounts(refreshToken string, accounts core_config.AccountsInfo) (Token, Token, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"bss_account":   accounts.AccountID,
		"ims_account":   accounts.IMSAccountID,
	}
	return auth.getToken(data)
}

func (auth *IAMAuthRepository) getToken(data map[string]string) (Token, Token, error) {
	tokenRequest := rest.PostRequest(IAMTokenEndpoint(auth.config.APIEndpoint())+"/oidc/token").
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

	var apiErr IAMError
	resp, err := auth.client.Do(tokenRequest, &tokenResponse, &apiErr)
	if err != nil {
		return Token{}, Token{}, err
	}

	if apiErr.ErrorCode != "" {
		if apiErr.ErrorCode == "BXNIM0407E" {
			return Token{}, Token{}, NewInvalidTokenError(apiErr.detailsOrMessage())
		}
		return Token{}, Token{}, NewServerError(resp.StatusCode, apiErr.ErrorCode, apiErr.detailsOrMessage())
	}

	iamToken := Token{
		AccessToken:  fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.AccessToken),
		RefreshToken: tokenResponse.RefreshToken,
	}
	uaaToken := Token{
		AccessToken:  fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.UAAAccessToken),
		RefreshToken: tokenResponse.UAARefreshToken,
	}
	return iamToken, uaaToken, nil
}

var domainRegExp = regexp.MustCompile(`(^https?://)?[^\.]+(\..+)+`)

func IAMTokenEndpoint(apiEndpoint string) string {
	if apiEndpoint == "" {
		return ""
	}

	endpoint := domainRegExp.ReplaceAllString(apiEndpoint, "${1}iam${2}")

	u, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}

	u.Scheme = "https"
	return u.String()
}
