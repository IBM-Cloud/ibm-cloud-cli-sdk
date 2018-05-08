package core_config

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
	"github.com/fatih/structs"
)

type raw map[string]interface{}

func (r raw) Marshal() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r raw) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, r)
}

type BXConfigData struct {
	APIEndpoint                string
	ConsoleEndpoint            string
	Region                     string
	RegionID                   string
	RegionType                 string
	IAMEndpoint                string
	IAMToken                   string
	IAMRefreshToken            string
	Account                    models.Account
	ResourceGroup              models.ResourceGroup
	CFEETargeted               bool
	PluginRepos                []models.PluginRepo
	SSLDisabled                bool
	Locale                     string
	Trace                      string
	ColorEnabled               string
	HTTPTimeout                int
	CLIInfoEndpoint            string
	CheckCLIVersionDisabled    bool
	UsageStatsDisabled         bool
	SDKVersion                 string
	UpdateCheckInterval        time.Duration
	UpdateRetryCheckInterval   time.Duration
	UpdateNotificationInterval time.Duration
	raw                        raw
}

func NewBXConfigData() *BXConfigData {
	data := new(BXConfigData)
	data.raw = make(map[string]interface{})
	return data
}

func (data *BXConfigData) Marshal() ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func (data *BXConfigData) Unmarshal(bytes []byte) error {
	err := json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}

	var raw raw
	err = json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	data.raw = raw

	return nil
}

type bxConfig struct {
	data      *BXConfigData
	persistor configuration.Persistor
	initOnce  *sync.Once
	lock      sync.RWMutex
	onError   func(error)
}

func createBluemixConfigFromPersistor(persistor configuration.Persistor, errHandler func(error)) *bxConfig {
	return &bxConfig{
		data:      NewBXConfigData(),
		persistor: persistor,
		initOnce:  new(sync.Once),
		onError:   errHandler,
	}
}

func (c *bxConfig) init() {
	c.initOnce.Do(func() {
		err := c.persistor.Load(c.data)
		if err != nil {
			c.onError(err)
		}
	})
}

func (c *bxConfig) read(cb func()) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	c.init()

	cb()
}

func (c *bxConfig) write(cb func()) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.init()

	cb()

	c.data.SDKVersion = bluemix.Version.String()
	c.data.raw = structs.Map(c.data)

	err := c.persistor.Save(c.data)
	if err != nil {
		c.onError(err)
	}
}

func (c *bxConfig) writeRaw(cb func()) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.init()

	cb()

	err := c.persistor.Save(c.data.raw)
	if err != nil {
		c.onError(err)
	}
}

func (c *bxConfig) APIEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.APIEndpoint
	})
	return
}

func (c *bxConfig) HasAPIEndpoint() bool {
	return c.APIEndpoint() != ""
}

func (c *bxConfig) IsSSLDisabled() (disabled bool) {
	c.read(func() {
		disabled = c.data.SSLDisabled
	})
	return
}

func (c *bxConfig) ConsoleEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.ConsoleEndpoint
	})
	return
}

func (c *bxConfig) CurrentRegion() (region models.Region) {
	c.read(func() {
		region = models.Region{
			ID:   c.data.RegionID,
			Name: c.data.Region,
			Type: c.data.RegionType,
		}
	})
	return
}

func (c *bxConfig) CloudName() string {
	regionID := c.CurrentRegion().ID
	if regionID == "" {
		return ""
	}

	splits := strings.Split(regionID, ":")
	if len(splits) != 3 {
		return ""
	}

	customer := splits[0]
	if customer != "ibm" {
		return customer
	}

	deployment := splits[1]
	switch {
	case deployment == "yp":
		return "bluemix"
	case strings.HasPrefix(deployment, "ys"):
		return "staging"
	default:
		return ""
	}
}

func (c *bxConfig) CloudType() string {
	return c.CurrentRegion().Type
}

func (c *bxConfig) IAMEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.IAMEndpoint
	})
	return
}

func (c *bxConfig) IAMToken() (token string) {
	c.read(func() {
		token = c.data.IAMToken
	})
	return
}

func (c *bxConfig) IAMRefreshToken() (token string) {
	c.read(func() {
		token = c.data.IAMRefreshToken
	})
	return
}

func (c *bxConfig) UserEmail() (email string) {
	c.read(func() {
		email = NewIAMTokenInfo(c.data.IAMToken).UserEmail
	})
	return
}

func (c *bxConfig) IAMID() (guid string) {
	c.read(func() {
		guid = NewIAMTokenInfo(c.data.IAMToken).IAMID
	})
	return
}

func (c *bxConfig) IsLoggedIn() (loggedIn bool) {
	c.read(func() {
		loggedIn = c.data.IAMToken != ""
	})
	return
}

func (c *bxConfig) CurrentAccount() (account models.Account) {
	c.read(func() {
		account = c.data.Account
	})
	return
}

func (c *bxConfig) HasTargetedAccount() bool {
	return c.CurrentAccount().GUID != ""
}

func (c *bxConfig) IMSAccountID() string {
	return NewIAMTokenInfo(c.IAMToken()).Accounts.IMSAccountID
}

func (c *bxConfig) CurrentResourceGroup() (group models.ResourceGroup) {
	c.read(func() {
		group = c.data.ResourceGroup
	})
	return
}

func (c *bxConfig) HasTargetedResourceGroup() (hasGroup bool) {
	c.read(func() {
		hasGroup = c.data.ResourceGroup.GUID != "" && c.data.ResourceGroup.Name != ""
	})
	return
}

func (c *bxConfig) PluginRepos() (repos []models.PluginRepo) {
	c.read(func() {
		repos = c.data.PluginRepos
	})
	return
}

func (c *bxConfig) PluginRepo(name string) (models.PluginRepo, bool) {
	for _, r := range c.PluginRepos() {
		if strings.EqualFold(r.Name, name) {
			return r, true
		}
	}
	return models.PluginRepo{}, false
}

func (c *bxConfig) Locale() (locale string) {
	c.read(func() {
		locale = c.data.Locale
	})
	return
}

func (c *bxConfig) Trace() (trace string) {
	c.read(func() {
		trace = c.data.Trace
	})
	return
}

func (c *bxConfig) ColorEnabled() (enabled string) {
	c.read(func() {
		enabled = c.data.ColorEnabled
	})
	return
}

func (c *bxConfig) HTTPTimeout() (timeout int) {
	c.read(func() {
		timeout = c.data.HTTPTimeout
	})
	return
}

func (c *bxConfig) CLIInfoEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.CLIInfoEndpoint
	})

	return endpoint
}

func (c *bxConfig) CheckCLIVersionDisabled() (disabled bool) {
	c.read(func() {
		disabled = c.data.CheckCLIVersionDisabled
	})
	return
}

func (c *bxConfig) UpdateCheckInterval() (interval time.Duration) {
	c.read(func() {
		interval = c.data.UpdateCheckInterval
	})
	return
}

func (c *bxConfig) UpdateRetryCheckInterval() (interval time.Duration) {
	c.read(func() {
		interval = c.data.UpdateRetryCheckInterval
	})
	return
}

func (c *bxConfig) UpdateNotificationInterval() (interval time.Duration) {
	c.read(func() {
		interval = c.data.UpdateNotificationInterval
	})
	return
}

func (c *bxConfig) UsageStatsDisabled() (disabled bool) {
	c.read(func() {
		disabled = c.data.UsageStatsDisabled
	})
	return
}

func (c *bxConfig) SDKVersion() (version string) {
	c.read(func() {
		version = c.data.SDKVersion
	})
	return
}

func (c *bxConfig) CFEETargeted() (targeted bool) {
	c.read(func() {
		targeted = c.data.CFEETargeted
	})
	return
}

func (c *bxConfig) SetAPIEndpoint(endpoint string) {
	c.write(func() {
		c.data.APIEndpoint = endpoint
	})
}

func (c *bxConfig) SetConsoleEndpoint(endpoint string) {
	c.write(func() {
		c.data.ConsoleEndpoint = endpoint
	})
}

func (c *bxConfig) SetRegion(region models.Region) {
	c.write(func() {
		c.data.Region = region.Name
		c.data.RegionID = region.ID
		c.data.RegionType = region.Type
	})
}

func (c *bxConfig) SetIAMEndpoint(endpoint string) {
	c.write(func() {
		c.data.IAMEndpoint = endpoint
	})
}

func (c *bxConfig) SetIAMToken(token string) {
	c.writeRaw(func() {
		c.data.IAMToken = token
		c.data.raw["IAMToken"] = token
	})
}

func (c *bxConfig) SetIAMRefreshToken(token string) {
	c.writeRaw(func() {
		c.data.IAMRefreshToken = token
		c.data.raw["IAMRefreshToken"] = token
	})
}

func (c *bxConfig) SetAccount(account models.Account) {
	c.write(func() {
		c.data.Account = account
	})
}

func (c *bxConfig) SetResourceGroup(group models.ResourceGroup) {
	c.write(func() {
		c.data.ResourceGroup = group
	})
}

func (c *bxConfig) SetPluginRepo(pluginRepo models.PluginRepo) {
	c.write(func() {
		c.data.PluginRepos = append(c.data.PluginRepos, pluginRepo)
	})
}

func (c *bxConfig) UnsetPluginRepo(repoName string) {
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

func (c *bxConfig) SetSSLDisabled(disabled bool) {
	c.write(func() {
		c.data.SSLDisabled = disabled
	})
}

func (c *bxConfig) SetHTTPTimeout(timeout int) {
	c.write(func() {
		c.data.HTTPTimeout = timeout
	})
}

func (c *bxConfig) SetCheckCLIVersionDisabled(disabled bool) {
	c.write(func() {
		c.data.CheckCLIVersionDisabled = disabled
	})
}

func (c *bxConfig) SetUpdateCheckInterval(interval time.Duration) {
	c.write(func() {
		c.data.UpdateCheckInterval = interval
	})
}

func (c *bxConfig) SetUpdateRetryCheckInterval(interval time.Duration) {
	c.write(func() {
		c.data.UpdateRetryCheckInterval = interval
	})
}

func (c *bxConfig) SetUpdateNotificationInterval(interval time.Duration) {
	c.write(func() {
		c.data.UpdateNotificationInterval = interval
	})
}

func (c *bxConfig) SetCLIInfoEndpoint(endpoint string) {
	c.write(func() {
		c.data.CLIInfoEndpoint = endpoint
	})
}

func (c *bxConfig) SetUsageStatsDisabled(disabled bool) {
	c.write(func() {
		c.data.UsageStatsDisabled = disabled
	})
}

func (c *bxConfig) SetColorEnabled(enabled string) {
	c.write(func() {
		c.data.ColorEnabled = enabled
	})
}

func (c *bxConfig) SetLocale(locale string) {
	c.write(func() {
		c.data.Locale = locale
	})
}

func (c *bxConfig) SetTrace(trace string) {
	c.write(func() {
		c.data.Trace = trace
	})
}

func (c *bxConfig) SetCFEETargeted(targeted bool) {
	c.write(func() {
		c.data.CFEETargeted = targeted
	})
}

func (c *bxConfig) ClearSession() {
	c.write(func() {
		c.data.IAMToken = ""
		c.data.IAMRefreshToken = ""
		c.data.Account = models.Account{}
		c.data.ResourceGroup = models.ResourceGroup{}
	})
}

func (c *bxConfig) UnsetAPI() {
	c.write(func() {
		c.data.APIEndpoint = ""
		c.data.Region = ""
		c.data.RegionID = ""
		c.data.RegionType = ""
		c.data.ConsoleEndpoint = ""
		c.data.IAMEndpoint = ""
	})
}
