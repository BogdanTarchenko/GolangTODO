package usecase

import (
	"log"
	"sort"
	"strings"
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

	// --- Status calculating ---
	if !task.IsCompleted {
		if task.Deadline != nil && now.After(*task.Deadline) {
			task.Status = model.StatusOverdue
		} else {
			task.Status = model.StatusActive
		}
	} else {
		if task.Deadline != nil && now.After(*task.Deadline) {
			task.Status = model.StatusLate
		} else {
			task.Status = model.StatusCompleted
		}
	}
	// --- Status calculating ---

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

func (u *taskUsecase) ListTasksWithFilter(filter *model.TaskFilter) ([]*model.Task, int, error) {
	if filter.Page <= 0 {
		return nil, 0, validation.NewValidationError("page must be greater than 0")
	}
	if filter.PageSize <= 0 {
		return nil, 0, validation.NewValidationError("page_size must be greater than 0")
	}

	tasks, err := u.repo.FindAll()
	if err != nil {
		return nil, 0, err
	}

	filtered := make([]*model.Task, 0)
	for _, t := range tasks {
		if filter.Status != "" && string(t.Status) != filter.Status {
			continue
		}
		if filter.Priority != "" && string(t.Priority) != filter.Priority {
			continue
		}
		filtered = append(filtered, t)
	}

	allowedSortFields := map[string]bool{
		"deadline":   true,
		"created_at": true,
		"priority":   true,
	}
	allowedSortOrders := map[string]bool{
		"asc":  true,
		"desc": true,
		"":     true,
	}

	var priorityOrder = map[model.TaskPriority]int{
		model.PriorityCritical: 4,
		model.PriorityHigh:     3,
		model.PriorityMedium:   2,
		model.PriorityLow:      1,
	}

	sortBy := filter.SortBy
	sortOrder := strings.ToLower(filter.SortOrder)

	if sortBy != "" && !allowedSortFields[sortBy] {
		return nil, 0, validation.NewValidationError("invalid sort_by field")
	}
	if !allowedSortOrders[sortOrder] {
		return nil, 0, validation.NewValidationError("invalid sort_order value")
	}

	if allowedSortFields[sortBy] {
		desc := sortOrder == "desc"
		sort.Slice(filtered, func(i, j int) bool {
			var less bool
			switch sortBy {
			case "deadline":
				if filtered[i].Deadline == nil || filtered[j].Deadline == nil {
					less = filtered[i].Deadline != nil
				} else {
					less = filtered[i].Deadline.Before(*filtered[j].Deadline)
				}
			case "created_at":
				less = filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
			case "priority":
				pi := priorityOrder[filtered[i].Priority]
				pj := priorityOrder[filtered[j].Priority]
				less = pi < pj
			}
			if desc {
				return !less
			}
			return less
		})
	}

	total := len(filtered)

	start := (filter.Page - 1) * filter.PageSize
	if start > total {
		start = total
	}
	end := start + filter.PageSize
	if filter.PageSize <= 0 || end > total {
		end = total
	}
	paged := filtered[start:end]

	return paged, total, nil
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
