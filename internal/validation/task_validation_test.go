package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todo/internal/domain/model"
)

// TestValidateTask_TitleTooShort checks that validation fails when task title
// is less than 4 characters long
func TestValidateTask_TitleTooShort(t *testing.T) {
	// Arrange
	task := &model.Task{
		Title:    "abc",
		Status:   model.StatusActive,
		Priority: model.PriorityMedium,
	}

	// Act
	err := ValidateTask(task)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title must be at least 4 characters")
}

// TestValidateTask_DeadlineInPast checks that validation fails when task deadline
// is set to a past date
func TestValidateTask_DeadlineInPast(t *testing.T) {
	// Arrange
	past := time.Now().Add(-time.Hour)
	task := &model.Task{
		Title:    "Valid title",
		Deadline: &past,
		Status:   model.StatusActive,
		Priority: model.PriorityMedium,
	}

	// Act
	err := ValidateTask(task)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deadline cannot be in the past")
}

// TestValidateTask_InvalidStatus checks that validation fails when task status
// is not one of the allowed values
func TestValidateTask_InvalidStatus(t *testing.T) {
	// Arrange
	task := &model.Task{
		Title:    "Valid title",
		Status:   "UNKNOWN",
		Priority: model.PriorityMedium,
	}

	// Act
	err := ValidateTask(task)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid task status")
}

// TestValidateTask_InvalidPriority checks that validation fails when task priority
// is not one of the allowed values
func TestValidateTask_InvalidPriority(t *testing.T) {
	// Arrange
	task := &model.Task{
		Title:    "Valid title",
		Status:   model.StatusActive,
		Priority: "UNKNOWN",
	}

	// Act
	err := ValidateTask(task)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid task priority")
}

// TestValidateTask_ValidTask verifies that validation passes for a task with
// valid title, future deadline, correct status and priority
func TestValidateTask_ValidTask(t *testing.T) {
	// Arrange
	future := time.Now().Add(time.Hour)
	task := &model.Task{
		Title:    "Valid title",
		Deadline: &future,
		Status:   model.StatusActive,
		Priority: model.PriorityMedium,
	}

	// Act
	err := ValidateTask(task)

	// Assert
	assert.NoError(t, err)
}
