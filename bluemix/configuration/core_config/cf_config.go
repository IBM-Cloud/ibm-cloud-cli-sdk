package core_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/uaa"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	"github.com/fatih/structs"
)

type CFConfigData struct {
	AccessToken              string
	APIVersion               string
	AsyncTimeout             uint
	AuthorizationEndpoint    string
	ColorEnabled             string
	ConfigVersion            int
	DopplerEndpoint          string
	Locale                   string
	LogCacheEndPoint         string
	MinCLIVersion            string
	MinRecommendedCLIVersion string
	OrganizationFields       models.OrganizationFields
	PluginRepos              []models.PluginRepo
	RefreshToken             string
	RoutingAPIEndpoint       string
	SpaceFields              models.SpaceFields
	SSHOAuthClient           string
	SSLDisabled              bool
	Target                   string
	Trace                    string
	UaaEndpoint              string
	LoginAt                  time.Time
	UAAGrantType             string
	UAAOAuthClient           string
	UAAOAuthClientSecret     string
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
	if data.ConfigVersion != 3 {
		data.ConfigVersion = 4
	}
	return json.MarshalIndent(data, "", "  ")
}

func (data *CFConfigData) Unmarshal(bytes []byte) error {
	err := json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}

	// clear out config if version is not 3 or 4
	if data.ConfigVersion < 3 || data.ConfigVersion > 4 {
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

func (c *cfConfig) refreshToken(refreshToken string) (uaa.Token, error) {
	auth := uaa.NewClient(uaa.DefaultConfig(c.data.AuthorizationEndpoint), rest.NewClient())
	refreshedToken, err := auth.GetToken(uaa.RefreshTokenRequest(refreshToken))
	if err != nil {
		return uaa.Token{}, err
	}

	// return an error if refreshed token is invalid
	refreshedTokenInfo := NewUAATokenInfo(refreshedToken.AccessToken)
	if !refreshedTokenInfo.exists() {
		return uaa.Token{}, errors.New("could not refresh token")
	}

	return *refreshedToken, nil
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

func (c *cfConfig) AsyncTimeout() (timeout uint) {
	c.read(func() {
		timeout = c.data.AsyncTimeout
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

// IsLoggedIn will check if the user is logged in. To determine if the user is logged in both the
// token and the refresh token will be checked
// If token is near expiration or expired, and a refresh token is present attempt to refresh the token.
// If token refresh was successful, check if the new UAA token is valid. If valid, user is logged in,
// otherwise user can be considered logged out. If refresh failed, then user is considered logged out.
// If no refresh token is present, and token is expired, then user is considered logged out.
func (c *cfConfig) IsLoggedIn() (loggedIn bool) {
	if token, refresh := c.UAAToken(), c.UAARefreshToken(); token != "" || refresh != "" {
		uaaTokenInfo := NewUAATokenInfo(token)
		if uaaTokenInfo.hasExpired() && refresh != "" {
			refreshedToken, err := c.refreshToken(token)
			if err != nil {
				return false
			}

			uaaToken := fmt.Sprintf("%s %s", refreshedToken.TokenType, refreshedToken.AccessToken)
			c.SetUAAToken(uaaToken)
			c.SetUAARefreshToken(refreshedToken.RefreshToken)

			return true
		} else if uaaTokenInfo.hasExpired() && refresh == "" {
			return false
		} else {
			return true
		}
	}

	return false
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
		c.data.DopplerEndpoint = ""
	})
}

func (c *cfConfig) LoginAt() (loginAt time.Time) {
	c.read(func() {
		loginAt = c.data.LoginAt
	})
	return
}

func (c *cfConfig) SetLoginAt(loginAt time.Time) {
	c.write(func() {
		c.data.LoginAt = loginAt
	})
}

func (c *cfConfig) ClearSession() {
	c.write(func() {
		c.data.AccessToken = ""
		c.data.RefreshToken = ""
		c.data.OrganizationFields = models.OrganizationFields{}
		c.data.SpaceFields = models.SpaceFields{}
		c.data.LoginAt = time.Time{}
	})
}
