package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
)

const (
	defaultWidth = 120
)

//go:embed README.md
var readmeContent []byte

func readLine(f *os.File) (string, error) {
	var line []byte
	buf := make([]byte, 1)

	for {
		n, err := f.Read(buf)
		if n > 0 {
			line = append(line, buf[0])
			if buf[0] == '\n' {
				return string(line), nil
			}
		}
		if err != nil {
			if err == io.EOF && len(line) > 0 {
				return string(line), nil
			}
			return string(line), err
		}
	}
}

func skipShebangIfNeeded(f *os.File) error {
	line, err := readLine(f)
	if err != nil && err != io.EOF {
		return err
	}

	appName := filepath.Base(os.Args[0])
	if !strings.HasPrefix(line, "#!") || !strings.Contains(line, appName) {
		_, err := f.Seek(0, 0)
		return err
	}

	// Skip following empty lines
	for {
		posBeforeRead, _ := f.Seek(0, io.SeekCurrent)

		line, err := readLine(f)
		if err != nil {
			break
		}

		if strings.TrimSpace(line) != "" {
			_, err = f.Seek(posBeforeRead, 0)
			return err
		}
	}

	return nil
}

func printVersion() {
	if info, ok := debug.ReadBuildInfo(); ok {
		out, _ := json.Marshal(info.Main)
		fmt.Println(string(out))
	} else {
		fmt.Println("{}")
	}
	os.Exit(0)
}

func printStyles() {
	styles := GetAvailableStyles()
	for _, style := range styles {
		fmt.Printf("%s\n", style)
	}
	os.Exit(0)
}

func printHelp(renderer *glamour.TermRenderer) {
	process(bytes.NewReader(readmeContent), renderer)
	os.Exit(0)
}

func main() {
	width, err := strconv.Atoi(os.Getenv("GLAMOUR_WIDTH"))
	if err != nil {
		width = defaultWidth
	}

	styleName := os.Getenv("GLAMOUR_STYLE")
	if styleName == "" {
		styleName = defaultStyle
	}

	customConfig := GetStyleConfig(styleName)

	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(customConfig),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating renderer: %s\n", err)
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		arg := os.Args[1]
		switch arg {
		case "--version", "-v":
			printVersion()
		case "--help", "-h":
			printHelp(r)
		case "--styles", "-s":
			printStyles()
		}
	}

	if len(os.Args) == 1 {
		process(os.Stdin, r)
		return
	}

	for _, filename := range os.Args[1:] {
		if filename == "-" {
			process(os.Stdin, r)
			continue
		}

		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %s\n", filename, err)
			continue
		}
		defer f.Close()

		fileInfo, err := f.Stat()
		if err == nil && fileInfo.Mode().Perm()&0111 != 0 {
			if err := skipShebangIfNeeded(f); err != nil {
				fmt.Fprintf(os.Stderr, "Error processing file %s: %s\n", filename, err)
				continue
			}
		}

		process(f, r)
	}
}

func process(reader io.Reader, renderer *glamour.TermRenderer) {
	in, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
		return
	}

	md, err := renderer.RenderBytes(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering markdown: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stdout, "%s", md)
}
