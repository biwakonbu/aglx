# internal/skill GEMINI

This package is responsible for strictly validating the **Agent Skills** specification.

## Responsibilities
- Parse `SKILL.md` YAML frontmatter.
- Verify directory structure (e.g., `scripts/`, `assets/` existence).
- Check `SKILL.md` body size for token efficiency.

## Key Files
- `validator.go`: Core validation logic.
- `types.go`: Frontmatter struct definitions.

## Performance
- Validation should be fast and non-destructive.
- Body token count estimation should be consistent.
