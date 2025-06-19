package department

import (
	"context"
	"math"
	"net/http"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/department"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Get retrieves a single department by ID
func (s *DepartmentService) Get(ctx context.Context, id int) (*ent.Department, error) {
	dept, err := s.Client.Department.Query().
		Where(department.ID(id)).
		WithOrganization().
		WithPositions().
		WithZaloDepartment().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Department not found",
			}
		}
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to get department",
		}
	}

	return dept, nil
}

// List retrieves departments with filtering, sorting, and pagination
func (s *DepartmentService) List(ctx context.Context, query dtos.DepartmentListQuery) (*dtos.DepartmentListResponse, error) {
	// Build the base query
	q := s.Client.Department.Query().
		Where(department.OrgID(query.OrgID)).
		WithOrganization().
		WithPositions().
		WithZaloDepartment()

	// Apply filters
	if query.Name != "" {
		q = q.Where(department.NameContains(query.Name))
	}
	if query.Code != "" {
		q = q.Where(department.CodeContains(query.Code))
	}

	// Handle pagination type
	if query.PaginationType == "cursor" {
		return s.listWithCursorPagination(ctx, q, query)
	}

	// Default to offset pagination
	return s.listWithOffsetPagination(ctx, q, query)
}

// listWithOffsetPagination handles offset-based pagination
func (s *DepartmentService) listWithOffsetPagination(ctx context.Context, q *ent.DepartmentQuery, query dtos.DepartmentListQuery) (*dtos.DepartmentListResponse, error) {
	// Set defaults
	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Apply sorting
	orderBy := query.OrderBy
	if orderBy == "" {
		orderBy = "id"
	}
	orderDir := query.OrderDir
	if orderDir == "" {
		orderDir = "asc"
	}

	switch orderBy {
	case "id":
		if orderDir == "desc" {
			q = q.Order(department.ByID(sql.OrderDesc()))
		} else {
			q = q.Order(department.ByID())
		}
	case "name":
		if orderDir == "desc" {
			q = q.Order(department.ByName(sql.OrderDesc()))
		} else {
			q = q.Order(department.ByName())
		}
	case "code":
		if orderDir == "desc" {
			q = q.Order(department.ByCode(sql.OrderDesc()))
		} else {
			q = q.Order(department.ByCode())
		}
	case "created_at":
		if orderDir == "desc" {
			q = q.Order(department.ByCreatedAt(sql.OrderDesc()))
		} else {
			q = q.Order(department.ByCreatedAt())
		}
	case "updated_at":
		if orderDir == "desc" {
			q = q.Order(department.ByUpdatedAt(sql.OrderDesc()))
		} else {
			q = q.Order(department.ByUpdatedAt())
		}
	default:
		q = q.Order(department.ByID())
	}

	// Get total count
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to count departments",
		}
	}

	// Apply pagination
	offset := (page - 1) * limit
	departments, err := q.Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch departments",
		}
	}

	// Build response data
	var responses []dtos.DepartmentResponse
	for _, dept := range departments {
		positionCount, _ := s.getPositionCount(ctx, dept.ID)
		responses = append(responses, s.buildDepartmentResponse(dept, positionCount))
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	pagination := dtos.DepartmentOffsetPagination{
		Type:        "page",
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		PerPage:     limit,
	}

	return &dtos.DepartmentListResponse{
		Data:       responses,
		Pagination: pagination,
	}, nil
}

// listWithCursorPagination handles cursor-based pagination
func (s *DepartmentService) listWithCursorPagination(ctx context.Context, q *ent.DepartmentQuery, query dtos.DepartmentListQuery) (*dtos.DepartmentListResponse, error) {
	limit := query.CursorLimit
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Apply sorting (cursor pagination works best with ID sorting)
	q = q.Order(department.ByID())

	// Handle cursor
	if query.Cursor != "" {
		cursorData, err := s.decodeCursor(query.Cursor)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid cursor format",
			}
		}

		if lastID, ok := cursorData["last_id"].(float64); ok {
			q = q.Where(department.IDGT(int(lastID)))
		}
	}

	// Fetch one extra to determine if there's a next page
	departments, err := q.Limit(limit + 1).All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch departments",
		}
	}

	// Determine if there's a next page
	hasNext := len(departments) > limit
	if hasNext {
		departments = departments[:limit]
	}

	// Build response data
	var responses []dtos.DepartmentResponse
	for _, dept := range departments {
		positionCount, _ := s.getPositionCount(ctx, dept.ID)
		responses = append(responses, s.buildDepartmentResponse(dept, positionCount))
	}

	// Build pagination info
	var nextCursor *string
	if hasNext && len(departments) > 0 {
		lastDept := departments[len(departments)-1]
		cursorData := map[string]interface{}{
			"last_id": lastDept.ID,
		}
		cursor := s.encodeCursor(cursorData)
		nextCursor = &cursor
	}

	pagination := dtos.DepartmentCursorPagination{
		Type:       "cursor",
		PerPage:    limit,
		HasNext:    hasNext,
		NextCursor: nextCursor,
	}

	return &dtos.DepartmentListResponse{
		Data:       responses,
		Pagination: pagination,
	}, nil
}
