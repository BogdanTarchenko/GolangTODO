package repository

import (
	"database/sql"
	"testing"
	"time"
	"todo/internal/domain/model"
	"todo/internal/domain/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// newTestTask creates a test task with default values for testing
func newTestTask() *model.Task {
	now := time.Now().UTC()
	return &model.Task{
		ID:          "test-id",
		Title:       "Test Task",
		Description: nil,
		Deadline:    nil,
		Status:      model.StatusActive,
		Priority:    model.PriorityMedium,
		CreatedAt:   now,
		UpdatedAt:   &now,
		IsCompleted: false,
	}
}

// TestTaskPgRepository_Create checks that a task is successfully created in the database
func TestTaskPgRepository_Create(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewTaskPgRepository(db)
	task := newTestTask()

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(
			task.ID, task.Title, task.Description, task.Deadline, task.Status,
			task.Priority, task.CreatedAt, task.UpdatedAt, task.IsCompleted,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Create(task)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskPgRepository_Update checks that a task is successfully updated in the database
func TestTaskPgRepository_Update(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewTaskPgRepository(db)
	task := newTestTask()

	mock.ExpectExec("UPDATE tasks").
		WithArgs(
			task.Title, task.Description, task.Deadline, task.Status,
			task.Priority, task.UpdatedAt, task.IsCompleted, task.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(task)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskPgRepository_Update_NotFound checks that updating non-existent task returns ErrTaskNotFound
func TestTaskPgRepository_Update_NotFound(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewTaskPgRepository(db)
	task := newTestTask()

	mock.ExpectExec("UPDATE tasks").
		WithArgs(
			task.Title, task.Description, task.Deadline, task.Status,
			task.Priority, task.UpdatedAt, task.IsCompleted, task.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(task)

	// Assert
	assert.ErrorIs(t, err, repository.ErrTaskNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskPgRepository_Delete checks that a task is successfully deleted from the database
func TestTaskPgRepository_Delete(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewTaskPgRepository(db)

	mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").
		WithArgs("test-id").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete("test-id")

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskPgRepository_Delete_NotFound checks that deleting non-existent task returns ErrTaskNotFound
func TestTaskPgRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewTaskPgRepository(db)

	mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").
		WithArgs("not-exist").
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete("not-exist")

	// Assert
	assert.ErrorIs(t, err, repository.ErrTaskNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskPgRepository_FindByID checks that a task is successfully retrieved by ID with all fields populated correctly
func TestTaskPgRepository_FindByID(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewTaskPgRepository(db)
	now := time.Now().UTC()
	updatedAt := now
	description := "desc"
	deadline := now.Add(24 * time.Hour)

	mock.ExpectQuery("SELECT id, title, description, deadline, status, priority, created_at, updated_at, is_completed FROM tasks WHERE id = \\$1").
		WithArgs("test-id").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "description", "deadline", "status", "priority", "created_at", "updated_at", "is_completed",
		}).AddRow(
			"test-id", "Test Task", description, deadline, model.StatusActive, model.PriorityMedium, now, updatedAt, false,
		))

	// Act
	task, err := repo.FindByID("test-id")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "test-id", task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.NotNil(t, task.Description)
	assert.Equal(t, description, *task.Description)
	assert.NotNil(t, task.Deadline)
	assert.WithinDuration(t, deadline, *task.Deadline, time.Second)
	assert.Equal(t, model.StatusActive, task.Status)
	assert.Equal(t, model.PriorityMedium, task.Priority)
	assert.WithinDuration(t, now, task.CreatedAt, time.Second)
	assert.NotNil(t, task.UpdatedAt)
	assert.WithinDuration(t, updatedAt, *task.UpdatedAt, time.Second)
	assert.False(t, task.IsCompleted)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskPgRepository_FindByID_NotFound checks that finding non-existent task returns ErrTaskNotFound
func TestTaskPgRepository_FindByID_NotFound(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewTaskPgRepository(db)

	mock.ExpectQuery("SELECT id, title, description, deadline, status, priority, created_at, updated_at, is_completed FROM tasks WHERE id = \\$1").
		WithArgs("not-exist").
		WillReturnError(sql.ErrNoRows)

	// Act
	task, err := repo.FindByID("not-exist")

	// Assert
	assert.ErrorIs(t, err, repository.ErrTaskNotFound)
	assert.Nil(t, task)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTaskPgRepository_FindAll checks that all tasks are successfully retrieved from the database
func TestTaskPgRepository_FindAll(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewTaskPgRepository(db)
	now := time.Now().UTC()
	updatedAt := now
	description := "desc"
	deadline := now.Add(24 * time.Hour)

	mock.ExpectQuery("SELECT id, title, description, deadline, status, priority, created_at, updated_at, is_completed FROM tasks").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "description", "deadline", "status", "priority", "created_at", "updated_at", "is_completed",
		}).AddRow(
			"test-id", "Test Task", description, deadline, model.StatusActive, model.PriorityMedium, now, updatedAt, false,
		))

	// Act
	tasks, err := repo.FindAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	task := tasks[0]
	assert.Equal(t, "test-id", task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.NotNil(t, task.Description)
	assert.Equal(t, description, *task.Description)
	assert.NotNil(t, task.Deadline)
	assert.WithinDuration(t, deadline, *task.Deadline, time.Second)
	assert.Equal(t, model.StatusActive, task.Status)
	assert.Equal(t, model.PriorityMedium, task.Priority)
	assert.WithinDuration(t, now, task.CreatedAt, time.Second)
	assert.NotNil(t, task.UpdatedAt)
	assert.WithinDuration(t, updatedAt, *task.UpdatedAt, time.Second)
	assert.False(t, task.IsCompleted)
	assert.NoError(t, mock.ExpectationsWereMet())
}
