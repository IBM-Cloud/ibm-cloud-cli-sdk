package terminal_test

import (
	"strings"
	"testing"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/stretchr/testify/assert"
)

func TestFormatWidthWithSpace(t *testing.T) {
	stringOf70 := "0123456789012345678901234567890123456789012345678901234567890123456789"
	stringOf10 := "abcdefghij"

	formattedString := FormatWidth(stringOf70+" "+stringOf10, 0)
	finalStrings := strings.Split(formattedString, "\n")
	assert.Equal(t, 2, len(finalStrings))
	assert.Equal(t, finalStrings[1], stringOf10)

	formattedString = FormatWidth(stringOf10+" "+stringOf70, 0)
	finalStrings = strings.Split(formattedString, "\n")
	assert.Equal(t, 2, len(finalStrings))
	assert.Equal(t, finalStrings[1], stringOf70)

	formattedString = FormatWidth(stringOf10+" "+stringOf70, 5)
	finalStrings = strings.Split(formattedString, "\n")
	assert.Equal(t, 2, len(finalStrings))
	assert.Equal(t, strings.Repeat(" ", 5)+stringOf70, finalStrings[1])

}

func TestFormatWidthLongString(t *testing.T) {
	longString := " 2. Reserved Enterprise : Enterprise plan for this offering has been deprecated. Please see the announcement here: https://www-01.ibm.com/common/ssi/ShowDoc.wss?docURL=/common/ssi/rep_ca/3/897/ENUS918-103/index.html&request_locale=en. Analytics Engine provides the ability to spin up and manage Spark clusters. We recommend using this for any production Spark workloads."

	formattedString := FormatWidth(longString, 4)
	expected := " 2. Reserved Enterprise : Enterprise plan for this offering has been \n    deprecated. Please see the announcement here: \n    https://www-01.ibm.com/common/ssi/ShowDoc.wss?docURL=/common/ssi/rep_ca/3/897/ENUS918-103/index.html&request_locale=en.\n    Analytics Engine provides the ability to spin up and manage Spark \n    clusters. We recommend using this for any production Spark workloads."
	assert.Equal(t, expected, formattedString)
}

func TestFuncWidithWithBreaks(t *testing.T) {
	longString := "this is a string that is so long that it is longer than 80 characters and contains a newline character at the end\n"
	expected := "this is a string that is so long that it is longer than 80 characters and \n    contains a newline character at the end\n"
	assert.Equal(t, expected, FormatWidth(longString, 4))

	anotherLongString := "this is a string that is so long that it is\n longer than 80 characters and contains a newline character in the middle and at the end\n"
	expected = "this is a string that is so long that it is\n longer than 80 characters and contains a newline character in the middle and \n    at the end\n"
	assert.Equal(t, expected, FormatWidth(anotherLongString, 4))

	lastLongString := "this is a string that is so long that it is\n longer than 80 characters and contains\n multiple newline characters in the middle and at the end\n"
	expected = "this is a string that is so long that it is\n longer than 80 characters and contains\n multiple newline characters in the middle and at the end"
	assert.Equal(t, expected, FormatWidth(lastLongString, 4))
}
