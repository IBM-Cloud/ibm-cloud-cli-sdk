package terminal

import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
)

func TestPrintableTable(t *testing.T) {
	for _, tc := range []struct{
		name string
		headers []string
		rows [][]string
		collapsibleColumns []string
		expected []string
	}{
		{
			name: "no collapsible columns",
			headers: []string{"name", "message", "comments"},
			rows: [][]string{
				{"foo", "message one", "this works fine"},
				{"bar", "message two", "tables print as expected"},
				{"baz", "message three", "no problem"},
			},
			expected: []string{
				"name", "message", "comments",
				"foo", "message one", "this works fine",
				"bar", "message two", "tables print as expected",
				"baz", "message three", "no problem",

			},
		},
		{
			name: "collapse no effect",
			headers: []string{"name", "message", "comments"},
			rows: [][]string{
				{"foo", "message one", "this works fine"},
				{"bar", "message two", "tables print as expected"},
				{"baz", "message three", "no problem"},
			},
			collapsibleColumns: []string{"name", "message"},
			expected: []string{
				"name", "message", "comments",
				"foo", "message one", "this works fine",
				"bar", "message two", "tables print as expected",
				"baz", "message three", "no problem",

			},
		},
		{
			name: "collapse",
			headers: []string{"name", "message", "comments", "other"},
			rows: [][]string{
				{"foo", "", "this works fine", ""},
				{"bar", "", "tables print as expected", "not empty"},
				{"baz", "", "no problem", ""},
			},
			collapsibleColumns: []string{"message", "other"},
			expected: []string{
				"name", "comments", "other",
				"foo", "this works fine",
				"bar", "tables print as expected", "not empty",
				"baz", "no problem",

			},
		},
		{
			name: "collapse last row",
			headers: []string{"name", "message", "comments", "other"},
			rows: [][]string{
				{"foo", "", "this works fine", ""},
				{"bar", "not empty", "tables print as expected", ""},
				{"baz", "", "no problem", ""},
			},
			collapsibleColumns: []string{"message", "other"},
			expected: []string{
				"name", "message", "comments",
				"foo", "this works fine",
				"bar", "not empty", "tables print as expected",
				"baz", "no problem",

			},
		},
	}{
		t.Run(tc.name, func(t *testing.T){
			w := testTableWriter{}
			table := NewTable(&w, tc.headers, WithCollapsibleColumns(tc.collapsibleColumns...))
			for _, row := range tc.rows {
				table.Add(row...)
			}
			table.Print()

			assert.Len(t, w.entries, len(tc.expected), "wrong number of observed entries in printed table")
			for i := range tc.expected {
				assert.Equal(t, tc.expected[i], w.entries[i], "table entry %d is incorrect", i)
			}
		})
	}
}

func filterSlice(sl []string, f func(string) bool) []string {
	newsl := make([]string, 0, len(sl))
	for _, elem := range sl {
		if !f(elem) {
			continue
		}
		newsl = append(newsl, elem)
	}
	return newsl
}

type testTableWriter struct {
	entries []string
}

func (w *testTableWriter) Write(p []byte) (n int, err error) {
	// split rows/headings to individual elems
	elems := strings.Split(string(p), "   ")

	// sanitize input to remove color/bold formatting and extra whitespace
	for i := range elems {
		elems[i] = Decolorize(strings.TrimSpace(elems[i]))
	}

	// filter out empty elements
	elems = filterSlice(elems, func(s string) bool {
		return s != ""
	})

	// record table entries
	w.entries = append(w.entries, elems...)
	return len(p), nil
}
