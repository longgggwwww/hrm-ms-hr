package services

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
)

type EmployeeService struct {
	Client     *ent.Client
	UserClient grpc_clients.UserServiceClient
}

// EmployeeListQuery gom các tham số truy vấn employee list
type EmployeeListQuery struct {
	Page     int
	Limit    int
	OrderBy  string
	OrderDir string
	OrgID    int
}

func NewEmployeeService(client *ent.Client, userClient grpc_clients.UserServiceClient) *EmployeeService {
	return &EmployeeService{
		Client:     client,
		UserClient: userClient,
	}
}

func (s *EmployeeService) Create(ctx context.Context, orgID int, input dtos.EmployeeCreateInput) (*ent.Employee, *grpc_clients.CreateUserResponse, error) {
	var joiningAt time.Time
	if input.JoiningAt != "" {
		var err error
		joiningAt, err = time.Parse(time.RFC3339, input.JoiningAt)
		if err != nil {
			return nil, nil, &ServiceError{Status: http.StatusBadRequest, Msg: "Invalid joining_at format, must be RFC3339"}
		}
	} else {
		joiningAt = time.Now()
	}

	tx, err := s.Client.Tx(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	status := employee.Status(input.Status)
	if status != employee.StatusActive && status != employee.StatusInactive {
		tx.Rollback()
		return nil, nil, &ServiceError{Status: http.StatusBadRequest, Msg: "Invalid status"}
	}

	employeeObj, err := tx.Employee.Create().
		SetCode(input.Code).
		SetPositionID(input.PositionID).
		SetJoiningAt(joiningAt).
		SetStatus(status).
		SetOrgID(orgID).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if s.UserClient == nil {
		tx.Rollback()
		return nil, nil, &ServiceError{Status: http.StatusInternalServerError, Msg: "User service unavailable"}
	}

	respb, err := s.UserClient.CreateUser(ctx, &grpc_clients.CreateUserRequest{
		FirstName: input.User.FirstName,
		LastName:  input.User.LastName,
		Email:     input.User.Email,
		Gender:    input.User.Gender,
		Phone:     input.User.Phone,
		Address:   input.User.Address,
		WardCode:  strconv.Itoa(input.User.WardCode),
		RoleIds:   input.User.RoleIds,
		PermIds:   input.User.PermIds,
		Account: &grpc_clients.Account{
			Username: input.User.Account.Username,
			Password: input.User.Account.Password,
		},
	})
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if respb != nil && respb.User != nil && respb.User.Id > 0 {
		userIDStr := strconv.FormatInt(int64(respb.User.Id), 10)
		_, err := tx.Employee.UpdateOneID(employeeObj.ID).SetUserID(userIDStr).Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return employeeObj, respb, nil
}

// List returns paginated employees, total count, and user info map
func (s *EmployeeService) List(ctx context.Context, q EmployeeListQuery) ([]*ent.Employee, int, map[int32]*grpc_clients.User, error) {
	query := s.Client.Employee.Query().
		Where(employee.OrgID(q.OrgID)).
		WithPosition(func(qp *ent.PositionQuery) {
			qp.WithDepartments()
		})

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	offset := (q.Page - 1) * q.Limit
	query = query.Offset(offset).Limit(q.Limit)

	switch q.OrderBy {
	case "id":
		if q.OrderDir == "asc" {
			query = query.Order(employee.ByID())
		} else {
			query = query.Order(employee.ByID(sql.OrderDesc()))
		}
	case "created_at":
		if q.OrderDir == "desc" {
			query = query.Order(employee.ByCreatedAt(sql.OrderDesc()))
		} else {
			query = query.Order(employee.ByCreatedAt())
		}
	case "updated_at":
		if q.OrderDir == "desc" {
			query = query.Order(employee.ByUpdatedAt(sql.OrderDesc()))
		} else {
			query = query.Order(employee.ByUpdatedAt())
		}
	}

	employees, err := query.All(ctx)
	if err != nil {
		return nil, 0, nil, err
	}
	total, err := s.Client.Employee.Query().Where(employee.OrgID(q.OrgID)).Count(ctx)
	if err != nil {
		return nil, 0, nil, err
	}

	var userIDs []int32
	for _, emp := range employees {
		if emp.UserID != "" {
			if id, err := strconv.Atoi(emp.UserID); err == nil {
				userIDs = append(userIDs, int32(id))
			}
		}
	}
	userInfoMap := make(map[int32]*grpc_clients.User)
	if s.UserClient != nil && len(userIDs) > 0 {
		resp, err := s.UserClient.GetUsersByIDs(ctx, &grpc_clients.GetUsersByIDsRequest{Ids: userIDs})
		if err == nil && resp != nil && len(resp.Users) > 0 {
			for _, u := range resp.Users {
				userInfoMap[u.Id] = u
			}
		}
	}
	return employees, total, userInfoMap, nil
}

// GetEmployeeWithUserInfo fetches a single employee by id, org, enriches with user info
func (s *EmployeeService) GetEmployeeWithUserInfo(
	ctx context.Context,
	id int,
	orgID int,
) (*ent.Employee, *grpc_clients.User, error) {
	emp, err := s.Client.Employee.Query().
		Where(employee.ID(id), employee.OrgID(orgID)).
		WithPosition(func(q *ent.PositionQuery) {
			q.WithDepartments()
		}).
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	var userInfo *grpc_clients.User
	if emp.UserID != "" && s.UserClient != nil {
		if uid, err := strconv.Atoi(emp.UserID); err == nil {
			resp, err := s.UserClient.GetUsersByIDs(
				ctx,
				&grpc_clients.GetUsersByIDsRequest{Ids: []int32{int32(uid)}},
			)
			if err == nil && resp != nil && len(resp.Users) > 0 {
				userInfo = resp.Users[0]
			}
		}
	}
	return emp, userInfo, nil
}

// DeleteById xóa employee theo id và org_id, nếu có user_id thì xóa luôn user qua userPb, tất cả trong 1 transaction. Trả về employee vừa xóa.
func (s *EmployeeService) DeleteById(ctx context.Context, id int, orgID int) (*ent.Employee, error) {
	tx, err := s.Client.Tx(ctx)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusInternalServerError, Msg: "Failed to start transaction"}
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	emp, err := tx.Employee.Query().Where(employee.ID(id), employee.OrgID(orgID)).Only(ctx)
	if err != nil {
		tx.Rollback()
		return nil, &ServiceError{Status: http.StatusNotFound, Msg: "#1 DeleteById: not found or not in your organization"}
	}
	userID := emp.UserID

	res, err := tx.Employee.Delete().Where(employee.ID(id), employee.OrgID(orgID)).Exec(ctx)
	if err != nil || res == 0 {
		tx.Rollback()
		return nil, &ServiceError{Status: http.StatusNotFound, Msg: "#2 DeleteById: not found or not in your organization"}
	}

	// Nếu có user_id thì gọi xóa user bên user service
	if userID != "" && s.UserClient != nil {
		uid, err := strconv.Atoi(userID)
		if err == nil {
			_, err := s.UserClient.DeleteUserByID(ctx, &grpc_clients.DeleteUserRequest{Id: int32(uid)})
			if err != nil {
				tx.Rollback()
				return nil, &ServiceError{Status: http.StatusInternalServerError, Msg: "#3 DeleteById: Failed to delete user in user service"}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, &ServiceError{Status: http.StatusInternalServerError, Msg: "#4 DeleteById: Failed to commit transaction"}
	}
	return emp, nil
}

func (s *EmployeeService) UpdateById(ctx context.Context, id int, orgID int, input dtos.EmployeeUpdateInput) (*ent.Employee, *grpc_clients.User, error) {
	tx, err := s.Client.Tx(ctx)
	if err != nil {
		return nil, nil, &ServiceError{Status: http.StatusInternalServerError, Msg: "#1 UpdateById: Failed to start transaction"}
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	emp, err := tx.Employee.Query().Where(employee.ID(id), employee.OrgID(orgID)).Only(ctx)
	if err != nil {
		tx.Rollback()
		return nil, nil, &ServiceError{Status: http.StatusNotFound, Msg: "#2 UpdateById: Employee not found or not in your organization"}
	}

	upd := tx.Employee.UpdateOneID(id)
	if input.Code != "" {
		upd.SetCode(input.Code)
	}
	if input.JoiningAt != "" {
		if joiningAt, err := time.Parse(time.RFC3339, input.JoiningAt); err == nil {
			upd.SetJoiningAt(joiningAt)
		}
	}
	if input.Status != "" {
		upd.SetStatus(employee.Status(input.Status))
	}
	updatedEmp, err := upd.Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, nil, &ServiceError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	var userInfo *grpc_clients.User

	if emp.UserID != "" && s.UserClient != nil {
		if uid, err := strconv.Atoi(emp.UserID); err == nil {
			userReq := &grpc_clients.UpdateUserRequest{Id: int32(uid)}
			u := input.User
			if u.FirstName != "" {
				userReq.FirstName = u.FirstName
			}
			if u.LastName != "" {
				userReq.LastName = u.LastName
			}
			if u.Email != "" {
				userReq.Email = u.Email
			}
			if u.Gender != "" {
				userReq.Gender = u.Gender
			}
			if u.Phone != "" {
				userReq.Phone = u.Phone
			}
			if u.Address != "" {
				userReq.Address = u.Address
			}
			if u.WardCode != 0 {
				userReq.WardCode = strconv.Itoa(u.WardCode)
			}
			if len(u.RoleIds) > 0 {
				userReq.RoleIds = u.RoleIds
			}
			if len(u.PermIds) > 0 {
				userReq.PermIds = u.PermIds
			}
			if u.Account.Username != "" || u.Account.Password != "" {
				userReq.Account = &grpc_clients.Account{
					Username: u.Account.Username,
					Password: u.Account.Password,
					Status:   u.Account.Status,
				}
			}
			userResp, err := s.UserClient.UpdateUser(ctx, userReq)

			if err != nil {
				tx.Rollback()
				return nil, nil, &ServiceError{Status: http.StatusInternalServerError, Msg: "#3 UpdateById: Failed to update user info: " + err.Error()}
			}
			if userResp != nil {
				userInfo = userResp.User
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, &ServiceError{Status: http.StatusInternalServerError, Msg: "#4 UpdateById: Failed to commit transaction"}
	}

	return updatedEmp, userInfo, nil
}
