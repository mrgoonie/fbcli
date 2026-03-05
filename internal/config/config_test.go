package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestDir(t *testing.T) {
	dir := Dir()
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".fbcli"), dir)
}

func TestFilePath(t *testing.T) {
	path := FilePath()
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	expected := filepath.Join(home, ".fbcli", "config.yaml")
	assert.Equal(t, expected, path)
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, ".fbcli")

	// Mock Dir() by changing home directory temporarily
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})

	os.Setenv("HOME", tmpDir)

	err := EnsureDir()
	assert.NoError(t, err)

	info, err := os.Stat(testDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
	// Check that directory has read/write/execute permissions for owner
	// (0700 & 0700 should equal 0700)
	assert.Equal(t, os.FileMode(0700), info.Mode().Perm())
}

func TestLoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "v24.0", cfg.APIVersion)
	assert.Empty(t, cfg.AppID)
	assert.Empty(t, cfg.AppSecret)
}

func TestLoadExisting(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	// Create config directory and file
	require.NoError(t, os.MkdirAll(Dir(), 0700))

	testCfg := &Config{
		AppID:         "test-app-id",
		AppSecret:     "test-secret",
		DefaultPageID: "12345",
		APIVersion:    "v20.0",
	}

	data, err := yaml.Marshal(testCfg)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(FilePath(), data, 0600))

	// Load and verify
	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "test-app-id", cfg.AppID)
	assert.Equal(t, "test-secret", cfg.AppSecret)
	assert.Equal(t, "12345", cfg.DefaultPageID)
	assert.Equal(t, "v20.0", cfg.APIVersion)
}

func TestLoadWithEnvOverride(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	oldAppID := os.Getenv("FBCLI_APP_ID")
	oldAppSecret := os.Getenv("FBCLI_APP_SECRET")

	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("FBCLI_APP_ID", oldAppID)
		os.Setenv("FBCLI_APP_SECRET", oldAppSecret)
	})

	os.Setenv("HOME", tmpDir)

	// Create config file with different values
	require.NoError(t, os.MkdirAll(Dir(), 0700))
	testCfg := &Config{
		AppID:     "file-app-id",
		AppSecret: "file-secret",
	}
	data, err := yaml.Marshal(testCfg)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(FilePath(), data, 0600))

	// Set env vars to override
	os.Setenv("FBCLI_APP_ID", "env-app-id")
	os.Setenv("FBCLI_APP_SECRET", "env-secret")

	cfg, err := Load()
	assert.NoError(t, err)
	// Env vars should override file values
	assert.Equal(t, "env-app-id", cfg.AppID)
	assert.Equal(t, "env-secret", cfg.AppSecret)
}

func TestLoadMalformed(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	require.NoError(t, os.MkdirAll(Dir(), 0700))
	// Write invalid YAML
	require.NoError(t, os.WriteFile(FilePath(), []byte("invalid: yaml: content:"), 0600))

	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	cfg := &Config{
		AppID:         "test-id",
		AppSecret:     "test-secret",
		DefaultPageID: "98765",
		APIVersion:    "v22.0",
	}

	err := Save(cfg)
	assert.NoError(t, err)

	// Verify file exists and has correct permissions
	info, err := os.Stat(FilePath())
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode())

	// Verify content
	data, err := os.ReadFile(FilePath())
	require.NoError(t, err)

	loadedCfg := &Config{}
	err = yaml.Unmarshal(data, loadedCfg)
	require.NoError(t, err)

	assert.Equal(t, cfg.AppID, loadedCfg.AppID)
	assert.Equal(t, cfg.AppSecret, loadedCfg.AppSecret)
	assert.Equal(t, cfg.DefaultPageID, loadedCfg.DefaultPageID)
	assert.Equal(t, cfg.APIVersion, loadedCfg.APIVersion)
}

func TestSaveCreatesDirIfMissing(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	os.Setenv("HOME", tmpDir)

	cfg := &Config{
		AppID:     "test-id",
		AppSecret: "test-secret",
	}

	// Ensure directory doesn't exist
	assert.NoFileExists(t, Dir())

	err := Save(cfg)
	assert.NoError(t, err)

	// Verify directory was created
	assert.DirExists(t, Dir())
	assert.FileExists(t, FilePath())
}
