#!/usr/bin/env python3

"""
Script to update the h3_functions.md file by marking implemented functions as completed.
This script reads the ported functions from internal/c2go/*.go files and updates the 
task list in internal/c2go/h3_functions.md accordingly.
"""

import re
import os
import sys
from pathlib import Path
from collections import defaultdict

def extract_ported_functions():
    """Extract all ported function names from internal/c2go/*.go files"""
    ported_functions = set()
    c2go_dir = Path(__file__).parent.parent / "internal" / "c2go"
    
    if not c2go_dir.exists():
        print(f"Error: {c2go_dir} does not exist")
        return ported_functions
    
    # Pattern to match "// Ported from H3 C: file.c::functionName" or "file.c::H3_EXPORT(functionName)"
    pattern1 = re.compile(r'// Ported from H3 C: [^:]+\.c::([^\s(]+)')
    pattern2 = re.compile(r'// Ported from H3 C: [^:]+\.c::H3_EXPORT\(([^)]+)\)')
    
    for go_file in c2go_dir.glob("*.go"):
        # Skip test files
        if go_file.name.endswith("_test.go") or go_file.name.endswith("_cgo.go"):
            continue
            
        try:
            with open(go_file, 'r', encoding='utf-8') as f:
                content = f.read()
                
                # Match direct function names
                matches1 = pattern1.findall(content)
                for match in matches1:
                    func_name = match.strip()
                    if func_name and func_name not in ['Ported', 'H3_EXPORT']:
                        ported_functions.add(func_name)
                
                # Match H3_EXPORT(functionName) format
                matches2 = pattern2.findall(content)
                for match in matches2:
                    func_name = match.strip()
                    if func_name and func_name not in ['Ported']:
                        ported_functions.add(func_name)
        except Exception as e:
            print(f"Warning: Could not read {go_file}: {e}")
    
    return ported_functions

def create_ported_functions_report(ported_functions):
    """Create a detailed report of all ported internal functions"""
    c2go_dir = Path(__file__).parent.parent / "internal" / "c2go"
    report_file = c2go_dir / "ported_functions_report.md"
    
    # Group functions by source file
    functions_by_file = {}
    pattern1 = re.compile(r'// Ported from H3 C: ([^:]+\.c)::([^\s(]+)')
    pattern2 = re.compile(r'// Ported from H3 C: ([^:]+\.c)::H3_EXPORT\(([^)]+)\)')
    
    for go_file in c2go_dir.glob("*.go"):
        if go_file.name.endswith("_test.go") or go_file.name.endswith("_cgo.go"):
            continue
            
        try:
            with open(go_file, 'r', encoding='utf-8') as f:
                content = f.read()
                
                # Match direct function names
                matches1 = pattern1.findall(content)
                for file_name, func_name in matches1:
                    if func_name.strip() and func_name not in ['Ported', 'H3_EXPORT']:
                        if file_name not in functions_by_file:
                            functions_by_file[file_name] = []
                        functions_by_file[file_name].append(func_name.strip())
                
                # Match H3_EXPORT(functionName) format
                matches2 = pattern2.findall(content)
                for file_name, func_name in matches2:
                    if func_name.strip() and func_name not in ['Ported']:
                        if file_name not in functions_by_file:
                            functions_by_file[file_name] = []
                        functions_by_file[file_name].append(func_name.strip())
        except Exception as e:
            print(f"Warning: Could not read {go_file}: {e}")
    
    # Generate the report
    total_functions = 0
    report_content = "# Ported H3 Internal Functions Report\n\n"
    report_content += "This report lists all H3 C functions that have been ported to Go in the internal/c2go package.\n"
    report_content += "These are primarily internal/static functions used to build the public API.\n\n"
    
    for file_name in sorted(functions_by_file.keys()):
        unique_funcs = sorted(set(functions_by_file[file_name]))
        total_functions += len(unique_funcs)
        report_content += f"## {file_name} ({len(unique_funcs)} functions)\n\n"
        
        for func in unique_funcs:
            report_content += f"- {func}\n"
        report_content += "\n"
    
    report_content += f"**Total ported internal functions: {total_functions}**\n\n"
    report_content += "---\n\n"
    report_content += "*This report is automatically generated. Run `./scripts/update-h3-status.sh` to refresh.*\n"
    
    # Write the report
    try:
        with open(report_file, 'w', encoding='utf-8') as f:
            f.write(report_content)
        print(f"Generated ported functions report: {report_file}")
        return True
    except Exception as e:
        print(f"Error writing report to {report_file}: {e}")
        return False

def update_h3_functions_md():
    """Update the h3_functions.md file to mark implemented functions as completed"""
    md_file = Path(__file__).parent.parent / "internal" / "c2go" / "h3_functions.md"
    
    if not md_file.exists():
        print(f"Error: {md_file} does not exist")
        return False
    
    # Get list of ported functions
    ported_functions = extract_ported_functions()
    print(f"Found {len(ported_functions)} ported functions")
    
    # Create the ported functions report
    create_ported_functions_report(ported_functions)
    
    # Read the current markdown file
    try:
        with open(md_file, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        print(f"Error reading {md_file}: {e}")
        return False
    
    # Track statistics
    total_functions = 0
    completed_functions = 0
    
    # Pattern to match task list items: "- [ ] functionName" or "- [x] functionName"
    def replace_function(match):
        nonlocal total_functions, completed_functions
        checkbox = match.group(1)  # [ ] or [x]
        func_name = match.group(2).strip()
        total_functions += 1
        
        if func_name in ported_functions:
            completed_functions += 1
            return f"- [x] {func_name}"
        else:
            return f"- [ ] {func_name}"
    
    # Replace task list items
    pattern = re.compile(r'^- \[([ x])\] (.+)$', re.MULTILINE)
    updated_content = pattern.sub(replace_function, content)
    
    # Update the summary section
    summary_pattern = re.compile(
        r'(\*\*Total Functions:\*\* \d+\n)'
        r'- \*\*Completed:\*\* \d+\n'
        r'- \*\*Remaining:\*\* \d+',
        re.MULTILINE
    )
    
    remaining = total_functions - completed_functions
    new_summary = (
        f"**Total Functions:** {total_functions}\n"
        f"- **Completed:** {completed_functions}\n"
        f"- **Remaining:** {remaining}"
    )
    
    updated_content = summary_pattern.sub(
        lambda m: m.group(1) + new_summary.split('\n', 1)[1],
        updated_content
    )
    
    # Write the updated content back to the file
    try:
        with open(md_file, 'w', encoding='utf-8') as f:
            f.write(updated_content)
        
        print(f"Successfully updated {md_file}")
        print(f"Status: {completed_functions}/{total_functions} functions completed ({completed_functions/total_functions*100:.1f}%)")
        return True
        
    except Exception as e:
        print(f"Error writing to {md_file}: {e}")
        return False

if __name__ == "__main__":
    success = update_h3_functions_md()
    sys.exit(0 if success else 1)