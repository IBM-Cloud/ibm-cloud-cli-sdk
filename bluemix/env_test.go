package bluemix_test

import (
	"os"
	"testing"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assert.Empty(t, EnvTrace.Get())

	os.Setenv("IBMCLOUD_TRACE", "true")
	assert.Equal(t, "true", EnvTrace.Get())

	os.Unsetenv("IBMCLOUD_TRACE")
	assert.Empty(t, EnvTrace.Get())

	os.Setenv("BLUEMIX_TRACE", "false")
	assert.Equal(t, "false", EnvTrace.Get())

	os.Unsetenv("BLUEMIX_TRACE")
	assert.Empty(t, EnvTrace.Get())
}

func TestSet(t *testing.T) {
	assert.Empty(t, os.Getenv("IBMCLOUD_COLOR"))
	assert.Empty(t, os.Getenv("BLUEMIX_COLOR"))
	assert.Empty(t, os.Getenv("IBMCLOUD_CR_TOKEN"))
	assert.Empty(t, os.Getenv("IBMCLOUD_CR_PROFILE"))
	assert.Empty(t, EnvColor.Get())

	EnvColor.Set("true")
	assert.Equal(t, "true", os.Getenv("IBMCLOUD_COLOR"))
	assert.Equal(t, "true", os.Getenv("BLUEMIX_COLOR"))
	assert.Equal(t, "true", EnvColor.Get())

	EnvCRTokenKey.Set("my_token")
	EnvCRProfile.Set("my_profile")

	assert.Equal(t, "my_token", os.Getenv("IBMCLOUD_CR_TOKEN"))
	assert.Equal(t, "my_profile", os.Getenv("IBMCLOUD_CR_PROFILE"))

	os.Unsetenv("IBMCLOUD_CR_TOKEN")
	os.Unsetenv("IBMCLOUD_CR_PROFILE")

	assert.Empty(t, os.Getenv("IBMCLOUD_CR_TOKEN"))
	assert.Empty(t, os.Getenv("IBMCLOUD_CR_PROFILE"))
}
