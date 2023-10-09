package motd

import (
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/api"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
)

type MODMessages struct {
	VersionRange string `json:"versionRange"`
	Region       string `json:"region"`
	OS           string `json:"os"`
	Message      string `json:"message"`
}

type MODResponse struct {
	Messages []MODMessages `json:"messages"`
}

func CheckMessageOftheDayForPlugin(pluginConfig plugin.PluginConfig) bool {
	currentMessagOfTheDay := pluginConfig.Get("MessageOfTheDay")
	if currentMessagOfTheDay == nil {
		return false
	}
	if currentMessagOfTheDayTimestamp, ok := currentMessagOfTheDay.(string); ok {
		lastCheckTime, parseErr := strconv.ParseInt(currentMessagOfTheDayTimestamp, 10, 64)
		if parseErr != nil {
			return false
		}
		return time.Since(time.Unix(lastCheckTime, 0)).Hours() < 24
	}
	return false
}

func DisplayMessageOfTheDay(client *rest.Client, config core_config.ReadWriter, pluginConfig plugin.PluginConfig, modURL string, ui terminal.UI, version string) {
	// the pluginConfig variable will be cast-able to a pointer to a plugin config type if the display is for a plugin

	if config != nil {
		if !config.CheckMessageOfTheDay() {
			return
		}
		defer config.SetMessageOfTheDayTime()
	} else {
		if !CheckMessageOftheDayForPlugin(pluginConfig) {
			return
		}
		defer pluginConfig.Set("MessageOfTheDayTime", time.Now().Unix())
	}

	var mod MODResponse

	_, err := client.Do(rest.GetRequest(modURL), &mod, nil)
	if err != nil {
		return
	}

	for _, mes := range mod.Messages {

		if mes.VersionRange != "" {
			constraint, err := api.NewSemverConstraint(mes.VersionRange)
			if err != nil || !constraint.Satisfied(bluemix.Version.String()) {
				continue
			}
		}

		// Only print message if targeted region matches or if message region is empty(print for all regions)
		if mes.Region != "" {
			skip := true
			for _, region := range strings.Split(mes.Region, ",") {
				region = strings.TrimSpace(region)
				if strings.EqualFold(region, config.CurrentRegion().Name) {
					skip = false
					break
				}
			}
			if skip {
				continue
			}
		}

		// Only print message if matching OS or if message's OS is empty(print for all platforms)
		if mes.OS != "" {
			skip := true
			for _, platform := range strings.Split(mes.OS, ",") {
				platform = strings.TrimSpace(platform)
				if strings.EqualFold(platform, core_config.DeterminePlatform()) {
					skip = false
					break
				}
			}
			if skip {
				continue
			}
		}

		ui.Warn(mes.Message)
	}
}
