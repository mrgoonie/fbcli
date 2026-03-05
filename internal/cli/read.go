package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read <post-id>",
	Short: "Show details of a specific post",
	Long: `Display full post content with engagement metrics.

Examples:
  fbcli read 123456789_001
  fbcli read 123456789_001 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runRead,
}

func runRead(cmd *cobra.Command, args []string) error {
	postID := args[0]

	client, err := requireAuth()
	if err != nil {
		return err
	}

	post, err := client.GetPost(postID)
	if err != nil {
		return fmt.Errorf("reading post: %w", err)
	}

	if isJSON() {
		printJSON(post)
		return nil
	}

	fmt.Printf("%s %s\n", bold("Post:"), post.ID)
	fmt.Printf("%s %s\n", bold("Date:"), post.CreatedTime.Format("January 2, 2006 at 3:04 PM"))

	if post.Message != "" {
		fmt.Printf("\n%s\n", post.Message)
	}

	fmt.Printf("\n%s %d  %s %d  %s %d\n",
		bold("Reactions:"), post.Reactions,
		bold("Comments:"), post.Comments,
		bold("Shares:"), post.Shares,
	)

	if post.URL != "" {
		fmt.Printf("%s %s\n", bold("URL:"), post.URL)
	}

	return nil
}
