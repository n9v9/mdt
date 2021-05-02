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
		Commands: []*cli.Command{
			{
				Name:            "from",
				Usage:           "Create markdown tables from data in a given format",
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "file",
						Usage: "Path to the `FILE` from which to read or - to read from stdin",
						Value: "-",
					},
					&cli.StringFlag{
						Name:  "align",
						Usage: "Sequence of alignment characters `dlrc `for each column. default (d), left (l), right (r) and center (c).",
					},
					&cli.BoolFlag{
						Name:  "no-header",
						Usage: "Do not interpret the first row as the table header",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:            "csv",
						Usage:           "Interpret the input data as csv",
						HideHelpCommand: true,
						Action:          fromCsv,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "delimiter",
								Usage: "CSV field delimiter character",
								Value: ",",
							},
						},
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
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
