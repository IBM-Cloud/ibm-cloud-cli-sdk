package vpc_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/vpc"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
)

const (
	instanceIdentityToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0aGVfYmVzdCI6IkVyaWNhIn0.c4C_BKtyZ4g78TB6wjdsX_MNx4KPoYj8YiikB1jO4o8"
	iamAccessToken        string = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6ImhlbGxvIiwicm9sZSI6InVzZXIiLCJwZXJtaXNzaW9ucyI6WyJhZG1pbmlzdHJhdG9yIiwiZGVwbG95bWVudF9hZG1pbiJdLCJzdWIiOiJoZWxsbyIsImlzcyI6IkpvaG4iLCJhdWQiOiJEU1giLCJ1aWQiOiI5OTkiLCJpYXQiOjE1NjAyNzcwNTEsImV4cCI6MTU2MDI4MTgxOSwianRpIjoiMDRkMjBiMjUtZWUyZC00MDBmLTg2MjMtOGNkODA3MGI1NDY4In0.cIodB4I6CCcX8vfIImz7Cytux3GpWyObt9Gkur5g1QI"

	applicationJson = "application/json"
	authorization   = "Authorization"
	contentType     = "Content-Type"

	mockProfileID  = "Profile-123"
	mockProfileCRN = "crn:iam-profile:12345"
)

var currentTime time.Time = time.Now()

func TestDefaultConfig(t *testing.T) {
	config := vpc.DefaultConfig("http://mockUrl", vpc.DefaultMetadataServiceVersion)
	assert.Equal(t, "http://mockUrl", config.VPCAuthEndpoint)
	assert.Equal(t, "http://mockUrl/instance_identity/v1/token", config.InstanceIdentiyTokenEndpoint)
	assert.Equal(t, "http://mockUrl/instance_identity/v1/iam_token", config.IAMTokenEndpoint)
	assert.Equal(t, vpc.DefaultMetadataServiceVersion, config.Version)

}
func TestIAMAccessTokenRequestWithProfileID(t *testing.T) {
	// build a mock token request with ProfileID
	req, err := vpc.NewIAMAccessTokenRequest(mockProfileID, "", instanceIdentityToken)
	assert.Nil(t, err)
	assert.NotNil(t, req)

	// validate the request
	assert.Equal(t, instanceIdentityToken, req.AccessToken)
	assert.NotNil(t, req.Body)
	assert.Empty(t, req.Body.CRN)
	assert.Equal(t, mockProfileID, req.Body.ID)
}

func TestIAMAccessTokenRequestWithProfileCRN(t *testing.T) {
	// build a mock token request with ProfileCRN
	req, err := vpc.NewIAMAccessTokenRequest("", mockProfileCRN, instanceIdentityToken)
	assert.Nil(t, err)
	assert.NotNil(t, req)

	// validate the request
	assert.Equal(t, instanceIdentityToken, req.AccessToken)
	assert.NotNil(t, req.Body)
	assert.Empty(t, req.Body.ID)
	assert.Equal(t, mockProfileCRN, req.Body.CRN)
}

func TestIAMAccessTokenRequestWithoutToken(t *testing.T) {
	// build a mock token request with ProfileCRN
	req, err := vpc.NewIAMAccessTokenRequest("", mockProfileCRN, "")
	assert.Nil(t, err)
	assert.NotNil(t, req)

	// validate the request
	assert.Empty(t, req.AccessToken)
	assert.NotNil(t, req.Body)
	assert.Empty(t, req.Body.ID)
	assert.Equal(t, mockProfileCRN, req.Body.CRN)
}

func TestIAMAccessTokenRequestFailure(t *testing.T) {
	req, err := vpc.NewIAMAccessTokenRequest(mockProfileID, mockProfileCRN, instanceIdentityToken)
	assert.NotNil(t, err)
	assert.Equal(t, "A Profile ID and Profile CRN cannot both be specified.", err.Error())
	assert.Nil(t, req)
}

func TestGetInstanceIdentiyTokenSuccess(t *testing.T) {
	server := startMockVPCServerForTokenExchange(t, http.StatusOK, false, false, true)
	defer server.Close()

	mockServerEndpoint := server.URL
	mockConfig := vpc.DefaultConfig(mockServerEndpoint, vpc.DefaultMetadataServiceVersion)
	mockClient := vpc.NewClient(mockConfig, rest.NewClient())

	// attempt to fetch the token
	identityToken, err := mockClient.GetInstanceIdentityToken()

	assert.Nil(t, err)
	assert.Equal(t, instanceIdentityToken, identityToken.AccessToken)
	assert.Equal(t, 300, identityToken.ExpiresIn)
	assert.NotEmpty(t, identityToken.ExpiresAt)

	// expiry should be 5 minutes in the future
	expectedExpiry := currentTime.Add(time.Second * 300).UTC().Format(time.RFC3339)
	assert.Equal(t, expectedExpiry, identityToken.ExpiresAt)

}

func TestGetInstanceIdentiyTokenFailure(t *testing.T) {
	server := startMockVPCServerForTokenExchange(t, http.StatusUnauthorized, false, false, true)
	defer server.Close()

	mockServerEndpoint := server.URL
	mockConfig := vpc.DefaultConfig(mockServerEndpoint, vpc.DefaultMetadataServiceVersion)
	mockClient := vpc.NewClient(mockConfig, rest.NewClient())

	// attempt to fetch the token
	identityToken, err := mockClient.GetInstanceIdentityToken()

	assert.NotNil(t, err)
	assert.Nil(t, identityToken)
}

func TestGetIAMAccessTokenSuccessWProfileID(t *testing.T) {
	server := startMockVPCServerForTokenExchange(t, http.StatusOK, false, true, true)
	defer server.Close()

	mockServerEndpoint := server.URL
	mockConfig := vpc.DefaultConfig(mockServerEndpoint, vpc.DefaultMetadataServiceVersion)
	mockClient := vpc.NewClient(mockConfig, rest.NewClient())

	// build the request using a mock profile ID
	req, err := vpc.NewIAMAccessTokenRequest(mockProfileID, "", instanceIdentityToken)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	// validate token
	assert.Equal(t, instanceIdentityToken, req.AccessToken)
	// validate the request body
	trustedProfile := req.Body
	assert.NotNil(t, trustedProfile)
	assert.NotNil(t, trustedProfile.ID)
	assert.Empty(t, trustedProfile.CRN)
	assert.Equal(t, mockProfileID, trustedProfile.ID)
	// attempt to fetch the token
	iamToken, err := mockClient.GetIAMAccessToken(req)

	assert.Nil(t, err)
	assert.Equal(t, iamAccessToken, iamToken.AccessToken)
	// expiry should be 1 hour in the future
	expectedExpiry := currentTime.Unix() + 3600
	assert.Equal(t, expectedExpiry, iamToken.Expiry.Unix())

}

func TestGetIAMAccessTokenSuccessWithoutBody1(t *testing.T) {
	server := startMockVPCServerForTokenExchange(t, http.StatusOK, false, false, false)
	defer server.Close()

	mockServerEndpoint := server.URL
	mockConfig := vpc.DefaultConfig(mockServerEndpoint, vpc.DefaultMetadataServiceVersion)
	mockClient := vpc.NewClient(mockConfig, rest.NewClient())

	// build the request manually without a body
	trustedProfile := vpc.TrustedProfileByIdOrCRN{}
	req := vpc.IAMAccessTokenRequest{
		AccessToken: instanceIdentityToken,
		Body:        &trustedProfile,
	}
	// attempt to fetch the token
	iamToken, err := mockClient.GetIAMAccessToken(&req)
	assert.Nil(t, err)
	assert.Equal(t, iamAccessToken, iamToken.AccessToken)
	// expiry should be 1 hour in the future
	expectedExpiry := currentTime.Unix() + 3600
	assert.Equal(t, expectedExpiry, iamToken.Expiry.Unix())

	// build the same request but using the `NewIAMAccessTokenRequest` method
	newReq, err := vpc.NewIAMAccessTokenRequest("", "", instanceIdentityToken)
	assert.Nil(t, err)
	assert.NotNil(t, newReq.Body)
	// attempt to fetch another token
	iamToken, err = mockClient.GetIAMAccessToken(newReq)
	assert.Nil(t, err)
	assert.Equal(t, iamAccessToken, iamToken.AccessToken)
}

func TestGetIAMAccessTokenSuccessWithoutBody2(t *testing.T) {
	server := startMockVPCServerForTokenExchange(t, http.StatusOK, false, false, false)
	defer server.Close()

	mockServerEndpoint := server.URL
	mockConfig := vpc.DefaultConfig(mockServerEndpoint, vpc.DefaultMetadataServiceVersion)
	mockClient := vpc.NewClient(mockConfig, rest.NewClient())

	// build the request manually without a body
	req := vpc.IAMAccessTokenRequest{
		AccessToken: instanceIdentityToken,
	}
	assert.Nil(t, req.Body)
	// attempt to fetch the token
	iamToken, err := mockClient.GetIAMAccessToken(&req)
	assert.Nil(t, err)
	assert.Equal(t, iamAccessToken, iamToken.AccessToken)
	// expiry should be 1 hour in the future
	expectedExpiry := currentTime.Unix() + 3600
	assert.Equal(t, expectedExpiry, iamToken.Expiry.Unix())

	// build the same request but using the `NewIAMAccessTokenRequest` method
	newReq, err := vpc.NewIAMAccessTokenRequest("", "", instanceIdentityToken)
	assert.Nil(t, err)
	assert.NotNil(t, newReq.Body)
	// attempt to fetch another token
	iamToken, err = mockClient.GetIAMAccessToken(newReq)
	assert.Nil(t, err)
	assert.Equal(t, iamAccessToken, iamToken.AccessToken)
}

func TestGetIAMAccessTokenFailureWProfileID(t *testing.T) {
	server := startMockVPCServerForTokenExchange(t, http.StatusBadRequest, false, true, true)
	defer server.Close()

	mockServerEndpoint := server.URL
	mockConfig := vpc.DefaultConfig(mockServerEndpoint, vpc.DefaultMetadataServiceVersion)
	mockClient := vpc.NewClient(mockConfig, rest.NewClient())

	// build the request using a mock profile ID
	req, err := vpc.NewIAMAccessTokenRequest(mockProfileID, "", "")
	assert.Nil(t, err)
	assert.NotNil(t, req)
	// validate token
	assert.Empty(t, req.AccessToken)
	// validate the request body
	trustedProfile := req.Body
	assert.NotNil(t, trustedProfile)
	assert.NotNil(t, trustedProfile.ID)
	assert.Empty(t, trustedProfile.CRN)
	assert.Equal(t, mockProfileID, trustedProfile.ID)
	// attempt to fetch the token
	iamToken, err := mockClient.GetIAMAccessToken(req)

	assert.NotNil(t, err)
	assert.Nil(t, iamToken)
}

func TestGetIAMAccessTokenSuccessWProfileCRN(t *testing.T) {
	server := startMockVPCServerForTokenExchange(t, http.StatusOK, true, false, true)
	defer server.Close()

	mockServerEndpoint := server.URL
	mockConfig := vpc.DefaultConfig(mockServerEndpoint, vpc.DefaultMetadataServiceVersion)
	mockClient := vpc.NewClient(mockConfig, rest.NewClient())

	// build the request using a mock profile CRN
	req, err := vpc.NewIAMAccessTokenRequest("", mockProfileCRN, instanceIdentityToken)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	// validate token
	assert.Equal(t, instanceIdentityToken, req.AccessToken)
	// validate the request body
	trustedProfile := req.Body
	assert.NotNil(t, trustedProfile)
	assert.NotNil(t, trustedProfile.CRN)
	assert.Empty(t, trustedProfile.ID)
	assert.Equal(t, mockProfileCRN, trustedProfile.CRN)
	// attempt to fetch the token
	iamToken, err := mockClient.GetIAMAccessToken(req)

	assert.Nil(t, err)
	assert.Equal(t, iamAccessToken, iamToken.AccessToken)
	assert.Equal(t, "Bearer", iamToken.TokenType)
	// expiry should be 1 hour in the future
	expectedExpiry := currentTime.Unix() + 3600
	assert.Equal(t, expectedExpiry, iamToken.Expiry.Unix())

}

func TestGetIAMAccessTokenFailureWProfileCRN(t *testing.T) {
	server := startMockVPCServerForTokenExchange(t, http.StatusBadRequest, true, false, true)
	defer server.Close()

	mockServerEndpoint := server.URL
	mockConfig := vpc.DefaultConfig(mockServerEndpoint, vpc.DefaultMetadataServiceVersion)
	mockClient := vpc.NewClient(mockConfig, rest.NewClient())

	// build the request using a mock profile CRN
	req, err := vpc.NewIAMAccessTokenRequest("", mockProfileCRN, "")
	assert.Nil(t, err)
	assert.NotNil(t, req)
	// validate token
	assert.Empty(t, req.AccessToken)
	// validate the request body
	trustedProfile := req.Body
	assert.NotNil(t, trustedProfile)
	assert.NotNil(t, trustedProfile.CRN)
	assert.Empty(t, trustedProfile.ID)
	assert.Equal(t, mockProfileCRN, trustedProfile.CRN)
	// attempt to fetch the token
	iamToken, err := mockClient.GetIAMAccessToken(req)

	assert.NotNil(t, err)
	assert.Nil(t, iamToken)

}

// startMockIAMServerForCRExchange will start a mock server endpoint that supports both the
// VPC operations for exchanging an instance identity token and an IAM TOKEN.
func startMockVPCServerForTokenExchange(t *testing.T, statusCode int, withCRN bool, withID bool, containsBody bool) *httptest.Server {
	// Create the mock server.
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		operationPath := req.URL.EscapedPath()

		if operationPath == "/instance_identity/v1/token" {
			// If this is an invocation of the VPC "create_access_token" operation,
			// then validate it and then send back a good response.
			assert.Equal(t, applicationJson, req.Header.Get("Accept"))
			assert.Equal(t, applicationJson, req.Header.Get(contentType))
			assert.Equal(t, "ibm", req.Header.Get(vpc.MetadataFlavor))
			// get query params and validate
			query := req.URL.Query()
			version := query.Get("version")

			assert.Equal(t, vpc.DefaultMetadataServiceVersion, version)
			// validate request body
			type createAccessTokenRequestBody struct {
				// Time in seconds before the access token expires.
				ExpiresIn *int64 `json:"expires_in,omitempty"`
			}
			// Unmarshal the request body.
			requestBody := &createAccessTokenRequestBody{}
			_ = json.NewDecoder(req.Body).Decode(requestBody)
			defer req.Body.Close()
			assert.NotNil(t, requestBody.ExpiresIn)

			expiration := currentTime.Add(time.Second * 300).UTC().Format(time.RFC3339)
			// convert expiration to valid Time.time value
			res.WriteHeader(statusCode)
			switch statusCode {
			case http.StatusOK:
				fmt.Fprintf(res, `{"access_token": "%s", "created_at": "%s", "expires_at": "%s", "expires_in": 300}`,
					instanceIdentityToken, time.Now().String(), expiration)
			case http.StatusBadRequest:
				fmt.Fprintf(res, `Sorry, bad request!`)

			case http.StatusUnauthorized:
				fmt.Fprintf(res, `Sorry, you are not authorized!`)
			}
		} else if operationPath == "/instance_identity/v1/iam_token" {
			// If this is an invocation of the VPC "create_iam_token" operation,
			// then validate it and then send back a good response.
			assert.Equal(t, applicationJson, req.Header.Get("Accept"))
			assert.Equal(t, applicationJson, req.Header.Get(contentType))
			authHeader := strings.TrimSpace(req.Header.Get(authorization))
			if authHeader != "Bearer" {
				assert.Equal(t, "Bearer "+instanceIdentityToken, req.Header.Get(authorization))
			}
			// get query params and validate
			query := req.URL.Query()
			version := query.Get("version")

			assert.Equal(t, vpc.DefaultMetadataServiceVersion, version)
			// Models the request body for the 'create_iam_token' operation.
			type createIamTokenRequestBody struct {
				TrustedProfile *vpc.TrustedProfileByIdOrCRN `json:"trusted_profile,omitempty"`
			}
			requestBody := &createIamTokenRequestBody{}
			_ = json.NewDecoder(req.Body).Decode(requestBody)
			defer req.Body.Close()

			if withID {
				assert.NotNil(t, requestBody.TrustedProfile.ID)
				assert.Equal(t, mockProfileID, requestBody.TrustedProfile.ID)
			} else if withCRN {
				assert.NotNil(t, requestBody.TrustedProfile.CRN)
				assert.Equal(t, mockProfileCRN, requestBody.TrustedProfile.CRN)
			} else if !containsBody {
				assert.NotNil(t, requestBody)
				assert.Nil(t, requestBody.TrustedProfile)
			}

			expiration := currentTime.Add(time.Second * 3600).UTC().Format(time.RFC3339)
			// convert expiration to valid Time.time value
			res.WriteHeader(statusCode)
			switch statusCode {
			case http.StatusOK:
				fmt.Fprintf(res, `{"access_token": "%s", "created_at": "%s", "expires_at": "%s", "expires_in": 3600}`,
					iamAccessToken, time.Now().String(), expiration)
			case http.StatusBadRequest:
				fmt.Fprintf(res, `Sorry, bad request!`)

			case http.StatusUnauthorized:
				fmt.Fprintf(res, `Sorry, you are not authorized!`)
			}
		} else {
			assert.Fail(t, "unknown operation path: "+operationPath)
		}
	}))
	return server
}
