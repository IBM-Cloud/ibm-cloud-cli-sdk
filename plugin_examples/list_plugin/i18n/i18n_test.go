package i18n_test

import (
	bxi18n "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin/pluginfakes"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/i18n"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("i18n", func() {
	var (
		pluginContext *pluginfakes.FakePluginContext
		t             bxi18n.TranslateFunc
		locale        string
	)

	BeforeEach(func() {
		pluginContext = new(pluginfakes.FakePluginContext)
	})

	Context("User Config", func() {

		Context("When config is set to zh_Hans", func() {
			BeforeEach(func() {
				locale = "zh_Hans"
				pluginContext.LocaleReturns(locale)
			})

			It("should translate message ID \"Created\" successfully", func() {
				t = i18n.Init(pluginContext)
				Expect(t("Created")).To(Equal("创建"))
			})
		})
	})
})
