package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/project"
)

type ProjectHandler struct {
	Client *ent.Client
}

func NewProjectHandler(client *ent.Client) *ProjectHandler {
	return &ProjectHandler{
		Client: client,
	}
}

func (h *ProjectHandler) RegisterRoutes(r *gin.Engine) {
	projs := r.Group("projects")
	{
		projs.POST("", h.Create)
		projs.GET("", h.List)
		projs.GET(":id", h.Get)
		projs.PATCH(":id", h.Update)
		projs.DELETE(":id", h.Delete)
	}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Code        string  `json:"code" binding:"required"`
		Description *string `json:"description"`
		StartAt     string  `json:"start_at" binding:"required"`
		EndAt       *string `json:"end_at"`
		CreatorID   int     `json:"creator_id" binding:"required"`
		UpdaterID   int     `json:"updater_id" binding:"required"`
		OrgID       int     `json:"org_id" binding:"required"`
		Process     *int    `json:"process"`
		Status      *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_at format, must be RFC3339"})
		return
	}

	var endAtPtr *time.Time
	if req.EndAt != nil {
		endAt, err := time.Parse(time.RFC3339, *req.EndAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid end_at format, must be RFC3339",
			})
			return
		}
		endAtPtr = &endAt
	}

	var statusVal project.Status
	if req.Status != nil {
		switch *req.Status {
		case string(project.StatusNotStarted),
			string(project.StatusInProgress),
			string(project.StatusCompleted):
			statusVal = project.Status(*req.Status)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status value",
			})
			return
		}
	} else {
		statusVal = project.StatusNotStarted
	}

	projectCreate := h.Client.Project.Create().
		SetName(req.Name).
		SetCode(req.Code).
		SetNillableDescription(req.Description).
		SetStartAt(startAt).
		SetNillableEndAt(endAtPtr).
		SetCreatorID(req.CreatorID).
		SetUpdaterID(req.UpdaterID).
		SetOrgID(req.OrgID).
		SetNillableProcess(req.Process).
		SetStatus(statusVal)

	row, err := projectCreate.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, row)
}

func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.Client.Project.Query().
		WithTasks().
		All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.Client.Project.Query().Where(project.ID(id)).Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Code        *string `json:"code"`
		Description *string `json:"description"`
		StartAt     *string `json:"start_at"`
		EndAt       *string `json:"end_at"`
		CreatorID   *int    `json:"creator_id"`
		UpdaterID   *int    `json:"updater_id"`
		OrgID       *int    `json:"org_id"`
		Process     *int    `json:"process"`
		Status      *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectUpdate := h.Client.Project.UpdateOneID(id)
	if req.Name != nil {
		projectUpdate.SetName(*req.Name)
	}
	if req.Code != nil {
		projectUpdate.SetCode(*req.Code)
	}
	if req.Description != nil {
		projectUpdate.SetDescription(*req.Description)
	}
	if req.StartAt != nil {
		startAt, err := time.Parse(time.RFC3339, *req.StartAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_at format, must be RFC3339"})
			return
		}
		projectUpdate.SetStartAt(startAt)
	}
	if req.EndAt != nil {
		endAt, err := time.Parse(time.RFC3339, *req.EndAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_at format, must be RFC3339"})
			return
		}
		projectUpdate.SetEndAt(endAt)
	}
	if req.CreatorID != nil {
		projectUpdate.SetCreatorID(*req.CreatorID)
	}
	if req.UpdaterID != nil {
		projectUpdate.SetUpdaterID(*req.UpdaterID)
	}
	if req.OrgID != nil {
		projectUpdate.SetOrgID(*req.OrgID)
	}
	if req.Process != nil {
		projectUpdate.SetProcess(*req.Process)
	}
	if req.Status != nil {
		switch *req.Status {
		case string(project.StatusNotStarted), string(project.StatusInProgress), string(project.StatusCompleted):
			projectUpdate.SetStatus(project.Status(*req.Status))
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}
	}

	_, err = projectUpdate.Save(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	project, err := h.Client.Project.Query().Where(project.ID(id)).Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"id": id})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	_, err = h.Client.Project.Delete().Where(project.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
