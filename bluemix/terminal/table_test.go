package terminal_test

import (
	"bytes"
	"os"
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
	assert.Equal(t, "\nrow1   row2\n", buf.String())
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
	assert.Equal(t, "col1   col2\nrow1\n       row2\n", buf.String())
}

func TestMoreColThanTerminalWidth(t *testing.T) {
	os.Setenv("TEST_TERMINAL_WIDTH", "1")
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"col1"})
	testTable.Add("row1", "row2")
	testTable.Print()
	assert.Contains(t, buf.String(), "row1")
	assert.Equal(t, "col1\nrow1   row2\n", buf.String())
	os.Unsetenv("TEST_TERMINAL_WIDTH")
}

func TestWideHeaderNames(t *testing.T) {
	buf := bytes.Buffer{}
	testTable := NewTable(&buf, []string{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt u", "NAME"})
	testTable.Add("col1", "col2")
	testTable.Print()
	assert.Contains(t, buf.String(), "Lorem ipsum dolor sit amet, consectetu")
	assert.Equal(t, "Lorem ipsum dolor sit amet, consectetu   NAME\nr adipiscing elit, sed do eiusmod temp\nor incididunt u\ncol1                                     col2\n", buf.String())
}

func TestWidestColumn(t *testing.T) {
	buf := bytes.Buffer{}
	id := "ABCDEFG-9b8babbd-f2ed-4371-b817-a839e4130332"
	testTable := NewTable(&buf, []string{"ID", "Name"})
	testTable.Add(id, "row2")
	testTable.Print()
	assert.Contains(t, buf.String(), id)
	assert.Equal(t, buf.String(), "ID                                             Name\nABCDEFG-9b8babbd-f2ed-4371-b817-a839e4130332   row2\n")
}

func TestMultiWideColumns(t *testing.T) {
	buf := bytes.Buffer{}
	id := "ABCDEFG-9b8babbd-f2ed-4371-b817-a839e4130332"
	desc := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut"
	testTable := NewTable(&buf, []string{"ID", "Description", "Name"})
	testTable.Add(id, desc, "col3")
	testTable.Print()
	assert.Contains(t, buf.String(), "ABCDEFG-9b8babbd-f2ed-4371-b817-a839")
	assert.Contains(t, buf.String(), "e4130332")
	assert.Equal(t, buf.String(), "ID                                     Description                           Name\nABCDEFG-9b8babbd-f2ed-4371-b817-a839   Lorem ipsum dolor sit amet, consect   col3\ne4130332                               etur adipiscing elit, sed do eiusmo\n                                       d tempor incididunt ut\n")
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
