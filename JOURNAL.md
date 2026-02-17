# Development Journal

## [2026-02-17] Project Initialization

### Changes Made
- Created Go module with `go mod init character-tool`
- Established project directory structure:
  - `parser/` - Markdown parsing logic
  - `converter/` - Dice roll and spell link conversion
  - `formatter/` - Output formatting
  - `data/` - Spell list data
  - `testdata/` - Example markdown files and test data
- Created initial documentation files:
  - `PLAN.md` - Implementation plan copied from design phase
  - `JOURNAL.md` - This development journal
  - `CLAUDE.md` - Development rules and guidelines

### Design Decisions
- **Module Structure**: Separated concerns into distinct packages (parser, converter, formatter) for maintainability and testability
- **Documentation First**: Creating documentation files before code to establish clear guidelines and reference material
- **TDD Approach**: Will follow test-driven development for all implementation phases

### Next Steps
- Create CLAUDE.md with development rules
- Create README.md with usage instructions
- Begin implementing parser module with TDD approach

## [2026-02-17] Parser Module Implementation

### Tests Written
- `TestParseMarkdown_EmptyInput` - Verifies empty input returns empty sections
- `TestParseMarkdown_SingleTrait` - Tests parsing a single trait with name and description
- `TestParseMarkdown_MultipleSections` - Tests parsing all four section types (Traits, Actions, Bonus Actions, Reactions)
- `TestParseMarkdown_MultipleAbilitiesInSection` - Tests multiple abilities in one section
- `TestParseMarkdown_PreservesFormatting` - Ensures italic and other formatting is preserved
- `TestParseMarkdown_IgnoresUnknownSections` - Verifies unknown sections are skipped
- `TestParseMarkdown_HandlesEmptySections` - Tests sections with no content

### Implementation
- Created `parser/parser.go` with:
  - `Ability` struct to represent individual abilities
  - `AbilityType` enum for Trait, Action, BonusAction, Reaction
  - `ParseResult` struct to organize parsed abilities
  - `ParseMarkdown()` function as main entry point
  - `splitBySections()` to extract markdown sections by ## headers
  - `getSectionType()` to map section names to types
  - `parseAbilities()` to extract **Name.** Description patterns

### Design Decisions
- **Paragraph-based splitting**: Split section content by `\n\n` to separate individual abilities, which is more reliable than complex regex with lookahead
- **Regex simplicity**: Used simple regex pattern `^\*\*([^*]+)\.\*\*\s*(.+)$` to match bold ability names followed by periods
- **Section mapping**: Case-insensitive matching for section names to handle "Actions" vs "actions"
- **Unknown sections**: Silently ignore sections that don't match the four expected types

### Issues Encountered
- **Initial regex issue**: First attempt used negative lookahead `(?!)` which isn't supported in Go's RE2 regex engine
- **Solution**: Simplified approach by splitting on paragraph breaks first, then applying regex to each paragraph

### Test Results
All 7 tests pass successfully.

### Next Steps
- Implement dice roll converter module with TDD

## [2026-02-17] Dice Roll Converter Implementation

### Tests Written
- `TestParseDiceNotation_Valid` - Tests valid dice notation parsing (1d20, 2d6+5, d20, etc.)
- `TestParseDiceNotation_Invalid` - Tests rejection of invalid notation (d3, incomplete, etc.)
- `TestConvertDiceRolls_Simple` - Tests basic "to hit: 1d20+5" conversion
- `TestConvertDiceRolls_Damage` - Tests damage roll conversion
- `TestConvertDiceRolls_NoModifier` - Tests rolls without modifiers (1d10)
- `TestConvertDiceRolls_MultipleRolls` - Tests text with multiple roll types
- `TestConvertDiceRolls_Healing` - Tests healing roll type
- `TestConvertDiceRolls_SaveDC` - Tests that DC values don't create rollables
- `TestConvertDiceRolls_NoRolls` - Tests text without any rolls
- `TestExtractModifier_*` - Tests modifier extraction (+5, -2, none)

### Implementation
- Created `converter/dice.go` with:
  - `ParseDiceNotation()` - Validates and normalizes dice notation
  - `ConvertDiceRolls()` - Main conversion function for text with roll keywords
  - `extractModifier()` - Helper to extract +X or -X modifiers
  - `RollableData` struct - JSON structure for rollable metadata
  - `validDice` map - Validates only d4, d6, d8, d10, d12, d20, d100

### Design Decisions
- **Regex pattern**: `(to hit|damage|healing|save):\s*(\d*d\d+[+-]?\d*)` to match roll keywords
- **Rollable format**: `[rollable]MODIFIER;JSON[/rollable]` where modifier is displayed and JSON contains full data
- **Validation**: Strictly validate dice types against D&D 5e standard dice
- **Normalization**: Convert implicit "d20" to "1d20" for consistency
- **DC handling**: Don't convert DC values to rollables (they're static numbers)
- **Error handling**: Return original text if dice notation is invalid

### Test Results
All 14 tests pass successfully.

### Next Steps
- Add D&D 5e spell list data
- Implement spell linker module with TDD

## [2026-02-17] Spell Linker Implementation

### Tests Written
- `TestLoadSpells` - Verifies spell list loads from JSON file
- `TestIsValidSpell` - Tests case-insensitive spell validation (8 cases)
- `TestConvertSpellLinks_Single` - Tests single spell conversion
- `TestConvertSpellLinks_Multiple` - Tests multiple spells in one text
- `TestConvertSpellLinks_InvalidSpell` - Tests warning for unknown spells
- `TestConvertSpellLinks_MixedValidInvalid` - Tests mix of valid and invalid spells
- `TestConvertSpellLinks_CaseInsensitive` - Tests case-insensitive matching
- `TestConvertSpellLinks_NoSpells` - Tests text without spell references
- `TestConvertSpellLinks_EmptySpellName` - Tests empty spell name handling

### Implementation
- Created `converter/spell.go` with:
  - `LoadSpells()` - Loads spell list from data/spells.json using runtime.Caller for path resolution
  - `IsValidSpell()` - Case-insensitive spell name validation
  - `ConvertSpellLinks()` - Converts {{spell:Name}} to [spell]Name[/spell] format
  - Returns warnings for unknown spells while still converting them

### Design Decisions
- **File loading**: Used `runtime.Caller(0)` to get source file location and navigate to data directory
- **Case-insensitive map**: Store all spell names as lowercase keys for O(1) lookup
- **Preserve original case**: Output uses the original case from input, not the canonical spell name
- **Regex pattern**: `\{\{spell:([^}]*)\}\}` to match spell syntax
- **Warning system**: Return list of warnings for invalid spells rather than failing
- **Empty names**: Warn about empty spell names but still convert

### Issues Encountered
- **Embed path issue**: Cannot use `//go:embed ../data/spells.json` with relative parent paths
- **Solution**: Use `runtime.Caller()` to dynamically resolve path to data directory

### Test Results
All 24 converter tests pass (14 dice + 10 spell).

### Next Steps
- Implement formatter module to combine all conversions

## [2026-02-17] Formatter Module Implementation

### Tests Written
- `TestFormatAbilities_EmptyList` - Tests empty input
- `TestFormatAbilities_SingleAbility` - Tests single ability formatting
- `TestFormatAbilities_MultipleAbilities` - Tests multiple abilities with blank line separation
- `TestFormatAbilities_WithDiceRolls` - Tests dice roll integration
- `TestFormatAbilities_WithSpellLinks` - Tests spell link integration
- `TestFormatAbilities_WithInvalidSpell` - Tests warning collection for invalid spells
- `TestFormatAbilities_WithDiceAndSpells` - Tests combining both conversions

### Implementation
- Created `formatter/formatter.go` with:
  - `FormatAbilities()` - Main formatting function that orchestrates all conversions
  - Converts spell links using converter.ConvertSpellLinks()
  - Converts dice rolls using converter.ConvertDiceRolls()
  - Formats as "Name. Description"
  - Joins multiple abilities with blank lines
  - Collects warnings from all conversions

### Design Decisions
- **Conversion order**: Apply spell links first, then dice rolls (order doesn't matter as they operate on different patterns)
- **Action name**: Use ability name as the action name for dice rolls
- **Warning aggregation**: Collect all warnings from both spell and dice conversions
- **Formatting**: Simple "Name. Description" format, blank line between abilities
- **Error handling**: Propagate errors from converters but continue processing

### Test Results
All 7 formatter tests pass.

### Next Steps
- Implement CLI with input/output handling
- Create example test file in testdata/
