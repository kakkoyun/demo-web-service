#!/usr/bin/env bash
#
# Runs all linters and formatters across the codebase
# Usage: ./tools/lint-all.sh [--fix]
#

set -euo pipefail

# Color variables for better output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly NC='\033[0m' # No Color

# Script directory for reference
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Parse arguments
FIX_MODE=false
for arg in "$@"; do
  case "${arg}" in
    --fix)
      FIX_MODE=true
      ;;
    *)
      echo -e "${RED}Unknown argument: ${arg}${NC}"
      echo "Usage: $0 [--fix]"
      exit 1
      ;;
  esac
done

# Check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Show instruction for installing a tool
show_install_instruction() {
  local tool=$1
  local instruction=$2
  echo -e "${YELLOW}Warning: ${tool} not found in PATH${NC}"
  echo -e "${YELLOW}${instruction}${NC}"
}

# Header for each section
print_header() {
  echo -e "\n${GREEN}==>${NC} ${YELLOW}$1${NC}"
}

# Main function that runs all linters
main() {
  cd "${REPO_ROOT}"

  print_header "Running Go linters"
  if ${FIX_MODE}; then
    go tool golangci-lint run ./... --fix
  else
    go tool golangci-lint run ./...
  fi

  print_header "Checking shell scripts with shellcheck"
  if command_exists shellcheck; then
    # Find all shell scripts and run shellcheck
    find . -name "*.sh" -not -path "./vendor/*" -exec shellcheck {} \;
  else
    show_install_instruction "shellcheck" "Install with your OS package manager: brew install shellcheck (macOS) or apt-get install shellcheck (Linux)"
    echo "Shell script linting will be skipped"
  fi

  if ${FIX_MODE}; then
    print_header "Formatting shell scripts with shfmt"
    # Use go tool with shfmt (Go 1.24 API)
    find . -name "*.sh" -not -path "./vendor/*" -exec go tool shfmt -i 2 -ci -w {} \;
  else
    print_header "Checking shell script formatting"
    find . -name "*.sh" -not -path "./vendor/*" -exec go tool shfmt -d -i 2 -ci {} \;
  fi

  if ${FIX_MODE}; then
    print_header "Formatting YAML files with yamlfmt"
    # Find all YAML files and format them
    find . \( -name "*.yml" -o -name "*.yaml" \) -not -path "./vendor/*" -exec go tool yamlfmt {} \;
  else
    print_header "Checking YAML formatting"
    # Check YAML formatting
    find . \( -name "*.yml" -o -name "*.yaml" \) -not -path "./vendor/*" -exec go tool yamlfmt -lint {} \;
  fi

  print_header "Checking if README code snippets are up-to-date"
  go tool embedmd -d README.md

  print_header "All checks completed!"
  echo -e "${GREEN}✓ Linting passed${NC}"
}

main "$@"
