package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// GetBranches trả về danh sách tất cả các chi nhánh
func (h *BranchHandler) GetBranches(c *gin.Context) {
	branches, err := h.Client.Branch.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch branches"})
		return
	}
	c.JSON(http.StatusOK, branches)
}

// GetBranchByID trả về thông tin chi nhánh theo ID
func (h *BranchHandler) GetBranchByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	branch, err := h.Client.Branch.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}
	c.JSON(http.StatusOK, branch)
}

// GetBranchFromToken trả về thông tin chi nhánh dựa vào branch_id trong JWT
func (h *BranchHandler) GetBranchFromToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	tokenString := parts[1]

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	branchIDStr, ok := claims["branch_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "branch_id not found in token"})
		return
	}

	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch_id in token"})
		return
	}

	branch, err := h.Client.Branch.Get(c.Request.Context(), branchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}
	c.JSON(http.StatusOK, branch)
}

// CreateBranch tạo mới một chi nhánh
func (h *BranchHandler) CreateBranch(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		CompanyID   string `json:"company_id" binding:"required"`
		Address     string `json:"address"`
		ContactInfo string `json:"contact_info"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company_id"})
		return
	}

	branch, err := h.Client.Branch.Create().
		SetName(req.Name).
		SetCode(req.Code).
		SetCompanyID(companyID).
		SetNillableAddress(&req.Address).
		SetNillableContactInfo(&req.ContactInfo).
		Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, branch)
}

// UpdateBranch cập nhật thông tin chi nhánh (chưa implement)
func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// DeleteBranch xóa chi nhánh theo ID
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
	branches := r.Group("/branches")
	{
		branches.GET("/", h.GetBranches)
		branches.GET("/:id", h.GetBranchByID)
		branches.GET("/from-token", h.GetBranchFromToken)
		branches.POST("/", h.CreateBranch)
		branches.PUT("/:id", h.UpdateBranch)
		branches.DELETE("/:id", h.DeleteBranch)
	}
}
