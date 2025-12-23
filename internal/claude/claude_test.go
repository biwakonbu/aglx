package claude

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindClaudeMd(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. None found
	if found := FindClaudeMd(tmpDir); found != "" {
		t.Errorf("expected empty string, got %q", found)
	}

	// 2. Root CLAUDE.md
	rootPath := filepath.Join(tmpDir, "CLAUDE.md")
	os.WriteFile(rootPath, []byte("root"), 0644)
	if found := FindClaudeMd(tmpDir); found != rootPath {
		t.Errorf("expected %q, got %q", rootPath, found)
	}

	// 3. .claude/CLAUDE.md (should take priority)
	claudeDir := filepath.Join(tmpDir, ".claude")
	os.Mkdir(claudeDir, 0755)
	nestedPath := filepath.Join(claudeDir, "CLAUDE.md")
	os.WriteFile(nestedPath, []byte("nested"), 0644)
	if found := FindClaudeMd(tmpDir); found != nestedPath {
		t.Errorf("expected %q, got %q", nestedPath, found)
	}
}

func TestParse(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "CLAUDE.md")

	t.Run("Plain Markdown", func(t *testing.T) {
		content := "# Hello\nWorld"
		os.WriteFile(filePath, []byte(content), 0644)
		skill, err := Parse(filePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if skill.HasFrontmatter {
			t.Error("expected no frontmatter")
		}
		if skill.Body != content {
			t.Errorf("expected body %q, got %q", content, skill.Body)
		}
	})

	t.Run("With Frontmatter", func(t *testing.T) {
		content := "---\nname: test\n---\n# Body"
		os.WriteFile(filePath, []byte(content), 0644)
		skill, err := Parse(filePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !skill.HasFrontmatter {
			t.Error("expected has frontmatter")
		}
		if skill.Frontmatter["name"] != "test" {
			t.Errorf("expected frontmatter name 'test', got %v", skill.Frontmatter["name"])
		}
		if skill.Body != "# Body" {
			t.Errorf("expected body '# Body', got %q", skill.Body)
		}
	})

	t.Run("Malformed Frontmatter", func(t *testing.T) {
		content := "---\nname: test\nNo closing"
		os.WriteFile(filePath, []byte(content), 0644)
		skill, err := Parse(filePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if skill.HasFrontmatter {
			t.Error("expected HasFrontmatter to be false for malformed input")
		}
		if skill.Body != content {
			t.Errorf("expected body to be full content, got %q", skill.Body)
		}
	})

	t.Run("Empty File", func(t *testing.T) {
		os.WriteFile(filePath, []byte(""), 0644)
		skill, err := Parse(filePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if skill.Body != "" {
			t.Errorf("expected empty body, got %q", skill.Body)
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("Small file", func(t *testing.T) {
		skill := &ClaudeSkill{BodySize: 100}
		result := Validate(skill)
		if result.HasWarnings() {
			t.Errorf("expected no warnings, got %v", result.Warnings)
		}
	})

	t.Run("Large file", func(t *testing.T) {
		skill := &ClaudeSkill{BodySize: 30000}
		result := Validate(skill)
		if !result.HasWarnings() {
			t.Error("expected warnings for large file")
		}
		found := false
		for _, w := range result.Warnings {
			if strings.Contains(w.Message, "moderately large") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected large file warning, got %v", result.Warnings)
		}
	})

	t.Run("Very large file", func(t *testing.T) {
		skill := &ClaudeSkill{BodySize: 60000}
		result := Validate(skill)
		found := false
		for _, w := range result.Warnings {
			if strings.Contains(w.Message, "very large") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected very large file warning, got %v", result.Warnings)
		}
	})

	t.Run("Empty file", func(t *testing.T) {
		skill := &ClaudeSkill{BodySize: 0}
		result := Validate(skill)
		if !result.HasWarnings() {
			t.Error("expected warning for empty file")
		}
	})
}
