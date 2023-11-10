package config_helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureAndPrepareEnv(a *assert.Assertions) ([]string, string) {
	env := os.Environ()

	userHome, err := ioutil.TempDir("", "config_dir_test")
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

func TestCompressEncodeDecompressDecode(t *testing.T) {
	assert := assert.New(t)

	data := "eyJraWQiOiIyMDE3MTAzMC0wMDowMDowMCIsImFsZyI6IlJTMjU2In0.eyJpYW1faWQiOiJJQk1pZC0yNzAwMDZWOEhNIiwiaWQiOiJJQk1pZC0yNzAwMDZWOEhNIiwicmVhbG1pZCI6IklCTWlkIiwiaWRlbnRpZmllciI6IjI3MDAwNlY4SE0iLCJnaXZlbl9uYW1lIjoiT0UgUnVudGltZXMiLCJmYW1pbHlfbmFtZSI6IlN5c3RlbSBVc2VyIiwibmFtZSI6Ik9FIFJ1bnRpbWVzIFN5c3RlbSBVc2VyIiwiZW1haWwiOiJydHN5c3VzckBjbi5pYm0uY29tIiwic3ViIjoicnRzeXN1c3JAY24uaWJtLmNvbSIsImFjY291bnQiOnsiYnNzIjoiOGQ2M2ZiMWNjNWU5OWU4NmRkNzIyOWRkZGZmYzA1YTUifSwiaWF0IjoxNTE2MTc0NjAzLCJleHAiOjE1MTYxNzgyMDMsImlzcyI6Imh0dHBzOi8vaWFtLmJsdWVtaXgubmV0L2lkZW50aXR5IiwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwic2NvcGUiOiJvcGVuaWQiLCJjbGllbnRfaWQiOiJieCJ9.gx-HQ1CSEwz5d4O1HXx4pusaYeEsqkQZgoBZ6esMBZG6wK6wQFPvC4D0Yvdi6CvKrVU-zV9PM_o3n5c-DFKjjTyTnRbQgrG0EPCRPmFW3bpepSb7eSw01S2YOLy5UTbz0cdM9hq-jafOu1S8pe9xeSMIMiA3-EFzCap5Z5CuoK9oIYJIFWseb1KsOyoiNOellbw1MaOmMzb4fsFz5Dr1Y8c1pNhoqp8M62E3y1yHe2jc6YepDab7Dqn2benK_e-MI3BlyWuBu4yo5mY2oCinJthr2E1YgbzWvcMy5a-ximnQIb4K6kscuUW_Yj_1GhDGJs4MP9u7M3-XdY1CNBGYeQp"

	encode, err := CompressEncode([]byte(data))
	assert.Nil(err)

	// expect the compressed encoded string to be smaller than the original data
	assert.True(len(encode) < len(data))

	decode, err := DecodeDecompress(encode)
	assert.Nil(err)
	assert.Equal(string(decode), data)
}
