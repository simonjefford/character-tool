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
		"1d3",    // d3 not valid
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

	expected := `[rollable]10(2d6+3);{"diceNotation":"2d6+3","rollType":"damage","rollAction":"Greatsword"}[/rollable]`
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

	// Non-d20 should show average and full notation in parentheses
	expected := `[rollable]6(1d10);{"diceNotation":"1d10","rollType":"damage","rollAction":"Fire Bolt"}[/rollable]`
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

func TestCalculateAverage(t *testing.T) {
	tests := []struct {
		notation string
		expected int
	}{
		{"1d6", 4},    // 3.5 -> 4
		{"2d6", 7},    // 7
		{"1d8+3", 8},  // 4.5 + 3 = 7.5 -> 8
		{"3d6+2", 13}, // 10.5 + 2 = 12.5 -> 13
		{"1d10", 6},   // 5.5 -> 6
		{"2d8-1", 8},  // 9 - 1 = 8
		{"1d4+5", 8},  // 2.5 + 5 = 7.5 -> 8
		{"4d6", 14},   // 14
	}

	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			result := calculateAverage(tt.notation)
			if result != tt.expected {
				t.Errorf("For %s: expected %d, got %d", tt.notation, tt.expected, result)
			}
		})
	}
}

func TestGetDisplayValue_D20Rolls(t *testing.T) {
	tests := []struct {
		notation string
		expected string
	}{
		{"1d20+5", "+5"},
		{"1d20", ""},
		{"1d20-2", "-2"},
		{"d20+10", "+10"},
	}

	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			result := getDisplayValue(tt.notation)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetDisplayValue_NonD20Rolls(t *testing.T) {
	tests := []struct {
		notation string
		expected string
	}{
		{"1d10+5", "11(1d10+5)"}, // 5.5 + 5 = 10.5 -> 11
		{"2d6+3", "10(2d6+3)"},   // 2*3.5 + 3 = 10
		{"1d8", "5(1d8)"},        // 4.5 -> 5
		{"1d4-1", "2(1d4-1)"},    // 2.5 - 1 = 1.5 -> 2
		{"3d6", "11(3d6)"},       // 3*3.5 = 10.5 -> 11
		{"1d12+4", "11(1d12+4)"}, // 6.5 + 4 = 10.5 -> 11
	}

	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			result := getDisplayValue(tt.notation)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestConvertDiceRolls_VariousDiceTypes(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		actionName   string
		expectedDisp string
	}{
		{
			name:         "D20 attack roll",
			input:        "to hit: 1d20+5",
			actionName:   "Longsword",
			expectedDisp: "+5",
		},
		{
			name:         "D8 damage roll",
			input:        "damage: 1d8+3",
			actionName:   "Longsword",
			expectedDisp: "8(1d8+3)",
		},
		{
			name:         "D6 damage roll",
			input:        "damage: 2d6+2",
			actionName:   "Greatsword",
			expectedDisp: "9(2d6+2)",
		},
		{
			name:         "D10 damage no modifier",
			input:        "damage: 1d10",
			actionName:   "Fire Bolt",
			expectedDisp: "6(1d10)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertDiceRolls(tt.input, tt.actionName)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if !strings.Contains(result, "[rollable]"+tt.expectedDisp+";") {
				t.Errorf("Expected display value '%s' in result, got:\n%s", tt.expectedDisp, result)
			}
		})
	}
}
