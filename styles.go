package main

import (
	"sort"

	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
)

func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }
func uintPtr(u uint) *uint       { return &u }

const (
	// defaultListIndent      = 2
	// defaultListLevelIndent = 4
	defaultMargin = 0
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
// and explicitly overrides the document margins and padding for flush CLI output.
func GetStyleConfig(styleName string) ansi.StyleConfig {
	configPtr, ok := styles.DefaultStyles[styleName]
	if !ok {
		configPtr = styles.DefaultStyles[defaultStyle]
	}

	config := *configPtr

	// To prevent mutating global default styles while stripping padding,
	// explicitly recreate the nested structs to override
	docColor := config.Document.Color
	config.Document = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: docColor,
		},
		Margin: uintPtr(defaultMargin),
	}

	config.CodeBlock.Margin = uintPtr(defaultMargin)

	return config
}
