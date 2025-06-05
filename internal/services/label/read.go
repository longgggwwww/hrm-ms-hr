package label

import (
	"context"
	"net/http"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/label"
	"github.com/longgggwwww/hrm-ms-hr/internal/dto"
)

// Get retrieves a single label by ID
func (s *LabelService) Get(ctx context.Context, id int) (*ent.Label, error) {
	labelObj, err := s.Client.Label.Query().
		Where(label.ID(id)).
		WithTasks().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Label not found",
			}
		}
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch label",
		}
	}

	return labelObj, nil
}

// List retrieves labels with filtering, sorting, and pagination
func (s *LabelService) List(ctx context.Context, query dto.LabelListQuery) (*dto.LabelListResponse, error) {
	labelQuery := s.Client.Label.Query()

	// Filter by org_id from token
	labelQuery = labelQuery.Where(label.OrgIDEQ(query.OrgID))

	// Apply filters
	if query.Name != "" {
		labelQuery = labelQuery.Where(label.NameContainsFold(query.Name))
	}

	if query.Description != "" {
		labelQuery = labelQuery.Where(label.DescriptionContainsFold(query.Description))
	}

	if query.Color != "" {
		labelQuery = labelQuery.Where(label.ColorEQ(query.Color))
	}

	// Determine pagination type
	if query.PaginationType == "cursor" {
		return s.listWithCursorPagination(ctx, labelQuery, query)
	}

	// Default: offset-based pagination
	return s.listWithOffsetPagination(ctx, labelQuery, query)
}

// listWithOffsetPagination handles offset-based pagination
func (s *LabelService) listWithOffsetPagination(ctx context.Context, labelQuery *ent.LabelQuery, query dto.LabelListQuery) (*dto.LabelListResponse, error) {
	// Apply sorting
	orderOption := s.getOrderOption(query.OrderBy, query.OrderDir)
	if orderOption == nil {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "Invalid order_by field. Valid fields: id, name, description, color, org_id, created_at, updated_at, task_count",
		}
	}
	labelQuery = labelQuery.Order(orderOption)

	// Set default pagination values
	page := query.Page
	if page <= 0 {
		page = 1
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	labelQuery = labelQuery.Offset(offset).Limit(limit)

	labels, err := labelQuery.All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	// Get total count for pagination info with same filters
	countQuery := s.Client.Label.Query()
	countQuery = countQuery.Where(label.OrgIDEQ(query.OrgID))

	if query.Name != "" {
		countQuery = countQuery.Where(label.NameContainsFold(query.Name))
	}
	if query.Description != "" {
		countQuery = countQuery.Where(label.DescriptionContainsFold(query.Description))
	}
	if query.Color != "" {
		countQuery = countQuery.Where(label.ColorEQ(query.Color))
	}

	total, err := countQuery.Count(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	totalPages := (total + limit - 1) / limit

	// Add task counts to labels
	labelsWithTaskCount, err := s.addTaskCountsToLabels(ctx, labels)
	if err != nil {
		return nil, err
	}

	response := &dto.LabelListResponse{
		Data: labelsWithTaskCount,
		Pagination: dto.OffsetPagination{
			Type:        "offset",
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalItems:  total,
			PerPage:     limit,
		},
	}

	return response, nil
}

// listWithCursorPagination handles cursor-based pagination
func (s *LabelService) listWithCursorPagination(ctx context.Context, labelQuery *ent.LabelQuery, query dto.LabelListQuery) (*dto.LabelListResponse, error) {
	// Apply sorting
	orderOption := s.getOrderOption(query.OrderBy, query.OrderDir)
	if orderOption == nil {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "Invalid order_by field",
		}
	}
	labelQuery = labelQuery.Order(orderOption)

	// Set cursor pagination limit
	limit := query.CursorLimit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Parse cursor if provided
	if query.Cursor != "" {
		cursorData, err := s.decodeCursor(query.Cursor)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid cursor format",
			}
		}

		// Apply cursor filter (assuming cursor based on ID)
		if id, ok := cursorData["id"].(float64); ok {
			labelQuery = labelQuery.Where(label.IDGT(int(id)))
		}
	}

	// Fetch one extra item to determine if there's a next page
	labels, err := labelQuery.Limit(limit + 1).All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	hasNext := len(labels) > limit
	if hasNext {
		labels = labels[:limit] // Remove the extra item
	}

	var nextCursor *string
	if hasNext && len(labels) > 0 {
		lastItem := labels[len(labels)-1]
		cursorStr := s.encodeCursor(map[string]interface{}{
			"id": lastItem.ID,
		})
		nextCursor = &cursorStr
	}

	// Add task counts to labels
	labelsWithTaskCount, err := s.addTaskCountsToLabels(ctx, labels)
	if err != nil {
		return nil, err
	}

	response := &dto.LabelListResponse{
		Data: labelsWithTaskCount,
		Pagination: dto.CursorPagination{
			Type:       "cursor",
			PerPage:    limit,
			HasNext:    hasNext,
			NextCursor: nextCursor,
		},
	}

	return response, nil
}

// getOrderOption returns the appropriate order option based on query parameters
func (s *LabelService) getOrderOption(orderBy, orderDir string) label.OrderOption {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if orderDir == "" {
		orderDir = "desc"
	}

	var orderOption label.OrderOption
	switch orderBy {
	case "id":
		if orderDir == "asc" {
			orderOption = label.ByID()
		} else {
			orderOption = label.ByID(sql.OrderDesc())
		}
	case "name":
		if orderDir == "asc" {
			orderOption = label.ByName()
		} else {
			orderOption = label.ByName(sql.OrderDesc())
		}
	case "description":
		if orderDir == "asc" {
			orderOption = label.ByDescription()
		} else {
			orderOption = label.ByDescription(sql.OrderDesc())
		}
	case "color":
		if orderDir == "asc" {
			orderOption = label.ByColor()
		} else {
			orderOption = label.ByColor(sql.OrderDesc())
		}
	case "org_id":
		if orderDir == "asc" {
			orderOption = label.ByOrgID()
		} else {
			orderOption = label.ByOrgID(sql.OrderDesc())
		}
	case "created_at":
		if orderDir == "asc" {
			orderOption = label.ByCreatedAt()
		} else {
			orderOption = label.ByCreatedAt(sql.OrderDesc())
		}
	case "updated_at":
		if orderDir == "asc" {
			orderOption = label.ByUpdatedAt()
		} else {
			orderOption = label.ByUpdatedAt(sql.OrderDesc())
		}
	case "task_count":
		// Sort by task count using ByTasksCount from Ent
		if orderDir == "asc" {
			orderOption = label.ByTasksCount()
		} else {
			orderOption = label.ByTasksCount(sql.OrderDesc())
		}
	default:
		return nil
	}

	return orderOption
}
