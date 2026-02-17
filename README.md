# D&D Beyond Character Tool

A Go CLI tool that converts markdown character descriptions into D&D Beyond-formatted blocks with spell links and rollable dice notation.

## Features

- Converts markdown character sheets to D&D Beyond format
- Generates rollable dice notation with proper metadata
- Creates clickable spell links
- Validates spell names and dice notation
- Outputs separate files for each ability section

## Installation

```bash
go build -o character-tool
```

## Usage

```bash
./character-tool -input character.md -output ./output
```

### Flags

- `-input`: Path to input markdown file (required)
- `-output`: Output directory for generated files (default: current directory)
- `-verbose`: Show detailed validation warnings

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

Use inline dice notation with keyword prefixes:

- `to hit: 1d20+5` - Attack rolls
- `damage: 2d6+3` - Damage rolls
- `healing: 1d8+4` - Healing rolls
- `DC 15` - Saving throw DCs
- `save: 1d20+2` - Saving throw rolls

Supported dice: d4, d6, d8, d10, d12, d20, d100

## Output

The tool generates four files:

- `traits.txt` - Character traits
- `actions.txt` - Actions
- `bonus-actions.txt` - Bonus actions
- `reactions.txt` - Reactions

Each file contains D&D Beyond-formatted text ready to paste into character sheets.

## Example Output

**traits.txt:**
```
Spellcasting. You can cast spells using [spell]Fireball[/spell] and [spell]Magic Missile[/spell].

Pack Tactics. You have advantage on attack rolls against a creature if at least one ally is within 5 feet.
```

**actions.txt:**
```
Quarterstaff. Melee Weapon Attack: [rollable]+2;{"diceNotation":"1d20+2","rollType":"to hit","rollAction":"Quarterstaff"}[/rollable], reach 5 ft., one target. Hit: [rollable]+2;{"diceNotation":"1d6+2","rollType":"damage","rollAction":"Quarterstaff"}[/rollable] bludgeoning damage.
```

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
