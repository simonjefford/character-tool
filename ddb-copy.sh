#!/bin/bash
# ddb-copy.sh - Run character-tool and copy output files to clipboard
#
# Usage: ./ddb-copy.sh <input-file.md>
#
# This script:
# 1. Runs character-tool on the input file (with --vault-mode)
# 2. Copies each generated .txt file to clipboard (in reverse order)
# 3. Files appear in clipboard history for pasting into D&D Beyond

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check arguments
if [ $# -eq 0 ]; then
    echo "Usage: $0 <input-file.md>"
    echo ""
    echo "Runs character-tool with --vault-mode (outputs to same directory as input)"
    echo ""
    echo "Example:"
    echo "  $0 ~/Documents/characters/fighter.md"
    exit 1
fi

INPUT_FILE="$1"

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file '$INPUT_FILE' not found"
    exit 1
fi

# Get absolute path to input file and its directory
INPUT_FILE_ABS=$(cd "$(dirname "$INPUT_FILE")" && pwd)/$(basename "$INPUT_FILE")
OUTPUT_DIR=$(dirname "$INPUT_FILE_ABS")

# Find character-tool binary (check script directory first, then PATH)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/character-tool" ]; then
    CHARACTER_TOOL="$SCRIPT_DIR/character-tool"
elif command -v character-tool >/dev/null 2>&1; then
    CHARACTER_TOOL="character-tool"
else
    echo "Error: character-tool binary not found"
    echo "  - Not found in script directory: $SCRIPT_DIR/character-tool"
    echo "  - Not found in PATH"
    echo ""
    echo "Build it with: go build -o character-tool"
    exit 1
fi

echo -e "${BLUE}Running character-tool...${NC}"
"$CHARACTER_TOOL" -i "$INPUT_FILE_ABS" --vault-mode

# Find generated .txt files (traits, actions, bonus-actions, reactions)
TXT_FILES=()
for filename in "traits.txt" "bonus-actions.txt" "reactions.txt" "actions.txt"; do
    filepath="$OUTPUT_DIR/$filename"
    if [ -f "$filepath" ]; then
        TXT_FILES+=("$filepath")
    fi
done

if [ ${#TXT_FILES[@]} -eq 0 ]; then
    echo -e "${YELLOW}No .txt files generated${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}Copying files to clipboard (in reverse order)...${NC}"

# Copy files in reverse order so they appear in correct order in clipboard history
for ((i=${#TXT_FILES[@]}-1; i>=0; i--)); do
    FILE="${TXT_FILES[$i]}"
    BASENAME=$(basename "$FILE")

    # Copy file content to clipboard
    cat "$FILE" | pbcopy

    echo -e "${GREEN}✓${NC} Copied: $BASENAME"

    # Small delay to ensure clipboard history registers each copy
    sleep 0.3
done

echo ""
echo -e "${GREEN}Done!${NC} ${#TXT_FILES[@]} file(s) copied to clipboard history"
echo "Paste them into D&D Beyond in order using your clipboard history app"
