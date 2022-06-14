package i18n_test

import (
	"os"

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

	Context("Default Config (English)", func() {
		BeforeEach(func() {
			pluginContext.LocaleReturns("")
		})
		It("should translate message ID \"Created\" successfully in English", func() {
			t = i18n.Init(pluginContext)
			Expect(t("Created")).To(Equal("Created"))
		})
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

	Context("Environment variables", func() {
		var origEnv string
		Context("LC_ALL is set", func() {
			BeforeEach(func() {
				pluginContext.LocaleReturns("")
				origEnv = os.Getenv("LC_ALL")
			})

			AfterEach(func() {
				os.Setenv("LC_ALL", origEnv)
			})

			Context("When config is set to zh_Hans", func() {
				BeforeEach(func() {
					locale = "zh_Hans"
					pluginContext.LocaleReturns("")
					os.Setenv("LC_ALL", locale)
				})

				It("should translate message ID \"Created\" successfully", func() {
					t = i18n.Init(pluginContext)
					Expect(t("Created")).To(Equal("创建"))
				})
			})
		})

		Context("LANG is set", func() {
			BeforeEach(func() {
				pluginContext.LocaleReturns("")
				origEnv = os.Getenv("LANG")
			})

			AfterEach(func() {
				os.Setenv("LANG", origEnv)
			})

			Context("When config is set to zh_Hans", func() {
				BeforeEach(func() {
					locale = "zh_Hans"
					pluginContext.LocaleReturns("")
					os.Setenv("LANG", locale)
				})

				It("should translate message ID \"Created\" successfully", func() {
					t = i18n.Init(pluginContext)
					Expect(t("Created")).To(Equal("创建"))
				})
			})
		})
	})

	Context("Missing translation ID in Chinese language", func() {
		var (
			translationID string
			translatedMsg string
		)
		BeforeEach(func() {
			locale = "zh_Hans"
			pluginContext.LocaleReturns(locale)
		})

		Context("Translation exist in default language (English)", func() {
			BeforeEach(func() {
				translationID = "Test translation only in English"
				translatedMsg = translationID
			})

			It("should translate message in english", func() {
				t = i18n.Init(pluginContext)
				Expect(t(translationID)).To(Equal(translatedMsg))
			})
		})
	})
})
