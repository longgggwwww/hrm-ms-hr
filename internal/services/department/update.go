package department

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/department"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Update updates an existing department
func (s *DepartmentService) Update(ctx context.Context, id int, input dtos.DepartmentUpdateInput) (*ent.Department, error) {
	// Check if department exists
	dept, err := s.Client.Department.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ServiceError{
				Status: http.StatusNotFound,
				Msg:    "Department not found",
			}
		}
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to get department",
		}
	}

	// Check if code is being updated and ensure uniqueness within organization
	if input.Code != nil && *input.Code != dept.Code {
		exists, err := s.Client.Department.Query().
			Where(
				department.Code(*input.Code),
				department.OrgID(dept.OrgID),
				department.IDNEQ(id),
			).
			Exist(ctx)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to check department code uniqueness",
			}
		}
		if exists {
			return nil, &ServiceError{
				Status: http.StatusConflict,
				Msg:    "Department code already exists in this organization",
			}
		}
	}

	// Start transaction if zalo_gid is being updated
	if input.ZaloGID != nil {
		tx, err := s.Client.Tx(ctx)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to start transaction",
			}
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Build update query for department
		update := tx.Department.UpdateOneID(id)

		// Apply department updates
		if input.Name != nil {
			update = update.SetName(*input.Name)
		}
		if input.Code != nil {
			update = update.SetCode(*input.Code)
		}
		if input.ZaloGID != nil {
			if *input.ZaloGID == "" {
				update = update.ClearZaloGid()
			} else {
				update = update.SetZaloGid(*input.ZaloGID)
			}
		}

		// Save the updated department
		updatedDept, err := update.Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to update department",
			}
		}

		if err := tx.Commit(); err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to commit transaction",
			}
		}

		// Reload with edges
		updatedDept, err = s.Client.Department.Query().
			Where(department.ID(id)).
			Only(ctx)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to reload department with edges",
			}
		}

		return updatedDept, nil
	}

	// Update department without zalo_gid changes
	update := s.Client.Department.UpdateOneID(id)

	// Apply updates
	if input.Name != nil {
		update = update.SetName(*input.Name)
	}
	if input.Code != nil {
		update = update.SetCode(*input.Code)
	}

	// Save the updated department
	updatedDept, err := update.Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to update department",
		}
	}

	return updatedDept, nil
}
