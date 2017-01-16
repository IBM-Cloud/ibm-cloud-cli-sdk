package core_config

// internal use only, plugin should use PluginContext

import (
	cfconfiguration "code.cloudfoundry.org/cli/cf/configuration"
	cfconfighelpers "code.cloudfoundry.org/cli/cf/configuration/confighelpers"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration/config_helpers"
)

type Reader interface {
	BXReader
	CFReader
	APICacheValid() bool
}

type ReadWriter interface {
	Reader

	BXWriter
	CFWriter
}

type configRepository struct {
	BXReadWriter
	CFReadWriter
}

func (c configRepository) APICacheValid() bool {
	return c.BXReadWriter.APICache().Target == c.CFReadWriter.APIEndpoint()
}

func NewCoreConfig(errHandler func(error)) ReadWriter {
	cfConfigPath, err := cfconfighelpers.DefaultFilePath()
	if err != nil {
		errHandler(err)
	}
	return NewCoreConfigFromPath(cfConfigPath, config_helpers.ConfigFilePath(), errHandler)
}

func NewCoreConfigFromPath(cfConfigPath string, bxConfigPath string, errHandler func(error)) configRepository {
	return NewCoreConfigFromPersistor(cfconfiguration.NewDiskPersistor(cfConfigPath), configuration.NewDiskPersistor(bxConfigPath), errHandler)
}

func NewCoreConfigFromPersistor(cfPersistor cfconfiguration.Persistor, bxPersistor configuration.Persistor, errHandler func(error)) configRepository {
	return configRepository{
		BXReadWriter: NewBluemixConfigFromPersistor(bxPersistor, errHandler),
		CFReadWriter: NewCFConfigAdapterFromPersistor(cfPersistor, errHandler),
	}
}
