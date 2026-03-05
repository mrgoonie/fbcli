# fbcli

A command-line tool for managing Facebook Pages. Create posts, upload media, schedule content — all from your terminal.

Inspired by [xurl](https://github.com/xdevplatform/xurl).

## Install

### Homebrew (macOS/Linux)

```bash
brew install mrgoonie/tap/fbcli
```

### Go

```bash
go install github.com/mrgoonie/fbcli/cmd/fbcli@latest
```

### From source

```bash
git clone https://github.com/mrgoonie/fbcli.git
cd fbcli
make build
```

## Setup

### 1. Create a Facebook App

1. Go to [Meta Developers](https://developers.facebook.com/)
2. Create a new app (type: **Business**)
3. Note your **App ID** and **App Secret**
4. Add **Facebook Login** product
5. Set redirect URI to `http://localhost:8910/callback`
6. Required permissions: `pages_manage_posts`, `pages_read_engagement`, `pages_show_list`

### 2. Authenticate

```bash
fbcli auth login
```

This will:
- Prompt for App ID and App Secret (first time only)
- Open your browser for Facebook OAuth
- Store your Page Access Token locally at `~/.fbcli/`

### 3. Verify

```bash
fbcli auth status
```

## Usage

### Post content

```bash
# Text post
fbcli post "Hello from the terminal!"

# Photo post
fbcli post -i photo.jpg "Check this out"

# Video post
fbcli post -v video.mp4 "Watch this"

# Link post
fbcli post -l https://example.com "Read this article"

# Scheduled post
fbcli post --schedule "2026-03-10T10:00" "Coming soon"

# Pipe from stdin
echo "Hello world" | fbcli post
```

### Manage posts

```bash
# List recent posts
fbcli list
fbcli list -n 20

# View post details
fbcli read POST_ID

# Delete a post
fbcli delete POST_ID
fbcli delete POST_ID --force
```

### Output formats

```bash
# JSON output (for scripting)
fbcli list --json
fbcli read POST_ID --json

# Verbose mode (debug)
fbcli post "Hello" --verbose
```

## CI/CD

Use environment variables to skip interactive login:

```bash
export FBCLI_PAGE_TOKEN="your-page-access-token"
export FBCLI_PAGE_ID="your-page-id"

fbcli post "Automated post from CI"
```

You can also set app credentials via environment:

```bash
export FBCLI_APP_ID="your-app-id"
export FBCLI_APP_SECRET="your-app-secret"
```

## Commands

| Command | Description |
|---------|-------------|
| `fbcli auth login` | Authenticate with Facebook |
| `fbcli auth status` | Show authentication status |
| `fbcli auth logout` | Clear stored tokens |
| `fbcli post` | Create a new post |
| `fbcli list` | List recent posts |
| `fbcli read` | View post details |
| `fbcli delete` | Delete a post |

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |
| `-V, --verbose` | Show API request/response details |
| `--page` | Override default page ID |
| `-v, --version` | Show version |

## Token Storage

Tokens are stored at `~/.fbcli/` with `0600` permissions:
- `config.yaml` - App ID, App Secret, default page
- `tokens.yaml` - Page Access Token (never expires until revoked)

## License

MIT
