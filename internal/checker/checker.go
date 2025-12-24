// Package checker provides a skill validator for Agent Skills and Claude Code Skills.
package checker

import (
	"github.com/biwakonbu/aglx/internal/skill"
)

// Status represents the validation status.
type Status int

const (
	StatusPass Status = iota
	StatusFail
	StatusNotFound
	StatusWarning
)

func (s Status) String() string {
	switch s {
	case StatusPass:
		return "PASS"
	case StatusFail:
		return "FAIL"
	case StatusNotFound:
		return "N/A"
	case StatusWarning:
		return "WARN"
	default:
		return "UNKNOWN"
	}
}

// CheckOptions configures the validation behavior.
type CheckOptions struct {
	// Spec specifies which specification to validate against for SKILL.md.
	// SpecAuto (default): validates against both specifications
	// SpecAgentSkills: validates against agentskills.io specification only
	// SpecClaudeCode: validates against Claude Code specification only
	Spec skill.Spec
}

// SpecResult holds the validation result for a single specification.
type SpecResult struct {
	ValidationResult *skill.ValidationResult
	Status           Status
}

// Result holds the validation result for SKILL.md.
type Result struct {
	Path string

	// Parsed skill (shared between specs)
	Skill      *skill.Skill
	ParseError error

	// Results per specification
	AgentSkillsResult *SpecResult
	ClaudeCodeResult  *SpecResult
}

// Check validates SKILL.md in the given directory.
func Check(dirPath string) *Result {
	return CheckWithOptions(dirPath, nil)
}

// CheckWithOptions validates SKILL.md with custom options.
func CheckWithOptions(dirPath string, opts *CheckOptions) *Result {
	result := &Result{
		Path: dirPath,
	}

	if opts == nil {
		opts = &CheckOptions{}
	}

	// Parse SKILL.md
	parsedSkill, err := skill.Parse(dirPath)
	if err != nil {
		result.ParseError = err
		return result
	}

	result.Skill = parsedSkill

	// Validate based on spec option
	switch opts.Spec {
	case skill.SpecAuto:
		// Validate against both specifications
		result.AgentSkillsResult = validateWithSpec(parsedSkill, skill.SpecAgentSkills)
		result.ClaudeCodeResult = validateWithSpec(parsedSkill, skill.SpecClaudeCode)
	case skill.SpecAgentSkills:
		result.AgentSkillsResult = validateWithSpec(parsedSkill, skill.SpecAgentSkills)
	case skill.SpecClaudeCode:
		result.ClaudeCodeResult = validateWithSpec(parsedSkill, skill.SpecClaudeCode)
	}

	return result
}

func validateWithSpec(parsedSkill *skill.Skill, spec skill.Spec) *SpecResult {
	validationResult := skill.ValidateWithOptions(parsedSkill, &skill.ValidationOptions{
		Spec: spec,
	})

	var status Status
	if validationResult.IsValid() {
		if validationResult.HasWarnings() {
			status = StatusWarning
		} else {
			status = StatusPass
		}
	} else {
		status = StatusFail
	}

	return &SpecResult{
		ValidationResult: validationResult,
		Status:           status,
	}
}

// CheckMultiple validates multiple directories and returns all results.
func CheckMultiple(dirPaths []string) []*Result {
	return CheckMultipleWithOptions(dirPaths, nil)
}

// CheckMultipleWithOptions validates multiple directories with custom options.
func CheckMultipleWithOptions(dirPaths []string, opts *CheckOptions) []*Result {
	var results []*Result
	for _, dirPath := range dirPaths {
		results = append(results, CheckWithOptions(dirPath, opts))
	}
	return results
}
