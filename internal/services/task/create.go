package task

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/label"
	"github.com/longgggwwww/hrm-ms-hr/ent/project"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Create creates a new task
func (s *TaskService) Create(ctx context.Context, userID, employeeID int, input dtos.TaskCreateInput) (*ent.Task, error) {
	// Validate project membership if project_id is provided
	if input.ProjectID != nil {
		// Check if the employee is a member of the specified project
		projectExists, err := s.Client.Project.Query().
			Where(project.ID(*input.ProjectID)).
			Where(project.HasMembersWith(employee.ID(employeeID))).
			Exist(ctx)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to validate project membership",
			}
		}
		if !projectExists {
			return nil, &ServiceError{
				Status: http.StatusForbidden,
				Msg:    "You are not a member of this project",
			}
		}
	}

	// Parse start_at if provided
	var startAtPtr *time.Time
	if input.StartAt != nil {
		startAt, err := time.Parse(time.RFC3339, *input.StartAt)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid start_at format, must be RFC3339",
			}
		}
		startAtPtr = &startAt
	}

	// Parse due_date if provided
	var dueDatePtr *time.Time
	if input.DueDate != nil {
		dueDate, err := time.Parse(time.RFC3339, *input.DueDate)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid due_date format, must be RFC3339",
			}
		}
		dueDatePtr = &dueDate
	}

	// Validate and set type
	var typeVal task.Type
	if input.Type != nil {
		switch *input.Type {
		case string(task.TypeTask),
			string(task.TypeFeature),
			string(task.TypeBug),
			string(task.TypeAnother):
			typeVal = task.Type(*input.Type)
		default:
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid type value. Valid values: task, feature, bug, another",
			}
		}
	} else {
		typeVal = task.TypeTask
	}

	// Auto-generate code if not provided (GitHub-style: #x)
	var taskCode string
	if input.Code != nil && *input.Code != "" {
		taskCode = *input.Code

		// Check if code already exists globally
		exists, err := s.Client.Task.Query().
			Where(task.CodeEQ(taskCode)).
			Exist(ctx)
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

	// Create task with basic fields
	taskCreate := s.Client.Task.Create().
		SetName(input.Name).
		SetCode(taskCode).
		SetType(typeVal).
		SetNillableStartAt(startAtPtr).
		SetNillableDueDate(dueDatePtr).
		SetNillableProjectID(input.ProjectID).
		SetCreatorID(userID).
		SetUpdaterID(userID)

	// Validate and add labels if provided
	if len(input.LabelIDs) > 0 {
		// Check if all label IDs exist in the label table
		existingLabels, err := s.Client.Label.Query().
			Where(label.IDIn(input.LabelIDs...)).
			Select(label.FieldID).
			All(ctx)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to validate label IDs",
			}
		}

		// Create map of existing label IDs for validation
		existingIDs := make(map[int]bool)
		for _, lbl := range existingLabels {
			existingIDs[lbl.ID] = true
		}

		// Check if all requested label IDs exist
		var invalidIDs []int
		for _, id := range input.LabelIDs {
			if !existingIDs[id] {
				invalidIDs = append(invalidIDs, id)
			}
		}

		if len(invalidIDs) > 0 {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Some label IDs do not exist",
			}
		}

		taskCreate = taskCreate.AddLabelIDs(input.LabelIDs...)
	}

	// Validate and add assignees if provided
	if len(input.AssigneeIDs) > 0 {
		// Check if all assignee IDs exist in the employee table
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

		taskCreate = taskCreate.AddAssigneeIDs(input.AssigneeIDs...)
	}

	row, err := taskCreate.Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to create task",
		}
	}

	// Get the created task with all edges
	createdTask, err := s.Client.Task.Query().
		Where(task.ID(row.ID)).
		WithProject().
		WithLabels().
		WithAssignees().
		WithReports().
		Only(ctx)
	if err != nil {
		// Return the basic task if we can't fetch with edges
		return row, nil
	}

	return createdTask, nil
}
