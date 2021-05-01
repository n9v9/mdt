// Package mdt provides an API to work with markdown tables.
package mdt

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

// TableAlignment specifies the alignment of columns in a table.
type TableAlignment int

// Aligments for columns in a table.
const (
	AlignDefault TableAlignment = iota
	AlignLeft
	AlignCenter
	AlignRight
)

// Table represents a markdown table with a given alignment.
type Table struct {
	Rows [][]string
	// Determines whether the first entry of Rows represents the table header.
	NoHeader bool
	// The first alignment corresponds to the first column in Rows, the second alignment to the
	// second column in Rows and so on.
	//
	// If Alignments is nil or contains no elements, AlignDefault will be used for all columns.
	// If only one alignment is given, it will be used for all columns.
	// If more than one alignment is given but not enough for all columns then AlignDefault will be used
	// for the remaining columns.
	Alignments []TableAlignment
}

func (t *Table) String() string {
	// Map each column index to its maximum length:
	// 0 -> max len of all strings in column 0
	// 1 -> max len of all strings in column 1
	// ...
	headerSepLen := make(map[int]int)

	// Even though we use tabwriter to get the correct width for the table, we have to
	// calculate the width for the separator ourselves.
	for _, record := range t.Rows {
		for i, v := range record {
			// The pipe character `|` must be escaped because it is used as a column delimiter.
			// Because we don't modify rows, we count its occurrences here and take them into
			// account when calculating the width.
			// The escaping takes place when printing further down.
			columnLen := len(v) + strings.Count(v, "|")
			if length, ok := headerSepLen[i]; !ok {
				// Each separator must be at least 3 chars long.
				// This way, alignment can be correctly specified:
				// --- default
				// :-- left
				// :-: center
				// --: right
				if columnLen < 3 {
					headerSepLen[i] = 3
				} else {
					headerSepLen[i] = columnLen
				}
			} else if length < columnLen {
				headerSepLen[i] = columnLen
			}
		}
	}

	// Create the dashes based on the width and alignment.
	dashes := func(sepLen int, columnIdx int) string {
		var alignment TableAlignment

		switch v := len(t.Alignments); {
		case v == 1:
			alignment = t.Alignments[0]
		case columnIdx < v:
			alignment = t.Alignments[columnIdx]
		default:
			alignment = AlignDefault
		}

		switch alignment {
		case AlignDefault:
			return strings.Repeat("-", sepLen)
		case AlignLeft:
			return ":" + strings.Repeat("-", sepLen-1)
		case AlignRight:
			return strings.Repeat("-", sepLen-1) + ":"
		case AlignCenter:
			return ":" + strings.Repeat("-", sepLen-2) + ":"
		default:
			panic(fmt.Sprintf("mdt: invalid TableAlignment: %v", alignment))
		}
	}

	const (
		column           = "| %s \t"
		columnEnd        = "|\t"
		columnEndNewLine = "|\t\n"
	)

	buf := &bytes.Buffer{}
	tw := tabwriter.NewWriter(buf, 0, 0, 0, ' ', 0)
	writeHeader := !t.NoHeader

	for rowIdx, row := range t.Rows {
		for _, v := range row {
			if writeHeader && rowIdx == 1 {
				writeHeader = false
				for i := 0; i < len(row); i++ {
					fmt.Fprintf(tw, column, dashes(headerSepLen[i], i))
				}
				fmt.Fprint(tw, columnEndNewLine)
			}
			fmt.Fprintf(tw, column, strings.ReplaceAll(v, "|", "\\|"))
		}
		if rowIdx == len(t.Rows)-1 {
			// After the last row should be no new line.
			fmt.Fprint(tw, columnEnd)
		} else {
			fmt.Fprint(tw, columnEndNewLine)
		}
	}

	// Ignore error as Flush calls Write on the buffer and Write on bytes.Buffer always returns nil.
	_ = tw.Flush()

	return buf.String()
}
