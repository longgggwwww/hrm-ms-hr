package label

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/label"
	"github.com/longgggwwww/hrm-ms-hr/internal/dto"
)

// Create creates a new label
func (s *LabelService) Create(ctx context.Context, orgID int, input dto.LabelCreateInput) (*dto.LabelResponse, error) {
	create := s.Client.Label.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetColor(input.Color).
		SetOrgID(orgID)

	labelObj, err := create.Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to create label",
		}
	}

	// Fetch created label with relationships
	createdLabel, err := s.Client.Label.Query().
		Where(label.ID(labelObj.ID)).
		WithTasks().
		Only(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch created label",
		}
	}

	// Add task count to response
	labelsWithTaskCount, err := s.addTaskCountsToLabels(ctx, []*ent.Label{createdLabel})
	if err != nil {
		return nil, err
	}

	return &labelsWithTaskCount[0], nil
}

// CreateBulk creates multiple labels in bulk
func (s *LabelService) CreateBulk(ctx context.Context, orgID int, req dto.LabelBulkCreateInput) ([]*ent.Label, error) {
	if len(req.Labels) == 0 {
		return nil, &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "No labels provided",
		}
	}

	bulk := make([]*ent.LabelCreate, len(req.Labels))
	for i, labelInput := range req.Labels {
		create := s.Client.Label.Create().
			SetName(labelInput.Name).
			SetDescription(labelInput.Description).
			SetColor(labelInput.Color).
			SetOrgID(orgID)

		bulk[i] = create
	}

	labels, err := s.Client.Label.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to create labels",
		}
	}

	// Fetch created labels with relationships
	labelIDs := make([]int, len(labels))
	for i, lbl := range labels {
		labelIDs[i] = lbl.ID
	}

	createdLabels, err := s.Client.Label.Query().
		Where(label.IDIn(labelIDs...)).
		WithTasks().
		All(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to fetch created labels",
		}
	}

	return createdLabels, nil
}
