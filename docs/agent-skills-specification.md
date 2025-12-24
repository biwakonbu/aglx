# Agent Skills Specification

This document summarizes the official specification for [Agent Skills](https://agentskills.io/).

## Overview

Agent Skills is a simple, open format for giving AI agents new capabilities and expertise.

### What Agent Skills Can Do

- **Domain Expertise**: Package specialized knowledge—from legal review processes to data analysis pipelines—into reusable instructions.
- **New Capabilities**: Grant agents new abilities (e.g., creating presentations, building MCP servers, analyzing datasets).
- **Reproducible Workflows**: Convert multi-step tasks into consistent, auditable workflows.
- **Interoperability**: Reuse the same skills across different skill-compatible agent products.

---

## Directory Structure

In its simplest form, only a `SKILL.md` file is required:

```
skill-name/
└── SKILL.md              # Required: Instructions + Metadata
```

Full configuration:

```
my-skill/
├── SKILL.md              # Required: Instructions + Metadata
├── scripts/              # Optional: Executable code
├── references/           # Optional: Documentation
└── assets/               # Optional: Templates, resources
```

---

## SKILL.md Format

The `SKILL.md` file consists of YAML frontmatter and Markdown body content.

### Frontmatter (Required)

```yaml
---
name: skill-name
description: A description of what this skill does and when to use it.
---
```

#### Example with All Fields

```yaml
---
name: pdf-processing
description: Extract text and tables from PDF files, fill forms, merge documents.
license: Apache-2.0
compatibility: Designed for Claude Code (or similar products)
allowed-tools: Bash(git:*) Bash(jq:*) Read
metadata:
  author: example-org
  version: "1.0"
---
```

---

## Frontmatter Field Details

### `name` Field (Required)

A short identifier for the skill.

**Constraints:**
- 1–64 characters
- Only Unicode lowercase alphanumeric characters and hyphens allowed (`a-z`, `0-9`, `-`)
- Must not start or end with a hyphen
- Must not contain consecutive hyphens (`--`)
- **Must match the parent directory name**

**Valid Examples:**
```yaml
name: pdf-processing
name: data-analysis
name: code-review
```

**Invalid Examples:**
```yaml
name: PDF-Processing      # Uppercase not allowed
name: -pdf                 # Cannot start with a hyphen
name: pdf--processing      # Consecutive hyphens not allowed
```

---

### `description` Field (Required)

Explains what the skill does and when to use it.

**Constraints:**
- 1–1024 characters
- Should describe both the skill's functionality and usage timing
- Should include specific keywords to help agents identify relevant tasks

**Good Example:**
```yaml
description: Extracts text and tables from PDF files, fills PDF forms, and merges multiple PDFs. Use when working with PDF documents or when the user mentions PDFs, forms, or document extraction.
```

**Bad Example:**
```yaml
description: Helps with PDFs.  # Insufficient detail
```

---

### `license` Field (Optional)

Specifies the license applied to the skill.

```yaml
license: Apache-2.0
license: Proprietary. LICENSE.txt has complete terms
```

---

### `compatibility` Field (Optional)

Used when the skill has specific environmental requirements.

**Constraints:**
- 1–500 characters (if provided)

**Example:**
```yaml
compatibility: Designed for Claude Code (or similar products)
compatibility: Requires git, docker, jq, and access to the internet
```

---

### `metadata` Field (Optional)

A map from string keys to string values. Used to store additional properties not defined in the Agent Skills specification.

```yaml
metadata:
  author: example-org
  version: "1.0"
```

---

### `allowed-tools` Field (Optional, Experimental)

A list of pre-approved tools. The format varies by implementation:

> ⚠️ **Experimental Feature**: Support for this field may vary depending on the agent implementation.

**Agent Skills Specification (agentskills.io):**

Space-separated format:

```yaml
allowed-tools: Bash(git:*) Bash(jq:*) Read
```

**Claude Code (code.claude.com):**

Comma-separated format:

```yaml
allowed-tools: Read, Grep, Glob
```

> **Note**: `aglx` supports both formats by default. Use `-spec=agent-skills` or `--spec=claude-code` to enforce strict format validation.

---

## Body Content

The Markdown body after the frontmatter should include:

- Step-by-step instructions
- Input and output examples
- Common edge cases

### Complete SKILL.md Example

```markdown
---
name: pdf-processing
description: Extract text and tables from PDF files, fill forms, merge documents.
---

# PDF Processing

## When to use this skill

Use this skill when the user needs to work with PDF files...

## How to extract text

1. Use pdfplumber for text extraction...

## How to fill forms

...
```

---

## Optional Directories

### `scripts/`

Stores executable scripts.

**Best Practices:**
- Ensure they are self-contained or clearly document dependencies
- Include helpful error messages
- Properly handle edge cases

### `references/`

Stores supplementary documentation.

**Common Files:**
- `REFERENCE.md` - Detailed technical reference
- `FORMS.md` - Form templates or structured data formats
- Domain-specific files (`finance.md`, `legal.md`, etc.)

### `assets/`

Stores static resources.

**Examples:**
- Templates (document templates, configuration templates)
- Images (diagrams, examples)
- Data files (lookup tables, schemas)

---

## Progressive Disclosure

Skills are loaded progressively for token efficiency:

1. **Metadata (approx. 100 tokens)**: Load the `name` and `description` fields of all skills at startup.
2. **Instructions (recommended max 5000 tokens)**: Load the full `SKILL.md` body when a skill is activated.
3. **Resources (on-demand)**: Load files in `scripts/`, `references/`, and `assets/` only when needed.

---

## File References

You can reference other files within `SKILL.md`:

```markdown
See [the reference guide](references/REFERENCE.md) for details.

Run the extraction script: scripts/extract.py
```

---

## Skill Execution Flow

1. **Discovery**: At startup, the agent only loads the name and description of each available skill.
2. **Activation**: When a task matches the skill's description, the agent loads the full `SKILL.md` instructions into context.
3. **Execution**: The agent follows the instructions, reading reference files or executing bundled code as needed.

---

## Agent Integration

### Integration Summary

1. Discover skills in configured directories.
2. Load metadata (name and description) at startup.
3. Match user tasks to relevant skills.
4. Activate the skill by loading full instructions.
5. Execute scripts and access resources as needed.

### Metadata Parsing (Pseudo-Code)

```
function parseMetadata(skillPath):
    content = readFile(skillPath + "/SKILL.md")
    frontmatter = extractYAMLFrontmatter(content)
    return {
        name: frontmatter.name,
        description: frontmatter.description,
        path: skillPath
    }
```

### Context Injection (XML Format)

```xml
<available_skills>
  <skill>
    <name>pdf-processing</name>
    <description>Extracts text and tables from PDF files, fills forms, merges documents.</description>
    <location>/path/to/skills/pdf-processing/SKILL.md</location>
  </skill>
  <skill>
    <name>data-analysis</name>
    <description>Analyzes datasets, generates charts, and creates summary reports.</description>
    <location>/path/to/skills/data-analysis/SKILL.md</location>
  </skill>
</available_skills>
```

---

## Security Considerations

- **Sandboxing**: Run scripts in isolated environments.
- **Allowlisting**: Only execute scripts from trusted skills.
- **Confirmation**: Ask for user confirmation before performing potentially dangerous operations.
- **Logging**: Record all script executions for auditing.

---

## Validation

You can validate skills using the official reference library [skills-ref](https://github.com/agentskills/agentskills/tree/main/skills-ref):

```bash
skills-ref validate ./my-skill
```

Generate prompt XML:

```bash
skills-ref to-prompt <path>...
```

---

## Links

- [Official Agent Skills Website](https://agentskills.io/)
- [Specification](https://agentskills.io/specification)
- [Integrating Skills](https://agentskills.io/integrate-skills)
- [Sample Skills (GitHub)](https://github.com/anthropics/skills)
- [Reference Library (skills-ref)](https://github.com/agentskills/agentskills/tree/main/skills-ref)
- [Best Practices](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices)
