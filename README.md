# D&D Beyond Character Tool

A Go CLI tool that converts markdown character descriptions into D&D Beyond-formatted blocks with spell links and rollable dice notation.

## Features

- **Markdown to D&D Beyond format** - Convert character sheets with proper formatting
- **Rollable dice notation** - Attack rolls show modifiers, damage rolls show averages
- **Average damage calculation** - DMs can use averages (e.g., `8(1d8+3)`) for quick resolution
- **Spell links** - Auto-generates `[spell]SpellName[/spell]` tags with validation
- **Plain text support** - Include context paragraphs alongside named abilities
- **Clipboard workflow** - macOS script to copy outputs directly to clipboard history
- **Separate output files** - One file per section (traits, actions, bonus actions, reactions)

## Installation

### Quick Install (Recommended)

Install to `~/bin` with the provided installer:

```bash
./install.sh
```

This will:
1. Build the `character-tool` binary
2. Copy both `character-tool` and `ddb-copy.sh` to `~/bin`
3. Make them executable and available system-wide

### Manual Build

```bash
go build -o character-tool
```

## Usage

```bash
character-tool --input character.md --output ./output
```

Or using short flags:

```bash
character-tool -i character.md -o ./output
```

### Flags

- `-i, --input`: Path to input markdown file (required)
- `-o, --output`: Output directory for generated files (default: current directory)
- `--vault-mode`: Output files to same directory as input file (useful for Obsidian)
- `-v, --verbose`: Show detailed validation warnings
- `-h, --help`: Show help message

## Input Format

Create a markdown file with the following structure:

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

### Spell Links

Use `{{spell:SpellName}}` syntax to create spell links. The tool validates against the D&D 5e spell list.

### Dice Rolls

Dice notation must be preceded by one of the following keywords to be converted to rollable format:

- `to hit: 1d20+5` - Attack rolls
- `damage: 2d6+3` - Damage rolls
- `healing: 1d8+4` - Healing rolls
- `save: 1d20+2` - Saving throw rolls

**Important**: Only dice notation with these specific keywords will be made rollable. Other text with dice notation (like `things: 20d20`) will remain as plain text.

Examples:
```markdown
**Good**: Regain healing: 1d10+5 hit points.
**Bad**: Regain 1d10+5 hit points. (no keyword)
**Bad**: Regain things: 1d10+5 hit points. (wrong keyword)
```

Supported dice types: d4, d6, d8, d10, d12, d20, d100

#### Dice Roll Display Format

The tool formats dice rolls differently based on type:

**Attack rolls (d20):** Show only the modifier
- Input: `to hit: 1d20+5`
- Output: `[rollable]+5;{...}[/rollable]`

**Damage/Healing rolls (non-d20):** Show average and full notation
- Input: `damage: 1d8+3`
- Output: `[rollable]8(1d8+3);{...}[/rollable]`
- Input: `damage: 2d6+7`
- Output: `[rollable]14(2d6+7);{...}[/rollable]`

The average is calculated as: `(number_of_dice × average_per_die) + modifier`, rounded to nearest integer. This allows DMs to quickly use average damage instead of rolling.

### Plain Text Paragraphs

You can include plain text paragraphs (without the `**Name.**` format) in any section to provide context or explanations:

```markdown
## Traits

This character has enhanced abilities due to their elven heritage.

**Darkvision.** You can see in dim light within 60 feet.

The following traits come from their scholar background.

**Researcher.** When attempting to learn information, you can consult your notes.
```

Plain text paragraphs are preserved in the output without the "Name." prefix.

## Output

The tool generates four files:

- `traits.txt` - Character traits
- `actions.txt` - Actions
- `bonus-actions.txt` - Bonus actions
- `reactions.txt` - Reactions

Each file contains D&D Beyond-formatted text ready to paste into character sheets.

## Clipboard Workflow (macOS)

The `ddb-copy.sh` script automates copying output files to your clipboard for easy pasting into D&D Beyond.

### Setup

Install with `./install.sh` (recommended) or ensure `character-tool` is in your PATH.

The script will look for the binary in:
1. The same directory as the script
2. Your PATH (e.g., `~/bin`)

### Usage

```bash
ddb-copy.sh path/to/character.md
```

This script:
1. Runs character-tool with `--vault-mode` (outputs next to your input file)
2. Copies each generated .txt file to clipboard in reverse order
3. Files appear in your clipboard history app (like Paste, Maccy, etc.)
4. Paste them into D&D Beyond in order from your clipboard history

### Requirements

- macOS (uses `pbcopy`)
- Clipboard history app (recommended: Paste, Maccy, CopyClip)

### Example Workflow

```bash
# 1. Run the script on your character file
ddb-copy.sh ~/Documents/fighter.md

# Output:
# ✓ Copied: traits.txt
# ✓ Copied: reactions.txt
# ✓ Copied: actions.txt
# Done! 3 file(s) copied to clipboard history

# 2. Open D&D Beyond character sheet
# 3. Open your clipboard history app
# 4. Paste each section into the corresponding D&D Beyond field
```

## Obsidian Integration

You can use character-tool directly from Obsidian using the Shell Commands community plugin. This allows you to format character files with a single command or keyboard shortcut.

### Quick Start

1. Install the [Shell Commands](https://github.com/Taitava/obsidian-shellcommands) plugin in Obsidian
2. Add a new command:
   ```bash
   /path/to/character-tool --input "{{file_path:absolute}}" --vault-mode
   ```
3. Set a keyboard shortcut (optional)
4. Run the command on any character markdown file

The formatted `.txt` files will appear in the same directory as your character file.

### Full Setup Guide

See [docs/OBSIDIAN_INTEGRATION.md](docs/OBSIDIAN_INTEGRATION.md) for detailed installation instructions, configuration options, and troubleshooting.

## Example Output

**Input markdown:**
```markdown
## Traits

This half-elf has unique abilities from their heritage.

**Darkvision.** You can see in dim light within 60 feet.

## Actions

**Longsword.** Melee Weapon Attack: to hit: 1d20+5, reach 5 ft. Hit: damage: 1d8+3 slashing.

**Fire Bolt.** Ranged Spell Attack: to hit: 1d20+5, range 120 ft. Hit: damage: 1d10 fire damage.
```

**traits.txt:**
```
This half-elf has unique abilities from their heritage.

Darkvision. You can see in dim light within 60 feet.
```

**actions.txt:**
```
Longsword. Melee Weapon Attack: [rollable]+5;{"diceNotation":"1d20+5","rollType":"to hit","rollAction":"Longsword"}[/rollable], reach 5 ft. Hit: [rollable]8(1d8+3);{"diceNotation":"1d8+3","rollType":"damage","rollAction":"Longsword"}[/rollable] slashing.

Fire Bolt. Ranged Spell Attack: [rollable]+5;{"diceNotation":"1d20+5","rollType":"to hit","rollAction":"Fire Bolt"}[/rollable], range 120 ft. Hit: [rollable]6(1d10);{"diceNotation":"1d10","rollType":"damage","rollAction":"Fire Bolt"}[/rollable] fire damage.
```

Notice:
- Plain text paragraph preserved without "Name." prefix
- Attack rolls show modifier only: `+5`
- Damage rolls show average and notation: `8(1d8+3)` and `6(1d10)`

## Development

See `PLAN.md` for implementation details and `JOURNAL.md` for development history.

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o character-tool
```

## License

MIT
