package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mrgoonie/fbcli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestTokensFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	path := TokensFilePath()
	expected := filepath.Join(tmpDir, ".fbcli", "tokens.yaml")
	assert.Equal(t, expected, path)
}

func TestLoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	ts, err := Load()
	assert.NoError(t, err)
	assert.Nil(t, ts)
}

func TestLoadExisting(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	// Create config directory
	require.NoError(t, os.MkdirAll(config.Dir(), 0700))

	// Create token store
	now := time.Now()
	testStore := &TokenStore{
		PageToken: "test-page-token",
		UserToken: "test-user-token",
		PageID:    "12345",
		PageName:  "My Page",
		UpdatedAt: now,
		ExpiresAt: now.Add(24 * time.Hour),
	}

	data, err := yaml.Marshal(testStore)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(TokensFilePath(), data, 0600))

	// Load and verify
	ts, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, ts)
	assert.Equal(t, "test-page-token", ts.PageToken)
	assert.Equal(t, "test-user-token", ts.UserToken)
	assert.Equal(t, "12345", ts.PageID)
	assert.Equal(t, "My Page", ts.PageName)
}

func TestLoadMalformed(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	require.NoError(t, os.MkdirAll(config.Dir(), 0700))
	// Write invalid YAML
	require.NoError(t, os.WriteFile(TokensFilePath(), []byte("invalid: yaml: content:"), 0600))

	ts, err := Load()
	assert.Error(t, err)
	assert.Nil(t, ts)
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	store := &TokenStore{
		PageToken: "test-token",
		PageID:    "54321",
		PageName:  "Test Page",
	}

	err := Save(store)
	assert.NoError(t, err)

	// Verify file exists with correct permissions
	info, err := os.Stat(TokensFilePath())
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode())

	// Verify content
	data, err := os.ReadFile(TokensFilePath())
	require.NoError(t, err)

	loaded := &TokenStore{}
	err = yaml.Unmarshal(data, loaded)
	require.NoError(t, err)

	assert.Equal(t, "test-token", loaded.PageToken)
	assert.Equal(t, "54321", loaded.PageID)
	assert.Equal(t, "Test Page", loaded.PageName)
	// UpdatedAt should be set to recent time
	assert.NotZero(t, loaded.UpdatedAt)
	assert.True(t, time.Since(loaded.UpdatedAt) < time.Second)
}

func TestSaveCreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	store := &TokenStore{
		PageToken: "test-token",
		PageID:    "54321",
	}

	// Ensure directory doesn't exist
	assert.NoFileExists(t, config.Dir())

	err := Save(store)
	assert.NoError(t, err)

	// Verify directory was created
	assert.DirExists(t, config.Dir())
	assert.FileExists(t, TokensFilePath())
}

func TestClear(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	// Create and save token store
	require.NoError(t, os.MkdirAll(config.Dir(), 0700))
	store := &TokenStore{PageToken: "test"}
	require.NoError(t, Save(store))

	// Verify file exists
	assert.FileExists(t, TokensFilePath())

	// Clear tokens
	err := Clear()
	assert.NoError(t, err)

	// Verify file is deleted
	assert.NoFileExists(t, TokensFilePath())
}

func TestClearNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	// Clear when file doesn't exist should not error
	err := Clear()
	assert.NoError(t, err)
}

func TestIsValidNil(t *testing.T) {
	var ts *TokenStore
	assert.False(t, ts.IsValid())
}

func TestIsValidEmptyToken(t *testing.T) {
	ts := &TokenStore{
		PageToken: "",
		PageID:    "12345",
	}
	assert.False(t, ts.IsValid())
}

func TestIsValidExpired(t *testing.T) {
	ts := &TokenStore{
		PageToken: "test-token",
		PageID:    "12345",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	assert.False(t, ts.IsValid())
}

func TestIsValidNotExpired(t *testing.T) {
	ts := &TokenStore{
		PageToken: "test-token",
		PageID:    "12345",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	assert.True(t, ts.IsValid())
}

func TestIsValidNoExpiry(t *testing.T) {
	ts := &TokenStore{
		PageToken: "test-token",
		PageID:    "12345",
		// ExpiresAt is zero
	}
	assert.True(t, ts.IsValid())
}

func TestResolveTokenFromEnv(t *testing.T) {
	oldToken := os.Getenv("FBCLI_PAGE_TOKEN")
	oldPageID := os.Getenv("FBCLI_PAGE_ID")

	t.Cleanup(func() {
		os.Setenv("FBCLI_PAGE_TOKEN", oldToken)
		os.Setenv("FBCLI_PAGE_ID", oldPageID)
	})

	os.Setenv("FBCLI_PAGE_TOKEN", "env-token")
	os.Setenv("FBCLI_PAGE_ID", "env-page-123")

	token, pageID, err := ResolveToken()
	assert.NoError(t, err)
	assert.Equal(t, "env-token", token)
	assert.Equal(t, "env-page-123", pageID)
}

func TestResolveTokenFromEnvDefaultPageID(t *testing.T) {
	oldToken := os.Getenv("FBCLI_PAGE_TOKEN")
	oldPageID := os.Getenv("FBCLI_PAGE_ID")

	t.Cleanup(func() {
		os.Setenv("FBCLI_PAGE_TOKEN", oldToken)
		os.Setenv("FBCLI_PAGE_ID", oldPageID)
	})

	os.Setenv("FBCLI_PAGE_TOKEN", "env-token")
	os.Setenv("FBCLI_PAGE_ID", "")

	token, pageID, err := ResolveToken()
	assert.NoError(t, err)
	assert.Equal(t, "env-token", token)
	assert.Equal(t, "me", pageID)
}

func TestResolveTokenFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	oldToken := os.Getenv("FBCLI_PAGE_TOKEN")

	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("FBCLI_PAGE_TOKEN", oldToken)
	})

	os.Setenv("HOME", tmpDir)
	os.Setenv("FBCLI_PAGE_TOKEN", "")

	// Create and save token store
	require.NoError(t, os.MkdirAll(config.Dir(), 0700))
	store := &TokenStore{
		PageToken: "file-token",
		PageID:    "file-page-id",
	}
	require.NoError(t, Save(store))

	token, pageID, err := ResolveToken()
	assert.NoError(t, err)
	assert.Equal(t, "file-token", token)
	assert.Equal(t, "file-page-id", pageID)
}

func TestResolveTokenNotAuthenticated(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	oldToken := os.Getenv("FBCLI_PAGE_TOKEN")

	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("FBCLI_PAGE_TOKEN", oldToken)
	})

	os.Setenv("HOME", tmpDir)
	os.Setenv("FBCLI_PAGE_TOKEN", "")

	// No token file exists
	token, pageID, err := ResolveToken()
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Empty(t, pageID)
	assert.Contains(t, err.Error(), "not authenticated")
}

func TestResolveTokenExpired(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	oldToken := os.Getenv("FBCLI_PAGE_TOKEN")

	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("FBCLI_PAGE_TOKEN", oldToken)
	})

	os.Setenv("HOME", tmpDir)
	os.Setenv("FBCLI_PAGE_TOKEN", "")

	// Create expired token
	require.NoError(t, os.MkdirAll(config.Dir(), 0700))
	store := &TokenStore{
		PageToken: "expired-token",
		PageID:    "page-id",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	require.NoError(t, Save(store))

	token, pageID, err := ResolveToken()
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Empty(t, pageID)
	assert.Contains(t, err.Error(), "not authenticated")
}

func TestResolveTokenEnvOverridesFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	oldToken := os.Getenv("FBCLI_PAGE_TOKEN")
	oldPageID := os.Getenv("FBCLI_PAGE_ID")

	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("FBCLI_PAGE_TOKEN", oldToken)
		os.Setenv("FBCLI_PAGE_ID", oldPageID)
	})

	os.Setenv("HOME", tmpDir)

	// Create file-based token
	require.NoError(t, os.MkdirAll(config.Dir(), 0700))
	store := &TokenStore{
		PageToken: "file-token",
		PageID:    "file-page-id",
	}
	require.NoError(t, Save(store))

	// Set env token with page ID
	os.Setenv("FBCLI_PAGE_TOKEN", "env-token")
	os.Setenv("FBCLI_PAGE_ID", "env-page-id")

	token, pageID, err := ResolveToken()
	assert.NoError(t, err)
	// Env should take priority
	assert.Equal(t, "env-token", token)
	assert.Equal(t, "env-page-id", pageID)
}
