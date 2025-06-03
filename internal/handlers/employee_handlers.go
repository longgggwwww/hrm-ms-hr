package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	userPb "github.com/huynhthanhthao/hrm_user_service/proto/user"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/position"
)

type EmployeeHandler struct {
	Client     *ent.Client
	UserClient userPb.UserServiceClient
}

func NewEmployeeHandler(client *ent.Client, userClient userPb.UserServiceClient) *EmployeeHandler {
	return &EmployeeHandler{
		Client:     client,
		UserClient: userClient,
	}
}

func (h *EmployeeHandler) RegisterRoutes(r *gin.Engine) {
	employees := r.Group("employees")
	{
		employees.POST("", h.Create)
		employees.GET("", h.List)
		employees.GET(":id", h.Get)
		employees.PATCH(":id", h.Update)
		employees.DELETE(":id", h.Delete)
	}
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	log.Println("Creating new employee")
	type EmployeeCreateInput struct {
		Code       string   `json:"code" binding:"required"`
		FirstName  string   `json:"first_name" binding:"required"`
		LastName   string   `json:"last_name" binding:"required"`
		Gender     string   `json:"gender" binding:"required,oneof=male female"`
		Phone      string   `json:"phone" binding:"required"`
		Email      string   `json:"email" binding:"omitempty"`
		Address    string   `json:"address" binding:"omitempty"`
		WardCode   int      `json:"ward_code" binding:"omitempty"`
		AvatarURL  string   `json:"avatar_url" binding:"omitempty"`
		PositionID int      `json:"position_id" binding:"required"`
		JoiningAt  string   `json:"joining_at" binding:"required"` // ISO8601 string
		Status     string   `json:"status" binding:"omitempty,oneof=active inactive"`
		Username   string   `json:"username" binding:"required"`
		Password   string   `json:"password" binding:"required"`
		RoleIds    []string `json:"role_ids" binding:"omitempty"`
		PermIds    []string `json:"perm_ids" binding:"omitempty"`
	}
	var input EmployeeCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	positionObj, err := h.Client.Position.Query().
		Where(position.ID(input.PositionID)).
		WithDepartments().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid position_id",
		})
		return
	}

	joiningAt, err := time.Parse(time.RFC3339, input.JoiningAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid joining_at format, must be RFC3339",
		})
		return
	}

	tx, err := h.Client.Tx(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to start transaction",
		})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	status := employee.Status(input.Status)
	if status != employee.StatusActive && status != employee.StatusInactive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status value, must be 'active' or 'inactive'",
		})
		return
	}

	employeeObj, err := tx.Employee.Create().
		SetCode(input.Code).
		SetPositionID(input.PositionID).
		SetOrgID(positionObj.Edges.Departments.OrgID).
		SetJoiningAt(joiningAt).
		SetStatus(status).
		Save(c.Request.Context())
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if h.UserClient == nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "UserClient is not initialized",
		})
		return
	}

	respb, err := h.UserClient.CreateUser(c.Request.Context(), &userPb.CreateUserRequest{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Gán userID cho employee vừa tạo
	if respb != nil && respb.User != nil && respb.User.Id > 0 {
		userIDStr := strconv.FormatInt(int64(respb.User.Id), 10)
		_, err := tx.Employee.UpdateOneID(employeeObj.ID).SetUserID(userIDStr).Save(c.Request.Context())
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update employee with userID",
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"employee": employeeObj, "user": respb})
}

func (h *EmployeeHandler) List(c *gin.Context) {
	employees, err := h.Client.Employee.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to fetch employees",
		})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid employee ID",
		})
		return
	}

	employeeObj, err := h.Client.Employee.Query().
		Where(employee.ID(id)).
		WithPosition(func(q *ent.PositionQuery) {
			q.WithDepartments()
		}).
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employee not found",
		})
		return
	}
	c.JSON(http.StatusOK, employeeObj)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid employee ID",
		})
		return
	}
	type EmployeeUpdateInput struct {
		Code       *string   `json:"code"`
		FirstName  *string   `json:"first_name"`
		LastName   *string   `json:"last_name"`
		Gender     *string   `json:"gender"`
		Phone      *string   `json:"phone"`
		Email      *string   `json:"email"`
		Address    *string   `json:"address"`
		WardCode   *int      `json:"ward_code"`
		AvatarURL  *string   `json:"avatar_url"`
		PositionID *int      `json:"position_id"`
		JoiningAt  *string   `json:"joining_at"`
		Status     *string   `json:"status"`
		Username   *string   `json:"username"`
		Password   *string   `json:"password"`
		RoleIds    *[]string `json:"role_ids"`
		PermIds    *[]string `json:"perm_ids"`
	}
	var input EmployeeUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, err := h.Client.Tx(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to start transaction",
		})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	employeeObj, err := tx.Employee.Get(c.Request.Context(), id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employee not found",
		})
		return
	}

	update := tx.Employee.UpdateOneID(id)
	if input.Code != nil {
		update.SetCode(*input.Code)
	}
	if input.PositionID != nil {
		positionObj, err := h.Client.Position.Query().Where(position.ID(*input.PositionID)).WithDepartments().Only(c.Request.Context())
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid position_id",
			})
			return
		}
		update.SetPositionID(*input.PositionID)
		update.SetOrgID(positionObj.Edges.Departments.OrgID)
	}
	if input.JoiningAt != nil {
		joiningAt, err := time.Parse(time.RFC3339, *input.JoiningAt)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid joining_at format, must be RFC3339",
			})
			return
		}
		update.SetJoiningAt(joiningAt)
	}
	if input.Status != nil {
		status := employee.Status(*input.Status)
		if status != employee.StatusActive && status != employee.StatusInactive {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value, must be 'active' or 'inactive'",
			})
			return
		}
		update.SetStatus(status)
	}
	// update các trường user liên quan qua UserService
	if h.UserClient != nil && (input.FirstName != nil || input.LastName != nil || input.Email != nil || input.Gender != nil || input.Phone != nil || input.Address != nil || input.WardCode != nil || input.Username != nil || input.Password != nil || input.RoleIds != nil || input.PermIds != nil) {
		userID := employeeObj.UserID
		if userID == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Employee does not have a linked user",
			})
			return
		}
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid userID format",
			})
			return
		}
		wardCodeStr := ""
		if input.WardCode != nil {
			wardCodeStr = strconv.Itoa(*input.WardCode)
		}
		_, err = h.UserClient.UpdateUserByID(c.Request.Context(), &userPb.UpdateUserRequest{
			Id:        int32(userIDInt),
			FirstName: derefStr(input.FirstName),
			LastName:  derefStr(input.LastName),
			Email:     derefStr(input.Email),
			Gender:    derefStr(input.Gender),
			Phone:     derefStr(input.Phone),
			Address:   derefStr(input.Address),
			WardCode:  wardCodeStr,
			RoleIds:   derefStrSlice(input.RoleIds),
			PermIds:   derefStrSlice(input.PermIds),
			Account: &userPb.Account{
				Username: derefStr(input.Username),
				Password: derefStr(input.Password),
				Status:   *input.Status,
			},
		})
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	employeeObj, err = update.Save(c.Request.Context())
	if err != nil {
		tx.Rollback()
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Employee not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update employee",
		})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusOK, employeeObj)
}

// Helper functions for pointer deref
func derefStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func derefStrSlice(s *[]string) []string {
	if s != nil {
		return *s
	}
	return []string{} // Trả về mảng rỗng thay vì nil
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid employee ID",
		})
		return
	}
	_, err = h.Client.Employee.Delete().Where(employee.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employee not found",
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
