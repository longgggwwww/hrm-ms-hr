package handlers

import (
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/project"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type ProjectHandler struct {
	Client *ent.Client
}

func NewProjectHandler(client *ent.Client) *ProjectHandler {
	return &ProjectHandler{
		Client: client,
	}
}

func (h *ProjectHandler) RegisterRoutes(r *gin.Engine) {
	projs := r.Group("projects")
	{
		projs.POST("", h.Create)
		projs.GET("", h.List)
		projs.GET(":id", h.Get)
		projs.PATCH(":id", h.Update)
		projs.DELETE(":id", h.Delete)
		projs.DELETE("", h.BulkDelete)
	}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Code        string  `json:"code" binding:"required"`
		Description *string `json:"description"`
		StartAt     string  `json:"start_at" binding:"required"`
		EndAt       *string `json:"end_at"`
		OrgID       int     `json:"org_id" binding:"required"`
		Process     *int    `json:"process"`
		Status      *string `json:"status"`
		Visibility  *string `json:"visibility"`
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

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_at format, must be RFC3339"})
		return
	}

	var endAtPtr *time.Time
	if req.EndAt != nil {
		endAt, err := time.Parse(time.RFC3339, *req.EndAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid end_at format, must be RFC3339",
			})
			return
		}
		endAtPtr = &endAt
	}

	var statusVal project.Status
	if req.Status != nil {
		switch *req.Status {
		case string(project.StatusNotStarted),
			string(project.StatusInProgress),
			string(project.StatusCompleted):
			statusVal = project.Status(*req.Status)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value",
			})
			return
		}
	} else {
		statusVal = project.StatusNotStarted
	}

	var visibilityVal project.Visibility
	if req.Visibility != nil {
		switch *req.Visibility {
		case string(project.VisibilityPrivate),
			string(project.VisibilityPublic),
			string(project.VisibilityInternal):
			visibilityVal = project.Visibility(*req.Visibility)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid visibility value. Valid values: private, public, internal",
			})
			return
		}
	} else {
		visibilityVal = project.VisibilityPrivate
	}

	projectCreate := h.Client.Project.Create().
		SetName(req.Name).
		SetCode(req.Code).
		SetNillableDescription(req.Description).
		SetStartAt(startAt).
		SetNillableEndAt(endAtPtr).
		SetCreatorID(userID).
		SetUpdaterID(userID).
		SetOrgID(req.OrgID).
		SetNillableProcess(req.Process).
		SetStatus(statusVal).
		SetVisibility(visibilityVal)

	row, err := projectCreate.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Get the created project with all edges
	project, err := h.Client.Project.Query().
		Where(project.ID(row.ID)).
		WithTasks().
		WithOrganization().
		WithCreator().
		WithUpdater().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusCreated, row)
		return
	}
	c.JSON(http.StatusCreated, project)
}

// List retrieves projects with filtering, sorting, and pagination support.
//
// Query Parameters:
// - name: Filter by project name (contains search)
// - code: Filter by project code (contains search)
// - status: Filter by status (not_started, in_progress, completed)
// - visibility: Filter by visibility (private, public, internal)
// - org_id: Filter by organization ID
// - creator_id: Filter by creator ID
// - process: Filter by process percentage
// - start_date_from: Filter projects that start from this date (RFC3339 format)
// - start_date_to: Filter projects that start before this date (RFC3339 format)
// - order_by: Sort field (id, name, code, start_at, end_at, status, visibility, process, org_id, creator_id, created_at, updated_at, tasks_count)
// - order_dir: Sort direction (asc, desc) - default: desc
// - page: Page number (default: 1)
// - limit: Items per page (default: 10, max: 100)
//
// Example: GET /projects?name=example&status=in_progress&visibility=public&order_by=name&order_dir=asc&page=1&limit=20
func (h *ProjectHandler) List(c *gin.Context) {
	query := h.Client.Project.Query().
		WithTasks().
		WithOrganization().
		WithCreator().
		WithUpdater()

	// Filter by name
	if name := c.Query("name"); name != "" {
		query = query.Where(project.NameContains(name))
	}

	// Filter by code
	if code := c.Query("code"); code != "" {
		query = query.Where(project.CodeContains(code))
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		switch status {
		case string(project.StatusNotStarted),
			string(project.StatusInProgress),
			string(project.StatusCompleted):
			query = query.Where(project.StatusEQ(project.Status(status)))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value. Valid values: not_started, in_progress, completed",
			})
			return
		}
	}

	// Filter by visibility
	if visibility := c.Query("visibility"); visibility != "" {
		switch visibility {
		case string(project.VisibilityPrivate),
			string(project.VisibilityPublic),
			string(project.VisibilityInternal):
			query = query.Where(project.VisibilityEQ(project.Visibility(visibility)))
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid visibility value. Valid values: private, public, internal",
			})
			return
		}
	}

	// Filter by org_id
	if orgIDStr := c.Query("org_id"); orgIDStr != "" {
		orgID, err := strconv.Atoi(orgIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid org_id format",
			})
			return
		}
		query = query.Where(project.OrgIDEQ(orgID))
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
		query = query.Where(project.CreatorIDEQ(creatorID))
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
		query = query.Where(project.ProcessEQ(process))
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
		query = query.Where(project.StartAtGTE(startDate))
	}

	if startDateStr := c.Query("start_date_to"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid start_date_to format, must be RFC3339",
			})
			return
		}
		query = query.Where(project.StartAtLTE(startDate))
	}

	// Order by field and direction
	orderBy := c.DefaultQuery("order_by", "created_at")
	orderDir := c.DefaultQuery("order_dir", "desc")

	var orderOption project.OrderOption
	switch orderBy {
	case "id":
		if orderDir == "asc" {
			orderOption = project.ByID()
		} else {
			orderOption = project.ByID(sql.OrderDesc())
		}
	case "name":
		if orderDir == "asc" {
			orderOption = project.ByName()
		} else {
			orderOption = project.ByName(sql.OrderDesc())
		}
	case "code":
		if orderDir == "asc" {
			orderOption = project.ByCode()
		} else {
			orderOption = project.ByCode(sql.OrderDesc())
		}
	case "start_at":
		if orderDir == "asc" {
			orderOption = project.ByStartAt()
		} else {
			orderOption = project.ByStartAt(sql.OrderDesc())
		}
	case "end_at":
		if orderDir == "asc" {
			orderOption = project.ByEndAt()
		} else {
			orderOption = project.ByEndAt(sql.OrderDesc())
		}
	case "status":
		if orderDir == "asc" {
			orderOption = project.ByStatus()
		} else {
			orderOption = project.ByStatus(sql.OrderDesc())
		}
	case "visibility":
		if orderDir == "asc" {
			orderOption = project.ByVisibility()
		} else {
			orderOption = project.ByVisibility(sql.OrderDesc())
		}
	case "process":
		if orderDir == "asc" {
			orderOption = project.ByProcess()
		} else {
			orderOption = project.ByProcess(sql.OrderDesc())
		}
	case "org_id":
		if orderDir == "asc" {
			orderOption = project.ByOrgID()
		} else {
			orderOption = project.ByOrgID(sql.OrderDesc())
		}
	case "creator_id":
		if orderDir == "asc" {
			orderOption = project.ByCreatorID()
		} else {
			orderOption = project.ByCreatorID(sql.OrderDesc())
		}
	case "created_at":
		if orderDir == "asc" {
			orderOption = project.ByCreatedAt()
		} else {
			orderOption = project.ByCreatedAt(sql.OrderDesc())
		}
	case "updated_at":
		if orderDir == "asc" {
			orderOption = project.ByUpdatedAt()
		} else {
			orderOption = project.ByUpdatedAt(sql.OrderDesc())
		}
	case "tasks_count":
		if orderDir == "asc" {
			orderOption = project.ByTasksCount()
		} else {
			orderOption = project.ByTasksCount(sql.OrderDesc())
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order_by field. Valid fields: id, name, code, start_at, end_at, status, visibility, process, org_id, creator_id, created_at, updated_at, tasks_count",
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

	projects, err := query.All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get total count for pagination info
	total, err := h.Client.Project.Query().Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := gin.H{
		"data": projects,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.Client.Project.Query().
		Where(project.ID(id)).
		WithTasks().
		WithOrganization().
		WithCreator().
		WithUpdater().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Code        *string `json:"code"`
		Description *string `json:"description"`
		StartAt     *string `json:"start_at"`
		EndAt       *string `json:"end_at"`
		CreatorID   *int    `json:"creator_id"`
		UpdaterID   *int    `json:"updater_id"`
		OrgID       *int    `json:"org_id"`
		Process     *int    `json:"process"`
		Status      *string `json:"status"`
		Visibility  *string `json:"visibility"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectUpdate := h.Client.Project.UpdateOneID(id)
	if req.Name != nil {
		projectUpdate.SetName(*req.Name)
	}
	if req.Code != nil {
		projectUpdate.SetCode(*req.Code)
	}
	if req.Description != nil {
		projectUpdate.SetDescription(*req.Description)
	}
	if req.StartAt != nil {
		startAt, err := time.Parse(time.RFC3339, *req.StartAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_at format, must be RFC3339"})
			return
		}
		projectUpdate.SetStartAt(startAt)
	}
	if req.EndAt != nil {
		endAt, err := time.Parse(time.RFC3339, *req.EndAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_at format, must be RFC3339"})
			return
		}
		projectUpdate.SetEndAt(endAt)
	}
	if req.CreatorID != nil {
		projectUpdate.SetCreatorID(*req.CreatorID)
	}
	if req.UpdaterID != nil {
		projectUpdate.SetUpdaterID(*req.UpdaterID)
	}
	if req.OrgID != nil {
		projectUpdate.SetOrgID(*req.OrgID)
	}
	if req.Process != nil {
		projectUpdate.SetProcess(*req.Process)
	}
	if req.Status != nil {
		switch *req.Status {
		case string(project.StatusNotStarted), string(project.StatusInProgress), string(project.StatusCompleted):
			projectUpdate.SetStatus(project.Status(*req.Status))
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}
	}
	if req.Visibility != nil {
		switch *req.Visibility {
		case string(project.VisibilityPrivate), string(project.VisibilityPublic), string(project.VisibilityInternal):
			projectUpdate.SetVisibility(project.Visibility(*req.Visibility))
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visibility value. Valid values: private, public, internal"})
			return
		}
	}

	_, err = projectUpdate.Save(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Get the updated project with all edges
	project, err := h.Client.Project.Query().
		Where(project.ID(id)).
		WithTasks().
		WithOrganization().
		WithCreator().
		WithUpdater().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"id": id})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	_, err = h.Client.Project.Delete().Where(project.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// BulkDelete deletes multiple projects by their IDs.
//
// Request body should contain:
//
//	{
//	  "ids": [1, 2, 3, 4, 5]
//	}
//
// Response will include:
// - deleted_count: number of successfully deleted projects
// - failed_ids: array of IDs that failed to delete (if any)
// - errors: array of error messages for failed deletions (if any)
func (h *ProjectHandler) BulkDelete(c *gin.Context) {
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

	// Check which projects exist before attempting deletion
	existingProjects, err := h.Client.Project.Query().
		Where(project.IDIn(req.IDs...)).
		Select(project.FieldID).
		All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Create a map of existing project IDs for quick lookup
	existingIDs := make(map[int]bool)
	for _, p := range existingProjects {
		existingIDs[p.ID] = true
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
		deletedCount, err = h.Client.Project.Delete().
			Where(project.IDIn(validIDs...)).
			Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}

	// Add not found IDs to failed IDs
	failedIDs = append(failedIDs, notFoundIDs...)
	for _, id := range notFoundIDs {
		errors = append(errors, "Project ID "+strconv.Itoa(id)+" not found")
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
