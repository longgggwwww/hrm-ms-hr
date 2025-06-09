package dtos

import (
	"github.com/go-playground/validator/v10"
)

// TaskCreateInput represents the input for creating a task
type TaskCreateInput struct {
	Name        string  `json:"name" binding:"required" validate:"required,min=1,max=200"`
	Code        *string `json:"code" validate:"omitempty,min=1,max=50"`
	Type        *string `json:"type" validate:"omitempty,oneof=task feature bug another"`
	StartAt     *string `json:"start_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DueDate     *string `json:"due_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	ProjectID   *int    `json:"project_id" validate:"omitempty,min=1"`
	LabelIDs    []int   `json:"label_ids" validate:"omitempty,dive,min=1"`
	AssigneeIDs []int   `json:"assignee_ids" validate:"omitempty,dive,min=1"`
}

// TaskUpdateInput represents the input for updating a task
type TaskUpdateInput struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=200"`
	Code        *string `json:"code" validate:"omitempty,min=1,max=50"`
	Type        *string `json:"type" validate:"omitempty,oneof=task feature bug another"`
	StartAt     *string `json:"start_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DueDate     *string `json:"due_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	ProjectID   *int    `json:"project_id" validate:"omitempty,min=1"`
	LabelIDs    []int   `json:"label_ids" validate:"omitempty,dive,min=1"`
	AssigneeIDs []int   `json:"assignee_ids" validate:"omitempty,dive,min=1"`
	Status      *string `json:"status" validate:"omitempty,oneof=not_received received in_progress completed cancelled"`
	Process     *int    `json:"process" validate:"omitempty,min=0,max=100"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
}

// TaskUpdateProgressInput represents the input for updating task progress
type TaskUpdateProgressInput struct {
	Status  *string `json:"status" validate:"omitempty,oneof=in_progress completed cancelled"`
	Process *int    `json:"process" validate:"omitempty,min=0,max=100"`
}

// TaskListQuery represents query parameters for listing tasks
type TaskListQuery struct {
	Name           string `form:"name"`
	Code           string `form:"code"`
	Status         string `form:"status" validate:"omitempty,oneof=not_received received in_progress completed cancelled"`
	Type           string `form:"type" validate:"omitempty,oneof=task feature bug another"`
	ProjectID      string `form:"project_id"`
	CreatorID      string `form:"creator_id"`
	Process        string `form:"process"`
	StartDateFrom  string `form:"start_date_from" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	StartDateTo    string `form:"start_date_to" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DueDateFrom    string `form:"due_date_from" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DueDateTo      string `form:"due_date_to" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	OrderBy        string `form:"order_by" validate:"omitempty,oneof=id name code status type process project_id creator_id start_at due_date created_at updated_at"`
	OrderDir       string `form:"order_dir" validate:"omitempty,oneof=asc desc"`
	Page           int    `form:"page" validate:"omitempty,min=1"`
	Limit          int    `form:"limit" validate:"omitempty,min=1,max=100"`
	Cursor         string `form:"cursor"`
	CursorLimit    int    `form:"cursor_limit" validate:"omitempty,min=1,max=100"`
	PaginationType string `form:"pagination_type" validate:"omitempty,oneof=page cursor"`
}

// TaskBulkDeleteInput represents input for bulk deleting tasks
type TaskBulkDeleteInput struct {
	IDs []int `json:"ids" binding:"required,min=1" validate:"required,min=1,max=100,dive,min=1"`
}

// TaskBulkDeleteResponse represents the response for bulk delete operations
type TaskBulkDeleteResponse struct {
	DeletedCount int      `json:"deleted_count"`
	FailedIDs    []int    `json:"failed_ids,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// TaskResponse represents a task with additional computed fields
type TaskResponse struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Description string      `json:"description"`
	Process     int         `json:"process"`
	Status      string      `json:"status"`
	Type        string      `json:"type"`
	StartAt     string      `json:"start_at,omitempty"`
	DueDate     string      `json:"due_date,omitempty"`
	ProjectID   int         `json:"project_id,omitempty"`
	CreatorID   int         `json:"creator_id"`
	UpdaterID   int         `json:"updater_id"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Edges       interface{} `json:"edges,omitempty"`
}

// TaskListResponse represents the response for task list with pagination
type TaskListResponse struct {
	Data       []TaskResponse `json:"data"`
	Pagination interface{}    `json:"pagination"`
}

// TaskOffsetPagination represents offset-based pagination info
type TaskOffsetPagination struct {
	Type        string `json:"type"`
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
	TotalItems  int    `json:"total_items"`
	PerPage     int    `json:"per_page"`
}

// TaskCursorPagination represents cursor-based pagination info
type TaskCursorPagination struct {
	Type       string  `json:"type"`
	PerPage    int     `json:"per_page"`
	HasNext    bool    `json:"has_next"`
	NextCursor *string `json:"next_cursor"`
}

// RegisterTaskValidators registers custom validators for task DTOs
func RegisterTaskValidators(v *validator.Validate) {
	// Custom validators can be added here if needed
	// For example, custom validation for task codes, date ranges, etc.
}
