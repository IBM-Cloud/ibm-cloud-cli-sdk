package config_helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/file_helpers"

	"github.com/stretchr/testify/assert"
)

func TestMigrateFromOldConfig(t *testing.T) {
	assert := assert.New(t)

	err := prepareBluemixHome()
	assert.NoError(err)
	defer clearBluemixHome()

	err = os.MkdirAll(oldConfigDir(), 0700)
	assert.NoError(err)
	oldConfigPath := filepath.Join(oldConfigDir(), "config.json")
	err = ioutil.WriteFile(oldConfigPath, []byte("old"), 0600)
	assert.NoError(err)

	err = MigrateFromOldConfig()
	assert.NoError(err)

	newConfigPath := filepath.Join(ConfigDir(), "config.json")
	assert.True(file_helpers.FileExists(newConfigPath))
	content, err := ioutil.ReadFile(newConfigPath)
	assert.NoError(err)
	assert.Equal([]byte("old"), content, "should copy old config file")

	assert.False(file_helpers.FileExists(oldConfigDir()), "old config dir should be deleted")
}

func TestMigrateFromOldConfig_NewConfigExist(t *testing.T) {
	assert := assert.New(t)

	err := prepareBluemixHome()
	assert.NoError(err)
	defer clearBluemixHome()

	err = os.MkdirAll(oldConfigDir(), 0700)
	assert.NoError(err)
	oldConfigPath := filepath.Join(oldConfigDir(), "config.json")
	err = ioutil.WriteFile(oldConfigPath, []byte("old"), 0600)
	assert.NoError(err)

	err = os.MkdirAll(ConfigDir(), 0700)
	assert.NoError(err)
	newConfigPath := filepath.Join(ConfigDir(), "config.json")
	err = ioutil.WriteFile(newConfigPath, []byte("new"), 0600)
	assert.NoError(err)

	err = MigrateFromOldConfig()
	assert.NoError(err)

	content, err := ioutil.ReadFile(newConfigPath)
	assert.NoError(err)
	assert.Equal([]byte("new"), content, "should not copy old config file")
}

func TestMigrateFromOldConfig_OldConfigNotExist(t *testing.T) {
	assert := assert.New(t)

	err := prepareBluemixHome()
	assert.NoError(err)
	defer clearBluemixHome()

	err = MigrateFromOldConfig()
	assert.NoError(err)
}

func prepareBluemixHome() error {
	temp, err := ioutil.TempDir("", "IBMCloudSDKConfigTest")
	if err != nil {
		return err
	}
	os.Setenv("BLUEMIX_HOME", temp)
	return nil
}

func clearBluemixHome() {
	if homeDir := os.Getenv("BLUEMIX_HOME"); homeDir != "" {
		os.RemoveAll(homeDir)
		os.Unsetenv("BLUEMIX_HOME")
	}
}
