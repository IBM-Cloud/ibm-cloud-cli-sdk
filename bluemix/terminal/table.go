package terminal

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"io"
	"strings"
)

type Table interface {
	Add(row ...string)
	Print()
}

type PrintableTable struct {
	writer             io.Writer
	headers            []string
	headerPrinted      bool
	maxSizes           []int
	rows               [][]string //each row is single line
	collapsibleColumns []string
}

type TableOpt func(Table)

func WithCollapsibleColumns(headers ...string) TableOpt {
	return func(t Table) {
		pt, ok := t.(*PrintableTable)
		if !ok {
			return
		}
		pt.collapsibleColumns = append([]string{}, headers...)
	}
}

func NewTable(w io.Writer, headers []string, opts ...TableOpt) Table {
	table := &PrintableTable{
		writer:   w,
		headers:  headers,
		maxSizes: make([]int, len(headers)),
	}
	for _, opt := range opts {
		opt(table)
	}
	return table
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
}

func (t *PrintableTable) Print() {
	if len(t.collapsibleColumns) > 0 {
		t.collapse()
	}

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
	if col < len(t.headers)-1 {
		padding = strings.Repeat(" ", t.maxSizes[col]-runewidth.StringWidth(Decolorize(value)))
	}
	return fmt.Sprintf("%s%s   ", value, padding)
}

func (t *PrintableTable) collapse() {
	jend := len(t.headers)
	for j := 0; j < jend; j++ {
		if !t.isCollapsible(t.headers[j]) {
			continue
		}

		// "shift" left to delete jth column while preserving order
		t.headers = append(t.headers[:j], t.headers[j+1:jend]...)
		for i := range t.rows {
			t.rows[i] = append(t.rows[i][:j],t.rows[i][j+1:jend]...)
		}

		j--
		jend--
	}
}

func (t *PrintableTable) isCollapsible(header string) bool {
	var collapsible bool
	for _, cc := range t.collapsibleColumns {
		if cc == header {
			collapsible = true
			break
		}
	}
	if !collapsible {
		return false
	}

	for j, col := range t.headers {
		if col != header {
			continue
		}

		for _, row := range t.rows {
			if row[j] != "" {
				return false
			}
		}
		return true
	}
	return false
}
