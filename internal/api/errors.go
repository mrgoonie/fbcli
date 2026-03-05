package api

import (
	"fmt"

	fb "github.com/huandu/facebook/v2"
)

// APIError represents a structured Facebook API error
type APIError struct {
	Code    int
	Message string
	Type    string
	Subcode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Facebook API error %d: %s", e.Code, e.Message)
}

// Hint returns a user-friendly hint for common error codes
func (e *APIError) Hint() string {
	switch e.Code {
	case 190:
		return "Token is invalid or expired. Run: fbcli auth login"
	case 100:
		return "Invalid parameter in request"
	case 200:
		return "Permission denied. Ensure your app has the required permissions."
	case 368:
		return "Content blocked by Facebook policy"
	case 4:
		return "Rate limit exceeded. Wait a few minutes and try again."
	default:
		return ""
	}
}

// wrapFBError converts a facebook SDK error into our APIError
func wrapFBError(err error) error {
	if err == nil {
		return nil
	}

	if fbErr, ok := err.(*fb.Error); ok {
		apiErr := &APIError{
			Code:    fbErr.Code,
			Message: fbErr.Message,
			Type:    fbErr.Type,
		}
		hint := apiErr.Hint()
		if hint != "" {
			return fmt.Errorf("%w\nHint: %s", apiErr, hint)
		}
		return apiErr
	}

	return err
}
