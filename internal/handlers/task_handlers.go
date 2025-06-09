package handlers

import (
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
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
	query := h.Client.Task.Query().
		WithProject().
		WithLabels().
		WithAssignees()

	// Filter by name
	if name := c.Query("name"); name != "" {
		query = query.Where(task.NameContains(name))
	}

	// Filter by code
	if code := c.Query("code"); code != "" {
		query = query.Where(task.CodeContains(code))
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		switch status {
		case string(task.StatusNotReceived),
			string(task.StatusReceived),
			string(task.StatusInProgress),
			string(task.StatusCompleted),
			string(task.StatusCancelled):
			query = query.Where(task.StatusEQ(task.Status(status)))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value. Valid values: not_received, received, in_progress, completed, cancelled",
			})
			return
		}
	}

	// Filter by type
	if taskType := c.Query("type"); taskType != "" {
		switch taskType {
		case string(task.TypeTask),
			string(task.TypeFeature),
			string(task.TypeBug),
			string(task.TypeAnother):
			query = query.Where(task.TypeEQ(task.Type(taskType)))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid type value. Valid values: task, feature, bug, another",
			})
			return
		}
	}

	// Filter by project_id
	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid project_id format",
			})
			return
		}
		query = query.Where(task.ProjectIDEQ(projectID))
	}

	// Filter by creator_id
	if creatorIDStr := c.Query("creator_id"); creatorIDStr != "" {
		creatorID, err := strconv.Atoi(creatorIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid creator_id format",
			})
			return
		}
		query = query.Where(task.CreatorIDEQ(creatorID))
	}

	// Filter by process (percentage)
	if processStr := c.Query("process"); processStr != "" {
		process, err := strconv.Atoi(processStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid process format",
			})
			return
		}
		query = query.Where(task.ProcessEQ(process))
	}

	// Date range filtering
	if startDateStr := c.Query("start_date_from"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid start_date_from format, must be RFC3339",
			})
			return
		}
		query = query.Where(task.StartAtGTE(startDate))
	}

	if startDateStr := c.Query("start_date_to"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid start_date_to format, must be RFC3339",
			})
			return
		}
		query = query.Where(task.StartAtLTE(startDate))
	}

	// Due date range filtering
	if dueDateStr := c.Query("due_date_from"); dueDateStr != "" {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid due_date_from format, must be RFC3339",
			})
			return
		}
		query = query.Where(task.DueDateGTE(dueDate))
	}

	if dueDateStr := c.Query("due_date_to"); dueDateStr != "" {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid due_date_to format, must be RFC3339",
			})
			return
		}
		query = query.Where(task.DueDateLTE(dueDate))
	}

	// Order by field and direction
	orderBy := c.DefaultQuery("order_by", "created_at")
	orderDir := c.DefaultQuery("order_dir", "desc")

	var orderOption task.OrderOption
	switch orderBy {
	case "id":
		if orderDir == "asc" {
			orderOption = task.ByID()
		} else {
			orderOption = task.ByID(sql.OrderDesc())
		}
	case "name":
		if orderDir == "asc" {
			orderOption = task.ByName()
		} else {
			orderOption = task.ByName(sql.OrderDesc())
		}
	case "code":
		if orderDir == "asc" {
			orderOption = task.ByCode()
		} else {
			orderOption = task.ByCode(sql.OrderDesc())
		}
	case "status":
		if orderDir == "asc" {
			orderOption = task.ByStatus()
		} else {
			orderOption = task.ByStatus(sql.OrderDesc())
		}
	case "type":
		if orderDir == "asc" {
			orderOption = task.ByType()
		} else {
			orderOption = task.ByType(sql.OrderDesc())
		}
	case "process":
		if orderDir == "asc" {
			orderOption = task.ByProcess()
		} else {
			orderOption = task.ByProcess(sql.OrderDesc())
		}
	case "project_id":
		if orderDir == "asc" {
			orderOption = task.ByProjectID()
		} else {
			orderOption = task.ByProjectID(sql.OrderDesc())
		}
	case "creator_id":
		if orderDir == "asc" {
			orderOption = task.ByCreatorID()
		} else {
			orderOption = task.ByCreatorID(sql.OrderDesc())
		}
	case "start_at":
		if orderDir == "asc" {
			orderOption = task.ByStartAt()
		} else {
			orderOption = task.ByStartAt(sql.OrderDesc())
		}
	case "due_date":
		if orderDir == "asc" {
			orderOption = task.ByDueDate()
		} else {
			orderOption = task.ByDueDate(sql.OrderDesc())
		}
	case "created_at":
		if orderDir == "asc" {
			orderOption = task.ByCreatedAt()
		} else {
			orderOption = task.ByCreatedAt(sql.OrderDesc())
		}
	case "updated_at":
		if orderDir == "asc" {
			orderOption = task.ByUpdatedAt()
		} else {
			orderOption = task.ByUpdatedAt(sql.OrderDesc())
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order_by field. Valid fields: id, name, code, status, type, process, project_id, creator_id, start_at, due_date, created_at, updated_at",
		})
		return
	}

	// Apply ordering
	query = query.Order(orderOption)

	// Pagination
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	tasks, err := query.All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get total count for pagination info
	total, err := h.Client.Task.Query().Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := gin.H{
		"data": tasks,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.Client.Task.Query().
		Where(task.ID(id)).
		WithProject().
		WithLabels().
		WithAssignees().
		WithReports().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
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

	taskUpdate := h.Client.Task.UpdateOneID(id).SetUpdaterID(userID)

	if req.Name != nil {
		taskUpdate.SetName(*req.Name)
	}

	// Handle code update with auto-generation if needed
	if req.Code != nil {
		var taskCode string
		if *req.Code != "" {
			taskCode = *req.Code

			// Check if the new code already exists (excluding current task)
			var existsQuery *ent.TaskQuery
			if req.ProjectID != nil {
				// Check uniqueness within the project
				existsQuery = h.Client.Task.Query().
					Where(task.CodeEQ(taskCode)).
					Where(task.ProjectIDEQ(*req.ProjectID)).
					Where(task.IDNEQ(id))
			} else {
				// Check global uniqueness if no project
				existsQuery = h.Client.Task.Query().
					Where(task.CodeEQ(taskCode)).
					Where(task.IDNEQ(id))
			}

			exists, err := existsQuery.Exist(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate task code"})
				return
			}
			if exists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Task code already exists"})
				return
			}
		} else {
			// Auto-generate code in GitHub-style format: #x
			// Get the total count of tasks to generate the next sequence number
			count, err := h.Client.Task.Query().
				Where(task.IDNEQ(id)). // Exclude current task
				Count(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate task code"})
				return
			}

			sequence := count + 1
			taskCode = "#" + strconv.Itoa(sequence)

			// Double-check uniqueness (in case of concurrent requests)
			for {
				exists, err := h.Client.Task.Query().
					Where(task.CodeEQ(taskCode)).
					Where(task.IDNEQ(id)).
					Exist(c.Request.Context())
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate generated task code"})
					return
				}
				if !exists {
					break
				}
				sequence++
				taskCode = "#" + strconv.Itoa(sequence)
			}
		}
		taskUpdate.SetCode(taskCode)
	}
	if req.Description != nil {
		taskUpdate.SetDescription(*req.Description)
	}
	if req.Process != nil {
		taskUpdate.SetProcess(*req.Process)
	}
	if req.StartAt != nil {
		startAt, err := time.Parse(time.RFC3339, *req.StartAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_at format, must be RFC3339"})
			return
		}
		taskUpdate.SetStartAt(startAt)
	}
	if req.DueDate != nil {
		dueDate, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format, must be RFC3339"})
			return
		}
		taskUpdate.SetDueDate(dueDate)
	}
	if req.ProjectID != nil {
		taskUpdate.SetProjectID(*req.ProjectID)
	}
	if req.Status != nil {
		switch *req.Status {
		case string(task.StatusNotReceived),
			string(task.StatusReceived),
			string(task.StatusInProgress),
			string(task.StatusCompleted),
			string(task.StatusCancelled):
			taskUpdate.SetStatus(task.Status(*req.Status))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value. Valid values: not_received, received, in_progress, completed, cancelled",
			})
			return
		}
	}
	if req.Type != nil {
		switch *req.Type {
		case string(task.TypeTask),
			string(task.TypeFeature),
			string(task.TypeBug),
			string(task.TypeAnother):
			taskUpdate.SetType(task.Type(*req.Type))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid type value. Valid values: task, feature, bug, another",
			})
			return
		}
	}

	// Handle label assignments
	if req.LabelIDs != nil {
		// Clear existing labels and add new ones
		taskUpdate.ClearLabels()
		if len(req.LabelIDs) > 0 {
			taskUpdate.AddLabelIDs(req.LabelIDs...)
		}
	}

	// Handle assignee assignments
	if req.AssigneeIDs != nil {
		// Clear existing assignees first
		taskUpdate.ClearAssignees()

		if len(req.AssigneeIDs) > 0 {
			// Validate that all assignee IDs exist in the employee table
			existingEmployees, err := h.Client.Employee.Query().
				Where(employee.IDIn(req.AssigneeIDs...)).
				Select(employee.FieldID).
				All(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to validate assignee IDs"})
				return
			}

			// Create map of existing employee IDs for validation
			existingIDs := make(map[int]bool)
			for _, emp := range existingEmployees {
				existingIDs[emp.ID] = true
			}

			// Check if all requested assignee IDs exist
			var invalidIDs []int
			for _, id := range req.AssigneeIDs {
				if !existingIDs[id] {
					invalidIDs = append(invalidIDs, id)
				}
			}

			if len(invalidIDs) > 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":       "Some assignee IDs do not exist",
					"invalid_ids": invalidIDs,
				})
				return
			}

			taskUpdate.AddAssigneeIDs(req.AssigneeIDs...)
		}
	}

	_, err = taskUpdate.Save(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Get the updated task with all edges
	task, err := h.Client.Task.Query().
		Where(task.ID(id)).
		WithProject().
		WithLabels().
		WithAssignees().
		WithReports().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"id": id})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	_, err = h.Client.Task.Delete().Where(task.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
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

	// Validate maximum number of IDs to prevent abuse
	if len(req.IDs) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 100 IDs allowed per bulk delete operation"})
		return
	}

	// Check which tasks exist before attempting deletion
	existingTasks, err := h.Client.Task.Query().
		Where(task.IDIn(req.IDs...)).
		Select(task.FieldID).
		All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Create a map of existing task IDs for quick lookup
	existingIDs := make(map[int]bool)
	for _, t := range existingTasks {
		existingIDs[t.ID] = true
	}

	// Separate existing and non-existing IDs
	var validIDs []int
	var notFoundIDs []int
	for _, id := range req.IDs {
		if existingIDs[id] {
			validIDs = append(validIDs, id)
		} else {
			notFoundIDs = append(notFoundIDs, id)
		}
	}

	// Perform bulk deletion for valid IDs
	var deletedCount int
	var failedIDs []int
	var errors []string

	if len(validIDs) > 0 {
		deletedCount, err = h.Client.Task.Delete().
			Where(task.IDIn(validIDs...)).
			Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}

	// Add not found IDs to failed IDs
	failedIDs = append(failedIDs, notFoundIDs...)
	for _, id := range notFoundIDs {
		errors = append(errors, "Task ID "+strconv.Itoa(id)+" not found")
	}

	response := dtos.TaskBulkDeleteResponse{
		DeletedCount: deletedCount,
	}

	// Include failed IDs and errors if any
	if len(failedIDs) > 0 {
		response.FailedIDs = failedIDs
		response.Errors = errors
	}

	// Determine appropriate status code
	if deletedCount == 0 {
		c.JSON(http.StatusNotFound, response)
	} else if len(failedIDs) > 0 {
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

	// Get the task with assignees to check if user is assigned
	taskEntity, err := h.Client.Task.Query().
		Where(task.ID(id)).
		WithAssignees().
		Only(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Check if the current user is assigned to this task
	isAssigned := false
	for _, assignee := range taskEntity.Edges.Assignees {
		if assignee.ID == userID {
			isAssigned = true
			break
		}
	}

	if !isAssigned {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not assigned to this task"})
		return
	}

	// Check if task is in the correct status to be received
	if taskEntity.Status != task.StatusNotReceived {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Task can only be received when status is 'not_received'",
			"current_status": string(taskEntity.Status),
		})
		return
	}

	// Update task status to "received"
	updatedTask, err := h.Client.Task.UpdateOneID(id).
		SetStatus(task.StatusReceived).
		SetUpdaterID(userID).
		Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Get the updated task with all edges
	taskWithEdges, err := h.Client.Task.Query().
		Where(task.IDEQ(updatedTask.ID)).
		WithProject().
		WithLabels().
		WithAssignees().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, updatedTask)
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

	// Validate at least one field is provided
	if req.Status == nil && req.Process == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one field (status or process) must be provided",
		})
		return
	}

	// Get the task with assignees to check if user is assigned
	taskEntity, err := h.Client.Task.Query().
		Where(task.ID(id)).
		WithAssignees().
		Only(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Check if the current user is assigned to this task
	isAssigned := false
	for _, assignee := range taskEntity.Edges.Assignees {
		if assignee.ID == userID {
			isAssigned = true
			break
		}
	}

	if !isAssigned {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not assigned to this task"})
		return
	}

	// Check if task is in a valid status for progress updates
	if taskEntity.Status == task.StatusNotReceived {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Task must be received before updating progress",
			"current_status": string(taskEntity.Status),
		})
		return
	}

	if taskEntity.Status == task.StatusCompleted || taskEntity.Status == task.StatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Cannot update progress for completed or cancelled tasks",
			"current_status": string(taskEntity.Status),
		})
		return
	}

	taskUpdate := h.Client.Task.UpdateOneID(id).SetUpdaterID(userID)

	// Validate and set status if provided
	if req.Status != nil {
		switch *req.Status {
		case string(task.StatusInProgress),
			string(task.StatusCompleted),
			string(task.StatusCancelled):
			taskUpdate.SetStatus(task.Status(*req.Status))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value. Valid values for progress update: in_progress, completed, cancelled",
			})
			return
		}
	}

	// Validate and set process if provided
	if req.Process != nil {
		if *req.Process < 0 || *req.Process > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Process must be between 0 and 100",
			})
			return
		}
		taskUpdate.SetProcess(*req.Process)

		// Auto-update status based on process
		if *req.Process == 100 && req.Status == nil {
			taskUpdate.SetStatus(task.StatusCompleted)
		} else if *req.Process > 0 && *req.Process < 100 && req.Status == nil && taskEntity.Status == task.StatusReceived {
			taskUpdate.SetStatus(task.StatusInProgress)
		}
	}

	updatedTask, err := taskUpdate.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Get the updated task with all edges
	taskWithEdges, err := h.Client.Task.Query().
		Where(task.IDEQ(updatedTask.ID)).
		WithProject().
		WithLabels().
		WithAssignees().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, updatedTask)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task progress updated successfully",
		"task":    taskWithEdges,
	})
}
