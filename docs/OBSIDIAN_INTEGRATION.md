# Obsidian Integration Guide

This guide shows you how to use character-tool directly from Obsidian using the Shell Commands community plugin. This integration allows you to format your D&D character markdown files with a single command or keyboard shortcut.

## Overview

The integration works by:
1. Installing the Shell Commands plugin in Obsidian
2. Configuring a command that calls character-tool on the current file
3. Using the command from Obsidian's command palette or a keyboard shortcut

The formatted files will appear in the same directory as your character markdown file, ready to copy into D&D Beyond.

## Prerequisites

### 1. Install Go (if not already installed)

character-tool is written in Go, so you need Go installed to build it:

- **macOS**: `brew install go`
- **Linux**: `sudo apt install golang` or `sudo dnf install golang`
- **Windows**: Download from [golang.org/dl](https://golang.org/dl)

Verify installation: `go version`

### 2. Build character-tool

Clone and build the tool:

```bash
git clone https://github.com/yourusername/character-tool.git
cd character-tool
go build
```

This creates an executable file named `character-tool` (or `character-tool.exe` on Windows).

### 3. Note the absolute path

Get the full path to the executable:

```bash
# macOS/Linux
pwd
# Shows something like: /Users/yourname/character-tool

# Windows (PowerShell)
Get-Location
# Shows something like: C:\Users\yourname\character-tool
```

The full path to your executable will be:
- macOS/Linux: `/Users/yourname/character-tool/character-tool`
- Windows: `C:\Users\yourname\character-tool\character-tool.exe`

## Installing Shell Commands Plugin

1. Open Obsidian Settings (gear icon)
2. Navigate to **Community plugins**
3. Click **Browse** and search for "Shell Commands"
4. Install the plugin by **Mithril0x**
5. Enable the plugin

## Configuring Shell Commands

### Method 1: Manual Configuration

1. Open Obsidian Settings
2. Go to **Shell Commands** in the left sidebar
3. Click **New shell command**
4. Enter the command (replace `/path/to/character-tool` with your actual path):

   ```bash
   /path/to/character-tool --input "{{file_path:absolute}}" --vault-mode
   ```

5. Give it an alias like "Format D&D Character"
6. Click the check mark to save

### Method 2: Import Configuration (Advanced)

See [shell-commands-config.json](shell-commands-config.json) for an importable configuration file.

## Command Options

### Basic Command (Recommended)
```bash
/path/to/character-tool --input "{{file_path:absolute}}" --vault-mode
```
- Uses vault-mode to output files in the same directory as your markdown file
- Silent on success, shows errors if something goes wrong

### Verbose Command
```bash
/path/to/character-tool --input "{{file_path:absolute}}" --vault-mode --verbose
```
- Shows validation warnings (unknown spells, invalid dice notation)
- Useful for troubleshooting

### Custom Output Directory
```bash
/path/to/character-tool --input "{{file_path:absolute}}" --output "{{folder_path:absolute}}/formatted"
```
- Outputs to a specific subdirectory
- Omit `--vault-mode` when using custom output path

## Setting Up a Keyboard Shortcut

1. In Shell Commands settings, find your command
2. Click the keyboard icon next to it
3. Click **Add hotkey**
4. Press your desired key combination (e.g., `Cmd+Shift+F` or `Ctrl+Shift+F`)
5. Click outside the hotkey field to save

## Usage

### From Command Palette
1. Open a character markdown file
2. Press `Cmd+P` (macOS) or `Ctrl+P` (Windows/Linux) to open command palette
3. Type "Format D&D Character"
4. Press Enter

### From Keyboard Shortcut
1. Open a character markdown file
2. Press your configured hotkey (e.g., `Cmd+Shift+F`)

### Expected Output

The tool creates four files in your vault:
- `traits.txt` - Character traits and features
- `actions.txt` - Actions the character can take
- `bonus-actions.txt` - Bonus actions
- `reactions.txt` - Reactions

These files appear in the same folder as your character markdown file (when using `--vault-mode`).

## Example Workflow

1. Create a character file in Obsidian: `Characters/Wizard.md`
2. Add abilities using the structured format:
   ```markdown
   # Traits
   ## Arcane Recovery
   Once per day regain spell slots...

   # Actions
   ## Fire Bolt
   Ranged spell attack: {{spell:Fire Bolt}}
   Attack: to hit: 1d20+5
   Damage: damage: 1d10 fire
   ```
3. Run the format command (palette or hotkey)
4. Open the generated `.txt` files
5. Copy formatted text into D&D Beyond

## File Organization

Consider organizing your vault like this:

```
YourVault/
├── Characters/
│   ├── Wizard.md          (your source file)
│   ├── traits.txt         (generated)
│   ├── actions.txt        (generated)
│   ├── bonus-actions.txt  (generated)
│   └── reactions.txt      (generated)
└── Templates/
    └── character-template.md
```

You can add `*.txt` to your `.gitignore` or Obsidian's excluded files if you don't want to track generated files.

## Troubleshooting

### Command not found
- Verify the path to character-tool is correct and absolute
- Test the command in your terminal first
- On Windows, use backslashes: `C:\path\to\character-tool.exe`

### Permission denied
- Make the tool executable: `chmod +x /path/to/character-tool`
- On macOS, you may need to allow the app in System Preferences > Security

### No output files created
- Check that your markdown file has the correct headers: `# Traits`, `# Actions`, etc.
- Run with `--verbose` to see validation warnings
- Verify the input file path is correct

### Spell validation warnings
- The tool validates spell names against the SRD spell list
- Use exact spell names: `{{spell:Fire Bolt}}` not `{{spell:Firebolt}}`
- Custom spells will show warnings but still format correctly

### Output appears in wrong location
- Make sure you're using `--vault-mode` flag
- The flag outputs to the same directory as the input file
- Without the flag, output goes to current working directory

## Advanced Configuration

### Multiple Commands

You can create multiple commands for different scenarios:

1. **Quick Format**: Standard command with `--vault-mode`
2. **Verbose Format**: Same but with `--verbose` flag
3. **Export Format**: Custom output directory for sharing

### Automation

While the Shell Commands plugin doesn't support auto-execution on file save, you can:
- Set up a convenient keyboard shortcut
- Run the command before copying to D&D Beyond
- Create a workflow where you format before committing changes

## Support

For issues with:
- **character-tool**: Open an issue on GitHub
- **Shell Commands plugin**: Check the plugin's documentation
- **Obsidian**: Visit the Obsidian community forum

## Additional Resources

- [Character Tool Documentation](../README.md)
- [Shell Commands Plugin](https://github.com/Taitava/obsidian-shellcommands)
- [Markdown Format Guide](../README.md#markdown-format)
