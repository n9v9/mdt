# mdt CLI

CLI to format markdown tables and convert between markdown tables and the CSV
format.

I often find myself working with markdown files and I think building and
especially modifying markdown tables is clumsy so I made this CLI. I use it
inside vim with the following keybindings:

```vim
vnoremap <silent><leader>mdf :!mdt fmt<CR>
vnoremap <silent><leader>mdt :!mdt md<CR>
vnoremap <silent><leader>mdc :!mdt csv<CR>
```

## Examples

### Markdown

Convert data from the CSV format into markdown tables:

![Example: md command](svgs/example_md.svg)

### CSV

Convert markdown tables into the CSV format:

![Example: csv command](svgs/example_csv.svg)

### Format

Format markdown tables for better readability:

![Example: fmt command](svgs/example_fmt.svg)
