package core_config

import (
	"encoding/json"
	"sync"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/configuration"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/models"
	"github.com/fatih/structs"
)

type CFConfigData struct {
	ConfigVersion            int
	Target                   string
	APIVersion               string
	AuthorizationEndpoint    string
	LoggregatorEndpoint      string
	DopplerEndpoint          string
	UaaEndpoint              string
	RoutingAPIEndpoint       string
	AccessToken              string
	RefreshToken             string
	UAAOAuthClient           string
	UAAOAuthClientSecret     string
	SSHOAuthClient           string
	OrganizationFields       models.OrganizationFields
	SpaceFields              models.SpaceFields
	SSLDisabled              bool
	AsyncTimeout             uint
	Trace                    string
	ColorEnabled             string
	Locale                   string
	PluginRepos              []models.PluginRepo
	MinCLIVersion            string
	MinRecommendedCLIVersion string
	raw                      raw
}

func NewCFConfigData() *CFConfigData {
	data := new(CFConfigData)
	data.raw = make(map[string]interface{})

	data.UAAOAuthClient = "cf"
	data.UAAOAuthClientSecret = ""

	return data
}

func (data *CFConfigData) Marshal() ([]byte, error) {
	data.ConfigVersion = 3
	return json.MarshalIndent(data, "", "  ")
}

func (data *CFConfigData) Unmarshal(bytes []byte) error {
	err := json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}

	if data.ConfigVersion != 3 {
		*data = CFConfigData{raw: make(map[string]interface{})}
		return nil
	}

	var raw raw
	err = json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	data.raw = raw
	return nil
}

type cfConfig struct {
	data      *CFConfigData
	persistor configuration.Persistor
	initOnce  *sync.Once
	lock      sync.RWMutex
	onError   func(error)
}

func createCFConfigFromPersistor(persistor configuration.Persistor, errHandler func(error)) *cfConfig {
	data := NewCFConfigData()
	if !persistor.Exists() {
		data.PluginRepos = []models.PluginRepo{
			{
				Name: "CF-Community",
				URL:  "https://plugins.cloudfoundry.org",
			},
		}
	}

	return &cfConfig{
		data:      data,
		persistor: persistor,
		initOnce:  new(sync.Once),
		onError:   errHandler,
	}
}

func (c *cfConfig) init() {
	c.initOnce.Do(func() {
		err := c.persistor.Load(c.data)
		if err != nil {
			c.onError(err)
		}
	})
}

func (c *cfConfig) read(cb func()) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	c.init()

	cb()
}

func (c *cfConfig) write(cb func()) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.init()

	cb()

	c.data.raw = structs.Map(c.data)

	err := c.persistor.Save(c.data)
	if err != nil {
		c.onError(err)
	}
}

func (c *cfConfig) writeRaw(cb func()) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.init()

	cb()

	err := c.persistor.Save(c.data.raw)
	if err != nil {
		c.onError(err)
	}
}

func (c *cfConfig) APIVersion() (version string) {
	c.read(func() {
		version = c.data.APIVersion
	})
	return
}

func (c *cfConfig) APIEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.Target
	})
	return
}

func (c *cfConfig) HasAPIEndpoint() (hasEndpoint bool) {
	c.read(func() {
		hasEndpoint = c.data.APIVersion != "" && c.data.Target != ""
	})
	return
}

func (c *cfConfig) AuthenticationEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.AuthorizationEndpoint
	})
	return
}

func (c *cfConfig) DopplerEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.DopplerEndpoint
	})
	return
}

func (c *cfConfig) LoggregatorEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.LoggregatorEndpoint
	})
	return
}

func (c *cfConfig) UAAEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.UaaEndpoint
	})
	return
}

func (c *cfConfig) RoutingAPIEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.RoutingAPIEndpoint
	})
	return
}

// func (c *cfConfig) UAAOAuthClient() (client string) {
// 	c.read(func() {
// 		client = c.data.UAAOAuthClient
// 	})
// 	return
// }

// func (c *cfConfig) UAAOAuthClientSecret() (secret string) {
// 	c.read(func() {
// 		secret = c.data.UAAOAuthClientSecret
// 	})
// 	return
// }

func (c *cfConfig) SSHOAuthClient() (client string) {
	c.read(func() {
		client = c.data.SSHOAuthClient
	})
	return
}

func (c *cfConfig) MinCFCLIVersion() (version string) {
	c.read(func() {
		version = c.data.MinCLIVersion
	})
	return
}

func (c *cfConfig) MinRecommendedCFCLIVersion() (version string) {
	c.read(func() {
		version = c.data.MinRecommendedCLIVersion
	})
	return
}

func (c *cfConfig) UAAToken() (token string) {
	c.read(func() {
		token = c.data.AccessToken
	})
	return
}

func (c *cfConfig) UAARefreshToken() (token string) {
	c.read(func() {
		token = c.data.RefreshToken
	})
	return
}

func (c *cfConfig) UserEmail() (email string) {
	c.read(func() {
		email = NewUAATokenInfo(c.data.AccessToken).Email
	})
	return
}

func (c *cfConfig) UserGUID() (guid string) {
	c.read(func() {
		guid = NewUAATokenInfo(c.data.AccessToken).UserGUID
	})
	return
}

func (c *cfConfig) Username() (name string) {
	c.read(func() {
		name = NewUAATokenInfo(c.data.AccessToken).Username
	})
	return
}

func (c *cfConfig) IsLoggedIn() (loggedIn bool) {
	c.read(func() {
		loggedIn = c.data.AccessToken != ""
	})
	return
}

func (c *cfConfig) CurrentOrganization() (org models.OrganizationFields) {
	c.read(func() {
		org = c.data.OrganizationFields
	})
	return
}

func (c *cfConfig) HasTargetedOrganization() (hasOrg bool) {
	c.read(func() {
		hasOrg = c.data.OrganizationFields.GUID != "" && c.data.OrganizationFields.Name != ""
	})
	return
}

func (c *cfConfig) CurrentSpace() (space models.SpaceFields) {
	c.read(func() {
		space = c.data.SpaceFields
	})
	return
}

func (c *cfConfig) HasTargetedSpace() (hasSpace bool) {
	c.read(func() {
		hasSpace = c.data.SpaceFields.GUID != "" && c.data.SpaceFields.Name != ""
	})
	return
}

func (c *cfConfig) IsSSLDisabled() (isSSLDisabled bool) {
	c.read(func() {
		isSSLDisabled = c.data.SSLDisabled
	})
	return
}

func (c *cfConfig) Trace() (trace string) {
	c.read(func() {
		trace = c.data.Trace
	})
	return
}

func (c *cfConfig) ColorEnabled() (enabled string) {
	c.read(func() {
		enabled = c.data.ColorEnabled
	})
	return
}

func (c *cfConfig) Locale() (locale string) {
	c.read(func() {
		locale = c.data.Locale
	})
	return
}

func (c *cfConfig) SetAPIVersion(version string) {
	c.write(func() {
		c.data.APIVersion = version
	})
}

func (c *cfConfig) SetAPIEndpoint(endpoint string) {
	c.write(func() {
		c.data.Target = endpoint
	})
}

func (c *cfConfig) SetAuthenticationEndpoint(endpoint string) {
	c.write(func() {
		c.data.AuthorizationEndpoint = endpoint
	})
}

func (c *cfConfig) SetLoggregatorEndpoint(endpoint string) {
	c.write(func() {
		c.data.LoggregatorEndpoint = endpoint
	})
}

func (c *cfConfig) SetDopplerEndpoint(endpoint string) {
	c.write(func() {
		c.data.DopplerEndpoint = endpoint
	})
}

func (c *cfConfig) SetUAAEndpoint(endpoint string) {
	c.write(func() {
		c.data.UaaEndpoint = endpoint
	})
}

func (c *cfConfig) SetRoutingAPIEndpoint(endpoint string) {
	c.write(func() {
		c.data.RoutingAPIEndpoint = endpoint
	})
}

func (c *cfConfig) SetSSHOAuthClient(client string) {
	c.write(func() {
		c.data.SSHOAuthClient = client
	})
}

func (c *cfConfig) SetMinCFCLIVersion(version string) {
	c.write(func() {
		c.data.MinCLIVersion = version
	})
}

func (c *cfConfig) SetMinRecommendedCFCLIVersion(version string) {
	c.write(func() {
		c.data.MinRecommendedCLIVersion = version
	})
}

func (c *cfConfig) SetUAAToken(token string) {
	c.writeRaw(func() {
		c.data.AccessToken = token
		c.data.raw["AccessToken"] = token
	})
}

func (c *cfConfig) SetUAARefreshToken(token string) {
	c.writeRaw(func() {
		c.data.RefreshToken = token
		c.data.raw["RefreshToken"] = token
	})
}

func (c *cfConfig) SetOrganization(org models.OrganizationFields) {
	c.write(func() {
		c.data.OrganizationFields = org
	})
}

func (c *cfConfig) SetSpace(space models.SpaceFields) {
	c.write(func() {
		c.data.SpaceFields = space
	})
}

func (c *cfConfig) SetLocale(locale string) {
	c.write(func() {
		c.data.Locale = locale
	})
}

func (c *cfConfig) SetSSLDisabled(sslDisabled bool) {
	c.write(func() {
		c.data.SSLDisabled = sslDisabled
	})
}

func (c *cfConfig) SetTrace(trace string) {
	c.write(func() {
		c.data.Trace = trace
	})
}

func (c *cfConfig) SetColorEnabled(colorEnabled string) {
	c.write(func() {
		c.data.ColorEnabled = colorEnabled
	})
}

func (c *cfConfig) UnsetAPI() {
	c.write(func() {
		c.data.APIVersion = ""
		c.data.Target = ""
		c.data.AuthorizationEndpoint = ""
		c.data.UaaEndpoint = ""
		c.data.RoutingAPIEndpoint = ""
		c.data.LoggregatorEndpoint = ""
		c.data.DopplerEndpoint = ""
	})
}

func (c *cfConfig) ClearSession() {
	c.write(func() {
		c.data.AccessToken = ""
		c.data.RefreshToken = ""
		c.data.OrganizationFields = models.OrganizationFields{}
		c.data.SpaceFields = models.SpaceFields{}
	})
}
