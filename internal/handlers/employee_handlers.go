package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/huynhthanhthao/hrm_user_service/generated"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
)

type EmployeeHandler struct {
	Client     *ent.Client
	UserClient *pb.UserServiceClient
}

func NewEmployeeHandler(client *ent.Client, userClient *pb.UserServiceClient) *EmployeeHandler {
	return &EmployeeHandler{
		Client:     client,
		UserClient: userClient,
	}
}

// GetEmployees trả về danh sách tất cả nhân viên
func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	employees, err := h.Client.Employee.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch employees"})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// GetEmployeeByID trả về thông tin nhân viên theo ID
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
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

// UpdateEmployee cập nhật thông tin nhân viên
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
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
		// Add other fields as needed
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

// DeleteEmployee xóa nhân viên theo ID
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
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

func (h *EmployeeHandler) RegisterRoutes(r *gin.Engine) {
	employees := r.Group("/employees")
	{
		employees.GET("/", h.GetEmployees)
		employees.GET("/:id", h.GetEmployeeByID)
		employees.PUT("/:id", h.UpdateEmployee)
		employees.DELETE("/:id", h.DeleteEmployee)
	}
}
