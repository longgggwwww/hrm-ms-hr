package task

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// UpdateProgress allows an assigned employee to update task status and progress
// Only employees who are assigned to the task can update its progress
func (s *TaskService) UpdateProgress(ctx context.Context, taskID, userID int, input dtos.TaskUpdateProgressInput) (*ent.Task, error) {
	// Validate at least one field is provided
	if input.Status == nil && input.Process == nil {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "At least one field (status or process) must be provided",
		}
	}

	// Get the task with assignees to check if user is assigned
	taskEntity, err := s.Client.Task.Query().
		Where(task.ID(taskID)).
		WithAssignees().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Task not found",
			}
		}
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch task",
		}
	}

	// Check if the current user is assigned to this task
	isAssigned := false
	for _, assignee := range taskEntity.Edges.Assignees {
		if assignee.ID == userID {
			isAssigned = true
			break
		}
	}

	if !isAssigned {
		return nil, &ServiceError{
			Status: http.StatusForbidden,
			Msg:    "You are not assigned to this task",
		}
	}

	// Check if task is in a valid status for progress updates
	if taskEntity.Status == task.StatusNotReceived {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "Task must be received before updating progress. Current status: " + string(taskEntity.Status),
		}
	}

	if taskEntity.Status == task.StatusCompleted || taskEntity.Status == task.StatusCancelled {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "Cannot update progress for completed or cancelled tasks. Current status: " + string(taskEntity.Status),
		}
	}

	taskUpdate := s.Client.Task.UpdateOneID(taskID).SetUpdaterID(userID)

	// Validate and set status if provided
	if input.Status != nil {
		switch *input.Status {
		case string(task.StatusInProgress),
			string(task.StatusCompleted),
			string(task.StatusCancelled):
			taskUpdate.SetStatus(task.Status(*input.Status))
		default:
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid status value. Valid values for progress update: in_progress, completed, cancelled",
			}
		}
	}

	// Validate and set process if provided
	if input.Process != nil {
		if *input.Process < 0 || *input.Process > 100 {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Process must be between 0 and 100",
			}
		}
		taskUpdate.SetProcess(*input.Process)

		// Auto-update status based on process
		if *input.Process == 100 && input.Status == nil {
			taskUpdate.SetStatus(task.StatusCompleted)
		} else if *input.Process > 0 && *input.Process < 100 && input.Status == nil && taskEntity.Status == task.StatusReceived {
			taskUpdate.SetStatus(task.StatusInProgress)
		}
	}

	updatedTask, err := taskUpdate.Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to update task progress",
		}
	}

	// Get the updated task with all edges
	taskWithEdges, err := s.Client.Task.Query().
		Where(task.IDEQ(updatedTask.ID)).
		WithProject().
		WithLabels().
		WithAssignees().
		Only(ctx)
	if err != nil {
		// If we can't fetch with edges, return the updated task without edges
		return updatedTask, nil
	}

	return taskWithEdges, nil
}
