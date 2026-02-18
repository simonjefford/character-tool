#!/bin/bash
# install.sh - Install character-tool and ddb-copy.sh to ~/bin
#
# This script:
# 1. Builds the character-tool binary
# 2. Copies both character-tool and ddb-copy.sh to ~/bin
# 3. Ensures ~/bin is in PATH

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

INSTALL_DIR="$HOME/bin"

echo -e "${BLUE}Character Tool Installer${NC}"
echo ""

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo -e "${RED}Error: Must be run from character-tool directory${NC}"
    echo "This script should be run from the root of the character-tool repository"
    exit 1
fi

# Create ~/bin if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
    echo -e "${BLUE}Creating $INSTALL_DIR...${NC}"
    mkdir -p "$INSTALL_DIR"
    echo -e "${GREEN}✓${NC} Created $INSTALL_DIR"
    echo ""
fi

# Build the binary
echo -e "${BLUE}Building character-tool...${NC}"
go build -o character-tool
echo -e "${GREEN}✓${NC} Built character-tool binary"
echo ""

# Copy files to ~/bin
echo -e "${BLUE}Installing to $INSTALL_DIR...${NC}"
cp character-tool "$INSTALL_DIR/character-tool"
echo -e "${GREEN}✓${NC} Installed character-tool"

cp ddb-copy.sh "$INSTALL_DIR/ddb-copy.sh"
chmod +x "$INSTALL_DIR/ddb-copy.sh"
echo -e "${GREEN}✓${NC} Installed ddb-copy.sh"
echo ""

# Check if ~/bin is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}Warning: $INSTALL_DIR is not in your PATH${NC}"
    echo ""
    echo "Add it to your PATH by adding this line to your shell config:"
    echo ""

    # Detect shell
    if [ -n "$ZSH_VERSION" ]; then
        echo "  echo 'export PATH=\"\$HOME/bin:\$PATH\"' >> ~/.zshrc"
        echo "  source ~/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        echo "  echo 'export PATH=\"\$HOME/bin:\$PATH\"' >> ~/.bashrc"
        echo "  source ~/.bashrc"
    else
        echo "  export PATH=\"\$HOME/bin:\$PATH\""
    fi
    echo ""
else
    echo -e "${GREEN}✓${NC} $INSTALL_DIR is already in PATH"
    echo ""
fi

echo -e "${GREEN}Installation complete!${NC}"
echo ""
echo "Usage:"
echo "  character-tool -i character.md --vault-mode"
echo "  ddb-copy.sh character.md"
