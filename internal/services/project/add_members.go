package project

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/project"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
)

// AddMembers adds members to an existing project
func (s *ProjectService) AddMembers(ctx context.Context, id, employeeID, orgID int, input dtos.ProjectAddMembersInput) (map[string]interface{}, error) {
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
			Msg:    "Only the project creator can add members to this project",
		}
	}

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

	// Get current project members to avoid duplicates
	currentMembers, err := s.Client.Project.Query().
		Where(project.ID(id)).
		QueryMembers().
		Select(employee.FieldID).
		All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to get current project members",
		}
	}

	// Create map of current member IDs
	currentMemberIDs := make(map[int]bool)
	for _, member := range currentMembers {
		currentMemberIDs[member.ID] = true
	}

	// Filter out members that are already in the project
	var newMemberIDs []int
	for _, memberID := range input.MemberIDs {
		if !currentMemberIDs[memberID] {
			newMemberIDs = append(newMemberIDs, memberID)
		}
	}

	// If no new members to add, return current project
	if len(newMemberIDs) == 0 {
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
				Msg:    "Failed to fetch project",
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

	// Add new members to the project
	projectUpdate := s.Client.Project.UpdateOneID(id)
	projectUpdate.AddMemberIDs(newMemberIDs...)

	_, err = projectUpdate.Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to add members to project",
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
