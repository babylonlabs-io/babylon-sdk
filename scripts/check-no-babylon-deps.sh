#!/bin/bash
set -e

echo "Checking that demo/ and x/ directories don't use babylonlabs-io/babylon dependencies..."

# Function to check for babylon dependencies in a go.mod file
check_go_mod() {
    local mod_file="$1"
    local dir_name="$2"

    if [ ! -f "$mod_file" ]; then
        echo "ERROR: $mod_file not found"
        return 1
    fi

    echo "Checking $mod_file..."

    # Check for any babylon dependencies (excluding babylon-sdk and local workspace references)
    babylon_deps=$(grep -E "github\.com/babylonlabs-io/babylon($|[^-])" "$mod_file" | grep -v "// local work dir" | grep -v "=> \.\." || true)

    if [ -n "$babylon_deps" ]; then
        echo "ERROR: Found babylonlabs-io/babylon dependencies in $dir_name/:"
        echo "$babylon_deps"
        return 1
    fi

    echo "✓ No babylonlabs-io/babylon dependencies found in $dir_name/"
}

# Check demo/go.mod
check_go_mod "demo/go.mod" "demo"

# Check x/go.mod
check_go_mod "x/go.mod" "x"

echo "✓ All checks passed"
