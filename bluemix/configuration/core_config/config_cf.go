package core_config

import (
	cfconfiguration "code.cloudfoundry.org/cli/cf/configuration"
	cfcoreconfig "code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	cfmodels "code.cloudfoundry.org/cli/cf/models"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/models"
)

type CFReader interface {
	APIEndpoint() string
	HasAPIEndpoint() bool
	APIVersion() string
	AuthenticationEndpoint() string
	LoggregatorEndpoint() string
	DopplerEndpoint() string
	UaaEndpoint() string
	RoutingAPIEndpoint() string
	SSHOAuthClient() string

	Username() string
	UserGUID() string
	UserEmail() string
	IsLoggedIn() bool
	AccessToken() string
	RefreshToken() string

	OrganizationFields() models.OrganizationFields
	HasOrganization() bool
	SpaceFields() models.SpaceFields
	HasSpace() bool

	IsSSLDisabled() bool
}

type CFWriter interface {
	SetAPIVersion(string)
	SetAPIEndpoint(string)
	SetAuthenticationEndpoint(string)
	SetLoggregatorEndpoint(string)
	SetDopplerEndpoint(string)
	SetUaaEndpoint(string)
	SetRoutingAPIEndpoint(string)
	SetSSHOAuthClient(string)
	SetAccessToken(string)
	SetRefreshToken(string)
	SetOrganizationFields(models.OrganizationFields)
	SetSpaceFields(models.SpaceFields)
	SetSSLDisabled(bool)
	Reload()
}

type CFReadWriter interface {
	CFReader
	CFWriter
}

type cfConfigAdapter struct {
	persistor  cfconfiguration.Persistor
	errHandler func(error)

	cfcoreconfig.Repository
}

func NewCFConfigAdapterFromPath(filepath string, errHandler func(error)) *cfConfigAdapter {
	return NewCFConfigAdapterFromPersistor(cfconfiguration.NewDiskPersistor(filepath), errHandler)
}

func NewCFConfigAdapterFromPersistor(persistor cfconfiguration.Persistor, errHandler func(error)) *cfConfigAdapter {
	return &cfConfigAdapter{
		persistor:  persistor,
		errHandler: errHandler,
		Repository: cfcoreconfig.NewRepositoryFromPersistor(persistor, errHandler),
	}
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

func (c *cfConfigAdapter) Reload() {
	c.Repository.Close()
	c.Repository = cfcoreconfig.NewRepositoryFromPersistor(c.persistor, c.errHandler)
}
