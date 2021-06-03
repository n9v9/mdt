// Package mdt implements the building and parsing of markdown tables.
package mdt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// TableAlignment specifies the alignment of columns in a table.
type TableAlignment int

// Alignments for columns in a table.
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

// String returns the table as a markdown string.
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

// ParseTable reads markdown from r and parses it into a Table.
// r must provide a valid markdown table representation.
// If noHeader is true then the second row, the alignment row, will be parsed as normal data.
func ParseTable(r io.Reader, noHeader bool) (*Table, error) {
	var (
		rows       [][]string
		alignments []TableAlignment
		s          = bufio.NewScanner(r)
		rowIdx     int
	)

	for s.Scan() {
		row := strings.TrimSpace(s.Text())
		if rowIdx == 1 && !noHeader {
			alignments = parseAlignment(row)
			rowIdx++
			continue
		}

		var (
			cols        []string
			col         strings.Builder
			pipeEscaped bool
		)
		// Skip the first character as it must always be the pipe character.
		for i := 1; i < len(row); i++ {
			switch v := row[i]; v {
			case '\\':
				// To show a literal pipe character in row, it must be escaped.
				// When parsing we don't escape it.
				if row[i+1] == '|' {
					pipeEscaped = true
				} else {
					col.WriteByte(v)
				}
			case '|':
				if pipeEscaped {
					pipeEscaped = false
					col.WriteByte(v)
				} else {
					cols = append(cols, strings.TrimSpace(col.String()))
					col.Reset()
				}
			default:
				col.WriteByte(v)
			}
		}

		rows = append(rows, cols)
		rowIdx++
	}

	return &Table{
		Rows:       rows,
		Alignments: alignments,
		NoHeader:   noHeader,
	}, s.Err()
}

func parseAlignment(row string) []TableAlignment {
	var alignments []TableAlignment

	cols := strings.Split(row, "|")

	// Skip the first and last items as they are empty because the row must start and end with
	// the pipe character.
	for _, col := range cols[1 : len(cols)-1] {
		col := strings.TrimSpace(col)
		if col[0] == ':' && col[len(col)-1] == ':' {
			alignments = append(alignments, AlignCenter)
		} else if col[0] == ':' {
			alignments = append(alignments, AlignLeft)
		} else if col[len(col)-1] == ':' {
			alignments = append(alignments, AlignRight)
		} else {
			alignments = append(alignments, AlignDefault)
		}
	}

	return alignments
}
