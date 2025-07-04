package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/organization"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type OrgHandler struct {
	Client      *ent.Client
	UserService *grpc_clients.UserServiceClient
}

func NewOrgHandler(client *ent.Client, userService *grpc_clients.UserServiceClient) *OrgHandler {
	return &OrgHandler{
		Client:      client,
		UserService: userService,
	}
}

func (h *OrgHandler) RegisterRoutes(r *gin.Engine) {
	orgs := r.Group("/orgs")
	{
		orgs.POST("/", h.Create)
		orgs.GET("/", h.List)
		orgs.GET("/:id", h.Get)
		orgs.GET("/from-token", h.GetOrgFromToken)
		orgs.PATCH("/:id", h.Update)
		orgs.DELETE("/:id", h.Delete)
		orgs.DELETE("/", h.DeleteBulk)
	}
}

func (h *OrgHandler) Create(c *gin.Context) {
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
	row, err := orgCreate.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	org, err := h.Client.Organization.Query().
		Where(organization.ID(row.ID)).
		WithParent().
		WithChildren().
		WithDepartments().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusCreated, row)
		return
	}
	c.JSON(http.StatusCreated, org)
}

func (h *OrgHandler) List(c *gin.Context) {
	orgs, err := h.Client.Organization.Query().
		WithParent().
		WithChildren().
		WithDepartments().
		All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch orgs"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func (h *OrgHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid org ID"})
		return
	}

	org, err := h.Client.Organization.Query().
		Where(organization.ID(id)).
		WithParent().
		WithChildren().
		WithDepartments().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Org not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}

func (h *OrgHandler) Update(c *gin.Context) {
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

func (h *OrgHandler) Delete(c *gin.Context) {
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

func (h *OrgHandler) DeleteBulk(c *gin.Context) {
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

	_, err := h.Client.Organization.Delete().
		Where(organization.IDIn(req.IDs...)).
		Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *OrgHandler) GetOrgFromToken(c *gin.Context) {
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orgID, ok := ids["org_id"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org_id not found in token"})
		return
	}

	org, err := h.Client.Organization.Query().
		Where(organization.ID(orgID)).
		WithParent().
		WithChildren().
		WithDepartments().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Org not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}
