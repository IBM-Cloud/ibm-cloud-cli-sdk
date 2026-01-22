package config_helpers

import (
	"encoding/base64"
	gourl "net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureAndPrepareEnv(a *assert.Assertions) ([]string, string) {
	env := os.Environ()

	userHome, err := os.MkdirTemp("", "config_dir_test")
	a.NoError(err)

	os.Unsetenv("IBMCLOUD_CONFIG_HOME")
	os.Unsetenv("IBMCLOUD_HOME")
	os.Unsetenv("BLUEMIX_HOME")
	os.Setenv("HOME", userHome)
	// UserHomeDir() uses HOMEDRIVE + HOMEPATH for windows
	if os.Getenv("OS") == "Windows_NT" {
		// ioutil.TempDir has the drive letter in the path, so we need to remove it when we set HOMEDRIVE
		os.Setenv("HOMEPATH", strings.Replace(userHome, os.Getenv("HOMEDRIVE"), "", -1))
	}
	a.NoError(os.RemoveAll(userHome))

	return env, userHome
}

func resetEnv(env []string) {
	os.Clearenv()
	for _, e := range env {
		pair := strings.Split(e, "=")
		os.Setenv(pair[0], pair[1])
	}
}

// If $USER_HOME/.ibmcloud does not exist, $USER_HOME/.bluemix should be used
func TestConfigDir_NothingSet_NothingExists(t *testing.T) {
	assert := assert.New(t)

	env, userHome := captureAndPrepareEnv(assert)
	defer resetEnv(env)
	defer os.RemoveAll(userHome)

	// directory should not exist - will use bluemix
	assert.Equal(filepath.Join(userHome, ".bluemix"), ConfigDir())
}

// If $USER_HOME/.ibmcloud exists, it should be used
func TestConfigDir_NothingSet_IBMCloudExists(t *testing.T) {
	assert := assert.New(t)

	env, userHome := captureAndPrepareEnv(assert)
	defer resetEnv(env)
	defer os.RemoveAll(userHome)

	// create a .ibmcloud directory and it should be returned
	ibmcloudDir := filepath.Join(userHome, ".ibmcloud")
	assert.NoError(os.MkdirAll(ibmcloudDir, 0700))
	assert.Equal(ibmcloudDir, ConfigDir())
}

// If only BLUEMIX_HOME is set, $BLUEMIX_HOME/.bluemix should be used
func TestConfigDir_BluemixHomeSet_NothingExists(t *testing.T) {
	assert := assert.New(t)

	env, userHome := captureAndPrepareEnv(assert)
	defer resetEnv(env)
	defer os.RemoveAll(userHome)

	// if only BLUEMIX_HOME is set, BLUEMIX_HOME is used
	os.Setenv("BLUEMIX_HOME", "/my_bluemix_home")
	assert.Equal(filepath.Join("/my_bluemix_home", ".bluemix"), ConfigDir())
}

// If BLUEMIX_HOME and IBMCLOUD_HOME are set and $IBMCLOUD_HOME/.ibmcloud does not exist, $IBMCLOUD_HOME/.bluemix should be used
func TestConfigDir_BluemixHomesAndIbmCloudHomeSet_NothingExists(t *testing.T) {
	assert := assert.New(t)

	env, userHome := captureAndPrepareEnv(assert)
	defer resetEnv(env)
	defer os.RemoveAll(userHome)

	// if BLUEMIX_HOME and IBMCLOUD_HOME is set, IBMCLOUD_HOME is used
	os.Setenv("BLUEMIX_HOME", "/my_bluemix_home")
	os.Setenv("IBMCLOUD_HOME", "/my_ibmcloud_home")
	assert.Equal(filepath.Join("/my_ibmcloud_home", ".bluemix"), ConfigDir())
}

// If IBMCLOUD_CONFIG_HOME is set, $IBMCLOUD_CONFIG_HOME should be used
func TestConfigDir_IbmCloudConfigHomeSet_Exists(t *testing.T) {
	assert := assert.New(t)

	env, userHome := captureAndPrepareEnv(assert)
	defer resetEnv(env)
	defer os.RemoveAll(userHome)

	// if IBMCLOUD_CONFIG_HOME is set and exists, IBMCLOUD_CONFIG_HOME is used
	os.Setenv("IBMCLOUD_CONFIG_HOME", userHome)
	assert.Equal(userHome, ConfigDir())
}

func TestIsValidPaginationNextURL(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name              string
		nextURL           string
		encodedQueryParam string
		expectedQueries   gourl.Values
		isValid           bool
	}{
		{
			name:              "return true for matching expected queries in pagination url",
			nextURL:           "/api/example?cursor=" + base64.RawURLEncoding.EncodeToString([]byte("limit=100&active=true")),
			encodedQueryParam: "cursor",
			expectedQueries: gourl.Values{
				"limit":  []string{"100"},
				"active": []string{"true"},
			},
			isValid: true,
		},
		{
			name:              "return true for matching expected queries with extraneous queries in pagination url",
			nextURL:           "/api/example?cursor=" + base64.RawURLEncoding.EncodeToString([]byte("limit=100&active=true&extra=foo")),
			encodedQueryParam: "cursor",
			expectedQueries: gourl.Values{
				"limit":  []string{"100"},
				"active": []string{"true"},
			},
			isValid: true,
		},
		{
			name:              "return false for different limit in pagination url",
			nextURL:           "/api/example?cursor=" + base64.RawURLEncoding.EncodeToString([]byte("limit=200")),
			encodedQueryParam: "cursor",
			expectedQueries: gourl.Values{
				"limit": []string{"100"},
			},
			isValid: false,
		},
		{
			name:              "return false for different query among multiple parameters in the pagination url",
			nextURL:           "/api/example?cursor=" + base64.RawURLEncoding.EncodeToString([]byte("limit=100&active=true")),
			encodedQueryParam: "cursor",
			expectedQueries: gourl.Values{
				"limit":  []string{"100"},
				"active": []string{"false"},
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(_ *testing.T) {
			isValid := IsValidPaginationNextURL(tc.nextURL, tc.encodedQueryParam, tc.expectedQueries)
			assert.Equal(tc.isValid, isValid)
		})
	}
}
