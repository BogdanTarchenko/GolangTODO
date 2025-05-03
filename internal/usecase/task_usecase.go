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

func (u *taskUsecase) UpdateTask(task *model.Task) error {
	existing, err := u.repo.FindByID(task.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return repository.ErrTaskNotFound
	}

	if err := validation.ValidateTask(task); err != nil {
		return err
	}

	now := time.Now()
	task.UpdatedAt = &now

	if task.Priority == "" {
		task.Priority = model.PriorityMedium
	}

	return u.repo.Update(task)
}

func (u *taskUsecase) DeleteTask(id string) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return repository.ErrTaskNotFound
	}

	return u.repo.Delete(id)
}

func (u *taskUsecase) GetTask(id string) (*model.Task, error) {
	task, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, repository.ErrTaskNotFound
	}
	return task, nil
}
