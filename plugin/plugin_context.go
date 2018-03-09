package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/authentication"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/consts"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
)

type pluginContext struct {
	core_config.ReadWriter
	cfConfig     cfConfigWrapper
	pluginConfig PluginConfig
	pluginPath   string
}

type cfConfigWrapper struct {
	core_config.CFConfig
}

func (c cfConfigWrapper) RefreshUAAToken() (string, error) {
	if !c.HasAPIEndpoint() {
		return "", fmt.Errorf("CloudFoundry API endpoint is not set")
	}

	config := &authentication.UAAConfig{UAAEndpoint: c.AuthenticationEndpoint()}
	auth := authentication.NewUAARepository(config, rest.NewClient())
	token, err := auth.RefreshToken(c.UAARefreshToken())
	if err != nil {
		return "", err
	}

	c.SetUAAToken(token.Token())
	c.SetUAARefreshToken(token.RefreshToken)
	return token.Token(), nil
}

func createPluginContext(pluginPath string, coreConfig core_config.ReadWriter) *pluginContext {
	return &pluginContext{
		pluginPath:   pluginPath,
		pluginConfig: loadPluginConfigFromPath(filepath.Join(pluginPath, "config.json")),
		ReadWriter:   coreConfig,
		cfConfig:     cfConfigWrapper{coreConfig.CFConfig()},
	}
}

func (c *pluginContext) PluginDirectory() string {
	return c.pluginPath
}

func (c *pluginContext) PluginConfig() PluginConfig {
	return c.pluginConfig
}

func (c *pluginContext) RefreshIAMToken() (string, error) {
	endpoint := os.Getenv("IAM_ENDPOINT")
	if endpoint == "" {
		endpoint = c.IAMEndpoint()
	}
	if endpoint == "" {
		return "", fmt.Errorf("IAM endpoint is not set")
	}

	config := &authentication.IAMConfig{TokenEndpoint: endpoint + "/identity/token"}
	auth := authentication.NewIAMAuthRepository(config, rest.NewClient())
	iamToken, err := auth.RefreshToken(c.IAMRefreshToken())
	if err != nil {
		return "", err
	}

	c.SetIAMToken(iamToken.Token())
	c.SetIAMRefreshToken(iamToken.RefreshToken)

	return iamToken.Token(), nil
}

func (c *pluginContext) Trace() string {
	return getFromEnvOrConfig(consts.ENV_BLUEMIX_TRACE, c.ReadWriter.Trace())
}

func (c *pluginContext) ColorEnabled() string {
	return getFromEnvOrConfig(consts.ENV_BLUEMIX_COLOR, c.ReadWriter.ColorEnabled())
}

func (c *pluginContext) VersionCheckEnabled() bool {
	return !c.CheckCLIVersionDisabled()
}

func (c *pluginContext) CF() CFContext {
	return c.cfConfig
}

func (c *pluginContext) HasTargetedCF() bool {
	return c.cfConfig.HasAPIEndpoint()
}

func getFromEnvOrConfig(envKey string, config string) string {
	if envVal := os.Getenv(envKey); envVal != "" {
		return envVal
	}
	return config
}

func (c *pluginContext) CommandNamespace() string {
	return os.Getenv(consts.ENV_BLUEMIX_PLUGIN_NAMESPACE)
}

func (c *pluginContext) CLIName() string {
	cliName := os.Getenv(consts.ENV_BLUEMIX_CLI)
	if cliName == "" {
		cliName = "bx"
	}
	return cliName
}
