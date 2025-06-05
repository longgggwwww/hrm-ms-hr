package project

import (
	"context"
	"net/http"
	"strconv"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/project"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Delete removes a single project by ID
func (s *ProjectService) Delete(ctx context.Context, id int) error {
	// Check if project exists
	exists, err := s.Client.Project.Query().
		Where(project.ID(id)).
		Exist(ctx)
	if err != nil {
		return &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to check project existence",
		}
	}
	if !exists {
		return &ServiceError{
			Status: http.StatusNotFound,
			Msg:    "Project not found",
		}
	}

	// Delete the project
	_, err = s.Client.Project.Delete().
		Where(project.ID(id)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Project not found",
			}
		}
		return &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to delete project",
		}
	}

	return nil
}

// DeleteBulk removes multiple projects by their IDs
func (s *ProjectService) DeleteBulk(ctx context.Context, input dtos.ProjectDeleteBulkInput) (*dtos.ProjectBulkDeleteResponse, error) {
	// Validate maximum number of IDs to prevent abuse
	if len(input.IDs) > 100 {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "Maximum 100 IDs allowed per bulk delete operation",
		}
	}

	if len(input.IDs) == 0 {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "No project IDs provided",
		}
	}

	// Check which projects exist before attempting deletion
	existingProjects, err := s.Client.Project.Query().
		Where(project.IDIn(input.IDs...)).
		Select(project.FieldID).
		All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to validate project IDs",
		}
	}

	// Create a map of existing project IDs for quick lookup
	existingIDs := make(map[int]bool)
	for _, p := range existingProjects {
		existingIDs[p.ID] = true
	}

	// Separate existing and non-existing IDs
	var validIDs []int
	var notFoundIDs []int
	for _, id := range input.IDs {
		if existingIDs[id] {
			validIDs = append(validIDs, id)
		} else {
			notFoundIDs = append(notFoundIDs, id)
		}
	}

	// Perform bulk deletion for valid IDs
	var deletedCount int
	var failedIDs []int
	var errors []string

	if len(validIDs) > 0 {
		deletedCount, err = s.Client.Project.Delete().
			Where(project.IDIn(validIDs...)).
			Exec(ctx)
		if err != nil {
			// If deletion fails, add all valid IDs to failed IDs
			failedIDs = append(failedIDs, validIDs...)
			errors = append(errors, "Failed to delete projects: "+err.Error())
		}
	}

	// Add not found IDs to failed IDs
	failedIDs = append(failedIDs, notFoundIDs...)
	for _, id := range notFoundIDs {
		errors = append(errors, "Project ID "+strconv.Itoa(id)+" not found")
	}

	return &dtos.ProjectBulkDeleteResponse{
		DeletedCount: deletedCount,
		FailedIDs:    failedIDs,
		Errors:       errors,
	}, nil
}
