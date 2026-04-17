package main

import (
	"sort"
	"strconv"

	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
)

func uintPtr(u uint) *uint { return &u }

const (
	// defaultListIndent      = 2
	// defaultListLevelIndent = 4
	defaultMargin = "0"
	defaultStyle  = "dark"
)

// GetAvailableStyles returns a sorted list of all available styles
func GetAvailableStyles() []string {
	var available []string
	for k := range styles.DefaultStyles {
		available = append(available, k)
	}
	sort.Strings(available)
	return available
}

// GetStyleConfig loads a named style from Glamour's default styles
// and explicitly overrides the document margins and padding for flush CLI output
func GetStyleConfig() ansi.StyleConfig {
	styleName := getEnv("GLAMOUR_STYLE", defaultStyle)

	configPtr, ok := styles.DefaultStyles[styleName]
	if !ok {
		configPtr = styles.DefaultStyles[defaultStyle]
	}

	config := *configPtr

	parsedMargin, err := strconv.ParseUint(getEnv("GLAMOUR_OVERRIDE_MARGIN", defaultMargin), 10, 32)
	if err != nil {
		parsedMargin, _ = strconv.ParseUint(defaultMargin, 10, 32)
	}
	margin := uint(parsedMargin)

	parsedCodeMargin, err := strconv.ParseUint(getEnv("GLAMOUR_OVERRIDE_MARGIN_CODE", defaultMargin), 10, 32)
	if err != nil {
		parsedCodeMargin, _ = strconv.ParseUint(defaultMargin, 10, 32)
	}
	codeMargin := uint(parsedCodeMargin)

	// To prevent mutating global default styles while stripping padding,
	// explicitly recreate the nested structs to override
	docColor := config.Document.Color
	config.Document = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: docColor,
		},
		Margin: uintPtr(margin),
	}

	config.CodeBlock.Margin = uintPtr(codeMargin)

	return config
}
