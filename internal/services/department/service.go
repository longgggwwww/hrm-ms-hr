package department

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/position"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// ServiceError represents a service-level error
type ServiceError struct {
	Status int
	Msg    string
}

func (e *ServiceError) Error() string {
	return e.Msg
}

// DepartmentService provides business logic for department operations
type DepartmentService struct {
	Client *ent.Client
}

// NewDepartmentService creates a new department service
func NewDepartmentService(client *ent.Client) *DepartmentService {
	return &DepartmentService{
		Client: client,
	}
}

// encodeCursor encodes cursor data to base64
func (s *DepartmentService) encodeCursor(data map[string]interface{}) string {
	jsonData, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(jsonData)
}

// decodeCursor decodes base64 cursor to map data
func (s *DepartmentService) decodeCursor(cursor string) (map[string]interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// buildDepartmentResponse creates a DepartmentResponse from an ent.Department
func (s *DepartmentService) buildDepartmentResponse(dept *ent.Department, positionCount int) dtos.DepartmentResponse {
	return dtos.DepartmentResponse{
		ID:            dept.ID,
		Name:          dept.Name,
		Code:          dept.Code,
		OrgID:         dept.OrgID,
		CreatedAt:     dept.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     dept.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		PositionCount: positionCount,
		ZaloGID:       dept.ZaloGid,
	}
}

// getPositionCount gets the count of positions for a department
func (s *DepartmentService) getPositionCount(ctx context.Context, departmentID int) (int, error) {
	count, err := s.Client.Position.Query().
		Where(position.DepartmentID(departmentID)).
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}
