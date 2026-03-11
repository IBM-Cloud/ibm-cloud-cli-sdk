package terminal

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"

	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
	PrintCsv() error
}

type PrintableTable struct {
	writer   io.Writer
	headers  []string
	maxSizes []int
	rows     [][]string //each row is single line
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

func isWideColumn(col string) bool {
	// list of common columns that are usually wide
	largeColumnTypes := []string{T("ID"), T("Description")}

	for _, largeColn := range largeColumnTypes {
		if strings.Contains(largeColn, col) {
			return true
		}
	}

	return false

}

func terminalWidth() int {
	var err error
	terminalWidth, _, err := term.GetSize(int(os.Stdin.Fd()))

	if err != nil {
		// Assume normal 80 char width line
		terminalWidth = 80
	}

	testTerminalWidth, envSet := os.LookupEnv("TEST_TERMINAL_WIDTH")
	if envSet {
		envWidth, err := strconv.Atoi(testTerminalWidth)
		if err == nil {
			terminalWidth = envWidth
		}
	}
	return terminalWidth
}

func (t *PrintableTable) Print() {
	for _, row := range append(t.rows, t.headers) {
		t.calculateMaxSize(row)
	}

	tbl := table.NewWriter()
	tbl.SetOutputMirror(t.writer)
	tbl.SuppressTrailingSpaces()
	// remove padding from the left to keep the table aligned to the left
	tbl.Style().Box.PaddingLeft = ""
	tbl.Style().Box.PaddingRight = strings.Repeat(" ", minSpace)
	// remove all border and column and row separators
	tbl.Style().Options.DrawBorder = false
	tbl.Style().Options.SeparateColumns = false
	tbl.Style().Options.SeparateFooter = false
	tbl.Style().Options.SeparateHeader = false
	tbl.Style().Options.SeparateRows = false
	tbl.Style().Format.Header = text.FormatDefault
	tbl.Style().Color.Header = text.Colors{text.Bold}

	headerRow, rows := t.createPrettyRowsAndHeaders()
	columnConfig := t.createColumnConfigs()

	tbl.SetColumnConfigs(columnConfig)
	tbl.AppendHeader(headerRow)
	tbl.AppendRows(rows)
	tbl.Render()
}

func (t *PrintableTable) createColumnConfigs() []table.ColumnConfig {
	// there must be at row in order to configure column
	if len(t.rows) == 0 {
		return []table.ColumnConfig{}
	}

	colCount := len(t.rows[0])
	var (
		widestColIndicies []int
		terminalWidth     = 1000
		// total amount padding space that a row will take up
		totalPaddingSpace = (colCount - 1) * minSpace
		remainingSpace    = max(0, terminalWidth-totalPaddingSpace)
		// the estimated max column width by dividing the remaining space evenly across the columns
		maxColWidth = remainingSpace / colCount
	)
	columnConfig := make([]table.ColumnConfig, colCount)

	for colIndex := range columnConfig {
		columnConfig[colIndex] = table.ColumnConfig{
			AlignHeader: text.AlignLeft,
			Align:       text.AlignLeft,
			WidthMax:    maxColWidth,
			Number:      colIndex + 1,
			WidthMaxEnforcer: func(str string, warpLen int) string {
				// This effectively disables wrapping.
				// https://github.com/jedib0t/go-pretty/blob/main/table/config.go#L68C2-L68C18
				return str
			},
		}

		// assuming the table has headers: store columns with wide content where the max width may need to be adjusted
		// using the remaining space
		if t.maxSizes[colIndex] > maxColWidth && (colIndex < len(t.headers) && isWideColumn(t.headers[colIndex])) {
			widestColIndicies = append(widestColIndicies, colIndex)
		} else if t.maxSizes[colIndex] < maxColWidth {
			// use the max column width instead of the estimated max column width
			// if it is shorter
			columnConfig[colIndex].WidthMax = t.maxSizes[colIndex]
			remainingSpace -= t.maxSizes[colIndex]
		} else {
			remainingSpace -= maxColWidth
		}
	}

	// if only one wide column use the remaining space as the max column width
	if len(widestColIndicies) == 1 {
		widestColIndx := widestColIndicies[0]
		columnConfig[widestColIndx].WidthMax = remainingSpace
	}

	// if more than one wide column, spread the remaining space between the columns
	if len(widestColIndicies) > 1 {
		remainingSpace /= len(widestColIndicies)
		for _, columnCfgIdx := range widestColIndicies {
			columnConfig[columnCfgIdx].WidthMax = remainingSpace
		}

		origRemainingSpace := remainingSpace
		moreRemainingSpace := origRemainingSpace % len(widestColIndicies)
		if moreRemainingSpace != 0 {
			columnConfig[0].WidthMax += moreRemainingSpace
		}
	}

	return columnConfig
}

func (t *PrintableTable) createPrettyRowsAndHeaders() (headerRow table.Row, rows []table.Row) {
	for _, header := range t.headers {
		headerRow = append(headerRow, header)
	}

	for i := range t.rows {
		var row, emptyRow table.Row
		for j, cell := range t.rows[i] {
			if j == 0 {
				cell = TableContentHeaderColor(cell)
			}
			row = append(row, cell)
			emptyRow = append(emptyRow, "")
		}
		if i == 0 && len(t.headers) == 0 {
			rows = append(rows, emptyRow)
		}
		rows = append(rows, row)
	}

	return
}

func (t *PrintableTable) calculateMaxSize(row []string) {
	for index, value := range row {
		cellLength := runewidth.StringWidth(Decolorize(value))
		if t.maxSizes[index] < cellLength {
			t.maxSizes[index] = cellLength
		}
	}
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

func (t *PrintableTable) PrintCsv() error {
	csvwriter := csv.NewWriter(t.writer)
	err := csvwriter.Write(t.headers)
	if err != nil {
		return fmt.Errorf(T("Failed, header could not convert to csv format"), err.Error())
	}
	err = csvwriter.WriteAll(t.rows)
	if err != nil {
		return fmt.Errorf(T("Failed, rows could not convert to csv format"), err.Error())
	}
	return nil
}
