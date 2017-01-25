package i18n

import (
	"os"
	"path"
	"sort"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/plugin"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/plugin_examples/list_plugin/resources"
)

const (
	defaultLocale   = "en_US"
	resourcesSuffix = ".all.json"
)

var T i18n.TranslateFunc

func Init(context plugin.PluginContext) i18n.TranslateFunc {
	loadAsset("i18n/resources/" + defaultLocale + resourcesSuffix)
	defaultTfunc := i18n.MustTfunc(defaultLocale)

	supportedLocales := supportedLocales()

	sources := []string{
		context.Locale(),
		os.Getenv("LC_ALL"), // can also use jibber_jabber.DetectIETF()
		os.Getenv("LANG"),   // can also use jibber_jabber.DetectLanguage()
	}

	for _, source := range sources {
		if source == "" {
			continue
		}

		for _, l := range language.Parse(source) {
			switch l.Tag {
			case "zh-cn", "zh-sg":
				l.Tag = "zh-hans"
			case "zh-hk", "zh-tw":
				l.Tag = "zh-hant"
			}

			tags := l.MatchingTags()
			sort.Strings(tags)

			var matchedLocale string
			for i := len(tags) - 1; i >= 0; i-- {
				for l, _ := range supportedLocales {
					if strings.HasPrefix(l, tags[i]) {
						matchedLocale = l
					}
				}
			}

			if matchedLocale != "" {
				loadAsset(supportedLocales[matchedLocale])

				t := i18n.MustTfunc(matchedLocale)
				return func(translationID string, args ...interface{}) string {
					if translated := t(translationID, args...); translated != translationID {
						return translated
					}

					return defaultTfunc(translationID, args...)
				}
			}
		}
	}

	return defaultTfunc
}

func supportedLocales() map[string]string {
	m := make(map[string]string)
	for _, assetName := range resources.AssetNames() {
		locale := normalizeLocale(strings.TrimSuffix(path.Base(assetName), resourcesSuffix))
		m[locale] = assetName
	}
	return m
}

func normalizeLocale(locale string) string {
	return strings.ToLower(strings.Replace(locale, "_", "-", -1))
}

func loadAsset(assetName string) {
	bytes, err := resources.Asset(assetName)
	if err == nil {
		err = i18n.ParseTranslationFileBytes(assetName, bytes)
	}
	if err != nil {
		panic(err)
	}
}
