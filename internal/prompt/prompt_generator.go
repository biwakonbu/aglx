// Package prompt provides utilities for generating XML prompts from Agent Skills.
package prompt

import (
	"encoding/xml"
	"fmt"

	"github.com/biwakonbu/aglx/internal/skill"
)

// AvailableSkills represents the root XML element for the skills prompt.
type AvailableSkills struct {
	XMLName xml.Name `xml:"available_skills"`
	Skills  []Skill  `xml:"skill"`
}

// Skill represents the XML element for a single skill.
type Skill struct {
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Location    string `xml:"location"`
}

// GenerateXMLPrompt takes a slice of Skills and returns an XML string.
func GenerateXMLPrompt(skills []*skill.Skill) (string, error) {
	prompt := AvailableSkills{
		Skills: make([]Skill, 0, len(skills)),
	}

	for _, s := range skills {
		prompt.Skills = append(prompt.Skills, Skill{
			Name:        s.Name,
			Description: s.Description,
			Location:    s.Path + "/SKILL.md",
		})
	}

	output, err := xml.MarshalIndent(prompt, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal XML: %w", err)
	}

	return string(output), nil
}
