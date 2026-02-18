package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Valid dice types in D&D 5e
var validDice = map[string]bool{
	"d4":   true,
	"d6":   true,
	"d8":   true,
	"d10":  true,
	"d12":  true,
	"d20":  true,
	"d100": true,
}

// RollableData represents the JSON data embedded in rollable tags
type RollableData struct {
	DiceNotation string `json:"diceNotation"`
	RollType     string `json:"rollType"`
	RollAction   string `json:"rollAction"`
}

// ParseDiceNotation validates and normalizes dice notation
func ParseDiceNotation(notation string) (string, error) {
	notation = strings.TrimSpace(notation)

	// Match pattern: optional count + d + sides + optional modifier
	re := regexp.MustCompile(`^(\d*)d(\d+)([+-]\d+)?$`)
	matches := re.FindStringSubmatch(notation)

	if matches == nil {
		return "", errors.New("invalid dice notation format")
	}

	count := matches[1]
	if count == "" {
		count = "1"
	}

	sides := matches[2]
	modifier := matches[3]

	// Validate dice type
	diceType := "d" + sides
	if !validDice[diceType] {
		return "", fmt.Errorf("invalid dice type: %s (must be d4, d6, d8, d10, d12, d20, or d100)", diceType)
	}

	// Reconstruct normalized notation
	result := count + "d" + sides
	if modifier != "" {
		result += modifier
	}

	return result, nil
}

// ConvertDiceRolls converts dice notation with roll type keywords to D&D Beyond format
func ConvertDiceRolls(text string, actionName string) (string, error) {
	// Pattern to match roll type keywords followed by dice notation
	// Supports: to hit:, damage:, healing:, save:
	rollPattern := regexp.MustCompile(`(to hit|damage|healing|save):\s*(\d*d\d+[+-]?\d*)`)

	result := rollPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Extract roll type and dice notation
		parts := strings.SplitN(match, ":", 2)
		if len(parts) != 2 {
			return match
		}

		rollType := strings.TrimSpace(parts[0])
		diceNotation := strings.TrimSpace(parts[1])

		// Validate dice notation
		normalized, err := ParseDiceNotation(diceNotation)
		if err != nil {
			// Return original if invalid
			return match
		}

		// Get display value based on dice type
		displayValue := getDisplayValue(normalized)

		// Create rollable data
		data := RollableData{
			DiceNotation: normalized,
			RollType:     rollType,
			RollAction:   actionName,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return match
		}

		// Format: [rollable]DISPLAY_VALUE;JSON[/rollable]
		// For d20: displays modifier only (e.g., "+5")
		// For non-d20: displays full notation (e.g., "1d8+5")
		if displayValue == "" {
			return fmt.Sprintf("[rollable];%s[/rollable]", jsonData)
		}
		return fmt.Sprintf("[rollable]%s;%s[/rollable]", displayValue, jsonData)
	})

	return result, nil
}

// extractModifier extracts the +X or -X modifier from dice notation
func extractModifier(notation string) string {
	re := regexp.MustCompile(`[+-]\d+`)
	match := re.FindString(notation)
	return match
}

// isD20Roll checks if the dice notation uses d20 dice
func isD20Roll(notation string) bool {
	re := regexp.MustCompile(`^\d*d20([+-]\d+)?$`)
	return re.MatchString(notation)
}

// getDisplayValue returns the display value for a rollable tag.
// For d20 rolls: returns only the modifier (e.g., "+5" or "")
// For non-d20 rolls: returns average and notation (e.g., "8(1d8+5)")
func getDisplayValue(notation string) string {
	if isD20Roll(notation) {
		return extractModifier(notation)
	}
	avg := calculateAverage(notation)
	return fmt.Sprintf("%d(%s)", avg, notation)
}

// calculateAverage calculates the average result of a dice notation
// e.g., "2d6+3" -> (2 * 3.5) + 3 = 10
func calculateAverage(notation string) int {
	// Parse: count d sides +/- modifier
	re := regexp.MustCompile(`^(\d+)d(\d+)([+-]\d+)?$`)
	match := re.FindStringSubmatch(notation)

	if len(match) < 3 {
		return 0
	}

	count := 0
	sides := 0
	modifier := 0

	fmt.Sscanf(match[1], "%d", &count)
	fmt.Sscanf(match[2], "%d", &sides)

	if len(match) >= 4 && match[3] != "" {
		fmt.Sscanf(match[3], "%d", &modifier)
	}

	// Average of a die: (sides + 1) / 2.0
	// Total average: count * avg_per_die + modifier
	avgPerDie := float64(sides+1) / 2.0
	total := float64(count)*avgPerDie + float64(modifier)

	// Round to nearest integer
	return int(total + 0.5)
}
