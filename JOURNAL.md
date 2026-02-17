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
