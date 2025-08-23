#!/bin/bash

# Script to list all ported H3 functions grouped by source file

cd "$(dirname "$0")/../internal/c2go" || exit 1

echo "# Ported H3 Functions by Source File"
echo ""

# Extract and process ported function comments
grep -h "// Ported from H3 C:" *.go | \
  grep -v "_test.go" | \
  grep -v "_cgo.go" | \
  sed 's/^.*: \([^:]*\.c\)::\([^(]*\).*$/\1:\2/' | \
  sort | uniq | \
  awk -F: '{print $1 " " $2}' | \
  python3 -c "
import sys
from collections import defaultdict

funcs_by_file = defaultdict(list)
for line in sys.stdin:
    parts = line.strip().split(' ', 1)
    if len(parts) == 2:
        file_name, func_name = parts
        if func_name not in ['Ported', '', 'H3_EXPORT']:
            funcs_by_file[file_name].append(func_name)

total_funcs = 0
for file_name in sorted(funcs_by_file.keys()):
    if file_name != '//':
        unique_funcs = sorted(set(funcs_by_file[file_name]))
        total_funcs += len(unique_funcs)
        print(f'## {file_name} ({len(unique_funcs)} functions)')
        print()
        for func in unique_funcs:
            print(f'- {func}')
        print()

print(f'**Total ported functions: {total_funcs}**')
"