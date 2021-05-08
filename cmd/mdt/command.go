package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/n9v9/mdt"
	"github.com/urfave/cli/v2"
)

func fromCsv(ctx *cli.Context) (err error) {
	align := ctx.String("align")
	alignments, err := parseAlignments(align)
	if err != nil {
		return err
	}

	delimiter, _ := utf8.DecodeRuneInString(ctx.String("delimiter"))

	rc := inputFile(ctx.Context)
	defer func() {
		closeErr := rc.Close()
		if err == nil {
			err = closeErr
		}
	}()

	r := csv.NewReader(rc)
	r.TrimLeadingSpace = true
	r.Comma = delimiter

	records, err := r.ReadAll()
	if err != nil {
		return
	}

	table := &mdt.Table{
		Rows:       records,
		Alignments: alignments,
		NoHeader:   ctx.Bool("no-header"),
	}

	_, err = fmt.Println(table)

	return
}

func fromMarkdown(ctx *cli.Context) (err error) {
	rc := inputFile(ctx.Context)
	defer func() {
		closeErr := rc.Close()
		if err == nil {
			err = closeErr
		}
	}()

	table, err := mdt.ParseTable(rc, ctx.Bool("no-header"))
	if err != nil {
		return
	}

	w := csv.NewWriter(os.Stdout)
	w.Comma = rune(ctx.String("delimiter")[0])

	return w.WriteAll(table.Rows)
}

func fmtMarkdown(ctx *cli.Context) (err error) {
	rc := inputFile(ctx.Context)
	defer func() {
		closeErr := rc.Close()
		if err == nil {
			err = closeErr
		}
	}()

	table, err := mdt.ParseTable(rc, ctx.Bool("no-header"))
	if err != nil {
		return err
	}

	// Only change the alignment when explicitly set; otherwise we already got the correct
	// alignment by parsing the table.
	if align := ctx.String("align"); align != "" {
		alignments, err := parseAlignments(align)
		if err != nil {
			return err
		}
		table.Alignments = alignments
	}

	_, err = fmt.Println(table)
	return
}

func parseAlignments(alignRunes string) ([]mdt.TableAlignment, error) {
	alignments := make([]mdt.TableAlignment, 0, len(alignRunes))

	for _, v := range alignRunes {
		switch v {
		case 'd':
			alignments = append(alignments, mdt.AlignDefault)
		case 'l':
			alignments = append(alignments, mdt.AlignLeft)
		case 'r':
			alignments = append(alignments, mdt.AlignRight)
		case 'c':
			alignments = append(alignments, mdt.AlignCenter)
		default:
			return nil, fmt.Errorf("invalid alignment character %c", v)
		}
	}

	return alignments, nil
}
