package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin/pluginfakes"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/i18n"
)

func TestCommands(t *testing.T) {
	i18n.T = i18n.Init(new(pluginfakes.FakePluginContext))

	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}
