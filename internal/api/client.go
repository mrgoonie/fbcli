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

// get performs a GET request to the Graph API
func (c *Client) get(path string, params fb.Params) (fb.Result, error) {
	if c.verbose {
		fmt.Printf("[GET] /%s %v\n", path, params)
	}

	res, err := c.session.Get(path, params)
	if err != nil {
		return nil, wrapFBError(err)
	}

	if c.verbose {
		fmt.Printf("[RESPONSE] %v\n", res)
	}

	return res, nil
}

// post performs a POST request to the Graph API
func (c *Client) post(path string, params fb.Params) (fb.Result, error) {
	if c.verbose {
		fmt.Printf("[POST] /%s %v\n", path, params)
	}

	res, err := c.session.Post(path, params)
	if err != nil {
		return nil, wrapFBError(err)
	}

	if c.verbose {
		fmt.Printf("[RESPONSE] %v\n", res)
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
