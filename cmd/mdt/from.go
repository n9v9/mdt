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
	alignments := make([]mdt.TableAlignment, 0, len(align))

	for _, v := range align {
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
			return fmt.Errorf("invalid alignment character %c", v)
		}
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
