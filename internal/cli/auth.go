package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mrgoonie/fbcli/internal/api"
	"github.com/mrgoonie/fbcli/internal/auth"
	"github.com/mrgoonie/fbcli/internal/config"
	"github.com/mrgoonie/fbcli/internal/store"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "Login, check status, or logout from your Facebook Page.",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Facebook",
	RunE:  runLogin,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE:  runStatus,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored authentication",
	RunE:  runLogout,
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(logoutCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Prompt for App ID/Secret if not configured
	if cfg.AppID == "" || cfg.AppSecret == "" {
		reader := bufio.NewReader(os.Stdin)

		if cfg.AppID == "" {
			fmt.Print("Facebook App ID: ")
			cfg.AppID, _ = reader.ReadString('\n')
			cfg.AppID = strings.TrimSpace(cfg.AppID)
		}
		if cfg.AppSecret == "" {
			fmt.Print("Facebook App Secret: ")
			cfg.AppSecret, _ = reader.ReadString('\n')
			cfg.AppSecret = strings.TrimSpace(cfg.AppSecret)
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		printSuccess("App credentials saved")
	}

	port := 8910
	oauthCfg := auth.OAuthConfig{
		AppID:       cfg.AppID,
		AppSecret:   cfg.AppSecret,
		RedirectURI: auth.RedirectURI(port),
		Scopes:      auth.DefaultScopes(),
	}

	// Build auth URL and open browser
	authURL := auth.AuthURL(oauthCfg)
	fmt.Printf("\nOpening browser for authentication...\n")
	fmt.Printf("If browser doesn't open, visit:\n%s\n\n", cyan(authURL))

	openBrowser(authURL)

	// Listen for callback
	ctx := context.Background()
	fmt.Println("Waiting for authentication...")
	code, err := auth.ListenForCallback(ctx, port)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Exchange code for token
	fmt.Println("Exchanging authorization code...")
	shortToken, err := auth.ExchangeCode(ctx, oauthCfg, code)
	if err != nil {
		return err
	}

	// Extend to long-lived token
	fmt.Println("Getting long-lived token...")
	longToken, err := auth.ExtendToken(ctx, oauthCfg, shortToken)
	if err != nil {
		printWarning("Could not extend token, using short-lived token")
		longToken = shortToken
	}

	// Fetch user's pages
	fmt.Println("Fetching your pages...")
	pages, err := api.FetchUserPageTokens(longToken, isVerbose())
	if err != nil {
		return fmt.Errorf("fetching pages: %w", err)
	}

	if len(pages) == 0 {
		return fmt.Errorf("no Facebook Pages found. Make sure you manage at least one Page")
	}

	// Select page (auto-select if only one)
	var selected api.PageTokenInfo
	if len(pages) == 1 {
		selected = pages[0]
		fmt.Printf("Found page: %s\n", bold(selected.Name))
	} else {
		// Simple selection for now (Bubble Tea can be added later)
		fmt.Println("\nYour pages:")
		for i, p := range pages {
			fmt.Printf("  %d. %s (%s)\n", i+1, p.Name, p.Category)
		}
		fmt.Print("\nSelect page number: ")
		var choice int
		fmt.Scanln(&choice)
		if choice < 1 || choice > len(pages) {
			return fmt.Errorf("invalid selection")
		}
		selected = pages[choice-1]
	}

	// Store tokens
	ts := &store.TokenStore{
		PageToken: selected.AccessToken,
		UserToken: longToken,
		PageID:    selected.ID,
		PageName:  selected.Name,
	}

	if err := store.Save(ts); err != nil {
		return fmt.Errorf("saving tokens: %w", err)
	}

	// Save default page in config
	cfg.DefaultPageID = selected.ID
	cfg.DefaultPageName = selected.Name
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	printSuccess(fmt.Sprintf("Logged in as %s", bold(selected.Name)))
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	token, pageID, err := store.ResolveToken()
	if err != nil {
		printError("Not authenticated")
		fmt.Println("Run: fbcli auth login")
		return nil
	}

	client := api.NewClient(token, pageID, isVerbose())
	info, err := client.ValidateToken()
	if err != nil {
		printError("Token is invalid or expired")
		fmt.Println("Run: fbcli auth login")
		return nil
	}

	if isJSON() {
		printJSON(info)
		return nil
	}

	fmt.Printf("%s Authenticated\n", green("✓"))
	fmt.Printf("  Page: %s\n", bold(info.Name))
	fmt.Printf("  ID:   %s\n", info.ID)

	// Show token source
	if os.Getenv("FBCLI_PAGE_TOKEN") != "" {
		fmt.Printf("  Source: %s\n", cyan("FBCLI_PAGE_TOKEN env var"))
	} else {
		ts, _ := store.Load()
		if ts != nil {
			fmt.Printf("  Updated: %s\n", relativeTime(ts.UpdatedAt))
		}
	}

	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	if err := store.Clear(); err != nil {
		return err
	}
	printSuccess("Logged out")
	return nil
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}
