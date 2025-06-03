package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/label"
	"github.com/longgggwwww/hrm-ms-hr/ent/task"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type LabelHandler struct {
	Client *ent.Client
}

func NewLabelHandler(client *ent.Client) *LabelHandler {
	return &LabelHandler{
		Client: client,
	}
}

func (h *LabelHandler) RegisterRoutes(r *gin.Engine) {
	labels := r.Group("/labels")
	{
		labels.POST("/", h.Create)
		labels.POST("/bulk", h.CreateBulk)
		labels.GET("/", h.List)
		labels.GET("/:id", h.Get)
		labels.PATCH("/:id", h.Update)
		labels.DELETE("/:id", h.Delete)
		labels.DELETE("/", h.DeleteBulk)
	}
}

// List lấy danh sách nhãn với các tùy chọn lọc, sắp xếp và phân trang nâng cao
// Tham số truy vấn:
// - name: Lọc theo tên nhãn (tìm kiếm chứa)
// - description: Lọc theo mô tả (tìm kiếm chứa)
// - color: Lọc theo màu chính xác
// - order_by: Trường sắp xếp (id, name, description, color, org_id, created_at, updated_at, task_count) - mặc định: created_at
// - order_dir: Hướng sắp xếp (asc, desc) - mặc định: desc
// - page: Số trang cho phân trang offset (mặc định: 1)
// - limit: Số mục trên mỗi trang cho phân trang offset (mặc định: 10, tối đa: 100)
// - cursor: Con trỏ cho phân trang dựa trên cursor (ID được mã hóa base64)
// - cursor_limit: Số mục trên mỗi trang cursor (mặc định: 10, tối đa: 100)
// - pagination_type: Loại phân trang (page, cursor) - mặc định: page
//
// Ví dụ:
// - Phân trang offset: GET /labels?name=urgent&order_by=name&order_dir=asc&page=1&limit=20
// - Sắp xếp theo số lượng task: GET /labels?order_by=task_count&order_dir=desc&page=1&limit=10
// - Phân trang cursor: GET /labels?pagination_type=cursor&cursor_limit=10&cursor=eyJpZCI6MTB9
// - Lọc theo tổ chức: Nhãn được tự động lọc theo org_id từ JWT token
func (h *LabelHandler) List(c *gin.Context) {
	// Trích xuất org_id từ JWT token
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract org_id from token: " + err.Error(),
		})
		return
	}

	orgID, ok := tokenData["org_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "org_id not found in token",
		})
		return
	}

	query := h.Client.Label.Query()

	// Lọc theo org_id từ token
	query = query.Where(label.OrgIDEQ(orgID))

	// Lọc dữ liệu
	if name := c.Query("name"); name != "" {
		query = query.Where(label.NameContainsFold(name))
	}

	if description := c.Query("description"); description != "" {
		query = query.Where(label.DescriptionContainsFold(description))
	}

	// Xác định loại phân trang
	paginationType := c.DefaultQuery("pagination_type", "page")

	if paginationType == "cursor" {
		h.listWithCursorPagination(c, query)
		return
	}

	// Mặc định: phân trang dựa trên offset
	h.listWithOffsetPagination(c, query)
}

func (h *LabelHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid label ID",
		})
		return
	}

	labelObj, err := h.Client.Label.Query().
		Where(label.ID(id)).
		WithTasks().
		Only(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Label not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch label",
			})
			log.Println("Error fetching label:", err)
		}
		return
	}

	c.JSON(http.StatusOK, labelObj)
}

func (h *LabelHandler) Create(c *gin.Context) {
	type LabelInput struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color" binding:"required"`
	}
	var input LabelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trích xuất org_id từ JWT token
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract org_id from token: " + err.Error(),
		})
		return
	}

	orgID, ok := tokenData["org_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "org_id not found in token",
		})
		return
	}

	create := h.Client.Label.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetColor(input.Color).
		SetOrgID(orgID)

	labelObj, err := create.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create label"})
		return
	}

	// Lấy nhãn đã tạo với các mối quan hệ edge
	createdLabel, err := h.Client.Label.Query().
		Where(label.ID(labelObj.ID)).
		WithTasks().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch created label",
		})
		return
	}

	// Thêm task_count vào response
	labelsWithTaskCount, err := h.addTaskCountsToLabels(c, []*ent.Label{createdLabel})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task count"})
		return
	}

	c.JSON(http.StatusCreated, labelsWithTaskCount[0])
}

func (h *LabelHandler) CreateBulk(c *gin.Context) {
	type LabelInput struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color" binding:"required"`
	}

	var req struct {
		Labels []LabelInput `json:"labels" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Labels) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No labels provided"})
		return
	}

	// Trích xuất org_id từ JWT token
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract org_id from token: " + err.Error(),
		})
		return
	}

	orgID, ok := tokenData["org_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "org_id not found in token",
		})
		return
	}

	bulk := make([]*ent.LabelCreate, len(req.Labels))
	for i, labelInput := range req.Labels {
		create := h.Client.Label.Create().
			SetName(labelInput.Name).
			SetDescription(labelInput.Description).
			SetColor(labelInput.Color).
			SetOrgID(orgID)

		bulk[i] = create
	}

	labels, err := h.Client.Label.CreateBulk(bulk...).Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create labels"})
		return
	}

	// Lấy các nhãn đã tạo với các mối quan hệ edge
	labelIDs := make([]int, len(labels))
	for i, lbl := range labels {
		labelIDs[i] = lbl.ID
	}

	createdLabels, err := h.Client.Label.Query().
		Where(label.IDIn(labelIDs...)).
		WithTasks().
		All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created labels"})
		return
	}

	c.JSON(http.StatusCreated, createdLabels)
}

func (h *LabelHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid label ID",
		})
		return
	}
	type LabelUpdateInput struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Color       *string `json:"color"`
		OrgID       *int    `json:"org_id"`
	}
	var input LabelUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	update := h.Client.Label.UpdateOneID(id)
	update.SetNillableName(input.Name).
		SetNillableDescription(input.Description).
		SetNillableColor(input.Color)

	if input.OrgID != nil {
		update = update.SetNillableOrgID(input.OrgID)
	}

	labelObj, err := update.Save(c.Request.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update label"})
		return
	}

	// Lấy nhãn đã cập nhật với các mối quan hệ edge
	updatedLabel, err := h.Client.Label.Query().
		Where(label.ID(labelObj.ID)).
		WithTasks().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated label"})
		return
	}

	// Thêm task_count vào response
	labelsWithTaskCount, err := h.addTaskCountsToLabels(c, []*ent.Label{updatedLabel})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task count"})
		return
	}

	c.JSON(http.StatusOK, labelsWithTaskCount[0])
}

func (h *LabelHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}
	_, err = h.Client.Label.Delete().Where(label.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *LabelHandler) DeleteBulk(c *gin.Context) {
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

	_, err := h.Client.Label.Delete().
		Where(label.IDIn(req.IDs...)).
		Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *LabelHandler) listWithOffsetPagination(c *gin.Context, query *ent.LabelQuery) {
	// Trích xuất org_id từ JWT token cho truy vấn đếm
	tokenData, err := utils.ExtractIDsFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract org_id from token: " + err.Error(),
		})
		return
	}

	orgID, ok := tokenData["org_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "org_id not found in token",
		})
		return
	}

	// Áp dụng sắp xếp
	orderOption := h.getOrderOption(c)
	if orderOption == nil {
		return // Error response already sent
	}
	query = query.Order(orderOption)

	// Phân trang
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

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	labels, err := query.All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Lấy tổng số lượng cho thông tin phân trang với cùng các bộ lọc
	countQuery := h.Client.Label.Query()

	// Áp dụng bộ lọc org_id từ token
	countQuery = countQuery.Where(label.OrgIDEQ(orgID))

	// Áp dụng các bộ lọc giống như truy vấn chính
	if name := c.Query("name"); name != "" {
		countQuery = countQuery.Where(label.NameContainsFold(name))
	}
	if description := c.Query("description"); description != "" {
		countQuery = countQuery.Where(label.DescriptionContainsFold(description))
	}

	total, err := countQuery.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	// Thêm số lượng task vào nhãn
	labelsWithTaskCount, err := h.addTaskCountsToLabels(c, labels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add task counts",
		})
		return
	}

	response := gin.H{
		"data": labelsWithTaskCount,
		"pagination": gin.H{
			"type":         "offset",
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

// listWithCursorPagination xử lý phân trang dựa trên cursor
func (h *LabelHandler) listWithCursorPagination(c *gin.Context, query *ent.LabelQuery) {
	// Áp dụng sắp xếp
	orderOption := h.getOrderOption(c)
	if orderOption == nil {
		return // Error response already sent
	}
	query = query.Order(orderOption)

	// Thiết lập phân trang cursor
	limit := 10
	if limitStr := c.Query("cursor_limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Phân tích cursor nếu được cung cấp
	if cursor := c.Query("cursor"); cursor != "" {
		cursorData, err := h.decodeCursor(cursor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid cursor format",
			})
			return
		}

		// Áp dụng bộ lọc cursor (giả định cursor dựa trên ID)
		if id, ok := cursorData["id"].(float64); ok {
			query = query.Where(label.IDGT(int(id)))
		}
	}

	// Lấy thêm một mục để xác định có trang tiếp theo không
	labels, err := query.Limit(limit + 1).All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	hasNext := len(labels) > limit
	if hasNext {
		labels = labels[:limit] // Loại bỏ mục thêm
	}

	var nextCursor *string
	if hasNext && len(labels) > 0 {
		lastItem := labels[len(labels)-1]
		cursorStr := h.encodeCursor(map[string]interface{}{
			"id": lastItem.ID,
		})
		nextCursor = &cursorStr
	}

	// Thêm số lượng task vào nhãn
	labelsWithTaskCount, err := h.addTaskCountsToLabels(c, labels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add task counts",
		})
		return
	}

	response := gin.H{
		"data": labelsWithTaskCount,
		"pagination": gin.H{
			"type":        "cursor",
			"per_page":    limit,
			"has_next":    hasNext,
			"next_cursor": nextCursor,
		},
	}

	c.JSON(http.StatusOK, response)
}

// getOrderOption trả về tùy chọn sắp xếp phù hợp dựa trên tham số truy vấn
func (h *LabelHandler) getOrderOption(c *gin.Context) label.OrderOption {
	orderBy := c.DefaultQuery("order_by", "created_at")
	orderDir := c.DefaultQuery("order_dir", "desc")

	var orderOption label.OrderOption
	switch orderBy {
	case "id":
		if orderDir == "asc" {
			orderOption = label.ByID()
		} else {
			orderOption = label.ByID(sql.OrderDesc())
		}
	case "name":
		if orderDir == "asc" {
			orderOption = label.ByName()
		} else {
			orderOption = label.ByName(sql.OrderDesc())
		}
	case "description":
		if orderDir == "asc" {
			orderOption = label.ByDescription()
		} else {
			orderOption = label.ByDescription(sql.OrderDesc())
		}
	case "color":
		if orderDir == "asc" {
			orderOption = label.ByColor()
		} else {
			orderOption = label.ByColor(sql.OrderDesc())
		}
	case "org_id":
		if orderDir == "asc" {
			orderOption = label.ByOrgID()
		} else {
			orderOption = label.ByOrgID(sql.OrderDesc())
		}
	case "created_at":
		if orderDir == "asc" {
			orderOption = label.ByCreatedAt()
		} else {
			orderOption = label.ByCreatedAt(sql.OrderDesc())
		}
	case "updated_at":
		if orderDir == "asc" {
			orderOption = label.ByUpdatedAt()
		} else {
			orderOption = label.ByUpdatedAt(sql.OrderDesc())
		}
	case "task_count":
		// Sắp xếp theo số lượng task sử dụng ByTasksCount của Ent
		if orderDir == "asc" {
			orderOption = label.ByTasksCount()
		} else {
			orderOption = label.ByTasksCount(sql.OrderDesc())
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order_by field. Valid fields: id, name, description, color, org_id, created_at, updated_at, task_count",
		})
		return nil
	}

	return orderOption
}

// encodeCursor mã hóa dữ liệu cursor thành base64
func (h *LabelHandler) encodeCursor(data map[string]any) string {
	jsonData, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(jsonData)
}

// decodeCursor giải mã cursor base64 thành map dữ liệu
func (h *LabelHandler) decodeCursor(cursor string) (map[string]any, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// addTaskCountsToLabels thêm trường task_count vào mỗi nhãn
func (h *LabelHandler) addTaskCountsToLabels(c *gin.Context, labels []*ent.Label) ([]map[string]any, error) {
	result := make([]map[string]any, len(labels))

	for i, labelEntity := range labels {
		// Đếm task cho nhãn này
		taskCount, err := h.Client.Task.Query().
			Where(task.HasLabelsWith(label.IDEQ(labelEntity.ID))).
			Count(c.Request.Context())
		if err != nil {
			return nil, err
		}

		// Chuyển đổi nhãn thành map và thêm task_count
		labelMap := map[string]any{
			"id":          labelEntity.ID,
			"name":        labelEntity.Name,
			"description": labelEntity.Description,
			"color":       labelEntity.Color,
			"org_id":      labelEntity.OrgID,
			"created_at":  labelEntity.CreatedAt,
			"updated_at":  labelEntity.UpdatedAt,
			"task_count":  taskCount,
		}

		result[i] = labelMap
	}

	return result, nil
}
