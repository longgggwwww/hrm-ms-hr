package task

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// List retrieves tasks with filtering, sorting, and pagination support
func (s *TaskService) List(ctx context.Context, query dtos.TaskListQuery) (map[string]interface{}, error) {
	taskQuery := s.Client.Task.Query().
		WithProject().
		WithLabels().
		WithAssignees()

	// Filter by name
	if query.Name != "" {
		taskQuery = taskQuery.Where(task.NameContains(query.Name))
	}

	// Filter by code
	if query.Code != "" {
		taskQuery = taskQuery.Where(task.CodeContains(query.Code))
	}

	// Filter by status
	if query.Status != "" {
		switch query.Status {
		case string(task.StatusNotReceived),
			string(task.StatusReceived),
			string(task.StatusInProgress),
			string(task.StatusCompleted),
			string(task.StatusCancelled):
			taskQuery = taskQuery.Where(task.StatusEQ(task.Status(query.Status)))
		default:
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid status value. Valid values: not_received, received, in_progress, completed, cancelled",
			}
		}
	}

	// Filter by type
	if query.Type != "" {
		switch query.Type {
		case string(task.TypeTask),
			string(task.TypeFeature),
			string(task.TypeBug),
			string(task.TypeAnother):
			taskQuery = taskQuery.Where(task.TypeEQ(task.Type(query.Type)))
		default:
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid type value. Valid values: task, feature, bug, another",
			}
		}
	}

	// Filter by project_id
	if query.ProjectID != "" {
		projectID, err := strconv.Atoi(query.ProjectID)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid project_id format",
			}
		}
		taskQuery = taskQuery.Where(task.ProjectIDEQ(projectID))
	}

	// Filter by creator_id
	if query.CreatorID != "" {
		creatorID, err := strconv.Atoi(query.CreatorID)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid creator_id format",
			}
		}
		taskQuery = taskQuery.Where(task.CreatorIDEQ(creatorID))
	}

	// Filter by process (percentage)
	if query.Process != "" {
		process, err := strconv.Atoi(query.Process)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid process format",
			}
		}
		taskQuery = taskQuery.Where(task.ProcessEQ(process))
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
		taskQuery = taskQuery.Where(task.StartAtGTE(startDate))
	}

	if query.StartDateTo != "" {
		startDate, err := time.Parse(time.RFC3339, query.StartDateTo)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid start_date_to format, must be RFC3339",
			}
		}
		taskQuery = taskQuery.Where(task.StartAtLTE(startDate))
	}

	// Due date range filtering
	if query.DueDateFrom != "" {
		dueDate, err := time.Parse(time.RFC3339, query.DueDateFrom)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid due_date_from format, must be RFC3339",
			}
		}
		taskQuery = taskQuery.Where(task.DueDateGTE(dueDate))
	}

	if query.DueDateTo != "" {
		dueDate, err := time.Parse(time.RFC3339, query.DueDateTo)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid due_date_to format, must be RFC3339",
			}
		}
		taskQuery = taskQuery.Where(task.DueDateLTE(dueDate))
	}

	// Order by field and direction
	orderBy := query.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	orderDir := query.OrderDir
	if orderDir == "" {
		orderDir = "desc"
	}

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
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "Invalid order_by field. Valid fields: id, name, code, status, type, process, project_id, creator_id, start_at, due_date, created_at, updated_at",
		}
	}

	// Apply ordering
	taskQuery = taskQuery.Order(orderOption)

	// Pagination
	page := 1
	limit := 10
	if query.Page > 0 {
		page = query.Page
	}
	if query.Limit > 0 && query.Limit <= 100 {
		limit = query.Limit
	}

	offset := (page - 1) * limit
	taskQuery = taskQuery.Offset(offset).Limit(limit)

	tasks, err := taskQuery.All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch tasks",
		}
	}

	// Get total count for pagination info
	total, err := s.Client.Task.Query().Count(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to count total tasks",
		}
	}

	totalPages := (total + limit - 1) / limit

	response := map[string]interface{}{
		"data": tasks,
		"pagination": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	}

	return response, nil
}
