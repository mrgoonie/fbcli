# fbcli

A command-line tool for managing Facebook Pages. Create posts, upload media, schedule content — all from your terminal.

Inspired by [xurl](https://github.com/xdevplatform/xurl).

## Install

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

1. Go to [Meta Developers](https://developers.facebook.com/) and create a new app (type: **Business**)
2. Go to **Use Cases** → **Add use case** → select **"Manage Everything on your Page"**
3. Click **Customize** → find `pages_manage_posts` in the "Permissions & features" column → click **Add**

> No redirect URI setup needed — development mode works with localhost automatically.

### 2. Get a Page Token

1. Open [Graph API Explorer](https://developers.facebook.com/tools/explorer/)
2. Select your app from the **Meta App** dropdown
3. Under **Permissions**, add: `pages_manage_posts`, `pages_read_engagement`, `pages_show_list`
4. Change **User or Page** to your Facebook Page (not "User Token")
5. Click **Generate Access Token** and approve the permissions
6. Copy the generated token

### 3. Authenticate

```bash
fbcli auth token "<your-page-token>"
```

### 4. Verify

```bash
# Check auth status
fbcli auth status

# List recent posts
fbcli list

# Try posting
fbcli post "Hello from the terminal!" --verbose
```

> **Alternative:** You can also use `fbcli auth login` for browser-based OAuth flow.
> This requires configuring a redirect URI (`http://localhost:8910/callback`) in your app's Facebook Login settings.

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
| `fbcli auth token` | Set a Page Access Token manually |
| `fbcli auth login` | Authenticate via browser OAuth flow |
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
| `-V, --verbose` | Show API request/response details (sensitive data redacted) |
| `--page` | Override default page ID |
| `-v, --version` | Show version |

## Token Storage

Tokens are stored at `~/.fbcli/` with `0600` permissions:
- `config.yaml` - App ID, App Secret, default page
- `tokens.yaml` - Page Access Token (never expires until revoked)

## License

MIT
