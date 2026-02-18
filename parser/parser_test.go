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

func TestParseMarkdown_PeriodPlacement(t *testing.T) {
	input := `## Bonus Actions

**Second Wind.** Regain healing: 1d10+5 hit points.

**Fireball**. Kill all the things: 20d20 hit points.`

	result, err := ParseMarkdown(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.BonusActions) != 2 {
		t.Fatalf("Expected 2 bonus actions, got %d", len(result.BonusActions))
	}

	// Test period inside bold markers
	if result.BonusActions[0].Name != "Second Wind" {
		t.Errorf("Expected name 'Second Wind', got '%s'", result.BonusActions[0].Name)
	}

	// Test period outside bold markers
	if result.BonusActions[1].Name != "Fireball" {
		t.Errorf("Expected name 'Fireball', got '%s'", result.BonusActions[1].Name)
	}
}

func TestParseMarkdown_PlainTextParagraphs(t *testing.T) {
	input := `## Traits

This character has a mysterious background that grants them unique abilities.

**Darkvision.** You can see in dim light within 60 feet.

Due to their elven heritage, they also gain the following benefits.

**Fey Ancestry.** You have advantage on saving throws against being charmed.`

	result, err := ParseMarkdown(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Traits) != 4 {
		t.Fatalf("Expected 4 traits (2 named + 2 plain text), got %d", len(result.Traits))
	}

	// First should be plain text
	if result.Traits[0].Name != "" {
		t.Errorf("Expected first trait to have no name, got '%s'", result.Traits[0].Name)
	}
	if !strings.Contains(result.Traits[0].Description, "mysterious background") {
		t.Error("Expected first trait to contain plain text description")
	}

	// Second should be named ability
	if result.Traits[1].Name != "Darkvision" {
		t.Errorf("Expected second trait name 'Darkvision', got '%s'", result.Traits[1].Name)
	}

	// Third should be plain text
	if result.Traits[2].Name != "" {
		t.Errorf("Expected third trait to have no name, got '%s'", result.Traits[2].Name)
	}
	if !strings.Contains(result.Traits[2].Description, "elven heritage") {
		t.Error("Expected third trait to contain plain text description")
	}

	// Fourth should be named ability
	if result.Traits[3].Name != "Fey Ancestry" {
		t.Errorf("Expected fourth trait name 'Fey Ancestry', got '%s'", result.Traits[3].Name)
	}
}
