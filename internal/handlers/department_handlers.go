package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userpb "github.com/huynhthanhthao/hrm_user_service/proto/user"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/department"
)

type DepartmentHandler struct {
	Client *ent.Client
	UserPb *userpb.UserServiceClient
}

func NewDeptHandler(client *ent.Client, userpb *userpb.UserServiceClient) *DepartmentHandler {
	return &DepartmentHandler{
		Client: client,
		UserPb: userpb,
	}
}

func (h *DepartmentHandler) RegisterRoutes(r *gin.Engine) {
	depts := r.Group("departments")
	{
		depts.POST("", h.Create)
		depts.GET("", h.List)
		depts.GET(":id", h.Get)
		depts.PUT(":id", h.Update)
		depts.DELETE(":id", h.Delete)
	}
}

func (h *DepartmentHandler) List(c *gin.Context) {
	depts, err := h.Client.Department.Query().
		WithOrganization().
		WithPositions().
		All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch departments"})
		return
	}
	c.JSON(http.StatusOK, depts)
}

func (h *DepartmentHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dept, err := h.Client.Department.Query().
		Where(department.ID(id)).
		WithOrganization().
		WithPositions().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dept)
}

func (h *DepartmentHandler) Create(c *gin.Context) {
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

func (h *DepartmentHandler) Update(c *gin.Context) {
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

func (h *DepartmentHandler) Delete(c *gin.Context) {
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
