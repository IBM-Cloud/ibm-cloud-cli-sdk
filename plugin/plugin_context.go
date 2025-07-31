package plugin

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/endpoints"
)

type pluginContext struct {
	core_config.ReadWriter
	pluginConfig PluginConfig
	pluginPath   string
}

func createPluginContext(pluginPath string, coreConfig core_config.ReadWriter) *pluginContext {
	return &pluginContext{
		pluginPath:   pluginPath,
		pluginConfig: loadPluginConfigFromPath(filepath.Join(pluginPath, "config.json")),
		ReadWriter:   coreConfig,
	}
}

func (c *pluginContext) APIEndpoint() string {
	return c.ReadWriter.APIEndpoint()
}

func (c *pluginContext) ConsoleEndpoint() string {
	if c.IsPrivateEndpointEnabled() {
		if c.IsAccessFromVPC() {
			// return VPC endpoint
			return c.ConsoleEndpoints().PrivateVPCEndpoint
		} else {
			// return CSE endpoint
			return c.ConsoleEndpoints().PrivateEndpoint
		}
	}
	return c.ConsoleEndpoints().PublicEndpoint
}

// GetEndpoint returns the private or public endpoint for a requested service
func (c *pluginContext) GetEndpoint(svc endpoints.Service) (string, error) {
	if c.CloudType() != "public" {
		return "", fmt.Errorf("only public cloud is supported")
	}

	if !c.HasAPIEndpoint() {
		return "", nil
	}

	var cloudDomain string
	switch cname := c.CloudName(); cname {
	case "bluemix":
		cloudDomain = "cloud.ibm.com"
	case "staging":
		cloudDomain = "test.cloud.ibm.com"
	default:
		return "", fmt.Errorf("unknown cloud name '%s'", cname)
	}

	return endpoints.Endpoint(svc, cloudDomain, c.CurrentRegion().Name, c.IsPrivateEndpointEnabled(), c.IsAccessFromVPC())
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

func (c *pluginContext) Trace() string {
	return envOrConfig(bluemix.EnvTrace, c.ReadWriter.Trace())
}

func (c *pluginContext) ColorEnabled() string {
	return envOrConfig(bluemix.EnvColor, c.ReadWriter.ColorEnabled())
}

func (c *pluginContext) VersionCheckEnabled() bool {
	return !c.CheckCLIVersionDisabled()
}

func (c *pluginContext) MCPEnabled() bool {
	return bluemix.EnvMCP.Get() != ""
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
