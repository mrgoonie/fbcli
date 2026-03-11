package cli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTruncateShortString(t *testing.T) {
	result := truncate("hello", 10)
	assert.Equal(t, "hello", result)
}

func TestTruncateLongString(t *testing.T) {
	result := truncate("hello world this is a long string", 10)
	assert.Equal(t, "hello w...", result)
	assert.Len(t, result, 10)
}

func TestTruncateExactLength(t *testing.T) {
	result := truncate("hello", 5)
	assert.Equal(t, "hello", result)
}

func TestTruncateOneCharMore(t *testing.T) {
	result := truncate("hello!", 5)
	assert.Equal(t, "he...", result)
	assert.Len(t, result, 5)
}

func TestTruncateNewlines(t *testing.T) {
	text := "hello\nworld\nthis is a long string"
	result := truncate(text, 20)
	// Newlines should be replaced with spaces
	assert.NotContains(t, result, "\n")
	assert.Contains(t, result, " ")
}

func TestTruncateNewlineReplacement(t *testing.T) {
	result := truncate("hello\nworld", 20)
	assert.Equal(t, "hello world", result)
}

func TestTruncateMinLength(t *testing.T) {
	// When max < 4, return string as-is to avoid panic on negative slice index
	result := truncate("test", 3)
	assert.Equal(t, "test", result)

	result = truncate("test", 0)
	assert.Equal(t, "test", result)

	// max=4 is the minimum that supports truncation with "..."
	result = truncate("hello", 4)
	assert.Equal(t, "h...", result)
}

func TestTruncateEmpty(t *testing.T) {
	result := truncate("", 10)
	assert.Equal(t, "", result)
}

func TestRelativeTimeJustNow(t *testing.T) {
	// Within last minute
	t1 := time.Now().Add(-10 * time.Second)
	result := relativeTime(t1)
	assert.Equal(t, "just now", result)
}

func TestRelativeTimeMinutesAgo(t *testing.T) {
	t1 := time.Now().Add(-30 * time.Minute)
	result := relativeTime(t1)
	assert.Equal(t, "30m ago", result)
}

func TestRelativeTimeMinutesAgoSingle(t *testing.T) {
	t1 := time.Now().Add(-1 * time.Minute)
	result := relativeTime(t1)
	assert.Equal(t, "1m ago", result)
}

func TestRelativeTimeHoursAgo(t *testing.T) {
	t1 := time.Now().Add(-5 * time.Hour)
	result := relativeTime(t1)
	assert.Equal(t, "5h ago", result)
}

func TestRelativeTimeHourAgoSingle(t *testing.T) {
	t1 := time.Now().Add(-1 * time.Hour)
	result := relativeTime(t1)
	assert.Equal(t, "1h ago", result)
}

func TestRelativeTimeDaysAgo(t *testing.T) {
	t1 := time.Now().Add(-7 * 24 * time.Hour)
	result := relativeTime(t1)
	assert.Equal(t, "7d ago", result)
}

func TestRelativeTimeDayAgoSingle(t *testing.T) {
	t1 := time.Now().Add(-1 * 24 * time.Hour)
	result := relativeTime(t1)
	assert.Equal(t, "1d ago", result)
}

func TestRelativeTimeJustBefore24Hours(t *testing.T) {
	t1 := time.Now().Add(-23 * time.Hour)
	result := relativeTime(t1)
	assert.Equal(t, "23h ago", result)
}

func TestRelativeTimeJustAfter24Hours(t *testing.T) {
	// Just after 24 hours should show days
	t1 := time.Now().Add(-24 * time.Hour - 1*time.Hour)
	result := relativeTime(t1)
	assert.Equal(t, "1d ago", result)
}

func TestRelativeTime30DaysAgo(t *testing.T) {
	t1 := time.Now().Add(-30 * 24 * time.Hour)
	result := relativeTime(t1)
	// At exactly 30 days, the boundary condition applies: d < 30*24*time.Hour is false
	// So it uses date format instead
	assert.NotContains(t, result, "d ago")
	assert.NotContains(t, result, "ago")
}

func TestRelativeTimeOver30Days(t *testing.T) {
	t1 := time.Now().Add(-40 * 24 * time.Hour)
	result := relativeTime(t1)
	// Should show full date format
	assert.NotContains(t, result, "d ago")
	assert.NotContains(t, result, "ago")
	// Check it has date format (with potential variations in day number)
	// e.g., "Jan 2, 2026" (11 chars) or "Jan 24, 2026" (12 chars)
	assert.GreaterOrEqual(t, len(result), 11)
	assert.LessOrEqual(t, len(result), 12)
}

func TestRelativeTimeMonthsOld(t *testing.T) {
	t1 := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	result := relativeTime(t1)
	// Should use date format
	assert.Equal(t, "Jan 1, 2020", result)
}

func TestRelativeTimeEdgeCaseJust59SecAgo(t *testing.T) {
	t1 := time.Now().Add(-59 * time.Second)
	result := relativeTime(t1)
	assert.Equal(t, "just now", result)
}

func TestRelativeTimeEdgeCaseJust61SecAgo(t *testing.T) {
	t1 := time.Now().Add(-61 * time.Second)
	result := relativeTime(t1)
	assert.Equal(t, "1m ago", result)
}

func TestRelativeTimeEdgeCaseJust59Min(t *testing.T) {
	t1 := time.Now().Add(-59 * time.Minute)
	result := relativeTime(t1)
	assert.Equal(t, "59m ago", result)
}

func TestRelativeTimeEdgeCaseJust61Min(t *testing.T) {
	t1 := time.Now().Add(-61 * time.Minute)
	result := relativeTime(t1)
	assert.Equal(t, "1h ago", result)
}

func TestRelativeTimeEdgeCaseJust23Hour59Min(t *testing.T) {
	t1 := time.Now().Add(-(23*time.Hour + 59*time.Minute))
	result := relativeTime(t1)
	assert.Equal(t, "23h ago", result)
}

func TestRelativeTimeEdgeCaseJust24Hour1Min(t *testing.T) {
	t1 := time.Now().Add(-(24*time.Hour + 1*time.Minute))
	result := relativeTime(t1)
	assert.Equal(t, "1d ago", result)
}

func TestTruncateMultipleNewlines(t *testing.T) {
	text := "line1\n\n\nline2\n\nline3"
	result := truncate(text, 50)
	assert.Equal(t, "line1   line2  line3", result)
}

func TestTruncateMixedContent(t *testing.T) {
	text := "This is a test\nwith newlines\nand long content that should be truncated"
	result := truncate(text, 30)
	assert.Equal(t, 30, len(result))
	assert.True(t, len(text) > len(result))
	// After newline replacement "This is a test with newlines and long..."
	// Truncated to 30 chars becomes "This is a test with newline..."
	assert.Equal(t, "This is a test with newline...", result)
}
