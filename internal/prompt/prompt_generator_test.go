package prompt

import (
	"strings"
	"testing"

	"github.com/biwakonbu/aglx/internal/skill"
)

func TestGenerateXMLPrompt(t *testing.T) {
	skills := []*skill.Skill{
		{
			Name:        "pdf-processing",
			Description: "PDF extraction",
			Path:        "/path/to/pdf",
		},
		{
			Name:        "data-analysis",
			Description: "Data analysis",
			Path:        "/path/to/data",
		},
	}

	got, err := GenerateXMLPrompt(skills)
	if err != nil {
		t.Fatalf("GenerateXMLPrompt() error = %v", err)
	}

	expectedParts := []string{
		"<available_skills>",
		"<skill>",
		"<name>pdf-processing</name>",
		"<description>PDF extraction</description>",
		"<location>/path/to/pdf/SKILL.md</location>",
		"<name>data-analysis</name>",
		"</skill>",
		"</available_skills>",
	}

	for _, part := range expectedParts {
		if !strings.Contains(got, part) {
			t.Errorf("GenerateXMLPrompt() missing part %q, got:\n%s", part, got)
		}
	}
}
