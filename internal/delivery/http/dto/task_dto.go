package dto

import "time"

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=4" example:"Купить продукты"`
	Description *string    `json:"description" example:"Купить хлеб, молоко и яйца"`
	Deadline    *time.Time `json:"deadline" example:"2025-06-01T18:00:00Z"`
	Priority    string     `json:"priority" example:"MEDIUM"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title" example:"Обновлённая задача"`
	Description *string    `json:"description" example:"Новое описание"`
	Deadline    *time.Time `json:"deadline" example:"2025-06-02T18:00:00Z"`
	Priority    *string    `json:"priority" example:"HIGH"`
}

type TaskResponse struct {
	ID          string     `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title       string     `json:"title" example:"Купить продукты"`
	Description *string    `json:"description" example:"Купить хлеб, молоко и яйца"`
	Deadline    *time.Time `json:"deadline" example:"2025-06-01T18:00:00Z"`
	Status      string     `json:"status" example:"ACTIVE"`
	Priority    string     `json:"priority" example:"MEDIUM"`
	CreatedAt   time.Time  `json:"created_at" example:"2025-05-04T21:00:00Z"`
	UpdatedAt   *time.Time `json:"updated_at" example:"2025-05-04T21:30:00Z"`
}

type UpdateTaskStatusRequest struct {
	IsCompleted bool `json:"is_completed" example:"true"`
}
