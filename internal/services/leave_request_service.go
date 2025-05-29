package services

import (
	"context"
	"net/http"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/leaveapproval"
	"github.com/longgggwwww/hrm-ms-hr/ent/leaverequest"
)

// DTO cho tạo mới đơn nghỉ phép
type LeaveRequestCreateDTO struct {
	TotalDays  float64
	StartAt    time.Time
	EndAt      time.Time
	Reason     string
	Type       string
	EmployeeID int
	OrgID      int
}

// Service xử lý logic tạo đơn nghỉ phép
func Create(ctx context.Context, client *ent.Client, dto LeaveRequestCreateDTO) (*ent.LeaveRequest, error) {
	leaveType := leaverequest.Type(dto.Type)
	leave, err := client.LeaveRequest.Create().
		SetTotalDays(dto.TotalDays).
		SetStartAt(dto.StartAt).
		SetEndAt(dto.EndAt).
		SetReason(dto.Reason).
		SetType(leaveType).
		SetEmployeeID(dto.EmployeeID).
		SetOrgID(dto.OrgID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// Eager load applicant và leave_approves
	return client.LeaveRequest.Query().
		Where(leaverequest.ID(leave.ID)).
		WithApplicant().
		WithLeaveApproves().
		Only(ctx)
}

// Service xử lý logic lấy 1 đơn nghỉ phép
func GetLeaveRequest(ctx context.Context, client *ent.Client, id int) (*ent.LeaveRequest, error) {
	return client.LeaveRequest.Query().
		Where(leaverequest.ID(id)).
		WithApplicant().
		WithLeaveApproves().
		Only(ctx)
}

// Service lấy danh sách đơn nghỉ phép với filter, phân trang offset
func ListLeaveRequests(ctx context.Context, client *ent.Client, filter map[string]interface{}, page, limit int, orderBy, orderDir string) ([]*ent.LeaveRequest, int, error) {
	query := client.LeaveRequest.Query().
		WithApplicant().
		WithLeaveApproves()
	if status, ok := filter["status"].(string); ok && status != "" {
		query = query.Where(leaverequest.StatusEQ(leaverequest.Status(status)))
	}
	if employeeID, ok := filter["employee_id"].(int); ok && employeeID > 0 {
		query = query.Where(leaverequest.EmployeeIDEQ(employeeID))
	}
	if orgID, ok := filter["org_id"].(int); ok && orgID > 0 {
		query = query.Where(leaverequest.OrgIDEQ(orgID))
	}
	switch orderBy {
	case "id":
		if orderDir == "asc" {
			query = query.Order(leaverequest.ByID())
		} else {
			query = query.Order(leaverequest.ByID(sql.OrderDesc()))
		}
	case "created_at":
		if orderDir == "desc" {
			query = query.Order(leaverequest.ByCreatedAt(sql.OrderDesc()))
		} else {
			query = query.Order(leaverequest.ByCreatedAt())
		}
	case "updated_at":
		if orderDir == "desc" {
			query = query.Order(leaverequest.ByUpdatedAt(sql.OrderDesc()))
		} else {
			query = query.Order(leaverequest.ByUpdatedAt())
		}
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)
	list, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}
	total, err := client.LeaveRequest.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Service duyệt đơn nghỉ phép, tạo LeaveApproval nếu chưa có
func ApproveLeaveRequest(ctx context.Context, client *ent.Client, id int, reviewerID int) (*ent.LeaveRequest, error) {
	// Check đã có LeaveApproval cho leave_request này chưa
	exists, err := client.LeaveApproval.Query().Where(
		leaveapproval.LeaveRequestID(id),
	).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &ServiceError{Status: http.StatusBadRequest, Msg: "#1 ApproveLeaveRequest: Leave request already approved/rejected"}
	}
	// Tạo LeaveApproval
	_, err = client.LeaveApproval.Create().
		SetLeaveRequestID(id).
		SetReviewerID(reviewerID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// Cập nhật trạng thái đơn nghỉ phép
	leave, err := client.LeaveRequest.UpdateOneID(id).
		SetStatus(leaverequest.StatusApproved).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// Eager load applicant và leave_approves
	return client.LeaveRequest.Query().
		Where(leaverequest.ID(leave.ID)).
		WithApplicant().
		WithLeaveApproves().
		Only(ctx)
}

// Service huỷ đơn nghỉ phép, tạo LeaveApproval nếu chưa có
func Reject(ctx context.Context, client *ent.Client, id int, reviewerID int) (*ent.LeaveRequest, error) {
	exists, err := client.LeaveApproval.Query().Where(
		leaveapproval.LeaveRequestID(id),
	).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &ServiceError{Status: http.StatusBadRequest, Msg: "#1 Reject: Leave request already approved/rejected"}
	}
	_, err = client.LeaveApproval.Create().
		SetLeaveRequestID(id).
		SetReviewerID(reviewerID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	leave, err := client.LeaveRequest.UpdateOneID(id).
		SetStatus(leaverequest.StatusRejected).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// Eager load applicant và leave_approves
	return client.LeaveRequest.Query().
		Where(leaverequest.ID(leave.ID)).
		WithApplicant().
		WithLeaveApproves().
		Only(ctx)
}

// Service huỷ đơn nghỉ phép bởi employee (chỉ khi đúng employee và trạng thái pending)
func RejectByEmployee(ctx context.Context, client *ent.Client, id int, employeeID int) (*ent.LeaveRequest, error) {
	leave, err := client.LeaveRequest.Query().Where(leaverequest.ID(id), leaverequest.EmployeeIDEQ(employeeID)).Only(ctx)
	if err != nil {
		return nil, err
	}
	if leave.Status != leaverequest.StatusPending {
		return nil, &ServiceError{Status: http.StatusBadRequest, Msg: "#1 RejectByEmployee: Leave request is not pending, cannot reject"}
	}
	return client.LeaveRequest.UpdateOneID(leave.ID).
		SetStatus(leaverequest.StatusRejected).
		Save(ctx)
}

// ServiceError dùng để trả lỗi kèm http status cho handler
// Handler sẽ dùng utils.RespondWithError để trả về client

type ServiceError struct {
	Status int
	Msg    string
}

func (e *ServiceError) Error() string {
	return e.Msg
}
