package main

import (
	"character-tool/converter"
	"character-tool/formatter"
	"character-tool/parser"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Define flags
	inputFile := flag.String("input", "", "Path to input markdown file (required)")
	outputDir := flag.String("output", ".", "Output directory for generated files (default: current directory)")
	verbose := flag.Bool("verbose", false, "Show detailed validation warnings")

	flag.Parse()

	// Validate input flag
	if *inputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: -input flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Run the tool
	if err := run(*inputFile, *outputDir, *verbose); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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
			fmt.Printf("  %d warning(s) found. Use -verbose flag for details.\n", len(allWarnings))
		}
	}

	fmt.Println("\nConversion complete!")
	return nil
}
