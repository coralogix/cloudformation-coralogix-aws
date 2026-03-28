!/bin/bash
# Wrapper script to build and run the Go version of generate_fields
# Usage: ./generate_fields_wrapper.sh <integration_name> <new_version> <integration_definitions_path>

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

go build -o generate_fields generate_fields.go

./generate_fields "$@"
rm -f generate_fields
