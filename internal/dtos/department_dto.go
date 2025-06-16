package dtos

import (
	"github.com/go-playground/validator/v10"
)

// DepartmentCreateInput represents the input for creating a department
type DepartmentCreateInput struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=100"`
	Code string `json:"code" binding:"required" validate:"required,min=1,max=50"`
}

// DepartmentUpdateInput represents the input for updating a department
type DepartmentUpdateInput struct {
	Name *string `json:"name" validate:"omitempty,min=1,max=100"`
	Code *string `json:"code" validate:"omitempty,min=1,max=50"`
}

// DepartmentBulkCreateInput represents the input for bulk creating departments
type DepartmentBulkCreateInput struct {
	Departments []DepartmentCreateInput `json:"departments" binding:"required" validate:"required,min=1,max=100,dive"`
}

// DepartmentListQuery represents query parameters for listing departments
type DepartmentListQuery struct {
	Name           string `form:"name"`
	Code           string `form:"code"`
	OrderBy        string `form:"order_by" validate:"omitempty,oneof=id name code created_at updated_at"`
	OrderDir       string `form:"order_dir" validate:"omitempty,oneof=asc desc"`
	Page           int    `form:"page" validate:"omitempty,min=1"`
	Limit          int    `form:"limit" validate:"omitempty,min=1,max=100"`
	Cursor         string `form:"cursor"`
	CursorLimit    int    `form:"cursor_limit" validate:"omitempty,min=1,max=100"`
	PaginationType string `form:"pagination_type" validate:"omitempty,oneof=page cursor"`
	OrgID          int    // From JWT token
}

// DepartmentDeleteBulkInput represents input for bulk deleting departments
type DepartmentDeleteBulkInput struct {
	IDs []int `json:"ids" binding:"required" validate:"required,min=1,max=100,dive,min=1"`
}

// DepartmentResponse represents a department with additional computed fields
type DepartmentResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	OrgID         int    `json:"org_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PositionCount int    `json:"position_count"`
}

// DepartmentListResponse represents the response for department list with pagination
type DepartmentListResponse struct {
	Data       []DepartmentResponse `json:"data"`
	Pagination interface{}          `json:"pagination"`
}

// DepartmentBulkDeleteResponse represents the response for bulk delete operations
type DepartmentBulkDeleteResponse struct {
	DeletedCount int      `json:"deleted_count"`
	FailedIDs    []int    `json:"failed_ids,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// DepartmentOffsetPagination represents offset-based pagination info
type DepartmentOffsetPagination struct {
	Type        string `json:"type"`
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
	TotalItems  int    `json:"total_items"`
	PerPage     int    `json:"per_page"`
}

// DepartmentCursorPagination represents cursor-based pagination info
type DepartmentCursorPagination struct {
	Type       string  `json:"type"`
	PerPage    int     `json:"per_page"`
	HasNext    bool    `json:"has_next"`
	NextCursor *string `json:"next_cursor"`
}

// RegisterDepartmentValidators registers custom validators for department DTOs
func RegisterDepartmentValidators(v *validator.Validate) {
	// Custom validators can be added here if needed
	// For example, custom validation for department codes, etc.
}
