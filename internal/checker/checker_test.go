package checker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheck(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Empty dir (N/A for both)
	result := Check(tmpDir)
	if result.AgentStatus != StatusNotFound {
		t.Errorf("expected Agent NotFound, got %v", result.AgentStatus)
	}
	if result.ClaudeStatus != StatusNotFound {
		t.Errorf("expected Claude NotFound, got %v", result.ClaudeStatus)
	}

	// 2. Only Agent Skills
	agentPath := filepath.Join(tmpDir, "SKILL.md")
	content := "---\nname: my-skill\ndescription: test\n---"
	os.WriteFile(agentPath, []byte(content), 0644)

	// Directory must match name
	agentDir := filepath.Join(tmpDir, "my-skill")
	os.Mkdir(agentDir, 0755)
	os.WriteFile(filepath.Join(agentDir, "SKILL.md"), []byte(content), 0644)

	result = Check(agentDir)
	if result.AgentStatus != StatusPass {
		t.Errorf("expected Agent Pass, got %v (ParseError: %v)", result.AgentStatus, result.AgentParseError)
	}
	if result.ClaudeStatus != StatusNotFound {
		t.Errorf("expected Claude NotFound, got %v", result.ClaudeStatus)
	}

	// 3. Only Claude Skills
	claudeDir := filepath.Join(tmpDir, "project")
	os.Mkdir(claudeDir, 0755)
	os.Mkdir(filepath.Join(claudeDir, ".claude"), 0755)
	os.WriteFile(filepath.Join(claudeDir, ".claude", "CLAUDE.md"), []byte("# Hello"), 0644)

	result = Check(claudeDir)
	if result.AgentStatus != StatusNotFound {
		t.Errorf("expected Agent NotFound, got %v", result.AgentStatus)
	}
	if result.ClaudeStatus != StatusPass {
		t.Errorf("expected Claude Pass, got %v", result.ClaudeStatus)
	}

	// 4. Both
	os.WriteFile(filepath.Join(agentDir, "CLAUDE.md"), []byte("# Hello"), 0644)
	result = Check(agentDir)
	if result.AgentStatus != StatusPass {
		t.Errorf("expected Agent Pass, got %v", result.AgentStatus)
	}
	if result.ClaudeStatus != StatusPass {
		t.Errorf("expected Claude Pass, got %v", result.ClaudeStatus)
	}
}

func TestStatus_String(t *testing.T) {
	if StatusPass.String() != "PASS" {
		t.Error("StatusPass.String() failed")
	}
}
