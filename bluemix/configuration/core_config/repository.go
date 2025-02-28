// Package core_config provides functions to load core configuration.
// The package is for internal only.
package core_config

import (
	"fmt"
	"os"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/iam"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/vpc"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/config_helpers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
)

type Repository interface {
	APIEndpoint() string
	HasAPIEndpoint() bool
	IsPrivateEndpointEnabled() bool
	IsAccessFromVPC() bool
	ConsoleEndpoints() models.Endpoints
	IAMEndpoint() string
	IAMEndpoints() models.Endpoints
	CloudName() string
	CloudType() string
	CurrentRegion() models.Region
	HasTargetedRegion() bool
	IAMToken() string
	IAMRefreshToken() string
	IsLoggedIn() bool
	IsLoggedInWithServiceID() bool
	IsLoggedInAsProfile() bool
	IsLoggedInAsCRI() bool
	UserEmail() string
	// UserDisplayText is the human readable ID for logged-in users which include non-human IDs
	UserDisplayText() string
	IAMID() string
	CurrentAccount() models.Account
	HasTargetedAccount() bool
	HasTargetedProfile() bool
	HasTargetedComputeResource() bool
	IMSAccountID() string
	CurrentProfile() models.Profile
	CurrentResourceGroup() models.ResourceGroup
	// CRIType returns the type of compute resource the user logged in as, if applicable. Valid values are `IKS`, `VPC`, or `OTHER`
	CRIType() string
	HasTargetedResourceGroup() bool
	PluginRepos() []models.PluginRepo
	PluginRepo(string) (models.PluginRepo, bool)
	IsSSLDisabled() bool
	TypeOfSSO() string
	AlphaCommandsEnabled() string
	AssumedTrustedProfileId() string
	FallbackAccount() models.Account
	FallbackIAMToken() string
	FallbackIAMRefreshToken() string
	HTTPTimeout() int
	CLIInfoEndpoint() string
	CheckCLIVersionDisabled() bool
	UpdateCheckInterval() time.Duration
	UpdateRetryCheckInterval() time.Duration
	UpdateNotificationInterval() time.Duration
	// VPCCRITokenURL() returns the value specified by the environment variable 'IBMCLOUD_CR_VPC_URL', if set.
	// Otherwise, the default VPC auth url specified by the constant `DefaultServerEndpoint` is returned
	VPCCRITokenURL() string

	// UsageSatsDisabled returns whether the usage statistics data collection is disabled or not
	// Deprecated: use UsageSatsEnabled instead. We change to disable usage statistics by default,
	// So this property will not be used anymore
	UsageStatsDisabled() bool
	// UsageSatsEnabled returns whether the usage statistics data collection is enabled or not
	UsageStatsEnabled() bool
	// UsageStatsEnabledLastUpdate returns last time when `UsageStatsEnabled` was updated
	UsageStatsEnabledLastUpdate() time.Time
	Locale() string
	LoginAt() time.Time
	Trace() string
	ColorEnabled() string
	SDKVersion() string

	UnsetAPI()
	RefreshIAMToken() (string, error)
	SetAPIEndpoint(string)
	SetPrivateEndpointEnabled(bool)
	SetAccessFromVPC(bool)
	SetConsoleEndpoints(models.Endpoints)
	SetIAMEndpoints(models.Endpoints)
	SetCloudType(string)
	SetCloudName(string)
	SetCRIType(string)
	SetIsLoggedInAsCRI(bool)
	SetRegion(models.Region)
	SetIAMToken(string)
	SetIAMRefreshToken(string)
	SetFallbackIAMTokens(string, string)
	SetFallbackAccount(string, string, string)
	SetAssumedTrustedProfileId(string)
	ClearSession()
	SetAccount(models.Account)
	SetProfile(models.Profile)
	SetResourceGroup(models.ResourceGroup)
	SetLoginAt(loginAt time.Time)
	SetCheckCLIVersionDisabled(bool)
	SetCLIInfoEndpoint(string)
	SetPluginRepo(models.PluginRepo)
	UnsetPluginRepo(string)
	SetSSLDisabled(bool)
	SetTypeOfSSO(string)
	SetHTTPTimeout(int)
	// SetUsageSatsDisabled disable or enable usage statistics data collection
	// Deprecated: use SetUsageSatsEnabled instead
	SetUsageStatsDisabled(bool)
	// SetUsageSatsEnabled enable or disable usage statistics data collection
	SetUsageStatsEnabled(bool)
	SetUpdateCheckInterval(time.Duration)
	SetUpdateRetryCheckInterval(time.Duration)
	SetUpdateNotificationInterval(time.Duration)
	SetLocale(string)
	SetTrace(string)
	SetColorEnabled(string)
	SetAlphaCommandsEnabled(string)

	CheckMessageOfTheDay() bool
	SetMessageOfTheDayTime()

	SetLastSessionUpdateTime()
	LastSessionUpdateTime() (session int64)

	SetPaginationURLs(paginationURLs []models.PaginationURL)
	ClearPaginationURLs()
	AddPaginationURL(lastIndex int, nextURL string)
	PaginationURLs() []models.PaginationURL
}

// Deprecated
type ReadWriter interface {
	Repository
}

type repository struct {
	*bxConfig
}

func newRepository(bx *bxConfig) repository {
	return repository{
		bxConfig: bx,
	}
}

func (c repository) IsLoggedIn() bool {
	return c.bxConfig.IsLoggedIn()
}

func (c repository) IsLoggedInWithServiceID() bool {
	return c.bxConfig.IsLoggedIn() && NewIAMTokenInfo(c.IAMToken()).SubjectType == SubjectTypeServiceID
}

func (c repository) IsLoggedInAsProfile() bool {
	return c.bxConfig.IsLoggedIn() && NewIAMTokenInfo(c.IAMToken()).SubjectType == SubjectTypeTrustedProfile
}

func (c repository) VPCCRITokenURL() string {
	if env := bluemix.EnvCRVpcUrl.Get(); env != "" {
		return env
	}

	// default server endpoint is a constant value in vpc authenticator
	return vpc.DefaultServerEndpoint
}

func (c repository) IAMEndpoint() string {
	if c.IsPrivateEndpointEnabled() {
		if c.IsAccessFromVPC() {
			// return VPC endpoint
			return c.IAMEndpoints().PrivateVPCEndpoint
		} else {
			// return CSE endpoint
			return c.IAMEndpoints().PrivateEndpoint
		}
	}
	return c.IAMEndpoints().PublicEndpoint
}

func (c repository) RefreshIAMToken() (string, error) {
	var ret string
	// confirm user is logged in as a VPC compute resource identity
	isLoggedInAsCRI := c.IsLoggedInAsCRI()
	criType := c.CRIType()
	if isLoggedInAsCRI && criType == "VPC" {
		token, err := c.fetchNewIAMTokenUsingVPCAuth()
		if err != nil {
			return "", err
		}

		ret = fmt.Sprintf("%s %s", token.TokenType, token.AccessToken)
		c.SetIAMToken(ret)
		// this should be empty for vpc vsi tokens
		c.SetIAMRefreshToken(token.RefreshToken)
	} else {
		iamEndpoint := os.Getenv("IAM_ENDPOINT")
		if iamEndpoint == "" {
			iamEndpoint = c.IAMEndpoint()
		}
		if iamEndpoint == "" {
			return "", fmt.Errorf("IAM endpoint is not set")
		}

		auth := iam.NewClient(iam.DefaultConfig(iamEndpoint), rest.NewClient())
		token, err := auth.GetToken(iam.RefreshTokenRequest(c.IAMRefreshToken()))
		if err != nil {
			return "", err
		}

		ret = fmt.Sprintf("%s %s", token.TokenType, token.AccessToken)
		c.SetIAMToken(ret)
		c.SetIAMRefreshToken(token.RefreshToken)
	}

	c.SetLastSessionUpdateTime()

	return ret, nil
}

func (c repository) fetchNewIAMTokenUsingVPCAuth() (*iam.Token, error) {
	// create a vpc client using default configuration
	client := vpc.NewClient(vpc.DefaultConfig(c.VPCCRITokenURL(), vpc.DefaultMetadataServiceVersion), rest.NewClient())
	// fetch an instance identity token from the metadata server
	identityToken, err := client.GetInstanceIdentityToken()
	if err != nil {
		return nil, err
	}

	// get the existing targeted IAM trusted profile ID of the CLI session
	targetProfile := c.CurrentProfile()
	profileID := targetProfile.ID
	if profileID == "" {
		return nil, fmt.Errorf("Trusted profile not set in configuration")
	}
	// prepare IAM token request using the existing targeted profile.
	req, err := vpc.NewIAMAccessTokenRequest(profileID, "", identityToken.AccessToken)
	if err != nil {
		return nil, err
	}

	// get the new access token
	iamToken, err := client.GetIAMAccessToken(req)
	if err != nil {
		return nil, err
	}

	return iamToken, nil
}

func (c repository) UserEmail() string {
	return c.bxConfig.UserEmail()
}

func (c repository) SetSSLDisabled(disabled bool) {
	c.bxConfig.SetSSLDisabled(disabled)
}

func (c repository) SetColorEnabled(enabled string) {
	c.bxConfig.SetColorEnabled(enabled)
}

func (c repository) SetTrace(trace string) {
	c.bxConfig.SetTrace(trace)
}

func (c repository) SetLocale(locale string) {
	c.bxConfig.SetLocale(locale)
}

func (c repository) UnsetAPI() {
	c.bxConfig.UnsetAPI()
}

func (c repository) ClearSession() {
	c.bxConfig.ClearSession()
}

func (c repository) LastSessionUpdateTime() (session int64) {
	return c.bxConfig.LastSessionUpdateTime()
}

func (c repository) SetLastSessionUpdateTime() {
	c.bxConfig.SetLastSessionUpdateTime()
}

func (c repository) PaginationURLs() []models.PaginationURL {
	return c.bxConfig.PaginationURLs()
}

func (c repository) AddPaginationURL(index int, url string) {
	c.bxConfig.AddPaginationURL(index, url)
}

func (c repository) SetPaginationURLs(paginationURLs []models.PaginationURL) {
	c.bxConfig.SetPaginationURLs(paginationURLs)
}

func (c repository) ClearPaginationURLs() {
	c.bxConfig.ClearPaginationURLs()
}

func NewCoreConfig(errHandler func(error)) ReadWriter {
	return NewCoreConfigFromPath(config_helpers.ConfigFilePath(), errHandler)
}

func NewCoreConfigFromPath(bxConfigPath string, errHandler func(error)) ReadWriter {
	return NewCoreConfigFromPersistor(configuration.NewDiskPersistor(bxConfigPath), errHandler)
}

func NewCoreConfigFromPersistor(bxPersistor configuration.Persistor, errHandler func(error)) ReadWriter {
	return newRepository(createBluemixConfigFromPersistor(bxPersistor, errHandler))
}
