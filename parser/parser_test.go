package parser

import (
	"strings"
	"testing"
)

func TestParseMarkdown_EmptyInput(t *testing.T) {
	input := ""
	result, err := ParseMarkdown(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Traits) != 0 || len(result.Actions) != 0 || len(result.BonusActions) != 0 || len(result.Reactions) != 0 {
		t.Error("Expected empty sections for empty input")
	}
}

func TestParseMarkdown_SingleTrait(t *testing.T) {
	input := `## Traits

**Darkvision.** You can see in dim light within 60 feet of you as if it were bright light.`

	result, err := ParseMarkdown(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Traits) != 1 {
		t.Fatalf("Expected 1 trait, got %d", len(result.Traits))
	}

	trait := result.Traits[0]
	if trait.Name != "Darkvision" {
		t.Errorf("Expected name 'Darkvision', got '%s'", trait.Name)
	}

	expectedDesc := "You can see in dim light within 60 feet of you as if it were bright light."
	if trait.Description != expectedDesc {
		t.Errorf("Expected description '%s', got '%s'", expectedDesc, trait.Description)
	}
}

func TestParseMarkdown_MultipleSections(t *testing.T) {
	input := `## Traits

**Pack Tactics.** You have advantage on attack rolls.

## Actions

**Bite.** Melee Weapon Attack: to hit: 1d20+4.

## Bonus Actions

**Dash.** You can take the Dash action.

## Reactions

**Parry.** Add 2 to your AC.`

	result, err := ParseMarkdown(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Traits) != 1 {
		t.Errorf("Expected 1 trait, got %d", len(result.Traits))
	}

	if len(result.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(result.Actions))
	}

	if len(result.BonusActions) != 1 {
		t.Errorf("Expected 1 bonus action, got %d", len(result.BonusActions))
	}

	if len(result.Reactions) != 1 {
		t.Errorf("Expected 1 reaction, got %d", len(result.Reactions))
	}
}

func TestParseMarkdown_MultipleAbilitiesInSection(t *testing.T) {
	input := `## Actions

**Bite.** Melee Weapon Attack: to hit: 1d20+4.

**Claw.** Melee Weapon Attack: to hit: 1d20+2.

**Tail.** Melee Weapon Attack: to hit: 1d20+3.`

	result, err := ParseMarkdown(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Actions) != 3 {
		t.Fatalf("Expected 3 actions, got %d", len(result.Actions))
	}

	expectedNames := []string{"Bite", "Claw", "Tail"}
	for i, action := range result.Actions {
		if action.Name != expectedNames[i] {
			t.Errorf("Expected action name '%s', got '%s'", expectedNames[i], action.Name)
		}
	}
}

func TestParseMarkdown_PreservesFormatting(t *testing.T) {
	input := `## Traits

**Spellcasting.** You can cast *fireball* and _magic missile_.`

	result, err := ParseMarkdown(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Traits) != 1 {
		t.Fatalf("Expected 1 trait, got %d", len(result.Traits))
	}

	// Should preserve italic markers
	if !strings.Contains(result.Traits[0].Description, "*fireball*") && !strings.Contains(result.Traits[0].Description, "_magic missile_") {
		t.Error("Expected formatting to be preserved")
	}
}

func TestParseMarkdown_IgnoresUnknownSections(t *testing.T) {
	input := `## Some Random Section

**Something.** This should be ignored.

## Actions

**Bite.** Melee Weapon Attack.`

	result, err := ParseMarkdown(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(result.Actions))
	}

	// Should not have anything from the unknown section
	if len(result.Traits) != 0 || len(result.BonusActions) != 0 || len(result.Reactions) != 0 {
		t.Error("Expected only Actions to be parsed")
	}
}

func TestParseMarkdown_HandlesEmptySections(t *testing.T) {
	input := `## Traits

## Actions

**Bite.** Melee Weapon Attack.`

	result, err := ParseMarkdown(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Traits) != 0 {
		t.Errorf("Expected 0 traits for empty section, got %d", len(result.Traits))
	}

	if len(result.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(result.Actions))
	}
}
