package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	userPb "github.com/huynhthanhthao/hrm_user_service/generated"

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
		employees.PUT(":id", h.Update)
		employees.DELETE(":id", h.Delete)
	}
}

func (h *EmployeeHandler) Create(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	positionObj, err := h.Client.Position.Query().
		Where(position.ID(input.PositionID)).
		WithDepartments().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position_id"})
		return
	}

	joiningAt, err := time.Parse(time.RFC3339, input.JoiningAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid joining_at format, must be RFC3339"})
		return
	}

	tx, err := h.Client.Tx(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	status := employee.Status(input.Status)
	if status != employee.StatusActive && status != employee.StatusInactive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value, must be 'active' or 'inactive'"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserClient is not initialized"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gán userID cho employee vừa tạo
	if respb != nil && respb.User != nil && respb.User.Id > 0 {
		userIDStr := strconv.FormatInt(int64(respb.User.Id), 10)
		_, err := tx.Employee.UpdateOneID(employeeObj.ID).SetUserID(userIDStr).Save(c.Request.Context())
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee with userID"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"employee": employeeObj, "user": respb})
}

func (h *EmployeeHandler) List(c *gin.Context) {
	employees, err := h.Client.Employee.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch employees"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employeeObj, err := h.Client.Employee.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, employeeObj)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}
	type EmployeeUpdateInput struct {
		Name       *string `json:"name"`
		Code       *string `json:"code"`
		Department *int    `json:"department_id"`
		Position   *int    `json:"position_id"`
	}
	var input EmployeeUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update employee"})
	// update := h.Client.Employee.UpdateOneID(id)
	// if input.Name != nil {
	// 	update.SetName(*input.Name)
	// }
	// if input.Code != nil {
	// 	update.SetCode(*input.Code)
	// }
	// if input.Department != nil {
	// 	update.SetDepartmentID(*input.Department)
	// }
	// if input.Position != nil {
	// 	update.SetPositionID(*input.Position)
	// }
	// employeeObj, err := update.Save(c.Request.Context())
	// if err != nil {
	// 	if ent.IsNotFound(err) {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
	// 	return
	// }
	// c.JSON(http.StatusOK, employeeObj)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}
	_, err = h.Client.Employee.Delete().Where(employee.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
