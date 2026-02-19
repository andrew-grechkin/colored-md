# colored-md

`colored-md` is a command-line utility for rendering Markdown files with syntax highlighting and formatting directly in your terminal. It's built on top of the `charmbracelet/glamour` library, providing a simple wrapper to bring rich, colored Markdown output to your CLI workflows.

## Features

- Renders Markdown to ANSI-colored terminal output.
- Customizable output width via `GLAMOUR_WIDTH` environment variable.
- Processes both piped input and specified file paths.

## Installation

To install `colored-md`:

```bash
go install github.com/andrew-grechkin/colored-md@latest
```

## Usage

### Piping input

You can pipe Markdown content directly to `colored-md`:

```bash
echo "# Hello World

This is **bold** text." | colored-md
```

### Processing files

Specify one or more Markdown files as arguments:

```bash
colored-md README.md my_document.md
```

To read from standard input while also processing files, use `-` as a filename:

```bash
cat my_file.md | colored-md - another_file.md
```

### Customizing output width

Set the `GLAMOUR_WIDTH` environment variable to control the word wrap width:

```bash
GLAMOUR_WIDTH=80 colored-md my_file.md
```

### Styling

`colored-md` uses a predefined `DarkStyleConfig`. While `glamour` typically allows setting styles via `GLAMOUR_STYLE`
environment variable, `colored-md` overrides this with its internal `DarkStyleConfig`. Future versions might expose more
style customization options. This is done because of weird un-configurable marging of several spaces all default
`glamour` styles have. They also always add empty lines around document, which is undesirable.

## Why `colored-md` over `glow`?

While `glow` is a popular Markdown renderer, `colored-md` was created to address specific shortcomings when used as a CLI filter:

- **Undesired directory traversal**: `glow` has functionality to traverse directories and find Markdown files, which is often not desired when simply piping content or processing specific files.
- **Forced paging in pipelines**: `glow` tends to force paging (e.g., using `less`) even when used in pipelines, which can disrupt CLI workflows where direct output is expected.

`colored-md` aims to be a simpler, more predictable CLI filter for Markdown rendering, following strictly UNIX philosophy.

## Author

- Andrew Grechkin

## License

This project is licensed under the GNU General Public License Version 2 (GPLv2). See the `LICENSE` file for details.
