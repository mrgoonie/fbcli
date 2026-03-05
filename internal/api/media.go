package api

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	fb "github.com/huandu/facebook/v2"
)

var (
	supportedImageExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	supportedVideoExts = map[string]bool{".mp4": true, ".mov": true, ".avi": true, ".flv": true, ".wmv": true, ".3gp": true}
	maxImageSize       = int64(10 * 1024 * 1024)  // 10MB
	maxVideoSize       = int64(1024 * 1024 * 1024) // 1GB
	chunkSize          = int64(4 * 1024 * 1024)    // 4MB
)

// UploadPhoto uploads a photo to the page
func (c *Client) UploadPhoto(imagePath, caption string) (string, error) {
	// Validate file
	ext := strings.ToLower(filepath.Ext(imagePath))
	if !supportedImageExts[ext] {
		return "", fmt.Errorf("unsupported image format: %s (supported: jpg, png, gif)", ext)
	}

	info, err := os.Stat(imagePath)
	if err != nil {
		return "", fmt.Errorf("image file not found: %s", imagePath)
	}
	if info.Size() > maxImageSize {
		return "", fmt.Errorf("image too large: %d MB (max: 10 MB)", info.Size()/1024/1024)
	}

	// Upload via multipart
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("opening image: %w", err)
	}
	defer file.Close()

	params := fb.Params{
		"message": caption,
		"source":  fb.File(imagePath),
	}

	res, err := c.post(c.pageID+"/photos", params)
	if err != nil {
		return "", err
	}

	id := res.Get("id").(string)
	return id, nil
}

// UploadVideo uploads a video using the resumable upload protocol
func (c *Client) UploadVideo(videoPath, caption string) (string, error) {
	// Validate file
	ext := strings.ToLower(filepath.Ext(videoPath))
	if !supportedVideoExts[ext] {
		return "", fmt.Errorf("unsupported video format: %s (supported: mp4, mov, avi, flv)", ext)
	}

	info, err := os.Stat(videoPath)
	if err != nil {
		return "", fmt.Errorf("video file not found: %s", videoPath)
	}
	if info.Size() > maxVideoSize {
		return "", fmt.Errorf("video too large: %d MB (max: 1024 MB)", info.Size()/1024/1024)
	}

	fileSize := info.Size()

	// Phase 1: Start upload
	startRes, err := c.post(c.pageID+"/videos", fb.Params{
		"upload_phase": "start",
		"file_size":    fileSize,
	})
	if err != nil {
		return "", fmt.Errorf("starting video upload: %w", err)
	}

	uploadSessionID := startRes.Get("upload_session_id").(string)
	videoID := startRes.Get("video_id")

	// Convert video ID to string
	var videoIDStr string
	switch v := videoID.(type) {
	case string:
		videoIDStr = v
	case float64:
		videoIDStr = fmt.Sprintf("%.0f", v)
	default:
		videoIDStr = fmt.Sprintf("%v", v)
	}

	// Phase 2: Transfer chunks
	file, err := os.Open(videoPath)
	if err != nil {
		return "", fmt.Errorf("opening video: %w", err)
	}
	defer file.Close()

	startOffset := int64(0)
	for startOffset < fileSize {
		currentChunkSize := chunkSize
		if startOffset+currentChunkSize > fileSize {
			currentChunkSize = fileSize - startOffset
		}

		chunk := make([]byte, currentChunkSize)
		n, err := file.ReadAt(chunk, startOffset)
		if err != nil && int64(n) != currentChunkSize {
			return "", fmt.Errorf("reading video chunk at offset %d: %w", startOffset, err)
		}

		if c.verbose {
			pct := float64(startOffset+currentChunkSize) / float64(fileSize) * 100
			fmt.Printf("  Uploading: %.0f%%\n", pct)
		}

		transferRes, err := c.post(c.pageID+"/videos", fb.Params{
			"upload_phase":      "transfer",
			"upload_session_id": uploadSessionID,
			"start_offset":     startOffset,
			"video_file_chunk":  fb.Data(filepath.Base(videoPath), bytes.NewReader(chunk)),
		})
		if err != nil {
			return "", fmt.Errorf("transferring video chunk: %w", err)
		}

		// Parse next start offset from response
		if so, ok := transferRes["start_offset"]; ok {
			switch v := so.(type) {
			case string:
				fmt.Sscanf(v, "%d", &startOffset)
			case float64:
				startOffset = int64(v)
			}
		} else {
			startOffset += currentChunkSize
		}
	}

	// Phase 3: Finish upload
	_, err = c.post(c.pageID+"/videos", fb.Params{
		"upload_phase":      "finish",
		"upload_session_id": uploadSessionID,
		"title":             caption,
	})
	if err != nil {
		return "", fmt.Errorf("finishing video upload: %w", err)
	}

	return videoIDStr, nil
}
