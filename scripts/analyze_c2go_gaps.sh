#!/bin/bash

# Script to analyze gaps between C functions inventory and current TODO
# This helps identify which functions are not yet tracked for c2go porting

set -euo pipefail

INVENTORY_FILE="${1:-c_functions_inventory.md}"
TODO_FILE="${2:-internal/c2go/TODO.md}"
OUTPUT_FILE="${3:-c2go_gaps_analysis.md}"

if [[ ! -f "$INVENTORY_FILE" ]]; then
    echo "Error: Inventory file not found: $INVENTORY_FILE"
    echo "Usage: $0 [inventory_file] [todo_file] [output_file]"
    exit 1
fi

if [[ ! -f "$TODO_FILE" ]]; then
    echo "Error: TODO file not found: $TODO_FILE"
    echo "Usage: $0 [inventory_file] [todo_file] [output_file]"
    exit 1
fi

echo "Analyzing gaps between inventory and TODO..."
echo "📄 Inventory: $INVENTORY_FILE"
echo "📄 TODO: $TODO_FILE"
echo "📄 Output: $OUTPUT_FILE"

# Extract function names from inventory
echo "Extracting function names from inventory..."
inventory_functions=$(grep -E '^\- \*\*Line [0-9]+\*\*: `[^`]+`' "$INVENTORY_FILE" | sed 's/.*`\([^`]*\)`.*/\1/' | sort -u || echo "")

# Extract function names from TODO (look for function names in backticks followed by parentheses)
echo "Extracting function names from TODO..."
todo_functions=$(grep -o -E '`[_a-zA-Z][a-zA-Z0-9_]*\(' "$TODO_FILE" | sed 's/`//g' | sed 's/(.*//' | grep -v '^cap$' | grep -v '^H3_EXPORT$' | sort -u || echo "")

# Create output file
cat > "$OUTPUT_FILE" << 'EOF'
# C2Go Gaps Analysis

This file identifies C functions that exist in the H3 source but are not yet tracked in TODO.md for c2go porting.

## Summary

EOF

# Count totals (handle empty results)
if [[ -n "$inventory_functions" ]]; then
    inventory_count=$(echo "$inventory_functions" | wc -l)
else
    inventory_count=0
fi

if [[ -n "$todo_functions" ]]; then
    todo_count=$(echo "$todo_functions" | wc -l)
else
    todo_count=0
fi

echo "- **Total functions in inventory:** $inventory_count" >> "$OUTPUT_FILE"
echo "- **Total functions in TODO:** $todo_count" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Find missing functions (in inventory but not in TODO)
echo "Finding functions that are in inventory but missing from TODO..."
if [[ -n "$inventory_functions" && -n "$todo_functions" ]]; then
    missing_functions=$(comm -23 <(echo "$inventory_functions" | sort) <(echo "$todo_functions" | sort))
elif [[ -n "$inventory_functions" && -z "$todo_functions" ]]; then
    missing_functions="$inventory_functions"
else
    missing_functions=""
fi
missing_count=$(echo "$missing_functions" | grep -c . || echo "0")

echo "- **Functions missing from TODO:** $missing_count" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

if [[ $missing_count -gt 0 ]]; then
    echo "## Missing Functions" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    echo "These functions exist in H3 source but are not tracked in TODO.md:" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    
    # List missing functions (simplified to avoid performance issues)
    echo "$missing_functions" | while IFS= read -r func; do
        if [[ -n "$func" ]]; then
            echo "- \`$func\`" >> "$OUTPUT_FILE"
        fi
    done
fi

# Find functions in TODO but not in inventory (possibly already implemented or removed)
echo "" >> "$OUTPUT_FILE"
echo "## Functions in TODO but not in Current Inventory" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "These might be already implemented, renamed, or removed from H3 source:" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

if [[ -n "$inventory_functions" && -n "$todo_functions" ]]; then
    todo_only_functions=$(comm -13 <(echo "$inventory_functions" | sort) <(echo "$todo_functions" | sort))
elif [[ -z "$inventory_functions" && -n "$todo_functions" ]]; then
    todo_only_functions="$todo_functions"
else
    todo_only_functions=""
fi
todo_only_count=$(echo "$todo_only_functions" | grep -c . || echo "0")

if [[ $todo_only_count -gt 0 ]]; then
    for func in $todo_only_functions; do
        if [[ -n "$func" ]]; then
            echo "- \`$func\`" >> "$OUTPUT_FILE"
        fi
    done
else
    echo "*(None found - all TODO functions exist in current inventory)*" >> "$OUTPUT_FILE"
fi

# Prioritization suggestions
echo "" >> "$OUTPUT_FILE"
echo "## Suggested Next Targets" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "Based on simplicity and usefulness, here are some recommended functions to implement next:" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Look for simple patterns in missing functions (single-line functions, getters, constants)
simple_candidates=""
if [[ $missing_count -gt 0 ]]; then
    for func in $missing_functions; do
        if [[ -n "$func" ]]; then
            # Check if it looks like a simple function (getter, boolean check, constant return)
            if [[ "$func" =~ ^get[A-Z] ]] || [[ "$func" =~ ^is[A-Z] ]] || [[ "$func" =~ ^_is[A-Z] ]] || [[ "$func" =~ Count$ ]] || [[ "$func" =~ ^_get ]]; then
                source_file=$(grep -B 20 -E "^\- \*\*Line [0-9]+\*\*: \`$func\`" "$INVENTORY_FILE" | grep "^## " | tail -1 | sed 's/^## //')
                simple_candidates="$simple_candidates- \`$func\` (from $source_file) - looks like a getter/checker/counter\n"
            fi
        fi
    done
    
    if [[ -n "$simple_candidates" ]]; then
        echo "**Simple utility functions (getters, booleans, counters):**" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
        echo -e "$simple_candidates" >> "$OUTPUT_FILE"
    else
        echo "*(No obviously simple functions identified in missing list)*" >> "$OUTPUT_FILE"
    fi
fi

echo "" >> "$OUTPUT_FILE"
echo "---" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "*Generated on $(date) by c2go gaps analysis script*" >> "$OUTPUT_FILE"

echo ""
echo "✅ Gap analysis complete!"
echo "📊 Found $missing_count functions missing from TODO"
echo "📄 Report written to: $OUTPUT_FILE"
echo ""
echo "Next steps:"
echo "1. Review the gap analysis report"
echo "2. Add missing functions to TODO.md"
echo "3. Prioritize simple functions for next iterations"