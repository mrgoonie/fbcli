package api

import (
	"fmt"

	fb "github.com/huandu/facebook/v2"
)

// Client wraps the Facebook Graph API
type Client struct {
	app     *fb.App
	session *fb.Session
	pageID  string
	verbose bool
}

// NewClient creates a new Facebook API client
func NewClient(pageToken, pageID string, verbose bool) *Client {
	app := fb.New("", "")
	session := app.Session(pageToken)
	session.Version = "v24.0"

	return &Client{
		app:     app,
		session: session,
		pageID:  pageID,
		verbose: verbose,
	}
}

// PageID returns the configured page ID
func (c *Client) PageID() string {
	return c.pageID
}

// getString safely extracts a string value from an API result
func getString(res fb.Result, key string) (string, error) {
	v, ok := res.Get(key).(string)
	if !ok || v == "" {
		return "", fmt.Errorf("unexpected API response: missing or invalid '%s' field", key)
	}
	return v, nil
}

// sensitiveKeys lists parameter names that should be redacted in verbose output
var sensitiveKeys = map[string]bool{
	"access_token":  true,
	"client_secret": true,
}

// redactParams returns a copy of params with sensitive values masked
func redactParams(p fb.Params) fb.Params {
	safe := fb.Params{}
	for k, v := range p {
		if sensitiveKeys[k] {
			safe[k] = "***"
		} else {
			safe[k] = v
		}
	}
	return safe
}

// get performs a GET request to the Graph API
func (c *Client) get(path string, params fb.Params) (fb.Result, error) {
	if c.verbose {
		fmt.Printf("[GET] /%s %v\n", path, redactParams(params))
	}

	res, err := c.session.Get(path, params)
	if err != nil {
		return nil, wrapFBError(err)
	}

	if c.verbose {
		fmt.Printf("[RESPONSE] /%s OK\n", path)
	}

	return res, nil
}

// post performs a POST request to the Graph API
func (c *Client) post(path string, params fb.Params) (fb.Result, error) {
	if c.verbose {
		fmt.Printf("[POST] /%s %v\n", path, redactParams(params))
	}

	res, err := c.session.Post(path, params)
	if err != nil {
		return nil, wrapFBError(err)
	}

	if c.verbose {
		fmt.Printf("[RESPONSE] /%s OK\n", path)
	}

	return res, nil
}

// del performs a DELETE request to the Graph API
func (c *Client) del(path string) error {
	if c.verbose {
		fmt.Printf("[DELETE] /%s\n", path)
	}

	_, err := c.session.Delete(path, nil)
	if err != nil {
		return wrapFBError(err)
	}

	return nil
}
