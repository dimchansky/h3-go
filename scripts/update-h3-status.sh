#!/bin/bash

# Convenience script to update the H3 function implementation status

SCRIPT_DIR="$(dirname "$0")"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Updating H3 function implementation status..."
python3 "$SCRIPT_DIR/update-function-status.py"

echo ""
echo "Current implementation status:"
echo "============================="

# Show completed functions
completed_count=$(grep -c "\- \[x\]" "$PROJECT_ROOT/internal/c2go/h3_functions.md")
total_count=$(grep -c "^- \[" "$PROJECT_ROOT/internal/c2go/h3_functions.md")
remaining_count=$((total_count - completed_count))

echo "✅ Completed: $completed_count"
echo "⏳ Remaining: $remaining_count"
echo "📊 Total: $total_count"
echo "📈 Progress: $(python3 -c "print(f'{$completed_count/$total_count*100:.1f}%')")"

echo ""
echo "Completed functions:"
echo "===================="
grep "\- \[x\]" "$PROJECT_ROOT/internal/c2go/h3_functions.md" | sed 's/- \[x\] /✅ /'

echo ""
echo "📄 Files generated/updated:"
echo "   • internal/c2go/h3_functions.md (Public API status)"
echo "   • internal/c2go/ported_functions_report.md (Internal functions report)"
echo ""
echo "🔧 Other commands:"
echo "   • ./scripts/list-ported-functions.sh (Live ported functions list)"