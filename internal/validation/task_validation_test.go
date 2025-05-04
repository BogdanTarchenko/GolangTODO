package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todo/internal/domain/model"
)

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
