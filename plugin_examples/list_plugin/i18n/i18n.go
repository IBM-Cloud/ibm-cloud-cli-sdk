package i18n

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/resources"
)

const (
	defaultLocale   = "en_US"
	resourcePrefix  = "all."
	resourcesSuffix = ".json"
)

var (
	bundle        *goi18n.Bundle
	T             i18n.TranslateFunc
	RESOURCE_PATH = filepath.Join("i18n", "resources")
)

func Init(context plugin.PluginContext) i18n.TranslateFunc {
	bundle = i18n.Bundle()
	resource := resourcePrefix + defaultLocale + resourcesSuffix
	loadAsset(filepath.Join(RESOURCE_PATH, resource))
	defaultLocalizer := goi18n.NewLocalizer(bundle, defaultLocale)
	defaultTfunc := i18n.Translate(defaultLocalizer)

	supportedLocales, supportedLocalToAsssetMap := supportedLocales()

	sources := []string{
		context.Locale(),
		os.Getenv("LC_ALL"), // can also use jibber_jabber.DetectIETF()
		os.Getenv("LANG"),   // can also use jibber_jabber.DetectLanguage()
	}

	for _, source := range sources {
		if source == "" {
			continue
		}

		// Handle chinese language variants
		// (eg. Chinese (Simplified, Singapore), Chinese (Traditional, Hong Kong S.A.R.)
		switch source {
		case "zh-cn", "zh-sg":
			source = "zh-hans"
		case "zh-hk", "zh-tw":
			source = "zh-hant"
		}

		lang, _ := language.Parse(source)
		matcher := language.NewMatcher(supportedLocales)
		tag, _ := language.MatchStrings(matcher, lang.String())
		assetName, found := supportedLocalToAsssetMap[tag.String()]

		if found {
			loadAsset(assetName)
			localizer := goi18n.NewLocalizer(bundle, source)

			t := i18n.Translate(localizer)
			return func(translationID string, args ...interface{}) string {
				if translated := t(translationID, args...); translated != translationID {
					return translated
				}

				return defaultTfunc(translationID, args...)
			}
		}
	}

	return defaultTfunc
}

func supportedLocales() ([]language.Tag, map[string]string) {
	l := []language.Tag{language.English}
	m := make(map[string]string)
	for _, assetName := range resources.AssetNames() {
		// Remove the "all." prefix and ".json" suffix to get language/locale
		locale := normalizeLocale(strings.TrimSuffix(path.Base(assetName), ".json"))
		locale = strings.TrimPrefix(locale, "all.")

		if !strings.Contains(locale, normalizeLocale(defaultLocale)) {
			lang, _ := language.Parse(locale)
			l = append(l, lang)
			m[lang.String()] = assetName
		}
	}
	return l, m
}

func normalizeLocale(locale string) string {
	return strings.ToLower(strings.Replace(locale, "_", "-", -1))
}

func loadAsset(assetName string) {
	bytes, err := resources.Asset(assetName)
	if err == nil {
		_, err = bundle.ParseMessageFileBytes(bytes, assetName)
	}
	if err != nil {
		panic(err)
	}
}
