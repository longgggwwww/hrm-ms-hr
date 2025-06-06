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

// RemoveMembers removes members from an existing project
func (s *ProjectService) RemoveMembers(ctx context.Context, id, employeeID, orgID int, input dtos.ProjectRemoveMembersInput) (map[string]interface{}, error) {
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
			Msg:    "Only the project creator can remove members from this project",
		}
	}

	// Get current project members
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

	// Filter members to remove - only remove members that are actually in the project
	var membersToRemove []int
	var notFoundMembers []int
	for _, memberID := range input.MemberIDs {
		if currentMemberIDs[memberID] {
			membersToRemove = append(membersToRemove, memberID)
		} else {
			notFoundMembers = append(notFoundMembers, memberID)
		}
	}

	// If some members are not in the project, return an error
	if len(notFoundMembers) > 0 {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "Some specified members are not currently in this project",
		}
	}

	// If no members to remove, return current project
	if len(membersToRemove) == 0 {
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

	// Remove members from the project
	projectUpdate := s.Client.Project.UpdateOneID(id)
	projectUpdate.RemoveMemberIDs(membersToRemove...)

	_, err = projectUpdate.Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to remove members from project",
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
