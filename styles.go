package main

import (
	"sort"

	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
)

//go:fix inline
func uintPtr(u uint) *uint { return new(u) }

const (
	// defaultListIndent      = 2
	// defaultListLevelIndent = 4
	defaultMargin = 0
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
	styleName := getEnv("GLAMOUR_STYLE", detectStyleByTerminalBrightness)

	configPtr, ok := styles.DefaultStyles[styleName]
	if !ok {
		configPtr = styles.DefaultStyles[defaultStyleDark]
	}

	config := *configPtr

	margin := uint(getEnvUint("GLAMOUR_OVERRIDE_MARGIN", func() uint64 {
		return defaultMargin
	}))

	codeMargin := uint(getEnvUint("GLAMOUR_OVERRIDE_MARGIN_CODE", func() uint64 {
		return defaultMargin
	}))

	// To prevent mutating global default styles while stripping padding,
	// explicitly recreate the nested structs to override
	docColor := config.Document.Color
	config.Document = ansi.StyleBlock{
		Color:  docColor,
		Margin: new(margin),
	}

	config.CodeBlock.Margin = new(codeMargin)

	return config
}
