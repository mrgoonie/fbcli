# System Architecture

## Overview
```
User → Cobra CLI → API Client → Facebook Graph API v24.0
                 → Auth Module → OAuth 2.0 → Token Store (~/.fbcli/)
```

## Package Structure
```
cmd/fbcli/main.go        # Entry point
internal/
├── cli/                  # Cobra commands (root, auth, post, list, read, delete)
│   ├── root.go          # Root command, global flags
│   ├── auth.go          # auth login|status|logout
│   ├── post.go          # post command with media flags
│   ├── list.go          # list with table output
│   ├── read.go          # read post details
│   ├── delete.go        # delete with confirmation
│   └── output.go        # Shared formatting helpers
├── api/                  # Facebook Graph API client
│   ├── client.go        # HTTP session wrapper
│   ├── posts.go         # CRUD operations
│   ├── media.go         # Photo + chunked video upload
│   ├── pages.go         # Page info + token validation
│   └── errors.go        # Error types with hints
├── auth/                 # OAuth 2.0 flow
│   ├── oauth.go         # URL generation, token exchange
│   └── callback_listener.go  # Local HTTP callback server
├── config/               # Configuration management
│   └── config.go        # YAML config at ~/.fbcli/config.yaml
└── store/                # Token persistence
    └── tokens.go        # YAML tokens at ~/.fbcli/tokens.yaml
```

## Authentication Flow
1. `fbcli auth login` → prompts for App ID/Secret (first time)
2. Generate CSRF state parameter for OAuth security
3. Opens browser → Facebook OAuth dialog
4. Local HTTP server receives callback on 127.0.0.1:8910 (secure localhost binding)
5. Verify CSRF state parameter, then exchange code via POST (secrets in body, not URL)
6. GET /me/accounts → page access tokens
7. Store selected page token in ~/.fbcli/tokens.yaml

## Token Resolution Priority
1. `FBCLI_PAGE_TOKEN` env var (CI/CD)
2. `~/.fbcli/tokens.yaml` file (interactive)

## Key Design Decisions
- Page Access Tokens never expire (until revoked)
- Video upload uses 3-phase chunked protocol (4MB chunks)
- Scheduled posts via `published=false` + `scheduled_publish_time`
- No external tablewriter dep - uses stdlib tabwriter
