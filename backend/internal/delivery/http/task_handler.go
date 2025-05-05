package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo/internal/delivery/http/dto"
	"todo/internal/domain/model"
	"todo/internal/domain/usecase"
)

type TaskHandler struct {
	usecase usecase.TaskUsecase
}

func NewTaskHandler(u usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{usecase: u}
}

func (h *TaskHandler) RegisterRoutes(r *gin.Engine) {
	tasks := r.Group("/api/tasks")
	{
		tasks.POST("", h.CreateTask)
		tasks.GET("", h.ListTasks)
		tasks.GET("/:id", h.GetTask)
		tasks.PATCH("/:id", h.UpdateTask)
		tasks.PATCH("/:id/status", h.UpdateTaskStatus)
		tasks.DELETE("/:id", h.DeleteTask)
	}
}

// CreateTask godoc
// @Summary     Create a new task
// @Description Creates a new task with the provided data
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       task  body      dto.CreateTaskRequest  true  "New task data"
// @Success     201   {object}  dto.TaskResponse
// @Failure     400   {object}  map[string]string   // Invalid input
// @Failure     500   {object}  map[string]string   // Internal server error
// @Router      /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Priority:    model.TaskPriority(req.Priority),
	}

	createdTask, err := h.usecase.CreateTask(task)
	if err != nil {
		c.Error(err)
		return
	}

	resp := dto.TaskResponse{
		ID:          createdTask.ID,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		Deadline:    createdTask.Deadline,
		Status:      string(createdTask.Status),
		Priority:    string(createdTask.Priority),
		CreatedAt:   createdTask.CreatedAt,
		UpdatedAt:   createdTask.UpdatedAt,
	}

	c.JSON(http.StatusCreated, resp)
}

// ListTasks godoc
// @Summary     List all tasks
// @Description Returns a list of all existing tasks with optional filters and sorting
// @Tags        tasks
// @Produce     json
// @Param       status     query     string  false  "Task status"
// @Param       priority   query     string  false  "Task priority"
// @Param       sort_by    query     string  false  "Sort by field: deadline, created_at, priority"
// @Param       sort_order query     string  false  "Sort order: asc or desc"
// @Param       page       query     int     false  "Page number"
// @Param       page_size  query     int     false  "Page size"
// @Success     200  {object}  dto.PaginatedTasksResponse
// @Failure     500  {object}  map[string]string   // Internal server error
// @Router      /tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	var query dto.ListTasksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	if c.Query("page") == "" {
		query.Page = 1
	}
	if c.Query("page_size") == "" {
		query.PageSize = 10
	}

	filter := &model.TaskFilter{
		Status:    query.Status,
		Priority:  query.Priority,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
		Page:      query.Page,
		PageSize:  query.PageSize,
	}

	tasks, total, err := h.usecase.ListTasksWithFilter(filter)
	if err != nil {
		c.Error(err)
		return
	}

	var respItems []dto.TaskResponse
	for _, t := range tasks {
		respItems = append(respItems, dto.TaskResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Deadline:    t.Deadline,
			Status:      string(t.Status),
			Priority:    string(t.Priority),
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	totalPages := (total + filter.PageSize - 1) / filter.PageSize

	resp := dto.PaginatedTasksResponse{
		Items: respItems,
		Meta: dto.PaginationMeta{
			Total:      total,
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// GetTask godoc
// @Summary     Get a task by ID
// @Description Returns a task by its identifier
// @Tags        tasks
// @Produce     json
// @Param       id   path      string  true  "Task ID"
// @Success     200  {object}  dto.TaskResponse
// @Failure     404  {object}  map[string]string   // Task not found
// @Failure     500  {object}  map[string]string   // Internal server error
// @Router      /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := h.usecase.GetTask(id)
	if err != nil {
		c.Error(err)
		return
	}

	resp := dto.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Deadline:    task.Deadline,
		Status:      string(task.Status),
		Priority:    string(task.Priority),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateTask godoc
// @Summary     Update a task
// @Description Updates an existing task by ID
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id    path      string                 true  "Task ID"
// @Param       task  body      dto.UpdateTaskRequest  true  "Updated task data"
// @Success     200   {object}  dto.TaskResponse
// @Failure     400   {object}  map[string]string   // Invalid input
// @Failure     404   {object}  map[string]string   // Task not found
// @Failure     500   {object}  map[string]string   // Internal server error
// @Router      /tasks/{id} [patch]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	existing, err := h.usecase.GetTask(id)
	if err != nil {
		c.Error(err)
		return
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.Deadline != nil {
		existing.Deadline = req.Deadline
	}
	if req.Priority != nil {
		existing.Priority = model.TaskPriority(*req.Priority)
	}

	updatedTask, err := h.usecase.UpdateTask(existing)
	if err != nil {
		c.Error(err)
		return
	}

	resp := dto.TaskResponse{
		ID:          updatedTask.ID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Deadline:    updatedTask.Deadline,
		Status:      string(updatedTask.Status),
		Priority:    string(updatedTask.Priority),
		CreatedAt:   updatedTask.CreatedAt,
		UpdatedAt:   updatedTask.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteTask godoc
// @Summary     Delete a task
// @Description Deletes a task by its identifier
// @Tags        tasks
// @Produce     json
// @Param       id   path      string  true  "Task ID"
// @Success     204  "Task successfully deleted"
// @Failure     404  {object}  map[string]string   // Task not found
// @Failure     500  {object}  map[string]string   // Internal server error
// @Router      /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	if err := h.usecase.DeleteTask(id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateTaskStatus godoc
// @Summary     Mark task as completed or not completed
// @Description Updates the is_completed flag and recalculates the status
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id   path      string                     true  "Task ID"
// @Param       body body      dto.UpdateTaskStatusRequest true  "Completion status"
// @Success     200  {object}  dto.TaskResponse
// @Failure     400  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /tasks/{id}/status [patch]
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	task, err := h.usecase.GetTask(id)
	if err != nil {
		c.Error(err)
		return
	}

	task.IsCompleted = req.IsCompleted

	updatedTask, err := h.usecase.SetTaskCompletion(task)
	if err != nil {
		c.Error(err)
		return
	}

	resp := dto.TaskResponse{
		ID:          updatedTask.ID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Deadline:    updatedTask.Deadline,
		Status:      string(updatedTask.Status),
		Priority:    string(updatedTask.Priority),
		CreatedAt:   updatedTask.CreatedAt,
		UpdatedAt:   updatedTask.UpdatedAt,
	}
	c.JSON(http.StatusOK, resp)
}
