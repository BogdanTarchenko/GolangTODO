package validation

import (
	"errors"
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
		return errors.New("title must be at least 4 characters")
	}

	if t.Deadline != nil {
		now := time.Now()
		if t.Deadline.Before(now) {
			return errors.New("deadline cannot be in the past")
		}
	}

	if !isValidStatus(t.Status) {
		return errors.New("invalid task status")
	}

	if !isValidPriority(t.Priority) {
		return errors.New("invalid task priority")
	}

	return nil
}
