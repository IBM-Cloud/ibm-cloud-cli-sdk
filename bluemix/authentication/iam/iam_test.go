package iam

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
)

const (
	crAuthMockIAMProfileName string = "iam-user-123"
	crAuthMockIAMProfileID   string = "iam-id-123"
	crAuthMockIAMProfileCRN  string = "crn:v1:bluemix:public:iam-identity::a/123456::profile:Profile-9fd84246-7df4-4667-94e4-8ecde51d5ac5"
	crAuthTestCRToken1       string = "cr-token-1"
	crAuthTestAccessToken1   string = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6ImhlbGxvIiwicm9sZSI6InVzZXIiLCJwZXJtaXNzaW9ucyI6WyJhZG1pbmlzdHJhdG9yIiwiZGVwbG95bWVudF9hZG1pbiJdLCJzdWIiOiJoZWxsbyIsImlzcyI6IkpvaG4iLCJhdWQiOiJEU1giLCJ1aWQiOiI5OTkiLCJpYXQiOjE1NjAyNzcwNTEsImV4cCI6MTU2MDI4MTgxOSwianRpIjoiMDRkMjBiMjUtZWUyZC00MDBmLTg2MjMtOGNkODA3MGI1NDY4In0.cIodB4I6CCcX8vfIImz7Cytux3GpWyObt9Gkur5g1QI"
	crAuthTestAccessToken2   string = "3yJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6ImhlbGxvIiwicm9sZSI6InVzZXIiLCJwZXJtaXNzaW9ucyI6WyJhZG1pbmlzdHJhdG9yIiwiZGVwbG95bWVudF9hZG1pbiJdLCJzdWIiOiJoZWxsbyIsImlzcyI6IkpvaG4iLCJhdWQiOiJEU1giLCJ1aWQiOiI5OTkiLCJpYXQiOjE1NjAyNzcwNTEsImV4cCI6MTU2MDI4MTgxOSwianRpIjoiMDRkMjBiMjUtZWUyZC00MDBmLTg2MjMtOGNkODA3MGI1NDY4In0.cIodB4I6CCcX8vfIImz7Cytux3GpWyObt9Gkur5g1QI"

	APPLICATION_JSON        = "application/json"
	FORM_URL_ENCODED_HEADER = "application/x-www-form-urlencoded"
)

func TestCRTokenRequestWithProfileID(t *testing.T) {
	// build a mock token request with ProfileID
	tokenReq := CRTokenRequest(crAuthTestCRToken1, crAuthMockIAMProfileID, "")

	// validate the request
	assert.Equal(t, tokenReq.GrantType(), GrantTypeCRToken)
	assert.Equal(t, tokenReq.GetTokenParam(profileIDParam), crAuthMockIAMProfileID)
	assert.Equal(t, tokenReq.GetTokenParam(profileNameParam), "")
	assert.Equal(t, tokenReq.GetTokenParam(profileCRNParam), "")
}

func TestCRTokenRequestWithProfileName(t *testing.T) {
	// build a mock token request with ProfileName
	tokenReq := CRTokenRequest(crAuthTestCRToken1, "", crAuthMockIAMProfileName)

	// validate the request
	assert.Equal(t, tokenReq.GrantType(), GrantTypeCRToken)
	assert.Equal(t, tokenReq.GetTokenParam(profileNameParam), crAuthMockIAMProfileName)
	assert.Equal(t, tokenReq.GetTokenParam(profileIDParam), "")
	assert.Equal(t, tokenReq.GetTokenParam(profileCRNParam), "")
}

func TestCRTokenRequestWithProfileCRN(t *testing.T) {
	// build a mock token request with ProfileCRN
	tokenReq := CRTokenRequestWithCRN(crAuthTestCRToken1, "", "", crAuthMockIAMProfileCRN)

	// validate the request
	assert.Equal(t, tokenReq.GrantType(), GrantTypeCRToken)
	assert.Equal(t, tokenReq.GetTokenParam(profileCRNParam), crAuthMockIAMProfileCRN)
	assert.Equal(t, tokenReq.GetTokenParam(profileIDParam), "")
	assert.Equal(t, tokenReq.GetTokenParam(profileNameParam), "")
}

func TestCRTokenRequestWithProfileNameAndID(t *testing.T) {
	// build a mock token request with both ProfileID and ProfileName
	tokenReq := CRTokenRequest(crAuthTestCRToken1, crAuthMockIAMProfileID, crAuthMockIAMProfileName)

	// validate the request
	assert.Equal(t, tokenReq.GrantType(), GrantTypeCRToken)
	assert.Equal(t, tokenReq.GetTokenParam(profileNameParam), crAuthMockIAMProfileName)
	assert.Equal(t, tokenReq.GetTokenParam(profileIDParam), crAuthMockIAMProfileID)
	assert.Equal(t, tokenReq.GetTokenParam(profileCRNParam), "")
}

func TestCRTokenRequestWithProfileNameAndIDandCRN(t *testing.T) {
	// build a mock token request with both ProfileID and ProfileName
	tokenReq := CRTokenRequestWithCRN(crAuthTestCRToken1, crAuthMockIAMProfileID, crAuthMockIAMProfileName, crAuthMockIAMProfileCRN)

	// validate the request
	assert.Equal(t, tokenReq.GrantType(), GrantTypeCRToken)
	assert.Equal(t, tokenReq.GetTokenParam(profileNameParam), crAuthMockIAMProfileName)
	assert.Equal(t, tokenReq.GetTokenParam(profileIDParam), crAuthMockIAMProfileID)
	assert.Equal(t, tokenReq.GetTokenParam(profileCRNParam), crAuthMockIAMProfileCRN)
}

func TestGetTokenOneFromServerSuccessWithProfileID(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 1, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequest(crAuthTestCRToken1, crAuthMockIAMProfileID, "")
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken1, IAMToken.AccessToken)
}

func TestGetTokenTwoFromServerSuccessWithProfileID(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 2, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequest(crAuthTestCRToken1, crAuthMockIAMProfileID, "")
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken2, IAMToken.AccessToken)
}

func TestGetTokenOneFromServerSuccessWithProfileName(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 1, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequest(crAuthTestCRToken1, "", crAuthMockIAMProfileName)
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken1, IAMToken.AccessToken)
}

func TestGetTokenTwoFromServerSuccessWithProfileName(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 2, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequest(crAuthTestCRToken1, "", crAuthMockIAMProfileName)
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken2, IAMToken.AccessToken)
}

func TestGetTokenOneFromServerSuccessWithProfileCRN(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 1, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequestWithCRN(crAuthTestCRToken1, "", "", crAuthMockIAMProfileCRN)
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken1, IAMToken.AccessToken)
}

func TestGetTokenTwoFromServerSuccessWithProfileCRN(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 2, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequestWithCRN(crAuthTestCRToken1, "", "", crAuthMockIAMProfileCRN)
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken2, IAMToken.AccessToken)
}

func TestGetTokenOneFromServerSuccessWithProfileNameAndID(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 1, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequest(crAuthTestCRToken1, crAuthMockIAMProfileID, crAuthMockIAMProfileName)
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken1, IAMToken.AccessToken)
}

func TestGetTokenTwoFromServerSuccessWithProfileNameAndID(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 2, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequest(crAuthTestCRToken1, crAuthMockIAMProfileID, crAuthMockIAMProfileName)
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken2, IAMToken.AccessToken)
}

func TestGetTokenOneFromServerSuccessWithProfileNameAndIDAndCRN(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 1, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequestWithCRN(crAuthTestCRToken1, crAuthMockIAMProfileID, crAuthMockIAMProfileName, crAuthMockIAMProfileCRN)
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken1, IAMToken.AccessToken)
}

func TestGetTokenTwoFromServerSuccessWithProfileNameAndIDAndCRN(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 2, http.StatusOK, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequestWithCRN(crAuthTestCRToken1, crAuthMockIAMProfileID, crAuthMockIAMProfileName, crAuthMockIAMProfileCRN)
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.Nil(t, err)
	assert.Equal(t, crAuthTestAccessToken2, IAMToken.AccessToken)
}

func TestGetTokenOneFromServerFailureWithProfileNameAndIDAndCRN(t *testing.T) {
	server := startMockIAMServerForCRExchange(t, 1, http.StatusUnauthorized, "")
	defer server.Close()

	mockIAMEndpoint := server.URL
	mockConfig := DefaultConfig(mockIAMEndpoint)
	mockClient := NewClient(mockConfig, rest.NewClient())

	// build the request, call fetch token, and verify response
	tokenReq := CRTokenRequestWithCRN(crAuthTestCRToken1, crAuthMockIAMProfileID, crAuthMockIAMProfileName, crAuthMockIAMProfileCRN)
	// Force the first fetch and verify we got the first access token.
	IAMToken, err := mockClient.GetToken(tokenReq)
	assert.NotNil(t, err)
	assert.Nil(t, IAMToken)
	assert.Contains(t, err.Error(), "Sorry, you are not authorized!")
}

func TestGetTokenOneFromServerApiErrorWithProfileNameAndID(t *testing.T) {
	errorCases := []struct {
		errorCode string
		errorMessage string

	}{
		{
			errorCode: InvalidTokenErrorCode,
			errorMessage: "invalid token",
		},
		{
			errorCode: RefreshTokenExpiryErrorCode,
			errorMessage: "refresh token expired",
		},
		{
			errorCode: ExternalAuthenticationErrorCode,
			errorMessage: "External authentication failed",
		},
		{
			errorCode: SessionInactiveErrorCode,
			errorMessage: "sdf",
		},
	}


	for _, errorCase := range errorCases {
		errorJson := fmt.Sprintf(`{"errorCode": "%s", "errorMessage": "%s", "errorDetails": "", "requirements": {"code": "", "error": ""}}`, errorCase.errorCode, errorCase.errorMessage)
		server := startMockIAMServerForCRExchange(t, 1, http.StatusUnauthorized, errorJson)
		defer server.Close()

		mockIAMEndpoint := server.URL
		mockConfig := DefaultConfig(mockIAMEndpoint)
		mockClient := NewClient(mockConfig, rest.NewClient())

		// build the request, call fetch token, and verify response
		tokenReq := CRTokenRequestWithCRN(crAuthTestCRToken1, crAuthMockIAMProfileID, crAuthMockIAMProfileName, crAuthMockIAMProfileCRN)
		// Force the first fetch and verify we got the first access token.
		IAMToken, err := mockClient.GetToken(tokenReq)
		assert.NotNil(t, err)
		assert.Nil(t, IAMToken)
		assert.Contains(t, err.Error(),errorCase.errorMessage)


	}



}

// startMockIAMServerForCRExchange will start a mock server endpoint that supports both the
// IAM operations that we'll need to call.
func startMockIAMServerForCRExchange(t *testing.T, call int, statusCode int, errorJson string) *httptest.Server {
	// Create the mock server.
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		operationPath := req.URL.EscapedPath()

		if operationPath == "/identity/token" {
			// If this is an invocation of the IAM "get_token" operation,
			// then validate it a bit and then send back a good response.
			assert.Equal(t, APPLICATION_JSON, req.Header.Get("Accept"))
			assert.Equal(t, FORM_URL_ENCODED_HEADER, req.Header.Get("Content-Type"))
			assert.Equal(t, crAuthTestCRToken1, req.FormValue("cr_token"))
			assert.Equal(t, GrantTypeCRToken.String(), req.FormValue("grant_type"))
			username, password, ok := req.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, defaultClientID, username)
			assert.Equal(t, defaultClientSecret, password)

			iamProfileID := req.FormValue("profile_id")
			iamProfileName := req.FormValue("profile_name")
			iamProfileCRN := req.FormValue("profile_crn")
			assert.True(t, iamProfileName != "" || iamProfileID != "" || iamProfileCRN != "")

			// This is the access token we'll send back in the mock response.
			// We'll default to token 2, then see if the caller asked for token 1
			// via the call param.
			var accessToken string
			if call == 1 {
				accessToken = crAuthTestAccessToken1
			} else {
				accessToken = crAuthTestAccessToken2
			}

			expiration := time.Now().Unix() + 3600
			res.WriteHeader(statusCode)
			mockErrorJson := errorJson
			switch statusCode {
			case http.StatusOK:
				fmt.Fprintf(res, `{"access_token": "%s", "token_type": "Bearer", "expires_in": 3600, "expiration": %d, "refresh_token": ""}`,
					accessToken, expiration)
			case http.StatusBadRequest:
				if errorJson == "" {
					mockErrorJson = "Sorry, bad request!"
				}
				fmt.Fprint(res,  mockErrorJson)

			case http.StatusUnauthorized:
				if errorJson == "" {
					mockErrorJson = "Sorry, you are not authorized!"
				}
				fmt.Fprint(res, mockErrorJson)
			}
		} else {
			assert.Fail(t, "unknown operation path: "+operationPath)
		}
	}))
	return server
}
