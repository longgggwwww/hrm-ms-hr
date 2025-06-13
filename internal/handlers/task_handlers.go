package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/kafka"
	taskService "github.com/longgggwwww/hrm-ms-hr/internal/services/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type TaskHandler struct {
	Client      *ent.Client
	TaskService *taskService.TaskService
}

func NewTaskHandler(client *ent.Client) *TaskHandler {
	return &TaskHandler{
		Client:      client,
		TaskService: taskService.NewTaskService(client),
	}
}

// NewTaskHandlerWithKafka creates a new task handler with Kafka support
func NewTaskHandlerWithKafka(client *ent.Client, kafkaClient *kafka.KafkaClient) *TaskHandler {
	taskSvc := taskService.NewTaskService(client)
	taskSvc.SetKafkaClient(kafkaClient)

	return &TaskHandler{
		Client:      client,
		TaskService: taskSvc,
	}
}

func (h *TaskHandler) RegisterRoutes(r *gin.Engine) {
	tasks := r.Group("tasks")
	{
		tasks.POST("", h.Create)
		tasks.GET("", h.List)
		tasks.GET(":id", h.Get)
		tasks.PATCH(":id", h.Update)
		tasks.DELETE(":id", h.Delete)
		tasks.DELETE("", h.BulkDelete)
		tasks.PATCH(":id/receive", h.ReceiveTask)
		tasks.PATCH(":id/progress", h.UpdateProgress)
	}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req dtos.TaskCreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID and employee ID from JWT token
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userID := ids["user_id"]
	employeeID := ids["employee_id"]

	// Call the task service to create the task
	task, err := h.TaskService.Create(c.Request.Context(), userID, employeeID, req)
	if err != nil {
		// Check if it's a ServiceError
		if serviceErr, ok := err.(*taskService.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{"error": serviceErr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, task)
}

// List retrieves tasks with filtering, sorting, and pagination support.
//
// Query Parameters:
// - name: Filter by task name (contains search)
// - code: Filter by task code (contains search)
// - status: Filter by status (not_received, received, in_progress, completed, cancelled)
// - type: Filter by type (task, feature, bug, another)
// - project_id: Filter by project ID
// - creator_id: Filter by creator ID
// - process: Filter by process percentage
// - start_date_from: Filter tasks that start from this date (RFC3339 format)
// - start_date_to: Filter tasks that start before this date (RFC3339 format)
// - due_date_from: Filter tasks with due date from this date (RFC3339 format)
// - due_date_to: Filter tasks with due date before this date (RFC3339 format)
// - order_by: Sort field (id, name, code, status, type, process, project_id, creator_id, start_at, due_date, created_at, updated_at)
// - order_dir: Sort direction (asc, desc) - default: desc
// - page: Page number (default: 1)
// - limit: Items per page (default: 10, max: 100)
//
// Example: GET /tasks?name=example&status=in_progress&type=feature&order_by=name&order_dir=asc&page=1&limit=20
func (h *TaskHandler) List(c *gin.Context) {
	// Parse query parameters
	var query dtos.TaskListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values for pagination
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}
	if query.OrderBy == "" {
		query.OrderBy = "created_at"
	}
	if query.OrderDir == "" {
		query.OrderDir = "desc"
	}

	// Call service
	response, err := h.TaskService.List(c.Request.Context(), query)
	if err != nil {
		if serviceErr, ok := err.(*taskService.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{"error": serviceErr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Call service
	task, err := h.TaskService.Get(c.Request.Context(), id)
	if err != nil {
		if serviceErr, ok := err.(*taskService.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{"error": serviceErr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req dtos.TaskUpdateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from JWT token
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userID := ids["user_id"]

	// Call service
	updatedTask, err := h.TaskService.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		if serviceErr, ok := err.(*taskService.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{"error": serviceErr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Call service
	err = h.TaskService.Delete(c.Request.Context(), id)
	if err != nil {
		if serviceErr, ok := err.(*taskService.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{"error": serviceErr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// BulkDelete deletes multiple tasks by their IDs.
//
// Request body should contain:
//
//	{
//	  "ids": [1, 2, 3, 4, 5]
//	}
//
// Response will include:
// - deleted_count: number of successfully deleted tasks
// - failed_ids: array of IDs that failed to delete (if any)
// - errors: array of error messages for failed deletions (if any)
func (h *TaskHandler) BulkDelete(c *gin.Context) {
	var req dtos.TaskBulkDeleteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	response, err := h.TaskService.BulkDelete(c.Request.Context(), req)
	if err != nil {
		if serviceErr, ok := err.(*taskService.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{"error": serviceErr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Determine appropriate status code
	if response.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, response)
	} else if len(response.FailedIDs) > 0 {
		c.JSON(http.StatusPartialContent, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// ReceiveTask allows an assigned employee to receive/accept a task.
// Only employees who are assigned to the task can receive it.
// This changes the task status from "not_received" to "received".
func (h *TaskHandler) ReceiveTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Extract user ID from JWT token
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userID := ids["user_id"]

	// Call service
	taskWithEdges, err := h.TaskService.ReceiveTask(c.Request.Context(), id, userID)
	if err != nil {
		if serviceErr, ok := err.(*taskService.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{"error": serviceErr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task received successfully",
		"task":    taskWithEdges,
	})
}

// UpdateProgress allows an assigned employee to update task status and progress.
// Only employees who are assigned to the task can update its progress.
//
// Request body:
//
//	{
//	  "status": "in_progress|completed|cancelled", // optional
//	  "process": 50 // optional, percentage (0-100)
//	}
func (h *TaskHandler) UpdateProgress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Extract user ID from JWT token
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userID := ids["user_id"]

	var req dtos.TaskUpdateProgressInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	taskWithEdges, err := h.TaskService.UpdateProgress(c.Request.Context(), id, userID, req)
	if err != nil {
		if serviceErr, ok := err.(*taskService.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{"error": serviceErr.Msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task progress updated successfully",
		"task":    taskWithEdges,
	})
}
