package validation

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
	"todo/internal/domain/model"
)

// TestParseTaskMacros_PriorityMacro checks that priority macro (!1) is correctly parsed,
// sets CRITICAL priority and removes macro from title
func TestParseTaskMacros_PriorityMacro(t *testing.T) {
	// Arrange
	title := "Do something !1"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Do something", result.Title)
	assert.NotNil(t, result.Priority)
	assert.Equal(t, model.PriorityCritical, *result.Priority)
	assert.Nil(t, result.Deadline)
}

// TestParseTaskMacros_DeadlineMacro_DotFormat checks that deadline macro with dot format (DD.MM.YYYY)
// is correctly parsed and sets the deadline date
func TestParseTaskMacros_DeadlineMacro_DotFormat(t *testing.T) {
	// Arrange
	title := "Finish report !before 31.12.2099"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Finish report", result.Title)
	assert.Nil(t, result.Priority)
	assert.NotNil(t, result.Deadline)
	assert.Equal(t, 2099, result.Deadline.Year())
	assert.Equal(t, time.December, result.Deadline.Month())
	assert.Equal(t, 31, result.Deadline.Day())
}

// TestParseTaskMacros_DeadlineMacro_DashFormat checks that deadline macro with dash format (DD-MM-YYYY)
// is correctly parsed and sets the deadline date
func TestParseTaskMacros_DeadlineMacro_DashFormat(t *testing.T) {
	// Arrange
	title := "Finish report !before 31-12-2099"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Finish report", result.Title)
	assert.Nil(t, result.Priority)
	assert.NotNil(t, result.Deadline)
	assert.Equal(t, 2099, result.Deadline.Year())
	assert.Equal(t, time.December, result.Deadline.Month())
	assert.Equal(t, 31, result.Deadline.Day())
}

// TestParseTaskMacros_PriorityAndDeadline checks that both priority (!2) and deadline macros
// can be used together and are correctly parsed
func TestParseTaskMacros_PriorityAndDeadline(t *testing.T) {
	// Arrange
	title := "Urgent task !2 !before 01.01.2100"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Urgent task", result.Title)
	assert.NotNil(t, result.Priority)
	assert.Equal(t, model.PriorityHigh, *result.Priority)
	assert.NotNil(t, result.Deadline)
	assert.Equal(t, 2100, result.Deadline.Year())
	assert.Equal(t, time.January, result.Deadline.Month())
	assert.Equal(t, 1, result.Deadline.Day())
}

// TestParseTaskMacros_NoMacros verifies that title without any macros
// remains unchanged and no priority or deadline is set
func TestParseTaskMacros_NoMacros(t *testing.T) {
	// Arrange
	title := "Just a regular task"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Just a regular task", result.Title)
	assert.Nil(t, result.Priority)
	assert.Nil(t, result.Deadline)
}

// TestParseTaskMacros_InvalidDate checks that invalid date format in deadline macro
// is ignored and no deadline is set
func TestParseTaskMacros_InvalidDate(t *testing.T) {
	// Arrange
	title := "Task with bad date !before 99.99.9999"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Task with bad date", result.Title)
	assert.Nil(t, result.Priority)
	assert.Nil(t, result.Deadline)
}

// TestParseTaskMacros_MultiplePriorityMacros verifies that when multiple priority macros exist,
// only one priority is set and one macro remains in the title
func TestParseTaskMacros_MultiplePriorityMacros(t *testing.T) {
	// Arrange
	title := "Task !3 !1"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.NotNil(t, result.Priority)

	assert.True(t, *result.Priority == model.PriorityCritical || *result.Priority == model.PriorityMedium,
		"Priority should be either CRITICAL (!1) or MEDIUM (!3)")

	assert.False(t, strings.Contains(result.Title, "!1") && strings.Contains(result.Title, "!3"),
		"Title should not contain both priority macros")

	assert.Nil(t, result.Deadline)
}

// TestParseTaskMacros_ExtraSpaces checks that leading, trailing and extra spaces
// between macros are properly trimmed from title
func TestParseTaskMacros_ExtraSpaces(t *testing.T) {
	// Arrange
	title := "   Important !4   !before 02.02.2025   "

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Important", result.Title)
	assert.NotNil(t, result.Priority)
	assert.Equal(t, model.PriorityLow, *result.Priority)
	assert.NotNil(t, result.Deadline)
	assert.Equal(t, 2025, result.Deadline.Year())
	assert.Equal(t, time.February, result.Deadline.Month())
	assert.Equal(t, 2, result.Deadline.Day())
}

// TestParseTaskMacros_InvalidPriority verifies that invalid priority number (!5)
// is ignored and remains in title without setting priority
func TestParseTaskMacros_InvalidPriority(t *testing.T) {
	// Arrange
	title := "Task !5"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Task !5", result.Title)
	assert.Nil(t, result.Priority)
	assert.Nil(t, result.Deadline)
}

// TestParseTaskMacros_MultipleDeadlines checks that when multiple deadline macros exist,
// the first one is used and sets the corresponding deadline
func TestParseTaskMacros_MultipleDeadlines(t *testing.T) {
	// Arrange
	title := "Task !before 01.01.2024 !before 02.02.2024"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Task  !before 02.02.2024", result.Title)
	assert.NotNil(t, result.Deadline)
	assert.Equal(t, 2024, result.Deadline.Year())
	assert.Equal(t, time.January, result.Deadline.Month())
	assert.Equal(t, 1, result.Deadline.Day())
}
