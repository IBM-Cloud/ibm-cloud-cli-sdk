package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/iam"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/uaa"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
)

type pluginContext struct {
	core_config.ReadWriter
	cfConfig     cfConfigWrapper
	pluginConfig PluginConfig
	pluginPath   string
}

type cfConfigWrapper struct {
	core_config.CFConfig
}

func (c cfConfigWrapper) RefreshUAAToken() (string, error) {
	if !c.HasAPIEndpoint() {
		return "", fmt.Errorf("CloudFoundry API endpoint is not set")
	}

	auth := uaa.NewClient(uaa.DefaultConfig(c.AuthenticationEndpoint()), rest.NewClient())
	token, err := auth.GetToken(uaa.RefreshTokenRequest(c.UAARefreshToken()))
	if err != nil {
		return "", err
	}

	ret := fmt.Sprintf("%s %s", token.TokenType, token.AccessToken)
	c.SetUAAToken(ret)
	c.SetUAARefreshToken(token.RefreshToken)
	return ret, nil
}

func createPluginContext(pluginPath string, coreConfig core_config.ReadWriter) *pluginContext {
	return &pluginContext{
		pluginPath:   pluginPath,
		pluginConfig: loadPluginConfigFromPath(filepath.Join(pluginPath, "config.json")),
		ReadWriter:   coreConfig,
		cfConfig:     cfConfigWrapper{coreConfig.CFConfig()},
	}
}

func (c *pluginContext) APIEndpoint() string {
	if compareVersion(c.SDKVersion(), "0.1.1") < 0 {
		return c.ReadWriter.CFConfig().APIEndpoint()
	}
	return c.ReadWriter.APIEndpoint()
}

func (c *pluginContext) IAMEndpoint() string {
	if c.IsPrivateEndpointEnabled() {
		return c.IAMEndpoints().PrivateEndpoint
	}
	return c.IAMEndpoints().PublicEndpoint
}

func (c *pluginContext) ConsoleEndpoint() string {
	if c.IsPrivateEndpointEnabled() {
		return c.ConsoleEndpoints().PrivateEndpoint
	}
	return c.ConsoleEndpoints().PublicEndpoint
}

func compareVersion(v1, v2 string) int {
	s1 := strings.Split(v1, ".")
	s2 := strings.Split(v2, ".")

	n := len(s1)
	if len(s2) > n {
		n = len(s2)
	}

	for i := 0; i < n; i++ {
		var p1, p2 int
		if len(s1) > i {
			p1, _ = strconv.Atoi(s1[i])
		}
		if len(s2) > i {
			p2, _ = strconv.Atoi(s2[i])
		}
		if p1 > p2 {
			return 1
		}
		if p1 < p2 {
			return -1
		}
	}
	return 0
}

func (c *pluginContext) HasAPIEndpoint() bool {
	return c.APIEndpoint() != ""
}

func (c *pluginContext) PluginDirectory() string {
	return c.pluginPath
}

func (c *pluginContext) PluginConfig() PluginConfig {
	return c.pluginConfig
}

func (c *pluginContext) RefreshIAMToken() (string, error) {
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

	ret := fmt.Sprintf("%s %s", token.TokenType, token.AccessToken)
	c.SetIAMToken(ret)
	c.SetIAMRefreshToken(token.RefreshToken)
	return ret, nil
}

func (c *pluginContext) Trace() string {
	return envOrConfig(bluemix.EnvTrace, c.ReadWriter.Trace())
}

func (c *pluginContext) ColorEnabled() string {
	return envOrConfig(bluemix.EnvColor, c.ReadWriter.ColorEnabled())
}

func (c *pluginContext) VersionCheckEnabled() bool {
	return !c.CheckCLIVersionDisabled()
}

func (c *pluginContext) CF() CFContext {
	return c.cfConfig
}

func envOrConfig(env bluemix.Env, config string) string {
	if v := env.Get(); v != "" {
		return v
	}
	return config
}

func (c *pluginContext) CommandNamespace() string {
	return bluemix.EnvPluginNamespace.Get()
}

func (c *pluginContext) CLIName() string {
	return bluemix.EnvCLIName.Get()
}
