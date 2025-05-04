package usecase

import (
	"log"
	"time"

	"todo/internal/domain/model"
	"todo/internal/domain/repository"
	"todo/internal/validation"

	"github.com/google/uuid"
)

type taskUsecase struct {
	repo repository.TaskRepository
}

func NewTaskUsecase(repo repository.TaskRepository) *taskUsecase {
	return &taskUsecase{repo: repo}
}

func (u *taskUsecase) CreateTask(task *model.Task) (*model.Task, error) {
	now := time.Now().UTC()
	task.ID = uuid.New().String()

	// --- Macro parsing ---
	macros := validation.ParseTaskMacros(task.Title)
	task.Title = macros.Title
	if task.Priority == "" && macros.Priority != nil {
		task.Priority = *macros.Priority
	}
	if task.Deadline == nil && macros.Deadline != nil {
		task.Deadline = macros.Deadline
	}
	// --- Macro parsing ---

	if task.Status == "" {
		task.Status = model.StatusActive
	}
	if task.Priority == "" {
		task.Priority = model.PriorityMedium
	}
	task.CreatedAt = now

	if err := validation.ValidateTask(task); err != nil {
		return nil, err
	}

	if err := u.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (u *taskUsecase) UpdateTask(task *model.Task) (*model.Task, error) {
	existing, err := u.repo.FindByID(task.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, repository.ErrTaskNotFound
	}

	// --- Macro parsing ---
	macros := validation.ParseTaskMacros(task.Title)
	task.Title = macros.Title
	if task.Priority == "" && macros.Priority != nil {
		task.Priority = *macros.Priority
	}
	if task.Deadline == nil && macros.Deadline != nil {
		task.Deadline = macros.Deadline
	}
	// --- Macro parsing ---

	if err := validation.ValidateTask(task); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	task.UpdatedAt = &now

	if task.Priority == "" {
		task.Priority = model.PriorityMedium
	}

	if err := u.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
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

func (u *taskUsecase) ListTasks() ([]*model.Task, error) {
	return u.repo.FindAll()
}

func (u *taskUsecase) SetTaskCompletion(task *model.Task) (*model.Task, error) {
	now := time.Now().UTC()
	task.UpdatedAt = &now

	if task.IsCompleted {
		if task.Deadline != nil && now.After(*task.Deadline) {
			task.Status = model.StatusLate
		} else {
			task.Status = model.StatusCompleted
		}
	} else {
		if task.Deadline != nil && now.After(*task.Deadline) {
			task.Status = model.StatusOverdue
		} else {
			task.Status = model.StatusActive
		}
	}

	if err := u.repo.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (u *taskUsecase) UpdateOverdueTasks() error {
	tasks, err := u.repo.FindAll()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	for _, task := range tasks {
		if !task.IsCompleted && task.Deadline != nil && now.After(*task.Deadline) && task.Status == model.StatusActive {
			task.Status = model.StatusOverdue
			task.UpdatedAt = &now
			if err := u.repo.Update(task); err != nil {
				log.Printf("[CRON] Failed to update task %s to Overdue: %v", task.ID, err)
			} else {
				log.Printf("[CRON] Task %s marked as Overdue", task.ID)
			}
		}
	}
	return nil
}
