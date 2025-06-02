package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/internal/services"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type TaskReportHandler struct {
	Client *ent.Client
}

func NewTaskReportHandler(client *ent.Client) *TaskReportHandler {
	return &TaskReportHandler{Client: client}
}

func (h *TaskReportHandler) RegisterRoutes(r *gin.Engine) {
	taskReports := r.Group("/task-reports")
	taskReports.POST("", h.Create)
	taskReports.PATCH(":id", h.Update)
}

func (h *TaskReportHandler) Create(c *gin.Context) {
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	reporterID := ids["employee_id"]
	type Input struct {
		TaskID  int    `json:"task_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	report, err := services.CreateTaskReport(c.Request.Context(), h.Client, input.TaskID, reporterID, input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, report)
}

func (h *TaskReportHandler) Update(c *gin.Context) {
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	reporterID := ids["employee_id"]
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task report ID"})
		return
	}
	type Input struct {
		Content string `json:"content" binding:"required"`
	}
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	report, err := services.UpdateTaskReport(c.Request.Context(), h.Client, id, reporterID, input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}
