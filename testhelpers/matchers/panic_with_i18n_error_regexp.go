package matchers

import (
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/onsi/gomega"
)

type TranslateErrorMatcher struct {
	expected string
}

// PanicWithi18nErrorRegexp matches when a panic is due to message id not found in the translation files
// If match is not found this method will panic.
// @see https://github.com/onsi/gomega/issues/471 for more information on reporting errors using WithTransform
func PanicWithi18nErrorRegexp(message string) gomega.OmegaMatcher {
	return gomega.PanicWith(gomega.WithTransform(func(actual error) string {
		var e *i18n.MessageNotFoundErr
		if errors.As(actual, &e) {
			return actual.Error()
		}
		return ""
	}, gomega.MatchRegexp(message)))
}
