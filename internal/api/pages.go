package api

import (
	"fmt"

	fb "github.com/huandu/facebook/v2"
)

// PageInfo holds basic page information
type PageInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category,omitempty"`
	FanCount int    `json:"fan_count,omitempty"`
	Link     string `json:"link,omitempty"`
}

// ValidateToken checks if the page token is valid
func (c *Client) ValidateToken() (*PageInfo, error) {
	res, err := c.get("me", fb.Params{
		"fields": "id,name",
	})
	if err != nil {
		return nil, err
	}

	return &PageInfo{
		ID:   res.Get("id").(string),
		Name: res.Get("name").(string),
	}, nil
}

// GetPageInfo returns detailed page information
func (c *Client) GetPageInfo() (*PageInfo, error) {
	res, err := c.get(c.pageID, fb.Params{
		"fields": "id,name,category,fan_count,link",
	})
	if err != nil {
		return nil, err
	}

	info := &PageInfo{
		ID: res.Get("id").(string),
	}

	if name, ok := res["name"].(string); ok {
		info.Name = name
	}
	if cat, ok := res["category"].(string); ok {
		info.Category = cat
	}
	if count, ok := res["fan_count"].(float64); ok {
		info.FanCount = int(count)
	}
	if link, ok := res["link"].(string); ok {
		info.Link = link
	}

	return info, nil
}

// FetchUserPages fetches pages the user manages (using user token, not page token)
func FetchUserPages(userToken string, verbose bool) ([]PageInfo, error) {
	app := fb.New("", "")
	session := app.Session(userToken)
	session.Version = "v24.0"

	if verbose {
		fmt.Println("[GET] /me/accounts")
	}

	res, err := session.Get("/me/accounts", fb.Params{
		"fields": "id,name,access_token,category",
	})
	if err != nil {
		return nil, wrapFBError(err)
	}

	var response struct {
		Data []struct {
			ID          string `facebook:"id"`
			Name        string `facebook:"name"`
			AccessToken string `facebook:"access_token"`
			Category    string `facebook:"category"`
		} `facebook:"data"`
	}

	if err := res.Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding pages: %w", err)
	}

	pages := make([]PageInfo, 0, len(response.Data))
	for _, p := range response.Data {
		pages = append(pages, PageInfo{
			ID:       p.ID,
			Name:     p.Name,
			Category: p.Category,
		})
	}

	return pages, nil
}

// PageTokenInfo holds page info with its access token (used during auth)
type PageTokenInfo struct {
	PageInfo
	AccessToken string `json:"-"`
}

// FetchUserPageTokens fetches pages with their access tokens
func FetchUserPageTokens(userToken string, verbose bool) ([]PageTokenInfo, error) {
	app := fb.New("", "")
	session := app.Session(userToken)
	session.Version = "v24.0"

	if verbose {
		fmt.Println("[GET] /me/accounts")
	}

	res, err := session.Get("/me/accounts", fb.Params{
		"fields": "id,name,access_token,category",
	})
	if err != nil {
		return nil, wrapFBError(err)
	}

	var response struct {
		Data []struct {
			ID          string `facebook:"id"`
			Name        string `facebook:"name"`
			AccessToken string `facebook:"access_token"`
			Category    string `facebook:"category"`
		} `facebook:"data"`
	}

	if err := res.Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding pages: %w", err)
	}

	pages := make([]PageTokenInfo, 0, len(response.Data))
	for _, p := range response.Data {
		pages = append(pages, PageTokenInfo{
			PageInfo: PageInfo{
				ID:       p.ID,
				Name:     p.Name,
				Category: p.Category,
			},
			AccessToken: p.AccessToken,
		})
	}

	return pages, nil
}
