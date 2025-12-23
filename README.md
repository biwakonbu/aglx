# aglx

**Agent sKiLls eXaminer** - Validation tool for Agent Skills (SKILL.md) and Claude Skills (CLAUDE.md).

`aglx` is a CLI tool written in Go that validates `SKILL.md` files against the [Agent Skills](https://agentskills.io/) specification and checks `CLAUDE.md` for common issues.

## Installation

```bash
go install github.com/biwakonbu/aglx/cmd/aglx@latest
```

Alternatively, build from source:

```bash
git clone https://github.com/biwakonbu/aglx.git
cd aglx
go build -o aglx ./cmd/aglx
```

## Usage

### Skill Validation

```bash
# Validate a single skill
aglx validate ./my-skill

# Validate multiple skills
aglx validate ./skill1 ./skill2

# Output in JSON format
aglx validate ./skills/* --json

# Quiet mode (only display errors and warnings)
aglx validate ./skills/* --quiet
```

### Prompt Generation

```bash
# Generate XML prompt from skill metadata
aglx to-prompt ./skills/*
```

### Output Example

**Success:**
```
=== ./testdata/valid/pdf-processing ===

--- Agent Skills (SKILL.md) ---
  ✓ pdf-processing

--- Claude Skills (CLAUDE.md) ---
  - CLAUDE.md not found

=== Summary ===
./testdata/valid/pdf-processing: Agent Skills: ✓ PASS | Claude Skills: - N/A
```

**Error:**
```
✗ empty-scripts
  - scripts: must not be empty if present
```

### Commands

```
aglx validate <path>...    Validate SKILL.md and CLAUDE.md simultaneously
aglx to-prompt <path>...   Generate XML prompt for AI agents
aglx version               Show version information
aglx help                  Show this help message
```

## Validation Items

| Field            | Check Description                                              |
| ----------------- | ------------------------------------------------------------ |
| `name`            | Required, 1-64 characters, lowercase alphanumeric + hyphen only |
| `name`            | No leading/trailing hyphens, no consecutive hyphens          |
| `name`            | Must match parent directory name                             |
| `description`     | Required, 1-1024 characters                                  |
| `compatibility`   | Optional, 1-500 characters                                    |
| `allowed-tools`   | Optional, format check (alphanumeric or `Tool(args)`)       |
| `scripts/` etc.   | Optional, must be a directory and not empty                  |
| File Existence    | Verifies `SKILL.md` exists                                    |

## Specification

For detailed Agent Skills specification, see [docs/agent-skills-specification.md](docs/agent-skills-specification.md).

## Development

### Run Tests

```bash
go test ./... -v
```

### Build

```bash
goreleaser build --snapshot --clean
```

## License

MIT License
