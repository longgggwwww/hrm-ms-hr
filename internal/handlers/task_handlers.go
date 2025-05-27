package handlers

import (
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type TaskHandler struct {
	Client *ent.Client
}

func NewTaskHandler(client *ent.Client) *TaskHandler {
	return &TaskHandler{
		Client: client,
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
	}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Code        string  `json:"code" binding:"required"`
		Description *string `json:"description"`
		Process     *int    `json:"process"`
		Status      *string `json:"status"`
		StartAt     *string `json:"start_at"`
		ProjectID   *int    `json:"project_id"`
		Type        *string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from JWT token
	userID, err := utils.ExtractUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Parse start_at if provided
	var startAtPtr *time.Time
	if req.StartAt != nil {
		startAt, err := time.Parse(time.RFC3339, *req.StartAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid start_at format, must be RFC3339",
			})
			return
		}
		startAtPtr = &startAt
	}

	// Validate and set status
	var statusVal task.Status
	if req.Status != nil {
		switch *req.Status {
		case string(task.StatusNotReceived),
			string(task.StatusReceived),
			string(task.StatusInProgress),
			string(task.StatusCompleted),
			string(task.StatusCancelled):
			statusVal = task.Status(*req.Status)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value. Valid values: not_received, received, in_progress, completed, cancelled",
			})
			return
		}
	} else {
		statusVal = task.StatusNotReceived
	}

	// Validate and set type
	var typeVal task.Type
	if req.Type != nil {
		switch *req.Type {
		case string(task.TypeTask),
			string(task.TypeFeature),
			string(task.TypeBug),
			string(task.TypeAnother):
			typeVal = task.Type(*req.Type)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid type value. Valid values: task, feature, bug, another",
			})
			return
		}
	} else {
		typeVal = task.TypeTask
	}

	// Set default process if not provided
	if req.Process == nil {
		defaultProcess := 0
		req.Process = &defaultProcess
	}

	taskCreate := h.Client.Task.Create().
		SetName(req.Name).
		SetCode(req.Code).
		SetNillableDescription(req.Description).
		SetProcess(*req.Process).
		SetStatus(statusVal).
		SetNillableStartAt(startAtPtr).
		SetNillableProjectID(req.ProjectID).
		SetCreatorID(userID).
		SetUpdaterID(userID).
		SetType(typeVal)

	row, err := taskCreate.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Get the created task with all edges
	task, err := h.Client.Task.Query().
		Where(task.ID(row.ID)).
		WithProject().
		WithLabels().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusCreated, row)
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
// - order_by: Sort field (id, name, code, status, type, process, project_id, creator_id, start_at, created_at, updated_at)
// - order_dir: Sort direction (asc, desc) - default: desc
// - page: Page number (default: 1)
// - limit: Items per page (default: 10, max: 100)
//
// Example: GET /tasks?name=example&status=in_progress&type=feature&order_by=name&order_dir=asc&page=1&limit=20
func (h *TaskHandler) List(c *gin.Context) {
	query := h.Client.Task.Query().
		WithProject().
		WithLabels()

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
			"error": "Invalid order_by field. Valid fields: id, name, code, status, type, process, project_id, creator_id, start_at, created_at, updated_at",
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

	var req struct {
		Name        *string `json:"name"`
		Code        *string `json:"code"`
		Description *string `json:"description"`
		Process     *int    `json:"process"`
		Status      *string `json:"status"`
		StartAt     *string `json:"start_at"`
		ProjectID   *int    `json:"project_id"`
		CreatorID   *int    `json:"creator_id"`
		UpdaterID   *int    `json:"updater_id"`
		Type        *string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskUpdate := h.Client.Task.UpdateOneID(id)
	if req.Name != nil {
		taskUpdate.SetName(*req.Name)
	}
	if req.Code != nil {
		taskUpdate.SetCode(*req.Code)
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
	if req.ProjectID != nil {
		taskUpdate.SetProjectID(*req.ProjectID)
	}
	if req.CreatorID != nil {
		taskUpdate.SetCreatorID(*req.CreatorID)
	}
	if req.UpdaterID != nil {
		taskUpdate.SetUpdaterID(*req.UpdaterID)
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
	var req struct {
		IDs []int `json:"ids" binding:"required,min=1"`
	}
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

	response := gin.H{
		"deleted_count": deletedCount,
	}

	// Include failed IDs and errors if any
	if len(failedIDs) > 0 {
		response["failed_ids"] = failedIDs
		response["errors"] = errors
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
