# cmd/ GEMINI

This directory contains the CLI implementation for `aglx`.

## CLI Design
- **Subcommands**: Use `cobra` or standard `flag` for subcommands. Current commands: `validate`, `to-prompt`.
- **Output**: Support both human-readable text and JSON output using the `--json` flag.

## Dependencies
- Depends on all internal packages in `internal/`.
- Should not contain core business logic; only CLI plumbing.

## Future Commands
- `init`: Scaffold a new skill directory.
- `bundle`: Package assets for distribution.
