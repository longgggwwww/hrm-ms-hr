package department

import (
	"context"
	"net/http"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/department"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
)

// Create creates a new department
func (s *DepartmentService) Create(ctx context.Context, orgID int, input dtos.DepartmentCreateInput) (*ent.Department, error) {
	// Check if department code already exists in the organization
	exists, err := s.Client.Department.Query().
		Where(
			department.Code(input.Code),
			department.OrgID(orgID),
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

	// Create the department
	dept, err := s.Client.Department.Create().
		SetName(input.Name).
		SetCode(input.Code).
		SetOrgID(orgID).
		Save(ctx)
	if err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to create department",
		}
	}

	return dept, nil
}

// CreateBulk creates multiple departments
func (s *DepartmentService) CreateBulk(ctx context.Context, orgID int, input dtos.DepartmentBulkCreateInput) ([]dtos.DepartmentResponse, error) {
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

	var responses []dtos.DepartmentResponse

	for _, deptInput := range input.Departments {
		// Check if department code already exists in the organization
		exists, err := tx.Department.Query().
			Where(
				department.Code(deptInput.Code),
				department.OrgID(orgID),
			).
			Exist(ctx)
		if err != nil {
			tx.Rollback()
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to check department code uniqueness for: " + deptInput.Code,
			}
		}
		if exists {
			tx.Rollback()
			return nil, &ServiceError{
				Status: http.StatusConflict,
				Msg:    "Department code already exists in this organization: " + deptInput.Code,
			}
		}

		// Create the department
		dept, err := tx.Department.Create().
			SetName(deptInput.Name).
			SetCode(deptInput.Code).
			SetOrgID(orgID).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, &ServiceError{
				Status: http.StatusInternalServerError,
				Msg:    "Failed to create department: " + deptInput.Name,
			}
		}

		// Get position count (will be 0 for newly created departments)
		positionCount, _ := s.getPositionCount(ctx, dept.ID)
		responses = append(responses, s.buildDepartmentResponse(dept, positionCount))
	}

	if err := tx.Commit(); err != nil {
		return nil, &ServiceError{
			Status: http.StatusInternalServerError,
			Msg:    "Failed to commit transaction",
		}
	}

	return responses, nil
}
