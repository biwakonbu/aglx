# internal/claude GEMINI

This package handles validation for **Claude Skills** (`CLAUDE.md`).

## Responsibilities
- Ensure `CLAUDE.md` exists and is properly formatted.
- Warn if file sizes in the skill package exceed Claude's context limits.
- Validate the internal structure against the Claude Skills spec.

## Implementation Notes
- Focus on "Warnings" for non-breaking but inefficient patterns.
- Keep standard Claude patterns in mind (e.g., project knowledge).
