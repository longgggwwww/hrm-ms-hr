package label

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent/label"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Delete deletes a label by ID
func (s *LabelService) Delete(ctx context.Context, id int) error {
	_, err := s.Client.Label.Delete().Where(label.ID(id)).Exec(ctx)
	if err != nil {
		return &ServiceError{
			Status: http.StatusNotFound,
			Msg:    "Label not found",
		}
	}
	return nil
}

// DeleteBulk deletes multiple labels by IDs
func (s *LabelService) DeleteBulk(ctx context.Context, req dtos.LabelDeleteBulkInput) error {
	if len(req.IDs) == 0 {
		return &ServiceError{
			Status: http.StatusBadRequest,
			Msg:    "No IDs provided",
		}
	}

	_, err := s.Client.Label.Delete().
		Where(label.IDIn(req.IDs...)).
		Exec(ctx)
	if err != nil {
		return &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    err.Error(),
		}
	}

	return nil
}
