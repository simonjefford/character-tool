package formatter

import (
	"character-tool/converter"
	"character-tool/parser"
	"strings"
)

// FormatAbilities formats a list of abilities with dice rolls and spell links converted
func FormatAbilities(abilities []parser.Ability, spells map[string]bool) (string, []string, error) {
	if len(abilities) == 0 {
		return "", []string{}, nil
	}

	var formatted []string
	var allWarnings []string

	for _, ability := range abilities {
		var text string
		if ability.Name != "" {
			// Named ability: format as "Name. Description"
			text = ability.Name + ". " + ability.Description
		} else {
			// Plain text paragraph: just the description
			text = ability.Description
		}

		// Convert spell links first
		text, spellWarnings := converter.ConvertSpellLinks(text, spells)
		allWarnings = append(allWarnings, spellWarnings...)

		// Convert dice rolls (use ability name as action name, or empty string for plain text)
		text, err := converter.ConvertDiceRolls(text, ability.Name)
		if err != nil {
			return "", allWarnings, err
		}

		formatted = append(formatted, text)
	}

	// Join abilities with blank lines
	result := strings.Join(formatted, "\n\n")

	return result, allWarnings, nil
}
