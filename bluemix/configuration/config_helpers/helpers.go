// Package config_helpers provides helper functions to locate configuration files.
package config_helpers

import (
	"encoding/base64"
	gourl "net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/file_helpers"
)

func ConfigDir() string {
	if dir := bluemix.EnvConfigDir.Get(); dir != "" {
		return dir
	}
	// TODO: switched to the new default config after all plugin has bumped SDK
	if new := defaultConfigDirNew(); file_helpers.FileExists(new) {
		return new
	}
	return defaultConfigDirOld()
}

func defaultConfigDirNew() string {
	return filepath.Join(homeDir(), ".ibmcloud")
}

func defaultConfigDirOld() string {
	return filepath.Join(homeDir(), ".bluemix")
}

func homeDir() string {
	if homeDir := bluemix.EnvConfigHome.Get(); homeDir != "" {
		return homeDir
	}
	return UserHomeDir()
}

func TempDir() string {
	return filepath.Join(ConfigDir(), "tmp")
}

func ConfigFilePath() string {
	return filepath.Join(ConfigDir(), "config.json")
}

func PluginRepoDir() string {
	return filepath.Join(ConfigDir(), "plugins")
}

func PluginRepoCacheDir() string {
	return filepath.Join(PluginRepoDir(), ".cache")
}

func PluginsConfigFilePath() string {
	return filepath.Join(PluginRepoDir(), "config.json")
}

func PluginDir(pluginName string) string {
	return filepath.Join(PluginRepoDir(), pluginName)
}

func PluginBinaryLocation(pluginName string) string {
	executable := filepath.Join(PluginDir(pluginName), pluginName)
	if runtime.GOOS == "windows" {
		executable = executable + ".exe"
	}
	return executable
}

func CFHome() string {
	return ConfigDir()
}

func CFConfigDir() string {
	return filepath.Join(CFHome(), ".cf")
}

func CFConfigFilePath() string {
	return filepath.Join(CFConfigDir(), "config.json")
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}

// IsValidPaginationNextURL will return true if the provided nextURL has the expected queries provided
func IsValidPaginationNextURL(nextURL string, cursorQueryParamName string, expectedQueries gourl.Values) bool {
	parsedURL, parseErr := gourl.Parse(nextURL)
	// NOTE: ignore handling error(s) since if there error(s)
	// we can assume the url is invalid
	if parseErr != nil {
		return false
	}

	// retrive encoded cursor
	// eg. /api?cursor=<encode_string>
	queries := parsedURL.Query()
	encodedQuery := queries.Get(cursorQueryParamName)
	if encodedQuery == "" {
		return false

	}
	// decode string and parse encoded queries
	decodedQuery, decodedErr := base64.RawURLEncoding.DecodeString(encodedQuery)
	if decodedErr != nil {
		return false
	}
	queries, parsedErr := gourl.ParseQuery(string(decodedQuery))
	if parsedErr != nil {
		return false
	}

	// compare expected queries that should match
	// NOTE: assume queries are single value queries.
	// if multi-value queries will check the first query
	for expectedQuery := range expectedQueries {
		paginationQueryValue := queries.Get(expectedQuery)
		expectedQueryValue := expectedQueries.Get(expectedQuery)
		if paginationQueryValue != expectedQueryValue {
			return false
		}

	}

	return true
}
