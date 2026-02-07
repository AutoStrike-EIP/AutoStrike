package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteYAMLFiles_Basic(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID:        "T1082",
			Name:      "System Info",
			Tactic:    "discovery",
			Tactics:   []string{"discovery"},
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{
				{Name: "systeminfo", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 120},
			},
			IsSafe: true,
		},
	}

	result, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	if result.FilesWritten != 1 {
		t.Errorf("FilesWritten = %d, want 1", result.FilesWritten)
	}
	if result.TotalTechniques != 1 {
		t.Errorf("TotalTechniques = %d, want 1", result.TotalTechniques)
	}

	// Check file exists
	path := filepath.Join(tmpDir, "discovery.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "T1082") {
		t.Error("Output should contain technique ID T1082")
	}
	if !strings.Contains(content, "systeminfo") {
		t.Error("Output should contain command 'systeminfo'")
	}
	if !strings.Contains(content, "Auto-generated") {
		t.Error("Output should contain header comment")
	}
}

func TestWriteYAMLFiles_MultiTactic(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID:        "T1082",
			Name:      "Discovery Tech",
			Tactic:    "discovery",
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{
				{Name: "test", Type: "cmd", Platform: "windows", Command: "whoami", Timeout: 60},
			},
		},
		{
			ID:        "T1059",
			Name:      "Execution Tech",
			Tactic:    "execution",
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{
				{Name: "test", Type: "cmd", Platform: "windows", Command: "cmd /c dir", Timeout: 60},
			},
		},
	}

	result, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	if result.FilesWritten != 2 {
		t.Errorf("FilesWritten = %d, want 2", result.FilesWritten)
	}

	// Check both files exist
	if _, err := os.Stat(filepath.Join(tmpDir, "discovery.yaml")); os.IsNotExist(err) {
		t.Error("discovery.yaml not created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "execution.yaml")); os.IsNotExist(err) {
		t.Error("execution.yaml not created")
	}
}

func TestWriteYAMLFiles_MultiTacticTechnique(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID:        "T1078",
			Name:      "Valid Accounts",
			Tactic:    "initial-access",
			Tactics:   []string{"initial-access", "persistence"},
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{
				{Name: "test", Type: "cmd", Platform: "windows", Command: "net user", Timeout: 60},
			},
		},
	}

	_, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "initial-access.yaml"))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	content := string(data)
	// Should have tactics array since it's multi-tactic
	if !strings.Contains(content, "tactics:") {
		t.Error("Multi-tactic technique should have tactics field in YAML")
	}
}

func TestWriteYAMLFiles_EmptyTechniques(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := WriteYAMLFiles(nil, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	if result.FilesWritten != 0 {
		t.Errorf("FilesWritten = %d, want 0", result.FilesWritten)
	}
}

func TestTacticFilename(t *testing.T) {
	tests := []struct {
		tactic string
		want   string
	}{
		{"discovery", "discovery.yaml"},
		{"initial-access", "initial-access.yaml"},
		{"command-and-control", "command-and-control.yaml"},
		{"EXECUTION", "execution.yaml"},
	}

	for _, tt := range tests {
		got := tacticFilename(tt.tactic)
		if got != tt.want {
			t.Errorf("tacticFilename(%s) = %s, want %s", tt.tactic, got, tt.want)
		}
	}
}

func TestToYAMLTechnique(t *testing.T) {
	tech := &MergedTechnique{
		ID:          "T1082",
		Name:        "System Info",
		Description: "test description",
		Tactic:      "discovery",
		Tactics:     []string{"discovery", "collection"},
		Platforms:   []string{"windows"},
		Executors: []MergedExecutor{
			{Name: "test", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 120},
		},
		References: []string{"https://example.com"},
		IsSafe:     true,
	}

	yt := toYAMLTechnique(tech)

	if yt.ID != "T1082" {
		t.Errorf("ID = %s, want T1082", yt.ID)
	}
	if len(yt.Tactics) != 2 {
		t.Errorf("Tactics length = %d, want 2 (multi-tactic)", len(yt.Tactics))
	}
	if len(yt.Executors) != 1 {
		t.Errorf("Executors length = %d, want 1", len(yt.Executors))
	}
}

func TestToYAMLTechnique_SingleTactic(t *testing.T) {
	tech := &MergedTechnique{
		ID:      "T1082",
		Tactic:  "discovery",
		Tactics: []string{"discovery"}, // Single tactic
		Executors: []MergedExecutor{
			{Name: "test", Type: "cmd", Command: "whoami", Timeout: 60},
		},
	}

	yt := toYAMLTechnique(tech)

	// Single tactic should NOT include tactics array
	if yt.Tactics != nil {
		t.Errorf("Single-tactic technique should have nil Tactics, got %v", yt.Tactics)
	}
}

func TestWriteYAMLFiles_DeterministicOrder(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID: "T1083", Name: "File Discovery", Tactic: "discovery",
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{{Name: "t", Type: "cmd", Platform: "windows", Command: "dir", Timeout: 60}},
		},
		{
			ID: "T1057", Name: "Process Discovery", Tactic: "discovery",
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{{Name: "t", Type: "cmd", Platform: "windows", Command: "tasklist", Timeout: 60}},
		},
		{
			ID: "T1082", Name: "System Info", Tactic: "discovery",
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{{Name: "t", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 60}},
		},
	}

	_, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "discovery.yaml"))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	content := string(data)
	// T1057 should come before T1082 which should come before T1083
	idx1057 := strings.Index(content, "T1057")
	idx1082 := strings.Index(content, "T1082")
	idx1083 := strings.Index(content, "T1083")

	if idx1057 > idx1082 || idx1082 > idx1083 {
		t.Error("Techniques should be sorted by ID")
	}
}
