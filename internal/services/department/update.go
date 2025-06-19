package department

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/department"
	"github.com/longgggwwww/hrm-ms-hr/ent/zalodepartment"
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

	// Start transaction if group_id is being updated
	if input.GroupID != nil {
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

		// Save the updated department
		updatedDept, err := update.Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to update department",
			}
		}

		// Handle zalo_department update
		if *input.GroupID == "" {
			// Remove zalo_department if group_id is empty
			_, err = tx.ZaloDepartment.Delete().
				Where(zalodepartment.DepartmentID(id)).
				Exec(ctx)
			if err != nil {
				tx.Rollback()
				return nil, &ServiceError{
					Status: http.StatusInternalServerError,
					Msg:    "Failed to remove zalo department mapping",
				}
			}
		} else {
			// Check if zalo_department exists
			exists, err := tx.ZaloDepartment.Query().
				Where(zalodepartment.DepartmentID(id)).
				Exist(ctx)
			if err != nil {
				tx.Rollback()
				return nil, &ServiceError{
					Status: http.StatusInternalServerError,
					Msg:    "Failed to check zalo department existence",
				}
			}

			if exists {
				// Update existing zalo_department
				_, err = tx.ZaloDepartment.Update().
					Where(zalodepartment.DepartmentID(id)).
					SetGroupID(*input.GroupID).
					Save(ctx)
				if err != nil {
					tx.Rollback()
					return nil, &ServiceError{
						Status: http.StatusInternalServerError,
						Msg:    "Failed to update zalo department mapping",
					}
				}
			} else {
				// Create new zalo_department
				_, err = tx.ZaloDepartment.Create().
					SetGroupID(*input.GroupID).
					SetDepartmentID(id).
					Save(ctx)
				if err != nil {
					tx.Rollback()
					return nil, &ServiceError{
						Status: http.StatusInternalServerError,
						Msg:    "Failed to create zalo department mapping",
					}
				}
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
			WithZaloDepartment().
			Only(ctx)
		if err != nil {
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to reload department with edges",
			}
		}

		return updatedDept, nil
	}

	// Update department without group_id changes
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
