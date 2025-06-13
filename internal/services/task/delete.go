package task

import (
	"context"
	"net/http"
	"strconv"

	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Delete removes a single task by ID
func (s *TaskService) Delete(ctx context.Context, id int) error {
	affected, err := s.Client.Task.Delete().Where(task.ID(id)).Exec(ctx)
	if err != nil {
		return &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to delete task",
		}
	}

	if affected == 0 {
		return &ServiceError{
			Status: http.StatusNotFound,
			Msg:    "Task not found",
		}
	}

	return nil
}

// BulkDelete removes multiple tasks by their IDs
func (s *TaskService) BulkDelete(ctx context.Context, input dtos.TaskBulkDeleteInput) (*dtos.TaskBulkDeleteResponse, error) {
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
			Msg:    "No task IDs provided",
		}
	}

	// Check which tasks exist before attempting deletion
	existingTasks, err := s.Client.Task.Query().
		Where(task.IDIn(input.IDs...)).
		Select(task.FieldID).
		All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to validate task IDs",
		}
	}

	// Create a map of existing task IDs for quick lookup
	existingIDs := make(map[int]bool)
	for _, t := range existingTasks {
		existingIDs[t.ID] = true
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
		deletedCount, err = s.Client.Task.Delete().
			Where(task.IDIn(validIDs...)).
			Exec(ctx)
		if err != nil {
			// If deletion fails, add all valid IDs to failed IDs
			failedIDs = append(failedIDs, validIDs...)
			errors = append(errors, "Failed to delete tasks: "+err.Error())
		}
	}

	// Add not found IDs to failed IDs
	failedIDs = append(failedIDs, notFoundIDs...)
	for _, id := range notFoundIDs {
		errors = append(errors, "Task ID "+strconv.Itoa(id)+" not found")
	}

	return &dtos.TaskBulkDeleteResponse{
		DeletedCount: deletedCount,
		FailedIDs:    failedIDs,
		Errors:       errors,
	}, nil
}
