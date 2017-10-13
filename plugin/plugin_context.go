package plugin

import (
	"os"
	"path/filepath"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/authentication"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/consts"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/models"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
)

type pluginContext struct {
	coreConfig   core_config.ReadWriter
	pluginConfig PluginConfig
	pluginPath   string
}

func createPluginContext(pluginPath string, coreConfig core_config.ReadWriter) *pluginContext {
	return &pluginContext{
		pluginPath:   pluginPath,
		pluginConfig: loadPluginConfigFromPath(filepath.Join(pluginPath, "config.json")),
		coreConfig:   coreConfig,
	}
}

func (c *pluginContext) PluginDirectory() string {
	return c.pluginPath
}

func (c *pluginContext) PluginConfig() PluginConfig {
	return c.pluginConfig
}

func (c *pluginContext) APIVersion() string {
	return c.coreConfig.APIVersion()
}

func (c *pluginContext) APIEndpoint() string {
	return c.coreConfig.APIEndpoint()
}

func (c *pluginContext) HasAPIEndpoint() bool {
	return c.coreConfig.HasAPIEndpoint()
}

func (c *pluginContext) Region() models.Region {
	return c.coreConfig.Region()
}

func (c *pluginContext) CloudName() string {
	return c.coreConfig.CloudName()
}

func (c *pluginContext) CloudType() string {
	return c.coreConfig.CloudType()
}

func (c *pluginContext) DopplerEndpoint() string {
	return c.coreConfig.DopplerEndpoint()
}

func (c *pluginContext) ConsoleEndpoint() string {
	return c.coreConfig.ConsoleEndpoint()
}

// deprecate loggergator endpoint, use Doppler endpoint instead
//
// func (c *pluginContext) LoggregatorEndpoint() string {
// 	return c.coreConfig.LoggregatorEndpoint()
// }

func (c *pluginContext) UAAEndpoint() string {
	return c.coreConfig.UAAEndpoint()
}

func (c *pluginContext) UAAToken() string {
	return c.coreConfig.UAAToken()
}

func (c *pluginContext) UAARefreshToken() string {
	return c.coreConfig.UAARefreshToken()
}

func (c *pluginContext) RefreshUAAToken() (string, error) {
	auth := authentication.NewUAARepository(c.coreConfig, rest.NewClient())
	token, err := auth.RefreshToken(c.UAARefreshToken())
	if err != nil {
		return "", err
	}

	c.coreConfig.SetUAAToken(token.AccessToken)
	c.coreConfig.SetUAARefreshToken(token.RefreshToken)

	return token.AccessToken, nil
}

func (c *pluginContext) IAMTokenEndpoint() string {
	return c.coreConfig.IAMEndpoint()
}

func (c *pluginContext) IAMToken() string {
	return c.coreConfig.IAMToken()
}

func (c *pluginContext) IAMRefreshToken() string {
	return c.coreConfig.IAMRefreshToken()
}

func (c *pluginContext) RefreshIAMToken() (string, error) {
	auth := authentication.NewIAMAuthRepository(c.coreConfig, rest.NewClient())
	iamToken, uaaToken, err := auth.RefreshToken(c.IAMRefreshToken())
	if err != nil {
		return "", err
	}

	c.coreConfig.SetIAMToken(iamToken.AccessToken)
	c.coreConfig.SetIAMRefreshToken(iamToken.RefreshToken)
	c.coreConfig.SetUAAToken(uaaToken.AccessToken)
	c.coreConfig.SetUAARefreshToken(uaaToken.RefreshToken)

	return iamToken.AccessToken, nil
}

func (c *pluginContext) IsLoggedIn() bool {
	return c.coreConfig.IsLoggedIn()
}

func (c *pluginContext) UserEmail() string {
	return c.coreConfig.UserEmail()
}

func (c *pluginContext) UserGUID() string {
	return c.coreConfig.UserGUID()
}

func (c *pluginContext) Username() string {
	return c.coreConfig.Username()
}

func (c *pluginContext) AccountID() string {
	return c.Account().GUID
}

func (c *pluginContext) Account() models.Account {
	return c.coreConfig.Account()
}

func (c *pluginContext) IMSAccountID() string {
	return c.coreConfig.IMSAccountID()
}

func (c *pluginContext) ResourceGroup() models.ResourceGroup {
	return c.coreConfig.ResourceGroup()
}

func (c *pluginContext) CurrentOrg() models.OrganizationFields {
	return c.coreConfig.OrganizationFields()
}

func (c *pluginContext) HasOrganization() bool {
	return c.coreConfig.HasOrganization()
}

func (c *pluginContext) CurrentSpace() models.SpaceFields {
	return c.coreConfig.SpaceFields()
}

func (c *pluginContext) HasSpace() bool {
	return c.coreConfig.HasSpace()
}

//TODO: return locale based on both user configured locale and OS locale
func (c *pluginContext) Locale() string {
	return c.coreConfig.Locale()
}

func (c *pluginContext) IsSSLDisabled() bool {
	return c.coreConfig.IsSSLDisabled()
}

func (c *pluginContext) Trace() string {
	return getFromEnvOrConfig(consts.ENV_BLUEMIX_TRACE, c.coreConfig.Trace())
}

func (c *pluginContext) ColorEnabled() string {
	return getFromEnvOrConfig(consts.ENV_BLUEMIX_COLOR, c.coreConfig.ColorEnabled())
}

func (c *pluginContext) VersionCheckEnabled() bool {
	return !c.coreConfig.CheckCLIVersionDisabled()
}

func (c *pluginContext) HTTPTimeout() int {
	return c.coreConfig.HTTPTimeout()
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
