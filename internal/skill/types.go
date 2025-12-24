// Package skill provides types and utilities for parsing and validating Agent Skills.
package skill

import "strings"

// Spec represents the specification to validate against.
type Spec string

const (
	// SpecAuto automatically detects the format (allows both).
	SpecAuto Spec = ""
	// SpecAgentSkills validates against agentskills.io specification (space-separated allowed-tools).
	SpecAgentSkills Spec = "agent-skills"
	// SpecClaudeCode validates against Claude Code specification (comma-separated allowed-tools).
	SpecClaudeCode Spec = "claude-code"
)

// ValidationOptions configures additional validation behavior.
type ValidationOptions struct {
	// Spec specifies which specification to validate against.
	// If empty, both formats are accepted.
	Spec Spec
}

// Skill represents a parsed SKILL.md file.
type Skill struct {
	// Name is the skill identifier (required).
	// Must be 1-64 characters, lowercase alphanumeric and hyphens only.
	// Must match the parent directory name.
	Name string `yaml:"name"`

	// Description describes what the skill does and when to use it (required).
	// Must be 1-1024 characters.
	Description string `yaml:"description"`

	// License specifies the license applied to the skill (optional).
	License string `yaml:"license,omitempty"`

	// Compatibility describes environment requirements (optional).
	// Must be 1-500 characters if provided.
	Compatibility string `yaml:"compatibility,omitempty"`

	// AllowedTools is a space-delimited list of pre-approved tools (optional, experimental).
	AllowedTools string `yaml:"allowed-tools,omitempty"`

	// Metadata is a map for additional properties (optional).
	Metadata map[string]string `yaml:"metadata,omitempty"`

	// Body is the Markdown content after the frontmatter (not from YAML).
	Body string `yaml:"-"`

	// Path is the directory path containing this skill.
	Path string `yaml:"-"`
}

// ParsedAllowedTools returns the allowed-tools as a slice of strings.
// It supports both formats:
// - Agent Skills (agentskills.io): space-separated, e.g., "Bash(ls -la) Read"
// - Claude Code: comma-separated, e.g., "Read, Grep, Glob"
// It handles spaces within parentheses, e.g., "Bash(ls -la) Read" will return ["Bash(ls -la)", "Read"].
func (s *Skill) ParsedAllowedTools() []string {
	if s.AllowedTools == "" {
		return nil
	}

	// Detect format: if it contains comma followed by space or end, it's Claude Code format
	isCommaFormat := strings.Contains(s.AllowedTools, ", ") || strings.HasSuffix(s.AllowedTools, ",")

	if isCommaFormat {
		// Claude Code format: comma-separated
		parts := strings.Split(s.AllowedTools, ",")
		var tools []string
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				tools = append(tools, trimmed)
			}
		}
		return tools
	}

	// Agent Skills format: space-separated (with parentheses handling)
	var tools []string
	var current []rune
	inParens := false

	for _, c := range s.AllowedTools {
		switch c {
		case '(', '[', '{':
			inParens = true
			current = append(current, c)
		case ')', ']', '}':
			inParens = false
			current = append(current, c)
		case ' ', '\t':
			if inParens {
				current = append(current, c)
			} else {
				if len(current) > 0 {
					tools = append(tools, string(current))
					current = nil
				}
			}
		default:
			current = append(current, c)
		}
	}

	if len(current) > 0 {
		tools = append(tools, string(current))
	}

	return tools
}
