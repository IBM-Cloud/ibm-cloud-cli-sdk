package core_config

import (
	"os"
	"strings"
	"sync"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/models"
)

const (
	_DEFAULT_CLI_INFO_ENDPOINT = "https://clis.ng.bluemix.net/info"

	_DEFAULT_PLUGIN_REPO_NAME = "Bluemix"
	_DEFAULT_PLUGIN_REPO_URL  = "https://plugins.ng.bluemix.net"
)

type BXConfigReader interface {
	Region() string
	PluginRepos() []models.PluginRepo
	PluginRepo(string) (models.PluginRepo, bool)
	HTTPTimeout() int
	CLIInfoEndpoint() string
	CheckCLIVersionDisabled() bool
	UsageStatsDisabled() bool
}

type BXConfigReadWriter interface {
	BXConfigReader

	SetRegion(string)
	SetCheckCLIVersionDisabled(bool)
	SetCLIInfoEndpoint(string)
	SetPluginRepo(models.PluginRepo)
	UnSetPluginRepo(string)

	SetHTTPTimeout(int)
	SetUsageStatsDisabled(bool)
}

type BXConfigData struct {
	Region                  string
	PluginRepos             []models.PluginRepo
	Locale                  string
	Trace                   string
	ColorEnabled            string
	HTTPTimeout             int
	CLIInfoEndpoint         string
	CheckCLIVersionDisabled bool
	UsageStatsDisabled      bool
}

type bxConfigRepository struct {
	data      *BXConfigData
	persistor configuration.Persistor
	initOnce  *sync.Once
	lock      sync.RWMutex
	onError   func(error)
}

func createBluemixConfigFromPath(configPath string, errHandler func(error)) *bxConfigRepository {
	return createBluemixConfigFromPersistor(configuration.NewDiskPersistor(configPath), errHandler)
}

func createBluemixConfigFromPersistor(persistor configuration.Persistor, errHandler func(error)) *bxConfigRepository {
	return &bxConfigRepository{
		data:      new(BXConfigData),
		persistor: persistor,
		initOnce:  new(sync.Once),
		onError:   errHandler,
	}
}

func (c *bxConfigRepository) init() {
	c.initOnce.Do(func() {
		err := c.persistor.Load(c.data)

		if err != nil && os.IsNotExist(err) {
			c.setDefaults()
			err = c.persistor.Save(c.data)
		}

		if err != nil {
			c.onError(err)
		}
	})
}

func (c *bxConfigRepository) setDefaults() {
	c.data.PluginRepos = []models.PluginRepo{
		{
			Name: _DEFAULT_PLUGIN_REPO_NAME,
			URL:  _DEFAULT_PLUGIN_REPO_URL,
		},
	}
}

func (c *bxConfigRepository) read(cb func()) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	c.init()

	cb()
}

func (c *bxConfigRepository) write(cb func()) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.init()

	cb()

	err := c.persistor.Save(c.data)
	if err != nil {
		c.onError(err)
	}
}

func (c *bxConfigRepository) Region() (region string) {
	c.read(func() {
		region = c.data.Region
	})
	return
}

func (c *bxConfigRepository) PluginRepos() (repos []models.PluginRepo) {
	c.read(func() {
		repos = c.data.PluginRepos
	})
	return
}

func (c *bxConfigRepository) PluginRepo(name string) (models.PluginRepo, bool) {
	for _, r := range c.PluginRepos() {
		if strings.EqualFold(r.Name, name) {
			return r, true
		}
	}
	return models.PluginRepo{}, false
}

func (c *bxConfigRepository) HTTPTimeout() (timeout int) {
	c.read(func() {
		timeout = c.data.HTTPTimeout
	})
	return
}

func (c *bxConfigRepository) CLIInfoEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.CLIInfoEndpoint
	})

	if endpoint != "" {
		return endpoint
	}
	return _DEFAULT_CLI_INFO_ENDPOINT
}

func (c *bxConfigRepository) CheckCLIVersionDisabled() (disabled bool) {
	c.read(func() {
		disabled = c.data.CheckCLIVersionDisabled
	})
	return
}

func (c *bxConfigRepository) UsageStatsDisabled() (disabled bool) {
	c.read(func() {
		disabled = c.data.UsageStatsDisabled
	})
	return
}

func (c *bxConfigRepository) SetRegion(region string) {
	c.write(func() {
		c.data.Region = region
	})
}

func (c *bxConfigRepository) SetPluginRepo(pluginRepo models.PluginRepo) {
	c.write(func() {
		c.data.PluginRepos = append(c.data.PluginRepos, pluginRepo)
	})
}

func (c *bxConfigRepository) UnSetPluginRepo(repoName string) {
	c.write(func() {
		i := 0
		for ; i < len(c.data.PluginRepos); i++ {
			if strings.ToLower(c.data.PluginRepos[i].Name) == strings.ToLower(repoName) {
				break
			}
		}
		if i != len(c.data.PluginRepos) {
			c.data.PluginRepos = append(c.data.PluginRepos[:i], c.data.PluginRepos[i+1:]...)
		}
	})
}

func (c *bxConfigRepository) SetHTTPTimeout(timeout int) {
	c.write(func() {
		c.data.HTTPTimeout = timeout
	})
}

func (c *bxConfigRepository) SetCheckCLIVersionDisabled(disabled bool) {
	c.write(func() {
		c.data.CheckCLIVersionDisabled = disabled
	})
}

func (c *bxConfigRepository) SetCLIInfoEndpoint(endpoint string) {
	c.write(func() {
		c.data.CLIInfoEndpoint = endpoint
	})
}

func (c *bxConfigRepository) SetUsageStatsDisabled(disabled bool) {
	c.write(func() {
		c.data.UsageStatsDisabled = disabled
	})
}
