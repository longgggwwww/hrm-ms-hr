package project

import (
	"context"
	"net/http"
	"time"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/project"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
)

// Update updates an existing project
func (s *ProjectService) Update(ctx context.Context, id, employeeID, orgID int, input dtos.ProjectUpdateInput) (map[string]interface{}, error) {
	// Check if the project exists and if the current user is the creator
	existingProject, err := s.Client.Project.Query().
		Where(project.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Project not found",
			}
		}
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to get project",
		}
	}

	// Check if current employee is the creator of the project
	if existingProject.CreatorID != employeeID {
		return nil, &ServiceError{
			Status: http.StatusForbidden,
			Msg:    "Only the project creator can update this project",
		}
	}

	projectUpdate := s.Client.Project.UpdateOneID(id)

	// Update fields if provided
	if input.Name != nil {
		projectUpdate.SetName(*input.Name)
	}

	if input.Code != nil {
		// Check if the new code already exists in the organization (excluding current project)
		existingProject, err := s.Client.Project.Query().
			Where(project.CodeEQ(*input.Code)).
			Where(project.OrgIDEQ(orgID)).
			Where(project.IDNEQ(id)). // Exclude current project
			First(ctx)
		if err == nil && existingProject != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Project code already exists in your organization",
			}
		}
		if err != nil && !ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to validate project code",
			}
		}
		projectUpdate.SetCode(*input.Code)
	}

	if input.Description != nil {
		projectUpdate.SetDescription(*input.Description)
	}

	if input.StartAt != nil {
		startAt, err := time.Parse(time.RFC3339, *input.StartAt)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid start_at format, must be RFC3339",
			}
		}
		projectUpdate.SetStartAt(startAt)
	}

	if input.EndAt != nil {
		endAt, err := time.Parse(time.RFC3339, *input.EndAt)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid end_at format, must be RFC3339",
			}
		}
		projectUpdate.SetEndAt(endAt)
	}

	if input.CreatorID != nil {
		projectUpdate.SetCreatorID(*input.CreatorID)
	}

	if input.UpdaterID != nil {
		projectUpdate.SetUpdaterID(*input.UpdaterID)
	}

	if input.OrgID != nil {
		projectUpdate.SetOrgID(*input.OrgID)
	}

	if input.Process != nil {
		projectUpdate.SetProcess(*input.Process)
	}

	if input.Status != nil {
		switch *input.Status {
		case string(project.StatusNotStarted), string(project.StatusInProgress), string(project.StatusCompleted):
			projectUpdate.SetStatus(project.Status(*input.Status))
		default:
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid status value",
			}
		}
	}

	// Handle member assignments
	if input.MemberIDs != nil {
		// Clear existing members first
		projectUpdate.ClearMembers()

		if len(input.MemberIDs) > 0 {
			// Validate that all member IDs exist in the employee table and belong to the same organization
			existingEmployees, err := s.Client.Employee.Query().
				Where(employee.IDIn(input.MemberIDs...)).
				Where(employee.OrgID(orgID)).
				Select(employee.FieldID).
				All(ctx)
			if err != nil {
				return nil, &ServiceError{
					Status: http.StatusInternalServerError,
					Msg:    "Failed to validate member IDs",
				}
			}

			// Create map of existing employee IDs for validation
			existingIDs := make(map[int]bool)
			for _, emp := range existingEmployees {
				existingIDs[emp.ID] = true
			}

			// Check if all requested member IDs exist
			var invalidIDs []int
			for _, id := range input.MemberIDs {
				if !existingIDs[id] {
					invalidIDs = append(invalidIDs, id)
				}
			}

			if len(invalidIDs) > 0 {
				return nil, &ServiceError{
					Status: http.StatusBadRequest,
					Msg:    "Some member IDs do not exist or do not belong to your organization",
				}
			}

			projectUpdate.AddMemberIDs(input.MemberIDs...)
		}
	}

	_, err = projectUpdate.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Project not found",
			}
		}
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to update project",
		}
	}

	// Get the updated project with all edges
	updatedProject, err := s.Client.Project.Query().
		Where(project.ID(id)).
		WithTasks().
		WithOrganization().
		WithCreator().
		WithUpdater().
		WithMembers().
		Only(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch updated project",
		}
	}

	// Collect user IDs from the project
	projects := []*ent.Project{updatedProject}
	userIDs := s.collectUserIDsFromProjects(projects)

	// Fetch user information
	userMap, err := s.getUserInfoMap(userIDs)
	if err != nil {
		// Log error but continue without user enrichment
		userMap = make(map[int32]*grpc_clients.User)
	}

	// Enrich project with user information and tasks
	enrichedProject := s.enrichProjectWithUserInfoForGet(updatedProject, userMap)

	return enrichedProject, nil
}
