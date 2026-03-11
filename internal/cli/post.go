package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mrgoonie/fbcli/internal/api"
	"github.com/mrgoonie/fbcli/internal/store"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	imageFlag    string
	videoFlag    string
	scheduleFlag string
	linkFlag     string
)

var postCmd = &cobra.Command{
	Use:   "post [message]",
	Short: "Create a new post on your Facebook Page",
	Long: `Create a text, photo, video, or link post on your Facebook Page.

Examples:
  fbcli post "Hello world!"
  fbcli post -i photo.jpg "Check this out"
  fbcli post -v video.mp4 "Watch this"
  fbcli post --schedule "2026-03-10T10:00" "Coming soon"
  fbcli post -l https://example.com "Read this article"
  echo "Hello" | fbcli post`,
	RunE: runPost,
}

func init() {
	postCmd.Flags().StringVarP(&imageFlag, "image", "i", "", "Path to image file (jpg, png, gif)")
	postCmd.Flags().StringVarP(&videoFlag, "video", "v", "", "Path to video file (mp4, mov, avi)")
	postCmd.Flags().StringVarP(&scheduleFlag, "schedule", "s", "", "Schedule time (ISO 8601: 2026-03-10T10:00)")
	postCmd.Flags().StringVarP(&linkFlag, "link", "l", "", "URL to attach")
}

func runPost(cmd *cobra.Command, args []string) error {
	// Get message from args or stdin
	var message string
	if len(args) > 0 {
		message = args[0]
	} else if !term.IsTerminal(int(os.Stdin.Fd())) {
		data, err := io.ReadAll(io.LimitReader(os.Stdin, 64*1024)) // 64KB max
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
		message = string(data)
	}

	if message == "" && imageFlag == "" && videoFlag == "" {
		return fmt.Errorf("message is required. Usage: fbcli post \"your message\"")
	}

	// Validate flags
	if imageFlag != "" && videoFlag != "" {
		return fmt.Errorf("cannot use both --image and --video in the same post")
	}

	// Parse schedule time
	var scheduleTime *time.Time
	if scheduleFlag != "" {
		t, err := time.ParseInLocation("2006-01-02T15:04", scheduleFlag, time.Local)
		if err != nil {
			return fmt.Errorf("invalid schedule time format. Use: 2026-03-10T10:00")
		}
		if t.Before(time.Now().Add(10 * time.Minute)) {
			return fmt.Errorf("schedule time must be at least 10 minutes in the future")
		}
		if t.After(time.Now().Add(6 * 30 * 24 * time.Hour)) {
			return fmt.Errorf("schedule time must be within 6 months")
		}
		scheduleTime = &t
	}

	client, err := requireAuth()
	if err != nil {
		return err
	}

	opts := api.PostOptions{
		ImagePath:    imageFlag,
		VideoPath:    videoFlag,
		ScheduleTime: scheduleTime,
		Link:         linkFlag,
	}

	result, err := client.CreatePost(message, opts)
	if err != nil {
		return fmt.Errorf("creating post: %w", err)
	}

	if isJSON() {
		printJSON(result)
		return nil
	}

	if result.Scheduled {
		printSuccess(fmt.Sprintf("Post scheduled: %s", result.URL))
	} else {
		printSuccess(fmt.Sprintf("Post created: %s", result.URL))
	}

	return nil
}

// requireAuth resolves token and creates API client
func requireAuth() (*api.Client, error) {
	token, pageID, err := store.ResolveToken()
	if err != nil {
		return nil, err
	}

	// Override page ID if flag set
	if pageIDFlg != "" {
		pageID = pageIDFlg
	}

	return api.NewClient(token, pageID, isVerbose()), nil
}
