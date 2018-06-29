package bluemix

import (
	"os"
)

var (
	EnvTrace        = newEnv("IBMCLOUD_TRACE", "BLUEMIX_TRACE")
	EnvColor        = newEnv("IBMCLOUD_COLOR", "BLUEMIX_COLOR")
	EnvVersionCheck = newEnv("IBMCLOUD_VERSION_CHECK", "BLUEMIX_VERSION_CHECK")
	EnvAnalytics    = newEnv("IBMCLOUD_ANALYTICS", "BLUEMIX_ANALYTICS")
	EnvHTTPTimeout  = newEnv("IBMCLOUD_HTTP_TIMEOUT", "BLUEMIX_HTTP_TIMEOUT")
	EnvAPIKey       = newEnv("IBMCLOUD_API_KEY", "BLUEMIX_API_KEY")
	EnvConfigHome   = newEnv("IBMCLOUD_HOME", "BLUEMIX_HOME")
	EnvConfigDir    = newEnv("IBMCLOUD_CONFIG_HOME")

	// for internal use
	EnvCLIName         = newEnv("IBMCLOUD_CLI", "BLUEMIX_CLI")
	EnvPluginNamespace = newEnv("IBMCLOUD_PLUGIN_NAMESPACE", "BLUEMIX_PLUGIN_NAMESPACE")
)

type Env struct {
	names []string
}

func (e Env) Get() string {
	for _, n := range e.names {
		if v := os.Getenv(n); v != "" {
			return v
		}
	}
	return ""
}

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
