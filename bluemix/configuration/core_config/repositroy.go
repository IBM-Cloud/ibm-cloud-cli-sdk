// Package core_config provides functions to load core configuration.
// The package is for internal only.
package core_config

import (
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/config_helpers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
)

type Repository interface {
	APIEndpoint() string
	HasAPIEndpoint() bool
	ConsoleEndpoint() string
	IAMEndpoint() string
	CloudName() string
	CloudType() string
	CurrentRegion() models.Region
	IAMToken() string
	IAMRefreshToken() string
	IsLoggedIn() bool
	UserEmail() string
	IAMID() string
	CurrentAccount() models.Account
	HasTargetedAccount() bool
	IMSAccountID() string
	CurrentResourceGroup() models.ResourceGroup
	HasTargetedResourceGroup() bool
	PluginRepos() []models.PluginRepo
	PluginRepo(string) (models.PluginRepo, bool)
	IsSSLDisabled() bool
	HTTPTimeout() int
	CLIInfoEndpoint() string
	CheckCLIVersionDisabled() bool
	UpdateCheckInterval() time.Duration
	UpdateRetryCheckInterval() time.Duration
	UpdateNotificationInterval() time.Duration
	UsageStatsDisabled() bool
	Locale() string
	Trace() string
	ColorEnabled() string
	SDKVersion() string

	UnsetAPI()
	SetAPIEndpoint(string)
	SetConsoleEndpoint(string)
	SetIAMEndpoint(string)
	SetRegion(models.Region)
	SetIAMToken(string)
	SetIAMRefreshToken(string)
	ClearSession()
	SetAccount(models.Account)
	SetResourceGroup(models.ResourceGroup)
	SetCheckCLIVersionDisabled(bool)
	SetCLIInfoEndpoint(string)
	SetPluginRepo(models.PluginRepo)
	UnsetPluginRepo(string)
	SetSSLDisabled(bool)
	SetHTTPTimeout(int)
	SetUsageStatsDisabled(bool)
	SetUpdateCheckInterval(time.Duration)
	SetUpdateRetryCheckInterval(time.Duration)
	SetUpdateNotificationInterval(time.Duration)
	SetLocale(string)
	SetTrace(string)
	SetColorEnabled(string)

	CFConfig() CFConfig
	HasTargetedCF() bool
	HasTargetedCFEE() bool
	SetCFEETargeted(bool)
}

// Deprecated
type ReadWriter interface {
	Repository
}

type CFConfig interface {
	APIVersion() string
	APIEndpoint() string
	HasAPIEndpoint() bool
	AuthenticationEndpoint() string
	UAAEndpoint() string
	LoggregatorEndpoint() string
	DopplerEndpoint() string
	RoutingAPIEndpoint() string
	SSHOAuthClient() string
	MinCFCLIVersion() string
	MinRecommendedCFCLIVersion() string
	Username() string
	UserGUID() string
	UserEmail() string
	IsLoggedIn() bool
	UAAToken() string
	UAARefreshToken() string
	CurrentOrganization() models.OrganizationFields
	HasTargetedOrganization() bool
	CurrentSpace() models.SpaceFields
	HasTargetedSpace() bool

	UnsetAPI()
	SetAPIVersion(string)
	SetAPIEndpoint(string)
	SetAuthenticationEndpoint(string)
	SetLoggregatorEndpoint(string)
	SetDopplerEndpoint(string)
	SetUAAEndpoint(string)
	SetRoutingAPIEndpoint(string)
	SetSSHOAuthClient(string)
	SetMinCFCLIVersion(string)
	SetMinRecommendedCFCLIVersion(string)
	SetUAAToken(string)
	SetUAARefreshToken(string)
	SetOrganization(models.OrganizationFields)
	SetSpace(models.SpaceFields)
	ClearSession()
}

type repository struct {
	*bxConfig
	cfConfig *cfConfig
}

func (c repository) CFConfig() CFConfig {
	return c.cfConfig
}

func (c repository) HasTargetedCF() bool {
	return c.cfConfig.HasAPIEndpoint()
}

func (c repository) HasTargetedCFEE() bool {
	return c.HasTargetedCF() && c.CFEETargeted()
}

func (c repository) SetSSLDisabled(disabled bool) {
	c.bxConfig.SetSSLDisabled(disabled)
	c.cfConfig.SetSSLDisabled(disabled)
}

func (c repository) SetColorEnabled(enabled string) {
	c.bxConfig.SetColorEnabled(enabled)
	c.cfConfig.SetColorEnabled(enabled)
}

func (c repository) SetTrace(trace string) {
	c.bxConfig.SetTrace(trace)
	c.cfConfig.SetTrace(trace)
}

func (c repository) SetLocale(locale string) {
	c.bxConfig.SetLocale(locale)
	c.cfConfig.SetLocale(locale)
}

func (c repository) UnsetAPI() {
	c.bxConfig.UnsetAPI()
	c.bxConfig.SetCFEETargeted(false)
	c.cfConfig.UnsetAPI()
}

func (c repository) ClearSession() {
	c.bxConfig.ClearSession()
	c.cfConfig.ClearSession()
}

func NewCoreConfig(errHandler func(error)) ReadWriter {
	return NewCoreConfigFromPath(config_helpers.CFConfigFilePath(), config_helpers.ConfigFilePath(), errHandler)
}

func NewCoreConfigFromPath(cfConfigPath string, bxConfigPath string, errHandler func(error)) ReadWriter {
	return NewCoreConfigFromPersistor(configuration.NewDiskPersistor(cfConfigPath), configuration.NewDiskPersistor(bxConfigPath), errHandler)
}

func NewCoreConfigFromPersistor(cfPersistor configuration.Persistor, bxPersistor configuration.Persistor, errHandler func(error)) ReadWriter {
	return repository{
		cfConfig: createCFConfigFromPersistor(cfPersistor, errHandler),
		bxConfig: createBluemixConfigFromPersistor(bxPersistor, errHandler),
	}
}
