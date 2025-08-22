#!/bin/bash

# Script to extract all C functions from H3 source files
# This helps identify candidates for c2go porting

set -euo pipefail

H3_SRC_DIR="${1:-testref/h3-4.3.0/src/h3lib/lib}"
OUTPUT_FILE="${2:-c_functions_inventory.md}"

if [[ ! -d "$H3_SRC_DIR" ]]; then
    echo "Error: H3 source directory not found: $H3_SRC_DIR"
    echo "Usage: $0 [h3_src_dir] [output_file]"
    exit 1
fi

echo "Extracting C functions from: $H3_SRC_DIR"
echo "Output file: $OUTPUT_FILE"

# Create output file with header
cat > "$OUTPUT_FILE" << 'EOF'
# H3 C Functions Inventory

This file contains all C functions found in the H3 source code, grouped by source file.
Generated automatically for c2go porting planning.

**Legend:**
- `static` functions are internal to the file
- `H3_EXPORT` functions are part of the public API
- Functions without prefix are internal but available for porting

EOF

# Process each C file
for c_file in "$H3_SRC_DIR"/*.c; do
    if [[ ! -f "$c_file" ]]; then
        continue
    fi
    
    filename=$(basename "$c_file")
    echo "Processing $filename..."
    
    # Extract function definitions - look for patterns that start with return type and function name
    # This pattern specifically looks for function definitions that start a line
    functions=$(grep -n -E '^[a-zA-Z_][a-zA-Z0-9_*[:space:]]+[a-zA-Z_][a-zA-Z0-9_]*[[:space:]]*\([^)]*\)[[:space:]]*\{' "$c_file" || true)
    
    # Also catch one-liner functions
    oneline_functions=""
    
    if [[ -n "$functions" ]] || [[ -n "$oneline_functions" ]]; then
        echo "" >> "$OUTPUT_FILE"
        echo "## $filename" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
        
        # Combine and sort by line number
        {
            if [[ -n "$functions" ]]; then
                echo "$functions"
            fi
            if [[ -n "$oneline_functions" ]]; then
                echo "$oneline_functions"
            fi
        } | sort -n | while IFS=':' read -r line_no line_content; do
            # Clean up the line and extract key information
            clean_line=$(echo "$line_content" | sed 's/^[[:space:]]*//' | sed 's/[[:space:]]*{.*$//')
            
            # Determine function characteristics
            is_static=""
            is_export=""
            if [[ "$clean_line" =~ ^static[[:space:]] ]]; then
                is_static="[STATIC] "
                clean_line=$(echo "$clean_line" | sed 's/^static[[:space:]]*//')
            fi
            if [[ "$clean_line" =~ H3_EXPORT ]]; then
                is_export="[EXPORT] "
                clean_line=$(echo "$clean_line" | sed 's/H3_EXPORT[[:space:]]*([^)]*)//g' | sed 's/^[[:space:]]*//')
            fi
            
            # Extract function name (last word before parentheses)
            func_name=$(echo "$clean_line" | sed -n 's/.*[[:space:]]\([a-zA-Z_][a-zA-Z0-9_]*\)[[:space:]]*(.*/\1/p')
            if [[ -z "$func_name" ]]; then
                # Try alternative pattern for cases where there's no space before function name
                func_name=$(echo "$clean_line" | sed -n 's/.*\([a-zA-Z_][a-zA-Z0-9_]*\)[[:space:]]*(.*/\1/p')
            fi
            
            # Format output
            echo "- **Line $line_no**: \`$func_name\` $is_static$is_export" >> "$OUTPUT_FILE"
            echo "  \`\`\`c" >> "$OUTPUT_FILE"
            echo "  $clean_line" >> "$OUTPUT_FILE"
            echo "  \`\`\`" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        done
    fi
done

# Add summary section
echo "" >> "$OUTPUT_FILE"
echo "## Summary" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "**File counts:**" >> "$OUTPUT_FILE"

total_files=0
total_functions=0
for c_file in "$H3_SRC_DIR"/*.c; do
    if [[ ! -f "$c_file" ]]; then
        continue
    fi
    filename=$(basename "$c_file")
    
    # Count functions in this file
    func_count=$(grep -c -E '^[a-zA-Z_][a-zA-Z0-9_*[:space:]]+[a-zA-Z_][a-zA-Z0-9_]*[[:space:]]*\([^)]*\)[[:space:]]*\{' "$c_file" 2>/dev/null || echo "0")
    file_total=$func_count
    
    if [[ "$file_total" -gt 0 ]]; then
        echo "- $filename: $file_total functions" >> "$OUTPUT_FILE"
        total_files=$(( $total_files + 1 ))
        total_functions=$(( $total_functions + $file_total ))
    fi
done

echo "" >> "$OUTPUT_FILE"
echo "**Total: $total_functions functions across $total_files files**" >> "$OUTPUT_FILE"

echo ""
echo "✅ Function extraction complete!"
echo "📄 Output written to: $OUTPUT_FILE"
echo "📊 Found $total_functions functions across $total_files C files"
echo ""
echo "Next steps:"
echo "1. Review the generated inventory"
echo "2. Compare against current TODO.md to identify gaps"
echo "3. Prioritize simple functions for next c2go iterations"