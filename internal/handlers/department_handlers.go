package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/huynhthanhthao/hrm_user_service/generated"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/department"
)

type DepartmentHandler struct {
	Client     *ent.Client
	UserClient *pb.UserServiceClient
}

func NewDepartmentHandler(client *ent.Client, userClient *pb.UserServiceClient) *DepartmentHandler {
	return &DepartmentHandler{
		Client:     client,
		UserClient: userClient,
	}
}

// GetDepartments trả về danh sách tất cả phòng ban
func (h *DepartmentHandler) GetDepartments(c *gin.Context) {
	departments, err := h.Client.Department.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch departments"})
		return
	}
	c.JSON(http.StatusOK, departments)
}

// GetDepartmentByID trả về thông tin phòng ban theo ID
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}
	departmentObj, err := h.Client.Department.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}
	c.JSON(http.StatusOK, departmentObj)
}

// CreateDepartment tạo mới phòng ban
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	type DepartmentInput struct {
		Name  string `json:"name" binding:"required"`
		Code  string `json:"code" binding:"required"`
		OrgID int    `json:"org_id" binding:"required"`
	}
	var input DepartmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	departmentObj, err := h.Client.Department.Create().
		SetName(input.Name).
		SetCode(input.Code).
		SetOrgID(input.OrgID).
		Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}
	c.JSON(http.StatusCreated, departmentObj)
}

// UpdateDepartment cập nhật thông tin phòng ban
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}
	type DepartmentUpdateInput struct {
		Name  *string `json:"name"`
		Code  *string `json:"code"`
		OrgID *int    `json:"org_id"`
	}
	var input DepartmentUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	update := h.Client.Department.UpdateOneID(id)
	if input.Name != nil {
		update.SetName(*input.Name)
	}
	if input.Code != nil {
		update.SetCode(*input.Code)
	}
	if input.OrgID != nil {
		update.SetOrgID(*input.OrgID)
	}
	departmentObj, err := update.Save(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department"})
		return
	}
	c.JSON(http.StatusOK, departmentObj)
}

// DeleteDepartment xóa phòng ban theo ID
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}
	_, err = h.Client.Department.Delete().Where(department.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *DepartmentHandler) RegisterRoutes(r *gin.Engine) {
	departments := r.Group("/departments")
	{
		departments.GET("", h.GetDepartments)
		departments.GET(":id", h.GetDepartmentByID)
		departments.POST("", h.CreateDepartment)
		departments.PUT(":id", h.UpdateDepartment)
		departments.DELETE(":id", h.DeleteDepartment)
	}
}
