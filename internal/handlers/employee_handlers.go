package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huynhthanhthao/hrm-ms-shared/middleware"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/internal/constants"
	"github.com/longgggwwww/hrm-ms-hr/internal/dtos"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
	"github.com/longgggwwww/hrm-ms-hr/internal/services"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type EmployeeHandler struct {
	Client     *ent.Client
	UserClient grpc_clients.UserServiceClient
}

func NewEmployeeHandler(client *ent.Client, userClient grpc_clients.UserServiceClient) *EmployeeHandler {
	return &EmployeeHandler{
		Client:     client,
		UserClient: userClient,
	}
}

func (h *EmployeeHandler) RegisterRoutes(r *gin.Engine) {
	employees := r.Group("employees")
	{
		employees.POST("", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.EmployeeCreate},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.Create(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
		employees.GET("", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.EmployeeRead},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.List(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
		employees.GET(":id", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.EmployeeRead},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.GetById(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
		employees.PATCH(":id", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.EmployeeUpdate},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.UpdateById(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
		employees.PATCH(":id/position", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.EmployeeUpdate},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.UpdatePositionById(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
		employees.DELETE(":id", func(c *gin.Context) {
			middleware.AuthMiddleware([]string{constants.EmployeeDelete},
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.Request = r
					h.DeleteById(c)
				})).ServeHTTP(c.Writer, c.Request)
		})
		employees.POST("/root", h.CreateOrgAndRootEmployee)
	}
}

// CreateOrgAndRootEmployee tạo mới tổ chức và employee root đầu tiên
func (h *EmployeeHandler) CreateOrgAndRootEmployee(c *gin.Context) {
	var input dtos.CreateOrgAndRootEmployeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := services.NewEmployeeService(h.Client, h.UserClient)
	org, emp, userInfo, err := svc.CreateOrgAndRootEmployee(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"organization": gin.H{
			"id":   org.ID,
			"name": org.Name,
			"code": org.Code,
		},
		"employee": gin.H{
			"id":          emp.ID,
			"code":        emp.Code,
			"status":      emp.Status,
			"position_id": emp.PositionID,
			"joining_at":  emp.JoiningAt,
			"org_id":      emp.OrgID,
			"created_at":  emp.CreatedAt,
			"updated_at":  emp.UpdatedAt,
			"user_id":     emp.UserID,
		},
		"user_info": normalizeUserInfo(userInfo),
	})
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var input dtos.EmployeeCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil || ids["org_id"] == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing org_id in token"})
		return
	}
	orgID := ids["org_id"]

	svc := services.NewEmployeeService(h.Client, h.UserClient)
	employeeObj, userResp, err := svc.Create(c.Request.Context(), orgID, input)
	if err != nil {
		if serr, ok := err.(*services.ServiceError); ok {
			c.JSON(serr.Status, gin.H{"error": serr.Msg})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userInfo *grpc_clients.User
	if userResp != nil && userResp.User != nil {
		userInfo = userResp.User
	}

	resp := gin.H{
		"id":          employeeObj.ID,
		"code":        employeeObj.Code,
		"status":      employeeObj.Status,
		"position_id": employeeObj.PositionID,
		"joining_at":  employeeObj.JoiningAt,
		"org_id":      employeeObj.OrgID,
		"created_at":  employeeObj.CreatedAt,
		"updated_at":  employeeObj.UpdatedAt,
		"edges":       gin.H{},
		"user_info":   normalizeUserInfo(userInfo),
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *EmployeeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orderBy := c.DefaultQuery("order_by", "id")
	orderDir := c.DefaultQuery("order_dir", "asc")

	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil || ids["org_id"] == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "#1 List: Invalid or missing org_id in token"})
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
		c.JSON(http.StatusBadGateway, gin.H{"error": "#2 List: Failed to fetch employees"})
		return
	}

	var data []gin.H
	// In List handler, ensure user_info is always normalized
	for _, emp := range employees {
		var userInfo *grpc_clients.User = nil
		if emp.UserID != "" {
			if id, err := strconv.Atoi(emp.UserID); err == nil {
				userInfo = userInfoMap[int32(id)]
			}
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
			"edges":       emp.Edges,
			"user_info":   normalizeUserInfo(userInfo),
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

func (h *EmployeeHandler) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "#1 GetById: Invalid employee ID"})
		return
	}
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil || ids["org_id"] == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "#2 GetById: Invalid or missing org_id in token"})
		return
	}
	orgID := ids["org_id"]

	svc := services.NewEmployeeService(h.Client, h.UserClient)
	emp, userInfo, err := svc.GetEmployeeWithUserInfo(c.Request.Context(), id, orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "#3 GetById: Employee not found"})
		return
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
		"edges":       emp.Edges,
		"user_info":   normalizeUserInfo(userInfo),
	}
	c.JSON(http.StatusOK, resp)
}

func (h *EmployeeHandler) UpdateById(c *gin.Context) {
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

	var input dtos.EmployeeUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := services.NewEmployeeService(h.Client, h.UserClient)
	emp, userInfo, err := svc.UpdateById(c.Request.Context(), id, orgID, input)
	if err != nil {
		if serr, ok := err.(*services.ServiceError); ok {
			c.JSON(serr.Status, gin.H{"error": serr.Msg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := gin.H{
		"id":          emp.ID,
		"code":        emp.Code,
		"status":      emp.Status,
		"position_id": emp.PositionID,
		"joining_at":  emp.JoiningAt,
		"org_id":      emp.OrgID,
		"created_at":  emp.CreatedAt,
		"updated_at":  emp.UpdatedAt,
		"user_id":     emp.UserID,
		"edges":       gin.H{},
		"user_info":   normalizeUserInfo(userInfo),
	}

	c.JSON(http.StatusOK, resp)
}

// UpdatePositionById cập nhật position cho employee và lưu lịch sử bổ nhiệm
func (h *EmployeeHandler) UpdatePositionById(c *gin.Context) {
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

	var input dtos.UpdateEmployeePositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := services.NewEmployeeService(h.Client, h.UserClient)
	emp, err := svc.UpdatePositionAndLogHistory(
		c.Request.Context(),
		orgID,
		id,
		input,
	)
	if err != nil {
		if serr, ok := err.(*services.ServiceError); ok {
			c.JSON(serr.Status, gin.H{"error": serr.Msg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          emp.ID,
		"code":        emp.Code,
		"status":      emp.Status,
		"position_id": emp.PositionID,
		"org_id":      emp.OrgID,
		"joining_at":  emp.JoiningAt,
		"created_at":  emp.CreatedAt,
		"updated_at":  emp.UpdatedAt,
		"user_id":     emp.UserID,
	})
}

func (h *EmployeeHandler) DeleteById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "#1 DeleteById: Invalid employee ID"})
		return
	}
	ids, err := utils.ExtractIDsFromToken(c)
	if err != nil || ids["org_id"] == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "#2 DeleteById: Invalid or missing org_id in token"})
		return
	}
	orgID := ids["org_id"]

	svc := services.NewEmployeeService(h.Client, h.UserClient)
	emp, err := svc.DeleteById(c.Request.Context(), id, orgID)
	if err != nil {
		if serr, ok := err.(*services.ServiceError); ok {
			c.JSON(serr.Status, gin.H{"error": serr.Msg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := gin.H{
		"id":          emp.ID,
		"code":        emp.Code,
		"status":      emp.Status,
		"position_id": emp.PositionID,
		"joining_at":  emp.JoiningAt,
		"org_id":      emp.OrgID,
		"created_at":  emp.CreatedAt,
		"updated_at":  emp.UpdatedAt,
		"user_id":     emp.UserID,
		"edges":       gin.H{},
	}

	c.JSON(http.StatusOK, resp)
}

func normalizeUserInfo(user *grpc_clients.User) gin.H {
	if user == nil {
		return gin.H{}
	}
	info := gin.H{
		"id":         user.Id,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"gender":     user.Gender,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	if user.Phone != nil {
		info["phone"] = user.Phone.Value
	}
	if user.Email != nil {
		info["email"] = user.Email.Value
	}
	if user.WardCode != nil {
		info["ward_code"] = user.WardCode.Value
	}
	if user.Address != nil {
		info["address"] = user.Address.Value
	}
	if user.Avatar != nil {
		info["avatar"] = user.Avatar.Value
	}
	return info
}
