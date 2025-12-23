// Package errors provides custom error types for the aglx tool.
package errors

import "fmt"

// ExitCode represents the exit code for the CLI.
type ExitCode int

const (
	ExitSuccess         ExitCode = 0
	ExitValidationError ExitCode = 1
	ExitParseError      ExitCode = 2
	ExitUsageError      ExitCode = 64
)

// CLIError represents an error with an associated exit code.
type CLIError struct {
	Message  string
	ExitCode ExitCode
}

func (e *CLIError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error.
func NewValidationError(msg string) *CLIError {
	return &CLIError{
		Message:  msg,
		ExitCode: ExitValidationError,
	}
}

// NewParseError creates a new parse error.
func NewParseError(msg string) *CLIError {
	return &CLIError{
		Message:  msg,
		ExitCode: ExitParseError,
	}
}

// NewUsageError creates a new usage error.
func NewUsageError(msg string) *CLIError {
	return &CLIError{
		Message:  fmt.Sprintf("usage error: %s", msg),
		ExitCode: ExitUsageError,
	}
}
