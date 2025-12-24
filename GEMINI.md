# GEMINI - Project Guidelines

This document outlines the architectural overview, development rules, and guidelines for the `aglx` project.

## Project Mission
`aglx` (Agent sKiLls eXaminer) aims to provide a robust, dual-validation tool for **Agent Skills** (`SKILL.md`) and **Claude Skills** (`CLAUDE.md`). It ensures that skill packages are properly formatted, token-efficient, and ready for discovery by AI agents.

## Development Rules

### 1. Language Policy
- **Always English**: All code, comments, documentation, and commit messages MUST be written in English.
- This ensures maximum accessibility and consistency for global development.

### 2. Go Standards
- Follow idiomatic Go patterns (Accept interfaces, return structs).
- Maintain high test coverage for all validation logic in `internal/`.
- Use `go fmt` and address all lint warnings from `staticcheck`.

### 3. Validation Philosophy
- **Errors**: Used for specification violations that would break discovery or execution.
- **Warnings**: Used for best practice recommendations (e.g., large body size, hidden files) that do not strictly violate the spec but may impact performance or security.

### 4. Knowledge Maintenance
- **Continuous Updates**: AI Agents MUST update `GEMINI.md` files (root and subdirectories) whenever significant architectural changes, new rules, or package responsibilities are added or modified.
- **Consistency**: Ensure that local `GEMINI.md` files are consistent with the root `GEMINI.md`.

---

## Technical Architecture

### Core Packages (`internal/`)
- **`skill`**: Strictly validates Agent Skills specification. Handles parsing of YAML frontmatter and directory structure verification.
- **`claude`**: Validates Claude Skills with a focus on file size warnings and structure.
- **`checker`**: Aggregates results from both validators into a unified status.
- **`prompt`**: Generates XML context snippets for agent discovery from validated skills.
- **`errors`**: Defines project-wide exit codes and common error types.

### CLI (`cmd/aglx`)
- The single entry point for the user. Uses subcommands (`validate`, `to-prompt`) to handle different workflows.
- Supports both human-readable text output and machine-readable JSON output.

---

## Guidelines for Adding Skills

1. **Directory Naming**: The directory name MUST exactly match the `name` field in the `SKILL.md` frontmatter.
2. **Body Efficiency**: Keep the `SKILL.md` body under 5000 tokens (approx. 20,000 characters) to ensure token efficiency during agent context injection.
3. **Hidden Files**: Avoid including hidden files (e.g., `.env`, `.DS_Store`) in `scripts/`, `assets/`, or `references/` directories.
4. **Verification**: Always run `aglx validate` locally before committing new skills.

## Verification Workflow
Before submitting a pull request:
1. Run all tests: `go test ./internal/...`
2. Build the tool: `go build -o aglx ./cmd/aglx`
3. Validate sample skills: `./aglx validate ./testdata/valid/*`

---

## Directory-Specific Guidelines
For more detailed information about specific components, refer to the local `GEMINI.md` files:
- [cmd/](file:///Users/biwakonbu/github/aglx/cmd/GEMINI.md): CLI implementation details.
- [internal/skill/](file:///Users/biwakonbu/github/aglx/internal/skill/GEMINI.md): Agent Skills validation.
- [internal/claude/](file:///Users/biwakonbu/github/aglx/internal/claude/GEMINI.md): Claude Skills validation.
- [internal/checker/](file:///Users/biwakonbu/github/aglx/internal/checker/GEMINI.md): Validation aggregation.
- [internal/prompt/](file:///Users/biwakonbu/github/aglx/internal/prompt/GEMINI.md): Prompt generation and XML logic.
- [.github/](file:///Users/biwakonbu/github/aglx/.github/GEMINI.md): CI/CD and GitHub configuration.
- [docs/](file:///Users/biwakonbu/github/aglx/docs/GEMINI.md): Specifications and documentation.
- [internal/errors/](file:///Users/biwakonbu/github/aglx/internal/errors/GEMINI.md): Error management and exit codes.
- [testdata/](file:///Users/biwakonbu/github/aglx/testdata/GEMINI.md): Test patterns and validation data.
