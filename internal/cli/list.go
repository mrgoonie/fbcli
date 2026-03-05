package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var limitFlag int

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent posts from your Facebook Page",
	Long: `List recent posts with ID, date, and message preview.

Examples:
  fbcli list
  fbcli list -n 20
  fbcli list --json`,
	RunE: runList,
}

func init() {
	listCmd.Flags().IntVarP(&limitFlag, "limit", "n", 10, "Number of posts to show")
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := requireAuth()
	if err != nil {
		return err
	}

	posts, err := client.ListPosts(limitFlag)
	if err != nil {
		return fmt.Errorf("listing posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found.")
		return nil
	}

	if isJSON() {
		printJSON(posts)
		return nil
	}

	headers := []string{"ID", "DATE", "MESSAGE"}
	rows := make([][]string, 0, len(posts))
	for _, p := range posts {
		rows = append(rows, []string{
			p.ID,
			relativeTime(p.CreatedTime),
			truncate(p.Message, 50),
		})
	}

	printTable(headers, rows)
	return nil
}
