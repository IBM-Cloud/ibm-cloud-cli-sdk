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

type BXConfigData struct {
	Region                  string
	IAMToken                string
	IAMRefreshToken         string
	Account                 models.Account
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

func (c *bxConfigRepository) IAMToken() (token string) {
	c.read(func() {
		token = c.data.IAMToken
	})
	return
}

func (c *bxConfigRepository) IAMRefreshToken() (token string) {
	c.read(func() {
		token = c.data.IAMRefreshToken
	})
	return
}

func (c *bxConfigRepository) Account() (account models.Account) {
	c.read(func() {
		account = c.data.Account
	})
	return
}

func (c *bxConfigRepository) HasAccount() bool {
	return c.Account().GUID != ""
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

func (c *bxConfigRepository) Locale() (locale string) {
	c.read(func() {
		locale = c.data.Locale
	})
	return
}

func (c *bxConfigRepository) Trace() (trace string) {
	c.read(func() {
		trace = c.data.Trace
	})
	return
}

func (c *bxConfigRepository) ColorEnabled() (enabled string) {
	c.read(func() {
		enabled = c.data.ColorEnabled
	})
	return
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

func (c *bxConfigRepository) SetIAMToken(token string) {
	c.write(func() {
		c.data.IAMToken = token
	})
}

func (c *bxConfigRepository) SetIAMRefreshToken(token string) {
	c.write(func() {
		c.data.IAMRefreshToken = token
	})
}

func (c *bxConfigRepository) SetAccount(account models.Account) {
	c.write(func() {
		c.data.Account = account
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

func (c *bxConfigRepository) SetColorEnabled(enabled string) {
	c.write(func() {
		c.data.ColorEnabled = enabled
	})
}

func (c *bxConfigRepository) SetLocale(locale string) {
	c.write(func() {
		c.data.Locale = locale
	})
}

func (c *bxConfigRepository) SetTrace(trace string) {
	c.write(func() {
		c.data.Trace = trace
	})
}

func (c *bxConfigRepository) ClearSession() {
	c.write(func() {
		c.data.IAMToken = ""
		c.data.IAMRefreshToken = ""
		c.data.Account = models.Account{}
	})
}

func (c *bxConfigRepository) ClearAPICache() {
	c.write(func() {
		c.data.Region = ""
	})
}
