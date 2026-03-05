package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	fb "github.com/huandu/facebook/v2"
)

func TestAPIErrorError(t *testing.T) {
	err := &APIError{
		Code:    190,
		Message: "Invalid token",
		Type:    "OAuthException",
	}

	expected := "Facebook API error 190: Invalid token"
	assert.Equal(t, expected, err.Error())
}

func TestAPIErrorHintInvalidToken(t *testing.T) {
	err := &APIError{
		Code:    190,
		Message: "Token is invalid or expired",
	}

	hint := err.Hint()
	assert.Equal(t, "Token is invalid or expired. Run: fbcli auth login", hint)
}

func TestAPIErrorHintInvalidParameter(t *testing.T) {
	err := &APIError{
		Code:    100,
		Message: "Invalid parameter",
	}

	hint := err.Hint()
	assert.Equal(t, "Invalid parameter in request", hint)
}

func TestAPIErrorHintPermissionDenied(t *testing.T) {
	err := &APIError{
		Code:    200,
		Message: "Permission denied",
	}

	hint := err.Hint()
	assert.Equal(t, "Permission denied. Ensure your app has the required permissions.", hint)
}

func TestAPIErrorHintContentBlocked(t *testing.T) {
	err := &APIError{
		Code:    368,
		Message: "Content blocked",
	}

	hint := err.Hint()
	assert.Equal(t, "Content blocked by Facebook policy", hint)
}

func TestAPIErrorHintRateLimit(t *testing.T) {
	err := &APIError{
		Code:    4,
		Message: "Rate limit exceeded",
	}

	hint := err.Hint()
	assert.Equal(t, "Rate limit exceeded. Wait a few minutes and try again.", hint)
}

func TestAPIErrorHintUnknownCode(t *testing.T) {
	err := &APIError{
		Code:    999,
		Message: "Unknown error",
	}

	hint := err.Hint()
	assert.Empty(t, hint)
}

func TestWrapFBErrorNil(t *testing.T) {
	err := wrapFBError(nil)
	assert.Nil(t, err)
}

func TestWrapFBErrorRegularError(t *testing.T) {
	originalErr := fmt.Errorf("some error")
	err := wrapFBError(originalErr)
	assert.Equal(t, originalErr, err)
}

func TestWrapFBErrorWithHint(t *testing.T) {
	fbErr := &fb.Error{
		Code:    190,
		Message: "Token invalid",
		Type:    "OAuthException",
	}

	err := wrapFBError(fbErr)
	assert.Error(t, err)
	errStr := err.Error()
	// Should contain the hint
	assert.Contains(t, errStr, "Token is invalid or expired")
	assert.Contains(t, errStr, "fbcli auth login")
}

func TestWrapFBErrorWithoutHint(t *testing.T) {
	fbErr := &fb.Error{
		Code:    999,
		Message: "Unknown error",
		Type:    "UnknownException",
	}

	err := wrapFBError(fbErr)
	assert.Error(t, err)

	// Should be APIError
	apiErr, ok := err.(*APIError)
	assert.True(t, ok)
	assert.Equal(t, 999, apiErr.Code)
	assert.Equal(t, "Unknown error", apiErr.Message)
	assert.Equal(t, "UnknownException", apiErr.Type)
}

func TestWrapFBErrorPermissionDenied(t *testing.T) {
	fbErr := &fb.Error{
		Code:    200,
		Message: "Permission denied",
		Type:    "OAuthException",
	}

	err := wrapFBError(fbErr)
	assert.Error(t, err)
	errStr := err.Error()
	assert.Contains(t, errStr, "Permission denied")
	assert.Contains(t, errStr, "required permissions")
}

func TestAPIErrorStructure(t *testing.T) {
	err := &APIError{
		Code:    100,
		Message: "Test message",
		Type:    "TestType",
		Subcode: 42,
	}

	assert.Equal(t, 100, err.Code)
	assert.Equal(t, "Test message", err.Message)
	assert.Equal(t, "TestType", err.Type)
	assert.Equal(t, 42, err.Subcode)
}
