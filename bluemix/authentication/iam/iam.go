package iam

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/types"
)

const (
	defaultClientID        = "bx"
	defaultClientSecret    = "bx"
	defaultUAAClientID     = "cf"
	defaultUAAClientSecret = ""
	crTokenParam           = "cr_token"
	profileIDParam         = "profile_id"
	profileNameParam       = "profile_name"
	profileCRNParam        = "profile_crn"
)

// Grant types
const (
	GrantTypePassword              authentication.GrantType = "password"                                 // #nosec G101 - this the API request grant type. Not a credential
	GrantTypeAPIKey                authentication.GrantType = "urn:ibm:params:oauth:grant-type:apikey"   // #nosec G101 - this the API request grant type. Not a credential
	GrantTypeOnetimePasscode       authentication.GrantType = "urn:ibm:params:oauth:grant-type:passcode" // #nosec G101 - this the API request grant type. Not a credential
	GrantTypeAuthorizationCode     authentication.GrantType = "authorization_code"
	GrantTypeRefreshToken          authentication.GrantType = "refresh_token"
	GrantTypeDelegatedRefreshToken authentication.GrantType = "urn:ibm:params:oauth:grant-type:delegated-refresh-token" // #nosec G101 - this the API request grant type. Not a credential
	GrantTypeIdentityCookie        authentication.GrantType = "urn:ibm:params:oauth:grant-type:identity-cookie"
	GrantTypeDerive                authentication.GrantType = "urn:ibm:params:oauth:grant-type:derive"
	GrantTypeCRToken               authentication.GrantType = "urn:ibm:params:oauth:grant-type:cr-token" // #nosec G101 - this the API request grant type. Not a credential
)

// Response types
const (
	ResponseTypeIAM                   authentication.ResponseType = "cloud_iam"
	ResponseTypeUAA                   authentication.ResponseType = "uaa"
	ResponseTypeIMS                   authentication.ResponseType = "ims_portal"
	ResponseTypeDelegatedRefreshToken authentication.ResponseType = "delegated_refresh_token" // #nosec G101 - this the API response grant type. Not a credential
)

const (
	InvalidTokenErrorCode           = "BXNIM0407E" // #nosec G101 - this an API error response code. Not a credential
	RefreshTokenExpiryErrorCode     = "BXNIM0408E" // #nosec G101 - this an API error response code. Not a credential
	ExternalAuthenticationErrorCode = "BXNIM0400E"
	SessionInactiveErrorCode        = "BXNIM0439E"
)

type MFAVendor string

func (m MFAVendor) String() string {
	return string(m)
}

// MFA vendors
const (
	MFAVendorVerisign    MFAVendor = "VERISIGN"
	MFAVendorTOTP        MFAVendor = "TOTP"
	MFAVendorPhoneFactor MFAVendor = "PHONE_FACTOR"
)

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
	r := authentication.NewTokenRequest(GrantTypeOnetimePasscode)
	r.SetTokenParam("passcode", passcode)
	for _, o := range opts {
		r.WithOption(o)
	}
	return r
}

func APIKeyTokenRequest(apikey string, opts ...authentication.TokenOption) *authentication.TokenRequest {
	r := authentication.NewTokenRequest(GrantTypeAPIKey)
	r.SetTokenParam("apikey", apikey)
	for _, o := range opts {
		r.WithOption(o)
	}
	return r
}

// CRTokenRequest builds a 'TokenRequest' struct from the user input. The value of 'crToken' is set as the value of the 'cr_token' form
// parameter of the request. 'profileID' and 'profileName' are optional parameters used to set the 'profile_id' and 'profile_name' form parameters
// in the request, respectively.
func CRTokenRequest(crToken string, profileID string, profileName string, opts ...authentication.TokenOption) *authentication.TokenRequest {
	return CRTokenRequestWithCRN(crToken, profileID, profileName, "", opts...)
}

// CRTokenRequestWithCRN builds a 'TokenRequest' struct from the user input. The value of 'crToken' is set as the value of the 'cr_token' form
// parameter of the request. 'profileID', 'profileName', and 'profileCRN' are optional parameters used to set the 'profile_id', 'profile_name',
// and 'profile_crn' form parameters in the request, respectively.
func CRTokenRequestWithCRN(crToken string, profileID string, profileName string, profileCRN string, opts ...authentication.TokenOption) *authentication.TokenRequest {
	r := authentication.NewTokenRequest(GrantTypeCRToken)
	r.SetTokenParam(crTokenParam, crToken)

	if profileID != "" {
		r.SetTokenParam(profileIDParam, profileID)
	}
	if profileName != "" {
		r.SetTokenParam(profileNameParam, profileName)
	}
	if profileCRN != "" {
		r.SetTokenParam(profileCRNParam, profileCRN)
	}

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

func AuthorizationTokenRequest(code string, redirectURI string, opts ...authentication.TokenOption) *authentication.TokenRequest {
	r := authentication.NewTokenRequest(GrantTypeAuthorizationCode)
	r.SetTokenParam("code", code)
	r.SetTokenParam("redirect_uri", redirectURI)
	for _, o := range opts {
		r.WithOption(o)
	}
	return r
}

func SetAccount(accountID string) authentication.TokenOption {
	return func(r *authentication.TokenRequest) {
		r.SetTokenParam("account", accountID)
	}
}

func SetIMSAccount(imsAccountID string) authentication.TokenOption {
	return func(r *authentication.TokenRequest) {
		r.SetTokenParam("ims_account", imsAccountID)
	}
}

func SetSecurityQuestion(questionID int, answer string) authentication.TokenOption {
	return func(r *authentication.TokenRequest) {
		r.SetTokenParam("security_question_id", strconv.Itoa(questionID))
		r.SetTokenParam("security_question_answer", answer)
	}
}

func SetVeriSignCode(code string) authentication.TokenOption {
	return SetSecurityCode(code, MFAVendorVerisign)
}

func SetTOTPCode(code string) authentication.TokenOption {
	return SetSecurityCode(code, MFAVendorTOTP)
}

func SetSecurityCode(code string, vendor MFAVendor) authentication.TokenOption {
	return func(r *authentication.TokenRequest) {
		r.SetTokenParam("security_code", code)
		r.SetTokenParam("vendor", vendor.String())
	}
}

func SetPhoneAuthToken(token string) authentication.TokenOption {
	return func(r *authentication.TokenRequest) {
		r.SetTokenParam("authentication_token", token)
		r.SetTokenParam("vendor", MFAVendorPhoneFactor.String())
	}
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	Scope        string    `json:"scope"`
	Expiry       time.Time `json:"expiration"`

	// Fields present when ResponseTypeUAA is set
	UAAToken        string `json:"uaa_token"`
	UAARefreshToken string `json:"uaa_refresh_token"`

	// Fields present when ResponseTypeIMS is set
	IMSUserID int64  `json:"ims_user_id"`
	IMSToken  string `json:"ims_token"`
}

type APIError struct {
	ErrorCode    string      `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
	ErrorDetails string      `json:"errorDetails"`
	Requirements Requirement `json:"requirements"`
}

type Requirement struct {
	ErrorCode    string `json:"code"`
	ErrorMessage string `json:"error"`
}

func (e APIError) errorMessage() string {
	if e.ErrorDetails != "" {
		return e.ErrorDetails
	}
	return e.ErrorMessage
}

type Endpoint struct {
	AuthURL     string `json:"authorization_endpoint"`
	TokenURL    string `json:"token_endpoint"`
	PasscodeURL string `json:"passcode_endpoint"`
}

type Interface interface {
	GetEndpoint() (*Endpoint, error)
	GetToken(req *authentication.TokenRequest) (*Token, error)
	InitiateIMSPhoneFactor(req *authentication.TokenRequest) (authToken string, err error)
}

type Config struct {
	IAMEndpoint     string
	TokenEndpoint   string // Optional. Default value is <IAMEndpoint>/identity/token
	ClientID        string
	ClientSecret    string
	UAAClientID     string
	UAAClientSecret string
}

func (c Config) tokenEndpoint() string {
	if c.TokenEndpoint != "" {
		return c.TokenEndpoint
	}
	return c.IAMEndpoint + "/identity/token"
}

func DefaultConfig(iamEndpoint string) Config {
	return Config{
		IAMEndpoint:     iamEndpoint,
		TokenEndpoint:   iamEndpoint + "/identity/token",
		ClientID:        defaultClientID,
		ClientSecret:    defaultClientSecret,
		UAAClientID:     defaultUAAClientID,
		UAAClientSecret: defaultUAAClientSecret,
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

	responseTypes := tokenReq.ResponseTypes()
	for _, t := range responseTypes {
		if t == ResponseTypeUAA {
			v.Set("uaa_client_id", c.config.UAAClientID)
			v.Set("uaa_client_secret", c.config.UAAClientSecret)
			break
		}
	}
	if len(responseTypes) == 0 {
		v.Set("response_type", ResponseTypeIAM.String())
	}

	r := rest.PostRequest(c.config.tokenEndpoint()).
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.ClientID, c.config.ClientSecret))))

	for k, ss := range v {
		for _, s := range ss {
			r.Field(k, s)
		}
	}

	var tokenResponse struct {
		Token
		Expiry types.UnixTime `json:"expiration"`
	}
	if err := c.doRequest(r, &tokenResponse); err != nil {
		return nil, err
	}

	ret := tokenResponse.Token
	ret.Expiry = time.Time(tokenResponse.Expiry)
	return &ret, nil
}

func (c *client) InitiateIMSPhoneFactor(tokenReq *authentication.TokenRequest) (authToken string, err error) {
	v := make(url.Values)
	tokenReq.SetValue(v)

	r := rest.PostRequest(c.config.IAMEndpoint+"/identity/initiate_ims_2fa").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.ClientID, c.config.ClientSecret))))

	for k, ss := range v {
		for _, s := range ss {
			r.Field(k, s)
		}
	}

	var resp struct {
		AuthenticationToken string `json:"authenticationToken"`
	}
	if err := c.doRequest(r, &resp); err != nil {
		return "", err
	}
	return resp.AuthenticationToken, nil
}

func (c *client) GetEndpoint() (*Endpoint, error) {
	var e Endpoint
	err := c.doRequest(rest.GetRequest(c.config.IAMEndpoint+"/identity/.well-known/openid-configuration"), &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (c *client) doRequest(r *rest.Request, respV interface{}) error {
	_, err := c.client.Do(r, respV, nil)
	switch err := err.(type) {
	case *rest.ErrorResponse:
		var apiErr APIError
		if jsonErr := json.Unmarshal([]byte(err.Message), &apiErr); jsonErr == nil {
			switch apiErr.ErrorCode {
			case "":
			case InvalidTokenErrorCode:
				return authentication.NewInvalidTokenError(apiErr.errorMessage())
			case RefreshTokenExpiryErrorCode:
				return authentication.NewRefreshTokenExpiryError(apiErr.errorMessage())
			case ExternalAuthenticationErrorCode:
				return &authentication.ExternalAuthenticationError{ErrorCode: apiErr.Requirements.ErrorCode, ErrorMessage: apiErr.Requirements.ErrorMessage}
			case SessionInactiveErrorCode:
				return authentication.NewSessionInactiveError(apiErr.errorMessage())
			default:
				return authentication.NewServerError(err.StatusCode, apiErr.ErrorCode, apiErr.errorMessage())
			}
		}
	}
	return err
}
