package label

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/label"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Update updates an existing label
func (s *LabelService) Update(ctx context.Context, id int, input dtos.LabelUpdateInput) (*dtos.LabelResponse, error) {
	update := s.Client.Label.UpdateOneID(id)
	update.SetNillableName(input.Name).
		SetNillableDescription(input.Description).
		SetNillableColor(input.Color)

	if input.OrgID != nil {
		update = update.SetNillableOrgID(input.OrgID)
	}

	labelObj, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Label not found",
			}
		}
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to update label",
		}
	}

	// Fetch updated label with relationships
	updatedLabel, err := s.Client.Label.Query().
		Where(label.ID(labelObj.ID)).
		WithTasks().
		Only(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch updated label",
		}
	}

	// Add task count to response
	labelsWithTaskCount, err := s.addTaskCountsToLabels(ctx, []*ent.Label{updatedLabel})
	if err != nil {
		return nil, err
	}

	return &labelsWithTaskCount[0], nil
}
