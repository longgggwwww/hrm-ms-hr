package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/longgggwwww/hrm-ms-hr/internal/dto"
)

// InitValidator initializes custom validators for the application
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators for label DTOs
		dto.RegisterCustomValidators(v)
	}
}
