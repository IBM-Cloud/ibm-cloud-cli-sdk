package bluemix

import (
	"os"
)

var (
	// EnvTrace is the environment variable `IBMCLOUD_TRACE` and `BLUEMIX_TRACE` (deprecated)
	EnvTrace = newEnv("IBMCLOUD_TRACE", "BLUEMIX_TRACE")
	// EnvColor is the environment variable `IBMCLOUD_COLOR` and `BLUEMIX_COLOR` (deprecated)
	EnvColor = newEnv("IBMCLOUD_COLOR", "BLUEMIX_COLOR")
	// EnvVersionCheck is the environment variable `IBMCLOUD_VERSION_CHECK` and `BLUEMIX_VERSION_CHECK` (deprecated)
	EnvVersionCheck = newEnv("IBMCLOUD_VERSION_CHECK", "BLUEMIX_VERSION_CHECK")
	// EnvAnalytics is the environment variable `IBMCLOUD_ANALYTICS` and `BLUEMIX_ANALYTICS` (deprecated)
	EnvAnalytics = newEnv("IBMCLOUD_ANALYTICS", "BLUEMIX_ANALYTICS")
	// EnvHTTPTimeout is the environment variable `IBMCLOUD_HTTP_TIMEOUT` and `BLUEMIX_HTTP_TIMEOUT` (deprecated)
	EnvHTTPTimeout = newEnv("IBMCLOUD_HTTP_TIMEOUT", "BLUEMIX_HTTP_TIMEOUT")
	// EnvAPIKey is the environment variable `IBMCLOUD_API_KEY` and `BLUEMIX_API_KEY` (deprecated)
	EnvAPIKey = newEnv("IBMCLOUD_API_KEY", "BLUEMIX_API_KEY")
	// EnvCRToken is the environment variable `IBMCLOUD_CR_TOKEN`
	EnvCRTokenKey = newEnv("IBMCLOUD_CR_TOKEN")
	// EnvCRProfile is the environment variable `IBMCLOUD_CR_PROFILE`
	EnvCRProfile = newEnv("IBMCLOUD_CR_PROFILE")
	// EnvCRVpcUrl is the environment variable `IBMCLOUD_CR_VPC_URL`
	EnvCRVpcUrl = newEnv("IBMCLOUD_CR_VPC_URL")
	// EnvConfigHome is the environment variable `IBMCLOUD_HOME` and `BLUEMIX_HOME` (deprecated)
	EnvConfigHome = newEnv("IBMCLOUD_HOME", "BLUEMIX_HOME")
	// EnvConfigDir is the environment variable `IBMCLOUD_CONFIG_HOME`
	EnvConfigDir = newEnv("IBMCLOUD_CONFIG_HOME")
	// EnvQuiet is the environment variable `IBMCLOUD_QUIET`
	EnvQuiet = newEnv("IBMCLOUD_QUIET")

	// for internal use
	EnvCLIName         = newEnv("IBMCLOUD_CLI", "BLUEMIX_CLI")
	EnvPluginNamespace = newEnv("IBMCLOUD_PLUGIN_NAMESPACE", "BLUEMIX_PLUGIN_NAMESPACE")
	EnvMCP             = newEnv("IBMCLOUD_MCP_ENABLED")
)

// Env is an environment variable supported by IBM Cloud CLI for specific purpose
// An Env could be bound to multiple environment variables due to historical reasons (i.e. renaming)
// Make sure you define the latest environment variable first
type Env struct {
	names []string
}

// Get will return the value of the environment variable, the first found non-empty value will be returned
func (e Env) Get() string {
	for _, n := range e.names {
		if v := os.Getenv(n); v != "" {
			return v
		}
	}
	return ""
}

// Set will set the value to **ALL** belonging environment variables
func (e Env) Set(val string) error {
	for _, n := range e.names {
		if err := os.Setenv(n, val); err != nil {
			return err
		}
	}
	return nil
}

func newEnv(names ...string) Env {
	return Env{names: names}
}
