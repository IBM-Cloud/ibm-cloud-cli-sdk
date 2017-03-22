package matchers

import (
	"fmt"
	"strings"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/terminal"
	"github.com/onsi/gomega"
)

type SliceMatcher struct {
	expected      [][]string
	failedAtIndex int
}

func ContainSubstrings(substrings ...[]string) gomega.OmegaMatcher {
	return &SliceMatcher{expected: substrings}
}

func (matcher *SliceMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("ContainSubstrings matcher expects a string, but it's actually a %T", actual)
	}

	actualStrings := strings.Split(actualString, "\n")
	allStringsMatched := make([]bool, len(matcher.expected))

	for index, expectedArray := range matcher.expected {
		for _, actualValue := range actualStrings {
			actualValue = terminal.Decolorize(actualValue)

			allStringsFound := true

			for _, expectedValue := range expectedArray {
				if !strings.Contains(actualValue, expectedValue) {
					allStringsFound = false
				}
			}

			if allStringsFound {
				allStringsMatched[index] = true
				break
			}
		}
	}

	for index, value := range allStringsMatched {
		if !value {
			matcher.failedAtIndex = index
			return false, nil
		}
	}

	return true, nil
}

func (matcher *SliceMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("expected to find \"%s\" in actual:\n\"%v\"\n", matcher.expected[matcher.failedAtIndex], actual)
}

func (matcher *SliceMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("expected to not find \"%s\" in actual:\n\"%v\"\n", matcher.expected[matcher.failedAtIndex], actual)
}
