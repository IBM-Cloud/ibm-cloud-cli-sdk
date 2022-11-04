package terminal

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
)

const (
	// minSpace is the number of spaces between the end of the longest value and
	// the start of the next row
	minSpace = 3
)

type Table interface {
	Add(row ...string)
	Print()
	PrintJson()
	PrintCsv()
}

type PrintableTable struct {
	writer        io.Writer
	headers       []string
	headerPrinted bool
	maxSizes      []int
	rows          [][]string //each row is single line
}

func NewTable(w io.Writer, headers []string) Table {
	return &PrintableTable{
		writer:   w,
		headers:  headers,
		maxSizes: make([]int, len(headers)),
	}
}

func (t *PrintableTable) Add(row ...string) {
	var maxLines int

	var columns [][]string
	for _, value := range row {
		lines := strings.Split(value, "\n")
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
		columns = append(columns, lines)
	}

	for i := 0; i < maxLines; i++ {
		var row []string
		for _, col := range columns {
			if i >= len(col) {
				row = append(row, "")
			} else {
				row = append(row, col[i])
			}
		}
		t.rows = append(t.rows, row)
	}

	// Incase we have more columns in a row than headers, need to update maxSizes
	if len(row) > len(t.maxSizes) {
		t.maxSizes = make([]int, len(row))
	}
}

func (t *PrintableTable) Print() {
	for _, row := range append(t.rows, t.headers) {
		t.calculateMaxSize(row)
	}

	if t.headerPrinted == false {
		t.printHeader()
		t.headerPrinted = true
	}

	for _, line := range t.rows {
		t.printRow(line)
	}

	t.rows = [][]string{}
}

func (t *PrintableTable) calculateMaxSize(row []string) {
	for index, value := range row {
		cellLength := runewidth.StringWidth(Decolorize(value))
		if t.maxSizes[index] < cellLength {
			t.maxSizes[index] = cellLength
		}
	}
}

func (t *PrintableTable) printHeader() {
	output := ""
	for col, value := range t.headers {
		output = output + t.cellValue(col, HeaderColor(value))
	}
	fmt.Fprintln(t.writer, output)
}

func (t *PrintableTable) printRow(row []string) {
	output := ""
	for columnIndex, value := range row {
		if columnIndex == 0 {
			value = TableContentHeaderColor(value)
		}

		output = output + t.cellValue(columnIndex, value)
	}
	fmt.Fprintln(t.writer, output)
}

func (t *PrintableTable) cellValue(col int, value string) string {
	padding := ""
	if col < len(t.maxSizes)-1 {
		padding = strings.Repeat(" ", t.maxSizes[col]-runewidth.StringWidth(Decolorize(value))+minSpace)
	}
	return fmt.Sprintf("%s%s", value, padding)
}

// Prints out a nicely/human formatted Json string instead of a table structure
func (t *PrintableTable) PrintJson() {
	total_col := len(t.headers) - 1
	total_row := len(t.rows) - 1
	fmt.Fprintln(t.writer, "[")
	// Iterate through the rows
	for i, row := range t.rows {
		fmt.Fprintln(t.writer, "\t{")
		// Iterate through the columns in a specific row
		for x, point := range row {
			cur_col := ""
			// Some rows might have more columns than headers
			// or empty headers
			if x > total_col || t.headers[x] == "" {
				cur_col = fmt.Sprintf("column_%d", (x + 1))
			} else {
				cur_col = t.headers[x]
			}
			entry := fmt.Sprintf("\t\t\"%s\": \"%s\"", cur_col, point)
			// emit a "," unless were at the last element
			if x != (len(row) - 1) {
				fmt.Fprintln(t.writer, fmt.Sprintf("%s,", entry))
			} else {
				fmt.Fprintln(t.writer, fmt.Sprintf("%s", entry))
			}
		}

		if i != total_row {
			fmt.Fprintln(t.writer, "\t},")
		} else {
			fmt.Fprintln(t.writer, "\t}")
		}
	}
	fmt.Fprintln(t.writer, "]")
	// mimic behavior of Print()
	t.rows = [][]string{}
}

func (t *PrintableTable) PrintCsv() {
	csvwriter := csv.NewWriter(t.writer)
	csvwriter.Write(t.headers)
	csvwriter.WriteAll(t.rows)
}
