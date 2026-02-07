package main

import (
	"testing"
)

func TestParseSTIXData_BasicTechnique(t *testing.T) {
	data := []byte(`{
		"objects": [
			{
				"type": "attack-pattern",
				"name": "System Information Discovery",
				"description": "An adversary may attempt to get detailed information about the OS.",
				"external_references": [
					{"source_name": "mitre-attack", "external_id": "T1082", "url": "https://attack.mitre.org/techniques/T1082"}
				],
				"kill_chain_phases": [
					{"kill_chain_name": "mitre-attack", "phase_name": "discovery"}
				],
				"x_mitre_platforms": ["Windows", "Linux", "macOS"],
				"x_mitre_deprecated": false,
				"revoked": false
			}
		]
	}`)

	techniques, err := ParseSTIXData(data)
	if err != nil {
		t.Fatalf("ParseSTIXData failed: %v", err)
	}

	if len(techniques) != 1 {
		t.Fatalf("Expected 1 technique, got %d", len(techniques))
	}

	tech := techniques["T1082"]
	if tech == nil {
		t.Fatal("T1082 not found")
	}

	if tech.Name != "System Information Discovery" {
		t.Errorf("Name = %s, want System Information Discovery", tech.Name)
	}
	if len(tech.Tactics) != 1 || tech.Tactics[0] != "discovery" {
		t.Errorf("Tactics = %v, want [discovery]", tech.Tactics)
	}
	if len(tech.Platforms) != 3 {
		t.Errorf("Platforms = %v, want 3 platforms", tech.Platforms)
	}
	if len(tech.References) != 1 || tech.References[0] != "https://attack.mitre.org/techniques/T1082" {
		t.Errorf("References = %v, want [https://attack.mitre.org/techniques/T1082]", tech.References)
	}
}

func TestParseSTIXData_MultiTactic(t *testing.T) {
	data := []byte(`{
		"objects": [
			{
				"type": "attack-pattern",
				"name": "Valid Accounts",
				"description": "Adversaries may obtain credentials.",
				"external_references": [
					{"source_name": "mitre-attack", "external_id": "T1078", "url": "https://attack.mitre.org/techniques/T1078"}
				],
				"kill_chain_phases": [
					{"kill_chain_name": "mitre-attack", "phase_name": "defense-evasion"},
					{"kill_chain_name": "mitre-attack", "phase_name": "persistence"},
					{"kill_chain_name": "mitre-attack", "phase_name": "privilege-escalation"},
					{"kill_chain_name": "mitre-attack", "phase_name": "initial-access"}
				],
				"x_mitre_platforms": ["Windows", "Linux"],
				"x_mitre_deprecated": false,
				"revoked": false
			}
		]
	}`)

	techniques, err := ParseSTIXData(data)
	if err != nil {
		t.Fatalf("ParseSTIXData failed: %v", err)
	}

	tech := techniques["T1078"]
	if tech == nil {
		t.Fatal("T1078 not found")
	}

	if len(tech.Tactics) != 4 {
		t.Errorf("Expected 4 tactics, got %d: %v", len(tech.Tactics), tech.Tactics)
	}
}

func TestParseSTIXData_SkipRevoked(t *testing.T) {
	data := []byte(`{
		"objects": [
			{
				"type": "attack-pattern",
				"name": "Revoked Tech",
				"description": "test",
				"external_references": [{"source_name": "mitre-attack", "external_id": "T9999"}],
				"kill_chain_phases": [{"kill_chain_name": "mitre-attack", "phase_name": "execution"}],
				"x_mitre_platforms": ["Windows"],
				"revoked": true
			}
		]
	}`)

	techniques, err := ParseSTIXData(data)
	if err != nil {
		t.Fatalf("ParseSTIXData failed: %v", err)
	}

	if len(techniques) != 0 {
		t.Errorf("Expected 0 techniques (revoked), got %d", len(techniques))
	}
}

func TestParseSTIXData_SkipDeprecated(t *testing.T) {
	data := []byte(`{
		"objects": [
			{
				"type": "attack-pattern",
				"name": "Deprecated Tech",
				"description": "test",
				"external_references": [{"source_name": "mitre-attack", "external_id": "T9998"}],
				"kill_chain_phases": [{"kill_chain_name": "mitre-attack", "phase_name": "execution"}],
				"x_mitre_platforms": ["Windows"],
				"x_mitre_deprecated": true,
				"revoked": false
			}
		]
	}`)

	techniques, err := ParseSTIXData(data)
	if err != nil {
		t.Fatalf("ParseSTIXData failed: %v", err)
	}

	if len(techniques) != 0 {
		t.Errorf("Expected 0 techniques (deprecated), got %d", len(techniques))
	}
}

func TestParseSTIXData_SkipCloudOnly(t *testing.T) {
	data := []byte(`{
		"objects": [
			{
				"type": "attack-pattern",
				"name": "Cloud Only",
				"description": "test",
				"external_references": [{"source_name": "mitre-attack", "external_id": "T9997"}],
				"kill_chain_phases": [{"kill_chain_name": "mitre-attack", "phase_name": "initial-access"}],
				"x_mitre_platforms": ["Azure AD", "Office 365", "SaaS"],
				"revoked": false
			}
		]
	}`)

	techniques, err := ParseSTIXData(data)
	if err != nil {
		t.Fatalf("ParseSTIXData failed: %v", err)
	}

	if len(techniques) != 0 {
		t.Errorf("Expected 0 techniques (cloud-only), got %d", len(techniques))
	}
}

func TestParseSTIXData_SkipNonAttackPattern(t *testing.T) {
	data := []byte(`{
		"objects": [
			{
				"type": "malware",
				"name": "Not an attack pattern"
			},
			{
				"type": "attack-pattern",
				"name": "Valid Tech",
				"description": "test",
				"external_references": [{"source_name": "mitre-attack", "external_id": "T1001"}],
				"kill_chain_phases": [{"kill_chain_name": "mitre-attack", "phase_name": "command-and-control"}],
				"x_mitre_platforms": ["Windows"],
				"revoked": false
			}
		]
	}`)

	techniques, err := ParseSTIXData(data)
	if err != nil {
		t.Fatalf("ParseSTIXData failed: %v", err)
	}

	if len(techniques) != 1 {
		t.Errorf("Expected 1 technique, got %d", len(techniques))
	}
}

func TestParseSTIXData_InvalidJSON(t *testing.T) {
	_, err := ParseSTIXData([]byte(`invalid json`))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestNormalizePlatforms(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  int
	}{
		{"all supported", []string{"Windows", "Linux", "macOS"}, 3},
		{"mixed case", []string{"WINDOWS", "linux", "MacOS"}, 3},
		{"cloud filtered", []string{"Azure AD", "Windows", "Office 365"}, 1},
		{"empty", []string{}, 0},
		{"duplicates", []string{"Windows", "windows", "WINDOWS"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizePlatforms(tt.input)
			if len(got) != tt.want {
				t.Errorf("normalizePlatforms(%v) = %d platforms, want %d", tt.input, len(got), tt.want)
			}
		})
	}
}

func TestParseSTIX_FileNotFound(t *testing.T) {
	_, err := ParseSTIX("/nonexistent/path.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}
