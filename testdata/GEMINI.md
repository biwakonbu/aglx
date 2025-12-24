# testdata/ GEMINI

This directory contains datasets used for automated validation testing.

## Structure
- `valid/`: Sample skill packages that MUST pass validation.
- `invalid/`: Sample skill packages designed to trigger specific validation errors or warnings.
- `claude/`: Specific variations for Claude Skills testing.

## Maintenance Rules
- Never delete test cases without ensuring the coverage is moved elsewhere.
- When adding new validation rules, always add a corresponding `invalid` test case.
- Ensure all test directories follow the naming conventions described in `GEMINI.md`.
