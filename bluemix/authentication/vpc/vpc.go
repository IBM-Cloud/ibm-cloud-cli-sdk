package vpc

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/iam"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
)

const (
	DefaultMetadataServiceVersion = "2022-06-30"
	DefaultServerEndpoint         = "http://169.254.169.254"

	defaultMetadataFlavor                 = "ibm"
	defaultInstanceIdentityTokenLifetime  = 300
	defaultOperationPathCreateAccessToken = "/instance_identity/v1/token"
	defaultOperationPathCreateIamToken    = "/instance_identity/v1/iam_token"

	// headers
	MetadataFlavor = "Metadata-Flavor"

	contentType   = "Content-Type"
	authorization = "Authorization"
	tokenType     = "Bearer"
)

// InstanceIdentityToken describes the response body for the 'GetInstanceIdentityToken'
// operation.
type InstanceIdentityToken struct {
	// The access token
	AccessToken string `json:"access_token"`
	// The date and time that the access token was created
	CreatedAt string `json:"created_at"`
	// The date and time that the access token will expire
	ExpiresAt string `json:"expires_at"`
	// Time in seconds before the access token expires
	ExpiresIn int `json:"expires_in"`
}

// getIAMAccessTokenResponse describes the response body for the 'GetIAMAccessToken'
// operation.
type getIAMAccessTokenResponse struct {
	// The access token
	AccessToken string `json:"access_token"`
	// The date and time that the access token was created
	CreatedAt string `json:"created_at"`
	// The date and time that the access token will expire
	ExpiresAt string `json:"expires_at"`
	// Time in seconds before the access token expires
	ExpiresIn int `json:"expires_in"`
}

type Interface interface {
	GetInstanceIdentityToken() (*InstanceIdentityToken, error)
	GetIAMAccessToken(req *IAMAccessTokenRequest) (*iam.Token, error)
}

type Config struct {
	VPCAuthEndpoint              string
	InstanceIdentiyTokenEndpoint string
	IAMTokenEndpoint             string
	Version                      string
}

// IAMAccessTokenRequest represents the request object for the `GetIAMAccessToken` operation.
// AccessToken - The instance identity token
// Body - The trusted profile ID/CRN represented as the `TrustedProfileByIdOrCRN` object
type IAMAccessTokenRequest struct {
	AccessToken string `json:"access_token"`
	Body        *TrustedProfileByIdOrCRN
}

type TrustedProfileByIdOrCRN struct {
	ID  string `json:"id"`
	CRN string `json:"crn"`
}

// NewIAMAccessTokenRequest builds the request body for the GetIAMAccessToken
// operation. The request body for this operation consists of either a trusted profile
// ID or CRN. If both a trusted profile ID and a trusted profile CRN are provided, then
// an error is returned.
func NewIAMAccessTokenRequest(profileID string, profileCRN string, token string) (*IAMAccessTokenRequest, error) {
	var trustedProfile TrustedProfileByIdOrCRN

	if profileID != "" && profileCRN != "" {
		return nil, fmt.Errorf("A Profile ID and Profile CRN cannot both be specified.")
	}

	if profileCRN != "" {
		trustedProfile.CRN = profileCRN
	}
	if profileID != "" {
		trustedProfile.ID = profileID
	}

	request := &IAMAccessTokenRequest{
		Body:        &trustedProfile,
		AccessToken: token,
	}

	return request, nil
}

func (c Config) instanceIdentityTokenEndpoint() string {
	if c.InstanceIdentiyTokenEndpoint != "" {
		return c.InstanceIdentiyTokenEndpoint
	}
	return c.VPCAuthEndpoint + defaultOperationPathCreateAccessToken
}

func (c Config) iamTokenEndpoint() string {
	if c.IAMTokenEndpoint != "" {
		return c.IAMTokenEndpoint
	}
	return c.VPCAuthEndpoint + defaultOperationPathCreateIamToken
}

func DefaultConfig(vpcAuthEndpoint string, apiVersion string) Config {
	return Config{
		VPCAuthEndpoint:              vpcAuthEndpoint,
		InstanceIdentiyTokenEndpoint: vpcAuthEndpoint + defaultOperationPathCreateAccessToken,
		IAMTokenEndpoint:             vpcAuthEndpoint + defaultOperationPathCreateIamToken,
		Version:                      apiVersion,
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

// GetInstanceIdentityToken retrieves the VPC/VSI compute resource's instance identity token specified at
// `c.config.VPCAuthEndpoint` using the "create_access_token" operation of the VPC Instance Metadata Service API.
func (c *client) GetInstanceIdentityToken() (*InstanceIdentityToken, error) {
	req := rest.PutRequest(c.config.instanceIdentityTokenEndpoint())
	// set headers
	req.Set(MetadataFlavor, defaultMetadataFlavor)
	// set query params
	req.Query("version", c.config.Version)

	// create body
	body := fmt.Sprintf("{\"expires_in\": %d}", defaultInstanceIdentityTokenLifetime)
	req.Body(body)

	var tokenResponse struct {
		InstanceIdentityToken
	}

	_, err := c.client.Do(req, &tokenResponse, nil)
	if err != nil {
		return nil, err
	}

	ret := tokenResponse.InstanceIdentityToken

	return &ret, nil
}

// GetIAMAccessToken exchanges the VPC/VSI compute resource's instance identity token for an IAM access token at
// `c.config.VPCAuthEndpoint` using the "create_iam_token" operation of the VPC Instance Metadata Service API.
// The request body can consist of a profile CRN or ID, but not both.
func (c *client) GetIAMAccessToken(tokenReq *IAMAccessTokenRequest) (*iam.Token, error) {
	req := rest.PostRequest(c.config.iamTokenEndpoint())
	// set instance identity token
	req.Set(authorization, "Bearer "+tokenReq.AccessToken)
	req.Query("version", c.config.Version)

	if tokenReq.Body != nil {
		var body string
		profileID := tokenReq.Body.ID
		profileCRN := tokenReq.Body.CRN

		if profileID != "" {
			body = fmt.Sprintf(`{"trusted_profile": {"id": "%s"}}`, profileID)
		} else if profileCRN != "" {
			body = fmt.Sprintf(`{"trusted_profile": {"crn": "%s"}}`, profileCRN)
		}

		if body != "" {
			req.Body(body)
		}

	}

	var tokenResponse struct {
		getIAMAccessTokenResponse
	}

	// error codes have not been defined, so return raw error from server
	_, err := c.client.Do(req, &tokenResponse, nil)
	if err != nil {
		return nil, err
	}

	// convert tokenResponse to an IAM token to maintain compatibility
	expiry, err := time.Parse(time.RFC3339, tokenResponse.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("Error parsing access token expiration %s", err)
	}
	iamToken := &iam.Token{
		AccessToken: tokenResponse.AccessToken,
		Expiry:      expiry,
		TokenType:   tokenType,
	}

	return iamToken, nil
}
