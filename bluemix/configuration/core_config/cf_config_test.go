package core_config_test

import (
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c := core_config.NewCFConfigData()

	assert.Equal(t, "cf", c.UAAOAuthClient)
	assert.Equal(t, "", c.UAAOAuthClientSecret)

}

var v4JSON = `{
  "AccessToken": "foo",
  "APIVersion": "5",
  "AsyncTimeout": 0,
  "AuthorizationEndpoint": "iam.test.cloud.ibm.com",
  "ColorEnabled": "",
  "ConfigVersion": 4,
  "DopplerEndpoint": "",
  "Locale": "",
  "LogCacheEndPoint": "",
  "MinCLIVersion": "",
  "MinRecommendedCLIVersion": "",
  "OrganizationFields": {
    "GUID": "",
    "Name": "",
    "QuotaDefinition": {
      "name": "",
      "memory_limit": 0,
      "instance_memory_limit": 0,
      "total_routes": 0,
      "total_services": 0,
      "non_basic_services_allowed": false,
      "app_instance_limit": 0
    }
  },
  "PluginRepos": null,
  "RefreshToken": "",
  "RoutingAPIEndpoint": "",
  "SpaceFields": {
    "GUID": "",
    "Name": "",
    "AllowSSH": false
  },
  "SSHOAuthClient": "",
  "SSLDisabled": false,
  "Target": "",
  "Trace": "",
  "UaaEndpoint": "",
  "LoginAt": "0001-01-01T00:00:00Z",
  "UAAGrantType": "",
  "UAAOAuthClient": "cf",
  "UAAOAuthClientSecret": ""
}`

func TestMarshal(t *testing.T) {
	c := core_config.NewCFConfigData()
	c.APIVersion = "5"
	c.AccessToken = "foo"
	c.AuthorizationEndpoint = "iam.test.cloud.ibm.com"

	bytes, err := c.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, string(bytes), v4JSON)
}

func TestUnmarshalV3(t *testing.T) {
	c := core_config.NewCFConfigData()
	assert.NoError(t, c.Unmarshal([]byte(v4JSON)))

	assert.Equal(t, 4, c.ConfigVersion)
	assert.Equal(t, "5", c.APIVersion)
	assert.Equal(t, "foo", c.AccessToken)
	assert.Equal(t, "iam.test.cloud.ibm.com", c.AuthorizationEndpoint)
}

func TestUnmarshalV4(t *testing.T) {
	var v3JSON = `{
  "AccessToken": "bar",
  "APIVersion": "4",
  "AsyncTimeout": 0,
  "AuthorizationEndpoint": "iam.test.cloud.ibm.com",
  "ColorEnabled": "",
  "ConfigVersion": 3,
  "DopplerEndpoint": "",
  "Locale": "",
  "LogCacheEndPoint": "",
  "MinCLIVersion": "",
  "MinRecommendedCLIVersion": "",
  "OrganizationFields": {
    "GUID": "",
    "Name": "",
    "QuotaDefinition": {
      "name": "",
      "memory_limit": 0,
      "instance_memory_limit": 0,
      "total_routes": 0,
      "total_services": 0,
      "non_basic_services_allowed": false,
      "app_instance_limit": 0
    }
  },
  "PluginRepos": null,
  "RefreshToken": "",
  "RoutingAPIEndpoint": "",
  "SpaceFields": {
    "GUID": "",
    "Name": "",
    "AllowSSH": false
  },
  "SSHOAuthClient": "",
  "SSLDisabled": false,
  "Target": "",
  "Trace": "",
  "UaaEndpoint": "",
  "LoginAt": "0001-01-01T00:00:00Z",
  "UAAGrantType": "",
  "UAAOAuthClient": "cf",
  "UAAOAuthClientSecret": ""
}`

	c := core_config.NewCFConfigData()
	assert.NoError(t, c.Unmarshal([]byte(v3JSON)))

	assert.Equal(t, 3, c.ConfigVersion)
	assert.Equal(t, "4", c.APIVersion)
	assert.Equal(t, "bar", c.AccessToken)
	assert.Equal(t, "iam.test.cloud.ibm.com", c.AuthorizationEndpoint)
}

func TestUnmarshalError(t *testing.T) {
	c := core_config.NewCFConfigData()
	assert.Error(t, c.Unmarshal([]byte(`{"db":cf}`)))
}
