package dto

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

// RegisterCustomValidators registers custom validators for label DTOs
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("color_hex", colorValidation)
}
