package token_refresher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/testhelpers/configuration"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken(t *testing.T) {
	assert := assert.New(t)

	uaaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(r.ParseForm())
		assert.Equal("refresh_token", r.Form.Get("grant_type"))
		assert.Equal("old-refresh-token", r.Form.Get("refresh_token"))
		assert.Equal("", r.Form.Get("scope"))

		fmt.Fprint(w, `{"refresh_token": "new-refresh-token",
                        "access_token": "new-access-token",
                        "token_type": "bearer"}`)

	}))
	defer uaaServer.Close()

	config := configuration.NewFakeCoreConfig()
	config.SetUaaEndpoint(uaaServer.URL)
	config.SetRefreshToken("old-refresh-token")

	tokenRefresher := NewTokenRefresher(config)
	newToken, newRefreshToken, err := tokenRefresher.RefreshAuthToken()
	assert.NoError(err)
	assert.Equal("bearer new-access-token", newToken)
	assert.Equal("new-refresh-token", newRefreshToken)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	assert := assert.New(t)

	uaaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		fmt.Fprint(w, `{"error": "invalid-token",
                        "error_description": "the-error-description"}`)
	}))
	defer uaaServer.Close()

	config := configuration.NewFakeCoreConfig()
	config.SetUaaEndpoint(uaaServer.URL)
	config.SetRefreshToken("old-refresh-token")

	tokenRefresher := NewTokenRefresher(config)
	_, _, err := tokenRefresher.RefreshAuthToken()

	assert.Error(err)
	assert.IsType(new(InvalidTokenError), err)
}
