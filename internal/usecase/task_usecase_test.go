package usecase

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todo/internal/domain/model"
)

type mockTaskRepo struct {
	tasks map[string]*model.Task

	FindByIDFunc func(id string) (*model.Task, error)
}

func newMockTaskRepo() *mockTaskRepo {
	return &mockTaskRepo{tasks: make(map[string]*model.Task)}
}

func (m *mockTaskRepo) Create(task *model.Task) error {
	if _, exists := m.tasks[task.ID]; exists {
		return errors.New("already exists")
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepo) Update(task *model.Task) error {
	if _, exists := m.tasks[task.ID]; !exists {
		return errors.New("not found")
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepo) Delete(id string) error {
	if _, exists := m.tasks[id]; !exists {
		return errors.New("not found")
	}
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepo) FindByID(id string) (*model.Task, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	task, exists := m.tasks[id]
	if !exists {
		return nil, nil
	}
	return task, nil
}

func (m *mockTaskRepo) FindAll() ([]*model.Task, error) {
	var result []*model.Task
	for _, t := range m.tasks {
		result = append(result, t)
	}
	return result, nil
}

// --- Tests ---

// TestCreateTask_SetsFieldsAndSaves checks that a task is created correctly,
// macros are parsed, fields are filled, status and priority are set as expected.
func TestCreateTask_SetsFieldsAndSaves(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: create a task with macros in the title
	task := &model.Task{
		Title: "Test task !1 !before 01.01.2100",
	}

	// Act: call CreateTask
	created, err := uc.CreateTask(task)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "Test task", created.Title)
	assert.Equal(t, model.PriorityCritical, created.Priority)
	assert.NotNil(t, created.Deadline)
	assert.Equal(t, model.StatusActive, created.Status)
	assert.False(t, created.IsCompleted)
	assert.WithinDuration(t, time.Now().UTC(), created.CreatedAt, time.Second*2)
}

// TestUpdateTask_ChangesDeadlineAndRecalculatesStatus checks that when the deadline is changed,
// the task status is recalculated accordingly.
func TestUpdateTask_ChangesDeadlineAndRecalculatesStatus(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: create a task with a past deadline directly in the repo
	past := time.Now().Add(-24 * time.Hour)
	task := &model.Task{
		ID:        "1",
		Title:     "Old task",
		Deadline:  &past,
		Status:    model.StatusActive,
		Priority:  model.PriorityMedium,
		CreatedAt: time.Now().Add(-48 * time.Hour),
	}
	_ = repo.Create(task)

	// Act: change the deadline to the future
	future := time.Now().Add(24 * time.Hour)
	task.Deadline = &future
	task.IsCompleted = false

	updated, err := uc.UpdateTask(task)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, model.StatusActive, updated.Status)
	assert.Equal(t, future, *updated.Deadline)
}

// TestSetTaskCompletion_CompletedBeforeDeadline checks that a task becomes COMPLETED if finished before the deadline.
func TestSetTaskCompletion_CompletedBeforeDeadline(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: task with a future deadline
	future := time.Now().Add(24 * time.Hour)
	task := &model.Task{
		ID:          "2",
		Title:       "Complete me",
		Deadline:    &future,
		Status:      model.StatusActive,
		Priority:    model.PriorityMedium,
		CreatedAt:   time.Now(),
		IsCompleted: false,
	}
	_ = repo.Create(task)

	// Act: mark the task as completed
	task.IsCompleted = true
	updated, err := uc.SetTaskCompletion(task)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, model.StatusCompleted, updated.Status)
}

// TestSetTaskCompletion_CompletedAfterDeadline checks that a task becomes LATE if finished after the deadline.
func TestSetTaskCompletion_CompletedAfterDeadline(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: task with a past deadline
	past := time.Now().Add(-24 * time.Hour)
	task := &model.Task{
		ID:          "3",
		Title:       "Late task",
		Deadline:    &past,
		Status:      model.StatusActive,
		Priority:    model.PriorityMedium,
		CreatedAt:   time.Now(),
		IsCompleted: false,
	}
	_ = repo.Create(task)

	// Act: mark the task as completed
	task.IsCompleted = true
	updated, err := uc.SetTaskCompletion(task)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, model.StatusLate, updated.Status)
}

// TestListTasksWithFilter_PaginationAndSorting checks filtering, sorting, and pagination logic.
func TestListTasksWithFilter_PaginationAndSorting(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: create 5 tasks with different creation times
	now := time.Now()
	for i := 1; i <= 5; i++ {
		task := &model.Task{
			ID:        string(rune('a' + i)),
			Title:     "Task",
			Status:    model.StatusActive,
			Priority:  model.TaskPriority("MEDIUM"),
			CreatedAt: now.Add(time.Duration(i) * time.Minute),
		}
		_ = repo.Create(task)
	}

	// Act: request the first page with 2 tasks, sorted by created_at descending
	filter := &model.TaskFilter{
		Page:      1,
		PageSize:  2,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	tasks, total, err := uc.ListTasksWithFilter(filter)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, tasks, 2)
	assert.True(t, tasks[0].CreatedAt.After(tasks[1].CreatedAt))
}

// TestUpdateTask_RepoError checks that an error from the repository update is returned.
func TestUpdateTask_RepoError(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: create a valid task
	task := &model.Task{
		ID:        "update",
		Title:     "Valid title",
		Status:    model.StatusActive,
		Priority:  model.PriorityMedium,
		CreatedAt: time.Now(),
	}
	_ = repo.Create(task)
	_ = repo.Delete("update") // Simulate repo error by deleting the task before update

	// Act
	task.Title = "Updated title"
	updated, err := uc.UpdateTask(task)

	// Assert
	assert.Nil(t, updated)
	assert.Error(t, err)
}

// TestDeleteTask_Success checks that deleting an existing task works.
func TestDeleteTask_Success(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: create a task to delete
	task := &model.Task{
		ID:    "del",
		Title: "To delete",
	}
	_ = repo.Create(task)

	// Act
	err := uc.DeleteTask("del")

	// Assert
	assert.NoError(t, err)
}

// TestDeleteTask_RepoError checks that an error from the repository delete is returned.
func TestDeleteTask_RepoError(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Act: try to delete a non-existent task
	err := uc.DeleteTask("not-exist")

	// Assert
	assert.Error(t, err)
}

// TestGetTask_Success checks that getting an existing task works.
func TestGetTask_Success(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: create a task to get
	task := &model.Task{
		ID:    "get",
		Title: "To get",
	}
	_ = repo.Create(task)

	// Act
	got, err := uc.GetTask("get")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, "To get", got.Title)
}

// TestGetTask_NotFound checks that getting a non-existent task returns ErrTaskNotFound.
func TestGetTask_NotFound(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Act: try to get a non-existent task
	task, err := uc.GetTask("not-exist")

	// Assert
	assert.Nil(t, task)
	assert.Error(t, err)
}

// TestListTasksWithFilter_EmptyList checks that filtering on an empty repo returns an empty list.
func TestListTasksWithFilter_EmptyList(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Act: filter on an empty repo
	filter := &model.TaskFilter{
		Page:     1,
		PageSize: 10,
	}
	tasks, total, err := uc.ListTasksWithFilter(filter)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, tasks, 0)
}

// TestListTasksWithFilter_PaginationEdgeCase checks pagination when offset is out of range.
func TestListTasksWithFilter_PaginationEdgeCase(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: add one task
	task := &model.Task{
		ID:    "t1",
		Title: "Task 1",
	}
	_ = repo.Create(task)

	// Act: request page 2 with page size 10 (should be empty)
	filter := &model.TaskFilter{
		Page:     2,
		PageSize: 10,
	}
	tasks, total, err := uc.ListTasksWithFilter(filter)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, tasks, 0)
}

// TestSetTaskCompletion_RepoError checks that an error from the repository update is returned.
func TestSetTaskCompletion_RepoError(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: task not added to repo, so update will fail
	task := &model.Task{
		ID:    "not-exist",
		Title: "No such task",
	}
	task.IsCompleted = true

	// Act
	updated, err := uc.SetTaskCompletion(task)

	// Assert
	assert.Nil(t, updated)
	assert.Error(t, err)
}

// TestCreateTask_ValidationError checks that creating a task with invalid data returns a validation error.
func TestCreateTask_ValidationError(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: task with too short title
	task := &model.Task{
		Title: "abc",
	}

	// Act
	created, err := uc.CreateTask(task)

	// Assert
	assert.Nil(t, created)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title must be at least 4 characters")
}

// TestCreateTask_InvalidStatus checks that creating a task with invalid status returns a validation error.
func TestCreateTask_InvalidStatus(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: task with invalid status
	task := &model.Task{
		Title:  "Valid title",
		Status: "UNKNOWN",
	}

	// Act
	created, err := uc.CreateTask(task)

	// Assert
	assert.Nil(t, created)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid task status")
}

// TestCreateTask_InvalidPriority checks that creating a task with invalid priority returns a validation error.
func TestCreateTask_InvalidPriority(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: task with invalid priority
	task := &model.Task{
		Title:    "Valid title",
		Priority: "UNKNOWN",
	}

	// Act
	created, err := uc.CreateTask(task)

	// Assert
	assert.Nil(t, created)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid task priority")
}

// TestUpdateTask_ValidationError checks that updating a task with invalid data returns a validation error.
func TestUpdateTask_ValidationError(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: create a valid task
	task := &model.Task{
		ID:        "valid",
		Title:     "Valid title",
		Status:    model.StatusActive,
		Priority:  model.PriorityMedium,
		CreatedAt: time.Now(),
	}
	_ = repo.Create(task)

	// Act: try to update with invalid title
	task.Title = "a"
	updated, err := uc.UpdateTask(task)

	// Assert
	assert.Nil(t, updated)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title must be at least 4 characters")
}

// TestUpdateTask_FindByIDError checks that an error from repo.FindByID is handled.
func TestUpdateTask_FindByIDError(t *testing.T) {
	repo := newMockTaskRepo()
	repo.FindByIDFunc = func(id string) (*model.Task, error) {
		return nil, errors.New("db error")
	}
	uc := NewTaskUsecase(repo)

	// Arrange: task to update
	task := &model.Task{
		ID:    "any",
		Title: "Valid title",
	}

	// Act
	updated, err := uc.UpdateTask(task)

	// Assert
	assert.Nil(t, updated)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

// TestDeleteTask_FindByIDError checks that an error from repo.FindByID is handled in DeleteTask.
func TestDeleteTask_FindByIDError(t *testing.T) {
	repo := newMockTaskRepo()
	repo.FindByIDFunc = func(id string) (*model.Task, error) {
		return nil, errors.New("db error")
	}
	uc := NewTaskUsecase(repo)

	// Act
	err := uc.DeleteTask("any")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

// TestGetTask_FindByIDError checks that an error from repo.FindByID is handled in GetTask.
func TestGetTask_FindByIDError(t *testing.T) {
	repo := newMockTaskRepo()
	repo.FindByIDFunc = func(id string) (*model.Task, error) {
		return nil, errors.New("db error")
	}
	uc := NewTaskUsecase(repo)

	// Act
	task, err := uc.GetTask("any")

	// Assert
	assert.Nil(t, task)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

// TestListTasksWithFilter_InvalidSortBy checks that invalid sort_by returns a validation error.
func TestListTasksWithFilter_InvalidSortBy(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: filter with invalid sort_by
	filter := &model.TaskFilter{
		Page:     1,
		PageSize: 10,
		SortBy:   "unknown",
	}

	// Act
	tasks, total, err := uc.ListTasksWithFilter(filter)

	// Assert
	assert.Nil(t, tasks)
	assert.Equal(t, 0, total)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort_by field")
}

// TestListTasksWithFilter_InvalidSortOrder checks that invalid sort_order returns a validation error.
func TestListTasksWithFilter_InvalidSortOrder(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: filter with invalid sort_order
	filter := &model.TaskFilter{
		Page:      1,
		PageSize:  10,
		SortOrder: "invalid",
	}

	// Act
	tasks, total, err := uc.ListTasksWithFilter(filter)

	// Assert
	assert.Nil(t, tasks)
	assert.Equal(t, 0, total)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort_order value")
}

// TestListTasksWithFilter_InvalidPage checks that invalid page returns a validation error.
func TestListTasksWithFilter_InvalidPage(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: filter with invalid page
	filter := &model.TaskFilter{
		Page:     0,
		PageSize: 10,
	}

	// Act
	tasks, total, err := uc.ListTasksWithFilter(filter)

	// Assert
	assert.Nil(t, tasks)
	assert.Equal(t, 0, total)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "page must be greater than 0")
}

// TestListTasksWithFilter_InvalidPageSize checks that invalid page_size returns a validation error.
func TestListTasksWithFilter_InvalidPageSize(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: filter with invalid page_size
	filter := &model.TaskFilter{
		Page:     1,
		PageSize: 0,
	}

	// Act
	tasks, total, err := uc.ListTasksWithFilter(filter)

	// Assert
	assert.Nil(t, tasks)
	assert.Equal(t, 0, total)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "page_size must be greater than 0")
}

// TestCreateTask_DefaultStatusAndPriority checks that default status and priority are set if not provided.
func TestCreateTask_DefaultStatusAndPriority(t *testing.T) {
	repo := newMockTaskRepo()
	uc := NewTaskUsecase(repo)

	// Arrange: task with no status and no priority
	task := &model.Task{
		Title: "Task with defaults",
	}

	// Act
	created, err := uc.CreateTask(task)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, model.StatusActive, created.Status)
	assert.Equal(t, model.PriorityMedium, created.Priority)
}
