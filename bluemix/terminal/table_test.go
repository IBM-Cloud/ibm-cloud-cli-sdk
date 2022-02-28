package terminal_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"

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
	assert.Equal(t, buf.String(), "test1   test2   \nrow1    row2   \n")
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
	assert.Equal(t, buf.String(), "          \nrow1   row2   \n")
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
	assert.Equal(t, buf.String(), "\nrow1   row2   \n")
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
	assert.Equal(t, buf.String(), "col1   col2   \nrow1   \n       row2   \n")
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
