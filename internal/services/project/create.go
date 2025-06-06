package project

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/project"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
)

// Create creates a new project
func (s *ProjectService) Create(ctx context.Context, orgID, employeeID int, input dtos.ProjectCreateInput) (map[string]interface{}, error) {
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

	var endAtPtr *time.Time
	if input.EndAt != nil {
		endAt, err := time.Parse(time.RFC3339, *input.EndAt)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid end_at format, must be RFC3339",
			}
		}
		endAtPtr = &endAt
	}

	// Auto-generate code if not provided
	var projectCode string
	if input.Code != nil && *input.Code != "" {
		projectCode = *input.Code

		// Check if code already exists in the organization
		existingProject, err := s.Client.Project.Query().
			Where(project.CodeEQ(projectCode)).
			Where(project.OrgIDEQ(orgID)).
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
	} else {
		// Auto-generate code: ORG{orgId}-PROJ-{sequence}
		// Find the highest sequence number for this organization
		latestProjects, err := s.Client.Project.Query().
			Where(project.OrgIDEQ(orgID)).
			Where(project.CodeContains("ORG" + strconv.Itoa(orgID) + "-PROJ-")).
			Order(project.ByCreatedAt(sql.OrderDesc())).
			Limit(1).
			All(ctx)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to generate project code",
			}
		}

		sequence := 1
		if len(latestProjects) > 0 {
			// Extract sequence number from the latest project code
			latestCode := latestProjects[0].Code
			prefix := "ORG" + strconv.Itoa(orgID) + "-PROJ-"
			if len(latestCode) > len(prefix) {
				sequenceStr := latestCode[len(prefix):]
				if seq, err := strconv.Atoi(sequenceStr); err == nil {
					sequence = seq + 1
				}
			}
		}

		projectCode = "ORG" + strconv.Itoa(orgID) + "-PROJ-" + fmt.Sprintf("%03d", sequence)

		// Double-check uniqueness (in case of concurrent requests)
		for {
			exists, err := s.Client.Project.Query().
				Where(project.CodeEQ(projectCode)).
				Where(project.OrgIDEQ(orgID)).
				Exist(ctx)
			if err != nil {
				return nil, &ServiceError{
					Status: http.StatusInternalServerError,
					Msg:    "Failed to validate generated project code",
				}
			}
			if !exists {
				break
			}
			sequence++
			projectCode = "ORG" + strconv.Itoa(orgID) + "-PROJ-" + fmt.Sprintf("%03d", sequence)
		}
	}

	// Prepare all member IDs (current employee + provided member IDs)
	allMemberIDs := []int{}

	// Add current employee ID from token
	if employeeID > 0 {
		allMemberIDs = append(allMemberIDs, employeeID)
	}

	// Add provided member IDs (avoid duplicates)
	memberIDSet := make(map[int]bool)
	memberIDSet[employeeID] = true // Mark current employee as already added

	for _, memberID := range input.MemberIDs {
		if !memberIDSet[memberID] {
			allMemberIDs = append(allMemberIDs, memberID)
			memberIDSet[memberID] = true
		}
	}

	// Validate all member IDs if there are any
	if len(allMemberIDs) > 0 {
		// Check if all member IDs exist and belong to the same organization
		memberCount, err := s.Client.Employee.Query().
			Where(employee.IDIn(allMemberIDs...)).
			Where(employee.OrgID(orgID)).
			Count(ctx)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to validate member IDs",
			}
		}
		if memberCount != len(allMemberIDs) {
			return nil, &ServiceError{
				Status: http.StatusBadRequest,
				Msg:    "One or more member IDs are invalid or do not belong to your organization",
			}
		}
	}

	projectCreate := s.Client.Project.Create().
		SetName(input.Name).
		SetCode(projectCode).
		SetNillableDescription(input.Description).
		SetNillableStartAt(startAtPtr).
		SetNillableEndAt(endAtPtr).
		SetCreatorID(employeeID).
		SetUpdaterID(employeeID).
		SetOrgID(orgID)

	// Add member IDs if there are any
	if len(allMemberIDs) > 0 {
		projectCreate = projectCreate.AddMemberIDs(allMemberIDs...)
	}

	row, err := projectCreate.Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to create project",
		}
	}

	// Get the created project with all edges
	createdProject, err := s.Client.Project.Query().
		Where(project.ID(row.ID)).
		WithTasks(func(q *ent.TaskQuery) {
			q.WithAssignees() // Load assignees for tasks
		}).
		WithOrganization().
		WithCreator().
		WithUpdater().
		WithMembers().
		Only(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch created project",
		}
	}

	// Collect user IDs from the project
	projects := []*ent.Project{createdProject}
	userIDs := s.collectUserIDsFromProjects(projects)

	// Fetch user information
	userMap, err := s.getUserInfoMap(userIDs)
	if err != nil {
		// Log error but continue without user enrichment
		userMap = make(map[int32]*grpc_clients.User)
	}

	// Enrich project with user information and tasks
	enrichedProject := s.enrichProjectWithUserInfoForGet(createdProject, userMap)

	return enrichedProject, nil
}
