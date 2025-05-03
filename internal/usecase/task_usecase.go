package usecase

import (
	"time"
	"todo/internal/domain/model"
	"todo/internal/domain/repository"
	"todo/internal/validation"
)

type taskUsecase struct {
	repo repository.TaskRepository
}

func NewTaskUsecase(repo repository.TaskRepository) *taskUsecase {
	return &taskUsecase{repo: repo}
}

func (u *taskUsecase) CreateTask(task *model.Task) error {
	if err := validation.ValidateTask(task); err != nil {
		return err
	}
	task.CreatedAt = time.Now()
	task.Status = model.StatusActive
	if task.Priority == "" {
		task.Priority = model.PriorityMedium
	}

	return u.repo.Create(task)
}
