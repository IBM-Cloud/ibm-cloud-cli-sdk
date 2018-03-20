package authentication

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

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
	AuthenticatePassword(username string, password string) (iamToken Token, err error)
	AuthenticateSSO(passcode string) (iamToken Token, err error)
	AuthenticateAPIKey(apiKey string) (iamToken Token, err error)
	GetUAAToken(iamAccessToken string) (uaaToken Token, err error)
	RefreshToken(refreshToken string) (iamToken Token, err error)
	RefreshTokenToLinkAccounts(refreshToken string, accounts core_config.AccountsInfo) (iamToken Token, err error)
	RefreshTokenToLinkAccountsAndGetUAAToken(refreshToken string, accounts core_config.AccountsInfo) (iamToken, uaaToken Token, err error)
}

type IAMConfig struct {
	// the token endpoint. for example: https://iam.example.com/indentity/token
	TokenEndpoint string
	// client ID and secret may be configurable in future
	// ClientID      string
	// ClientSecret  string
}

type iamAuthRepository struct {
	config *IAMConfig
	client *rest.Client
}

func NewIAMAuthRepository(config *IAMConfig, client *rest.Client) IAMAuthRepository {
	return &iamAuthRepository{
		config: config,
		client: client,
	}
}

func (auth *iamAuthRepository) AuthenticatePassword(username string, password string) (Token, error) {
	return auth.getIAMToken("password", map[string]string{"username": username, "password": password})
}

func (auth *iamAuthRepository) AuthenticateAPIKey(apiKey string) (Token, error) {
	return auth.getIAMToken("urn:ibm:params:oauth:grant-type:apikey", map[string]string{"apikey": apiKey})
}

func (auth *iamAuthRepository) AuthenticateSSO(passcode string) (Token, error) {
	return auth.getIAMToken("urn:ibm:params:oauth:grant-type:passcode", map[string]string{"passcode": passcode})
}

func (auth *iamAuthRepository) RefreshToken(refreshToken string) (Token, error) {
	return auth.getIAMToken("refresh_token", map[string]string{"refresh_token": refreshToken})
}

func (auth *iamAuthRepository) RefreshTokenToLinkAccounts(refreshToken string, accounts core_config.AccountsInfo) (Token, error) {
	return auth.getIAMToken("refresh_token", map[string]string{
		"refresh_token": refreshToken,
		"bss_account":   accounts.AccountID,
		"ims_account":   accounts.IMSAccountID,
	})
}

func (auth *iamAuthRepository) getIAMToken(grantType string, data map[string]string) (Token, error) {
	r := tokenRequest{
		iamTokenRequired: true,
		grantType:        grantType,
		data:             data,
	}

	tokens, err := auth.getToken(r)
	if err != nil {
		return Token{}, err
	}
	return tokens.iamToken(), nil
}

func (auth *iamAuthRepository) GetUAAToken(iamAccessToken string) (Token, error) {
	r := tokenRequest{
		uaaTokenRequired: true,
		grantType:        "urn:ibm:params:oauth:grant-type:derive",
		data:             map[string]string{"access_token": iamAccessToken},
	}

	tokens, err := auth.getToken(r)
	if err != nil {
		return Token{}, err
	}
	return tokens.uaaToken(), nil
}

func (auth *iamAuthRepository) RefreshTokenToLinkAccountsAndGetUAAToken(refreshToken string, accounts core_config.AccountsInfo) (Token, Token, error) {
	r := tokenRequest{
		iamTokenRequired: true,
		uaaTokenRequired: true,
		grantType:        "refresh_token",
		data:             map[string]string{"refresh_token": refreshToken},
	}

	tokens, err := auth.getToken(r)
	if err != nil {
		return Token{}, Token{}, err
	}
	return tokens.iamToken(), tokens.uaaToken(), nil
}

type tokenRequest struct {
	iamTokenRequired bool
	uaaTokenRequired bool
	grantType        string
	data             map[string]string
}

type tokenResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	UAAAccessToken  string `json:"uaa_token"`
	UAARefreshToken string `json:"uaa_refresh_token"`
	TokenType       string `json:"token_type"`
}

func (res tokenResponse) iamToken() Token {
	return Token{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		TokenType:    res.TokenType,
	}
}

func (res tokenResponse) uaaToken() Token {
	return Token{
		AccessToken:  res.UAAAccessToken,
		RefreshToken: res.UAARefreshToken,
		TokenType:    res.TokenType,
	}
}

func (auth *iamAuthRepository) getToken(r tokenRequest) (tokenResponse, error) {
	req := rest.PostRequest(auth.config.TokenEndpoint)
	req.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", iamClientID, iamClientSecret))))

	var grantTypes []string
	if r.iamTokenRequired {
		grantTypes = append(grantTypes, "cloud_iam")
	}
	if r.uaaTokenRequired {
		grantTypes = append(grantTypes, "uaa")
	}
	req.Field("response_type", strings.Join(grantTypes, ","))

	if r.uaaTokenRequired {
		req.Field("uaa_client_id", "cf").
			Field("uaa_client_secret", "")
	}

	req.Field("grant_type", r.grantType)

	for k, v := range r.data {
		req.Field(k, v)
	}

	var tokens tokenResponse
	err := auth.sendRequest(req, &tokens)
	if err != nil {
		return tokenResponse{}, err
	}
	return tokens, nil
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
