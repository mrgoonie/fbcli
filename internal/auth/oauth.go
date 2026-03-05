package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// OAuthConfig holds Facebook OAuth configuration
type OAuthConfig struct {
	AppID       string
	AppSecret   string
	RedirectURI string
	Scopes      []string
}

// DefaultScopes returns the required Facebook permissions
func DefaultScopes() []string {
	return []string{
		"pages_manage_posts",
		"pages_read_engagement",
		"pages_show_list",
		"public_profile",
	}
}

// AuthURL builds the Facebook OAuth authorization URL
func AuthURL(cfg OAuthConfig) string {
	scopes := ""
	for i, s := range cfg.Scopes {
		if i > 0 {
			scopes += ","
		}
		scopes += s
	}

	params := url.Values{
		"client_id":    {cfg.AppID},
		"redirect_uri": {cfg.RedirectURI},
		"scope":        {scopes},
		"response_type": {"code"},
	}

	return "https://www.facebook.com/v24.0/dialog/oauth?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for a user access token
func ExchangeCode(ctx context.Context, cfg OAuthConfig, code string) (string, error) {
	params := url.Values{
		"client_id":     {cfg.AppID},
		"redirect_uri":  {cfg.RedirectURI},
		"client_secret": {cfg.AppSecret},
		"code":          {code},
	}

	reqURL := "https://graph.facebook.com/v24.0/oauth/access_token?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("exchanging code: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		AccessToken string `json:"access_token"`
		Error       struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parsing token response: %w", err)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("token exchange failed: %s", result.Error.Message)
	}

	return result.AccessToken, nil
}

// ExtendToken exchanges a short-lived token for a long-lived one (~60 days)
func ExtendToken(ctx context.Context, cfg OAuthConfig, shortToken string) (string, error) {
	params := url.Values{
		"grant_type":        {"fb_exchange_token"},
		"client_id":         {cfg.AppID},
		"client_secret":     {cfg.AppSecret},
		"fb_exchange_token": {shortToken},
	}

	reqURL := "https://graph.facebook.com/v24.0/oauth/access_token?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("extending token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		AccessToken string `json:"access_token"`
		Error       struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parsing extended token: %w", err)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("token extension failed: %s", result.Error.Message)
	}

	return result.AccessToken, nil
}
