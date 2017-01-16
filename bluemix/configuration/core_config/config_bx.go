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

type BXAPICache struct {
	Target string
	Region string
}

type BXConfigData struct {
	BXAPICache

	PluginRepos             []models.PluginRepo
	Locale                  string
	Trace                   string
	ColorEnabled            string
	HTTPTimeout             int
	CLIInfoEndpoint         string
	CheckCLIVersionDisabled bool
	UsageStatsDisabled      bool
}

type BXReader interface {
	APICache() BXAPICache
	PluginRepos() []models.PluginRepo
	PluginRepo(string) (models.PluginRepo, bool)
	Locale() string
	Trace() string
	ColorEnabled() string
	HTTPTimeout() int
	CLIInfoEndpoint() string
	CheckCLIVersionDisabled() bool
	UsageStatsDisabled() bool
}

type BXWriter interface {
	ClearAPICache()
	SetAPICache(BXAPICache)
	SetCheckCLIVersionDisabled(bool)
	SetCLIInfoEndpoint(string)
	SetPluginRepo(models.PluginRepo)
	UnSetPluginRepo(string)
	SetLocale(string)
	SetTrace(string)
	SetColorEnabled(string)
	SetHTTPTimeout(int)
	SetUsageStatsDisabled(bool)
}

type BXReadWriter interface {
	BXReader
	BXWriter
}

type bxConfigRepository struct {
	configData *BXConfigData
	persistor  configuration.Persistor
	initOnce   *sync.Once
	lock       sync.RWMutex
	onError    func(error)
}

func NewBluemixConfigFromPath(configPath string, errHandler func(error)) BXReadWriter {
	return NewBluemixConfigFromPersistor(configuration.NewDiskPersistor(configPath), errHandler)
}

func NewBluemixConfigFromPersistor(persistor configuration.Persistor, errHandler func(error)) BXReadWriter {
	return &bxConfigRepository{
		configData: new(BXConfigData),
		persistor:  persistor,
		initOnce:   new(sync.Once),
		onError:    errHandler,
	}
}

func (bc *bxConfigRepository) init() {
	bc.initOnce.Do(func() {
		err := bc.persistor.Load(bc.configData)

		if err != nil && os.IsNotExist(err) {
			bc.setDefaults()
			err = bc.persistor.Save(bc.configData)
		}

		if err != nil {
			bc.onError(err)
		}
	})
}

func (bc *bxConfigRepository) setDefaults() {
	bc.configData.PluginRepos = []models.PluginRepo{
		{
			Name: _DEFAULT_PLUGIN_REPO_NAME,
			URL:  _DEFAULT_PLUGIN_REPO_URL,
		},
	}
}

func (bc *bxConfigRepository) read(cb func()) {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	bc.init()

	cb()
}

func (bc *bxConfigRepository) write(cb func()) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	bc.init()

	cb()

	err := bc.persistor.Save(bc.configData)
	if err != nil {
		bc.onError(err)
	}
}

func (bc *bxConfigRepository) APICache() (cache BXAPICache) {
	bc.read(func() {
		cache = bc.configData.BXAPICache
	})
	return
}

func (bc *bxConfigRepository) PluginRepos() (repos []models.PluginRepo) {
	bc.read(func() {
		repos = bc.configData.PluginRepos
	})
	return
}

func (bc *bxConfigRepository) PluginRepo(name string) (models.PluginRepo, bool) {
	for _, r := range bc.PluginRepos() {
		if strings.EqualFold(r.Name, name) {
			return r, true
		}
	}
	return models.PluginRepo{}, false
}

func (bc *bxConfigRepository) Locale() (locale string) {
	bc.read(func() {
		locale = bc.configData.Locale
	})
	return
}

func (bc *bxConfigRepository) Trace() (trace string) {
	bc.read(func() {
		trace = bc.configData.Trace
	})
	return
}

func (bc *bxConfigRepository) ColorEnabled() (colorEnabled string) {
	bc.read(func() {
		colorEnabled = bc.configData.ColorEnabled
	})
	return
}

func (bc *bxConfigRepository) HTTPTimeout() (timeout int) {
	bc.read(func() {
		timeout = bc.configData.HTTPTimeout
	})
	return
}

func (bc *bxConfigRepository) CLIInfoEndpoint() (endpoint string) {
	bc.read(func() {
		endpoint = bc.configData.CLIInfoEndpoint
	})

	if endpoint != "" {
		return endpoint
	}
	return _DEFAULT_CLI_INFO_ENDPOINT
}

func (bc *bxConfigRepository) CheckCLIVersionDisabled() (disabled bool) {
	bc.read(func() {
		disabled = bc.configData.CheckCLIVersionDisabled
	})
	return
}

func (bc *bxConfigRepository) UsageStatsDisabled() (disabled bool) {
	bc.read(func() {
		disabled = bc.configData.UsageStatsDisabled
	})
	return
}

func (bc *bxConfigRepository) ClearAPICache() {
	bc.write(func() {
		bc.configData.Region = ""
	})
}

func (bc *bxConfigRepository) SetAPICache(cache BXAPICache) {
	bc.write(func() {
		bc.configData.BXAPICache = cache
	})
}

func (bc *bxConfigRepository) SetPluginRepo(pluginRepo models.PluginRepo) {
	bc.write(func() {
		bc.configData.PluginRepos = append(bc.configData.PluginRepos, pluginRepo)
	})
}

func (bc *bxConfigRepository) UnSetPluginRepo(repoName string) {
	bc.write(func() {
		i := 0
		for ; i < len(bc.configData.PluginRepos); i++ {
			if strings.ToLower(bc.configData.PluginRepos[i].Name) == strings.ToLower(repoName) {
				break
			}
		}
		if i != len(bc.configData.PluginRepos) {
			bc.configData.PluginRepos = append(bc.configData.PluginRepos[:i], bc.configData.PluginRepos[i+1:]...)
		}
	})
}

func (bc *bxConfigRepository) SetLocale(locale string) {
	bc.write(func() {
		bc.configData.Locale = locale
	})
}

func (bc *bxConfigRepository) SetTrace(trace string) {
	bc.write(func() {
		bc.configData.Trace = trace
	})
}

func (bc *bxConfigRepository) SetColorEnabled(colorEnabled string) {
	bc.write(func() {
		bc.configData.ColorEnabled = colorEnabled
	})
}

func (bc *bxConfigRepository) SetHTTPTimeout(timeout int) {
	bc.write(func() {
		bc.configData.HTTPTimeout = timeout
	})
}

func (bc *bxConfigRepository) SetCheckCLIVersionDisabled(disabled bool) {
	bc.write(func() {
		bc.configData.CheckCLIVersionDisabled = disabled
	})
}

func (bc *bxConfigRepository) SetCLIInfoEndpoint(endpoint string) {
	bc.write(func() {
		bc.configData.CLIInfoEndpoint = endpoint
	})
}

func (bc *bxConfigRepository) SetUsageStatsDisabled(disabled bool) {
	bc.write(func() {
		bc.configData.UsageStatsDisabled = disabled
	})
}
