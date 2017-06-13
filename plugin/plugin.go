package plugin

import (
	"fmt"
	"strings"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/plugin/models"
)

type PluginMetadata struct {
	Name          string      // name of the plugin
	Version       VersionType // version of the plugin
	MinCliVersion VersionType // minimal CLI version required by the plugin
	Namespaces    []Namespace // command namespaces defined for the plugin
	Commands      []Command   // list of commands provided by the plugin
}

type VersionType struct {
	Major int
	Minor int
	Build int
}

func (v VersionType) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Build)
}

type Namespace struct {
	Name        string // name of the namespace
	Description string // description of the namespace
}

func (n Namespace) ParentName() string {
	i := strings.LastIndex(n.Name, " ")
	if i < 0 {
		return ""
	}
	return n.Name[:i]
}

type Command struct {
	Namespace   string // namespace of the command
	Name        string // command name
	Alias       string // command alias, usually the command's short name
	Description string // short description of the command
	Usage       string // usage detail to be displayed in command help
	Flags       []Flag // command options
}

func (c Command) FullName() string {
	return strings.TrimSpace(strings.Join([]string{c.Namespace, c.Name}, " "))
}

func (c Command) FullNames() []string {
	names := []string{c.FullName()}
	if c.Alias != "" {
		names = append(names, strings.TrimSpace(strings.Join([]string{c.Namespace, c.Alias}, " ")))
	}
	return names
}

// Command option
type Flag struct {
	Name        string // name of the option
	Description string // description of the option
	HasValue    bool   // whether the option requires a value or not
}

// Plugin is the interface of Bluemix CLI plugin.
type Plugin interface {
	// GetMetadata returns the metadata of the plugin.
	GetMetadata() PluginMetadata

	// Run runs the plugin command with plugin context and given arguments.
	// Note: the first argument is always the command name or alias no matter
	// the command has namespace or not.
	// To get command namespace, call PluginContext.CommandNamespace()
	Run(c PluginContext, args []string)
}

// PluginContext holds context to be passed into plugin's Run method.
type PluginContext interface {
	APIVersion() string
	APIEndpoint() string
	HasAPIEndpoint() bool
	// deprecate loggergator endpoint, use Doppler endpoint instead
	// LoggregatorEndpoint() string
	DopplerEndpoint() string
	ConsoleEndpoint() string
	UAAEndpoint() string
	UAAToken() string
	UAARefreshToken() string
	RefreshUAAToken() (string, error)
	IAMTokenEndpoint() string
	IAMToken() string
	IAMRefreshToken() string
	RefreshIAMToken() (string, error)
	Username() string
	UserGUID() string
	UserEmail() string
	IsLoggedIn() bool
	CurrentOrg() models.Organization
	HasOrganization() bool
	AccountID() string
	IMSAccountID() string
	CurrentSpace() models.Space
	HasSpace() bool
	Region() string   // region name
	RegionID() string // region unique ID
	Locale() string
	Trace() string
	ColorEnabled() string
	IsSSLDisabled() bool
	PluginDirectory() string
	HTTPTimeout() int
	VersionCheckEnabled() bool
	PluginConfig() PluginConfig
	CommandNamespace() string
}
