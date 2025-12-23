// Package checker provides a unified checker for both Agent Skills and Claude Skills.
package checker

import (
	"github.com/biwakonbu/aglx/internal/claude"
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

// Result holds the combined validation results for both specifications.
type Result struct {
	Path string

	// Agent Skills
	AgentSkill      *skill.Skill
	AgentResult     *skill.ValidationResult
	AgentParseError error
	AgentStatus     Status

	// Claude Skills
	ClaudeSkill      *claude.ClaudeSkill
	ClaudeResult     *claude.ValidationResult
	ClaudeParseError error
	ClaudeStatus     Status
}

// Check validates both Agent Skills and Claude Skills for the given directory.
func Check(dirPath string) *Result {
	result := &Result{
		Path: dirPath,
	}

	// Check Agent Skills (SKILL.md)
	checkAgentSkills(dirPath, result)

	// Check Claude Skills (CLAUDE.md)
	checkClaudeSkills(dirPath, result)

	return result
}

func checkAgentSkills(dirPath string, result *Result) {
	agentSkill, err := skill.Parse(dirPath)
	if err != nil {
		result.AgentParseError = err
		result.AgentStatus = StatusNotFound
		return
	}

	result.AgentSkill = agentSkill
	validationResult := skill.Validate(agentSkill)
	result.AgentResult = validationResult

	if validationResult.IsValid() {
		if validationResult.HasWarnings() {
			result.AgentStatus = StatusWarning
		} else {
			result.AgentStatus = StatusPass
		}
	} else {
		result.AgentStatus = StatusFail
	}
}

func checkClaudeSkills(dirPath string, result *Result) {
	claudeSkill, err := claude.ParseFromDir(dirPath)
	if err != nil {
		result.ClaudeParseError = err
		result.ClaudeStatus = StatusNotFound
		return
	}

	if claudeSkill == nil {
		result.ClaudeStatus = StatusNotFound
		return
	}

	result.ClaudeSkill = claudeSkill
	validationResult := claude.Validate(claudeSkill)
	result.ClaudeResult = validationResult

	if validationResult.HasWarnings() {
		result.ClaudeStatus = StatusWarning
	} else {
		result.ClaudeStatus = StatusPass
	}
}

// CheckMultiple validates multiple directories and returns all results.
func CheckMultiple(dirPaths []string) []*Result {
	var results []*Result
	for _, dirPath := range dirPaths {
		results = append(results, Check(dirPath))
	}
	return results
}
