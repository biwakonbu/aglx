package skill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidate_ValidSkill(t *testing.T) {
	skill := &Skill{
		Name:        "pdf-processing",
		Description: "Extract text and tables from PDF files.",
		Path:        "/path/to/pdf-processing",
	}

	result := Validate(skill)
	if !result.IsValid() {
		t.Errorf("expected valid skill, got errors: %v", result.Errors)
	}
}

func TestValidate_MissingName(t *testing.T) {
	skill := &Skill{
		Description: "Some description",
	}

	result := Validate(skill)
	if result.IsValid() {
		t.Error("expected invalid skill for missing name")
	}

	found := false
	for _, e := range result.Errors {
		if e.Field == "name" && e.Message == "is required" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'name is required' error")
	}
}

func TestValidate_UppercaseName(t *testing.T) {
	skill := &Skill{
		Name:        "PDF-Processing",
		Description: "Some description",
	}

	result := Validate(skill)
	if result.IsValid() {
		t.Error("expected invalid skill for uppercase name")
	}

	found := false
	for _, e := range result.Errors {
		if e.Field == "name" && e.Message == "must be lowercase (uppercase characters not allowed)" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected uppercase error, got: %v", result.Errors)
	}
}

func TestValidate_NameStartsWithHyphen(t *testing.T) {
	skill := &Skill{
		Name:        "-pdf-processing",
		Description: "Some description",
	}

	result := Validate(skill)
	found := false
	for _, e := range result.Errors {
		if e.Field == "name" && e.Message == "must not start with a hyphen" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'must not start with hyphen' error")
	}
}

func TestValidate_NameEndsWithHyphen(t *testing.T) {
	skill := &Skill{
		Name:        "pdf-processing-",
		Description: "Some description",
	}

	result := Validate(skill)
	found := false
	for _, e := range result.Errors {
		if e.Field == "name" && e.Message == "must not end with a hyphen" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'must not end with hyphen' error")
	}
}

func TestValidate_ConsecutiveHyphens(t *testing.T) {
	skill := &Skill{
		Name:        "pdf--processing",
		Description: "Some description",
	}

	result := Validate(skill)
	found := false
	for _, e := range result.Errors {
		if e.Field == "name" && e.Message == "must not contain consecutive hyphens (--)" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected consecutive hyphens error")
	}
}

func TestValidate_NameTooLong(t *testing.T) {
	longName := ""
	for i := 0; i < 65; i++ {
		longName += "a"
	}
	skill := &Skill{
		Name:        longName,
		Description: "Some description",
	}

	result := Validate(skill)
	found := false
	for _, e := range result.Errors {
		if e.Field == "name" && strings.Contains(e.Message, "1-64 characters") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected name length error, got: %v", result.Errors)
	}
}

func TestValidate_NameExactly64Chars(t *testing.T) {
	// 64 characters should be valid
	name := strings.Repeat("a", 64)
	skill := &Skill{
		Name:        name,
		Description: "Some description",
		Path:        "/path/to/" + name,
	}

	result := Validate(skill)
	if !result.IsValid() {
		t.Errorf("expected 64 char name to be valid, got errors: %v", result.Errors)
	}
}

func TestValidate_MissingDescription(t *testing.T) {
	skill := &Skill{
		Name: "test-skill",
	}

	result := Validate(skill)
	found := false
	for _, e := range result.Errors {
		if e.Field == "description" && e.Message == "is required" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'description is required' error")
	}
}

func TestValidate_DescriptionTooLong(t *testing.T) {
	longDesc := strings.Repeat("a", 1025)
	skill := &Skill{
		Name:        "test-skill",
		Description: longDesc,
	}

	result := Validate(skill)
	found := false
	for _, e := range result.Errors {
		if e.Field == "description" && strings.Contains(e.Message, "1-1024 characters") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected description length error, got: %v", result.Errors)
	}
}

func TestValidate_DescriptionExactly1024Chars(t *testing.T) {
	// 1024 characters should be valid
	desc := strings.Repeat("a", 1024)
	skill := &Skill{
		Name:        "test-skill",
		Description: desc,
		Path:        "/path/to/test-skill",
	}

	result := Validate(skill)
	if !result.IsValid() {
		t.Errorf("expected 1024 char description to be valid, got errors: %v", result.Errors)
	}
}

func TestValidate_CompatibilityTooLong(t *testing.T) {
	longCompat := strings.Repeat("a", 501)
	skill := &Skill{
		Name:          "test-skill",
		Description:   "Some description",
		Compatibility: longCompat,
	}

	result := Validate(skill)
	found := false
	for _, e := range result.Errors {
		if e.Field == "compatibility" && strings.Contains(e.Message, "1-500 characters") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected compatibility length error, got: %v", result.Errors)
	}
}

func TestValidate_CompatibilityExactly500Chars(t *testing.T) {
	// 500 characters should be valid
	compat := strings.Repeat("a", 500)
	skill := &Skill{
		Name:          "test-skill",
		Description:   "Some description",
		Compatibility: compat,
		Path:          "/path/to/test-skill",
	}

	result := Validate(skill)
	if !result.IsValid() {
		t.Errorf("expected 500 char compatibility to be valid, got errors: %v", result.Errors)
	}
}

func TestValidate_CompatibilityEmpty(t *testing.T) {
	// Empty compatibility should be valid (it's optional)
	skill := &Skill{
		Name:          "test-skill",
		Description:   "Some description",
		Compatibility: "",
		Path:          "/path/to/test-skill",
	}

	result := Validate(skill)
	if !result.IsValid() {
		t.Errorf("expected empty compatibility to be valid, got errors: %v", result.Errors)
	}
}

func TestValidate_DirectoryMismatch(t *testing.T) {
	skill := &Skill{
		Name:        "pdf-processing",
		Description: "Some description",
		Path:        "/path/to/different-name",
	}

	result := Validate(skill)
	found := false
	for _, e := range result.Errors {
		if e.Field == "name" && e.Message == `must match parent directory name (expected "different-name", got "pdf-processing")` {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected directory mismatch error, got: %v", result.Errors)
	}
}

func TestValidate_DirectoryMatchWithPath(t *testing.T) {
	skill := &Skill{
		Name:        "my-skill",
		Description: "Some description",
		Path:        "/some/nested/path/to/my-skill",
	}

	result := Validate(skill)
	if !result.IsValid() {
		t.Errorf("expected skill with matching directory to be valid, got: %v", result.Errors)
	}
}

func TestValidate_NameWithNumbers(t *testing.T) {
	skill := &Skill{
		Name:        "skill-v2",
		Description: "Some description",
		Path:        "/path/to/skill-v2",
	}

	result := Validate(skill)
	if !result.IsValid() {
		t.Errorf("expected name with numbers to be valid, got: %v", result.Errors)
	}
}

func TestValidate_NameOnlyNumbers(t *testing.T) {
	skill := &Skill{
		Name:        "123",
		Description: "Some description",
		Path:        "/path/to/123",
	}

	result := Validate(skill)
	if !result.IsValid() {
		t.Errorf("expected numeric name to be valid, got: %v", result.Errors)
	}
}

func TestValidate_NameWithSpecialChars(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"skill_underscore", true}, // underscore not allowed
		{"skill.dot", true},        // dot not allowed
		{"skill@at", true},         // @ not allowed
		{"skill space", true},      // space not allowed
		{"skill-valid", false},     // hyphen allowed
		{"skill123", false},        // numbers allowed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := &Skill{
				Name:        tt.name,
				Description: "Some description",
			}
			result := Validate(skill)
			if tt.wantErr && result.IsValid() {
				t.Errorf("expected %q to be invalid", tt.name)
			}
			if !tt.wantErr && !result.IsValid() {
				t.Errorf("expected %q to be valid, got: %v", tt.name, result.Errors)
			}
		})
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	skill := &Skill{
		Name:        "Invalid--Name-",
		Description: "", // missing
	}

	result := Validate(skill)
	if result.IsValid() {
		t.Error("expected multiple errors")
	}

	// Should have at least 3 errors: uppercase, consecutive hyphens, ends with hyphen, missing description
	if len(result.Errors) < 3 {
		t.Errorf("expected at least 3 errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

func TestValidate_AllowedTools(t *testing.T) {
	tests := []struct {
		tools   string
		valid   bool
		wantErr string
	}{
		{"Bash(git:*) Bash(jq:*) Read", true, ""},
		{"Bash(ls -la) Read", true, ""},
		{"ValidTool123", true, ""},
		{"ToolWith_Underscore", false, "invalid tool format"},
		{"Tool.With.Dot", false, "invalid tool format"},
		{"Tool@At", false, "invalid tool format"},
		{"Tool-With-Hyphen", false, "invalid tool format"}, // Spec says alphanumeric for tool name
	}

	for _, tt := range tests {
		t.Run(tt.tools, func(t *testing.T) {
			skill := &Skill{
				Name:         "test-skill",
				Description:  "Description",
				AllowedTools: tt.tools,
				Path:         "/path/to/test-skill",
			}

			result := Validate(skill)
			if tt.valid {
				if !result.IsValid() {
					t.Errorf("expected valid for %q, got: %v", tt.tools, result.Errors)
				}
			} else {
				if result.IsValid() {
					t.Errorf("expected invalid for %q", tt.tools)
				} else {
					found := false
					for _, e := range result.Errors {
						if e.Field == "allowed-tools" && strings.Contains(e.Message, tt.wantErr) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected error %q for %q, got: %v", tt.wantErr, tt.tools, result.Errors)
					}
				}
			}
		})
	}
}

func TestValidate_OptionalDirectories(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Valid: script exists and is not empty
	scriptsDir := filepath.Join(tempDir, "scripts")
	err := os.Mkdir(scriptsDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(scriptsDir, "test.sh"), []byte("echo hello"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	skill := &Skill{
		Name:        "temp-skill",
		Description: "Description",
		Path:        tempDir,
	}

	// We need to match the directory name to the skill name to avoid that error
	skill.Name = filepath.Base(tempDir)

	result := Validate(skill)
	if !result.IsValid() {
		t.Errorf("expected valid for non-empty scripts dir, got: %v", result.Errors)
	}

	// 2. Invalid: assets is a file, not a directory
	assetsFile := filepath.Join(tempDir, "assets")
	err = os.WriteFile(assetsFile, []byte("not a dir"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	result = Validate(skill)
	found := false
	for _, e := range result.Errors {
		if e.Field == "assets" && strings.Contains(e.Message, "must be a directory") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error for assets as file, got: %v", result.Errors)
	}
	os.Remove(assetsFile)

	// 3. Invalid: references exists but is empty
	refsDir := filepath.Join(tempDir, "references")
	err = os.Mkdir(refsDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	result = Validate(skill)
	found = false
	for _, e := range result.Errors {
		if e.Field == "references" && strings.Contains(e.Message, "must not be empty") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error for empty references dir, got: %v", result.Errors)
	}
}

func TestValidateMultiple(t *testing.T) {
	skills := []*Skill{
		{Name: "valid-skill", Description: "Valid", Path: "/path/to/valid-skill"},
		{Name: "Invalid", Description: "Invalid uppercase"},
		{Name: "also-valid", Description: "Also valid", Path: "/path/to/also-valid"},
	}

	results := ValidateMultiple(skills)
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	validCount := 0
	for _, r := range results {
		if r.IsValid() {
			validCount++
		}
	}
	if validCount != 2 {
		t.Errorf("expected 2 valid skills, got %d", validCount)
	}
}

// Integration tests using testdata directory
func TestValidate_Integration_ValidSkills(t *testing.T) {
	validPaths := []string{
		"../../testdata/valid/pdf-processing",
		"../../testdata/valid/simple-skill",
		"../../testdata/valid/with-metadata",
	}

	for _, path := range validPaths {
		t.Run(path, func(t *testing.T) {
			skill, err := Parse(path)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			result := Validate(skill)
			if !result.IsValid() {
				t.Errorf("expected valid skill, got errors: %v", result.Errors)
			}
		})
	}
}

func TestValidate_Integration_InvalidSkills(t *testing.T) {
	tests := []struct {
		path          string
		expectedField string
		expectedMsg   string
	}{
		{"../../testdata/invalid/uppercase-name", "name", "lowercase"},
		{"../../testdata/invalid/hyphen-start", "name", "start with a hyphen"},
		{"../../testdata/invalid/hyphen-end", "name", "end with a hyphen"},
		{"../../testdata/invalid/consecutive-hyphens", "name", "consecutive hyphens"},
		{"../../testdata/invalid/missing-name", "name", "required"},
		{"../../testdata/invalid/missing-description", "description", "required"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			skill, err := Parse(tt.path)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			result := Validate(skill)
			if result.IsValid() {
				t.Errorf("expected invalid skill for %s", tt.path)
				return
			}

			found := false
			for _, e := range result.Errors {
				if e.Field == tt.expectedField && strings.Contains(e.Message, tt.expectedMsg) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error containing field=%q msg=%q, got: %v",
					tt.expectedField, tt.expectedMsg, result.Errors)
			}
		})
	}
}

func TestValidate_Warnings(t *testing.T) {
	t.Run("Body size warning", func(t *testing.T) {
		// ~21000 chars should trigger warning (> 20000)
		body := strings.Repeat("A", 21000)
		skill := &Skill{
			Name:        "large-skill",
			Description: "Large body skill",
			Body:        body,
			Path:        "/path/to/large-skill",
		}
		result := Validate(skill)
		if !result.IsValid() {
			t.Errorf("expected skill to be valid (only warnings), got errors: %v", result.Errors)
		}
		if !result.HasWarnings() {
			t.Error("expected warnings for large body")
		}
		found := false
		for _, w := range result.Warnings {
			if w.Field == "body" && strings.Contains(w.Message, "very large") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected body size warning, got: %v", result.Warnings)
		}
	})

	t.Run("Hidden files warning", func(t *testing.T) {
		tmpDir := t.TempDir()
		skillDir := filepath.Join(tmpDir, "hidden-skill")
		os.Mkdir(skillDir, 0755)
		scriptsDir := filepath.Join(skillDir, "scripts")
		os.Mkdir(scriptsDir, 0755)

		// Create normal file and hidden file
		os.WriteFile(filepath.Join(scriptsDir, "script.sh"), []byte("echo hi"), 0644)
		os.WriteFile(filepath.Join(scriptsDir, ".env"), []byte("SECRET=123"), 0644)

		skill := &Skill{
			Name:        "hidden-skill",
			Description: "Skill with hidden files",
			Path:        skillDir,
		}
		result := Validate(skill)
		if !result.HasWarnings() {
			t.Error("expected warnings for hidden files")
		}
		found := false
		for _, w := range result.Warnings {
			if w.Field == "scripts" && strings.Contains(w.Message, "hidden file") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected hidden file warning, got: %v", result.Warnings)
		}
	})
}
