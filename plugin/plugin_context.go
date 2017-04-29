package plugin

import (
	"os"
	"path/filepath"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/authentication"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/config_helpers"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/consts"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/plugin/models"
)

type pluginContext struct {
	coreConfig   core_config.ReadWriter
	pluginConfig PluginConfig
	pluginPath   string
}

func NewPluginContext(pluginName string, coreConfig core_config.ReadWriter) *pluginContext {
	pluginPath := config_helpers.PluginDir(pluginName)
	c := &pluginContext{
		pluginPath:   pluginPath,
		pluginConfig: NewPluginConfig(filepath.Join(pluginPath, "config.json")),
	}
	c.coreConfig = coreConfig
	return c
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

func (c *pluginContext) Region() string {
	return c.coreConfig.Region()
}

func (c *pluginContext) AuthenticationEndpoint() string {
	return c.coreConfig.AuthenticationEndpoint()
}

func (c *pluginContext) DopplerEndpoint() string {
	return c.coreConfig.DopplerEndpoint()
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
	return authentication.NewUAARepository(c.coreConfig, new(rest.Client)).RefreshToken()
}

func (c *pluginContext) IAMToken() string {
	return c.coreConfig.IAMToken()
}

func (c *pluginContext) IAMRefreshToken() string {
	return c.coreConfig.IAMRefreshToken()
}

func (c *pluginContext) RefreshIAMToken() (string, error) {
	return authentication.NewIAMAuthRepository(c.coreConfig, new(rest.Client)).RefreshToken()
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
	token := c.coreConfig.IAMToken()
	if token != "" {
		return core_config.NewIAMTokenInfo(token).Accounts.AccountID
	}
	return c.coreConfig.Account().GUID
}

func (c *pluginContext) IMSAccountID() string {
	return c.coreConfig.IMSAccountID()
}

func (c *pluginContext) CurrentOrg() models.Organization {
	return models.Organization{
		OrganizationFields: c.coreConfig.OrganizationFields(),
	}
}

func (c *pluginContext) HasOrganization() bool {
	return c.coreConfig.HasOrganization()
}

func (c *pluginContext) CurrentSpace() models.Space {
	return models.Space{
		SpaceFields: c.coreConfig.SpaceFields(),
	}
}

func (c *pluginContext) HasSpace() bool {
	return c.coreConfig.HasSpace()
}

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
