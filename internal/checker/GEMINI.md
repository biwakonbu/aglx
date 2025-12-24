# internal/checker GEMINI

This package aggregates results from both the `skill` and `claude` validators.

## Responsibilities
- Provide a unified `Result` struct.
- Handle multi-directory validation passes.
- Summarize errors and warnings for the CLI layer.

## Future Plans
- Add parallel validation support for large numbers of skill packages.
