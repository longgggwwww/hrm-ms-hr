package department

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/position"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Delete deletes a department by ID
func (s *DepartmentService) Delete(ctx context.Context, id int) error {
	// Check if department exists
	_, err := s.Client.Department.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Department not found",
			}
		}
		return &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to get department",
		}
	}

	// Check if department has any positions
	hasPositions, err := s.Client.Position.Query().
		Where(position.DepartmentID(id)).
		Exist(ctx)
	if err != nil {
		return &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to check for associated positions",
		}
	}
	if hasPositions {
		return &ServiceError{
			Status: http.StatusConflict,
			Msg:    "Cannot delete department with existing positions",
		}
	}

	// Delete the department
	err = s.Client.Department.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to delete department",
		}
	}

	return nil
}

// DeleteBulk deletes multiple departments
func (s *DepartmentService) DeleteBulk(ctx context.Context, input dtos.DepartmentDeleteBulkInput) error {
	tx, err := s.Client.Tx(ctx)
	if err != nil {
		return &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to start transaction",
		}
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var errors []string
	deletedCount := 0

	for _, id := range input.IDs {
		// Check if department exists
		dept, err := tx.Department.Get(ctx, id)
		if err != nil {
			if ent.IsNotFound(err) {
				errors = append(errors, "Department not found: "+string(rune(id)))
				continue
			}
			tx.Rollback()
			return &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to get department: " + string(rune(id)),
			}
		}

		// Check if department has any positions
		hasPositions, err := tx.Position.Query().
			Where(position.DepartmentID(id)).
			Exist(ctx)
		if err != nil {
			tx.Rollback()
			return &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to check for associated positions: " + string(rune(id)),
			}
		}
		if hasPositions {
			errors = append(errors, "Cannot delete department with existing positions: "+dept.Name)
			continue
		}

		// Delete the department
		err = tx.Department.DeleteOneID(id).Exec(ctx)
		if err != nil {
			errors = append(errors, "Failed to delete department: "+dept.Name)
			continue
		}

		deletedCount++
	}

	if err := tx.Commit(); err != nil {
		return &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to commit transaction",
		}
	}

	// If there were any errors, return them
	if len(errors) > 0 {
		return &ServiceError{
			Status: http.StatusPartialContent,
			Msg:    "Some departments could not be deleted",
		}
	}

	return nil
}
