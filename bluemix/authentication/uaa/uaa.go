package uaa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/types"
)

const (
	defaultClientID     = "cf"
	defaultClientSecret = ""
)

// Grant types
const (
	GrantTypePassword          authentication.GrantType = "password"
	GrantTypeRefreshToken      authentication.GrantType = "refresh_token"
	GrantTypeAuthorizationCode authentication.GrantType = "authorization_code"
	GrantTypeIAMToken          authentication.GrantType = "iam_token"
)

// Response types
const (
	ResponseTypeUAA authentication.ResponseType = "uaa"
)

type APIError struct {
	ErrorCode   string `json:"error"`
	Description string `json:"error_description"`
}

func PasswordTokenRequest(username, password string, opts ...authentication.TokenOption) *authentication.TokenRequest {
	r := authentication.NewTokenRequest(GrantTypePassword)
	r.SetTokenParam("username", username)
	r.SetTokenParam("password", password)
	for _, o := range opts {
		r.WithOption(o)
	}
	return r
}

func OnetimePasscodeTokenRequest(passcode string, opts ...authentication.TokenOption) *authentication.TokenRequest {
	r := authentication.NewTokenRequest(GrantTypePassword)
	r.SetTokenParam("passcode", passcode)
	for _, o := range opts {
		r.WithOption(o)
	}
	return r
}

func APIKeyTokenRequest(apiKey string, opts ...authentication.TokenOption) *authentication.TokenRequest {
	r := authentication.NewTokenRequest(GrantTypePassword)
	r.SetTokenParam("username", "apikey")
	r.SetTokenParam("password", apiKey)
	for _, o := range opts {
		r.WithOption(o)
	}
	return r
}

func AuthorizationTokenRequest(code string, redirectURI string, opts ...authentication.TokenOption) *authentication.TokenRequest {
	r := authentication.NewTokenRequest(GrantTypeAuthorizationCode)
	r.SetTokenParam("code", code)
	r.SetTokenParam("redirect_uri", redirectURI)
	for _, o := range opts {
		r.WithOption(o)
	}
	return r
}

func RefreshTokenRequest(refreshToken string, opts ...authentication.TokenOption) *authentication.TokenRequest {
	r := authentication.NewTokenRequest(GrantTypeRefreshToken)
	r.SetTokenParam("refresh_token", refreshToken)
	for _, o := range opts {
		r.WithOption(o)
	}
	return r
}

func ConnectToIAM(iamAccessToken string) authentication.TokenOption {
	return func(r *authentication.TokenRequest) {
		r.SetTokenParam("connect_to_iam_token", iamAccessToken)
	}
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	Expiry       time.Time `json:"expires_in"`
	Scope        string    `json:"scope"`
}

type Interface interface {
	GetToken(req *authentication.TokenRequest) (*Token, error)
	ConnectToIAM(iamAccessToken string) (*Token, error)
	DisconnectIAM(uaaToken string) error
}

type Config struct {
	UAAEndpoint  string
	ClientID     string
	ClientSecret string
}

func DefaultConfig(uaaEndpoint string) Config {
	return Config{
		UAAEndpoint:  uaaEndpoint,
		ClientID:     defaultClientID,
		ClientSecret: defaultClientSecret,
	}
}

type client struct {
	config Config
	client *rest.Client
}

func NewClient(config Config, restClient *rest.Client) Interface {
	return &client{
		config: config,
		client: restClient,
	}
}

func (c *client) GetToken(tokenReq *authentication.TokenRequest) (*Token, error) {
	v := make(url.Values)
	tokenReq.SetValue(v)

	if len(tokenReq.ResponseTypes()) == 0 {
		v.Set("response_type", ResponseTypeUAA.String())
	}

	r := rest.PostRequest(c.config.UAAEndpoint+"/oauth/token").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.ClientID, c.config.ClientSecret)))).
		Field("scope", "")

	for k, ss := range v {
		for _, s := range ss {
			r.Field(k, s)
		}
	}

	var tokenResponse struct {
		Token
		Expiry types.UnixTime `json:"expires_in"`
	}
	if err := c.doRequest(r, &tokenResponse); err != nil {
		return nil, err
	}

	ret := tokenResponse.Token
	ret.Expiry = time.Time(tokenResponse.Expiry)
	return &ret, nil
}

func (c *client) ConnectToIAM(iamAccessToken string) (*Token, error) {
	r := authentication.NewTokenRequest(GrantTypeIAMToken)
	r.SetTokenParam("iam_token", iamAccessToken)
	return c.GetToken(r)
}

func (c *client) DisconnectIAM(uaaToken string) error {
	r := rest.PostRequest(c.config.UAAEndpoint+"/oauth/token/disconnect").
		Set("Authorization", uaaToken)
	return c.doRequest(r, nil)
}

func (c *client) doRequest(req *rest.Request, respV interface{}) error {
	_, err := c.client.Do(req, respV, nil)
	switch err := err.(type) {
	case *rest.ErrorResponse:
		var apiErr APIError
		if e := json.Unmarshal([]byte(err.Message), &apiErr); e == nil {
			switch apiErr.ErrorCode {
			case "":
			case "invalid_grant":
				return authentication.NewInvalidGrantTypeError(apiErr.Description)
			case "invalid-token":
				return authentication.NewInvalidTokenError(apiErr.Description)
			default:
				return authentication.NewServerError(err.StatusCode, apiErr.ErrorCode, apiErr.Description)
			}
		}
	}
	return err
}
