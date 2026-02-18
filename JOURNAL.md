# Development Journal

## [2026-02-18] Update Clipboard Script to Use Compiled Binary

### Description
Updated `ddb-copy.sh` to use the compiled `character-tool` binary instead of `go run`. Added `install.sh` script to install both the binary and script to `~/bin` for system-wide access.

### Changes
1. **Updated `ddb-copy.sh`:**
   - Removed `go run main.go` in favor of compiled binary
   - Removed directory check for main.go
   - Added logic to find binary in script directory or PATH
   - Provides helpful error message if binary not found

2. **Created `install.sh`:**
   - Builds the character-tool binary
   - Copies both `character-tool` and `ddb-copy.sh` to `~/bin`
   - Checks if `~/bin` is in PATH and provides instructions if not
   - Makes scripts executable
   - Color-coded output with status messages

3. **Updated README.md:**
   - Added "Quick Install (Recommended)" section
   - Updated clipboard workflow examples to use installed commands
   - Documented where script looks for binary

### Design Decisions
- **Compiled binary**: Much faster startup than `go run` (no compilation step)
- **~/bin location**: Standard location for user binaries on Unix systems
- **Flexible binary lookup**: Checks script directory first, then PATH for flexibility
- **Install script**: Automates setup process and provides PATH instructions

### Testing
Tested:
- Install script from project directory
- Running `ddb-copy.sh` from various directories with binary in PATH
- Error message when binary not found
- Both scripts work correctly after installation

### Issues Encountered
None - implementation straightforward.

## [2026-02-18] Add Clipboard Workflow Script for D&D Beyond

### Description
Created `ddb-copy.sh` wrapper script that automates the workflow of running character-tool and copying output files to clipboard for pasting into D&D Beyond. Uses `pbcopy` to populate clipboard history with each output file.

### Implementation
Created `ddb-copy.sh` with the following features:

1. **Vault-mode by default**: Runs character-tool with `--vault-mode` flag to output files next to the source markdown
2. **Smart file detection**: Only looks for specific character-tool output files (traits.txt, actions.txt, bonus-actions.txt, reactions.txt)
3. **Reverse order copying**: Copies files in reverse so they appear in correct order in clipboard history apps
4. **Graceful handling**: Works correctly when some sections are missing (e.g., no bonus actions or reactions)
5. **Visual feedback**: Color-coded output with green checkmarks and progress messages

### Script Usage
```bash
./ddb-copy.sh path/to/character.md
```

The script:
1. Runs character-tool with vault-mode
2. Finds generated .txt files in the same directory as input
3. Copies each file to clipboard with 0.3s delay between copies
4. Reports how many files were copied

### Design Decisions
- **macOS-specific**: Uses `pbcopy` which is macOS-only. Could be extended for Linux/Windows in future.
- **Clipboard history required**: Designed for use with clipboard managers (Paste, Maccy, CopyClip, etc.)
- **Reverse order**: Files copied in reverse so user can paste in forward order from clipboard history
- **Specific filenames**: Uses hardcoded list of expected output files rather than globbing to avoid issues in directories with many files
- **Small delays**: 0.3s sleep between copies ensures clipboard history registers each copy as separate entry

### Documentation
Updated README.md with new "Clipboard Workflow (macOS)" section including:
- Usage examples
- Requirements (macOS, clipboard history app)
- Example workflow showing full process
- Integration with existing Obsidian section

### Testing
Tested with:
- All four sections present (traits, actions, bonus-actions, reactions)
- Subset of sections (traits + actions only)
- Single section (actions only)
- Files generated in various directories (/tmp, project root, subdirectories)

All scenarios work correctly - script only copies files that exist.

### Issues Encountered
- Initial implementation used `find` with wildcard, which hung in /tmp due to too many files
- Solution: Changed to explicitly check for each expected filename rather than globbing

## [2026-02-18] Fix Dice Roll Display Format for Non-d20 Rolls

### Description
Changed dice roll display format so that non-d20 rolls show the full dice notation in the rollable tag, while d20 rolls continue to show only the modifier. This makes damage rolls more informative (e.g., `[rollable]1d8+5;...` instead of `[rollable]+5;...`).

### Tests Written/Modified
1. **New Tests in `converter/dice_test.go`:**
   - `TestGetDisplayValue_D20Rolls` - Verifies d20 rolls show only modifier (+5, empty string, -2, etc.)
   - `TestGetDisplayValue_NonD20Rolls` - Verifies non-d20 rolls show full notation (1d10+5, 2d6+3, 1d8, etc.)
   - `TestConvertDiceRolls_VariousDiceTypes` - Integration test for mixed d20 and non-d20 rolls

2. **Updated Tests in `converter/dice_test.go`:**
   - `TestConvertDiceRolls_Damage` - Changed expected from `[rollable]+3;` to `[rollable]2d6+3;`
   - `TestConvertDiceRolls_NoModifier` - Changed expected from `[rollable];` to `[rollable]1d10;`

3. **New Test in `formatter/formatter_test.go`:**
   - `TestFormatAbilities_D20VsDamageDisplay` - End-to-end test verifying d20 shows modifier only while damage shows full notation

### Implementation Details
Modified `converter/dice.go`:

1. **Added `isD20Roll()` helper function** (line 122-125):
   - Uses regex `^\d*d20([+-]\d+)?$` to identify d20 rolls
   - Returns true if notation uses d20 dice

2. **Added `getDisplayValue()` function** (line 127-135):
   - Returns only the modifier for d20 rolls (e.g., "+5" or "")
   - Returns full notation for non-d20 rolls (e.g., "1d8+5")
   - Centralizes display logic in one place

3. **Updated `ConvertDiceRolls()` function** (line 87-88):
   - Replaced `modifier := extractModifier(normalized)`
   - With `displayValue := getDisplayValue(normalized)`
   - Updated comments to reflect new behavior

### Design Decisions
- **Why differentiate d20 vs non-d20**: In D&D, attack rolls (d20) are familiar to players and the modifier is the key information. Damage rolls use varied dice (d4, d6, d8, d10, d12) and showing the full notation helps players understand what they're rolling.
- **Implementation approach**: Added helper functions rather than modifying existing `extractModifier()` to maintain backward compatibility and clear separation of concerns.
- **Regex pattern**: Used strict pattern `^\d*d20([+-]\d+)?$` to ensure only actual d20 rolls match, preventing edge cases.

### Test Results
All tests pass:
- `go test ./... -v` - All 30+ tests passing
- `go vet ./...` - No issues
- `gofmt -w .` - Code formatted

### Manual Verification
Created test markdown with mixed rolls:
- Longsword (d20+5 to hit, d8+5 damage) ✓
- Dagger (d20+3 to hit, d4+3 damage) ✓
- Greatsword (d20+7 to hit, 2d6+7 damage) ✓
- Fire Bolt (d20+5 to hit, 1d10 damage no modifier) ✓

Output correctly shows:
- To hit rolls: `[rollable]+5;...`, `[rollable]+3;...`, etc.
- Damage rolls: `[rollable]1d8+5;...`, `[rollable]1d4+3;...`, `[rollable]2d6+7;...`, `[rollable]1d10;...`

### Issues Encountered
None - implementation went smoothly following TDD approach.

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

## [2026-02-17] CLI Implementation

### Implementation
- Created `main.go` with CLI implementation:
  - Flag parsing for `-input`, `-output`, `-verbose`
  - File reading and markdown parsing
  - Spell list loading
  - Section processing (Traits, Actions, Bonus Actions, Reactions)
  - Output file writing with proper naming
  - Warning collection and display
  - Status reporting for each section
  - Empty section handling (skip file creation)

### Example File
- Created `testdata/example-character.md` with:
  - All four section types
  - Spell links ({{spell:Fireball}}, etc.)
  - Dice rolls with various roll types
  - Multiple abilities per section

### Testing
Manually tested with example file:
```
./character-tool -input testdata/example-character.md -output ./output -verbose
```

Verified output files:
- `traits.txt` - Spell links converted correctly
- `actions.txt` - Dice rolls with rollable tags and JSON metadata
- `bonus-actions.txt` - Healing roll converted correctly
- `reactions.txt` - Spell link converted correctly

All conversions working as expected!

### Next Steps
- Add .gitignore for build artifacts
- Finalize project documentation

## [2026-02-17] Project Summary

### Completed Components
1. ✅ Parser module - Markdown parsing with section extraction
2. ✅ Dice converter - Roll notation to D&D Beyond rollable format
3. ✅ Spell linker - {{spell:Name}} to [spell]Name[/spell] conversion
4. ✅ Formatter - Integration of all conversions
5. ✅ CLI - Command-line interface with file I/O
6. ✅ Documentation - PLAN.md, README.md, CLAUDE.md, JOURNAL.md
7. ✅ Example file - testdata/example-character.md

### Test Coverage
- Parser: 7 tests
- Dice converter: 14 tests
- Spell linker: 10 tests
- Formatter: 7 tests
- **Total: 38 unit tests, all passing**

### Project Status
**Implementation Complete!**

The tool successfully:
- Parses markdown character sheets
- Converts dice notation to rollable format
- Converts spell references to clickable links
- Validates spells and dice notation
- Outputs separate files per section
- Provides warning system for validation issues

All requirements from the plan have been implemented and tested.

## [2026-02-17] Refactor CLI to Use Cobra

### Changes Made
- Replaced stdlib `flag` package with `github.com/spf13/cobra`
- Added cobra and pflag dependencies to go.mod
- Restructured CLI with cobra command structure

### Implementation
- Created `rootCmd` with descriptive help text
- Added short flags: `-i` (input), `-o` (output), `-v` (verbose)
- Added long flags: `--input`, `--output`, `--verbose`, `--help`
- Improved help output with detailed description
- Maintained all existing functionality

### Benefits
- Better help output with automatic formatting
- Support for both short and long flags
- Built-in flag validation
- More extensible for future subcommands
- Industry-standard CLI framework

### Testing
- Verified `--help` output shows full documentation
- Tested with `-i` and `--input` flags (both work)
- Confirmed required flag validation
- Verified all output files generated correctly

All functionality preserved with improved CLI experience.

## [2026-02-17] Obsidian Integration

### Changes Made
- Added `--vault-mode` flag to CLI for Obsidian workflow
- Enhanced CLI output formatting for better readability
- Created comprehensive Obsidian integration documentation
- Added example Shell Commands plugin configuration
- Updated README with Obsidian integration section

### Implementation

#### CLI Enhancements (main.go)
**`--vault-mode` Flag:**
- Added new boolean flag that defaults output directory to input file's directory
- Simplifies Obsidian configuration by eliminating need to specify output path
- When enabled: `outputDir = filepath.Dir(inputFile)`
- Works seamlessly with Shell Commands plugin's `{{file_path:absolute}}` variable

**Output Formatting Improvements:**
- Changed from per-section status to summary format
- Shows success message: "✓ Formatted character abilities"
- Lists all created files with ability counts: `- path/file.txt (N abilities)`
- Improved warnings summary: Shows count in non-verbose mode
- Removed "Conversion complete!" trailing message for cleaner output
- Skips empty sections entirely (no status message, no file creation)

#### Documentation (docs/)
**OBSIDIAN_INTEGRATION.md:**
- Comprehensive setup guide covering prerequisites through troubleshooting
- Step-by-step Shell Commands plugin installation
- Configuration examples for basic and verbose modes
- Keyboard shortcut setup instructions
- Example workflow and file organization suggestions
- Troubleshooting section for common issues
- Links to additional resources

**shell-commands-config.json:**
- Ready-to-import Shell Commands configuration
- Two pre-configured commands:
  1. Standard format with `--vault-mode`
  2. Verbose format with `--vault-mode --verbose`
- Configured to show output in Obsidian notifications
- Uses `{{file_path:absolute}}` template variable

**README.md Updates:**
- Added `--vault-mode` flag to flags section
- New "Obsidian Integration" section with quick start
- Links to detailed integration guide
- Example command for Shell Commands plugin

### Design Decisions

**vault-mode Implementation:**
- Simple flag-based approach rather than detecting Obsidian environment
- Works anywhere, not just in Obsidian (useful for general same-directory output)
- No special handling of `{{vault}}` variable (Shell Commands resolves this)
- Clean separation: output logic doesn't need Obsidian awareness

**Output Format:**
- Summary-first approach better suited for Obsidian's notification system
- File paths shown with full absolute paths (Obsidian auto-links local files)
- Ability counts help users verify completeness at a glance
- Warnings shown as count unless verbose mode enabled

**Documentation Structure:**
- Separate detailed guide keeps README focused
- Example configuration file allows copy-paste setup
- Troubleshooting section addresses common path and permission issues
- Multiple command examples show verbose and custom output options

### Testing
Manual testing performed:
```bash
# Build tool
go build

# Test vault-mode flag
./character-tool --input testdata/example-character.md --vault-mode

# Verify output location
ls testdata/*.txt

# Test new output format
✓ Formatted character abilities

Output files:
  - testdata/traits.txt (2 abilities)
  - testdata/actions.txt (2 abilities)
  - testdata/bonus-actions.txt (1 abilities)
  - testdata/reactions.txt (1 abilities)
```

All tests successful:
- Files created in same directory as input
- Output format clean and informative
- Verbose mode shows warnings when present
- Empty sections handled correctly

### Benefits for Obsidian Users
1. **Simple Setup**: One-time Shell Commands configuration
2. **Fast Workflow**: Format with keyboard shortcut or command palette
3. **No Path Configuration**: vault-mode eliminates output directory complexity
4. **Clear Feedback**: Summary shows exactly what was created
5. **Vault Integration**: Output files appear immediately in same folder
6. **Copy-Ready**: Formatted text ready to paste into D&D Beyond

### Files Modified
- `main.go` - Added vaultMode flag, updated run() signature, improved output
- `README.md` - Added flag documentation and Obsidian section
- `JOURNAL.md` - This entry

### Files Created
- `docs/OBSIDIAN_INTEGRATION.md` - Complete setup guide
- `docs/shell-commands-config.json` - Example configuration

### Next Steps
- Consider adding auto-watch mode in future (re-format on save)
- Potential MCP server integration for AI agent access
- Template support for different character types

## [2026-02-18] Fix Parser Regex for Period Placement

### Issue
Parser only captured abilities with period inside bold markers (**Name.**) but failed to match period outside bold markers (**Name**.). This caused abilities to be silently skipped during parsing.

Example that failed:
```markdown
## Bonus Actions

**Second Wind.** Regain healing: 1d10+5 hit points.

**Fireball**. Kill all the things: 20d20 hit points.
```

Only "Second Wind" was parsed, "Fireball" was silently ignored.

### Tests Written
- `TestParseMarkdown_PeriodPlacement` - Tests both period placement styles (inside and outside bold markers)

### Implementation
Updated regex in `parser/parser.go`:
- **Old regex**: `^\*\*([^*]+)\.\*\*\s*(.+)$` (required period inside bold)
- **New regex**: `^\*\*([^*]+?)\.?\*\*\.?\s*(.+)$` (supports period inside or outside)

The new regex:
- `[^*]+?` - Non-greedy match for name (allows optional period after)
- `\.?` - Optional period inside bold markers
- `\*\*` - Closing bold markers
- `\.?` - Optional period outside bold markers
- Result: Matches both `**Name.**` and `**Name**.` patterns

### Testing
All existing tests pass, new test confirms both formats work:
```bash
go test ./parser -v
# All 8 tests PASS including new TestParseMarkdown_PeriodPlacement
```

Manual verification:
```bash
./character-tool --input test_input.md --output .
# Output: bonus-actions.txt (2 abilities)
# Previously: bonus-actions.txt (1 abilities)
```

### Design Decisions
- Made regex more permissive rather than enforcing single format
- Non-breaking change: existing markdown still works
- Better user experience: accepts both common bold+period styles
- No performance impact: regex still efficient with non-greedy matching

### Files Modified
- `parser/parser.go` - Updated ability regex pattern and comment
- `parser/parser_test.go` - Added TestParseMarkdown_PeriodPlacement test
- `JOURNAL.md` - This entry

### Root Cause
Original regex was too strict, assuming only one period placement style. Real-world markdown often mixes styles or uses period outside bold for aesthetic reasons.

## [2026-02-18] Document Dice Roll Keyword Requirements

### Issue
User reported dice notation not being converted to rollable format. Investigation revealed the text `things: 20d20` was not being converted because "things" is not a valid roll type keyword.

The dice converter only recognizes specific keywords:
- `to hit:` - Attack rolls
- `damage:` - Damage rolls
- `healing:` - Healing rolls
- `save:` - Saving throw rolls

Any other text followed by dice notation (like `things: 20d20`) remains as plain text.

### Changes Made
Updated documentation to clearly explain keyword requirements:

**README.md:**
- Expanded "Dice Rolls" section with keyword requirement
- Added "Important" callout explaining only specific keywords work
- Provided good/bad examples showing correct and incorrect usage
- Clarified that other text patterns won't be converted

**docs/OBSIDIAN_INTEGRATION.md:**
- Added troubleshooting section "Dice rolls not becoming rollable"
- Provided examples of working vs non-working patterns
- Listed valid keywords and supported dice types

### Design Decision
Chose to document behavior rather than expand keyword matching because:
1. **Specificity**: D&D Beyond expects specific roll types (to hit, damage, healing, save)
2. **Data Quality**: Precise keywords ensure correct rollType metadata in JSON
3. **User Control**: Users specify intent explicitly rather than tool guessing
4. **D&D Semantics**: Keywords map directly to D&D game mechanics
5. **No Breaking Changes**: Expanding keywords could cause unexpected conversions

### Files Modified
- `README.md` - Enhanced dice rolls section with examples
- `docs/OBSIDIAN_INTEGRATION.md` - Added troubleshooting entry
- `JOURNAL.md` - This entry

### User Education
Documentation now makes clear:
- Keywords are required, not optional
- Only four specific keywords are supported
- Provides correct usage patterns
- Shows common mistakes to avoid

## [2026-02-18] Refactor to Use strings.SplitSeq

### Issue
Linter warning: "Ranging over SplitSeq is more efficient" in parser.go:140

### Changes Made
Replaced `strings.Split()` with `strings.SplitSeq()` for more efficient iteration:
- **Before**: `paragraphs := strings.Split(content, "\n\n")` then `for _, paragraph := range paragraphs`
- **After**: `for paragraph := range strings.SplitSeq(content, "\n\n")`

### Benefits
- **Memory Efficiency**: SplitSeq returns an iterator, avoiding intermediate slice allocation
- **Performance**: Lazily generates values instead of allocating entire slice upfront
- **Modern Go**: Uses Go 1.23+ iterator pattern

### Testing
All existing parser tests pass without modification:
```bash
go test ./parser -v
# All 8 tests PASS
```

### Files Modified
- `parser/parser.go` - Updated parseAbilities to use SplitSeq
- `JOURNAL.md` - This entry
