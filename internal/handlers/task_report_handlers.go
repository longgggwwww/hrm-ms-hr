package handlers

import (
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/ent/taskreport"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type TaskReportHandler struct {
	Client *ent.Client
}

func NewTaskReportHandler(client *ent.Client) *TaskReportHandler {
	return &TaskReportHandler{
		Client: client,
	}
}

func (h *TaskReportHandler) RegisterRoutes(r *gin.Engine) {
	reports := r.Group("task-reports")
	{
		reports.POST("", h.Create)
		reports.GET("", h.List)
		reports.GET(":id", h.Get)
		reports.PATCH(":id", h.Update)
		reports.DELETE(":id", h.Delete)
		reports.DELETE("", h.BulkDelete)
	}
}

func (h *TaskReportHandler) Create(c *gin.Context) {
	var req struct {
		Title               string  `json:"title" binding:"required"`
		Content             string  `json:"content" binding:"required"`
		Status              *string `json:"status"`
		ProgressPercentage  *int    `json:"progress_percentage"`
		ReportedAt          *string `json:"reported_at"`
		IssuesEncountered   *string `json:"issues_encountered"`
		NextSteps           *string `json:"next_steps"`
		EstimatedCompletion *string `json:"estimated_completion"`
		TaskID              int     `json:"task_id" binding:"required"`
	}
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

	// Validate that the task exists
	taskExists, err := h.Client.Task.Query().Where(task.ID(req.TaskID)).Exist(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to validate task ID"})
		return
	}
	if !taskExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID does not exist"})
		return
	}

	// Validate and set status
	var statusVal taskreport.Status
	if req.Status != nil {
		switch *req.Status {
		case string(taskreport.StatusReceived),
			string(taskreport.StatusNotReceived),
			string(taskreport.StatusInProgress),
			string(taskreport.StatusCompleted),
			string(taskreport.StatusCancelled):
			statusVal = taskreport.Status(*req.Status)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value. Valid values: received, not_received, in_progress, completed, cancelled",
			})
			return
		}
	} else {
		statusVal = taskreport.StatusReceived // Default value
	}

	// Parse reported_at if provided
	var reportedAtPtr *time.Time
	if req.ReportedAt != nil {
		reportedAt, err := time.Parse(time.RFC3339, *req.ReportedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid reported_at format, must be RFC3339",
			})
			return
		}
		reportedAtPtr = &reportedAt
	}

	// Parse estimated_completion if provided
	var estimatedCompletionPtr *time.Time
	if req.EstimatedCompletion != nil {
		estimatedCompletion, err := time.Parse(time.RFC3339, *req.EstimatedCompletion)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid estimated_completion format, must be RFC3339",
			})
			return
		}
		estimatedCompletionPtr = &estimatedCompletion
	}

	// Validate progress percentage
	if req.ProgressPercentage != nil && (*req.ProgressPercentage < 0 || *req.ProgressPercentage > 100) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Progress percentage must be between 0 and 100",
		})
		return
	}

	// Create task report
	taskReportCreate := h.Client.TaskReport.Create().
		SetTitle(req.Title).
		SetContent(req.Content).
		SetStatus(statusVal).
		SetTaskID(req.TaskID).
		SetReporterID(userID)

	if req.ProgressPercentage != nil {
		taskReportCreate.SetProgressPercentage(*req.ProgressPercentage)
	}
	if reportedAtPtr != nil {
		taskReportCreate.SetReportedAt(*reportedAtPtr)
	}
	if req.IssuesEncountered != nil {
		taskReportCreate.SetIssuesEncountered(*req.IssuesEncountered)
	}
	if req.NextSteps != nil {
		taskReportCreate.SetNextSteps(*req.NextSteps)
	}
	if estimatedCompletionPtr != nil {
		taskReportCreate.SetEstimatedCompletion(*estimatedCompletionPtr)
	}

	row, err := taskReportCreate.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Get the created task report with all edges
	taskReport, err := h.Client.TaskReport.Query().
		Where(taskreport.ID(row.ID)).
		WithTask().
		WithReporter().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusCreated, row)
		return
	}
	c.JSON(http.StatusCreated, taskReport)
}

// List retrieves task reports with filtering, sorting, and pagination support.
//
// Query Parameters:
// - title: Filter by title (contains search)
// - content: Filter by content (contains search)
// - status: Filter by status (exact match)
// - task_id: Filter by task ID
// - reporter_id: Filter by reporter ID
// - progress_percentage: Filter by exact progress percentage
// - reported_at_from: Filter reports from this date (RFC3339)
// - reported_at_to: Filter reports to this date (RFC3339)
// - estimated_completion_from: Filter by estimated completion from date (RFC3339)
// - estimated_completion_to: Filter by estimated completion to date (RFC3339)
// - order_by: Field to order by (id, title, status, progress_percentage, reported_at, estimated_completion, created_at, updated_at)
// - order_dir: Order direction (asc, desc)
// - page: Page number (default: 1)
// - limit: Items per page (default: 10, max: 100)
func (h *TaskReportHandler) List(c *gin.Context) {
	query := h.Client.TaskReport.Query().
		WithTask().
		WithReporter()

	// Filter by title
	if title := c.Query("title"); title != "" {
		query = query.Where(taskreport.TitleContains(title))
	}

	// Filter by content
	if content := c.Query("content"); content != "" {
		query = query.Where(taskreport.ContentContains(content))
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		switch status {
		case string(taskreport.StatusReceived),
			string(taskreport.StatusNotReceived),
			string(taskreport.StatusInProgress),
			string(taskreport.StatusCompleted),
			string(taskreport.StatusCancelled):
			query = query.Where(taskreport.StatusEQ(taskreport.Status(status)))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value. Valid values: received, not_received, in_progress, completed, cancelled",
			})
			return
		}
	}

	// Filter by task_id
	if taskIDStr := c.Query("task_id"); taskIDStr != "" {
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid task_id format",
			})
			return
		}
		query = query.Where(taskreport.TaskIDEQ(taskID))
	}

	// Filter by reporter_id
	if reporterIDStr := c.Query("reporter_id"); reporterIDStr != "" {
		reporterID, err := strconv.Atoi(reporterIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid reporter_id format",
			})
			return
		}
		query = query.Where(taskreport.ReporterIDEQ(reporterID))
	}

	// Filter by progress_percentage
	if progressStr := c.Query("progress_percentage"); progressStr != "" {
		progress, err := strconv.Atoi(progressStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid progress_percentage format",
			})
			return
		}
		query = query.Where(taskreport.ProgressPercentageEQ(progress))
	}

	// Date range filtering for reported_at
	if reportedAtFromStr := c.Query("reported_at_from"); reportedAtFromStr != "" {
		reportedAtFrom, err := time.Parse(time.RFC3339, reportedAtFromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid reported_at_from format, must be RFC3339",
			})
			return
		}
		query = query.Where(taskreport.ReportedAtGTE(reportedAtFrom))
	}

	if reportedAtToStr := c.Query("reported_at_to"); reportedAtToStr != "" {
		reportedAtTo, err := time.Parse(time.RFC3339, reportedAtToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid reported_at_to format, must be RFC3339",
			})
			return
		}
		query = query.Where(taskreport.ReportedAtLTE(reportedAtTo))
	}

	// Date range filtering for estimated_completion
	if estimatedFromStr := c.Query("estimated_completion_from"); estimatedFromStr != "" {
		estimatedFrom, err := time.Parse(time.RFC3339, estimatedFromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid estimated_completion_from format, must be RFC3339",
			})
			return
		}
		query = query.Where(taskreport.EstimatedCompletionGTE(estimatedFrom))
	}

	if estimatedToStr := c.Query("estimated_completion_to"); estimatedToStr != "" {
		estimatedTo, err := time.Parse(time.RFC3339, estimatedToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid estimated_completion_to format, must be RFC3339",
			})
			return
		}
		query = query.Where(taskreport.EstimatedCompletionLTE(estimatedTo))
	}

	// Order by field and direction
	orderBy := c.DefaultQuery("order_by", "created_at")
	orderDir := c.DefaultQuery("order_dir", "desc")

	var orderOption taskreport.OrderOption
	switch orderBy {
	case "id":
		if orderDir == "asc" {
			orderOption = taskreport.ByID()
		} else {
			orderOption = taskreport.ByID(sql.OrderDesc())
		}
	case "title":
		if orderDir == "asc" {
			orderOption = taskreport.ByTitle()
		} else {
			orderOption = taskreport.ByTitle(sql.OrderDesc())
		}
	case "status":
		if orderDir == "asc" {
			orderOption = taskreport.ByStatus()
		} else {
			orderOption = taskreport.ByStatus(sql.OrderDesc())
		}
	case "progress_percentage":
		if orderDir == "asc" {
			orderOption = taskreport.ByProgressPercentage()
		} else {
			orderOption = taskreport.ByProgressPercentage(sql.OrderDesc())
		}
	case "reported_at":
		if orderDir == "asc" {
			orderOption = taskreport.ByReportedAt()
		} else {
			orderOption = taskreport.ByReportedAt(sql.OrderDesc())
		}
	case "estimated_completion":
		if orderDir == "asc" {
			orderOption = taskreport.ByEstimatedCompletion()
		} else {
			orderOption = taskreport.ByEstimatedCompletion(sql.OrderDesc())
		}
	case "task_id":
		if orderDir == "asc" {
			orderOption = taskreport.ByTaskID()
		} else {
			orderOption = taskreport.ByTaskID(sql.OrderDesc())
		}
	case "reporter_id":
		if orderDir == "asc" {
			orderOption = taskreport.ByReporterID()
		} else {
			orderOption = taskreport.ByReporterID(sql.OrderDesc())
		}
	case "created_at":
		if orderDir == "asc" {
			orderOption = taskreport.ByCreatedAt()
		} else {
			orderOption = taskreport.ByCreatedAt(sql.OrderDesc())
		}
	case "updated_at":
		if orderDir == "asc" {
			orderOption = taskreport.ByUpdatedAt()
		} else {
			orderOption = taskreport.ByUpdatedAt(sql.OrderDesc())
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order_by field. Valid fields: id, title, status, progress_percentage, reported_at, estimated_completion, task_id, reporter_id, created_at, updated_at",
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

	taskReports, err := query.All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get total count for pagination info
	total, err := h.Client.TaskReport.Query().Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := gin.H{
		"data": taskReports,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskReportHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task report ID"})
		return
	}

	taskReport, err := h.Client.TaskReport.Query().
		Where(taskreport.ID(id)).
		WithTask().
		WithReporter().
		Only(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task report not found"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, taskReport)
}

func (h *TaskReportHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task report ID"})
		return
	}

	var req struct {
		Title               *string `json:"title"`
		Content             *string `json:"content"`
		Status              *string `json:"status"`
		ProgressPercentage  *int    `json:"progress_percentage"`
		ReportedAt          *string `json:"reported_at"`
		IssuesEncountered   *string `json:"issues_encountered"`
		NextSteps           *string `json:"next_steps"`
		EstimatedCompletion *string `json:"estimated_completion"`
		TaskID              *int    `json:"task_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskReportUpdate := h.Client.TaskReport.UpdateOneID(id)

	if req.Title != nil {
		taskReportUpdate.SetTitle(*req.Title)
	}
	if req.Content != nil {
		taskReportUpdate.SetContent(*req.Content)
	}
	if req.Status != nil {
		switch *req.Status {
		case string(taskreport.StatusReceived),
			string(taskreport.StatusNotReceived),
			string(taskreport.StatusInProgress),
			string(taskreport.StatusCompleted),
			string(taskreport.StatusCancelled):
			taskReportUpdate.SetStatus(taskreport.Status(*req.Status))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value. Valid values: received, not_received, in_progress, completed, cancelled",
			})
			return
		}
	}
	if req.ProgressPercentage != nil {
		if *req.ProgressPercentage < 0 || *req.ProgressPercentage > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Progress percentage must be between 0 and 100",
			})
			return
		}
		taskReportUpdate.SetProgressPercentage(*req.ProgressPercentage)
	}
	if req.ReportedAt != nil {
		reportedAt, err := time.Parse(time.RFC3339, *req.ReportedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reported_at format, must be RFC3339"})
			return
		}
		taskReportUpdate.SetReportedAt(reportedAt)
	}
	if req.IssuesEncountered != nil {
		taskReportUpdate.SetIssuesEncountered(*req.IssuesEncountered)
	}
	if req.NextSteps != nil {
		taskReportUpdate.SetNextSteps(*req.NextSteps)
	}
	if req.EstimatedCompletion != nil {
		estimatedCompletion, err := time.Parse(time.RFC3339, *req.EstimatedCompletion)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid estimated_completion format, must be RFC3339"})
			return
		}
		taskReportUpdate.SetEstimatedCompletion(estimatedCompletion)
	}
	if req.TaskID != nil {
		// Validate that the task exists
		taskExists, err := h.Client.Task.Query().Where(task.ID(*req.TaskID)).Exist(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to validate task ID"})
			return
		}
		if !taskExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID does not exist"})
			return
		}
		taskReportUpdate.SetTaskID(*req.TaskID)
	}

	_, err = taskReportUpdate.Save(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task report not found"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Get the updated task report with all edges
	taskReport, err := h.Client.TaskReport.Query().
		Where(taskreport.ID(id)).
		WithTask().
		WithReporter().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"id": id})
		return
	}
	c.JSON(http.StatusOK, taskReport)
}

func (h *TaskReportHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task report ID"})
		return
	}

	_, err = h.Client.TaskReport.Delete().Where(taskreport.ID(id)).Exec(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task report not found"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *TaskReportHandler) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}

	count, err := h.Client.TaskReport.Delete().Where(taskreport.IDIn(req.IDs...)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Task reports deleted successfully",
		"deleted_count": count,
	})
}
