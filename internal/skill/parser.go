// Package skill provides types and utilities for parsing and validating Agent Skills.
package skill

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	frontmatterDelimiter = "---"
	skillFileName        = "SKILL.md"
)

// Parse reads and parses a SKILL.md file from the given directory path.
func Parse(dirPath string) (*Skill, error) {
	skillPath := filepath.Join(dirPath, skillFileName)

	file, err := os.Open(skillPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("SKILL.md not found in %s", dirPath)
		}
		return nil, fmt.Errorf("failed to open SKILL.md: %w", err)
	}
	defer file.Close()

	frontmatter, body, err := extractFrontmatter(file)
	if err != nil {
		return nil, fmt.Errorf("failed to extract frontmatter: %w", err)
	}

	var skill Skill
	if err := yaml.Unmarshal([]byte(frontmatter), &skill); err != nil {
		return nil, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	skill.Body = body
	skill.Path = dirPath

	return &skill, nil
}

// extractFrontmatter separates YAML frontmatter from Markdown body.
// Returns frontmatter content (without delimiters) and body content.
func extractFrontmatter(file *os.File) (string, string, error) {
	scanner := bufio.NewScanner(file)

	// Check for opening delimiter
	if !scanner.Scan() {
		return "", "", fmt.Errorf("empty file")
	}
	firstLine := strings.TrimSpace(scanner.Text())
	if firstLine != frontmatterDelimiter {
		return "", "", fmt.Errorf("missing opening frontmatter delimiter (---)")
	}

	// Read frontmatter lines until closing delimiter
	var frontmatterLines []string
	foundClosing := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == frontmatterDelimiter {
			foundClosing = true
			break
		}
		frontmatterLines = append(frontmatterLines, line)
	}

	if !foundClosing {
		return "", "", fmt.Errorf("missing closing frontmatter delimiter (---)")
	}

	// Read the rest as body
	var bodyLines []string
	for scanner.Scan() {
		bodyLines = append(bodyLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("error reading file: %w", err)
	}

	frontmatter := strings.Join(frontmatterLines, "\n")
	body := strings.Join(bodyLines, "\n")

	// Trim leading newlines from body
	body = strings.TrimLeft(body, "\n")

	return frontmatter, body, nil
}

// ParseMultiple parses multiple skill directories and returns all parsed skills.
// It continues parsing even if some directories fail, collecting errors.
func ParseMultiple(dirPaths []string) ([]*Skill, []error) {
	var skills []*Skill
	var errors []error

	for _, dirPath := range dirPaths {
		skill, err := Parse(dirPath)
		if err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", dirPath, err))
			continue
		}
		skills = append(skills, skill)
	}

	return skills, errors
}
