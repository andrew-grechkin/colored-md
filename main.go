package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

const (
	defaultStyleDark  = "dark"
	defaultStyleLight = "light"
	defaultWidth      = 120
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

func getEnv(key string, defaultFn func() string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultFn()
}

func getEnvInt(key string, defaultFn func() int64) int64 {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if parsed, err := strconv.ParseInt(val, 10, 32); err == nil {
			return parsed
		}
	}
	return defaultFn()
}

func getEnvUint(key string, defaultFn func() uint64) uint64 {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if parsed, err := strconv.ParseUint(val, 10, 32); err == nil {
			return parsed
		}
	}
	return defaultFn()
}

// queryTerminalBackground sends OSC 11 query and returns the terminal response
// Returns empty string if query fails or times out
func queryTerminalBackground() string {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return ""
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return ""
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	os.Stdout.Write([]byte("\033]11;?\007"))

	response := make([]byte, 64)
	readChan := make(chan int)
	go func() {
		n, _ := os.Stdin.Read(response)
		readChan <- n
	}()

	select {
	case n := <-readChan:
		return string(response[:n])
	case <-time.After(25 * time.Millisecond):
		return ""
	}
}

func hexToFloat(h string) float64 {
	val, _ := strconv.ParseUint(h, 16, 64)
	max := (1 << (uint(len(h)) * 4)) - 1
	return float64(val) / float64(max)
}

// parseTerminalBackgroundResponse parses OSC 11 response
func parseTerminalBackgroundResponse(res string) string {
	// The response format is usually: ^]11;rgb:rrrr/gggg/bbbb^G
	re := regexp.MustCompile(`rgb:([0-9a-fA-F]+)/([0-9a-fA-F]+)/([0-9a-fA-F]+)`)
	matches := re.FindStringSubmatch(res)
	if len(matches) == 4 {
		r := hexToFloat(matches[1])
		g := hexToFloat(matches[2])
		b := hexToFloat(matches[3])

		luma := (0.299 * r) + (0.587 * g) + (0.114 * b)
		if luma > 0.5 {
			return defaultStyleLight
		}
	}

	return defaultStyleDark
}

// detectStyleByTerminalBrightness determines terminal background brightness
func detectStyleByTerminalBrightness() string {
	if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
		parts := strings.Split(colorfgbg, ";")
		if len(parts) >= 2 {
			if bg, err := strconv.Atoi(parts[len(parts)-1]); err == nil && bg >= 7 {
				return defaultStyleLight
			}
			return defaultStyleDark
		}
	}

	if res := queryTerminalBackground(); res != "" {
		return parseTerminalBackgroundResponse(res)
	}

	return defaultStyleDark
}

func main() {
	width := int(getEnvInt("GLAMOUR_WIDTH", func() int64 {
		return defaultWidth
	}))

	customConfig := GetStyleConfig()

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
