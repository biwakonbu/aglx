// Package claude provides types and utilities for parsing and validating Claude Skills (CLAUDE.md).
package claude

// ClaudeSkill represents a parsed CLAUDE.md file.
type ClaudeSkill struct {
	// Path is the file path to the CLAUDE.md file.
	Path string

	// HasFrontmatter indicates if the file has YAML frontmatter.
	HasFrontmatter bool

	// Frontmatter contains the parsed frontmatter fields if present.
	Frontmatter map[string]interface{}

	// Body is the Markdown content (after frontmatter if present).
	Body string

	// BodySize is the size of the body in bytes.
	BodySize int
}

// Locations where CLAUDE.md files can be found.
const (
	// ClaudeFileName is the standard filename for Claude skills.
	ClaudeFileName = "CLAUDE.md"

	// ClaudeDir is the directory name for project-local Claude configuration.
	ClaudeDir = ".claude"
)
