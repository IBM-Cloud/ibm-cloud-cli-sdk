package i18n

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/resources"
)

const (
	defaultLocale   = "en_US"
	resourcesSuffix = ".all.json"
)

var T i18n.TranslateFunc

func init() {
	loadAsset("i18n/resources/" + defaultLocale + resourcesSuffix)
	T = Tfunc(defaultLocale)
}

func Tfunc(sources ...string) i18n.TranslateFunc {
	defaultTfunc := i18n.MustTfunc(defaultLocale)
	supportedLocales := supportedLocales()

	for _, source := range sources {
		if source == "" {
			continue
		}

		if source == defaultLocale {
			return defaultTfunc
		}

		for _, lang := range language.Parse(source) {
			switch lang.Tag {
			case "zh-cn", "zh-sg":
				lang.Tag = "zh-hans"
			case "zh-hk", "zh-tw":
				lang.Tag = "zh-hant"
			}

			tags := lang.MatchingTags()
			sort.Strings(tags)

			for i := len(tags) - 1; i >= 0; i-- {
				tag := tags[i]

				for locale, assetName := range supportedLocales {
					if strings.HasPrefix(locale, tag) {
						loadAsset(assetName)
						return i18n.MustTfunc(locale)
					}
				}
			}
		}
	}

	return defaultTfunc
}

func loadAsset(assetName string) {
	bytes, err := resources.Asset(assetName)
	if err != nil {
		panic(fmt.Sprintf("Could not load asset '%s': %s", assetName, err.Error()))
	}

	err = i18n.ParseTranslationFileBytes(assetName, bytes)
	if err != nil {
		panic(fmt.Sprintf("Could not load translations '%s': %s", assetName, err.Error()))
	}
}

func supportedLocales() map[string]string {
	m := make(map[string]string)
	for _, assetName := range resources.AssetNames() {
		locale := normalizeLocale(strings.TrimSuffix(path.Base(assetName), ".all.json"))
		m[locale] = assetName
	}
	return m
}

func normalizeLocale(locale string) string {
	return strings.ToLower(strings.Replace(locale, "_", "-", -1))
}
