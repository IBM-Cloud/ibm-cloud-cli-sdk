package authentication_test

import (
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/iam"
	"github.com/stretchr/testify/assert"
)

func TestGetTokenParam(t *testing.T) {
	req := authentication.NewTokenRequest(iam.GrantTypeCRToken)
	profileParam := "myProfile"
	req.SetTokenParam("profile", profileParam)

	parsedParam := req.GetTokenParam("profile")

	assert.Equal(t, profileParam, parsedParam)
}
