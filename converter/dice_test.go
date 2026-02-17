package converter

import (
	"strings"
	"testing"
)

func TestParseDiceNotation_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1d20", "1d20"},
		{"2d6", "2d6"},
		{"1d20+5", "1d20+5"},
		{"2d6-1", "2d6-1"},
		{"d20", "1d20"}, // implicit 1
		{"3d8+4", "3d8+4"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseDiceNotation(tt.input)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseDiceNotation_Invalid(t *testing.T) {
	tests := []string{
		"abc",
		"1d",
		"d",
		"1d3", // d3 not valid
		"1d100+", // incomplete modifier
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseDiceNotation(input)
			if err == nil {
				t.Errorf("Expected error for input %s, got nil", input)
			}
		})
	}
}

func TestConvertDiceRolls_Simple(t *testing.T) {
	input := "to hit: 1d20+5"
	actionName := "Longsword"

	result, err := ConvertDiceRolls(input, actionName)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := `[rollable]+5;{"diceNotation":"1d20+5","rollType":"to hit","rollAction":"Longsword"}[/rollable]`
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestConvertDiceRolls_Damage(t *testing.T) {
	input := "damage: 2d6+3"
	actionName := "Greatsword"

	result, err := ConvertDiceRolls(input, actionName)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := `[rollable]+3;{"diceNotation":"2d6+3","rollType":"damage","rollAction":"Greatsword"}[/rollable]`
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestConvertDiceRolls_NoModifier(t *testing.T) {
	input := "damage: 1d10"
	actionName := "Fire Bolt"

	result, err := ConvertDiceRolls(input, actionName)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// No modifier in display value
	expected := `[rollable];{"diceNotation":"1d10","rollType":"damage","rollAction":"Fire Bolt"}[/rollable]`
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestConvertDiceRolls_MultipleRolls(t *testing.T) {
	input := "Melee Attack: to hit: 1d20+4, reach 5 ft. Hit: damage: 2d6+2 slashing."
	actionName := "Greatsword"

	result, err := ConvertDiceRolls(input, actionName)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should contain both rolls
	if !strings.Contains(result, `"rollType":"to hit"`) {
		t.Error("Expected to hit roll in result")
	}
	if !strings.Contains(result, `"rollType":"damage"`) {
		t.Error("Expected damage roll in result")
	}
	if !strings.Contains(result, "Melee Attack:") {
		t.Error("Expected original text to be preserved")
	}
}

func TestConvertDiceRolls_Healing(t *testing.T) {
	input := "healing: 1d8+4"
	actionName := "Cure Wounds"

	result, err := ConvertDiceRolls(input, actionName)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, `"rollType":"healing"`) {
		t.Error("Expected healing roll type")
	}
}

func TestConvertDiceRolls_SaveDC(t *testing.T) {
	input := "save: DC 15"
	actionName := "Fireball"

	result, err := ConvertDiceRolls(input, actionName)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// DC shouldn't create a rollable, just preserve text
	if strings.Contains(result, "[rollable]") {
		t.Error("DC should not create rollable")
	}
}

func TestConvertDiceRolls_NoRolls(t *testing.T) {
	input := "You gain advantage on the next attack roll."
	actionName := "Help"

	result, err := ConvertDiceRolls(input, actionName)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != input {
		t.Errorf("Expected input unchanged when no rolls present")
	}
}

func TestExtractModifier_Positive(t *testing.T) {
	mod := extractModifier("1d20+5")
	if mod != "+5" {
		t.Errorf("Expected '+5', got '%s'", mod)
	}
}

func TestExtractModifier_Negative(t *testing.T) {
	mod := extractModifier("1d20-2")
	if mod != "-2" {
		t.Errorf("Expected '-2', got '%s'", mod)
	}
}

func TestExtractModifier_None(t *testing.T) {
	mod := extractModifier("1d20")
	if mod != "" {
		t.Errorf("Expected empty string, got '%s'", mod)
	}
}
