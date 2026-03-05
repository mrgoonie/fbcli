package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	verbose   bool
	jsonOut   bool
	pageIDFlg string
)

// SetVersion sets the version string from build-time ldflags
func SetVersion(v string) {
	version = v
}

var rootCmd = &cobra.Command{
	Use:   "fbcli",
	Short: "Facebook Page CLI - post and manage your Facebook Page from the terminal",
	Long: `fbcli is a command-line tool for managing Facebook Pages.
Create posts, upload media, schedule content, and view analytics — all from your terminal.`,
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Show detailed request/response info")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().StringVar(&pageIDFlg, "page", "", "Override page ID")

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(postCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(deleteCmd)
}

// Execute runs the root command
func Execute() error {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("fbcli %s\n", version))
	return rootCmd.Execute()
}

// isVerbose returns the verbose flag value
func isVerbose() bool {
	return verbose
}

// isJSON returns the json output flag value
func isJSON() bool {
	return jsonOut
}
