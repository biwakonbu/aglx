// Package claude provides types and utilities for parsing and validating Claude Skills (CLAUDE.md).
package claude

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const frontmatterDelimiter = "---"

// FindClaudeMd searches for CLAUDE.md in the given directory and its .claude subdirectory.
// Returns the path if found, empty string if not found.
func FindClaudeMd(dirPath string) string {
	// Check .claude/CLAUDE.md first
	claudeDirPath := filepath.Join(dirPath, ClaudeDir, ClaudeFileName)
	if _, err := os.Stat(claudeDirPath); err == nil {
		return claudeDirPath
	}

	// Check CLAUDE.md in the root directory
	rootPath := filepath.Join(dirPath, ClaudeFileName)
	if _, err := os.Stat(rootPath); err == nil {
		return rootPath
	}

	return ""
}

// Parse reads and parses a CLAUDE.md file from the given file path.
func Parse(filePath string) (*ClaudeSkill, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("CLAUDE.md not found at %s", filePath)
		}
		return nil, fmt.Errorf("failed to open CLAUDE.md: %w", err)
	}
	defer file.Close()

	skill := &ClaudeSkill{
		Path: filePath,
	}

	// Read all content
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(lines) == 0 {
		skill.Body = ""
		skill.BodySize = 0
		return skill, nil
	}

	// Check for frontmatter
	if strings.TrimSpace(lines[0]) == frontmatterDelimiter {
		frontmatter, bodyStart, err := extractFrontmatter(lines)
		if err != nil {
			// If frontmatter extraction fails, treat entire file as body
			skill.HasFrontmatter = false
			skill.Body = strings.Join(lines, "\n")
			skill.BodySize = len(skill.Body)
			return skill, nil
		}

		skill.HasFrontmatter = true

		// Parse frontmatter as YAML
		var fm map[string]interface{}
		if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
			// If YAML parsing fails, still record that frontmatter exists
			skill.Frontmatter = make(map[string]interface{})
		} else {
			skill.Frontmatter = fm
		}

		// Set body
		if bodyStart < len(lines) {
			skill.Body = strings.TrimLeft(strings.Join(lines[bodyStart:], "\n"), "\n")
		}
	} else {
		// No frontmatter, entire file is body
		skill.HasFrontmatter = false
		skill.Body = strings.Join(lines, "\n")
	}

	skill.BodySize = len(skill.Body)
	return skill, nil
}

// extractFrontmatter extracts YAML frontmatter from lines.
// Returns the frontmatter content, the line index where body starts, and any error.
func extractFrontmatter(lines []string) (string, int, error) {
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != frontmatterDelimiter {
		return "", 0, fmt.Errorf("no frontmatter delimiter at start")
	}

	// Find closing delimiter
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == frontmatterDelimiter {
			frontmatter := strings.Join(lines[1:i], "\n")
			return frontmatter, i + 1, nil
		}
	}

	return "", 0, fmt.Errorf("no closing frontmatter delimiter")
}

// ParseFromDir finds and parses CLAUDE.md from the given directory.
func ParseFromDir(dirPath string) (*ClaudeSkill, error) {
	filePath := FindClaudeMd(dirPath)
	if filePath == "" {
		return nil, nil // Not found, but not an error
	}
	return Parse(filePath)
}
