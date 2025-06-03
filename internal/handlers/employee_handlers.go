package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userPb "github.com/huynhthanhthao/hrm_user_service/proto/user"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/internal/services"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type EmployeeHandler struct {
	Client     *ent.Client
	UserClient userPb.UserServiceClient
}

func NewEmployeeHandler(client *ent.Client, userClient userPb.UserServiceClient) *EmployeeHandler {
	return &EmployeeHandler{
		Client:     client,
		UserClient: userClient,
	}
}

func (h *EmployeeHandler) RegisterRoutes(r *gin.Engine) {
	employees := r.Group("employees")
	{
		employees.POST("", h.Create)
		employees.GET("", h.List)
		employees.GET(":id", h.Get)
		employees.PATCH(":id", h.Update)
		employees.DELETE(":id", h.Delete)
	}
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	log.Println("Creating new employee")

	var input services.EmployeeCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	svc := services.NewEmployeeService(h.Client, h.UserClient)
	employeeObj, userResp, err := svc.CreateEmployee(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"employee": employeeObj, "user": userResp})
}

func (h *EmployeeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orderBy := c.DefaultQuery("order_by", "id")
	orderDir := c.DefaultQuery("order_dir", "asc")

	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil || ids["org_id"] == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing org_id in token"})
		return
	}
	orgID := ids["org_id"]

	svc := services.NewEmployeeService(h.Client, h.UserClient)
	query := services.EmployeeListQuery{
		Page:     page,
		Limit:    limit,
		OrderBy:  orderBy,
		OrderDir: orderDir,
		OrgID:    orgID,
	}
	employees, total, userInfoMap, err := svc.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch employees"})
		return
	}

	var data []gin.H
	for _, emp := range employees {
		var userInfo *userPb.User = nil
		if emp.UserID != "" {
			if id, err := strconv.Atoi(emp.UserID); err == nil {
				userInfo = userInfoMap[int32(id)]
			}
		}
		edges := gin.H{}
		if emp.Edges.Position != nil {
			pos := emp.Edges.Position
			posMap := gin.H{
				"id":            pos.ID,
				"name":          pos.Name,
				"code":          pos.Code,
				"department_id": pos.DepartmentID,
				"created_at":    pos.CreatedAt,
				"updated_at":    pos.UpdatedAt,
			}
			if pos.Edges.Departments != nil {
				dept := pos.Edges.Departments
				posMap["department"] = gin.H{
					"id":         dept.ID,
					"name":       dept.Name,
					"code":       dept.Code,
					"org_id":     dept.OrgID,
					"created_at": dept.CreatedAt,
					"updated_at": dept.UpdatedAt,
				}
			}
			edges["position"] = posMap
		}
		item := gin.H{
			"id":          emp.ID,
			"code":        emp.Code,
			"status":      emp.Status,
			"position_id": emp.PositionID,
			"org_id":      emp.OrgID,
			"joining_at":  emp.JoiningAt,
			"created_at":  emp.CreatedAt,
			"updated_at":  emp.UpdatedAt,
			"user_id":     emp.UserID,
			"edges":       edges,
			"user_info":   userInfo,
		}
		data = append(data, item)
	}
	totalPages := (total + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	})
}

func (h *EmployeeHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil || ids["org_id"] == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing org_id in token"})
		return
	}
	orgID := ids["org_id"]

	svc := services.NewEmployeeService(h.Client, h.UserClient)
	emp, userInfo, err := svc.GetEmployeeWithUserInfo(c.Request.Context(), id, orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	edges := gin.H{}
	if emp.Edges.Position != nil {
		pos := emp.Edges.Position
		posMap := gin.H{
			"id":            pos.ID,
			"name":          pos.Name,
			"code":          pos.Code,
			"department_id": pos.DepartmentID,
			"created_at":    pos.CreatedAt,
			"updated_at":    pos.UpdatedAt,
		}
		if pos.Edges.Departments != nil {
			dept := pos.Edges.Departments
			posMap["department"] = gin.H{
				"id":         dept.ID,
				"name":       dept.Name,
				"code":       dept.Code,
				"org_id":     dept.OrgID,
				"created_at": dept.CreatedAt,
				"updated_at": dept.UpdatedAt,
			}
		}
		edges["position"] = posMap
	}
	resp := gin.H{
		"id":          emp.ID,
		"code":        emp.Code,
		"status":      emp.Status,
		"position_id": emp.PositionID,
		"org_id":      emp.OrgID,
		"joining_at":  emp.JoiningAt,
		"created_at":  emp.CreatedAt,
		"updated_at":  emp.UpdatedAt,
		"user_id":     emp.UserID,
		"edges":       edges,
		"user_info":   userInfo,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented. Please use EmployeeService for business logic."})
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid employee ID",
		})
		return
	}
	_, err = h.Client.Employee.Delete().Where(employee.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employee not found",
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
