# internal/prompt GEMINI

This package generates XML context snippets for agent discovery.

## Responsibilities
- Convert validated skills into `<skill>` XML tags.
- Ensure XML is escaped and properly nested.
- Maintain token efficiency in the generated prompt.

## XML Structure
Skills are typically wrapped in `<agent_skills>` or similar tags for high-level agent awareness.
