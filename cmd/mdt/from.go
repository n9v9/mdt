package main

import (
	"encoding/csv"
	"fmt"
	"unicode/utf8"

	"github.com/n9v9/mdt"
	"github.com/urfave/cli/v2"
)

func fromCsv(ctx *cli.Context) error {
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
	defer rc.Close()

	r := csv.NewReader(rc)
	r.TrimLeadingSpace = true
	r.Comma = delimiter

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	table := &mdt.Table{
		Rows:       records,
		Alignments: alignments,
	}

	fmt.Println(table)

	return nil
}
