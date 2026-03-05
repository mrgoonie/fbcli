package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	defaultPort    = 8910
	callbackPath   = "/callback"
	listenTimeout  = 2 * time.Minute
)

// successHTML is served to the browser after successful auth
const successHTML = `<!DOCTYPE html>
<html><head><title>fbcli - Auth Success</title>
<style>body{font-family:system-ui;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;background:#f0f2f5}
.card{background:white;padding:2rem;border-radius:12px;box-shadow:0 2px 8px rgba(0,0,0,.1);text-align:center;max-width:400px}
h1{color:#1877f2;font-size:1.5rem}p{color:#65676b}</style></head>
<body><div class="card"><h1>✓ Authenticated</h1><p>You can close this tab and return to your terminal.</p></div></body></html>`

// RedirectURI returns the OAuth redirect URI for the given port
func RedirectURI(port int) string {
	return fmt.Sprintf("http://localhost:%d%s", port, callbackPath)
}

// ListenForCallback starts a local HTTP server and waits for the OAuth callback.
// Returns the authorization code or an error.
func ListenForCallback(ctx context.Context, port int) (string, error) {
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		// Check for errors from Facebook
		if errMsg := r.URL.Query().Get("error_description"); errMsg != "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s", errMsg)
			errCh <- fmt.Errorf("OAuth error: %s", errMsg)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Missing authorization code")
			errCh <- fmt.Errorf("missing authorization code in callback")
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, successHTML)
		codeCh <- code
	})

	// Find available port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// Try next ports
		for p := port + 1; p <= port+10; p++ {
			listener, err = net.Listen("tcp", fmt.Sprintf(":%d", p))
			if err == nil {
				port = p
				break
			}
		}
		if err != nil {
			return "", fmt.Errorf("no available port found near %d: %w", defaultPort, err)
		}
	}

	server := &http.Server{Handler: mux}

	// Start server in background
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for code, error, or timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, listenTimeout)
	defer cancel()

	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer shutdownCancel()
		server.Shutdown(shutdownCtx)
	}()

	select {
	case code := <-codeCh:
		return code, nil
	case err := <-errCh:
		return "", err
	case <-timeoutCtx.Done():
		return "", fmt.Errorf("authentication timed out after %v", listenTimeout)
	}
}
