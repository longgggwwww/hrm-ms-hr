package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/services/label"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type LabelHandler struct {
	Service *label.LabelService
}

func NewLabelHandler(client *ent.Client) *LabelHandler {
	return &LabelHandler{
		Service: label.NewLabelService(client),
	}
}

func (h *LabelHandler) RegisterRoutes(r *gin.Engine) {
	labels := r.Group("/labels")
	{
		labels.POST("/", h.Create)
		labels.POST("/bulk", h.CreateBulk)
		labels.GET("/", h.List)
		labels.GET("/:id", h.Get)
		labels.PATCH("/:id", h.Update)
		labels.DELETE("/:id", h.Delete)
		labels.DELETE("/", h.DeleteBulk)
	}
}

// List lấy danh sách nhãn với các tùy chọn lọc, sắp xếp và phân trang nâng cao
func (h *LabelHandler) List(c *gin.Context) {
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
	var query dtos.LabelListQuery
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
		if serviceErr, ok := err.(*label.ServiceError); ok {
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

func (h *LabelHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid label ID",
		})
		return
	}

	labelObj, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		if serviceErr, ok := err.(*label.ServiceError); ok {
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

	c.JSON(http.StatusOK, labelObj)
}

func (h *LabelHandler) Create(c *gin.Context) {
	var input dtos.LabelCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation using the validator
	validate := validator.New()
	dtos.RegisterCustomValidators(validate)
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
		if serviceErr, ok := err.(*label.ServiceError); ok {
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

func (h *LabelHandler) CreateBulk(c *gin.Context) {
	var req dtos.LabelBulkCreateInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation using the validator
	validate := validator.New()
	dtos.RegisterCustomValidators(validate)
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
	labels, err := h.Service.CreateBulk(c.Request.Context(), orgID, req)
	if err != nil {
		if serviceErr, ok := err.(*label.ServiceError); ok {
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

	c.JSON(http.StatusCreated, labels)
}

func (h *LabelHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid label ID",
		})
		return
	}

	var input dtos.LabelUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation using the validator
	validate := validator.New()
	dtos.RegisterCustomValidators(validate)
	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	response, err := h.Service.Update(c.Request.Context(), id, input)
	if err != nil {
		if serviceErr, ok := err.(*label.ServiceError); ok {
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

func (h *LabelHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	// Call service
	err = h.Service.Delete(c.Request.Context(), id)
	if err != nil {
		if serviceErr, ok := err.(*label.ServiceError); ok {
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

func (h *LabelHandler) DeleteBulk(c *gin.Context) {
	var req dtos.LabelDeleteBulkInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	err := h.Service.DeleteBulk(c.Request.Context(), req)
	if err != nil {
		if serviceErr, ok := err.(*label.ServiceError); ok {
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
