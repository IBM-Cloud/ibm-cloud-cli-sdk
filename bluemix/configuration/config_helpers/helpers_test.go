package config_helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// import (
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.ibm.com/bluemix-cli-release/build/src/github.ibm.com/Bluemix/bluemix-cli-common/file_helpers"
// )

func captureEnv() []string {
	return os.Environ()
}

func resetEnv(env []string) {
	os.Clearenv()
	for _, e := range env {
		pair := strings.Split(e, "=")
		os.Setenv(pair[0], pair[1])
	}
}

func TestConfigDir(t *testing.T) {
	assert := assert.New(t)

	env := captureEnv()
	defer resetEnv(env)

	userHome, err := ioutil.TempDir("", "config_dir_test")
	assert.NoError(err)

	os.Unsetenv("IBMCLOUD_CONFIG_HOME")
	os.Unsetenv("IBMCLOUD_HOME")
	os.Unsetenv("BLUEMIX_HOME")
	os.Setenv("HOME", userHome)

	assert.NoError(os.RemoveAll(userHome))

	// directory should not exist - will use bluemix
	assert.Equal(filepath.Join(userHome, ".bluemix"), ConfigDir())

	// create a .ibmcloud directory and it should be returned
	ibmcloudDir := filepath.Join(userHome, ".ibmcloud")
	assert.NoError(os.MkdirAll(ibmcloudDir, 0755))
	assert.Equal(ibmcloudDir, ConfigDir())

	// if only BLUEMIX_HOME is set, BLUEMIX_HOME is used
	os.Setenv("BLUEMIX_HOME", "/my_bluemix_home")
	assert.Equal(filepath.Join("/my_bluemix_home", ".bluemix"), ConfigDir())

	// if BLUEMIX_HOME and IBMCLOUD_HOME is set, IBMCLOUD_HOME is used
	os.Setenv("IBMCLOUD_HOME", "/my_ibmcloud_home")
	assert.Equal(filepath.Join("/my_ibmcloud_home", ".bluemix"), ConfigDir())

	// if IBMCLOUD_CONFIG_HOME is set but does not exist, IBMCLOUD_HOME is used
	os.Setenv("IBMCLOUD_CONFIG_HOME", "/my_ibmcloud_config_home")
	assert.Equal(filepath.Join("/my_ibmcloud_home", ".bluemix"), ConfigDir())

	// if IBMCLOUD_CONFIG_HOME is set and exists, IBMCLOUD_CONFIG_HOME is used
	os.Setenv("IBMCLOUD_CONFIG_HOME", userHome)
	assert.Equal(userHome, ConfigDir())
}

// func TestMigrateFromOldConfig(t *testing.T) {
// 	assert := assert.New(t)

// 	err := prepareBluemixHome()
// 	assert.NoError(err)
// 	defer clearBluemixHome()

// 	err = os.MkdirAll(oldConfigDir(), 0700)
// 	assert.NoError(err)
// 	oldConfigPath := filepath.Join(oldConfigDir(), "config.json")
// 	err = ioutil.WriteFile(oldConfigPath, []byte("old"), 0600)
// 	assert.NoError(err)

// 	err = MigrateFromOldConfig()
// 	assert.NoError(err)

// 	newConfigPath := filepath.Join(newConfigDir(), "config.json")
// 	assert.True(file_helpers.FileExists(newConfigPath))
// 	content, err := ioutil.ReadFile(newConfigPath)
// 	assert.NoError(err)
// 	assert.Equal([]byte("old"), content, "should copy old config file")

// 	assert.False(file_helpers.FileExists(oldConfigDir()), "old config dir should be deleted")
// }

// func TestMigrateFromOldConfig_NewConfigExist(t *testing.T) {
// 	assert := assert.New(t)

// 	err := prepareBluemixHome()
// 	assert.NoError(err)
// 	defer clearBluemixHome()

// 	err = os.MkdirAll(oldConfigDir(), 0700)
// 	assert.NoError(err)
// 	oldConfigPath := filepath.Join(oldConfigDir(), "config.json")
// 	err = ioutil.WriteFile(oldConfigPath, []byte("old"), 0600)
// 	assert.NoError(err)

// 	err = os.MkdirAll(newConfigDir(), 0700)
// 	assert.NoError(err)
// 	newConfigPath := filepath.Join(newConfigDir(), "config.json")
// 	err = ioutil.WriteFile(newConfigPath, []byte("new"), 0600)
// 	assert.NoError(err)

// 	err = MigrateFromOldConfig()
// 	assert.NoError(err)

// 	content, err := ioutil.ReadFile(newConfigPath)
// 	assert.NoError(err)
// 	assert.Equal([]byte("new"), content, "should not copy old config file")
// }

// func TestMigrateFromOldConfig_OldConfigNotExist(t *testing.T) {
// 	assert := assert.New(t)

// 	err := prepareBluemixHome()
// 	assert.NoError(err)
// 	defer clearBluemixHome()

// 	err = MigrateFromOldConfig()
// 	assert.NoError(err)
// }

// func prepareBluemixHome() error {
// 	temp, err := ioutil.TempDir("", "IBMCloudSDKConfigTest")
// 	if err != nil {
// 		return err
// 	}
// 	os.Setenv("BLUEMIX_HOME", temp)
// 	return nil
// }

// func clearBluemixHome() {
// 	if homeDir := os.Getenv("BLUEMIX_HOME"); homeDir != "" {
// 		os.RemoveAll(homeDir)
// 		os.Unsetenv("BLUEMIX_HOME")
// 	}
// }
