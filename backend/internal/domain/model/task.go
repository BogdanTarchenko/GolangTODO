package model

import (
	"time"
)

type TaskStatus string

const (
	StatusActive    TaskStatus = "ACTIVE"
	StatusCompleted TaskStatus = "COMPLETED"
	StatusOverdue   TaskStatus = "OVERDUE"
	StatusLate      TaskStatus = "LATE"
)

type TaskPriority string

const (
	PriorityLow      TaskPriority = "LOW"
	PriorityMedium   TaskPriority = "MEDIUM"
	PriorityHigh     TaskPriority = "HIGH"
	PriorityCritical TaskPriority = "CRITICAL"
)

type TaskFilter struct {
	Status    string
	Priority  string
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}

type Task struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	Deadline    *time.Time   `json:"deadline"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   *time.Time   `json:"updated_at"`
	IsCompleted bool         `json:"is_completed"`
}
