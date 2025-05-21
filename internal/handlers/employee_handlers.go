package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userPb "github.com/huynhthanhthao/hrm_user_service/generated"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
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

func (h *EmployeeHandler) Create(c *gin.Context) {
	// type EmployeeCreateInput struct {
	// 	UserID     string `json:"user_id" binding:"required"`
	// 	Code       string `json:"code" binding:"required"`
	// 	PositionID int    `json:"position_id" binding:"required"`
	// 	OrgID      int    `json:"org_id" binding:"required"`
	// 	JoiningAt  string `json:"joining_at" binding:"required"` // ISO8601 string
	// }
	// var input EmployeeCreateInput
	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// joiningAt, err := time.Parse(time.RFC3339, input.JoiningAt)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid joining_at format, must be RFC3339"})
	// 	return
	// }

	if h.UserClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserClient is not initialized"})
		return
	}

	res, err := h.UserClient.CreateUser(c.Request.Context(), &userPb.CreateUserRequest{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "jdakdja@gmail.com",
		Gender:    "male",
		Phone:     "123456789",
		Address:   "123 Main St",
		WardCode:  "3123213",
		RoleIds:   []string{"6c618337-b35b-48e1-8f98-dff55eb8eaf7"},
		Account: &userPb.Account{
			Username: "user01",
			Password: "user01",
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// employeeObj, err := h.Client.Employee.Create().
	// 	SetUserID(input.UserID).
	// 	SetCode(input.Code).
	// 	SetPositionID(input.PositionID).
	// 	SetOrgID(input.OrgID).
	// 	SetJoiningAt(joiningAt).
	// 	Save(c.Request.Context())
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
	// 	return
	// }
	c.JSON(http.StatusCreated, res)
}
