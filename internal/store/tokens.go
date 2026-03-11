package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mrgoonie/fbcli/internal/config"
	"gopkg.in/yaml.v3"
)

// TokenStore holds authentication tokens
type TokenStore struct {
	PageToken string    `yaml:"page_token"`
	UserToken string    `yaml:"user_token,omitempty"`
	PageID    string    `yaml:"page_id"`
	PageName  string    `yaml:"page_name"`
	ExpiresAt time.Time `yaml:"expires_at,omitempty"`
	UpdatedAt time.Time `yaml:"updated_at"`
}

// TokensFilePath returns the tokens file path
func TokensFilePath() string {
	return filepath.Join(config.Dir(), "tokens.yaml")
}

// Load reads tokens from disk
func Load() (*TokenStore, error) {
	data, err := os.ReadFile(TokensFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading tokens: %w", err)
	}

	var ts TokenStore
	if err := yaml.Unmarshal(data, &ts); err != nil {
		return nil, fmt.Errorf("parsing tokens: %w", err)
	}

	return &ts, nil
}

// Save writes tokens to disk with restricted permissions
func Save(ts *TokenStore) error {
	if err := config.EnsureDir(); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	ts.UpdatedAt = time.Now()

	data, err := yaml.Marshal(ts)
	if err != nil {
		return fmt.Errorf("marshaling tokens: %w", err)
	}

	if err := os.WriteFile(TokensFilePath(), data, 0600); err != nil {
		return fmt.Errorf("writing tokens: %w", err)
	}

	return nil
}

// Clear removes the tokens file
func Clear() error {
	err := os.Remove(TokensFilePath())
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing tokens: %w", err)
	}
	return nil
}

// IsValid checks if the token store has a usable page token
func (ts *TokenStore) IsValid() bool {
	if ts == nil || ts.PageToken == "" {
		return false
	}
	if !ts.ExpiresAt.IsZero() && time.Now().After(ts.ExpiresAt) {
		return false
	}
	return true
}

// ResolveToken returns page token and page ID from env var or stored tokens.
// Priority: FBCLI_PAGE_TOKEN env > stored tokens
func ResolveToken() (token, pageID string, err error) {
	// Check env var first (CI/CD)
	if t := os.Getenv("FBCLI_PAGE_TOKEN"); t != "" {
		pid := os.Getenv("FBCLI_PAGE_ID")
		if pid == "" {
			pid = "me"
		}
		return t, pid, nil
	}

	// Load stored tokens
	ts, err := Load()
	if err != nil {
		return "", "", err
	}

	if ts == nil || !ts.IsValid() {
		return "", "", fmt.Errorf("not authenticated. Run: fbcli auth login")
	}

	return ts.PageToken, ts.PageID, nil
}
