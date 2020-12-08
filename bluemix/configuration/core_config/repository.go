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
	IsPrivateEndpointEnabled() bool
	ConsoleEndpoints() models.Endpoints
	IAMEndpoints() models.Endpoints
	CloudName() string
	CloudType() string
	CurrentRegion() models.Region
	HasTargetedRegion() bool
	IAMToken() string
	IAMRefreshToken() string
	IsLoggedIn() bool
	IsLoggedInWithServiceID() bool
	UserEmail() string
	// UserDisplayText is the human readable ID for logged-in users which include non-human IDs
	UserDisplayText() string
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
	// UsageSatsDisabled returns whether the usage statistics data collection is disabled or not
	// Deprecated: use UsageSatsEnabled instead. We change to disable usage statistics by default,
	// So this property will not be used anymore
	UsageStatsDisabled() bool
	// UsageSatsEnabled returns whether the usage statistics data collection is enabled or not
	UsageStatsEnabled() bool
	// UsageStatsEnabledLastUpdate returns last time when `UsageStatsEnabled` was updated
	UsageStatsEnabledLastUpdate() time.Time
	Locale() string
	LoginAt() time.Time
	Trace() string
	ColorEnabled() string
	SDKVersion() string

	UnsetAPI()
	SetAPIEndpoint(string)
	SetPrivateEndpointEnabled(bool)
	SetConsoleEndpoints(models.Endpoints)
	SetIAMEndpoints(models.Endpoints)
	SetCloudType(string)
	SetCloudName(string)
	SetRegion(models.Region)
	SetIAMToken(string)
	SetIAMRefreshToken(string)
	ClearSession()
	SetAccount(models.Account)
	SetResourceGroup(models.ResourceGroup)
	SetLoginAt(loginAt time.Time)
	SetCheckCLIVersionDisabled(bool)
	SetCLIInfoEndpoint(string)
	SetPluginRepo(models.PluginRepo)
	UnsetPluginRepo(string)
	SetSSLDisabled(bool)
	SetHTTPTimeout(int)
	// SetUsageSatsDisabled disable or enable usage statistics data collection
	// Deprecated: use SetUsageSatsEnabled instead
	SetUsageStatsDisabled(bool)
	// SetUsageSatsEnabled enable or disable usage statistics data collection
	SetUsageStatsEnabled(bool)
	SetUpdateCheckInterval(time.Duration)
	SetUpdateRetryCheckInterval(time.Duration)
	SetUpdateNotificationInterval(time.Duration)
	SetLocale(string)
	SetTrace(string)
	SetColorEnabled(string)

	CFConfig() CFConfig
	HasTargetedCF() bool
	HasTargetedCFEE() bool
	HasTargetedPublicCF() bool
	SetCFEETargeted(bool)
	CFEEEnvID() string
	SetCFEEEnvID(string)
}

// Deprecated
type ReadWriter interface {
	Repository
}

type CFConfig interface {
	APIVersion() string
	APIEndpoint() string
	AsyncTimeout() uint
	ColorEnabled() string
	HasAPIEndpoint() bool
	AuthenticationEndpoint() string
	UAAEndpoint() string
	DopplerEndpoint() string
	RoutingAPIEndpoint() string
	SSHOAuthClient() string
	MinCFCLIVersion() string
	MinRecommendedCFCLIVersion() string
	Username() string
	UserGUID() string
	UserEmail() string
	Locale() string
	LoginAt() time.Time
	IsLoggedIn() bool
	SetLoginAt(loginAt time.Time)
	Trace() string
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
	cfConfig cfConfigWrapper
}

type cfConfigWrapper struct {
	*cfConfig
	bx *bxConfig
}

func (wrapper cfConfigWrapper) UnsetAPI() {
	wrapper.cfConfig.UnsetAPI()
	wrapper.bx.SetCFEEEnvID("")
	wrapper.bx.SetCFEETargeted(false)
}

func newRepository(bx *bxConfig, cf *cfConfig) repository {
	return repository{
		bxConfig: bx,
		cfConfig: cfConfigWrapper{cfConfig: cf, bx: bx},
	}
}

func (c repository) IsLoggedIn() bool {
	return c.bxConfig.IsLoggedIn() || c.cfConfig.IsLoggedIn()
}

func (c repository) IsLoggedInWithServiceID() bool {
	return c.bxConfig.IsLoggedIn() && NewIAMTokenInfo(c.IAMToken()).SubjectType == SubjectTypeServiceID
}

func (c repository) UserEmail() string {
	email := c.bxConfig.UserEmail()
	if email == "" {
		email = c.cfConfig.UserEmail()
	}
	return email
}

func (c repository) CFConfig() CFConfig {
	return c.cfConfig
}

func (c repository) HasTargetedCF() bool {
	return c.cfConfig.HasAPIEndpoint()
}

func (c repository) HasTargetedCFEE() bool {
	return c.HasTargetedCF() && c.bxConfig.HasTargetedCFEE()
}

func (c repository) HasTargetedPublicCF() bool {
	return c.HasTargetedCF() && !c.bxConfig.HasTargetedCFEE()
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
	c.bxConfig.SetCFEEEnvID("")
	c.cfConfig.UnsetAPI()
}

func (c repository) ClearSession() {
	c.bxConfig.ClearSession()
	c.cfConfig.ClearSession()
}

func NewCoreConfig(errHandler func(error)) ReadWriter {
	// config_helpers.MigrateFromOldConfig() // error ignored
	return NewCoreConfigFromPath(config_helpers.CFConfigFilePath(), config_helpers.ConfigFilePath(), errHandler)
}

func NewCoreConfigFromPath(cfConfigPath string, bxConfigPath string, errHandler func(error)) ReadWriter {
	return NewCoreConfigFromPersistor(configuration.NewDiskPersistor(cfConfigPath), configuration.NewDiskPersistor(bxConfigPath), errHandler)
}

func NewCoreConfigFromPersistor(cfPersistor configuration.Persistor, bxPersistor configuration.Persistor, errHandler func(error)) ReadWriter {
	return newRepository(createBluemixConfigFromPersistor(bxPersistor, errHandler), createCFConfigFromPersistor(cfPersistor, errHandler))
}
