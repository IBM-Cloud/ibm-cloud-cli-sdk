package terminal_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
)

// Happy path testing
func TestPrintTableSimple(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"test1", "test2"})
	testTable.Add("row1", "row2")
	testTable.Print()
	assert.Contains(t, buf.String(), "test2")
	assert.Contains(t, buf.String(), "row1")
	assert.Equal(t, "test1   test2\nrow1    row2\n", buf.String())
}

func TestPrintTableJson(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"test1", "test2"})
	testTable.Add("row1-col1", "row1-col2")
	testTable.Add("row2-col1", "row2-col2")
	testTable.PrintJson()
	assert.Contains(t, buf.String(), "\"test1\": \"row1-col1\"")
	assert.Contains(t, buf.String(), "\"test2\": \"row2-col2\"")
}

// Blank headers
func TestEmptyHeaderTable(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"", ""})
	testTable.Add("row1", "row2")
	testTable.Print()
	assert.Contains(t, buf.String(), "row1")
	assert.Equal(t, "       \nrow1   row2\n", buf.String())
}

func TestEmptyHeaderTableJson(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"", ""})
	testTable.Add("row1", "row2")
	testTable.PrintJson()
	assert.Contains(t, buf.String(), "\"column_2\": \"row2\"")
	assert.Contains(t, buf.String(), "\"column_1\": \"row1\"")
}

// Empty Headers / More rows than headers
func TestZeroHeadersTable(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{})
	testTable.Add("row1", "row2")
	testTable.Print()
	assert.Contains(t, buf.String(), "row1")
	assert.Equal(t, "\nrow1   row2\n", buf.String())
}

func TestZeroHeadersTableJson(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{})
	testTable.Add("row1", "row2")
	testTable.PrintJson()
	assert.Contains(t, buf.String(), "row1")
	assert.Contains(t, buf.String(), "\"column_2\": \"row2\"")
	assert.Contains(t, buf.String(), "\"column_1\": \"row1\"")
}

// Empty rows / More headers than rows

func TestNotEnoughRowEntires(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"col1", "col2"})
	testTable.Add("row1")
	testTable.Add("", "row2")
	testTable.Print()
	assert.Contains(t, buf.String(), "row1")
	assert.Equal(t, "col1   col2\nrow1   \n       row2\n", buf.String())
}

func TestNotEnoughRowEntiresJson(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{})
	testTable.Add("row1")
	testTable.Add("", "row2")
	testTable.PrintJson()
	assert.Contains(t, buf.String(), "row1")
	assert.Contains(t, buf.String(), "\"column_2\": \"row2\"")
	assert.Contains(t, buf.String(), "\"column_1\": \"row1\"")
	assert.Contains(t, buf.String(), "\"column_1\": \"\"")
}

func TestPrintCsvSimple(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"col1", "col2"})
	testTable.Add("row1-col1", "row1-col2")
	testTable.Add("row2-col1", "row2-col2")
	err := testTable.PrintCsv()
	assert.Equal(t, err, nil)
	assert.Contains(t, buf.String(), "col1,col2")
	assert.Contains(t, buf.String(), "row1-col1,row1-col2")
	assert.Contains(t, buf.String(), "row2-col1,row2-col2")
}

func TestNotEnoughColPrintCsv(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"", "col2"})
	testTable.Add("row1-col1", "row1-col2")
	testTable.Add("row2-col1", "row2-col2")
	err := testTable.PrintCsv()
	assert.Equal(t, err, nil)
	assert.Contains(t, buf.String(), ",col2")
	assert.Contains(t, buf.String(), "row1-col1,row1-col2")
	assert.Contains(t, buf.String(), "row2-col1,row2-col2")
}

func TestNotEnoughRowPrintCsv(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"col1", "col2"})
	testTable.Add("row1-col1", "row1-col2")
	testTable.Add("row2-col1", "")
	err := testTable.PrintCsv()
	assert.Equal(t, err, nil)
	assert.Contains(t, buf.String(), "col1,col2")
	assert.Contains(t, buf.String(), "row1-col1,row1-col2")
	assert.Contains(t, buf.String(), "row2-col1,")
}

func TestEmptyTable(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{})
	err := testTable.PrintCsv()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(strings.TrimSpace(buf.String())), 0)
}

// Text wrapping
func TestTableWrapRows(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"col1"})
	longString := " 2. Reserved Enterprise : Enterprise plan for this offering has been deprecated. Please see the announcement here: https://www-01.ibm.com/common/ssi/ShowDoc.wss?docURL=/common/ssi/rep_ca/3/897/ENUS918-103/index.html&request_locale=en. Analytics Engine provides the ability to spin up and manage Spark clusters. We recommend using this for any production Spark workloads."
	testTable.Add(longString)
	testTable.Print()
	formattedString := "2. Reserved Enterprise : Enterprise plan for this offering has been \ndeprecated. Please see the announcement here: \nhttps://www-01.ibm.com/common/ssi/ShowDoc.wss?docURL=/common/ssi/rep_ca/3/897/ENUS918-103/index.html&request_locale=en.\nAnalytics Engine provides the ability to spin up and manage Spark clusters. We \nrecommend using this for any production Spark workloads."
	assert.Contains(t, buf.String(), formattedString)
}
