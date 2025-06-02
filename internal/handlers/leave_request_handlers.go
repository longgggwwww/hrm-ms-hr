package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/internal/services"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type LeaveRequestHandler struct {
	Client *ent.Client
}

func NewLeaveRequestHandler(client *ent.Client) *LeaveRequestHandler {
	return &LeaveRequestHandler{
		Client: client,
	}
}

func (h *LeaveRequestHandler) RegisterRoutes(r *gin.Engine) {
	log.Println("vlxx")
	leaveRequests := r.Group("leave-requests")
	{
		leaveRequests.GET("/admin", h.ListAdmin)
		leaveRequests.GET("/employee", h.ListEmployee)
		leaveRequests.GET(":id/admin", h.GetAdmin)
		leaveRequests.GET(":id/employee", h.GetEmployee)
		leaveRequests.POST("", h.Create)
		leaveRequests.PATCH(":id/approve", h.Approve)
		leaveRequests.PATCH(":id/reject", h.Reject)
	}
}

func (h *LeaveRequestHandler) GetAdmin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "#1 GetAdmin: invalid leave request ID"})
		return
	}
	leaveRequest, err := services.GetLeaveRequest(c.Request.Context(), h.Client, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "#2 GetAdmin: Leave request not found"})
		return
	}
	c.JSON(http.StatusOK, leaveRequest)
}

func (h *LeaveRequestHandler) GetEmployee(c *gin.Context) {
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	employeeID := ids["employee_id"]
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "#1 GetEmployee: invalid leave request ID"})
		return
	}
	leaveRequest, err := services.GetLeaveRequest(c.Request.Context(), h.Client, id)
	if err != nil || leaveRequest.Edges.Applicant == nil || leaveRequest.Edges.Applicant.ID != employeeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "#2 GetEmployee: Leave request not found or not owned by employee"})
		return
	}
	c.JSON(http.StatusOK, leaveRequest)
}

func (h *LeaveRequestHandler) Create(c *gin.Context) {
	type LeaveRequestInput struct {
		TotalDays float64 `json:"total_days" binding:"required"`
		StartAt   string  `json:"start_at" binding:"required"`
		EndAt     string  `json:"end_at" binding:"required"`
		Reason    string  `json:"reason"`
		Type      string  `json:"type" binding:"required"`
	}
	var input LeaveRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	employeeID := ids["employee_id"]
	orgID := ids["org_id"]

	startAt, err := time.Parse(time.RFC3339, input.StartAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "#1 Create: invalid start_at format, must be RFC3339"})
		return
	}
	endAt, err := time.Parse(time.RFC3339, input.EndAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "#2 Create: invalid end_at format, must be RFC3339"})
		return
	}
	dto := services.LeaveRequestCreateDTO{
		TotalDays:  input.TotalDays,
		StartAt:    startAt,
		EndAt:      endAt,
		Reason:     input.Reason,
		Type:       input.Type,
		EmployeeID: employeeID,
		OrgID:      orgID,
	}

	leaveRequest, err := services.Create(c.Request.Context(), h.Client, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, leaveRequest)
}

func (h *LeaveRequestHandler) Approve(c *gin.Context) {
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, err)
		return
	}
	reviewerID := ids["employee_id"]
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("#1 Approve: invalid leave request ID"))
		return
	}
	leaveRequest, err := services.ApproveLeaveRequest(c.Request.Context(), h.Client, id, reviewerID)
	if err != nil {
		if svcErr, ok := err.(*services.ServiceError); ok {
			utils.RespondWithError(c, svcErr.Status, errors.New(svcErr.Msg))
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, errors.New("#2 Approve: failed to approve leave request"))
		return
	}
	c.JSON(http.StatusOK, leaveRequest)
}

func (h *LeaveRequestHandler) ListAdmin(c *gin.Context) {
	// Lấy filter, phân trang, order từ query
	filter := map[string]interface{}{}
	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	orderBy := c.DefaultQuery("order_by", "created_at")
	orderDir := c.DefaultQuery("order_dir", "desc")
	list, total, err := services.ListLeaveRequests(c.Request.Context(), h.Client, filter, page, limit, orderBy, orderDir)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	totalPages := (total + limit - 1) / limit
	response := gin.H{
		"data": list,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	}
	c.JSON(http.StatusOK, response)
}

func (h *LeaveRequestHandler) ListEmployee(c *gin.Context) {
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	employeeID := ids["employee_id"]
	filter := map[string]interface{}{"employee_id": employeeID}
	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	orderBy := c.DefaultQuery("order_by", "created_at")
	orderDir := c.DefaultQuery("order_dir", "desc")
	list, total, err := services.ListLeaveRequests(c.Request.Context(), h.Client, filter, page, limit, orderBy, orderDir)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	totalPages := (total + limit - 1) / limit
	response := gin.H{
		"data": list,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	}
	c.JSON(http.StatusOK, response)
}

func (h *LeaveRequestHandler) Reject(c *gin.Context) {
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, err)
		return
	}
	reviewerID := ids["employee_id"]
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("#1 Reject: invalid leave request ID"))
		return
	}
	leaveRequest, err := services.Reject(c.Request.Context(), h.Client, id, reviewerID)
	if err != nil {
		if svcErr, ok := err.(*services.ServiceError); ok {
			utils.RespondWithError(c, svcErr.Status, errors.New(svcErr.Msg))
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, errors.New(err.Error()))
		return
	}
	c.JSON(http.StatusOK, leaveRequest)
}
