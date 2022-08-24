package i18n

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/resources"
)

const (
	defaultLocale   = "en_US"
	resourcesSuffix = ".json"
	resourcesPrefix = "all."
)

var (
	bundle        *i18n.Bundle
	T             TranslateFunc
	RESOURCE_PATH = filepath.Join("i18n", "resources")
)

func init() {
	bundle = Bundle()
	resource := resourcesPrefix + defaultLocale + resourcesSuffix
	loadAsset(filepath.Join(RESOURCE_PATH, resource))
	T = MustTfunc(defaultLocale)
}

// Bundle returns an instance of i18n.bundle
func Bundle() *i18n.Bundle {
	if bundle == nil {
		bundle = i18n.NewBundle(language.AmericanEnglish)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	}
	return bundle
}

// Translate returns a method based on translate method signature in v1.3.0.
// To allow compatability between v1.3.0 and v2.0+, the `messageId` and `args` parameters are
// processed to fit with the new Localize API method.
// @see https://github.com/nicksnyder/go-i18n/blob/v1.3.0/i18n/bundle/bundle.go#L227-L257 for more
// information on the translate method
func Translate(loc *i18n.Localizer) TranslateFunc {
	return func(messageId string, args ...interface{}) string {
		var pluralCount interface{}
		var templateData interface{}

		/**
		 * For the common usecases we can expect two scenarios. Below are two examples:
		 *  1) T("Cannot login in region {{.REGION}}",  map[string]interface{}{"REGION": "us-south"})
		 *  2) T("login_fail_count", "2", map[string]interface{}{"Command": "ibmcloud login"})
		 *
		 * First paramter is always the `messageId`
		 * Second paramter can be either pluralCount or templateData.
		 * Third parameter can be templateData only if the second paramters is the plural count

		 * If we have 2 args than we should expect scenario 2, otherwise we will assume scenario 1
		 */
		if argc := len(args); argc > 0 {
			if isNumber(args[0]) {
				pluralCount = args[0]
				if argc > 1 {
					templateData = args[1]
				}
			} else {
				templateData = args[0]
			}

		}

		msg, err := loc.Localize(&i18n.LocalizeConfig{
			MessageID:    messageId,
			TemplateData: templateData,
			PluralCount:  pluralCount,
		})

		// If no message is returned we can assume that that
		// the translation could not be found in any of the files
		// panic and abort
		if msg == "" {
			panic(err)
		}
		return msg

	}
}

// TranslateFunc returns the translation of the string identified by translationID.
// @see https://github.com/nicksnyder/go-i18n/blob/v1.3.0/i18n/bundle/bundle.go#L19
type TranslateFunc func(translateID string, args ...interface{}) string

func MustTfunc(sources ...string) TranslateFunc {
	defaultLocalizer := i18n.NewLocalizer(bundle, defaultLocale)
	defaultTfunc := Translate(defaultLocalizer)

	supportedLocales, supportedLocalToAsssetMap := supportedLocales()

	for _, source := range sources {
		if source == "" {
			continue
		}

		if source == defaultLocale {
			return defaultTfunc
		}

		lang, _ := language.Parse(source)
		matcher := language.NewMatcher(supportedLocales)
		tag, _ := language.MatchStrings(matcher, lang.String())
		assetName, found := supportedLocalToAsssetMap[tag.String()]

		if found {
			loadAsset(assetName)
			localizer := i18n.NewLocalizer(bundle, source)
			return Translate(localizer)
		}

	}

	return defaultTfunc
}

func loadAsset(assetName string) {
	bytes, assetErr := resources.Asset(assetName)
	if assetErr != nil {
		panic(fmt.Sprintf("Could not load asset '%s': %s", assetName, assetErr.Error()))
	}

	if _, parseErr := bundle.ParseMessageFileBytes(bytes, assetName); parseErr != nil {
		panic(fmt.Sprintf("Could not load translations '%s': %s", assetName, parseErr.Error()))
	}
}

func supportedLocales() ([]language.Tag, map[string]string) {
	// When matching against supported language the first language is set as the fallback
	// so we will initialize the list with English as the first language
	// @see https://pkg.go.dev/golang.org/x/text/language#hdr-Matching_preferred_against_supported_languages for more information
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

func isNumber(n interface{}) bool {
	switch n.(type) {
	case int, int8, int16, int32, int64, string:
		return true
	}
	return false
}
