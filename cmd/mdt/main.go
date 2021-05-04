package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:            "mdt",
		Usage:           "Helps you work with markdown tables",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "Path to the `FILE` from which to read or - to read from stdin",
				Value:   "-",
			},
			&cli.BoolFlag{
				Name:  "no-header",
				Usage: "Do not interpret the first row as the table header",
			},
			&cli.StringFlag{
				Name:  "delimiter",
				Usage: "CSV field delimiter character",
				Value: ",",
			},
		},
		Commands: []*cli.Command{
			{
				Name:            "csv",
				Usage:           "Convert data in a CSV format into a markdown table",
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "align",
						Usage: "Sequence of alignment characters `dlrc `for each column. default (d), left (l), right (r) and center (c).",
					},
				},
				Action: fromCsv,
			},
			{
				Name:            "md",
				Usage:           "Convert a markdown table into CSV formatted data",
				HideHelpCommand: true,
				Action:          fromMarkdown,
			},
		},
		Before: func(ctx *cli.Context) error {
			switch v := ctx.String("file"); v {
			case "-":
				ctx.Context = withInputFile(ctx.Context, os.Stdin)
			default:
				f, err := os.Open(v)
				if err != nil {
					return err
				}
				ctx.Context = withInputFile(ctx.Context, f)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
