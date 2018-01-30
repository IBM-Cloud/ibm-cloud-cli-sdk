// Package core_config provides functions to load core configuration.
// The package is for internal only.
package core_config

import (
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/config_helpers"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/models"
)

type Reader interface {
	// CF config
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

	OrganizationFields() models.OrganizationFields
	HasOrganization() bool
	SpaceFields() models.SpaceFields
	HasSpace() bool

	IsSSLDisabled() bool
	Locale() string
	Trace() string
	ColorEnabled() string

	// bx config
	ConsoleEndpoint() string
	Region() models.Region
	CloudName() string
	CloudType() string
	IAMEndpoint() string
	IAMID() string
	IAMToken() string
	IAMRefreshToken() string
	Account() models.Account
	HasAccount() bool
	IMSAccountID() string
	ResourceGroup() models.ResourceGroup
	HasResourceGroup() bool

	PluginRepos() []models.PluginRepo
	PluginRepo(string) (models.PluginRepo, bool)
	HTTPTimeout() int
	CLIInfoEndpoint() string
	CheckCLIVersionDisabled() bool
	UsageStatsDisabled() bool
}

type ReadWriter interface {
	Reader

	// cf config
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

	SetOrganizationFields(models.OrganizationFields)
	SetSpaceFields(models.SpaceFields)

	SetSSLDisabled(bool)
	SetLocale(string)
	SetTrace(string)
	SetColorEnabled(string)

	// bx config
	SetConsoleEndpoint(string)
	SetRegion(models.Region)
	SetAccount(models.Account)
	SetIAMEndpoint(string)
	SetIAMToken(string)
	SetIAMRefreshToken(string)
	SetResourceGroup(models.ResourceGroup)
	SetCheckCLIVersionDisabled(bool)
	SetCLIInfoEndpoint(string)
	SetPluginRepo(models.PluginRepo)
	UnSetPluginRepo(string)

	SetHTTPTimeout(int)
	SetUsageStatsDisabled(bool)

	ClearAPIInfo()
	ClearSession()
}

type configRepository struct {
	*cfConfigRepository
	*bxConfigRepository
}

func (c configRepository) ColorEnabled() string {
	return c.bxConfigRepository.ColorEnabled()
}

func (c configRepository) Trace() string {
	return c.bxConfigRepository.Trace()
}

func (c configRepository) Locale() string {
	return c.bxConfigRepository.Locale()
}

func (c configRepository) SetColorEnabled(enabled string) {
	c.bxConfigRepository.SetColorEnabled(enabled)
	c.cfConfigRepository.SetColorEnabled(enabled)
}

func (c configRepository) SetTrace(trace string) {
	c.bxConfigRepository.SetTrace(trace)
	c.cfConfigRepository.SetTrace(trace)
}

func (c configRepository) SetLocale(locale string) {
	c.bxConfigRepository.SetLocale(locale)
	c.cfConfigRepository.SetLocale(locale)
}

func (c configRepository) ClearSession() {
	c.cfConfigRepository.ClearSession()
	c.bxConfigRepository.ClearSession()
}

func (c configRepository) ClearAPIInfo() {
	c.cfConfigRepository.SetAPIEndpoint("")
	c.bxConfigRepository.ClearAPICache()
}

func NewCoreConfig(errHandler func(error)) ReadWriter {
	return NewCoreConfigFromPath(config_helpers.CFConfigFilePath(), config_helpers.ConfigFilePath(), errHandler)
}

func NewCoreConfigFromPath(cfConfigPath string, bxConfigPath string, errHandler func(error)) configRepository {
	return NewCoreConfigFromPersistor(configuration.NewDiskPersistor(cfConfigPath), configuration.NewDiskPersistor(bxConfigPath), errHandler)
}

func NewCoreConfigFromPersistor(cfPersistor configuration.Persistor, bxPersistor configuration.Persistor, errHandler func(error)) configRepository {
	return configRepository{
		cfConfigRepository: createCFConfigFromPersistor(cfPersistor, errHandler),
		bxConfigRepository: createBluemixConfigFromPersistor(bxPersistor, errHandler),
	}
}
