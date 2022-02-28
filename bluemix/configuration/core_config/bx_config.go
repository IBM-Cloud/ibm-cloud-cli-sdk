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
	APIEndpoint                 string
	IsPrivate                   bool
	IsAccessFromVPC             bool
	ConsoleEndpoint             string
	ConsolePrivateEndpoint      string
	ConsolePrivateVPCEndpoint   string
	CloudType                   string
	CloudName                   string
	CRIType                     string
	Region                      string
	RegionID                    string
	IAMEndpoint                 string
	IAMPrivateEndpoint          string
	IAMPrivateVPCEndpoint       string
	IAMToken                    string
	IAMRefreshToken             string
	IsLoggedInAsCRI             bool
	Account                     models.Account
	Profile                     models.Profile
	ResourceGroup               models.ResourceGroup
	LoginAt                     time.Time
	CFEETargeted                bool
	CFEEEnvID                   string
	PluginRepos                 []models.PluginRepo
	SSLDisabled                 bool
	Locale                      string
	MessageOfTheDayTime         int64
	Trace                       string
	ColorEnabled                string
	HTTPTimeout                 int
	CLIInfoEndpoint             string // overwrite the cli info endpoint
	CheckCLIVersionDisabled     bool
	UsageStatsDisabled          bool // deprecated: use UsageStatsEnabled
	UsageStatsEnabled           bool
	UsageStatsEnabledLastUpdate time.Time
	SDKVersion                  string
	UpdateCheckInterval         time.Duration
	UpdateRetryCheckInterval    time.Duration
	UpdateNotificationInterval  time.Duration
	raw                         raw
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

func (c *bxConfig) IsPrivateEndpointEnabled() (isPrivate bool) {
	c.read(func() {
		isPrivate = c.data.IsPrivate
	})
	return
}

func (c *bxConfig) IsLoggedInAsCRI() (isCRI bool) {
	c.read(func() {
		isCRI = c.data.IsLoggedInAsCRI
	})
	return
}

func (c *bxConfig) IsAccessFromVPC() (isVPC bool) {
	c.read(func() {
		isVPC = c.data.IsAccessFromVPC
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

func (c *bxConfig) ConsoleEndpoints() (endpoints models.Endpoints) {
	c.read(func() {
		endpoints.PublicEndpoint = c.data.ConsoleEndpoint
		endpoints.PrivateEndpoint = c.data.ConsolePrivateEndpoint
		endpoints.PrivateVPCEndpoint = c.data.ConsolePrivateVPCEndpoint
	})
	return
}

func (c *bxConfig) CurrentRegion() (region models.Region) {
	c.read(func() {
		region = models.Region{
			MCCPID: c.data.RegionID,
			Name:   c.data.Region,
		}
	})
	return
}

func (c *bxConfig) HasTargetedRegion() bool {
	r := c.CurrentRegion()
	return r.Name != ""
}

func (c *bxConfig) CloudName() (cname string) {
	c.read(func() {
		cname = c.data.CloudName
	})
	return
}

func (c *bxConfig) CloudType() (ctype string) {
	c.read(func() {
		ctype = c.data.CloudType
	})
	return
}

func (c *bxConfig) CRIType() (criType string) {
	c.read(func() {
		criType = c.data.CRIType
	})
	return
}

func (c *bxConfig) IAMEndpoints() (endpoints models.Endpoints) {
	c.read(func() {
		endpoints.PublicEndpoint = c.data.IAMEndpoint
		endpoints.PrivateEndpoint = c.data.IAMPrivateEndpoint
		endpoints.PrivateVPCEndpoint = c.data.IAMPrivateVPCEndpoint
	})
	return
}

func (c *bxConfig) LoginAt() (loginAt time.Time) {
	c.read(func() {
		loginAt = c.data.LoginAt
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

func (c *bxConfig) UserDisplayText() (text string) {
	c.read(func() {
		token := NewIAMTokenInfo(c.data.IAMToken)
		if token.UserEmail != "" {
			text = token.UserEmail
		} else {
			text = token.Subject
		}
	})
	return
}

func (c *bxConfig) IAMID() (guid string) {
	c.read(func() {
		guid = NewIAMTokenInfo(c.data.IAMToken).IAMID
	})
	return
}

// IsLoggedIn will check if the user is logged in. To determine if the user is logged in both the
// token and the refresh token will be checked
// If token is near expiration or expired, and a refresh token is present attempt to refresh the token.
// If token refresh was successful, check if the new IAM token is valid. If valid, user is logged in,
// otherwise user can be considered logged out. If refresh failed, then user is considered logged out.
// If no refresh token is present, and token is expired, then user is considered logged out.
func (c *bxConfig) IsLoggedIn() bool {
	if token, refresh := c.IAMToken(), c.IAMRefreshToken(); token != "" || refresh != "" {
		iamTokenInfo := NewIAMTokenInfo(token)
		if iamTokenInfo.hasExpired() && refresh != "" {
			repo := newRepository(c, nil)
			if _, err := repo.RefreshIAMToken(); err != nil {
				return false
			}
			// Check again to make sure that the new token has not expired
			if iamTokenInfo = NewIAMTokenInfo(c.IAMToken()); iamTokenInfo.hasExpired() {
				return false
			}

			return true
		} else if iamTokenInfo.hasExpired() && refresh == "" {
			return false
		} else {
			return true
		}
	}

	return false
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

func (c *bxConfig) HasTargetedProfile() bool {
	return c.CurrentProfile().ID != ""
}

func (c *bxConfig) HasTargetedComputeResource() bool {
	authn := c.CurrentProfile().ComputeResource
	return authn.ID != ""
}

func (c *bxConfig) IMSAccountID() string {
	return NewIAMTokenInfo(c.IAMToken()).Accounts.IMSAccountID
}

func (c *bxConfig) CurrentProfile() (profile models.Profile) {
	c.read(func() {
		profile = c.data.Profile
	})
	return
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

// CheckMessageOfTheDay will return true if the message-of-the-day
// endpoint has not been check in the past 24 hours.
func (c *bxConfig) CheckMessageOfTheDay() bool {
	var lastCheck int64
	c.read(func() {
		lastCheck = c.data.MessageOfTheDayTime
	})

	if lastCheck <= 0 {
		return true
	}

	currentDate := time.Now()
	lastCheckDate := time.Unix(lastCheck, 0)

	// If MOD has not been checked in over 1 day
	return currentDate.Sub(lastCheckDate).Hours() > 24
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

func (c *bxConfig) UsageStatsEnabled() (enabled bool) {
	c.read(func() {
		enabled = !c.data.UsageStatsEnabledLastUpdate.IsZero() && c.data.UsageStatsEnabled
	})
	return
}

func (c *bxConfig) UsageStatsEnabledLastUpdate() (lastUpdate time.Time) {
	c.read(func() {
		lastUpdate = c.data.UsageStatsEnabledLastUpdate
	})
	return
}

func (c *bxConfig) SDKVersion() (version string) {
	c.read(func() {
		version = c.data.SDKVersion
	})
	return
}

func (c *bxConfig) HasTargetedCFEE() (targeted bool) {
	c.read(func() {
		targeted = c.data.CFEETargeted
	})
	return
}

func (c *bxConfig) CFEEEnvID() (envID string) {
	c.read(func() {
		envID = c.data.CFEEEnvID
	})
	return
}

func (c *bxConfig) SetAPIEndpoint(endpoint string) {
	c.write(func() {
		c.data.APIEndpoint = endpoint
	})
}

func (c *bxConfig) SetPrivateEndpointEnabled(isPrivate bool) {
	c.write(func() {
		c.data.IsPrivate = isPrivate
	})
}

func (c *bxConfig) SetAccessFromVPC(isVPC bool) {
	c.write(func() {
		c.data.IsAccessFromVPC = isVPC
	})
}

func (c *bxConfig) SetConsoleEndpoints(endpoint models.Endpoints) {
	c.write(func() {
		c.data.ConsoleEndpoint = endpoint.PublicEndpoint
		c.data.ConsolePrivateEndpoint = endpoint.PrivateEndpoint
		c.data.ConsolePrivateVPCEndpoint = endpoint.PrivateVPCEndpoint
	})
}

func (c *bxConfig) SetRegion(region models.Region) {
	c.write(func() {
		c.data.Region = region.Name
		c.data.RegionID = region.MCCPID
	})
}

func (c *bxConfig) SetIAMEndpoints(endpoints models.Endpoints) {
	c.write(func() {
		c.data.IAMEndpoint = endpoints.PublicEndpoint
		c.data.IAMPrivateEndpoint = endpoints.PrivateEndpoint
		c.data.IAMPrivateVPCEndpoint = endpoints.PrivateVPCEndpoint
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

func (c *bxConfig) SetProfile(profile models.Profile) {
	c.write(func() {
		c.data.Profile = profile
	})
}

func (c *bxConfig) SetCRIType(criType string) {
	c.write(func() {
		c.data.CRIType = criType
	})
}

func (c *bxConfig) SetIsLoggedInAsCRI(isCRI bool) {
	c.write(func() {
		c.data.IsLoggedInAsCRI = isCRI
	})
}

func (c *bxConfig) SetResourceGroup(group models.ResourceGroup) {
	c.write(func() {
		c.data.ResourceGroup = group
	})
}

func (c *bxConfig) SetLoginAt(loginAt time.Time) {
	c.write(func() {
		c.data.LoginAt = loginAt
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

func (c *bxConfig) SetUsageStatsEnabled(enabled bool) {
	c.write(func() {
		c.data.UsageStatsEnabled = enabled
		c.data.UsageStatsEnabledLastUpdate = time.Now()
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

func (c *bxConfig) SetCFEEEnvID(envID string) {
	c.write(func() {
		c.data.CFEEEnvID = envID
	})
}

func (c *bxConfig) SetCloudType(ctype string) {
	c.write(func() {
		c.data.CloudType = ctype
	})
}

func (c *bxConfig) SetCloudName(cname string) {
	c.write(func() {
		c.data.CloudName = cname
	})
}

func (c *bxConfig) SetMessageOfTheDayTime() {
	c.write(func() {
		c.data.MessageOfTheDayTime = time.Now().Unix()
	})
}

func (c *bxConfig) ClearSession() {
	c.write(func() {
		c.data.IAMToken = ""
		c.data.IAMRefreshToken = ""
		c.data.Account = models.Account{}
		c.data.Profile = models.Profile{}
		c.data.CRIType = ""
		c.data.IsLoggedInAsCRI = false
		c.data.ResourceGroup = models.ResourceGroup{}
		c.data.LoginAt = time.Time{}
	})
}

func (c *bxConfig) UnsetAPI() {
	c.write(func() {
		c.data.APIEndpoint = ""
		c.data.SSLDisabled = false
		c.data.IsPrivate = false
		c.data.IsAccessFromVPC = false
		c.data.Region = ""
		c.data.RegionID = ""
		c.data.ConsoleEndpoint = ""
		c.data.ConsolePrivateEndpoint = ""
		c.data.ConsolePrivateVPCEndpoint = ""
		c.data.IAMEndpoint = ""
		c.data.IAMPrivateEndpoint = ""
		c.data.IAMPrivateVPCEndpoint = ""
		c.data.CloudName = ""
		c.data.CloudType = ""
	})
}
