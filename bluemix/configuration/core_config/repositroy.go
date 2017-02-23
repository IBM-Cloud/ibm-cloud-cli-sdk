package core_config

import (
	cfconfiguration "code.cloudfoundry.org/cli/cf/configuration"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/config_helpers"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/models"
)

type Reader interface {
	CFConfigReader

	Region() string
	IAMToken() string
	IAMRefreshToken() string
	Account() models.Account
	HasAccount() bool
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
