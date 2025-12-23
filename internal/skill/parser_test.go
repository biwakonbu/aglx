package skill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse_ValidSkill(t *testing.T) {
	skill, err := Parse("../../testdata/valid/pdf-processing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if skill.Name != "pdf-processing" {
		t.Errorf("expected name 'pdf-processing', got %q", skill.Name)
	}

	if skill.Description == "" {
		t.Error("expected non-empty description")
	}

	if skill.License != "Apache-2.0" {
		t.Errorf("expected license 'Apache-2.0', got %q", skill.License)
	}

	if skill.Body == "" {
		t.Error("expected non-empty body")
	}
}

func TestParse_SimpleSkill(t *testing.T) {
	skill, err := Parse("../../testdata/valid/simple-skill")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if skill.Name != "simple-skill" {
		t.Errorf("expected name 'simple-skill', got %q", skill.Name)
	}

	if skill.License != "" {
		t.Errorf("expected empty license for simple skill, got %q", skill.License)
	}
}

func TestParse_WithMetadata(t *testing.T) {
	skill, err := Parse("../../testdata/valid/with-metadata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if skill.Name != "with-metadata" {
		t.Errorf("expected name 'with-metadata', got %q", skill.Name)
	}

	if skill.License != "MIT" {
		t.Errorf("expected license 'MIT', got %q", skill.License)
	}

	if skill.Compatibility != "Requires Python 3.9+ and internet access" {
		t.Errorf("unexpected compatibility: %q", skill.Compatibility)
	}

	if skill.AllowedTools != "Bash(python:*) Read Write" {
		t.Errorf("unexpected allowed-tools: %q", skill.AllowedTools)
	}

	// Check metadata
	if skill.Metadata == nil {
		t.Fatal("expected non-nil metadata")
	}
	if skill.Metadata["author"] != "test-author" {
		t.Errorf("expected author 'test-author', got %q", skill.Metadata["author"])
	}
	if skill.Metadata["version"] != "2.0" {
		t.Errorf("expected version '2.0', got %q", skill.Metadata["version"])
	}

	// Check ParsedAllowedTools
	tools := skill.ParsedAllowedTools()
	if len(tools) != 3 {
		t.Errorf("expected 3 tools, got %d: %v", len(tools), tools)
	}
}

func TestParse_MissingSkillMd(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := Parse(tmpDir)
	if err == nil {
		t.Error("expected error for missing SKILL.md")
	}
	if !strings.Contains(err.Error(), "SKILL.md not found") {
		t.Errorf("expected 'SKILL.md not found' error, got: %v", err)
	}
}

func TestParse_MissingFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte("# No frontmatter"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Parse(tmpDir)
	if err == nil {
		t.Error("expected error for missing frontmatter delimiter")
	}
}

func TestParse_UnclosedFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	content := `---
name: test
description: test
`
	if err := os.WriteFile(skillPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Parse(tmpDir)
	if err == nil {
		t.Error("expected error for unclosed frontmatter")
	}
	if !strings.Contains(err.Error(), "missing closing frontmatter") {
		t.Errorf("expected 'missing closing frontmatter' error, got: %v", err)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	_, err := Parse("../../testdata/invalid/invalid-yaml")
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
	if !strings.Contains(err.Error(), "YAML") {
		t.Errorf("expected YAML error, got: %v", err)
	}
}

func TestParse_NoFrontmatter(t *testing.T) {
	_, err := Parse("../../testdata/invalid/no-frontmatter")
	if err == nil {
		t.Error("expected error for missing frontmatter")
	}
}

func TestParse_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Parse(tmpDir)
	if err == nil {
		t.Error("expected error for empty file")
	}
}

func TestParseMultiple(t *testing.T) {
	skills, errs := ParseMultiple([]string{
		"../../testdata/valid/pdf-processing",
		"../../testdata/valid/simple-skill",
		"../../testdata/valid/with-metadata",
	})

	if len(skills) != 3 {
		t.Errorf("expected 3 skills, got %d", len(skills))
	}

	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %v", len(errs), errs)
	}
}

func TestParseMultiple_WithErrors(t *testing.T) {
	skills, errs := ParseMultiple([]string{
		"../../testdata/valid/pdf-processing",
		"../../testdata/invalid/invalid-yaml",
		"../../testdata/invalid/no-frontmatter",
	})

	if len(skills) != 1 {
		t.Errorf("expected 1 valid skill, got %d", len(skills))
	}

	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}

func TestParsedAllowedTools_Empty(t *testing.T) {
	skill := &Skill{AllowedTools: ""}
	tools := skill.ParsedAllowedTools()
	if tools != nil {
		t.Errorf("expected nil for empty allowed-tools, got %v", tools)
	}
}

func TestParsedAllowedTools_SingleTool(t *testing.T) {
	skill := &Skill{AllowedTools: "Read"}
	tools := skill.ParsedAllowedTools()
	if len(tools) != 1 || tools[0] != "Read" {
		t.Errorf("expected ['Read'], got %v", tools)
	}
}

func TestParsedAllowedTools_MultipleTools(t *testing.T) {
	skill := &Skill{AllowedTools: "Bash(git:*) Read Write"}
	tools := skill.ParsedAllowedTools()
	expected := []string{"Bash(git:*)", "Read", "Write"}
	if len(tools) != len(expected) {
		t.Fatalf("expected %d tools, got %d", len(expected), len(tools))
	}
	for i, tool := range tools {
		if tool != expected[i] {
			t.Errorf("expected %q at index %d, got %q", expected[i], i, tool)
		}
	}
}
