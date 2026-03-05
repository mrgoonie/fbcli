package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var forceFlag bool

var deleteCmd = &cobra.Command{
	Use:   "delete <post-id>",
	Short: "Delete a post from your Facebook Page",
	Long: `Delete a post by its ID. Prompts for confirmation unless --force is used.

Examples:
  fbcli delete 123456789_001
  fbcli delete 123456789_001 --force`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	deleteCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Skip confirmation prompt")
}

func runDelete(cmd *cobra.Command, args []string) error {
	postID := args[0]

	// Confirmation prompt
	if !forceFlag {
		fmt.Printf("Delete post %s? [y/N] ", postID)
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	client, err := requireAuth()
	if err != nil {
		return err
	}

	if err := client.DeletePost(postID); err != nil {
		return fmt.Errorf("deleting post: %w", err)
	}

	printSuccess("Post deleted.")
	return nil
}
