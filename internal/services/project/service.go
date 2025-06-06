package project

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
)

// ServiceError represents a service-level error
type ServiceError struct {
	Status int
	Msg    string
}

func (e *ServiceError) Error() string {
	return e.Msg
}

// ProjectService provides business logic for project operations
type ProjectService struct {
	Client     *ent.Client
	UserClient grpc_clients.UserServiceClient
}

// NewProjectService creates a new project service
func NewProjectService(client *ent.Client, userClient grpc_clients.UserServiceClient) *ProjectService {
	return &ProjectService{
		Client:     client,
		UserClient: userClient,
	}
}

// encodeCursor encodes cursor data to base64
func (s *ProjectService) encodeCursor(data map[string]interface{}) string {
	jsonData, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(jsonData)
}

// decodeCursor decodes base64 cursor to map data
func (s *ProjectService) decodeCursor(cursor string) (map[string]interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// getUserInfoMap fetches user information by user IDs and returns a map for quick lookup
func (s *ProjectService) getUserInfoMap(userIDs []int32) (map[int32]*grpc_clients.User, error) {
	if len(userIDs) == 0 || s.UserClient == nil {
		return make(map[int32]*grpc_clients.User), nil
	}

	// Use GetUsersByIDs for bulk fetch
	resp, err := s.UserClient.GetUsersByIDs(context.Background(), &grpc_clients.GetUsersByIDsRequest{
		Ids: userIDs,
	})
	if err != nil {
		return make(map[int32]*grpc_clients.User), nil
	}

	userMap := make(map[int32]*grpc_clients.User)
	if resp != nil && resp.Users != nil {
		for _, user := range resp.Users {
			userMap[user.Id] = user
		}
	}

	return userMap, nil
}

// normalizeUserInfo extracts user information from gRPC User struct
func (s *ProjectService) normalizeUserInfo(user *grpc_clients.User) map[string]interface{} {
	userInfo := map[string]interface{}{
		"id": user.Id,
	}
	if user.Email != nil && user.Email.Value != "" {
		userInfo["email"] = user.Email.Value
	}
	if user.Avatar != nil && user.Avatar.Value != "" {
		userInfo["avatar"] = user.Avatar.Value
	}
	if user.FirstName != "" {
		userInfo["first_name"] = user.FirstName
	}
	if user.LastName != "" {
		userInfo["last_name"] = user.LastName
	}
	if user.Phone != nil && user.Phone.Value != "" {
		userInfo["phone"] = user.Phone.Value
	}

	return userInfo
}

// collectUserIDsFromProjects collects all user IDs from project creators, updaters, and members
func (s *ProjectService) collectUserIDsFromProjects(projects []*ent.Project) []int32 {
	userIDSet := make(map[int32]bool)
	var userIDs []int32

	for _, proj := range projects {
		// Creator user ID
		if proj.Edges.Creator != nil && proj.Edges.Creator.UserID != "" {
			if userID, err := strconv.Atoi(proj.Edges.Creator.UserID); err == nil {
				if !userIDSet[int32(userID)] {
					userIDSet[int32(userID)] = true
					userIDs = append(userIDs, int32(userID))
				}
			}
		}

		// Updater user ID
		if proj.Edges.Updater != nil && proj.Edges.Updater.UserID != "" {
			if userID, err := strconv.Atoi(proj.Edges.Updater.UserID); err == nil {
				if !userIDSet[int32(userID)] {
					userIDSet[int32(userID)] = true
					userIDs = append(userIDs, int32(userID))
				}
			}
		}

		// Members user IDs
		if proj.Edges.Members != nil {
			for _, member := range proj.Edges.Members {
				if member.UserID != "" {
					if userID, err := strconv.Atoi(member.UserID); err == nil {
						if !userIDSet[int32(userID)] {
							userIDSet[int32(userID)] = true
							userIDs = append(userIDs, int32(userID))
						}
					}
				}
			}
		}
	}

	return userIDs
}

// addTaskCountsToProjects adds task_count field to each project for list operations
func (s *ProjectService) addTaskCountsToProjects(ctx context.Context, projects []*ent.Project) ([]dtos.ProjectResponse, error) {
	result := make([]dtos.ProjectResponse, len(projects))

	// Get task counts for all projects
	taskCounts := make(map[int]int)
	if len(projects) > 0 {
		for _, proj := range projects {
			count, err := s.Client.Task.Query().
				Where(task.ProjectIDEQ(proj.ID)).
				Count(ctx)
			if err == nil {
				taskCounts[proj.ID] = count
			}
		}
	}

	// Collect user IDs from all projects
	userIDs := s.collectUserIDsFromProjects(projects)

	// Fetch user information
	userMap, err := s.getUserInfoMap(userIDs)
	if err != nil {
		// Log error but continue without user enrichment
		userMap = make(map[int32]*grpc_clients.User)
	}

	// Enrich projects with user information
	for i, proj := range projects {
		taskCount := taskCounts[proj.ID] // Default to 0 if not found
		result[i] = s.enrichProjectWithUserInfo(proj, userMap, taskCount)
	}

	return result, nil
}

// enrichProjectWithUserInfo enriches a single project with user information for list operations
func (s *ProjectService) enrichProjectWithUserInfo(proj *ent.Project, userMap map[int32]*grpc_clients.User, taskCount int) dtos.ProjectResponse {
	response := dtos.ProjectResponse{
		ID:          proj.ID,
		Name:        proj.Name,
		Code:        proj.Code,
		Description: proj.Description, // Convert to pointer
		StartAt:     proj.StartAt,     // Convert to pointer
		EndAt:       proj.EndAt,       // Convert to pointer
		CreatorID:   proj.CreatorID,
		UpdaterID:   proj.UpdaterID,
		OrgID:       proj.OrgID,
		Process:     proj.Process,
		Status:      string(proj.Status),
		CreatedAt:   proj.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   proj.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		TaskCount:   taskCount,
	}

	// Create edges structure preserving original structure
	edges := make(map[string]interface{})

	// Add organization edge (unchanged)
	if proj.Edges.Organization != nil {
		edges["organization"] = proj.Edges.Organization
	}

	// Enrich creator with user_info while preserving original structure
	if proj.Edges.Creator != nil {
		creatorData := map[string]interface{}{
			"id":          proj.Edges.Creator.ID,
			"user_id":     proj.Edges.Creator.UserID,
			"code":        proj.Edges.Creator.Code,
			"position_id": proj.Edges.Creator.PositionID,
			"org_id":      proj.Edges.Creator.OrgID,
			"joining_at":  proj.Edges.Creator.JoiningAt,
			"status":      proj.Edges.Creator.Status,
			"created_at":  proj.Edges.Creator.CreatedAt,
			"updated_at":  proj.Edges.Creator.UpdatedAt,
		}

		// Add user_info if available
		if proj.Edges.Creator.UserID != "" {
			if userIDInt, err := strconv.Atoi(proj.Edges.Creator.UserID); err == nil {
				if userInfo, exists := userMap[int32(userIDInt)]; exists {
					creatorData["user_info"] = s.normalizeUserInfo(userInfo)
				}
			}
		}
		edges["creator"] = creatorData
	}

	// Enrich updater with user_info while preserving original structure
	if proj.Edges.Updater != nil {
		updaterData := map[string]interface{}{
			"id":          proj.Edges.Updater.ID,
			"user_id":     proj.Edges.Updater.UserID,
			"code":        proj.Edges.Updater.Code,
			"position_id": proj.Edges.Updater.PositionID,
			"org_id":      proj.Edges.Updater.OrgID,
			"joining_at":  proj.Edges.Updater.JoiningAt,
			"status":      proj.Edges.Updater.Status,
			"created_at":  proj.Edges.Updater.CreatedAt,
			"updated_at":  proj.Edges.Updater.UpdatedAt,
		}

		// Add user_info if available
		if proj.Edges.Updater.UserID != "" {
			if userIDInt, err := strconv.Atoi(proj.Edges.Updater.UserID); err == nil {
				if userInfo, exists := userMap[int32(userIDInt)]; exists {
					updaterData["user_info"] = s.normalizeUserInfo(userInfo)
				}
			}
		}
		edges["updater"] = updaterData
	}

	// Enrich members with user_info while preserving original structure
	if proj.Edges.Members != nil {
		var membersData []map[string]interface{}
		for _, member := range proj.Edges.Members {
			memberData := map[string]interface{}{
				"id":          member.ID,
				"user_id":     member.UserID,
				"code":        member.Code,
				"position_id": member.PositionID,
				"org_id":      member.OrgID,
				"joining_at":  member.JoiningAt,
				"status":      member.Status,
				"created_at":  member.CreatedAt,
				"updated_at":  member.UpdatedAt,
			}

			// Add position edge data if available
			if member.Edges.Position != nil {
				memberData["position"] = map[string]interface{}{
					"id":            member.Edges.Position.ID,
					"name":          member.Edges.Position.Name,
					"code":          member.Edges.Position.Code,
					"department_id": member.Edges.Position.DepartmentID,
					"created_at":    member.Edges.Position.CreatedAt,
					"updated_at":    member.Edges.Position.UpdatedAt,
				}
			}

			// Add user_info if available
			if member.UserID != "" {
				if userIDInt, err := strconv.Atoi(member.UserID); err == nil {
					if userInfo, exists := userMap[int32(userIDInt)]; exists {
						memberData["user_info"] = s.normalizeUserInfo(userInfo)
					}
				}
			}
			membersData = append(membersData, memberData)
		}
		edges["members"] = membersData
	}

	response.Edges = edges
	return response
}

// enrichProjectWithUserInfoForGet enriches a single project with user information and tasks (for Get method)
func (s *ProjectService) enrichProjectWithUserInfoForGet(proj *ent.Project, userMap map[int32]*grpc_clients.User) map[string]interface{} {
	result := map[string]interface{}{
		"id":          proj.ID,
		"name":        proj.Name,
		"code":        proj.Code,
		"description": proj.Description,
		"start_at":    proj.StartAt,
		"end_at":      proj.EndAt,
		"creator_id":  proj.CreatorID,
		"updater_id":  proj.UpdaterID,
		"org_id":      proj.OrgID,
		"process":     proj.Process,
		"status":      proj.Status,
		"created_at":  proj.CreatedAt,
		"updated_at":  proj.UpdatedAt,
	}

	// Create edges structure preserving original structure
	edges := make(map[string]interface{})

	// Add tasks array instead of task_count (for Get method)
	if proj.Edges.Tasks != nil {
		edges["tasks"] = proj.Edges.Tasks
	}

	// Add organization edge (unchanged)
	if proj.Edges.Organization != nil {
		edges["organization"] = proj.Edges.Organization
	}

	// Enrich creator with user_info while preserving original structure
	if proj.Edges.Creator != nil {
		creatorData := map[string]interface{}{
			"id":          proj.Edges.Creator.ID,
			"user_id":     proj.Edges.Creator.UserID,
			"code":        proj.Edges.Creator.Code,
			"position_id": proj.Edges.Creator.PositionID,
			"org_id":      proj.Edges.Creator.OrgID,
			"joining_at":  proj.Edges.Creator.JoiningAt,
			"status":      proj.Edges.Creator.Status,
			"created_at":  proj.Edges.Creator.CreatedAt,
			"updated_at":  proj.Edges.Creator.UpdatedAt,
		}

		// Add user_info if available
		if proj.Edges.Creator.UserID != "" {
			if userIDInt, err := strconv.Atoi(proj.Edges.Creator.UserID); err == nil {
				if userInfo, exists := userMap[int32(userIDInt)]; exists {
					creatorData["user_info"] = s.normalizeUserInfo(userInfo)
				}
			}
		}
		edges["creator"] = creatorData
	}

	// Enrich updater with user_info while preserving original structure
	if proj.Edges.Updater != nil {
		updaterData := map[string]interface{}{
			"id":          proj.Edges.Updater.ID,
			"user_id":     proj.Edges.Updater.UserID,
			"code":        proj.Edges.Updater.Code,
			"position_id": proj.Edges.Updater.PositionID,
			"org_id":      proj.Edges.Updater.OrgID,
			"joining_at":  proj.Edges.Updater.JoiningAt,
			"status":      proj.Edges.Updater.Status,
			"created_at":  proj.Edges.Updater.CreatedAt,
			"updated_at":  proj.Edges.Updater.UpdatedAt,
		}

		// Add user_info if available
		if proj.Edges.Updater.UserID != "" {
			if userIDInt, err := strconv.Atoi(proj.Edges.Updater.UserID); err == nil {
				if userInfo, exists := userMap[int32(userIDInt)]; exists {
					updaterData["user_info"] = s.normalizeUserInfo(userInfo)
				}
			}
		}
		edges["updater"] = updaterData
	}

	// Enrich members with user_info while preserving original structure
	if proj.Edges.Members != nil {
		var membersData []map[string]interface{}
		for _, member := range proj.Edges.Members {
			memberData := map[string]interface{}{
				"id":          member.ID,
				"user_id":     member.UserID,
				"code":        member.Code,
				"position_id": member.PositionID,
				"org_id":      member.OrgID,
				"joining_at":  member.JoiningAt,
				"status":      member.Status,
				"created_at":  member.CreatedAt,
				"updated_at":  member.UpdatedAt,
			}

			// Add position edge data if available
			if member.Edges.Position != nil {
				memberData["position"] = map[string]interface{}{
					"id":            member.Edges.Position.ID,
					"name":          member.Edges.Position.Name,
					"code":          member.Edges.Position.Code,
					"department_id": member.Edges.Position.DepartmentID,
					"created_at":    member.Edges.Position.CreatedAt,
					"updated_at":    member.Edges.Position.UpdatedAt,
				}
			}

			// Add user_info if available
			if member.UserID != "" {
				if userIDInt, err := strconv.Atoi(member.UserID); err == nil {
					if userInfo, exists := userMap[int32(userIDInt)]; exists {
						memberData["user_info"] = s.normalizeUserInfo(userInfo)
					}
				}
			}
			membersData = append(membersData, memberData)
		}
		edges["members"] = membersData
	}

	// Add the edges structure to the result
	result["edges"] = edges

	return result
}
