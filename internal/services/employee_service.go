package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	userPb "github.com/huynhthanhthao/hrm_user_service/proto/user"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
)

type EmployeeService struct {
	Client     *ent.Client
	UserClient userPb.UserServiceClient
}

type EmployeeCreateInput struct {
	Code       string   `json:"code" binding:"required"`
	FirstName  string   `json:"first_name" binding:"required"`
	LastName   string   `json:"last_name" binding:"required"`
	Gender     string   `json:"gender"`
	Phone      string   `json:"phone" binding:"required"`
	Email      string   `json:"email" binding:"required,email"`
	Address    string   `json:"address"`
	WardCode   int      `json:"ward_code"`
	AvatarURL  string   `json:"avatar_url"`
	PositionID int      `json:"position_id" binding:"required"`
	JoiningAt  string   `json:"joining_at"`
	Status     string   `json:"status" binding:"required"`
	Username   string   `json:"username" binding:"required"`
	Password   string   `json:"password" binding:"required"`
	RoleIds    []string `json:"role_ids"`
	PermIds    []string `json:"perm_ids"`
}

// EmployeeListQuery gom các tham số truy vấn employee list
type EmployeeListQuery struct {
	Page     int
	Limit    int
	OrderBy  string
	OrderDir string
	OrgID    int
}

func NewEmployeeService(client *ent.Client, userClient userPb.UserServiceClient) *EmployeeService {
	return &EmployeeService{
		Client:     client,
		UserClient: userClient,
	}
}

func (s *EmployeeService) Create(ctx context.Context, orgID int, input EmployeeCreateInput) (*ent.Employee, *userPb.CreateUserResponse, error) {
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
	fmt.Println(input.PositionID, 999)
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

	respb, err := s.UserClient.CreateUser(ctx, &userPb.CreateUserRequest{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Gender:    input.Gender,
		Phone:     input.Phone,
		Address:   input.Address,
		WardCode:  strconv.Itoa(input.WardCode),
		RoleIds:   input.RoleIds,
		PermIds:   input.PermIds,
		Account: &userPb.Account{
			Username: input.Username,
			Password: input.Password,
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
func (s *EmployeeService) List(ctx context.Context, q EmployeeListQuery) ([]*ent.Employee, int, map[int32]*userPb.User, error) {
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
	userInfoMap := make(map[int32]*userPb.User)
	if s.UserClient != nil && len(userIDs) > 0 {
		resp, err := s.UserClient.GetUsersByIDs(ctx, &userPb.GetUsersByIDsRequest{Ids: userIDs})
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
) (*ent.Employee, *userPb.User, error) {
	emp, err := s.Client.Employee.Query().
		Where(employee.ID(id), employee.OrgID(orgID)).
		WithPosition(func(q *ent.PositionQuery) {
			q.WithDepartments()
		}).
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	var userInfo *userPb.User
	if emp.UserID != "" && s.UserClient != nil {
		if uid, err := strconv.Atoi(emp.UserID); err == nil {
			resp, err := s.UserClient.GetUsersByIDs(
				ctx,
				&userPb.GetUsersByIDsRequest{Ids: []int32{int32(uid)}},
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
			_, err := s.UserClient.DeleteUserByID(ctx, &userPb.DeleteUserRequest{Id: int32(uid)})
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

type EmployeeUpdateInput struct {
	Code      string   `json:"code"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Gender    string   `json:"gender"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Address   string   `json:"address"`
	WardCode  int      `json:"ward_code"`
	AvatarURL string   `json:"avatar_url"`
	JoiningAt string   `json:"joining_at"`
	Status    string   `json:"status"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	RoleIds   []string `json:"role_ids"`
	PermIds   []string `json:"perm_ids"`
}

// UpdateById cập nhật employee (trừ position_id), gọi userPb để cập nhật user nếu có user_id
func (s *EmployeeService) UpdateById(ctx context.Context, id int, orgID int, input EmployeeUpdateInput) (*ent.Employee, *userPb.User, error) {
	emp, err := s.Client.Employee.Query().Where(employee.ID(id), employee.OrgID(orgID)).Only(ctx)
	if err != nil {
		return nil, nil, &ServiceError{Status: http.StatusNotFound, Msg: "#1 UpdateById: Employee not found or not in your organization"}
	}

	upd := s.Client.Employee.UpdateOneID(id)
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

	// Không cập nhật position_id, các trường không có trong ent/employee cũng không cập nhật
	updatedEmp, err := upd.Save(ctx)
	if err != nil {
		return nil, nil, &ServiceError{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	// Gọi userPb để cập nhật user nếu có user_id (nếu có method UpdateUser)
	var userInfo *userPb.User
	fmt.Println(000000, emp.UserID)

	if emp.UserID != "" && s.UserClient != nil {
		fmt.Println(111111, emp.UserID)
		if uid, err := strconv.Atoi(emp.UserID); err == nil {
			fmt.Println(222222, emp.UserID)
			userReq := &userPb.UpdateUserRequest{
				Id:        int32(uid),
				FirstName: input.FirstName,
				LastName:  input.LastName,
				Email:     input.Email,
				Gender:    input.Gender,
				Phone:     input.Phone,
				Address:   input.Address,
				WardCode:  strconv.Itoa(input.WardCode),
				RoleIds:   input.RoleIds,
				PermIds:   input.PermIds,
				Account: &userPb.Account{
					Username: input.Username,
					Password: input.Password,
				},
			}
			userResp, err := s.UserClient.UpdateUserByID(ctx, userReq)
			fmt.Println(333333, err)
			if err == nil {
				fmt.Println(123456, userResp)
				userInfo = userResp.User
			}
		}
	}

	return updatedEmp, userInfo, nil
}
