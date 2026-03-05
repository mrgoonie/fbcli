# Code Standards

## Language
- Go 1.22+
- Module path: `github.com/mrgoonie/fbcli`

## Naming
- Go standard: snake_case for files, CamelCase for exports
- Package names: lowercase, single word (`api`, `cli`, `auth`, `config`, `store`)

## Structure
- `cmd/` for entry points
- `internal/` for private packages (not importable by others)
- Test files alongside source: `*_test.go`

## Dependencies
- Cobra for CLI
- Viper for config (env var binding)
- huandu/facebook/v2 for Graph API
- testify for assertions
- No unnecessary external deps (use stdlib when possible)

## Error Handling
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- APIError type with Facebook error codes + user hints
- Never expose tokens in error messages

## Security
- Token files: `0600` permissions
- Config dir: `0700` permissions
- No secrets in logs (unless --verbose, and even then masked)
- Env var override for CI/CD (FBCLI_PAGE_TOKEN)

## Testing
- Use `t.TempDir()` for file tests
- No network calls in unit tests
- httptest for mocked API responses
- Target 70%+ coverage

## Commits
- Conventional commits: `feat:`, `fix:`, `docs:`, `test:`, `chore:`
- No AI references in messages
