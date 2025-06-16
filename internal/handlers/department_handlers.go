package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/services/department"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type DepartmentHandler struct {
	Service *department.DepartmentService
}

func NewDeptHandler(client *ent.Client) *DepartmentHandler {
	return &DepartmentHandler{
		Service: department.NewDepartmentService(client),
	}
}

func (h *DepartmentHandler) RegisterRoutes(r *gin.Engine) {
	depts := r.Group("/departments")
	{
		depts.POST("/", h.Create)
		depts.POST("/bulk", h.CreateBulk)
		depts.GET("/", h.List)
		depts.GET("/:id", h.Get)
		depts.PATCH("/:id", h.Update)
		depts.DELETE("/:id", h.Delete)
		depts.DELETE("/", h.DeleteBulk)
	}
}

// List lấy danh sách phòng ban với các tùy chọn lọc, sắp xếp và phân trang nâng cao
func (h *DepartmentHandler) List(c *gin.Context) {
	// Extract org_id from JWT token
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract org_id from token: " + err.Error(),
		})
		return
	}

	orgID, ok := tokenData["org_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "org_id not found in token",
		})
		return
	}

	// Parse query parameters
	var query dtos.DepartmentListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set org_id from token
	query.OrgID = orgID

	// Set defaults
	if query.PaginationType == "" {
		query.PaginationType = "page"
	}

	// Call service
	response, err := h.Service.List(c.Request.Context(), query)
	if err != nil {
		if serviceErr, ok := err.(*department.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{
				"error": serviceErr.Msg,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DepartmentHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid department ID",
		})
		return
	}

	departmentObj, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		if serviceErr, ok := err.(*department.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{
				"error": serviceErr.Msg,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, departmentObj)
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	var input dtos.DepartmentCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation using the validator
	validate := validator.New()
	dtos.RegisterDepartmentValidators(validate)
	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract org_id from JWT token
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract org_id from token: " + err.Error(),
		})
		return
	}

	orgID, ok := tokenData["org_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "org_id not found in token",
		})
		return
	}

	// Call service
	response, err := h.Service.Create(c.Request.Context(), orgID, input)
	if err != nil {
		if serviceErr, ok := err.(*department.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{
				"error": serviceErr.Msg,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *DepartmentHandler) CreateBulk(c *gin.Context) {
	var req dtos.DepartmentBulkCreateInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation using the validator
	validate := validator.New()
	dtos.RegisterDepartmentValidators(validate)
	if err := validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract org_id from JWT token
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract org_id from token: " + err.Error(),
		})
		return
	}

	orgID, ok := tokenData["org_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "org_id not found in token",
		})
		return
	}

	// Call service
	departments, err := h.Service.CreateBulk(c.Request.Context(), orgID, req)
	if err != nil {
		if serviceErr, ok := err.(*department.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{
				"error": serviceErr.Msg,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, departments)
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid department ID",
		})
		return
	}

	var input dtos.DepartmentUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation using the validator
	validate := validator.New()
	dtos.RegisterDepartmentValidators(validate)
	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	response, err := h.Service.Update(c.Request.Context(), id, input)
	if err != nil {
		if serviceErr, ok := err.(*department.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{
				"error": serviceErr.Msg,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	// Call service
	err = h.Service.Delete(c.Request.Context(), id)
	if err != nil {
		if serviceErr, ok := err.(*department.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{
				"error": serviceErr.Msg,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *DepartmentHandler) DeleteBulk(c *gin.Context) {
	var req dtos.DepartmentDeleteBulkInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	err := h.Service.DeleteBulk(c.Request.Context(), req)
	if err != nil {
		if serviceErr, ok := err.(*department.ServiceError); ok {
			c.JSON(serviceErr.Status, gin.H{
				"error": serviceErr.Msg,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
