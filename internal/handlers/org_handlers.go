package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	pb "github.com/huynhthanhthao/hrm_user_service/generated"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/organization"
)

type OrgHandler struct {
	Client     *ent.Client
	UserClient *pb.UserServiceClient
}

func NewOrgHandler(client *ent.Client, userClient *pb.UserServiceClient) *OrgHandler {
	return &OrgHandler{
		Client:     client,
		UserClient: userClient,
	}
}

// GetOrgs trả về danh sách tất cả các tổ chức
func (h *OrgHandler) GetOrgs(c *gin.Context) {
	orgs, err := h.Client.Organization.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch orgs"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

// GetOrgByID trả về thông tin tổ chức theo ID
func (h *OrgHandler) GetOrgByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid org ID"})
		return
	}

	org, err := h.Client.Organization.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Org not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}

// GetOrgFromToken trả về thông tin tổ chức dựa vào org_id trong JWT
func (h *OrgHandler) GetOrgFromToken(c *gin.Context) {
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

	orgIDStr, ok := claims["org_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org_id not found in token"})
		return
	}

	orgID, err := strconv.Atoi(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid org_id in token"})
		return
	}

	org, err := h.Client.Organization.Get(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Org not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}

// CreateOrg tạo mới một tổ chức
func (h *OrgHandler) CreateOrg(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Code     string `json:"code" binding:"required"`
		Address  string `json:"address"`
		LogoUrl  string `json:"logo_url"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Website  string `json:"website"`
		ParentID *int   `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orgCreate := h.Client.Organization.Create().
		SetName(req.Name).
		SetCode(req.Code).
		SetNillableAddress(&req.Address).
		SetNillableLogoURL(&req.LogoUrl).
		SetNillablePhone(&req.Phone).
		SetNillableEmail(&req.Email).
		SetNillableWebsite(&req.Website).
		SetNillableParentID(req.ParentID)
	org, err := orgCreate.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	// Lấy lại org kèm parent (nếu có)
	orgWithParent, err := h.Client.Organization.Query().
		Where(organization.ID(org.ID)).
		WithParent().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusCreated, org) // fallback nếu lỗi
		return
	}
	c.JSON(http.StatusCreated, orgWithParent)
}

// UpdateOrg cập nhật thông tin tổ chức
func (h *OrgHandler) UpdateOrg(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid org ID"})
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Code     *string `json:"code"`
		Address  *string `json:"address"`
		LogoUrl  *string `json:"logo_url"`
		Phone    *string `json:"phone"`
		Email    *string `json:"email"`
		Website  *string `json:"website"`
		ParentID *int    `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgUpdate := h.Client.Organization.UpdateOneID(id)
	if req.Name != nil {
		orgUpdate.SetName(*req.Name)
	}
	if req.Code != nil {
		orgUpdate.SetCode(*req.Code)
	}
	if req.Address != nil {
		orgUpdate.SetAddress(*req.Address)
	}
	if req.LogoUrl != nil {
		orgUpdate.SetLogoURL(*req.LogoUrl)
	}
	if req.Phone != nil {
		orgUpdate.SetPhone(*req.Phone)
	}
	if req.Email != nil {
		orgUpdate.SetEmail(*req.Email)
	}
	if req.Website != nil {
		orgUpdate.SetWebsite(*req.Website)
	}
	if req.ParentID != nil {
		orgUpdate.SetParentID(*req.ParentID)
	}

	_, err = orgUpdate.Save(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Org not found"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	// Lấy lại org kèm parent (nếu có)
	orgWithParent, err := h.Client.Organization.Query().
		Where(organization.ID(id)).
		WithParent().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"id": id}) // fallback nếu lỗi
		return
	}
	c.JSON(http.StatusOK, orgWithParent)
}

// DeleteOrg xóa tổ chức theo ID
func (h *OrgHandler) DeleteOrg(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid org ID"})
		return
	}
	_, err = h.Client.Organization.Delete().Where(organization.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Org not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// DeleteOrgs xóa nhiều tổ chức theo danh sách ID
func (h *OrgHandler) DeleteOrgs(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}
	_, err := h.Client.Organization.Delete().Where(organization.IDIn(req.IDs...)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *OrgHandler) RegisterRoutes(r *gin.Engine) {
	orgs := r.Group("/orgs")
	{
		orgs.GET("/", h.GetOrgs)
		orgs.GET("/:id", h.GetOrgByID)
		orgs.GET("/from-token", h.GetOrgFromToken)
		orgs.POST("/", h.CreateOrg)
		orgs.PATCH("/:id", h.UpdateOrg)
		orgs.DELETE("/:id", h.DeleteOrg)
		orgs.DELETE("/", h.DeleteOrgs) // Thêm endpoint xóa nhiều
	}
}
