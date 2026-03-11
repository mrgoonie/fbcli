# Code Standards

## Language
- Go 1.24+
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
- Use `getString()` helper for safe API response parsing (guards against type assertion panics)

## Security
- **OAuth:** CSRF state parameter required (`GenerateState()`)
- **Token exchange:** POST method (secrets in body, not URL query)
- **HTTP reads:** All responses limited to 1MB with `io.LimitReader`
- **Token files:** `0600` permissions
- **Config dir:** `0700` permissions
- **Verbose logging:** Redacts `access_token`, `client_secret` params, shows `OK` instead of raw responses
- **Env var override:** FBCLI_PAGE_TOKEN for CI/CD

## Testing
- Use `t.TempDir()` for file tests
- No network calls in unit tests
- httptest for mocked API responses
- Target 70%+ coverage

## Commits
- Conventional commits: `feat:`, `fix:`, `docs:`, `test:`, `chore:`
- No AI references in messages
