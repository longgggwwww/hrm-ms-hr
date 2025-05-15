package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pb "github.com/huynhthanhthao/hrm_user_service/generated"

	"github.com/longgggwwww/hrm-ms-hr/ent"
)

type BranchHandler struct {
	Client     *ent.Client
	UserClient *pb.UserServiceClient
}

func NewBranchHandler(client *ent.Client, userClient *pb.UserServiceClient) *BranchHandler {
	return &BranchHandler{
		Client:     client,
		UserClient: userClient,
	}
}

func (h *BranchHandler) GetBranchs(c *gin.Context) {
	employees, err := h.Client.Branch.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch branches"})
		return
	}
	c.JSON(200, employees)
}

func (h *BranchHandler) GetBranchByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	employee, err := h.Client.Branch.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Branch not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	// Implement the logic to update an existing employee
}

func (h *BranchHandler) DeleteBranch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}
	err = h.Client.Branch.DeleteOneID(id).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *BranchHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/branches")
	{
		group.GET("/", h.GetBranchs)
		group.GET("/:id", h.GetBranchByID)
		group.PUT("/:id", h.UpdateBranch)
		group.DELETE("/:id", h.DeleteBranch)
	}
}
