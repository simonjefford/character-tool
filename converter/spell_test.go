package converter

import (
	"strings"
	"testing"
)

func TestLoadSpells(t *testing.T) {
	spells, err := LoadSpells()
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	if len(spells) == 0 {
		t.Error("Expected spells to be loaded, got empty list")
	}

	// Check for some well-known spells (map uses lowercase keys)
	expectedSpells := []string{"fireball", "magic missile", "shield", "cure wounds"}
	for _, spell := range expectedSpells {
		if !spells[spell] {
			t.Errorf("Expected spell '%s' to be in list", spell)
		}
	}
}

func TestIsValidSpell(t *testing.T) {
	spells, err := LoadSpells()
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	tests := []struct {
		spell string
		valid bool
	}{
		{"Fireball", true},
		{"Magic Missile", true},
		{"Shield", true},
		{"fireball", true},              // case insensitive
		{"MAGIC MISSILE", true},         // case insensitive
		{"NotASpell", false},
		{"Random Spell Name", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.spell, func(t *testing.T) {
			result := IsValidSpell(tt.spell, spells)
			if result != tt.valid {
				t.Errorf("IsValidSpell(%q) = %v, want %v", tt.spell, result, tt.valid)
			}
		})
	}
}

func TestConvertSpellLinks_Single(t *testing.T) {
	spells, _ := LoadSpells()
	input := "You can cast {{spell:Fireball}}."

	result, warnings := ConvertSpellLinks(input, spells)

	expected := "You can cast [spell]Fireball[/spell]."
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}

func TestConvertSpellLinks_Multiple(t *testing.T) {
	spells, _ := LoadSpells()
	input := "Cast {{spell:Fireball}} or {{spell:Magic Missile}}."

	result, warnings := ConvertSpellLinks(input, spells)

	if !strings.Contains(result, "[spell]Fireball[/spell]") {
		t.Error("Expected Fireball to be converted")
	}

	if !strings.Contains(result, "[spell]Magic Missile[/spell]") {
		t.Error("Expected Magic Missile to be converted")
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}

func TestConvertSpellLinks_InvalidSpell(t *testing.T) {
	spells, _ := LoadSpells()
	input := "You can cast {{spell:NotASpell}}."

	result, warnings := ConvertSpellLinks(input, spells)

	// Should still convert but with warning
	expected := "You can cast [spell]NotASpell[/spell]."
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}

	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got %d", len(warnings))
	}

	if !strings.Contains(warnings[0], "NotASpell") {
		t.Errorf("Expected warning about NotASpell, got %s", warnings[0])
	}
}

func TestConvertSpellLinks_MixedValidInvalid(t *testing.T) {
	spells, _ := LoadSpells()
	input := "Cast {{spell:Fireball}} or {{spell:FakeSpell}} or {{spell:Shield}}."

	result, warnings := ConvertSpellLinks(input, spells)

	// All should be converted
	if !strings.Contains(result, "[spell]Fireball[/spell]") {
		t.Error("Expected Fireball to be converted")
	}
	if !strings.Contains(result, "[spell]FakeSpell[/spell]") {
		t.Error("Expected FakeSpell to be converted")
	}
	if !strings.Contains(result, "[spell]Shield[/spell]") {
		t.Error("Expected Shield to be converted")
	}

	// Should have warning only for FakeSpell
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got %d: %v", len(warnings), warnings)
	}

	if !strings.Contains(warnings[0], "FakeSpell") {
		t.Errorf("Expected warning about FakeSpell, got %s", warnings[0])
	}
}

func TestConvertSpellLinks_CaseInsensitive(t *testing.T) {
	spells, _ := LoadSpells()
	input := "Cast {{spell:fireball}} and {{spell:MAGIC MISSILE}}."

	result, warnings := ConvertSpellLinks(input, spells)

	// Should preserve original case in output
	if !strings.Contains(result, "[spell]fireball[/spell]") {
		t.Error("Expected fireball (lowercase) to be converted")
	}
	if !strings.Contains(result, "[spell]MAGIC MISSILE[/spell]") {
		t.Error("Expected MAGIC MISSILE (uppercase) to be converted")
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings for valid spells (case insensitive), got %v", warnings)
	}
}

func TestConvertSpellLinks_NoSpells(t *testing.T) {
	spells, _ := LoadSpells()
	input := "This text has no spell references."

	result, warnings := ConvertSpellLinks(input, spells)

	if result != input {
		t.Errorf("Expected input unchanged, got %s", result)
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}

func TestConvertSpellLinks_EmptySpellName(t *testing.T) {
	spells, _ := LoadSpells()
	input := "Cast {{spell:}}."

	result, warnings := ConvertSpellLinks(input, spells)

	// Should convert empty name but warn
	expected := "Cast [spell][/spell]."
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}

	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got %d", len(warnings))
	}
}
