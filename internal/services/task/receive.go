package task

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
)

// ReceiveTask allows an assigned employee to receive/accept a task
// Only employees who are assigned to the task can receive it
// This changes the task status from "not_received" to "received"
func (s *TaskService) ReceiveTask(ctx context.Context, taskID, userID int) (*ent.Task, error) {
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

	// Check if task is in the correct status to be received
	if taskEntity.Status != task.StatusNotReceived {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "Task can only be received when status is 'not_received'. Current status: " + string(taskEntity.Status),
		}
	}

	// Update task status to "received"
	updatedTask, err := s.Client.Task.UpdateOneID(taskID).
		SetStatus(task.StatusReceived).
		SetUpdaterID(userID).
		Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to update task status",
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
