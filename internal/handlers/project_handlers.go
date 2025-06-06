package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/huynhthanhthao/hrm-ms-shared/middleware"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/internal/constants"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
	"github.com/longgggwwww/hrm-ms-hr/internal/services/project"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type ProjectHandler struct {
	Service *project.ProjectService
}

func NewProjectHandler(client *ent.Client, userClient grpc_clients.UserServiceClient) *ProjectHandler {
	return &ProjectHandler{
		Service: project.NewProjectService(client, userClient),
	}
}

func (h *ProjectHandler) RegisterRoutes(r *gin.Engine) {
	projs := r.Group("/projects")
	{
		projs.POST("/", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.ProjectCreate},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.Create(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
		projs.GET("/", h.List)
		projs.GET("/:id", h.Get)
		projs.PATCH("/:id", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.ProjectUpdate},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.Update(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
		projs.DELETE("/:id", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.ProjectDelete},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.Delete(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
		projs.DELETE("/", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.ProjectDelete},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.DeleteBulk(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
	}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var input dtos.ProjectCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation using the validator
	validate := validator.New()
	dtos.RegisterProjectValidators(validate)
	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract token data
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract data from token: " + err.Error(),
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

	employeeID, ok := tokenData["employee_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "employee_id not found in token",
		})
		return
	}

	// Call service
	response, err := h.Service.Create(c.Request.Context(), orgID, employeeID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ProjectHandler) List(c *gin.Context) {
	// Extract employee_id from JWT token
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract employee_id from token: " + err.Error(),
		})
		return
	}

	employeeID, ok := tokenData["employee_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "employee_id not found in token",
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
	var query dtos.ProjectListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set employee_id and org_id from token
	query.EmployeeID = employeeID
	query.OrgID = orgID

	// Set defaults
	if query.PaginationType == "" {
		query.PaginationType = "page"
	}

	// Call service
	response, err := h.Service.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	project, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var input dtos.ProjectUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation using the validator
	validate := validator.New()
	dtos.RegisterProjectValidators(validate)
	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract employee_id from JWT token
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract employee_id from token: " + err.Error(),
		})
		return
	}

	employeeID, ok := tokenData["employee_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "employee_id not found in token",
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
	response, err := h.Service.Update(c.Request.Context(), id, employeeID, orgID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Call service
	err = h.Service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ProjectHandler) DeleteBulk(c *gin.Context) {
	var req dtos.ProjectDeleteBulkInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation using the validator
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	response, err := h.Service.DeleteBulk(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error: " + err.Error(),
		})
		return
	}

	// Return response based on results
	if response.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, response)
	} else if len(response.FailedIDs) > 0 {
		c.JSON(http.StatusMultiStatus, response) // Using 207 Multi-Status instead of 206 Partial Content
	} else {
		c.JSON(http.StatusOK, response)
	}
}
