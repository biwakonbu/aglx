# Release Procedures

This project uses [GoReleaser](https://goreleaser.com/) to build and distribute binaries.

## Prerequisites

- Go 1.23 or higher
- GoReleaser v2 (for local testing)

```bash
# macOS
brew install goreleaser
```

## Automated Release (Recommended)

When you push a tag to GitHub, GitHub Actions will automatically create a release.

```bash
# 1. Commit changes
git add .
git commit -m "feat: add new feature"

# 2. Create a tag
git tag v0.1.0

# 3. Push the tag -> Automated release starts
git push origin v0.1.0
```

### Generated Artifacts

| OS | Architecture | File Name |
|----|--------------|-----------|
| Linux | amd64 | `aglx_X.X.X_linux_amd64.tar.gz` |
| Linux | arm64 | `aglx_X.X.X_linux_arm64.tar.gz` |
| macOS | amd64 (Intel) | `aglx_X.X.X_darwin_amd64.tar.gz` |
| macOS | arm64 (Apple Silicon) | `aglx_X.X.X_darwin_arm64.tar.gz` |
| Windows | amd64 | `aglx_X.X.X_windows_amd64.zip` |
| Windows | arm64 | `aglx_X.X.X_windows_arm64.zip` |

## Local Testing

You can verify the configuration before releasing.

```bash
# Validate configuration file
goreleaser check

# Snapshot build (test without tags)
goreleaser build --snapshot --clean

# Verify build results
ls -la dist/

# Test binary
./dist/aglx_darwin_arm64_v8.0/aglx version
```

## Version Information

The following information is embedded during the build:

- `version`: Retrieved from the Git tag (e.g., `v0.1.0`)
- `commit`: Commit hash
- `date`: Build date and time

```bash
$ aglx version
aglx version 0.1.0 (commit: abc1234, built at: 2024-01-01T00:00:00Z)
```

## Configuration Files

| File | Description |
|----------|------|
| `.goreleaser.yaml` | GoReleaser configuration |
| `.github/workflows/release.yml` | GitHub Actions workflow |

## Distribution via Homebrew (Optional)

To distribute via a Homebrew tap, follow these steps:

1. Create a new repository named `homebrew-tap`.
2. Uncomment the `brews` section in `.goreleaser.yaml`.
3. Push tags and release as usual.

Users can then install using:

```bash
brew install biwakonbu/tap/aglx
```

## Troubleshooting

### Error in `goreleaser check`

Ensure the configuration matches the GoReleaser v2 specification. Warnings will be displayed for deprecated settings.

### GitHub Actions Release Fails

- Verify `GITHUB_TOKEN` permissions (requires `contents: write`).
- Verify the tag format (requires `v` prefix).
