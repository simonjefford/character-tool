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
	vaultMode bool
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
		return run(inputFile, outputDir, verbose, vaultMode)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "path to input markdown file (required)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "output directory for generated files")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show detailed validation warnings")
	rootCmd.Flags().BoolVar(&vaultMode, "vault-mode", false, "output files to same directory as input file (Obsidian integration)")
	rootCmd.MarkFlagRequired("input")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(inputFile, outputDir string, verbose, vaultMode bool) error {
	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// In vault mode, output to the same directory as input file
	if vaultMode {
		outputDir = filepath.Dir(inputFile)
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

	// Track all warnings and created files
	var allWarnings []string
	var createdFiles []string
	totalAbilities := 0

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
		if len(section.abilities) == 0 {
			continue
		}

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

		// Track created files
		createdFiles = append(createdFiles, fmt.Sprintf("%s (%d abilities)", outputPath, len(section.abilities)))
		totalAbilities += len(section.abilities)
	}

	// Display summary
	fmt.Println("✓ Formatted character abilities")
	fmt.Println("\nOutput files:")
	for _, file := range createdFiles {
		fmt.Printf("  - %s\n", file)
	}

	// Display warnings summary
	if len(allWarnings) > 0 {
		fmt.Println()
		if verbose {
			fmt.Println("Warnings:")
			for _, warning := range allWarnings {
				fmt.Printf("  ! %s\n", warning)
			}
		} else {
			fmt.Printf("Warnings: %d found (use --verbose for details)\n", len(allWarnings))
		}
	}

	return nil
}
