package services

import (
	"context"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	userPb "github.com/huynhthanhthao/hrm_user_service/proto/user"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/position"
)

type EmployeeService struct {
	Client     *ent.Client
	UserClient userPb.UserServiceClient
}

type EmployeeCreateInput struct {
	Code       string
	FirstName  string
	LastName   string
	Gender     string
	Phone      string
	Email      string
	Address    string
	WardCode   int
	AvatarURL  string
	PositionID int
	JoiningAt  string
	Status     string
	Username   string
	Password   string
	RoleIds    []string
	PermIds    []string
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

func (s *EmployeeService) CreateEmployee(ctx context.Context, input EmployeeCreateInput) (*ent.Employee, *userPb.CreateUserResponse, error) {
	positionObj, err := s.Client.Position.Query().
		Where(position.ID(input.PositionID)).
		WithDepartments().
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	joiningAt, err := time.Parse(time.RFC3339, input.JoiningAt)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	employeeObj, err := tx.Employee.Create().
		SetCode(input.Code).
		SetPositionID(input.PositionID).
		SetOrgID(positionObj.Edges.Departments.OrgID).
		SetJoiningAt(joiningAt).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if s.UserClient == nil {
		tx.Rollback()
		return nil, nil, err
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
func (s *EmployeeService) GetEmployeeWithUserInfo(ctx context.Context, id int, orgID int) (*ent.Employee, *userPb.User, error) {
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
			resp, err := s.UserClient.GetUsersByIDs(ctx, &userPb.GetUsersByIDsRequest{Ids: []int32{int32(uid)}})
			if err == nil && resp != nil && len(resp.Users) > 0 {
				userInfo = resp.Users[0]
			}
		}
	}
	return emp, userInfo, nil
}
