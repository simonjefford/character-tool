package parser

import (
	"regexp"
	"strings"
)

// AbilityType represents the type of ability
type AbilityType int

const (
	Trait AbilityType = iota
	Action
	BonusAction
	Reaction
)

// Ability represents a character ability, trait, action, etc.
type Ability struct {
	Name        string
	Description string
	Type        AbilityType
}

// ParseResult contains all parsed abilities organized by type
type ParseResult struct {
	Traits       []Ability
	Actions      []Ability
	BonusActions []Ability
	Reactions    []Ability
}

// ParseMarkdown parses a markdown string and extracts character abilities
func ParseMarkdown(markdown string) (*ParseResult, error) {
	result := &ParseResult{
		Traits:       []Ability{},
		Actions:      []Ability{},
		BonusActions: []Ability{},
		Reactions:    []Ability{},
	}

	if strings.TrimSpace(markdown) == "" {
		return result, nil
	}

	// Split by ## headers
	sections := splitBySections(markdown)

	for sectionName, content := range sections {
		abilityType, ok := getSectionType(sectionName)
		if !ok {
			// Skip unknown sections
			continue
		}

		abilities := parseAbilities(content, abilityType)

		switch abilityType {
		case Trait:
			result.Traits = append(result.Traits, abilities...)
		case Action:
			result.Actions = append(result.Actions, abilities...)
		case BonusAction:
			result.BonusActions = append(result.BonusActions, abilities...)
		case Reaction:
			result.Reactions = append(result.Reactions, abilities...)
		}
	}

	return result, nil
}

// splitBySections splits markdown by ## headers and returns a map of section name to content
func splitBySections(markdown string) map[string]string {
	sections := make(map[string]string)

	// Regex to match ## Header
	headerRegex := regexp.MustCompile(`(?m)^## (.+)$`)
	matches := headerRegex.FindAllStringSubmatchIndex(markdown, -1)

	if len(matches) == 0 {
		return sections
	}

	for i, match := range matches {
		headerEnd := match[1]
		nameStart := match[2]
		nameEnd := match[3]

		sectionName := markdown[nameStart:nameEnd]

		// Content starts after the header line
		contentStart := headerEnd
		if contentStart < len(markdown) && markdown[contentStart] == '\n' {
			contentStart++
		}

		// Content ends at the next header or end of string
		var contentEnd int
		if i < len(matches)-1 {
			contentEnd = matches[i+1][0]
		} else {
			contentEnd = len(markdown)
		}

		content := strings.TrimSpace(markdown[contentStart:contentEnd])
		sections[strings.TrimSpace(sectionName)] = content
	}

	return sections
}

// getSectionType maps section names to AbilityType
func getSectionType(sectionName string) (AbilityType, bool) {
	normalized := strings.ToLower(strings.TrimSpace(sectionName))

	switch normalized {
	case "traits":
		return Trait, true
	case "actions":
		return Action, true
	case "bonus actions":
		return BonusAction, true
	case "reactions":
		return Reaction, true
	default:
		return 0, false
	}
}

// parseAbilities extracts individual abilities from section content
func parseAbilities(content string, abilityType AbilityType) []Ability {
	abilities := []Ability{}

	if strings.TrimSpace(content) == "" {
		return abilities
	}

	// Split by paragraph breaks to separate abilities
	for paragraph := range strings.SplitSeq(content, "\n\n") {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		// Match **Name.** Description or **Name**. Description pattern
		// Supports period inside or outside bold markers
		abilityRegex := regexp.MustCompile(`^\*\*([^*]+?)\.?\*\*\.?\s*(.+)$`)
		match := abilityRegex.FindStringSubmatch(paragraph)

		if len(match) >= 3 {
			name := strings.TrimSpace(match[1])
			description := strings.TrimSpace(match[2])

			abilities = append(abilities, Ability{
				Name:        name,
				Description: description,
				Type:        abilityType,
			})
		}
	}

	return abilities
}
