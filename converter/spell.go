package converter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// LoadSpells loads the D&D 5e spell list from JSON file
func LoadSpells() (map[string]bool, error) {
	// Get the directory of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	// Navigate to data/spells.json relative to project root
	projectRoot := filepath.Join(filepath.Dir(filename), "..")
	spellsPath := filepath.Join(projectRoot, "data", "spells.json")

	// Read the file
	data, err := os.ReadFile(spellsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spells.json: %w", err)
	}

	var spellList []string
	if err := json.Unmarshal(data, &spellList); err != nil {
		return nil, fmt.Errorf("failed to parse spells.json: %w", err)
	}

	// Create a map for O(1) lookup, case-insensitive
	spells := make(map[string]bool)
	for _, spell := range spellList {
		spells[strings.ToLower(spell)] = true
	}

	return spells, nil
}

// IsValidSpell checks if a spell name is in the D&D 5e spell list (case-insensitive)
func IsValidSpell(spellName string, spells map[string]bool) bool {
	if spellName == "" {
		return false
	}
	return spells[strings.ToLower(spellName)]
}

// ConvertSpellLinks converts {{spell:Name}} syntax to [spell]Name[/spell]
// Returns the converted text and a list of warnings for invalid spells
func ConvertSpellLinks(text string, spells map[string]bool) (string, []string) {
	warnings := []string{}

	// Pattern to match {{spell:SpellName}}
	spellPattern := regexp.MustCompile(`\{\{spell:([^}]*)\}\}`)

	result := spellPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Extract spell name
		submatches := spellPattern.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		spellName := strings.TrimSpace(submatches[1])

		// Validate spell
		if !IsValidSpell(spellName, spells) {
			if spellName == "" {
				warnings = append(warnings, "Empty spell name in {{spell:}}")
			} else {
				warnings = append(warnings, fmt.Sprintf("Unknown spell: %q", spellName))
			}
		}

		// Convert to D&D Beyond format, preserving original case
		return fmt.Sprintf("[spell]%s[/spell]", spellName)
	})

	return result, warnings
}
