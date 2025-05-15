package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pb "github.com/huynhthanhthao/hrm_user_service/generated"

	"github.com/longgggwwww/hrm-ms-hr/ent"
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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.Client.Employee.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// UpdateEmployee cập nhật thông tin nhân viên (chưa implement)
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// DeleteEmployee xóa nhân viên theo ID
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}
	err = h.Client.Employee.DeleteOneID(id).Exec(c.Request.Context())
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
