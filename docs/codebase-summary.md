# Codebase Summary

## Stats
- **Language:** Go
- **Files:** 29
- **Lines:** ~3,200
- **Tests:** 80 (100% pass rate)
- **Binary size:** ~13MB

## Key Packages

### `internal/cli` (7 files)
CLI layer using Cobra. Each command is a separate file. `output.go` provides shared formatting (colors, tables, JSON, relative time).

### `internal/api` (5 files)
Facebook Graph API client. `client.go` wraps huandu/facebook session. `posts.go` handles CRUD. `media.go` implements photo upload and chunked video upload. `errors.go` maps Facebook error codes to user hints.

### `internal/auth` (2 files)
OAuth 2.0 flow. `oauth.go` builds auth URLs and exchanges tokens. `callback_listener.go` runs local HTTP server for OAuth redirect.

### `internal/config` (1 file)
YAML config at `~/.fbcli/config.yaml`. Supports env var overrides (FBCLI_APP_ID, FBCLI_APP_SECRET).

### `internal/store` (1 file)
Token persistence at `~/.fbcli/tokens.yaml`. `ResolveToken()` checks env var first, then file.

## Dependencies
- spf13/cobra - CLI framework
- huandu/facebook/v2 - Facebook SDK
- fatih/color - Terminal colors
- golang.org/x/oauth2 - OAuth support
- golang.org/x/term - Terminal detection
- gopkg.in/yaml.v3 - YAML parsing
- stretchr/testify - Test assertions
