package project

import (
	"context"
	"net/http"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/project"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
)

// List retrieves projects with filtering, sorting, and pagination support
func (s *ProjectService) List(ctx context.Context, query dtos.ProjectListQuery) (*dtos.ProjectListResponse, error) {
	// Build base query - only returns projects where the current employee is a member
	baseQuery := s.Client.Project.Query().
		WithOrganization().
		WithCreator().
		WithUpdater().
		WithMembers(func(q *ent.EmployeeQuery) {
			q.WithPosition() // Load position data for members
		}).
		Where(project.HasMembersWith(employee.IDEQ(query.EmployeeID)))

	// Apply filters
	if query.Name != "" {
		baseQuery = baseQuery.Where(project.NameContains(query.Name))
	}

	if query.Code != "" {
		baseQuery = baseQuery.Where(project.CodeContains(query.Code))
	}

	if query.Status != "" {
		switch query.Status {
		case string(project.StatusNotStarted),
			string(project.StatusInProgress),
			string(project.StatusCompleted):
			baseQuery = baseQuery.Where(project.StatusEQ(project.Status(query.Status)))
		default:
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid status value. Valid values: not_started, in_progress, completed",
			}
		}
	}

	if query.Process > 0 {
		baseQuery = baseQuery.Where(project.ProcessEQ(query.Process))
	}

	// Date range filtering
	if query.StartDateFrom != "" {
		startDate, err := time.Parse(time.RFC3339, query.StartDateFrom)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid start_date_from format, must be RFC3339",
			}
		}
		baseQuery = baseQuery.Where(project.StartAtGTE(startDate))
	}

	if query.StartDateTo != "" {
		startDate, err := time.Parse(time.RFC3339, query.StartDateTo)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid start_date_to format, must be RFC3339",
			}
		}
		baseQuery = baseQuery.Where(project.StartAtLTE(startDate))
	}

	// Handle pagination
	if query.PaginationType == "cursor" {
		return s.listWithCursorPagination(ctx, baseQuery, query)
	}

	return s.listWithOffsetPagination(ctx, baseQuery, query)
}

// listWithOffsetPagination handles offset-based pagination
func (s *ProjectService) listWithOffsetPagination(ctx context.Context, baseQuery *ent.ProjectQuery, query dtos.ProjectListQuery) (*dtos.ProjectListResponse, error) {
	// Set defaults
	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Apply ordering
	orderBy := query.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	orderDir := query.OrderDir
	if orderDir == "" {
		orderDir = "desc"
	}

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
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "Invalid order_by field. Valid fields: id, name, code, start_at, end_at, status, process, org_id, created_at, updated_at, tasks_count",
		}
	}

	// Create a copy of baseQuery for count
	countQuery := s.Client.Project.Query().
		Where(project.HasMembersWith(employee.IDEQ(query.EmployeeID)))

	// Apply the same filters to count query
	if query.Name != "" {
		countQuery = countQuery.Where(project.NameContains(query.Name))
	}
	if query.Code != "" {
		countQuery = countQuery.Where(project.CodeContains(query.Code))
	}
	if query.Status != "" {
		switch query.Status {
		case string(project.StatusNotStarted),
			string(project.StatusInProgress),
			string(project.StatusCompleted):
			countQuery = countQuery.Where(project.StatusEQ(project.Status(query.Status)))
		}
	}
	if query.Process > 0 {
		countQuery = countQuery.Where(project.ProcessEQ(query.Process))
	}
	if query.StartDateFrom != "" {
		if startDate, err := time.Parse(time.RFC3339, query.StartDateFrom); err == nil {
			countQuery = countQuery.Where(project.StartAtGTE(startDate))
		}
	}
	if query.StartDateTo != "" {
		if startDate, err := time.Parse(time.RFC3339, query.StartDateTo); err == nil {
			countQuery = countQuery.Where(project.StartAtLTE(startDate))
		}
	}

	// Get total count
	total, err := countQuery.Count(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to count projects",
		}
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	projects, err := baseQuery.
		Order(orderOption).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch projects",
		}
	}

	totalPages := (total + limit - 1) / limit

	// Enrich projects with task counts and user information
	enrichedProjects, err := s.addTaskCountsToProjects(ctx, projects)
	if err != nil {
		return nil, err
	}

	pagination := dtos.ProjectOffsetPagination{
		Type:        "page",
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		PerPage:     limit,
	}

	return &dtos.ProjectListResponse{
		Data:       enrichedProjects,
		Pagination: pagination,
	}, nil
}

// listWithCursorPagination handles cursor-based pagination
func (s *ProjectService) listWithCursorPagination(ctx context.Context, baseQuery *ent.ProjectQuery, query dtos.ProjectListQuery) (*dtos.ProjectListResponse, error) {
	limit := query.CursorLimit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Apply ordering (cursor pagination typically uses ID ordering)
	baseQuery = baseQuery.Order(project.ByID(sql.OrderDesc()))

	// Apply cursor if provided
	if query.Cursor != "" {
		cursorData, err := s.decodeCursor(query.Cursor)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid cursor format",
			}
		}

		if lastID, ok := cursorData["last_id"].(float64); ok {
			baseQuery = baseQuery.Where(project.IDLT(int(lastID)))
		}
	}

	// Fetch limit + 1 to check if there are more items
	projects, err := baseQuery.Limit(limit + 1).All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch projects",
		}
	}

	hasNext := len(projects) > limit
	if hasNext {
		projects = projects[:limit] // Remove the extra item
	}

	// Enrich projects with task counts and user information
	enrichedProjects, err := s.addTaskCountsToProjects(ctx, projects)
	if err != nil {
		return nil, err
	}

	// Generate next cursor if there are more items
	var nextCursor *string
	if hasNext && len(projects) > 0 {
		lastProject := projects[len(projects)-1]
		cursorData := map[string]interface{}{
			"last_id": lastProject.ID,
		}
		cursor := s.encodeCursor(cursorData)
		nextCursor = &cursor
	}

	pagination := dtos.ProjectCursorPagination{
		Type:       "cursor",
		PerPage:    limit,
		HasNext:    hasNext,
		NextCursor: nextCursor,
	}

	return &dtos.ProjectListResponse{
		Data:       enrichedProjects,
		Pagination: pagination,
	}, nil
}

// Get retrieves a single project by ID
func (s *ProjectService) Get(ctx context.Context, id int) (map[string]interface{}, error) {
	proj, err := s.Client.Project.Query().
		Where(project.ID(id)).
		WithOrganization().
		WithCreator().
		WithUpdater().
		WithMembers(func(q *ent.EmployeeQuery) {
			q.WithPosition() // Load position data for members
		}).
		WithTasks().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Project not found",
			}
		}
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch project",
		}
	}

	// Collect user IDs from the project
	projects := []*ent.Project{proj}
	userIDs := s.collectUserIDsFromProjects(projects)

	// Fetch user information
	userMap, err := s.getUserInfoMap(userIDs)
	if err != nil {
		// Log error but continue without user enrichment
		userMap = make(map[int32]*grpc_clients.User)
	}

	// Enrich project with user information and tasks
	enrichedProject := s.enrichProjectWithUserInfoForGet(proj, userMap)

	return enrichedProject, nil
}
