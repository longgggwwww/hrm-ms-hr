package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"github.com/longgggwwww/hrm-ms-hr/ent"
)

type CompanyHandler struct {
	Client     *ent.Client
	UserClient *pb.UserServiceClient
}

func NewCompanyHandler(client *ent.Client, userClient *pb.UserServiceClient) *CompanyHandler {
	return &CompanyHandler{
		Client:     client,
		UserClient: userClient,
	}
}

// GetCompanyFromToken trả về thông tin công ty dựa vào company_id trong JWT
func (h *CompanyHandler) GetCompanyFromToken(c *gin.Context) {
	fmt.Println("GetCompanyFromToken")
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

	companyIDStr, ok := claims["company_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id not found in token"})
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company_id in token"})
		return
	}

	company, err := h.Client.Company.Get(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}
	c.JSON(http.StatusOK, company)
}

// GetCompanies trả về danh sách tất cả các công ty
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	companies, err := h.Client.Company.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch companies"})
		return
	}
	c.JSON(http.StatusOK, companies)
}

// GetCompanyByID trả về thông tin công ty theo ID
func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
	fmt.Println("GetCompanyByID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	company, err := h.Client.Company.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}
	c.JSON(http.StatusOK, company)
}

// UpdateCompany cập nhật thông tin công ty (chưa implement)
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// DeleteCompany xóa công ty theo ID
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}
	err = h.Client.Company.DeleteOneID(id).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *CompanyHandler) RegisterRoutes(r *gin.Engine) {
	companies := r.Group("/companies")
	{
		companies.GET("/", h.GetCompanies)
		companies.GET("/:id", h.GetCompanyByID)
		companies.GET("/from-token", h.GetCompanyFromToken)
		companies.PUT("/:id", h.UpdateCompany)
		companies.DELETE("/:id", h.DeleteCompany)
	}
}
