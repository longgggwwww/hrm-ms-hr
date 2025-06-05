package label

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/label"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/dto"
)

// ServiceError represents a service-level error
type ServiceError struct {
	Status int
	Msg    string
}

func (e *ServiceError) Error() string {
	return e.Msg
}

// LabelService provides business logic for label operations
type LabelService struct {
	Client *ent.Client
}

// NewLabelService creates a new label service
func NewLabelService(client *ent.Client) *LabelService {
	return &LabelService{
		Client: client,
	}
}

// encodeCursor encodes cursor data to base64
func (s *LabelService) encodeCursor(data map[string]interface{}) string {
	jsonData, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(jsonData)
}

// decodeCursor decodes base64 cursor to map data
func (s *LabelService) decodeCursor(cursor string) (map[string]interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// addTaskCountsToLabels adds task_count field to each label
func (s *LabelService) addTaskCountsToLabels(ctx context.Context, labels []*ent.Label) ([]dto.LabelResponse, error) {
	result := make([]dto.LabelResponse, len(labels))

	for i, labelEntity := range labels {
		// Count tasks for this label
		taskCount, err := s.Client.Task.Query().
			Where(task.HasLabelsWith(label.IDEQ(labelEntity.ID))).
			Count(ctx)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to count tasks for label",
			}
		}

		// Convert label to response format
		result[i] = dto.LabelResponse{
			ID:          labelEntity.ID,
			Name:        labelEntity.Name,
			Description: labelEntity.Description,
			Color:       labelEntity.Color,
			OrgID:       labelEntity.OrgID,
			CreatedAt:   labelEntity.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   labelEntity.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			TaskCount:   taskCount,
		}
	}

	return result, nil
}
