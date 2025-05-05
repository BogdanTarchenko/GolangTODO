package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo/internal/delivery/http/dto"
	"todo/internal/delivery/http/middleware"
	"todo/internal/domain/model"
	"todo/internal/domain/repository"
	"todo/internal/validation"
)

// --- Mock Usecase ---

type mockTaskUsecase struct {
	CreateTaskFunc          func(*model.Task) (*model.Task, error)
	ListTasksWithFilterFunc func(*model.TaskFilter) ([]*model.Task, int, error)
	GetTaskFunc             func(string) (*model.Task, error)
	UpdateTaskFunc          func(*model.Task) (*model.Task, error)
	DeleteTaskFunc          func(string) error
	SetTaskCompletionFunc   func(*model.Task) (*model.Task, error)
	UpdateOverdueTasksFunc  func() error
}

func (m *mockTaskUsecase) CreateTask(t *model.Task) (*model.Task, error) {
	return m.CreateTaskFunc(t)
}
func (m *mockTaskUsecase) ListTasksWithFilter(f *model.TaskFilter) ([]*model.Task, int, error) {
	return m.ListTasksWithFilterFunc(f)
}
func (m *mockTaskUsecase) GetTask(id string) (*model.Task, error) {
	return m.GetTaskFunc(id)
}
func (m *mockTaskUsecase) UpdateTask(t *model.Task) (*model.Task, error) {
	return m.UpdateTaskFunc(t)
}
func (m *mockTaskUsecase) DeleteTask(id string) error {
	return m.DeleteTaskFunc(id)
}
func (m *mockTaskUsecase) SetTaskCompletion(t *model.Task) (*model.Task, error) {
	return m.SetTaskCompletionFunc(t)
}
func (m *mockTaskUsecase) UpdateOverdueTasks() error {
	if m.UpdateOverdueTasksFunc != nil {
		return m.UpdateOverdueTasksFunc()
	}
	return nil
}

// --- Helpers ---

func setupRouter(handler *TaskHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	handler.RegisterRoutes(r)
	return r
}

func newTestTask() *model.Task {
	now := time.Now().UTC()
	return &model.Task{
		ID:          "1",
		Title:       "Test",
		Description: nil,
		Deadline:    nil,
		Status:      model.StatusActive,
		Priority:    model.PriorityMedium,
		CreatedAt:   now,
		UpdatedAt:   &now,
		IsCompleted: false,
	}
}

// --- Tests ---

// TestTaskHandler_CreateTask_Success checks that a valid task is created successfully
func TestTaskHandler_CreateTask_Success(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		CreateTaskFunc: func(task *model.Task) (*model.Task, error) {
			task.ID = "1"
			return task, nil
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	reqBody := dto.CreateTaskRequest{Title: "Test"}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
}

// TestTaskHandler_CreateTask_ValidationError checks that invalid task data returns a validation error
func TestTaskHandler_CreateTask_ValidationError(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		CreateTaskFunc: func(task *model.Task) (*model.Task, error) {
			return nil, validation.NewValidationError("validation error")
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	reqBody := dto.CreateTaskRequest{Title: "bad"}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestTaskHandler_ListTasks_Success checks that tasks are listed successfully
func TestTaskHandler_ListTasks_Success(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		ListTasksWithFilterFunc: func(f *model.TaskFilter) ([]*model.Task, int, error) {
			return []*model.Task{newTestTask()}, 1, nil
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestTaskHandler_ListTasks_Error checks that internal error is handled properly
func TestTaskHandler_ListTasks_Error(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		ListTasksWithFilterFunc: func(f *model.TaskFilter) ([]*model.Task, int, error) {
			return nil, 0, errors.New("internal error")
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestTaskHandler_GetTask_Success checks that a task is retrieved successfully by ID
func TestTaskHandler_GetTask_Success(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		GetTaskFunc: func(id string) (*model.Task, error) {
			return newTestTask(), nil
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/1", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestTaskHandler_GetTask_NotFound checks that requesting non-existent task returns not found error
func TestTaskHandler_GetTask_NotFound(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		GetTaskFunc: func(id string) (*model.Task, error) {
			return nil, repository.ErrTaskNotFound
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/1", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestTaskHandler_UpdateTask_Success checks that a task is updated successfully
func TestTaskHandler_UpdateTask_Success(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		GetTaskFunc: func(id string) (*model.Task, error) {
			return newTestTask(), nil
		},
		UpdateTaskFunc: func(task *model.Task) (*model.Task, error) {
			task.Title = "Updated"
			return task, nil
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	reqBody := dto.UpdateTaskRequest{Title: ptr("Updated")}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestTaskHandler_UpdateTask_ValidationError checks that invalid update data returns a validation error
func TestTaskHandler_UpdateTask_ValidationError(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		GetTaskFunc: func(id string) (*model.Task, error) {
			return newTestTask(), nil
		},
		UpdateTaskFunc: func(task *model.Task) (*model.Task, error) {
			return nil, validation.NewValidationError("validation error")
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	reqBody := dto.UpdateTaskRequest{Title: ptr("bad")}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestTaskHandler_UpdateTask_NotFound checks that updating non-existent task returns not found error
func TestTaskHandler_UpdateTask_NotFound(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		GetTaskFunc: func(id string) (*model.Task, error) {
			return nil, repository.ErrTaskNotFound
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	reqBody := dto.UpdateTaskRequest{Title: ptr("Updated")}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestTaskHandler_DeleteTask_Success checks that a task is deleted successfully
func TestTaskHandler_DeleteTask_Success(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		DeleteTaskFunc: func(id string) error {
			return nil
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
}

// TestTaskHandler_DeleteTask_NotFound checks that deleting non-existent task returns not found error
func TestTaskHandler_DeleteTask_NotFound(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		DeleteTaskFunc: func(id string) error {
			return repository.ErrTaskNotFound
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestTaskHandler_UpdateTaskStatus_Success checks that task status is updated successfully
func TestTaskHandler_UpdateTaskStatus_Success(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		GetTaskFunc: func(id string) (*model.Task, error) {
			return newTestTask(), nil
		},
		SetTaskCompletionFunc: func(task *model.Task) (*model.Task, error) {
			task.IsCompleted = true
			task.Status = model.StatusCompleted
			return task, nil
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	reqBody := dto.UpdateTaskStatusRequest{IsCompleted: true}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestTaskHandler_UpdateTaskStatus_ValidationError checks that invalid status update returns a validation error
func TestTaskHandler_UpdateTaskStatus_ValidationError(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		GetTaskFunc: func(id string) (*model.Task, error) {
			return newTestTask(), nil
		},
		SetTaskCompletionFunc: func(task *model.Task) (*model.Task, error) {
			return nil, validation.NewValidationError("validation error")
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	reqBody := dto.UpdateTaskStatusRequest{IsCompleted: true}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestTaskHandler_UpdateTaskStatus_NotFound checks that updating status of non-existent task returns not found error
func TestTaskHandler_UpdateTaskStatus_NotFound(t *testing.T) {
	// Arrange
	mockUC := &mockTaskUsecase{
		GetTaskFunc: func(id string) (*model.Task, error) {
			return nil, repository.ErrTaskNotFound
		},
	}
	handler := NewTaskHandler(mockUC)
	router := setupRouter(handler)

	reqBody := dto.UpdateTaskStatusRequest{IsCompleted: true}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --- Utility ---

func ptr[T any](v T) *T {
	return &v
}
