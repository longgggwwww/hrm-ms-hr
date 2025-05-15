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

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	employees, err := h.Client.Employee.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch employees"})
		return
	}
	c.JSON(200, employees)
}

func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.Client.Employee.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	// Implement the logic to update an existing employee
}

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
	group := r.Group("/employees")
	{
		group.GET("/", h.GetEmployees)
		group.GET("/:id", h.GetEmployeeByID)
		group.PUT("/:id", h.UpdateEmployee)
		group.DELETE("/:id", h.DeleteEmployee)
	}
}
