package api

import (
	"fmt"
	"time"

	fb "github.com/huandu/facebook/v2"
)

// PostOptions configures a new post
type PostOptions struct {
	ImagePath    string
	VideoPath    string
	ScheduleTime *time.Time
	Link         string
}

// PostResult represents a created post
type PostResult struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Scheduled bool   `json:"scheduled"`
}

// PostSummary represents a post in list view
type PostSummary struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	CreatedTime time.Time `json:"created_time"`
	FullPicture string    `json:"full_picture,omitempty"`
}

// PostDetail represents a single post with engagement data
type PostDetail struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	CreatedTime time.Time `json:"created_time"`
	FullPicture string    `json:"full_picture,omitempty"`
	Shares      int       `json:"shares"`
	Reactions   int       `json:"reactions"`
	Comments    int       `json:"comments"`
	URL         string    `json:"url"`
}

// CreatePost publishes a new post to the page
func (c *Client) CreatePost(message string, opts PostOptions) (*PostResult, error) {
	// Handle media uploads
	if opts.ImagePath != "" {
		return c.createPhotoPost(message, opts.ImagePath, opts.ScheduleTime)
	}
	if opts.VideoPath != "" {
		return c.createVideoPost(message, opts.VideoPath, opts.ScheduleTime)
	}

	params := fb.Params{"message": message}

	if opts.Link != "" {
		params["link"] = opts.Link
	}

	if opts.ScheduleTime != nil {
		params["published"] = false
		params["scheduled_publish_time"] = opts.ScheduleTime.Unix()
	}

	res, err := c.post(c.pageID+"/feed", params)
	if err != nil {
		return nil, err
	}

	id, err := getString(res, "id")
	if err != nil {
		return nil, fmt.Errorf("creating post: %w", err)
	}
	return &PostResult{
		ID:        id,
		URL:       fmt.Sprintf("https://facebook.com/%s", id),
		Scheduled: opts.ScheduleTime != nil,
	}, nil
}

// createPhotoPost uploads a photo and creates a post
func (c *Client) createPhotoPost(message, imagePath string, scheduleTime *time.Time) (*PostResult, error) {
	photoID, err := c.UploadPhoto(imagePath, message)
	if err != nil {
		return nil, fmt.Errorf("uploading photo: %w", err)
	}

	return &PostResult{
		ID:        photoID,
		URL:       fmt.Sprintf("https://facebook.com/%s", photoID),
		Scheduled: scheduleTime != nil,
	}, nil
}

// createVideoPost uploads a video and creates a post
func (c *Client) createVideoPost(message, videoPath string, scheduleTime *time.Time) (*PostResult, error) {
	videoID, err := c.UploadVideo(videoPath, message)
	if err != nil {
		return nil, fmt.Errorf("uploading video: %w", err)
	}

	return &PostResult{
		ID:        videoID,
		URL:       fmt.Sprintf("https://facebook.com/%s", videoID),
		Scheduled: scheduleTime != nil,
	}, nil
}

// ListPosts returns recent posts from the page
func (c *Client) ListPosts(limit int) ([]PostSummary, error) {
	if limit <= 0 {
		limit = 10
	}

	res, err := c.get(c.pageID+"/posts", fb.Params{
		"fields": "id,message,created_time,full_picture",
		"limit":  limit,
	})
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []struct {
			ID          string `facebook:"id"`
			Message     string `facebook:"message"`
			CreatedTime string `facebook:"created_time"`
			FullPicture string `facebook:"full_picture"`
		} `facebook:"data"`
	}

	if err := res.Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding posts: %w", err)
	}

	posts := make([]PostSummary, 0, len(response.Data))
	for _, p := range response.Data {
		t, _ := time.Parse("2006-01-02T15:04:05-0700", p.CreatedTime)
		posts = append(posts, PostSummary{
			ID:          p.ID,
			Message:     p.Message,
			CreatedTime: t,
			FullPicture: p.FullPicture,
		})
	}

	return posts, nil
}

// GetPost returns detailed info about a single post
func (c *Client) GetPost(postID string) (*PostDetail, error) {
	res, err := c.get(postID, fb.Params{
		"fields": "id,message,created_time,full_picture,shares,reactions.summary(true),comments.summary(true),permalink_url",
	})
	if err != nil {
		return nil, err
	}

	resolvedID, err := getString(res, "id")
	if err != nil {
		return nil, fmt.Errorf("reading post: %w", err)
	}
	post := &PostDetail{
		ID: resolvedID,
	}

	if msg, ok := res["message"].(string); ok {
		post.Message = msg
	}
	if ct, ok := res["created_time"].(string); ok {
		post.CreatedTime, _ = time.Parse("2006-01-02T15:04:05-0700", ct)
	}
	if pic, ok := res["full_picture"].(string); ok {
		post.FullPicture = pic
	}
	if url, ok := res["permalink_url"].(string); ok {
		post.URL = url
	}

	// Parse shares count
	if shares, ok := res["shares"].(map[string]interface{}); ok {
		if count, ok := shares["count"].(float64); ok {
			post.Shares = int(count)
		}
	}

	// Parse reactions summary
	if reactions, ok := res["reactions"].(map[string]interface{}); ok {
		if summary, ok := reactions["summary"].(map[string]interface{}); ok {
			if count, ok := summary["total_count"].(float64); ok {
				post.Reactions = int(count)
			}
		}
	}

	// Parse comments summary
	if comments, ok := res["comments"].(map[string]interface{}); ok {
		if summary, ok := comments["summary"].(map[string]interface{}); ok {
			if count, ok := summary["total_count"].(float64); ok {
				post.Comments = int(count)
			}
		}
	}

	return post, nil
}

// DeletePost removes a post
func (c *Client) DeletePost(postID string) error {
	return c.del(postID)
}
