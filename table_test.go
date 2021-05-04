package mdt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func diff(t *testing.T, msg string, want, got interface{}) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("mismatch: %s (-want, +got):\n%s", msg, diff)
	}
}

func TestTable(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		got := (&Table{Rows: nil}).String()
		diff(t, `empty string`, "", got)
	})

	t.Run("no header", func(t *testing.T) {
		table := &Table{
			Rows: [][]string{
				{"id", "firstname"},
				{"1", "john"},
			},
			NoHeader: true,
		}
		want := "| id | firstname |\n| 1  | john      |"
		diff(t, "header", want, table.String())
	})

	t.Run("invalid alignment panics", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic")
			}
		}()
		table := &Table{
			Rows: [][]string{
				{"id", "firstname"},
				{"1", "john"},
			},
			// Valid values are 0 - 3 as defined by the TableAlignment constants.
			Alignments: []TableAlignment{TableAlignment(4)},
		}
		_ = table.String()
	})

	t.Run("alignment", func(t *testing.T) {
		testCases := []struct {
			name      string
			separator string
			align     []TableAlignment
		}{
			{
				"default",
				"| --- | --- | --- |",
				[]TableAlignment{AlignDefault},
			},
			{
				"left",
				"| :-- | :-- | :-- |",
				[]TableAlignment{AlignLeft},
			},
			{
				"center",
				"| :-: | :-: | :-: |",
				[]TableAlignment{AlignCenter},
			},
			{
				"right",
				"| --: | --: | --: |",
				[]TableAlignment{AlignRight},
			},
			{
				"no alignment",
				"| --- | --- | --- |",
				nil,
			},
			{
				"all different",
				"| :-- | :-: | --: |",
				[]TableAlignment{AlignLeft, AlignCenter, AlignRight},
			},
			{
				"only some given",
				"| :-- | :-: | --- |",
				[]TableAlignment{AlignLeft, AlignCenter},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var (
					rows = [][]string{
						{"a", "b", "c"},
						{"1", "2", "3"},
					}
					want = "| a   | b   | c   |" + "\n" +
						tc.separator + "\n" +
						"| 1   | 2   | 3   |"
				)
				got := (&Table{Rows: rows, Alignments: tc.align}).String()
				diff(t, "rows", want, got)
			})
		}
	})
}

func TestParseTable(t *testing.T) {
	t.Run("no header", func(t *testing.T) {
		markdown := `| id | firstname | lastname |
| 1 | john | doe |
| 2 | jane | doe |
| 3 | escaped \| pipe | test |`

		want := &Table{
			Rows: [][]string{
				{"id", "firstname", "lastname"},
				{"1", "john", "doe"},
				{"2", "jane", "doe"},
				{"3", "escaped | pipe", "test"},
			},
			NoHeader: true,
		}

		got, err := ParseTable(strings.NewReader(markdown), true)
		if err != nil {
			t.Fatal(err)
		}

		diff(t, "no header", want, got)
	})

	t.Run("alignment", func(*testing.T) {
		input := "| --- | :-- | :-: | --: | ----- | :---- | :---: | ----: |"
		want := []TableAlignment{
			AlignDefault,
			AlignLeft,
			AlignCenter,
			AlignRight,
			AlignDefault,
			AlignLeft,
			AlignCenter,
			AlignRight,
		}

		diff(t, "alignment", want, parseAlignment(input))
	})
}

func ExampleTable_default() {
	rows := [][]string{
		{"id", "firstname", "lastname"},
		{"1", "john", "doe"},
		{"2", "jane", "doe"},
		{"3", "max", "mustermann"},
	}

	t := &Table{Rows: rows}
	fmt.Print(t)
	// Output:
	// | id  | firstname | lastname   |
	// | --- | --------- | ---------- |
	// | 1   | john      | doe        |
	// | 2   | jane      | doe        |
	// | 3   | max       | mustermann |
}

func ExampleTable_alignSome() {
	rows := [][]string{
		{"id", "firstname", "lastname"},
		{"1", "john", "doe"},
		{"2", "jane", "doe"},
		{"3", "max", "mustermann"},
	}

	// Will align the first column to the left, the second to the right
	// and the third with the default alignment.
	t := &Table{Rows: rows, Alignments: []TableAlignment{AlignLeft, AlignRight}}
	fmt.Print(t)
	// Output:
	// | id  | firstname | lastname   |
	// | :-- | --------: | ---------- |
	// | 1   | john      | doe        |
	// | 2   | jane      | doe        |
	// | 3   | max       | mustermann |
}

func ExampleTable_alignOne() {
	rows := [][]string{
		{"id", "firstname", "lastname"},
		{"1", "john", "doe"},
		{"2", "jane", "doe"},
		{"3", "max", "mustermann"},
	}

	// Will align all columns according to the single alignment given.
	t := &Table{Rows: rows, Alignments: []TableAlignment{AlignCenter}}
	fmt.Print(t)
	// Output:
	// | id  | firstname | lastname   |
	// | :-: | :-------: | :--------: |
	// | 1   | john      | doe        |
	// | 2   | jane      | doe        |
	// | 3   | max       | mustermann |
}

func ExampleParseTable() {
	markdown := `| id  | firstname | lastname   |
| --- | :-------- | :--------: |
| 1   | john      | doe        |
| 2   | jane      | doe        |
| 3   | max       | mustermann |`

	t, _ := ParseTable(strings.NewReader(markdown), false)
	// t now contains the parsed rows and alignments
	// and can be converted back to markdown.
	fmt.Print(t)
	// Output:
	// | id  | firstname | lastname   |
	// | --- | :-------- | :--------: |
	// | 1   | john      | doe        |
	// | 2   | jane      | doe        |
	// | 3   | max       | mustermann |
}
