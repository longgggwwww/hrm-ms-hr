package dtos

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// projectCodeValidation validates that project code follows the expected format
func projectCodeValidation(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	// Regex pattern for ORG{orgId}-PROJ-{sequence} format
	pattern := `^ORG\d+-PROJ-\d{3,}$`
	matched, _ := regexp.MatchString(pattern, code)
	return matched
}

// ProjectCreateInput represents the input for creating a project
type ProjectCreateInput struct {
	Name        string  `json:"name" binding:"required" validate:"required,min=1,max=200"`
	Code        *string `json:"code" validate:"omitempty,project_code"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	StartAt     *string `json:"start_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndAt       *string `json:"end_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	MemberIDs   []int   `json:"member_ids" validate:"omitempty,dive,min=1"`
}

// ProjectUpdateInput represents the input for updating a project
type ProjectUpdateInput struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=200"`
	Code        *string `json:"code" validate:"omitempty,project_code"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	StartAt     *string `json:"start_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndAt       *string `json:"end_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	CreatorID   *int    `json:"creator_id" validate:"omitempty,min=1"`
	UpdaterID   *int    `json:"updater_id" validate:"omitempty,min=1"`
	OrgID       *int    `json:"org_id" validate:"omitempty,min=1"`
	Process     *int    `json:"process" validate:"omitempty,min=0,max=100"`
	Status      *string `json:"status" validate:"omitempty,oneof=not_started in_progress completed"`
	MemberIDs   []int   `json:"member_ids" validate:"omitempty,dive,min=1"`
}

// ProjectListQuery represents query parameters for listing projects
type ProjectListQuery struct {
	Name           string `form:"name"`
	Code           string `form:"code"`
	Status         string `form:"status" validate:"omitempty,oneof=not_started in_progress completed"`
	OrgID          int    `form:"org_id"`
	Process        int    `form:"process" validate:"omitempty,min=0,max=100"`
	StartDateFrom  string `form:"start_date_from" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	StartDateTo    string `form:"start_date_to" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	OrderBy        string `form:"order_by"`
	OrderDir       string `form:"order_dir" validate:"omitempty,oneof=asc desc"`
	Page           int    `form:"page" validate:"omitempty,min=1"`
	Limit          int    `form:"limit" validate:"omitempty,min=1,max=100"`
	Cursor         string `form:"cursor"`
	CursorLimit    int    `form:"cursor_limit" validate:"omitempty,min=1,max=100"`
	PaginationType string `form:"pagination_type" validate:"omitempty,oneof=page cursor"`
	EmployeeID     int    // From JWT token
}

// ProjectDeleteBulkInput represents input for bulk deleting projects
type ProjectDeleteBulkInput struct {
	IDs []int `json:"ids" binding:"required" validate:"required,min=1,max=100,dive,min=1"`
}

// ProjectBulkDeleteResponse represents the response for bulk delete operations
type ProjectBulkDeleteResponse struct {
	DeletedCount int      `json:"deleted_count"`
	FailedIDs    []int    `json:"failed_ids,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// ProjectResponse represents a project with additional computed fields
type ProjectResponse struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Description *string     `json:"description"`
	StartAt     *time.Time  `json:"start_at"`
	EndAt       *time.Time  `json:"end_at"`
	CreatorID   int         `json:"creator_id"`
	UpdaterID   int         `json:"updater_id"`
	OrgID       int         `json:"org_id"`
	Process     int         `json:"process"`
	Status      string      `json:"status"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	TaskCount   int         `json:"task_count"`
	Edges       interface{} `json:"edges,omitempty"`
}

// ProjectListResponse represents the response for project list with pagination
type ProjectListResponse struct {
	Data       []ProjectResponse `json:"data"`
	Pagination interface{}       `json:"pagination"`
}

// ProjectOffsetPagination represents offset-based pagination info
type ProjectOffsetPagination struct {
	Type        string `json:"type"`
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
	TotalItems  int    `json:"total_items"`
	PerPage     int    `json:"per_page"`
}

// ProjectCursorPagination represents cursor-based pagination info
type ProjectCursorPagination struct {
	Type       string  `json:"type"`
	PerPage    int     `json:"per_page"`
	HasNext    bool    `json:"has_next"`
	NextCursor *string `json:"next_cursor"`
}

// ProjectAddMembersInput represents the input for adding members to a project
type ProjectAddMembersInput struct {
	MemberIDs []int `json:"member_ids" binding:"required" validate:"required,dive,min=1"`
}

// ProjectRemoveMembersInput represents the input for removing members from a project
type ProjectRemoveMembersInput struct {
	MemberIDs []int `json:"member_ids" binding:"required" validate:"required,dive,min=1"`
}

// RegisterProjectValidators registers custom validators for project DTOs
func RegisterProjectValidators(v *validator.Validate) {
	v.RegisterValidation("project_code", projectCodeValidation)
}
