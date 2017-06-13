package core_config

import (
	cfconfiguration "code.cloudfoundry.org/cli/cf/configuration"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/config_helpers"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/models"
)

type Reader interface {
	CFConfigReader

	ConsoleEndpoint() string
	Region() string
	IAMToken() string
	IAMRefreshToken() string
	Account() models.Account
	HasAccount() bool
	IMSAccountID() string
	PluginRepos() []models.PluginRepo
	PluginRepo(string) (models.PluginRepo, bool)
	HTTPTimeout() int
	CLIInfoEndpoint() string
	CheckCLIVersionDisabled() bool
	UsageStatsDisabled() bool
}

type ReadWriter interface {
	Reader

	CFConfigWriter

	SetConsoleEndpoint(string)
	SetRegion(string)
	SetAccount(models.Account)
	SetIAMToken(string)
	SetIAMRefreshToken(string)
	SetCheckCLIVersionDisabled(bool)
	SetCLIInfoEndpoint(string)
	SetPluginRepo(models.PluginRepo)
	UnSetPluginRepo(string)

	SetHTTPTimeout(int)
	SetUsageStatsDisabled(bool)

	ClearAPIInfo()
	ClearSession()
}

type SessionClearer interface {
	ClearSession()
}

type configRepository struct {
	CFConfigReadWriter
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
	c.CFConfigReadWriter.SetColorEnabled(enabled)
}

func (c configRepository) SetTrace(trace string) {
	c.bxConfigRepository.SetTrace(trace)
	c.CFConfigReadWriter.SetTrace(trace)
}

func (c configRepository) SetLocale(locale string) {
	c.bxConfigRepository.SetLocale(locale)
	c.CFConfigReadWriter.SetLocale(locale)
}

func (c configRepository) ClearSession() {
	c.CFConfigReadWriter.(SessionClearer).ClearSession()
	c.bxConfigRepository.ClearSession()
}

func (c configRepository) ClearAPIInfo() {
	c.CFConfigReadWriter.SetAPIEndpoint("")
	c.bxConfigRepository.ClearAPICache()
}

func NewCoreConfig(errHandler func(error)) ReadWriter {
	return NewCoreConfigFromPath(config_helpers.CFConfigFilePath(), config_helpers.ConfigFilePath(), errHandler)
}

func NewCoreConfigFromPath(cfConfigPath string, bxConfigPath string, errHandler func(error)) configRepository {
	return NewCoreConfigFromPersistor(cfconfiguration.NewDiskPersistor(cfConfigPath), configuration.NewDiskPersistor(bxConfigPath), errHandler)
}

func NewCoreConfigFromPersistor(cfPersistor cfconfiguration.Persistor, bxPersistor configuration.Persistor, errHandler func(error)) configRepository {
	return configRepository{
		CFConfigReadWriter: createCFConfigAdapterFromPersistor(cfPersistor, errHandler),
		bxConfigRepository: createBluemixConfigFromPersistor(bxPersistor, errHandler),
	}
}
