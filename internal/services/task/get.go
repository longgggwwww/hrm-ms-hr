package task

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
)

// Get retrieves a single task by ID with all related data
func (s *TaskService) Get(ctx context.Context, id int) (*ent.Task, error) {
	task, err := s.Client.Task.Query().
		Where(task.ID(id)).
		WithProject().
		WithDepartment().
		WithLabels().
		WithAssignees().
		WithReports().
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

	return task, nil
}
