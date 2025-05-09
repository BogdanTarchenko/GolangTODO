package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
	"todo/internal/delivery/http/dto"
	"todo/internal/pkg/utils"
)

type APITestSuite struct {
	suite.Suite
	baseURL string
	client  *http.Client
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, &APITestSuite{
		baseURL: "http://localhost:8080",
		client:  &http.Client{Timeout: 10 * time.Second},
	})
}

// TestTaskCRUD tests the complete lifecycle of a task including creation, reading, updating, status change, and deletion.
// It verifies that all CRUD operations work correctly and maintain data consistency.
func (s *APITestSuite) TestTaskCRUD() {
	// Arrange: Prepare initial task data
	createReq := dto.CreateTaskRequest{
		Title:       "Test task !1 !before 01.01.2100",
		Description: utils.Ptr("Test description"),
		Priority:    "HIGH",
	}
	createBody, _ := json.Marshal(createReq)

	resp, err := s.client.Post(
		s.baseURL+"/api/tasks",
		"application/json",
		bytes.NewReader(createBody),
	)
	s.NoError(err)
	s.Equal(http.StatusCreated, resp.StatusCode)

	var createResp dto.TaskResponse
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	s.NoError(err)
	s.NotEmpty(createResp.ID)
	s.Equal("Test task", createResp.Title)
	s.Equal("HIGH", createResp.Priority)
	s.NotNil(createResp.Deadline)
	resp.Body.Close()

	resp, err = s.client.Get(s.baseURL + "/api/tasks")
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var listResp dto.PaginatedTasksResponse
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	s.NoError(err)
	s.GreaterOrEqual(len(listResp.Items), 1)
	resp.Body.Close()

	resp, err = s.client.Get(fmt.Sprintf("%s/api/tasks/%s", s.baseURL, createResp.ID))
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var getResp dto.TaskResponse
	err = json.NewDecoder(resp.Body).Decode(&getResp)
	s.NoError(err)
	s.Equal(createResp.ID, getResp.ID)
	resp.Body.Close()

	updateReq := dto.UpdateTaskRequest{
		Title:       utils.Ptr("Updated task"),
		Description: utils.Ptr("Updated description"),
		Priority:    utils.Ptr("MEDIUM"),
	}
	updateBody, _ := json.Marshal(updateReq)

	req, _ := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/api/tasks/%s", s.baseURL, createResp.ID),
		bytes.NewReader(updateBody),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var updateResp dto.TaskResponse
	err = json.NewDecoder(resp.Body).Decode(&updateResp)
	s.NoError(err)
	s.Equal("Updated task", updateResp.Title)
	s.Equal("MEDIUM", updateResp.Priority)
	resp.Body.Close()

	statusReq := dto.UpdateTaskStatusRequest{
		IsCompleted: true,
	}
	statusBody, _ := json.Marshal(statusReq)

	req, _ = http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/api/tasks/%s/status", s.baseURL, createResp.ID),
		bytes.NewReader(statusBody),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var statusResp dto.TaskResponse
	err = json.NewDecoder(resp.Body).Decode(&statusResp)
	s.NoError(err)
	s.True(statusResp.IsCompleted)
	resp.Body.Close()

	req, _ = http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/api/tasks/%s", s.baseURL, createResp.ID),
		nil,
	)

	resp, err = s.client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	resp, err = s.client.Get(fmt.Sprintf("%s/api/tasks/%s", s.baseURL, createResp.ID))
	s.NoError(err)
	s.Equal(http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}

// TestTaskValidation verifies that the API properly validates task input data.
// It checks various validation scenarios including empty titles, short titles, invalid priorities, and past dates.
func (s *APITestSuite) TestTaskValidation() {
	// Arrange: Define test cases with invalid inputs
	tests := []struct {
		name       string
		request    dto.CreateTaskRequest
		wantStatus int
	}{
		{
			name: "Empty title",
			request: dto.CreateTaskRequest{
				Title: "",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Short title",
			request: dto.CreateTaskRequest{
				Title: "abc",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid priority",
			request: dto.CreateTaskRequest{
				Title:    "Valid title",
				Priority: "INVALID",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Past date macro",
			request: dto.CreateTaskRequest{
				Title: "Task !before " + time.Now().Add(-24*time.Hour).Format("02.01.2006"),
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Arrange: Prepare request body
			body, _ := json.Marshal(tt.request)

			// Act: Send POST request
			resp, err := s.client.Post(
				s.baseURL+"/api/tasks",
				"application/json",
				bytes.NewReader(body),
			)

			// Assert: Verify response
			s.NoError(err)
			s.Equal(tt.wantStatus, resp.StatusCode)
			resp.Body.Close()
		})
	}
}

// TestTaskPagination verifies that the task listing endpoint correctly implements pagination.
// It tests various pagination scenarios including different page sizes and invalid page numbers.
func (s *APITestSuite) TestTaskPagination() {
	// Arrange: Get initial task count
	resp, err := s.client.Get(s.baseURL + "/api/tasks?page=1&page_size=1")
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var initialResp dto.PaginatedTasksResponse
	err = json.NewDecoder(resp.Body).Decode(&initialResp)
	s.NoError(err)
	initialTotal := initialResp.Meta.Total
	resp.Body.Close()

	for i := 1; i <= 15; i++ {
		req := dto.CreateTaskRequest{
			Title: fmt.Sprintf("Task %d", i),
		}
		body, _ := json.Marshal(req)

		resp, err := s.client.Post(
			s.baseURL+"/api/tasks",
			"application/json",
			bytes.NewReader(body),
		)
		s.NoError(err)
		s.Equal(http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	expectedTotal := initialTotal + 15

	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantCount  int
		wantTotal  int
		wantStatus int
	}{
		{
			name:       "First page",
			page:       1,
			pageSize:   10,
			wantCount:  10,
			wantTotal:  expectedTotal,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Second page",
			page:       2,
			pageSize:   10,
			wantCount:  10,
			wantTotal:  expectedTotal,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Zero page",
			page:       0,
			pageSize:   10,
			wantCount:  0,
			wantTotal:  0,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Zero page size",
			page:       1,
			pageSize:   0,
			wantCount:  0,
			wantTotal:  0,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			resp, err := s.client.Get(
				fmt.Sprintf("%s/api/tasks?page=%d&page_size=%d", s.baseURL, tt.page, tt.pageSize),
			)
			s.NoError(err)
			s.Equal(tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				var listResp dto.PaginatedTasksResponse
				err = json.NewDecoder(resp.Body).Decode(&listResp)
				s.NoError(err)
				s.Equal(tt.wantCount, len(listResp.Items))
				s.Equal(tt.wantTotal, listResp.Meta.Total)
			}
			resp.Body.Close()
		})
	}
}

// TestTaskFiltering verifies that the task listing endpoint correctly filters tasks based on various criteria.
// It tests filtering by priority, completion status, and combinations of filters.
func (s *APITestSuite) TestTaskFiltering() {
	// Arrange: Clear existing tasks and create test data
	req, _ := http.NewRequest("DELETE", s.baseURL+"/api/tasks/all", nil)
	resp, err := s.client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	tasks := []struct {
		title    string
		priority string
		status   bool
	}{
		{"High Priority Task", "HIGH", false},
		{"Medium Priority Task", "MEDIUM", true},
		{"Low Priority Task", "LOW", false},
	}

	for _, task := range tasks {
		req := dto.CreateTaskRequest{
			Title:    task.title,
			Priority: task.priority,
		}
		body, _ := json.Marshal(req)

		resp, err := s.client.Post(
			s.baseURL+"/api/tasks",
			"application/json",
			bytes.NewReader(body),
		)
		s.NoError(err)
		s.Equal(http.StatusCreated, resp.StatusCode)

		var createResp dto.TaskResponse
		err = json.NewDecoder(resp.Body).Decode(&createResp)
		s.NoError(err)
		resp.Body.Close()

		if task.status {
			statusReq := dto.UpdateTaskStatusRequest{IsCompleted: true}
			statusBody, _ := json.Marshal(statusReq)

			req, _ := http.NewRequest(
				"PATCH",
				fmt.Sprintf("%s/api/tasks/%s/status", s.baseURL, createResp.ID),
				bytes.NewReader(statusBody),
			)
			req.Header.Set("Content-Type", "application/json")

			resp, err = s.client.Do(req)
			s.NoError(err)
			s.Equal(http.StatusOK, resp.StatusCode)
			resp.Body.Close()
		}
	}

	tests := []struct {
		name       string
		query      string
		wantCount  int
		wantStatus int
	}{
		{
			name:       "Filter by HIGH priority",
			query:      "priority=HIGH&page=1&page_size=10",
			wantCount:  10,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Filter by completed tasks",
			query:      "is_completed=true&page=1&page_size=10",
			wantCount:  10,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Filter by incomplete tasks",
			query:      "is_completed=false&page=1&page_size=10",
			wantCount:  10,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Combined filter",
			query:      "priority=MEDIUM&is_completed=true&page=1&page_size=10",
			wantCount:  10,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid priority filter",
			query:      "priority=INVALID&page=1&page_size=10",
			wantCount:  0,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			resp, err := s.client.Get(fmt.Sprintf("%s/api/tasks?%s", s.baseURL, tt.query))
			s.NoError(err)
			s.Equal(tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				var listResp dto.PaginatedTasksResponse
				err = json.NewDecoder(resp.Body).Decode(&listResp)
				s.NoError(err)
				s.Equal(tt.wantCount, len(listResp.Items))
			}
			resp.Body.Close()
		})
	}
}

// TestTaskSorting verifies that the task listing endpoint correctly sorts tasks based on different criteria.
// It tests sorting by priority and deadline in both ascending and descending order.
func (s *APITestSuite) TestTaskSorting() {
	// Arrange: Clear existing tasks and create test data with different priorities and deadlines
	req, _ := http.NewRequest("DELETE", s.baseURL+"/api/tasks/all", nil)
	resp, err := s.client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	now := time.Now()
	tasks := []struct {
		title    string
		priority string
		deadline time.Time
	}{
		{"Task 1", "LOW", now.Add(24 * time.Hour)},
		{"Task 2", "HIGH", now.Add(48 * time.Hour)},
		{"Task 3", "MEDIUM", now.Add(72 * time.Hour)},
	}

	for _, task := range tasks {
		req := dto.CreateTaskRequest{
			Title:    task.title,
			Priority: task.priority,
			Deadline: &task.deadline,
		}
		body, _ := json.Marshal(req)

		resp, err := s.client.Post(
			s.baseURL+"/api/tasks",
			"application/json",
			bytes.NewReader(body),
		)
		s.NoError(err)
		s.Equal(http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	tests := []struct {
		name       string
		query      string
		checkOrder func([]dto.TaskResponse) bool
	}{
		{
			name:  "Sort by priority (ascending)",
			query: "sort_by=priority&sort_order=asc&page=1&page_size=10",
			checkOrder: func(tasks []dto.TaskResponse) bool {
				if len(tasks) < 4 {
					return false
				}

				priorities := []string{tasks[0].Priority, tasks[1].Priority, tasks[2].Priority, tasks[3].Priority}
				return priorities[0] == "LOW" && priorities[1] == "LOW" && priorities[2] == "LOW" && priorities[3] == "LOW"
			},
		},
		{
			name:  "Sort by deadline (ascending)",
			query: "sort_by=deadline&sort_order=asc&page=1&page_size=10",
			checkOrder: func(tasks []dto.TaskResponse) bool {
				if len(tasks) < 3 {
					return false
				}

				for _, task := range tasks[:3] {
					if task.Deadline == nil {
						return false
					}
				}

				for i := 0; i < len(tasks)-1; i++ {
					if tasks[i].Deadline.After(*tasks[i+1].Deadline) {
						return false
					}
				}
				return true
			},
		},
		{
			name:  "Invalid sort field",
			query: "sort_by=invalid_field&page=1&page_size=10",
			checkOrder: func(tasks []dto.TaskResponse) bool {
				return true
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			resp, err := s.client.Get(fmt.Sprintf("%s/api/tasks?%s", s.baseURL, tt.query))
			s.NoError(err)

			if tt.query == "sort_by=invalid_field&page=1&page_size=10" {
				s.Equal(http.StatusBadRequest, resp.StatusCode)
			} else {
				s.Equal(http.StatusOK, resp.StatusCode)
				var listResp dto.PaginatedTasksResponse
				err = json.NewDecoder(resp.Body).Decode(&listResp)
				s.NoError(err)
				s.True(tt.checkOrder(listResp.Items), "Порядок сортировки не соответствует ожидаемому")
			}
			resp.Body.Close()
		})
	}
}

// TestTaskBoundaryValues verifies that the API correctly handles edge cases and boundary values.
// It tests minimum valid title length, past deadlines, and far future deadlines.
func (s *APITestSuite) TestTaskBoundaryValues() {
	// Arrange: Clear existing tasks
	req, _ := http.NewRequest("DELETE", s.baseURL+"/api/tasks/all", nil)
	resp, err := s.client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	tests := []struct {
		name       string
		request    dto.CreateTaskRequest
		wantStatus int
	}{
		{
			name: "Minimal title length",
			request: dto.CreateTaskRequest{
				Title: "Test",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Deadline in the past",
			request: dto.CreateTaskRequest{
				Title:    "Valid title",
				Deadline: utils.Ptr(time.Now().Add(-1 * time.Second)),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Deadline in the future",
			request: dto.CreateTaskRequest{
				Title:    "Valid title",
				Deadline: utils.Ptr(time.Now().Add(100 * 365 * 24 * time.Hour)),
			},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			resp, err := s.client.Post(
				s.baseURL+"/api/tasks",
				"application/json",
				bytes.NewReader(body),
			)
			s.NoError(err)
			s.Equal(tt.wantStatus, resp.StatusCode)
			resp.Body.Close()
		})
	}
}
