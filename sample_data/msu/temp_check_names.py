#!/usr/bin/env python3
"""
Temporary script to check if all HTTPHeadingSource PrettyNames from registry/msu/msu.go
are present as keys in msu_competition_ids.go
"""

import re
import os

def extract_pretty_names_from_registry():
    """Extract PrettyName values from registry/msu/msu.go"""
    registry_path = "/home/yegor/Documents/Prest/analabit/core/registry/msu/msu.go"
    
    if not os.path.exists(registry_path):
        print(f"Registry file not found: {registry_path}")
        return set()
    
    pretty_names = set()
    
    with open(registry_path, 'r', encoding='utf-8') as f:
        content = f.read()
        
    # Find all PrettyName field assignments
    pattern = r'PrettyName:\s*"([^"]+)"'
    matches = re.findall(pattern, content)
    
    for match in matches:
        pretty_names.add(match)
    
    return pretty_names

def extract_keys_from_competition_ids():
    """Extract keys from msu_competition_ids.go"""
    ids_path = "/home/yegor/Documents/Prest/analabit/service/idmsu/resolver/msu_competition_ids.go"
    
    if not os.path.exists(ids_path):
        print(f"Competition IDs file not found: {ids_path}")
        return set()
    
    keys = set()
    
    with open(ids_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Find all map keys
    pattern = r'"([^"]+)":\s*\{'
    matches = re.findall(pattern, content)
    
    for match in matches:
        keys.add(match)
    
    return keys

def main():
    print("Checking if all PrettyNames from registry are present in competition IDs...")
    
    registry_names = extract_pretty_names_from_registry()
    competition_keys = extract_keys_from_competition_ids()
    
    print(f"\nFound {len(registry_names)} PrettyNames in registry/msu/msu.go:")
    for name in sorted(registry_names):
        print(f"  - {name}")
    
    print(f"\nFound {len(competition_keys)} keys in msu_competition_ids.go:")
    for key in sorted(competition_keys):
        print(f"  - {key}")
    
    # Check for missing keys
    missing_in_competition = registry_names - competition_keys
    missing_in_registry = competition_keys - registry_names
    
    print("\n=== ANALYSIS ===")
    
    if missing_in_competition:
        print(f"\n❌ PrettyNames from registry NOT found in competition IDs ({len(missing_in_competition)}):")
        for name in sorted(missing_in_competition):
            print(f"  - {name}")
    else:
        print("\n✅ All PrettyNames from registry are present in competition IDs")
    
    if missing_in_registry:
        print(f"\n⚠️  Keys in competition IDs NOT found in registry ({len(missing_in_registry)}):")
        for name in sorted(missing_in_registry):
            print(f"  - {name}")
    else:
        print("\n✅ All competition ID keys are present in registry")
    
    # Case sensitivity check
    print("\n=== CASE SENSITIVITY CHECK ===")
    registry_upper = {name.upper() for name in registry_names}
    competition_upper = {key.upper() for key in competition_keys}
    
    case_mismatches = []
    for reg_name in registry_names:
        for comp_key in competition_keys:
            if reg_name.upper() == comp_key.upper() and reg_name != comp_key:
                case_mismatches.append((reg_name, comp_key))
    
    if case_mismatches:
        print(f"\n⚠️  Found {len(case_mismatches)} case mismatches:")
        for reg_name, comp_key in case_mismatches:
            print(f"  Registry: '{reg_name}' vs Competition: '{comp_key}'")
    else:
        print("\n✅ No case mismatches found")

if __name__ == "__main__":
    main()