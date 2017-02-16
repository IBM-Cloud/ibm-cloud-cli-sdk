package core_config

import (
	cfconfiguration "code.cloudfoundry.org/cli/cf/configuration"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/config_helpers"
)

type Reader interface {
	CFConfigReader
	BXConfigReader
}

type ReadWriter interface {
	CFConfigReadWriter
	BXConfigReadWriter
}

type configRepository struct {
	CFConfigReadWriter
	BXConfigReadWriter
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
		BXConfigReadWriter: createBluemixConfigFromPersistor(bxPersistor, errHandler),
	}
}
