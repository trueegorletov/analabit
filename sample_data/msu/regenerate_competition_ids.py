#!/usr/bin/env python3
import json
import re
from typing import Dict, Optional

def extract_competition_id(url: str) -> str:
    """Extract competition ID from URL (last number after final slash)"""
    if not url:
        return ""
    match = re.search(r'/applicants/([0-9]+)$', url)
    return match.group(1) if match else ""

def generate_competition_ids_go_code():
    """Generate Go code for competition IDs with struct format"""
    
    # Read the JSON file
    with open('/home/yegor/Documents/Prest/analabit/sample_data/msu/msu_programs_data.json', 'r', encoding='utf-8') as f:
        programs = json.load(f)
    
    print("Found", len(programs), "programs in JSON file")
    
    # Generate Go struct definition and function
    go_code = '''package resolver

// CompetitionIDs represents the competition list IDs for different quota types
type CompetitionIDs struct {
	RegularBVI      string `json:"regular_bvi"`      // Regular and BVI (льготники) quota
	DedicatedQuota  string `json:"dedicated_quota"`  // Dedicated quota (целевая квота)
	SpecialQuota    string `json:"special_quota"`    // Special quota (особая квота)
	TargetQuota     string `json:"target_quota"`     // Target quota (квота для поступающих по договорам о целевом обучении)
}

// getMSUCompetitionIDs returns the mapping of MSU program names to their competition list IDs
// for different quota types. These IDs were extracted from the MSU admissions data.
func getMSUCompetitionIDs() map[string]CompetitionIDs {
	return map[string]CompetitionIDs{
'''
    
    # Process each program
    processed_programs = set()
    for program in programs:
        name = program['name']
        
        # Skip duplicates (some programs appear multiple times with different profiles)
        if name in processed_programs:
            continue
        processed_programs.add(name)
        
        # Extract competition IDs
        regular_bvi = extract_competition_id(program.get('regular_bvi_url', ''))
        dedicated_quota = extract_competition_id(program.get('dedicated_quota_url', ''))
        special_quota = extract_competition_id(program.get('special_quota_url', ''))
        target_quota = extract_competition_id(program.get('target_quota_url', ''))
        
        # Generate Go map entry
        go_code += f'\t\t"{name}": {{\n'
        go_code += f'\t\t\tRegularBVI:     "{regular_bvi}",\n'
        go_code += f'\t\t\tDedicatedQuota: "{dedicated_quota}",\n'
        go_code += f'\t\t\tSpecialQuota:   "{special_quota}",\n'
        go_code += f'\t\t\tTargetQuota:    "{target_quota}",\n'
        go_code += f'\t\t}},\n'
    
    go_code += '''\t}
}

// GetCompetitionIDsForProgram returns competition IDs for a given program name (case-insensitive)
func GetCompetitionIDsForProgram(programName string) (CompetitionIDs, bool) {
	competitionMap := getMSUCompetitionIDs()
	
	// Try exact match first
	if ids, exists := competitionMap[programName]; exists {
		return ids, true
	}
	
	// Try case-insensitive match
	upperProgramName := strings.ToUpper(programName)
	for name, ids := range competitionMap {
		if strings.ToUpper(name) == upperProgramName {
			return ids, true
		}
	}
	
	return CompetitionIDs{}, false
}
'''
    
    print(f"Processed {len(processed_programs)} unique programs")
    
    # Write the generated Go code to file
    with open('/home/yegor/Documents/Prest/analabit/service/idmsu/resolver/msu_competition_ids_new.go', 'w', encoding='utf-8') as f:
        f.write(go_code)
    
    print("Generated new competition IDs file: msu_competition_ids_new.go")
    
    # Show some statistics
    print("\nSample programs processed:")
    sample_programs = list(processed_programs)[:10]
    for prog in sample_programs:
        print(f"  - {prog}")
    
    return processed_programs

if __name__ == "__main__":
    programs = generate_competition_ids_go_code()
    print(f"\nTotal unique programs: {len(programs)}")