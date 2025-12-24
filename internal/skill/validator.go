// Package skill provides types and utilities for parsing and validating Agent Skills.
package skill

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError represents a single validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult holds the result of validating a skill.
type ValidationResult struct {
	Skill    *Skill
	Errors   []ValidationError
	Warnings []ValidationError
}

// IsValid returns true if there are no validation errors.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// HasWarnings returns true if there are any validation warnings.
func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// Validate checks if a skill conforms to the Agent Skills specification.
// Uses default options (auto-detect format).
func Validate(skill *Skill) *ValidationResult {
	return ValidateWithOptions(skill, nil)
}

// ValidateWithOptions checks if a skill conforms to the specification with custom options.
func ValidateWithOptions(skill *Skill, opts *ValidationOptions) *ValidationResult {
	result := &ValidationResult{Skill: skill}

	if opts == nil {
		opts = &ValidationOptions{}
	}

	// Validate name (required)
	validateName(skill, result, opts)

	// Validate description (required)
	validateDescription(skill, result, opts)

	// Validate compatibility (optional, but has constraints if present)
	validateCompatibility(skill, result)

	// Validate name matches directory name
	validateDirectoryMatch(skill, result)

	// Validate allowed-tools (optional, experimental)
	validateAllowedTools(skill, result, opts)

	// Validate optional directories
	validateOptionalDirectories(skill, result)

	// Validate body size (warning)
	validateBodySize(skill, result, opts)

	return result
}

func validateName(skill *Skill, result *ValidationResult, opts *ValidationOptions) {
	name := skill.Name

	if name == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Message: "is required",
		})
		return
	}

	// Length check: 1-64 characters
	if len(name) > 64 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("must be 1-64 characters (got %d)", len(name)),
		})
	}

	// Check for uppercase characters
	for _, r := range name {
		if unicode.IsUpper(r) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "name",
				Message: "must be lowercase (uppercase characters not allowed)",
			})
			break
		}
	}

	// Check for invalid characters (only lowercase alphanumeric and hyphens allowed)
	for _, r := range name {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '-' {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "name",
				Message: "may only contain lowercase alphanumeric characters (a-z, 0-9) and hyphens (-)",
			})
			break
		}
	}

	// Check for leading/trailing hyphens
	if strings.HasPrefix(name, "-") {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Message: "must not start with a hyphen",
		})
	}
	if strings.HasSuffix(name, "-") {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Message: "must not end with a hyphen",
		})
	}

	// Check for consecutive hyphens
	if strings.Contains(name, "--") {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Message: "must not contain consecutive hyphens (--)",
		})
	}

	// Claude Code specific validations
	if opts.Spec == SpecClaudeCode {
		// Check for XML tags (error)
		if containsXMLTags(name) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "name",
				Message: "must not contain XML tags",
			})
		}

		// Check for reserved words: "anthropic", "claude" (error)
		lowerName := strings.ToLower(name)
		if strings.Contains(lowerName, "anthropic") || strings.Contains(lowerName, "claude") {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "name",
				Message: "must not contain reserved words 'anthropic' or 'claude'",
			})
		}
	}
}

func validateDescription(skill *Skill, result *ValidationResult, opts *ValidationOptions) {
	desc := skill.Description

	if desc == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "description",
			Message: "is required",
		})
		return
	}

	// Length check: 1-1024 characters
	if len(desc) > 1024 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "description",
			Message: fmt.Sprintf("must be 1-1024 characters (got %d)", len(desc)),
		})
	}

	// Claude Code specific validations
	if opts.Spec == SpecClaudeCode {
		// Check for XML tags (error)
		if containsXMLTags(desc) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "description",
				Message: "must not contain XML tags",
			})
		}
	}
}

func validateCompatibility(skill *Skill, result *ValidationResult) {
	compat := skill.Compatibility

	if compat == "" {
		return // Optional field
	}

	// Length check: 1-500 characters
	if len(compat) > 500 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "compatibility",
			Message: fmt.Sprintf("must be 1-500 characters (got %d)", len(compat)),
		})
	}
}

func validateDirectoryMatch(skill *Skill, result *ValidationResult) {
	if skill.Path == "" || skill.Name == "" {
		return // Can't validate without path or name
	}

	dirName := filepath.Base(skill.Path)
	if dirName != skill.Name {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("must match parent directory name (expected %q, got %q)", dirName, skill.Name),
		})
	}
}

// toolPattern matches valid tool names: letters, digits, hyphens, and underscores.
// It can optionally include arguments in parentheses like ToolName(arg:*)
// Examples: Read, Bash(git:*), mcp__figma-desktop, mcp__chrome-devtools
var toolPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*(\(.*\))?$`)

func validateAllowedTools(skill *Skill, result *ValidationResult, opts *ValidationOptions) {
	if skill.AllowedTools == "" {
		return
	}

	// Check format based on specification
	if opts.Spec != SpecAuto {
		isCommaFormat := strings.Contains(skill.AllowedTools, ", ") || strings.Contains(skill.AllowedTools, ",")
		isSpaceFormat := !isCommaFormat

		switch opts.Spec {
		case SpecAgentSkills:
			if isCommaFormat {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "allowed-tools",
					Message: "must use space-separated format for Agent Skills specification (e.g., 'Read Glob Grep')",
				})
				return
			}
		case SpecClaudeCode:
			if isSpaceFormat && len(skill.ParsedAllowedTools()) > 1 {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "allowed-tools",
					Message: "must use comma-separated format for Claude Code specification (e.g., 'Read, Grep, Glob')",
				})
				return
			}
		}
	}

	// Validate individual tool names
	tools := skill.ParsedAllowedTools()
	for _, tool := range tools {
		if !toolPattern.MatchString(tool) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "allowed-tools",
				Message: fmt.Sprintf("invalid tool format: %q (must be alphanumeric or ToolName(args))", tool),
			})
		}
	}
}

func validateOptionalDirectories(skill *Skill, result *ValidationResult) {
	if skill.Path == "" {
		return
	}

	optionalDirs := []string{"scripts", "assets", "references"}
	for _, dir := range optionalDirs {
		dirPath := filepath.Join(skill.Path, dir)
		info, err := os.Stat(dirPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			// Other error, maybe permission?
			continue
		}

		if !info.IsDir() {
			result.Errors = append(result.Errors, ValidationError{
				Field:   dir,
				Message: "must be a directory if present (found a file)",
			})
			continue
		}

		// Check if directory is empty
		f, err := os.Open(dirPath)
		if err != nil {
			continue
		}
		defer f.Close()

		_, err = f.Readdirnames(1)
		if err == io.EOF {
			result.Errors = append(result.Errors, ValidationError{
				Field:   dir,
				Message: "must not be empty if present",
			})
		}
		f.Close() // Close early as we're in a loop

		// Check for hidden files (warning)
		checkForHiddenFiles(dirPath, dir, result)
	}
}

func checkForHiddenFiles(dirPath, rootField string, result *ValidationResult) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   rootField,
				Message: fmt.Sprintf("contains hidden file or directory: %q", f.Name()),
			})
		}
	}
}

const (
	// MaxBodyTokensRecommended is the recommended maximum body size in tokens.
	// 1 token is roughly 4 characters.
	MaxBodyTokensRecommended = 5000
	MaxBodyCharsRecommended  = MaxBodyTokensRecommended * 4

	// MaxBodyLinesClaudeCode is the recommended maximum body size in lines for Claude Code.
	MaxBodyLinesClaudeCode = 500
)

func validateBodySize(skill *Skill, result *ValidationResult, opts *ValidationOptions) {
	// Check token count (Agent Skills recommendation)
	if len(skill.Body) > MaxBodyCharsRecommended {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:   "body",
			Message: fmt.Sprintf("is very large (approximately %d tokens), recommendation is to keep it under 5000 tokens", len(skill.Body)/4),
		})
	}

	// Claude Code specific: check line count
	if opts.Spec == SpecClaudeCode {
		lineCount := strings.Count(skill.Body, "\n") + 1
		if lineCount > MaxBodyLinesClaudeCode {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "body",
				Message: fmt.Sprintf("exceeds recommended %d lines (got %d lines), consider splitting into separate files", MaxBodyLinesClaudeCode, lineCount),
			})
		}
	}
}

// ValidateMultiple validates multiple skills and returns all results.
func ValidateMultiple(skills []*Skill) []*ValidationResult {
	var results []*ValidationResult
	for _, skill := range skills {
		results = append(results, Validate(skill))
	}
	return results
}

// containsXMLTags checks if the string contains XML-like tags.
var xmlTagPattern = regexp.MustCompile(`<[a-zA-Z][^>]*>`)

func containsXMLTags(s string) bool {
	return xmlTagPattern.MatchString(s)
}
