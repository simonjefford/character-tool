package formatter

import (
	"character-tool/parser"
	"strings"
	"testing"
)

func TestFormatAbilities_EmptyList(t *testing.T) {
	abilities := []parser.Ability{}
	spells := make(map[string]bool)

	result, warnings, err := FormatAbilities(abilities, spells)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}

func TestFormatAbilities_SingleAbility(t *testing.T) {
	abilities := []parser.Ability{
		{
			Name:        "Darkvision",
			Description: "You can see in dim light within 60 feet.",
			Type:        parser.Trait,
		},
	}
	spells := make(map[string]bool)

	result, warnings, err := FormatAbilities(abilities, spells)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "Darkvision. You can see in dim light within 60 feet."
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}

func TestFormatAbilities_MultipleAbilities(t *testing.T) {
	abilities := []parser.Ability{
		{
			Name:        "Darkvision",
			Description: "You can see in dim light within 60 feet.",
			Type:        parser.Trait,
		},
		{
			Name:        "Pack Tactics",
			Description: "You have advantage on attack rolls.",
			Type:        parser.Trait,
		},
	}
	spells := make(map[string]bool)

	result, warnings, err := FormatAbilities(abilities, spells)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should have both abilities separated by blank line
	lines := strings.Split(result, "\n\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 abilities separated by blank lines, got %d: %v", len(lines), lines)
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}

func TestFormatAbilities_WithDiceRolls(t *testing.T) {
	abilities := []parser.Ability{
		{
			Name:        "Quarterstaff",
			Description: "Melee Weapon Attack: to hit: 1d20+2, reach 5 ft. Hit: damage: 1d6+2 bludgeoning.",
			Type:        parser.Action,
		},
	}
	spells := make(map[string]bool)

	result, warnings, err := FormatAbilities(abilities, spells)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should contain rollable tags
	if !strings.Contains(result, "[rollable]") {
		t.Error("Expected rollable tags in output")
	}

	if !strings.Contains(result, `"rollType":"to hit"`) {
		t.Error("Expected to hit roll in output")
	}

	if !strings.Contains(result, `"rollType":"damage"`) {
		t.Error("Expected damage roll in output")
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}

func TestFormatAbilities_WithSpellLinks(t *testing.T) {
	abilities := []parser.Ability{
		{
			Name:        "Spellcasting",
			Description: "You can cast {{spell:Fireball}} and {{spell:Shield}}.",
			Type:        parser.Trait,
		},
	}
	spells := map[string]bool{
		"fireball": true,
		"shield":   true,
	}

	result, warnings, err := FormatAbilities(abilities, spells)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should contain spell tags
	if !strings.Contains(result, "[spell]Fireball[/spell]") {
		t.Error("Expected Fireball spell tag in output")
	}

	if !strings.Contains(result, "[spell]Shield[/spell]") {
		t.Error("Expected Shield spell tag in output")
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings for valid spells, got %v", warnings)
	}
}

func TestFormatAbilities_WithInvalidSpell(t *testing.T) {
	abilities := []parser.Ability{
		{
			Name:        "Spellcasting",
			Description: "You can cast {{spell:NotASpell}}.",
			Type:        parser.Trait,
		},
	}
	spells := make(map[string]bool)

	result, warnings, err := FormatAbilities(abilities, spells)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should still convert spell
	if !strings.Contains(result, "[spell]NotASpell[/spell]") {
		t.Error("Expected NotASpell to be converted")
	}

	// Should have warning
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got %d: %v", len(warnings), warnings)
	}

	if !strings.Contains(warnings[0], "NotASpell") {
		t.Errorf("Expected warning about NotASpell, got %s", warnings[0])
	}
}

func TestFormatAbilities_WithDiceAndSpells(t *testing.T) {
	abilities := []parser.Ability{
		{
			Name:        "Magic Attack",
			Description: "Cast {{spell:Fire Bolt}} for to hit: 1d20+5, dealing damage: 1d10 fire damage.",
			Type:        parser.Action,
		},
	}
	spells := map[string]bool{
		"fire bolt": true,
	}

	result, warnings, err := FormatAbilities(abilities, spells)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should have both spell links and rollable dice
	if !strings.Contains(result, "[spell]Fire Bolt[/spell]") {
		t.Error("Expected Fire Bolt spell tag")
	}

	if !strings.Contains(result, "[rollable]") {
		t.Error("Expected rollable dice tags")
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}

func TestFormatAbilities_D20VsDamageDisplay(t *testing.T) {
	abilities := []parser.Ability{
		{
			Name:        "Longsword",
			Description: "Melee Weapon Attack: to hit: 1d20+5, reach 5 ft., one target. Hit: damage: 1d8+3 slashing.",
			Type:        parser.Action,
		},
		{
			Name:        "Dagger",
			Description: "Melee or Ranged Weapon Attack: to hit: 1d20+3, reach 5 ft. Hit: damage: 1d4+3 piercing.",
			Type:        parser.Action,
		},
	}
	spells := make(map[string]bool)

	result, warnings, err := FormatAbilities(abilities, spells)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// D20 rolls should show only modifier
	if !strings.Contains(result, "[rollable]+5;") {
		t.Error("Expected d20 roll to show '+5' modifier only")
	}
	if !strings.Contains(result, "[rollable]+3;") {
		t.Error("Expected d20 roll to show '+3' modifier only")
	}

	// Non-d20 damage rolls should show average and full notation in parentheses
	if !strings.Contains(result, "[rollable]8(1d8+3);") {
		t.Error("Expected damage roll to show '8(1d8+3)' with average")
	}
	if !strings.Contains(result, "[rollable]6(1d4+3);") {
		t.Error("Expected damage roll to show '6(1d4+3)' with average")
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}

func TestFormatAbilities_PlainTextParagraphs(t *testing.T) {
	abilities := []parser.Ability{
		{
			Name:        "",
			Description: "This character has enhanced abilities due to their training.",
			Type:        parser.Trait,
		},
		{
			Name:        "Enhanced Reflexes",
			Description: "You gain a +2 bonus to AC.",
			Type:        parser.Trait,
		},
		{
			Name:        "",
			Description: "The following abilities are granted by their magical armor.",
			Type:        parser.Trait,
		},
	}
	spells := make(map[string]bool)

	result, warnings, err := FormatAbilities(abilities, spells)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Plain text should not have "Name. " prefix
	lines := strings.Split(result, "\n\n")
	if len(lines) != 3 {
		t.Fatalf("Expected 3 paragraphs, got %d", len(lines))
	}

	// First paragraph - plain text
	if !strings.Contains(lines[0], "enhanced abilities") {
		t.Error("Expected first paragraph to contain plain text")
	}
	if strings.HasPrefix(lines[0], ". ") {
		t.Error("Plain text should not start with '. '")
	}

	// Second paragraph - named ability
	if !strings.HasPrefix(lines[1], "Enhanced Reflexes. ") {
		t.Error("Named ability should start with 'Enhanced Reflexes. '")
	}

	// Third paragraph - plain text
	if !strings.Contains(lines[2], "magical armor") {
		t.Error("Expected third paragraph to contain plain text")
	}
	if strings.HasPrefix(lines[2], ". ") {
		t.Error("Plain text should not start with '. '")
	}

	if len(warnings) != 0 {
		t.Errorf("Expected no warnings, got %v", warnings)
	}
}
