package core_config

import (
	cfconfiguration "code.cloudfoundry.org/cli/cf/configuration"
	cfcoreconfig "code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	cfmodels "code.cloudfoundry.org/cli/cf/models"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/models"
)

type CFConfigReader interface {
	APIEndpoint() string
	HasAPIEndpoint() bool
	APIVersion() string
	AuthenticationEndpoint() string
	// deprecate loggergator endpoint, use Doppler endpoint instead
	// LoggregatorEndpoint() string
	DopplerEndpoint() string
	UAAEndpoint() string
	RoutingAPIEndpoint() string
	SSHOAuthClient() string
	MinCFCLIVersion() string
	MinRecommendedCFCLIVersion() string

	Username() string
	UserGUID() string
	UserEmail() string
	IsLoggedIn() bool
	UAAToken() string
	UAARefreshToken() string

	OrganizationFields() models.OrganizationFields
	HasOrganization() bool
	SpaceFields() models.SpaceFields
	HasSpace() bool

	IsSSLDisabled() bool
	Locale() string
	Trace() string
	ColorEnabled() string
}

type CFConfigWriter interface {
	SetAPIVersion(string)
	SetAPIEndpoint(string)
	SetAuthenticationEndpoint(string)
	// deprecate loggergator endpoint, use Doppler endpoint instead
	// SetLoggregatorEndpoint(string)
	SetDopplerEndpoint(string)
	SetUAAEndpoint(string)
	SetRoutingAPIEndpoint(string)
	SetSSHOAuthClient(string)
	SetMinCFCLIVersion(string)
	SetMinRecommendedCFCLIVersion(string)

	SetUAAToken(string)
	SetUAARefreshToken(string)

	SetOrganizationFields(models.OrganizationFields)
	SetSpaceFields(models.SpaceFields)

	SetSSLDisabled(bool)
	SetLocale(string)
	SetTrace(string)
	SetColorEnabled(string)

	ReloadCF()
}

type CFConfigReadWriter interface {
	CFConfigReader
	CFConfigWriter
}

type cfConfigAdapter struct {
	persistor  cfconfiguration.Persistor
	errHandler func(error)

	cfcoreconfig.Repository
}

func createCFConfigAdapterFromPath(filepath string, errHandler func(error)) *cfConfigAdapter {
	return createCFConfigAdapterFromPersistor(cfconfiguration.NewDiskPersistor(filepath), errHandler)
}

func createCFConfigAdapterFromPersistor(persistor cfconfiguration.Persistor, errHandler func(error)) *cfConfigAdapter {
	return &cfConfigAdapter{
		persistor:  persistor,
		errHandler: errHandler,
		Repository: cfcoreconfig.NewRepositoryFromPersistor(persistor, errHandler),
	}
}

func (c *cfConfigAdapter) MinCFCLIVersion() string {
	return c.Repository.MinCLIVersion()
}

func (c *cfConfigAdapter) MinRecommendedCFCLIVersion() string {
	return c.Repository.MinRecommendedCLIVersion()
}

func (c *cfConfigAdapter) UAAEndpoint() string {
	return c.Repository.UaaEndpoint()
}

func (c *cfConfigAdapter) UAAToken() string {
	return c.Repository.AccessToken()
}

func (c *cfConfigAdapter) UAARefreshToken() string {
	return c.Repository.RefreshToken()
}

func (c *cfConfigAdapter) OrganizationFields() models.OrganizationFields {
	cforg := c.Repository.OrganizationFields()

	var org models.OrganizationFields
	org.Name = cforg.Name
	org.GUID = cforg.GUID
	org.QuotaDefinition.GUID = cforg.QuotaDefinition.GUID
	org.QuotaDefinition.Name = cforg.QuotaDefinition.Name
	org.QuotaDefinition.MemoryLimitInMB = cforg.QuotaDefinition.MemoryLimit
	org.QuotaDefinition.InstanceMemoryLimitInMB = cforg.QuotaDefinition.InstanceMemoryLimit
	org.QuotaDefinition.RoutesLimit = cforg.QuotaDefinition.RoutesLimit
	org.QuotaDefinition.ServicesLimit = cforg.QuotaDefinition.ServicesLimit
	org.QuotaDefinition.NonBasicServicesAllowed = cforg.QuotaDefinition.NonBasicServicesAllowed
	org.QuotaDefinition.AppInstanceLimit = cforg.QuotaDefinition.AppInstanceLimit
	return org
}

func (c *cfConfigAdapter) SpaceFields() models.SpaceFields {
	cfspace := c.Repository.SpaceFields()

	var space models.SpaceFields
	space.Name = cfspace.Name
	space.GUID = cfspace.GUID
	space.AllowSSH = cfspace.AllowSSH
	return space
}

func (c *cfConfigAdapter) SetMinCFCLIVersion(version string) {
	c.Repository.SetMinCLIVersion(version)
}

func (c *cfConfigAdapter) SetMinRecommendedCFCLIVersion(version string) {
	c.Repository.SetMinRecommendedCLIVersion(version)
}

func (c *cfConfigAdapter) SetUAAEndpoint(endpoint string) {
	c.Repository.SetUaaEndpoint(endpoint)
}

func (c *cfConfigAdapter) SetUAAToken(token string) {
	c.Repository.SetAccessToken(token)
}

func (c *cfConfigAdapter) SetUAARefreshToken(token string) {
	c.Repository.SetRefreshToken(token)
}

func (c *cfConfigAdapter) SetOrganizationFields(org models.OrganizationFields) {
	var cfOrg cfmodels.OrganizationFields

	cfOrg.GUID = org.GUID
	cfOrg.Name = org.Name
	cfOrg.QuotaDefinition.GUID = org.QuotaDefinition.GUID
	cfOrg.QuotaDefinition.Name = org.QuotaDefinition.Name
	cfOrg.QuotaDefinition.MemoryLimit = org.QuotaDefinition.MemoryLimitInMB
	cfOrg.QuotaDefinition.InstanceMemoryLimit = org.QuotaDefinition.InstanceMemoryLimitInMB
	cfOrg.QuotaDefinition.RoutesLimit = org.QuotaDefinition.RoutesLimit
	cfOrg.QuotaDefinition.ServicesLimit = org.QuotaDefinition.ServicesLimit
	cfOrg.QuotaDefinition.NonBasicServicesAllowed = org.QuotaDefinition.NonBasicServicesAllowed
	cfOrg.QuotaDefinition.AppInstanceLimit = org.QuotaDefinition.AppInstanceLimit

	c.Repository.SetOrganizationFields(cfOrg)
}

func (c *cfConfigAdapter) SetSpaceFields(space models.SpaceFields) {
	var cfSpace cfmodels.SpaceFields

	cfSpace.GUID = space.GUID
	cfSpace.Name = space.Name
	cfSpace.AllowSSH = space.AllowSSH

	c.Repository.SetSpaceFields(cfSpace)
}

func (c *cfConfigAdapter) ReloadCF() {
	c.Repository.Close()
	c.Repository = cfcoreconfig.NewRepositoryFromPersistor(c.persistor, c.errHandler)
}
