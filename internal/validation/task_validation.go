package validation

import (
	"strings"
	"time"
	"todo/internal/domain/model"
)

func isValidStatus(status model.TaskStatus) bool {
	switch status {
	case model.StatusActive, model.StatusCompleted, model.StatusOverdue, model.StatusLate:
		return true
	default:
		return false
	}
}

func isValidPriority(priority model.TaskPriority) bool {
	switch priority {
	case model.PriorityLow, model.PriorityMedium, model.PriorityHigh, model.PriorityCritical:
		return true
	default:
		return false
	}
}

func ValidateTask(t *model.Task) error {
	if len(strings.TrimSpace(t.Title)) < 4 {
		return NewValidationError("title must be at least 4 characters")
	}

	if t.Deadline != nil {
		if t.Deadline.Before(time.Now()) {
			return NewValidationError("deadline cannot be in the past")
		}
	}

	if !isValidStatus(t.Status) {
		return NewValidationError("invalid task status")
	}

	if !isValidPriority(t.Priority) {
		return NewValidationError("invalid task priority")
	}

	return nil
}
