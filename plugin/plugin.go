// Package plugin provides types and functions common among plugins.
//
// See examples in "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples".
package plugin

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/endpoints"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// PluginMetadata describes metadata of a plugin.
type PluginMetadata struct {
	Name          string      // name of the plugin
	Aliases       []string    // aliases of the plugin
	Version       VersionType // version of the plugin
	MinCliVersion VersionType // minimal version of CLI required by the plugin
	Namespaces    []Namespace // list of namespaces provided by the plugin
	Commands      []Command   // list of commands provided by the plugin

	// SDKVersion is SDK version used by the plugin.
	// It is set by the plugin framework to check SDK compatibility with the CLI.
	SDKVersion VersionType

	// If DelegateBashCompletion is true, plugin command's completion is handled by plugin.
	// The CLI will invoke '<plugin_binary> SendCompletion <args>'
	DelegateBashCompletion bool

	// Whether the plugin supports private endpoint
	PrivateEndpointSupported bool

	// Whether the plugin supports private endpoint access via VPC
	IsAccessFromVPC bool

	// Whether the plugin was built using Cobra
	IsCobraPlugin bool
}

func (p PluginMetadata) NameAndAliases() []string {
	return append([]string{p.Name}, p.Aliases...)
}

// VersionType describes the version info
type VersionType struct {
	Major int // major version
	Minor int // minimal version
	Build int // build number
}

func (v VersionType) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Build)
}

// Stage decribes the stage of the namespace or command and will be shown in the help text
type Stage string

// Valid stages
const (
	StageExperimental Stage = "experimental"
	StageBeta         Stage = "beta"
	StageDeprecated   Stage = "deprecated"
)

// Namespace represents a category of commands that have similar
// functionalities. A command under a namespace is run using 'bx [namespace]
// [command]'.
//
// Some namespaces are predefined by Bluemix CLI and sharable among plugins;
// others are non-shared and defined in each plugin. The Plugin can reference a
// predefined namespace or define a non-shared namespace of its own.
//
// Namespace also supports hierarchy. For example, namespace 'A' can have a sub
// namespace 'B'. And the full qualified name of namespace B is 'A B'.
type Namespace struct {
	ParentName  string   // full qualified name of the parent namespace
	Name        string   // base name
	Aliases     []string // aliases
	Description string   // description of the namespace
	Stage       Stage    // stage of the commands in the namespace
}

func (n Namespace) NameAndAliases() []string {
	return append([]string{n.Name}, n.Aliases...)
}

// Command describes the metadata of a plugin command
type Command struct {
	Namespace   string   // full qualified name of the command's namespace
	Name        string   // name of the command
	Alias       string   // Deprecated: use Aliases instead.
	Aliases     []string // command aliases
	Description string   // short description of the command
	Usage       string   // usage detail to be displayed in command help
	Flags       []Flag   // command options
	Hidden      bool     // true to hide the command in help text
	Stage       Stage    // stage of the command
}

func (c Command) NameAndAliases() []string {
	as := c.Aliases
	if len(as) == 0 && c.Alias != "" {
		as = []string{c.Alias}
	}
	return append([]string{c.Name}, as...)
}

// ConvertCobraFlagsToPluginFlags will convert flags defined by Cobra framework to Plugin Flags
// Method is used when defining the Flags in command metadata. @see Plugin#GetMetadata() for use case
func ConvertCobraFlagsToPluginFlags(cmd *cobra.Command) []Flag {
	var flags []Flag
	// NOTE: there is a strange behavior in Cobra where you need to call
	// either `InheritedFlags` or `LocalFlags` in order to include global
	// flags when calling `VisitAll`
	// see https://github.com/spf13/cobra/issues/412
	cmd.InheritedFlags()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		var name string
		if f.Shorthand != "" {
			name = f.Shorthand + "," + f.Name
		} else {
			name = f.Name
		}
		hasValue := true
		if f.Value.Type() == "bool" {
			hasValue = false
		}
		flags = append(flags, Flag{
			Name:        name,
			Description: f.Usage,
			HasValue:    hasValue,
			Hidden:      f.Hidden,
		})
	})

	return flags

}

// Flag describes a command option
type Flag struct {
	Name        string // name of the option
	Description string // description of the option
	HasValue    bool   // whether the option requires a value or not
	Hidden      bool   // true to hide the option in command help
}

// Plugin is an interface for Bluemix CLI plugins.
//
//go:generate counterfeiter . Plugin
type Plugin interface {
	// GetMetadata returns the metadata of the plugin.
	GetMetadata() PluginMetadata

	// Run runs the plugin with given plugin context and arguments.
	//
	// Note: the first arg in args is a command or alias no matter the command
	// has namespace or not. To get the namespace, call
	// PluginContext.CommandNamespace.
	Run(c PluginContext, args []string)
}

// PluginContext is a Bluemix CLI context passed to plugin's Run method. It
// carries service endpoints info, login session, user configuration, plugin
// configuration and provides utility methods.
//
//go:generate counterfeiter . PluginContext
type PluginContext interface {
	// APIEndpoint returns the targeted API endpoint of IBM Cloud
	APIEndpoint() string

	// HasAPIEndpoint() returns whether an IBM Cloud has been targeted
	HasAPIEndpoint() bool

	// IsPrivateEndpointEnabled returns whether use of the private endpoint has been chosen
	IsPrivateEndpointEnabled() bool

	// IsAccessFromVPC returns true if the access from VPC private endpoint is enabled.
	IsAccessFromVPC() bool

	// ConsoleEndpoint returns console's public endpoint if api endpoint is public, or returns
	// private endpoint if api endpoint is private.
	ConsoleEndpoint() string

	// ConsoleEndpoints returns both the public and private endpoint of console.
	ConsoleEndpoints() models.Endpoints

	// IAMEndpoint returns IAM's public endpoint if api endpoint is public, or returns private
	// endpoint if api endpoint is private.
	IAMEndpoint() string

	// IAMEndpoints returns both the public and private endpoint of IAM.
	IAMEndpoints() models.Endpoints

	// GetEndpoint is a utility method to return private or public endpoint for a requested service.
	// It supports public cloud only. For non public clouds, plugin needs its own way to determine endpoint.
	GetEndpoint(endpoints.Service) (string, error)

	// CloudName returns the name of the target cloud
	CloudName() string

	// CloudType returns the type of the target cloud (like 'public',
	// 'dedicated' etc)
	CloudType() string

	// Region returns the targeted region
	CurrentRegion() models.Region

	// CRIType returns the type of compute resource the user logged in as, if applicable. Valid values are `IKS`, `VPC`, or `OTHER`
	CRIType() string

	// VPCCRITokenURL() returns the value specified by the environment variable 'IBMCLOUD_CR_VPC_URL', if set.
	// Otherwise, the default VPC auth url specified by the constant `DefaultServerEndpoint` is returned
	VPCCRITokenURL() string

	// HasTargetedRegion() return whether a region is targeted
	HasTargetedRegion() bool

	// IAMToken returns the IAM access token
	IAMToken() string

	// IAMRefreshToken returns the IAM refresh token
	IAMRefreshToken() string

	// RefreshIAMToken refreshes and returns the IAM access token
	RefreshIAMToken() (string, error)

	// UserEmail returns the Email of the logged in user
	UserEmail() string

	// IsLoggedIn returns if a user has logged into IBM Cloud
	IsLoggedIn() bool

	// IsLoggedInWithServiceID returns if a user has logged into IBM Cloud using service ID.
	IsLoggedInWithServiceID() bool

	// IsLoggedInAsProfile returns true if a user logged into IBM Cloud using an IAM token pertaining to a trusted profile
	IsLoggedInAsProfile() bool

	// IsLoggedInAsCRI returns true if a user logged into IBM Cloud as a compute resource.
	IsLoggedInAsCRI() bool

	// IMSAccountID returns ID of the IMS account linked to the targeted BSS
	// account
	IMSAccountID() string

	// Account returns the targeted a BSS account
	CurrentAccount() models.Account

	// HasTargetedAccount returns whether an account has been targeted
	HasTargetedAccount() bool

	// HasTargetedProfile returns whether a profile has been targeted
	HasTargetedProfile() bool

	// CurrentProfile returns the targeted trusted profile
	CurrentProfile() models.Profile

	// HasTargetedComputeResource returns whether a compute resource has been targeted
	HasTargetedComputeResource() bool

	// ResourceGroup returns the targeted resource group
	CurrentResourceGroup() models.ResourceGroup

	// HasTargetedResourceGroup returns whether a resource group has been targeted
	HasTargetedResourceGroup() bool

	// Locale returns user specified locale
	Locale() string

	// Trace returns user specified trace setting.
	// The value is "true", "false" or path of the trace output file.
	Trace() string

	// ColorEnabled returns whether terminal displays color or not
	ColorEnabled() string

	// IsSSLDisabled returns whether skipping SSL validation or not
	IsSSLDisabled() bool

	// PluginDirectory returns the installation directory of the plugin
	PluginDirectory() string

	// HTTPTimeout returns a timeout for HTTP Client
	HTTPTimeout() int

	// VersionCheckEnabled() returns whether checking for update is performmed
	VersionCheckEnabled() bool

	// PluginConfig returns the plugin specific configuarion
	PluginConfig() PluginConfig

	// CommandNamespace returns the name of the parsed namespace
	CommandNamespace() string

	// CLIName returns binary name of the Bluemix CLI that is invoking the plugin
	CLIName() string

	// MCPEnabled returns true if the CLI is functioning as an MCP server
	MCPEnabled() bool
}
