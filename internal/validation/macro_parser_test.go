package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todo/internal/domain/model"
)

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

func TestParseTaskMacros_MultiplePriorityMacros(t *testing.T) {
	// Arrange
	title := "Task !3 !1"

	// Act
	result := ParseTaskMacros(title)

	// Assert
	assert.Equal(t, "Task !3", result.Title)
	assert.NotNil(t, result.Priority)
	assert.Equal(t, model.PriorityCritical, *result.Priority)
	assert.Nil(t, result.Deadline)
}

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
