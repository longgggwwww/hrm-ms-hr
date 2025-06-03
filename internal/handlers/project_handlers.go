package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	userpb "github.com/huynhthanhthao/hrm_user_service/proto/user"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/predicate"
	"github.com/longgggwwww/hrm-ms-hr/ent/project"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type ProjectHandler struct {
	Client     *ent.Client
	UserClient userpb.UserServiceClient
}

func NewProjectHandler(client *ent.Client, userClient userpb.UserServiceClient) *ProjectHandler {
	return &ProjectHandler{
		Client:     client,
		UserClient: userClient,
	}
}

func (h *ProjectHandler) RegisterRoutes(r *gin.Engine) {
	projs := r.Group("/projects")
	{
		projs.POST("/", h.Create)
		projs.GET("/", h.List)
		projs.GET("/:id", h.Get)
		projs.PATCH("/:id", h.Update)
		projs.DELETE("/:id", h.Delete)
		projs.DELETE("/", h.BulkDelete)
	}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Code        *string `json:"code"` // Made optional for auto-generation
		Description *string `json:"description"`
		StartAt     *string `json:"start_at"`
		EndAt       *string `json:"end_at"`
		Visibility  *string `json:"visibility"`
		MemberIDs   []int   `json:"member_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract org ID and employee ID from JWT token
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	orgID := ids["org_id"]
	employeeID := ids["employee_id"]

	var startAtPtr *time.Time
	if req.StartAt != nil {
		startAt, err := time.Parse(time.RFC3339, *req.StartAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_at format, must be RFC3339"})
			return
		}
		startAtPtr = &startAt
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

	// Auto-generate code if not provided
	var projectCode string
	if req.Code != nil && *req.Code != "" {
		projectCode = *req.Code

		// Check if code already exists in the organization
		existingProject, err := h.Client.Project.Query().
			Where(project.CodeEQ(projectCode)).
			Where(project.OrgIDEQ(orgID)).
			First(c.Request.Context())
		if err == nil && existingProject != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Project code already exists in your organization",
			})
			return
		}
		if err != nil && !ent.IsNotFound(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate project code"})
			return
		}
	} else {
		// Auto-generate code: ORG{orgId}-PROJ-{sequence}
		// Find the highest sequence number for this organization
		latestProjects, err := h.Client.Project.Query().
			Where(project.OrgIDEQ(orgID)).
			Where(project.CodeContains("ORG" + strconv.Itoa(orgID) + "-PROJ-")).
			Order(project.ByCreatedAt(sql.OrderDesc())).
			Limit(1).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate project code"})
			return
		}

		sequence := 1
		if len(latestProjects) > 0 {
			// Extract sequence number from the latest project code
			latestCode := latestProjects[0].Code
			prefix := "ORG" + strconv.Itoa(orgID) + "-PROJ-"
			if len(latestCode) > len(prefix) {
				sequenceStr := latestCode[len(prefix):]
				if seq, err := strconv.Atoi(sequenceStr); err == nil {
					sequence = seq + 1
				}
			}
		}

		projectCode = "ORG" + strconv.Itoa(orgID) + "-PROJ-" + fmt.Sprintf("%03d", sequence)

		// Double-check uniqueness (in case of concurrent requests)
		for {
			exists, err := h.Client.Project.Query().
				Where(project.CodeEQ(projectCode)).
				Where(project.OrgIDEQ(orgID)).
				Exist(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate generated project code"})
				return
			}
			if !exists {
				break
			}
			sequence++
			projectCode = "ORG" + strconv.Itoa(orgID) + "-PROJ-" + fmt.Sprintf("%03d", sequence)
		}
	}

	// Prepare all member IDs (current employee + provided member IDs)
	allMemberIDs := []int{}

	// Add current employee ID from token
	if employeeID > 0 {
		allMemberIDs = append(allMemberIDs, employeeID)
	}

	// Add provided member IDs (avoid duplicates)
	memberIDSet := make(map[int]bool)
	memberIDSet[employeeID] = true // Mark current employee as already added

	for _, memberID := range req.MemberIDs {
		if !memberIDSet[memberID] {
			allMemberIDs = append(allMemberIDs, memberID)
			memberIDSet[memberID] = true
		}
	}

	// Validate all member IDs if there are any
	if len(allMemberIDs) > 0 {
		// Check if all member IDs exist and belong to the same organization
		memberCount, err := h.Client.Employee.Query().
			Where(employee.IDIn(allMemberIDs...)).
			Where(employee.OrgID(orgID)).
			Count(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate member IDs"})
			return
		}
		if memberCount != len(allMemberIDs) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "One or more member IDs are invalid or do not belong to your organization",
			})
			return
		}
	}

	projectCreate := h.Client.Project.Create().
		SetName(req.Name).
		SetCode(projectCode).
		SetNillableDescription(req.Description).
		SetNillableStartAt(startAtPtr).
		SetNillableEndAt(endAtPtr).
		SetCreatorID(employeeID).
		SetUpdaterID(employeeID).
		SetOrgID(orgID).
		SetVisibility(visibilityVal)

	// Add member IDs if there are any
	if len(allMemberIDs) > 0 {
		projectCreate = projectCreate.AddMemberIDs(allMemberIDs...)
	}

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
		WithMembers().
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
	// Extract IDs from JWT token
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	orgID := ids["org_id"]
	employeeID := ids["employee_id"]

	// Build base query without hard filtering by org_id
	query := h.Client.Project.Query().
		WithOrganization().
		WithCreator().
		WithUpdater().
		WithMembers()

	// Apply visibility-based filtering
	visibilityConditions := []predicate.Project{}

	// Public projects - visible to everyone
	visibilityConditions = append(visibilityConditions, project.VisibilityEQ(project.VisibilityPublic))

	// Internal projects - visible if user belongs to the same org
	visibilityConditions = append(visibilityConditions,
		project.And(
			project.VisibilityEQ(project.VisibilityInternal),
			project.OrgIDEQ(orgID),
		),
	)

	// Private projects - visible only if user is a member
	visibilityConditions = append(visibilityConditions,
		project.And(
			project.VisibilityEQ(project.VisibilityPrivate),
			project.HasMembersWith(employee.IDEQ(employeeID)),
		),
	)

	// Apply OR condition for all visibility rules
	query = query.Where(project.Or(visibilityConditions...))

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

	// Filter by visibility (optional additional filter on top of access control)
	if visibility := c.Query("visibility"); visibility != "" {
		switch visibility {
		case string(project.VisibilityPrivate),
			string(project.VisibilityPublic),
			string(project.VisibilityInternal):
			// Apply additional filter on top of existing visibility rules
			additionalVisibilityConditions := []predicate.Project{}

			if visibility == string(project.VisibilityPublic) {
				additionalVisibilityConditions = append(additionalVisibilityConditions, project.VisibilityEQ(project.VisibilityPublic))
			} else if visibility == string(project.VisibilityInternal) {
				additionalVisibilityConditions = append(additionalVisibilityConditions,
					project.And(
						project.VisibilityEQ(project.VisibilityInternal),
						project.OrgIDEQ(orgID),
					),
				)
			} else if visibility == string(project.VisibilityPrivate) {
				additionalVisibilityConditions = append(additionalVisibilityConditions,
					project.And(
						project.VisibilityEQ(project.VisibilityPrivate),
						project.HasMembersWith(employee.IDEQ(employeeID)),
					),
				)
			}

			// Apply OR condition for additional visibility filter
			if len(additionalVisibilityConditions) > 0 {
				query = query.Where(project.Or(additionalVisibilityConditions...))
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid visibility value. Valid values: private, public, internal",
			})
			return
		}
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

	// Get task counts for all projects
	taskCounts := make(map[int]int)
	if len(projects) > 0 {
		var projectIDs []int
		for _, proj := range projects {
			projectIDs = append(projectIDs, proj.ID)
		}

		// Query task counts efficiently using a subquery approach
		for _, proj := range projects {
			count, err := h.Client.Task.Query().
				Where(task.ProjectIDEQ(proj.ID)).
				Count(c.Request.Context())
			if err == nil {
				taskCounts[proj.ID] = count
			}
		}
	}

	// Collect user IDs from all projects
	userIDs := h.collectUserIDsFromProjects(projects)

	// Fetch user information
	userMap, err := h.getUserInfoMap(userIDs)
	if err != nil {
		// Log error but continue without user enrichment
		// In production, you might want to use a proper logger
	}

	// Enrich projects with user information
	var enrichedProjects []map[string]interface{}
	for _, proj := range projects {
		taskCount := taskCounts[proj.ID] // Default to 0 if not found
		enrichedProject := h.enrichProjectWithUserInfo(proj, userMap, taskCount)
		enrichedProjects = append(enrichedProjects, enrichedProject)
	}

	response := gin.H{
		"data": enrichedProjects,
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
		WithOrganization().
		WithCreator().
		WithUpdater().
		WithMembers().
		WithTasks().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Collect user IDs from the project
	projects := []*ent.Project{project}
	userIDs := h.collectUserIDsFromProjects(projects)

	// Fetch user information
	userMap, err := h.getUserInfoMap(userIDs)
	if err != nil {
		// Log error but continue without user enrichment
		// In production, you might want to use a proper logger
	}

	// Enrich project with user information and tasks
	enrichedProject := h.enrichProjectWithUserInfoForGet(project, userMap)

	c.JSON(http.StatusOK, enrichedProject)
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
		// Check if the new code already exists in the organization (excluding current project)
		ids, err := utils.ExtractIDsFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		orgID := ids["org_id"]

		existingProject, err := h.Client.Project.Query().
			Where(project.CodeEQ(*req.Code)).
			Where(project.OrgIDEQ(orgID)).
			Where(project.IDNEQ(id)). // Exclude current project
			First(c.Request.Context())
		if err == nil && existingProject != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Project code already exists in your organization",
			})
			return
		}
		if err != nil && !ent.IsNotFound(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate project code"})
			return
		}

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
		WithMembers().
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

// getUserInfoMap fetches user information by user IDs and returns a map for quick lookup
func (h *ProjectHandler) getUserInfoMap(userIDs []int32) (map[int32]*userpb.User, error) {
	if h.UserClient == nil || len(userIDs) == 0 {
		return make(map[int32]*userpb.User), nil
	}

	// Remove duplicates
	uniqueIDs := make(map[int32]bool)
	var cleanIDs []int32
	for _, id := range userIDs {
		if !uniqueIDs[id] {
			uniqueIDs[id] = true
			cleanIDs = append(cleanIDs, id)
		}
	}

	userMap := make(map[int32]*userpb.User)

	// Fetch users individually (assuming GetUserById method exists)
	// If GetUsersByIds method exists, you can replace this with a bulk call
	for _, userID := range cleanIDs {
		response, err := h.UserClient.GetUserById(context.Background(), &userpb.GetUserByIdRequest{
			Id: userID,
		})
		if err != nil {
			// Log error but continue with other users
			continue
		}
		if response != nil && response.User != nil {
			userMap[userID] = response.User
		}
	}

	return userMap, nil
}

// enrichProjectWithUserInfo enriches a single project with user information
func (h *ProjectHandler) enrichProjectWithUserInfo(proj *ent.Project, userMap map[int32]*userpb.User, taskCount int) map[string]interface{} {
	result := map[string]interface{}{
		"id":          proj.ID,
		"name":        proj.Name,
		"code":        proj.Code,
		"description": proj.Description,
		"start_at":    proj.StartAt,
		"end_at":      proj.EndAt,
		"creator_id":  proj.CreatorID,
		"updater_id":  proj.UpdaterID,
		"org_id":      proj.OrgID,
		"process":     proj.Process,
		"status":      proj.Status,
		"visibility":  proj.Visibility,
		"created_at":  proj.CreatedAt,
		"updated_at":  proj.UpdatedAt,
	}

	// Create edges structure preserving original structure
	edges := make(map[string]interface{})

	// Add task_count instead of tasks array
	edges["task_count"] = taskCount

	// Add organization edge (unchanged)
	if proj.Edges.Organization != nil {
		edges["organization"] = proj.Edges.Organization
	}

	// Enrich creator with user_info while preserving original structure
	if proj.Edges.Creator != nil {
		// Start with the original creator data
		creatorData := map[string]interface{}{
			"id":          proj.Edges.Creator.ID,
			"code":        proj.Edges.Creator.Code,
			"position_id": proj.Edges.Creator.PositionID,
			"org_id":      proj.Edges.Creator.OrgID,
			"joining_at":  proj.Edges.Creator.JoiningAt,
			"status":      proj.Edges.Creator.Status,
			"created_at":  proj.Edges.Creator.CreatedAt,
			"updated_at":  proj.Edges.Creator.UpdatedAt,
		}

		// Add user_info if available
		if proj.Edges.Creator.UserID != "" {
			if userIDInt, err := strconv.Atoi(proj.Edges.Creator.UserID); err == nil {
				if userInfo, exists := userMap[int32(userIDInt)]; exists {
					creatorData["user_info"] = userInfo
				}
			}
		}
		edges["creator"] = creatorData
	}

	// Enrich updater with user_info while preserving original structure
	if proj.Edges.Updater != nil {
		// Start with the original updater data
		updaterData := map[string]interface{}{
			"id":          proj.Edges.Updater.ID,
			"code":        proj.Edges.Updater.Code,
			"position_id": proj.Edges.Updater.PositionID,
			"org_id":      proj.Edges.Updater.OrgID,
			"joining_at":  proj.Edges.Updater.JoiningAt,
			"status":      proj.Edges.Updater.Status,
			"created_at":  proj.Edges.Updater.CreatedAt,
			"updated_at":  proj.Edges.Updater.UpdatedAt,
		}

		// Add user_info if available
		if proj.Edges.Updater.UserID != "" {
			if userIDInt, err := strconv.Atoi(proj.Edges.Updater.UserID); err == nil {
				if userInfo, exists := userMap[int32(userIDInt)]; exists {
					updaterData["user_info"] = userInfo
				}
			}
		}
		edges["updater"] = updaterData
	}

	// Enrich members with user_info while preserving original structure
	if len(proj.Edges.Members) > 0 {
		var membersData []map[string]interface{}
		for _, member := range proj.Edges.Members {
			// Start with the original member data
			memberData := map[string]interface{}{
				"id":          member.ID,
				"code":        member.Code,
				"position_id": member.PositionID,
				"org_id":      member.OrgID,
				"joining_at":  member.JoiningAt,
				"status":      member.Status,
				"created_at":  member.CreatedAt,
				"updated_at":  member.UpdatedAt,
			}

			// Add user_info if available
			if member.UserID != "" {
				if userIDInt, err := strconv.Atoi(member.UserID); err == nil {
					if userInfo, exists := userMap[int32(userIDInt)]; exists {
						memberData["user_info"] = userInfo
					}
				}
			}
			membersData = append(membersData, memberData)
		}
		edges["members"] = membersData
	}

	// Add the edges structure to the result
	result["edges"] = edges

	return result
}

// enrichProjectWithUserInfoForGet enriches a single project with user information and tasks (for Get method)
func (h *ProjectHandler) enrichProjectWithUserInfoForGet(proj *ent.Project, userMap map[int32]*userpb.User) map[string]interface{} {
	result := map[string]interface{}{
		"id":          proj.ID,
		"name":        proj.Name,
		"code":        proj.Code,
		"description": proj.Description,
		"start_at":    proj.StartAt,
		"end_at":      proj.EndAt,
		"creator_id":  proj.CreatorID,
		"updater_id":  proj.UpdaterID,
		"org_id":      proj.OrgID,
		"process":     proj.Process,
		"status":      proj.Status,
		"visibility":  proj.Visibility,
		"created_at":  proj.CreatedAt,
		"updated_at":  proj.UpdatedAt,
	}

	// Create edges structure preserving original structure
	edges := make(map[string]interface{})

	// Add tasks array instead of task_count (for Get method)
	if proj.Edges.Tasks != nil {
		edges["tasks"] = proj.Edges.Tasks
	} else {
		edges["tasks"] = []interface{}{}
	}

	// Add organization edge (unchanged)
	if proj.Edges.Organization != nil {
		edges["organization"] = proj.Edges.Organization
	}

	// Enrich creator with user_info while preserving original structure
	if proj.Edges.Creator != nil {
		// Start with the original creator data
		creatorData := map[string]interface{}{
			"id":          proj.Edges.Creator.ID,
			"code":        proj.Edges.Creator.Code,
			"position_id": proj.Edges.Creator.PositionID,
			"org_id":      proj.Edges.Creator.OrgID,
			"joining_at":  proj.Edges.Creator.JoiningAt,
			"status":      proj.Edges.Creator.Status,
			"created_at":  proj.Edges.Creator.CreatedAt,
			"updated_at":  proj.Edges.Creator.UpdatedAt,
		}

		// Add user_info if available
		if proj.Edges.Creator.UserID != "" {
			if userIDInt, err := strconv.Atoi(proj.Edges.Creator.UserID); err == nil {
				if userInfo, exists := userMap[int32(userIDInt)]; exists {
					creatorData["user_info"] = userInfo
				}
			}
		}
		edges["creator"] = creatorData
	}

	// Enrich updater with user_info while preserving original structure
	if proj.Edges.Updater != nil {
		// Start with the original updater data
		updaterData := map[string]interface{}{
			"id":          proj.Edges.Updater.ID,
			"code":        proj.Edges.Updater.Code,
			"position_id": proj.Edges.Updater.PositionID,
			"org_id":      proj.Edges.Updater.OrgID,
			"joining_at":  proj.Edges.Updater.JoiningAt,
			"status":      proj.Edges.Updater.Status,
			"created_at":  proj.Edges.Updater.CreatedAt,
			"updated_at":  proj.Edges.Updater.UpdatedAt,
		}

		// Add user_info if available
		if proj.Edges.Updater.UserID != "" {
			if userIDInt, err := strconv.Atoi(proj.Edges.Updater.UserID); err == nil {
				if userInfo, exists := userMap[int32(userIDInt)]; exists {
					updaterData["user_info"] = userInfo
				}
			}
		}
		edges["updater"] = updaterData
	}

	// Enrich members with user_info while preserving original structure
	if len(proj.Edges.Members) > 0 {
		var membersData []map[string]interface{}
		for _, member := range proj.Edges.Members {
			// Start with the original member data
			memberData := map[string]interface{}{
				"id":          member.ID,
				"code":        member.Code,
				"position_id": member.PositionID,
				"org_id":      member.OrgID,
				"joining_at":  member.JoiningAt,
				"status":      member.Status,
				"created_at":  member.CreatedAt,
				"updated_at":  member.UpdatedAt,
			}

			// Add user_info if available
			if member.UserID != "" {
				if userIDInt, err := strconv.Atoi(member.UserID); err == nil {
					if userInfo, exists := userMap[int32(userIDInt)]; exists {
						memberData["user_info"] = userInfo
					}
				}
			}
			membersData = append(membersData, memberData)
		}
		edges["members"] = membersData
	}

	// Add the edges structure to the result
	result["edges"] = edges

	return result
}

// collectUserIDsFromProjects collects all user IDs from project creators, updaters, and members
func (h *ProjectHandler) collectUserIDsFromProjects(projects []*ent.Project) []int32 {
	userIDSet := make(map[int32]bool)
	var userIDs []int32

	for _, proj := range projects {
		// Creator user ID
		if proj.Edges.Creator != nil && proj.Edges.Creator.UserID != "" {
			if userID, err := strconv.Atoi(proj.Edges.Creator.UserID); err == nil {
				if !userIDSet[int32(userID)] {
					userIDSet[int32(userID)] = true
					userIDs = append(userIDs, int32(userID))
				}
			}
		}

		// Updater user ID
		if proj.Edges.Updater != nil && proj.Edges.Updater.UserID != "" {
			if userID, err := strconv.Atoi(proj.Edges.Updater.UserID); err == nil {
				if !userIDSet[int32(userID)] {
					userIDSet[int32(userID)] = true
					userIDs = append(userIDs, int32(userID))
				}
			}
		}

		// Members user IDs
		if proj.Edges.Members != nil {
			for _, member := range proj.Edges.Members {
				if member.UserID != "" {
					if userID, err := strconv.Atoi(member.UserID); err == nil {
						if !userIDSet[int32(userID)] {
							userIDSet[int32(userID)] = true
							userIDs = append(userIDs, int32(userID))
						}
					}
				}
			}
		}
	}

	return userIDs
}
