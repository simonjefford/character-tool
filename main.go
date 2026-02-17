package main

import (
	"character-tool/converter"
	"character-tool/formatter"
	"character-tool/parser"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	inputFile string
	outputDir string
	verbose   bool
)

var rootCmd = &cobra.Command{
	Use:   "character-tool",
	Short: "Convert D&D character markdown to D&D Beyond format",
	Long: `A tool that converts markdown character descriptions into D&D Beyond-formatted
blocks with spell links and rollable dice notation.

The tool parses markdown files with structured headers (Traits, Actions, Bonus Actions,
Reactions) and converts:
  - {{spell:SpellName}} syntax to clickable spell links
  - Dice notation (1d20+5) with keywords (to hit:, damage:) to rollable format
  - Validates spell names and dice notation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(inputFile, outputDir, verbose)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "path to input markdown file (required)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "output directory for generated files")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show detailed validation warnings")
	rootCmd.MarkFlagRequired("input")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(inputFile, outputDir string, verbose bool) error {
	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Parse markdown
	result, err := parser.ParseMarkdown(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse markdown: %w", err)
	}

	// Load spells
	spells, err := converter.LoadSpells()
	if err != nil {
		return fmt.Errorf("failed to load spell list: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Track all warnings
	var allWarnings []string

	// Process and write each section
	sections := map[string]struct {
		abilities []parser.Ability
		filename  string
	}{
		"Traits":        {result.Traits, "traits.txt"},
		"Actions":       {result.Actions, "actions.txt"},
		"Bonus Actions": {result.BonusActions, "bonus-actions.txt"},
		"Reactions":     {result.Reactions, "reactions.txt"},
	}

	for sectionName, section := range sections {
		formatted, warnings, err := formatter.FormatAbilities(section.abilities, spells)
		if err != nil {
			return fmt.Errorf("failed to format %s: %w", sectionName, err)
		}

		// Collect warnings
		for _, warning := range warnings {
			allWarnings = append(allWarnings, fmt.Sprintf("[%s] %s", sectionName, warning))
		}

		// Write output file
		outputPath := filepath.Join(outputDir, section.filename)
		if err := os.WriteFile(outputPath, []byte(formatted), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", section.filename, err)
		}

		// Report section status
		if len(section.abilities) > 0 {
			fmt.Printf("✓ %s: %d abilities written to %s\n", sectionName, len(section.abilities), section.filename)
		} else {
			fmt.Printf("- %s: empty (no file created)\n", sectionName)
			// Remove empty file if it exists
			os.Remove(outputPath)
		}
	}

	// Display warnings
	if len(allWarnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warning := range allWarnings {
			if verbose {
				fmt.Printf("  ! %s\n", warning)
			}
		}
		if !verbose {
			fmt.Printf("  %d warning(s) found. Use --verbose flag for details.\n", len(allWarnings))
		}
	}

	fmt.Println("\nConversion complete!")
	return nil
}
