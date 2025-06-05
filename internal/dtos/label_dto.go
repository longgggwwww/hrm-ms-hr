package dtos

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// colorValidation validates that color follows #xxxxxx format
func colorValidation(fl validator.FieldLevel) bool {
	color := fl.Field().String()
	// Regex pattern for #xxxxxx format (6 hex characters)
	pattern := `^#[a-fA-F0-9]{6}$`
	matched, _ := regexp.MatchString(pattern, color)
	return matched
}

// LabelCreateInput represents the input for creating a label
type LabelCreateInput struct {
	Name        string `json:"name" binding:"required" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
	Color       string `json:"color" binding:"required" validate:"required,color_hex"`
}

// LabelUpdateInput represents the input for updating a label
type LabelUpdateInput struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Color       *string `json:"color" validate:"omitempty,color_hex"`
	OrgID       *int    `json:"org_id"`
}

// LabelBulkCreateInput represents the input for bulk creating labels
type LabelBulkCreateInput struct {
	Labels []LabelCreateInput `json:"labels" binding:"required" validate:"required,min=1,max=100,dive"`
}

// LabelListQuery represents query parameters for listing labels
type LabelListQuery struct {
	Name           string `form:"name"`
	Description    string `form:"description"`
	Color          string `form:"color"`
	OrderBy        string `form:"order_by"`
	OrderDir       string `form:"order_dir"`
	Page           int    `form:"page"`
	Limit          int    `form:"limit"`
	Cursor         string `form:"cursor"`
	CursorLimit    int    `form:"cursor_limit"`
	PaginationType string `form:"pagination_type"`
	OrgID          int    // From JWT token
}

// LabelDeleteBulkInput represents input for bulk deleting labels
type LabelDeleteBulkInput struct {
	IDs []int `json:"ids" binding:"required"`
}

// LabelResponse represents a label with task count
type LabelResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	OrgID       int    `json:"org_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	TaskCount   int    `json:"task_count"`
}

// LabelListResponse represents the response for label list with pagination
type LabelListResponse struct {
	Data       []LabelResponse `json:"data"`
	Pagination interface{}     `json:"pagination"`
}

// OffsetPagination represents offset-based pagination info
type OffsetPagination struct {
	Type        string `json:"type"`
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
	TotalItems  int    `json:"total_items"`
	PerPage     int    `json:"per_page"`
}

// CursorPagination represents cursor-based pagination info
type CursorPagination struct {
	Type       string  `json:"type"`
	PerPage    int     `json:"per_page"`
	HasNext    bool    `json:"has_next"`
	NextCursor *string `json:"next_cursor"`
}

// RegisterCustomValidators registers custom validators for label DTOs
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("color_hex", colorValidation)
}
