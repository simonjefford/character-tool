# D&D Beyond Character Tool - Implementation Plan

## Overview
Create a Go CLI tool that converts markdown character descriptions into D&D Beyond-formatted blocks with spell links and rollable dice notation.

## Requirements Summary
- **Input**: Markdown file with structured headers (## Traits, ## Actions, ## Bonus Actions, ## Reactions)
- **Spell Links**: Explicit markup using `{{spell:SpellName}}` syntax
- **Dice Rolls**: Inline notation (e.g., `1d20+2`) with keyword prefixes (e.g., `to hit: 1d20+2`, `damage: 2d6+3`)
- **Output**: Separate files per section (traits.txt, actions.txt, bonus-actions.txt, reactions.txt)
- **Validation**: Validate spell names and dice notation, warn about issues

## Architecture

### Core Components

1. **Parser Module** (`parser/parser.go`)
   - Parse markdown file using structured headers
   - Extract sections: Traits, Actions, Bonus Actions, Reactions
   - Parse individual ability entries within each section

2. **Spell Linker** (`converter/spell.go`)
   - Detect `{{spell:SpellName}}` markup
   - Convert to D&D Beyond spell link format
   - Validate against D&D 5e spell list
   - Warn about unknown spells

3. **Dice Roller Converter** (`converter/dice.go`)
   - Parse inline dice notation (XdY+Z format)
   - Detect roll type keywords: "to hit:", "damage:", "DC", "save:", etc.
   - Extract action name from surrounding context
   - Generate D&D Beyond rollable format: `[rollable]+X;{"diceNotation":"...","rollType":"...","rollAction":"..."}[/rollable]`
   - Validate dice notation format

4. **Formatter** (`formatter/formatter.go`)
   - Combine parsed abilities with converted spell links and dice rolls
   - Format according to D&D Beyond text conventions
   - Generate separate output blocks for each section

5. **CLI** (`main.go`)
   - Accept input markdown file path
   - Accept output directory path (default: current directory)
   - Display validation warnings
   - Write separate output files

### Data Structures

```go
// Ability represents a trait, action, bonus action, or reaction
type Ability struct {
    Name        string
    Description string
    Type        AbilityType // Trait, Action, BonusAction, Reaction
}

// RollableInfo represents a dice roll with metadata
type RollableInfo struct {
    DiceNotation string
    Modifier     int
    RollType     string // "to hit", "damage", "save DC", "check"
    RollAction   string // ability/weapon name
}

// SpellReference represents a spell link
type SpellReference struct {
    SpellName string
    IsValid   bool
}
```

## Implementation Steps

### Phase 1: Project Setup
1. Initialize Go module (`go mod init character-tool`)
2. Create project documentation:
   - `PLAN.md` - Copy of this implementation plan
   - `JOURNAL.md` - Development journal (start with project initialization entry)
   - `CLAUDE.md` - Development rules (TDD, journal updates, conventional commits)
3. Create directory structure:
   - `parser/`
   - `converter/`
   - `formatter/`
   - `data/` (for spell list)
   - `testdata/` (for example markdown files)
4. Add dependencies:
   - `github.com/gomarkdown/markdown` for markdown parsing
   - `github.com/spf13/cobra` for CLI (optional, can use stdlib)

### Phase 2: Parser Implementation
1. Create markdown parser that:
   - Reads file content
   - Splits by `##` headers
   - Identifies section types (Traits, Actions, Bonus Actions, Reactions)
   - Extracts ability entries (typically `**Name:** description` or `### Name` format)
2. Handle multiple abilities per section
3. Preserve formatting (bold, italic, lists)

### Phase 3: Dice Roll Converter
1. Create regex patterns for:
   - Dice notation: `(\d+)?d(\d+)([+-]\d+)?`
   - Roll type keywords: `(to hit|damage|DC|save):\s*`
2. Implement context extraction:
   - Extract ability name from heading or bold text
   - Match roll type keyword before dice notation
3. Generate D&D Beyond rollable format
4. Validate dice notation (d4, d6, d8, d10, d12, d20, d100)

### Phase 4: Spell Linker
1. Load D&D 5e spell list (embedded data or JSON file)
2. Parse `{{spell:SpellName}}` syntax
3. Convert to D&D Beyond spell link format (research exact format)
4. Validate spell names against list
5. Generate warnings for unknown spells

### Phase 5: Formatter
1. Combine parsed content with converted links/rolls
2. Apply D&D Beyond text formatting conventions
3. Generate clean output blocks
4. Handle edge cases (empty sections, malformed input)

### Phase 6: CLI & Output
1. Implement main CLI:
   - Accept `-input` flag for markdown file
   - Accept `-output` flag for output directory
   - Accept `-verbose` flag for detailed warnings
2. Write separate output files:
   - `traits.txt`
   - `actions.txt`
   - `bonus-actions.txt`
   - `reactions.txt`
3. Display summary and warnings

### Phase 7: Testing & Examples
1. Create example markdown file in `testdata/`
2. Add unit tests for each component
3. Test with various edge cases:
   - Missing sections
   - Invalid dice notation
   - Unknown spells
   - Multiple rolls in one ability
4. Document expected input format in README

## Critical Files

**New Files to Create:**
- `main.go` - CLI entry point
- `parser/parser.go` - Markdown parsing
- `converter/dice.go` - Dice roll conversion
- `converter/spell.go` - Spell link conversion
- `formatter/formatter.go` - Output formatting
- `data/spells.json` - D&D 5e spell list
- `testdata/example-character.md` - Example input
- `go.mod` - Go module definition
- `README.md` - Usage documentation
- `PLAN.md` - Copy of this implementation plan
- `JOURNAL.md` - Development journal tracking all changes
- `CLAUDE.md` - Development rules and guidelines

## Example Input Format

```markdown
## Traits

**Spellcasting.** You can cast spells using {{spell:Fireball}} and {{spell:Magic Missile}}.

**Pack Tactics.** You have advantage on attack rolls against a creature if at least one ally is within 5 feet.

## Actions

**Quarterstaff.** Melee Weapon Attack: to hit: 1d20+2, reach 5 ft., one target. Hit: damage: 1d6+2 bludgeoning damage.

**Fire Bolt.** Ranged Spell Attack: to hit: 1d20+5, range 120 ft., one target. Hit: damage: 1d10 fire damage.

## Bonus Actions

**Second Wind.** Regain healing: 1d10+5 hit points.

## Reactions

**Shield.** Cast {{spell:Shield}} when hit by an attack, gaining +5 AC.
```

## Example Output Format

**traits.txt:**
```
Spellcasting. You can cast spells using [spell]Fireball[/spell] and [spell]Magic Missile[/spell].

Pack Tactics. You have advantage on attack rolls against a creature if at least one ally is within 5 feet.
```

**actions.txt:**
```
Quarterstaff. Melee Weapon Attack: [rollable]+2;{"diceNotation":"1d20+2","rollType":"to hit","rollAction":"Quarterstaff"}[/rollable], reach 5 ft., one target. Hit: [rollable]+2;{"diceNotation":"1d6+2","rollType":"damage","rollAction":"Quarterstaff"}[/rollable] bludgeoning damage.

Fire Bolt. Ranged Spell Attack: [rollable]+5;{"diceNotation":"1d20+5","rollType":"to hit","rollAction":"Fire Bolt"}[/rollable], range 120 ft., one target. Hit: [rollable];{"diceNotation":"1d10","rollType":"damage","rollAction":"Fire Bolt"}[/rollable] fire damage.
```

## Development Process

### Test-Driven Development (TDD)
All implementation must follow TDD:
1. Write failing test first
2. Implement minimal code to pass test
3. Refactor while keeping tests green
4. Document in JOURNAL.md

### Journal Updates
Every change must be recorded in JOURNAL.md with:
- Date and description of change
- Rationale for design decisions
- Tests written and implementation approach
- Any issues encountered and solutions

### Commit Messages
All commits must use conventional commit format:
- `feat: description` - New features
- `fix: description` - Bug fixes
- `test: description` - Test additions/changes
- `refactor: description` - Code refactoring
- `docs: description` - Documentation updates
- `chore: description` - Build/tooling changes

## Verification Plan

1. **Unit Tests**: Test each component (parser, dice converter, spell linker)
2. **Integration Test**: Run tool on `testdata/example-character.md`
3. **Manual Testing**:
   - Copy output blocks into D&D Beyond character sheet
   - Verify spell links are clickable
   - Verify dice rolls work correctly with proper labels
   - Test edge cases (missing sections, invalid syntax)
4. **Validation Testing**:
   - Verify warnings for unknown spells
   - Verify warnings for invalid dice notation
   - Verify graceful handling of malformed input

## Technical Decisions

1. **Markdown Parsing**: Use `github.com/gomarkdown/markdown` library for robust parsing
2. **Spell List**: Embed spell list as JSON file in `data/` directory
3. **Regex for Dice**: Simple regex sufficient for standard dice notation
4. **Error Handling**: Warn but don't fail on validation errors (best-effort conversion)
5. **CLI Framework**: Use stdlib `flag` package for simplicity (can upgrade to cobra later)
