package task

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Update updates an existing task
func (s *TaskService) Update(ctx context.Context, id, userID int, input dtos.TaskUpdateInput) (*ent.Task, error) {
	taskUpdate := s.Client.Task.UpdateOneID(id).SetUpdaterID(userID)

	if input.Name != nil {
		taskUpdate.SetName(*input.Name)
	}

	// Handle code update with auto-generation if needed
	if input.Code != nil {
		var taskCode string
		if *input.Code != "" {
			taskCode = *input.Code

			// Check if the new code already exists (excluding current task)
			var existsQuery *ent.TaskQuery
			if input.ProjectID != nil {
				// Check uniqueness within the project
				existsQuery = s.Client.Task.Query().
					Where(task.CodeEQ(taskCode)).
					Where(task.ProjectIDEQ(*input.ProjectID)).
					Where(task.IDNEQ(id))
			} else {
				// Check global uniqueness if no project
				existsQuery = s.Client.Task.Query().
					Where(task.CodeEQ(taskCode)).
					Where(task.IDNEQ(id))
			}

			exists, err := existsQuery.Exist(ctx)
			if err != nil {
				return nil, &ServiceError{
					Status: http.StatusInternalServerError,
					Msg:    "Failed to validate task code",
				}
			}
			if exists {
				return nil, &ServiceError{
					Status: http.StatusBadRequest,
					Msg:    "Task code already exists",
				}
			}
		} else {
			// Auto-generate code in GitHub-style format: #x
			// Get the total count of tasks to generate the next sequence number
			count, err := s.Client.Task.Query().
				Where(task.IDNEQ(id)). // Exclude current task
				Count(ctx)
			if err != nil {
				return nil, &ServiceError{
					Status: http.StatusInternalServerError,
					Msg:    "Failed to generate task code",
				}
			}

			sequence := count + 1
			taskCode = "#" + strconv.Itoa(sequence)

			// Double-check uniqueness (in case of concurrent requests)
			for {
				exists, err := s.Client.Task.Query().
					Where(task.CodeEQ(taskCode)).
					Where(task.IDNEQ(id)).
					Exist(ctx)
				if err != nil {
					return nil, &ServiceError{
						Status: http.StatusInternalServerError,
						Msg:    "Failed to validate generated task code",
					}
				}
				if !exists {
					break
				}
				sequence++
				taskCode = "#" + strconv.Itoa(sequence)
			}
		}
		taskUpdate.SetCode(taskCode)
	}

	if input.Description != nil {
		taskUpdate.SetDescription(*input.Description)
	}
	if input.Process != nil {
		taskUpdate.SetProcess(*input.Process)
	}
	if input.StartAt != nil {
		startAt, err := time.Parse(time.RFC3339, *input.StartAt)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid start_at format, must be RFC3339",
			}
		}
		taskUpdate.SetStartAt(startAt)
	}
	if input.DueDate != nil {
		dueDate, err := time.Parse(time.RFC3339, *input.DueDate)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid due_date format, must be RFC3339",
			}
		}
		taskUpdate.SetDueDate(dueDate)
	}
	if input.ProjectID != nil {
		taskUpdate.SetProjectID(*input.ProjectID)
	}
	if input.Status != nil {
		switch *input.Status {
		case string(task.StatusNotReceived),
			string(task.StatusReceived),
			string(task.StatusInProgress),
			string(task.StatusCompleted),
			string(task.StatusCancelled):
			taskUpdate.SetStatus(task.Status(*input.Status))
		default:
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid status value. Valid values: not_received, received, in_progress, completed, cancelled",
			}
		}
	}
	if input.Type != nil {
		switch *input.Type {
		case string(task.TypeTask),
			string(task.TypeFeature),
			string(task.TypeBug),
			string(task.TypeAnother):
			taskUpdate.SetType(task.Type(*input.Type))
		default:
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid type value. Valid values: task, feature, bug, another",
			}
		}
	}

	// Handle label assignments
	if input.LabelIDs != nil {
		// Clear existing labels and add new ones
		taskUpdate.ClearLabels()
		if len(input.LabelIDs) > 0 {
			taskUpdate.AddLabelIDs(input.LabelIDs...)
		}
	}

	// Handle assignee assignments
	if input.AssigneeIDs != nil {
		// Clear existing assignees first
		taskUpdate.ClearAssignees()

		if len(input.AssigneeIDs) > 0 {
			// Validate that all assignee IDs exist in the employee table
			existingEmployees, err := s.Client.Employee.Query().
				Where(employee.IDIn(input.AssigneeIDs...)).
				Select(employee.FieldID).
				All(ctx)
			if err != nil {
				return nil, &ServiceError{
					Status: http.StatusInternalServerError,
					Msg:    "Failed to validate assignee IDs",
				}
			}

			// Create map of existing employee IDs for validation
			existingIDs := make(map[int]bool)
			for _, emp := range existingEmployees {
				existingIDs[emp.ID] = true
			}

			// Check if all requested assignee IDs exist
			var invalidIDs []int
			for _, id := range input.AssigneeIDs {
				if !existingIDs[id] {
					invalidIDs = append(invalidIDs, id)
				}
			}

			if len(invalidIDs) > 0 {
				return nil, &ServiceError{
					Status: http.StatusBadRequest,
					Msg:    "Some assignee IDs do not exist",
				}
			}

			taskUpdate.AddAssigneeIDs(input.AssigneeIDs...)
		}
	}

	_, err := taskUpdate.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Task not found",
			}
		}
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to update task",
		}
	}

	// Get the updated task with all edges
	updatedTask, err := s.Client.Task.Query().
		Where(task.ID(id)).
		WithProject().
		WithLabels().
		WithAssignees().
		WithReports().
		Only(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch updated task",
		}
	}

	return updatedTask, nil
}
