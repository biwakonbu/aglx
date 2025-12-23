// Package claude provides types and utilities for parsing and validating Claude Skills (CLAUDE.md).
package claude

// ValidationWarning represents a non-fatal warning during validation.
type ValidationWarning struct {
	Field   string
	Message string
}

// ValidationResult holds the result of validating a Claude skill.
type ValidationResult struct {
	Skill    *ClaudeSkill
	Warnings []ValidationWarning
}

// IsValid returns true for Claude skills (they're always "valid" since rules are loose).
// Use HasWarnings to check for issues.
func (r *ValidationResult) IsValid() bool {
	return r.Skill != nil
}

// HasWarnings returns true if there are any warnings.
func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

const (
	// RecommendedMaxBodySize is the recommended maximum body size in bytes.
	// This is a soft limit based on typical context window constraints.
	RecommendedMaxBodySize = 50000 // ~50KB, roughly 12500 tokens

	// WarningBodySize is the size at which we warn about large files.
	WarningBodySize = 20000 // ~20KB
)

// Validate checks a Claude skill and returns warnings for potential issues.
// Unlike Agent Skills, Claude Skills validation is more lenient.
func Validate(skill *ClaudeSkill) *ValidationResult {
	result := &ValidationResult{Skill: skill}

	if skill == nil {
		return result
	}

	// Warning: Large body size
	if skill.BodySize > RecommendedMaxBodySize {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "body",
			Message: "file is very large (>50KB), may impact context window usage",
		})
	} else if skill.BodySize > WarningBodySize {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "body",
			Message: "file is moderately large (>20KB), consider splitting",
		})
	}

	// Warning: Empty body
	if skill.BodySize == 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "body",
			Message: "file is empty",
		})
	}

	return result
}
