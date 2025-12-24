package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/biwakonbu/aglx/internal/skill"
)

func TestCheck(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Empty dir (ParseError expected)
	result := Check(tmpDir)
	if result.ParseError == nil {
		t.Error("expected ParseError for empty dir")
	}

	// 2. Valid skill
	content := "---\nname: my-skill\ndescription: test\n---"

	// Directory must match name
	agentDir := filepath.Join(tmpDir, "my-skill")
	os.Mkdir(agentDir, 0755)
	os.WriteFile(filepath.Join(agentDir, "SKILL.md"), []byte(content), 0644)

	result = Check(agentDir)
	if result.ParseError != nil {
		t.Errorf("unexpected ParseError: %v", result.ParseError)
	}
	if result.Skill == nil {
		t.Error("expected Skill to be parsed")
	}
	if result.Skill.Name != "my-skill" {
		t.Errorf("expected name 'my-skill', got %q", result.Skill.Name)
	}
	// Default (auto) should have both results
	if result.AgentSkillsResult == nil || result.ClaudeCodeResult == nil {
		t.Error("expected both spec results in auto mode")
	}
}

func TestCheckWithOptions_AgentSkillsOnly(t *testing.T) {
	tmpDir := t.TempDir()

	content := "---\nname: test-skill\ndescription: test\n---"
	agentDir := filepath.Join(tmpDir, "test-skill")
	os.Mkdir(agentDir, 0755)
	os.WriteFile(filepath.Join(agentDir, "SKILL.md"), []byte(content), 0644)

	result := CheckWithOptions(agentDir, &CheckOptions{Spec: skill.SpecAgentSkills})
	if result.AgentSkillsResult == nil {
		t.Error("expected AgentSkillsResult")
	}
	if result.ClaudeCodeResult != nil {
		t.Error("expected no ClaudeCodeResult when spec=agent-skills")
	}
}

func TestCheckWithOptions_ClaudeCodeOnly(t *testing.T) {
	tmpDir := t.TempDir()

	content := "---\nname: test-skill\ndescription: test\n---"
	agentDir := filepath.Join(tmpDir, "test-skill")
	os.Mkdir(agentDir, 0755)
	os.WriteFile(filepath.Join(agentDir, "SKILL.md"), []byte(content), 0644)

	result := CheckWithOptions(agentDir, &CheckOptions{Spec: skill.SpecClaudeCode})
	if result.ClaudeCodeResult == nil {
		t.Error("expected ClaudeCodeResult")
	}
	if result.AgentSkillsResult != nil {
		t.Error("expected no AgentSkillsResult when spec=claude-code")
	}
}

func TestStatus_String(t *testing.T) {
	if StatusPass.String() != "PASS" {
		t.Error("StatusPass.String() failed")
	}
	if StatusFail.String() != "FAIL" {
		t.Error("StatusFail.String() failed")
	}
	if StatusWarning.String() != "WARN" {
		t.Error("StatusWarning.String() failed")
	}
	if StatusNotFound.String() != "N/A" {
		t.Error("StatusNotFound.String() failed")
	}
}
