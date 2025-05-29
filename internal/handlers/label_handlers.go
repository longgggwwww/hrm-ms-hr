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
	labels := r.Group("labels")
	{
		labels.POST("", h.Create)
		labels.POST("bulk", h.CreateBulk)
		labels.GET("", h.List)
		labels.GET(":id", h.Get)
		labels.PATCH(":id", h.Update)
		labels.DELETE(":id", h.Delete)
		labels.DELETE("", h.DeleteBulk)
	}
}

// List retrieves labels with advanced filtering, sorting and pagination options
// Query Parameters:
// - name: Filter by label name (contains search)
// - description: Filter by description (contains search)
// - color: Filter by exact color match
// - org_id: Filter by organization ID
// - order_by: Sort field (id, name, description, color, org_id, created_at, updated_at) - default: created_at
// - order_dir: Sort direction (asc, desc) - default: desc
// - page: Page number for offset pagination (default: 1)
// - limit: Items per page for offset pagination (default: 10, max: 100)
// - cursor: Cursor for cursor-based pagination (base64 encoded ID)
// - cursor_limit: Items per cursor page (default: 10, max: 100)
// - pagination_type: Type of pagination (page, cursor) - default: page
//
// Examples:
// - Offset pagination: GET /labels?name=urgent&order_by=name&order_dir=asc&page=1&limit=20
// - Cursor pagination: GET /labels?pagination_type=cursor&cursor_limit=10&cursor=eyJpZCI6MTB9
// - Filter by org: GET /labels?org_id=1&order_by=name
func (h *LabelHandler) List(c *gin.Context) {
	query := h.Client.Label.Query().WithTasks().WithOrganization()

	// Filtering
	if name := c.Query("name"); name != "" {
		query = query.Where(label.NameContains(name))
	}

	if description := c.Query("description"); description != "" {
		query = query.Where(label.DescriptionContains(description))
	}

	if color := c.Query("color"); color != "" {
		query = query.Where(label.ColorEQ(color))
	}

	if orgIDStr := c.Query("org_id"); orgIDStr != "" {
		if orgID, err := strconv.Atoi(orgIDStr); err == nil {
			query = query.Where(label.OrgIDEQ(orgID))
		}
	}

	// Determine pagination type
	paginationType := c.DefaultQuery("pagination_type", "page")

	if paginationType == "cursor" {
		h.listWithCursorPagination(c, query)
		return
	}

	// Default: offset-based pagination
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
		WithOrganization().
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
		OrgID       *int   `json:"org_id"`
	}
	var input LabelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	create := h.Client.Label.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetColor(input.Color)

	if input.OrgID != nil {
		create = create.SetOrgID(*input.OrgID)
	}

	labelObj, err := create.Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create label"})
		return
	}

	// Fetch the created label with edge relationships
	createdLabel, err := h.Client.Label.Query().
		Where(label.ID(labelObj.ID)).
		WithTasks().
		WithOrganization().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created label"})
		return
	}

	c.JSON(http.StatusCreated, createdLabel)
}

func (h *LabelHandler) CreateBulk(c *gin.Context) {
	type LabelInput struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color" binding:"required"`
		OrgID       *int   `json:"org_id"`
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

	bulk := make([]*ent.LabelCreate, len(req.Labels))
	for i, labelInput := range req.Labels {
		create := h.Client.Label.Create().
			SetName(labelInput.Name).
			SetDescription(labelInput.Description).
			SetColor(labelInput.Color)

		if labelInput.OrgID != nil {
			create = create.SetOrgID(*labelInput.OrgID)
		}

		bulk[i] = create
	}

	labels, err := h.Client.Label.CreateBulk(bulk...).Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create labels"})
		return
	}

	// Fetch the created labels with edge relationships
	labelIDs := make([]int, len(labels))
	for i, lbl := range labels {
		labelIDs[i] = lbl.ID
	}

	createdLabels, err := h.Client.Label.Query().
		Where(label.IDIn(labelIDs...)).
		WithTasks().
		WithOrganization().
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

	// Fetch the updated label with edge relationships
	updatedLabel, err := h.Client.Label.Query().
		Where(label.ID(labelObj.ID)).
		WithTasks().
		WithOrganization().
		Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated label"})
		return
	}

	c.JSON(http.StatusOK, updatedLabel)
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

// listWithOffsetPagination handles traditional offset-based pagination
func (h *LabelHandler) listWithOffsetPagination(c *gin.Context, query *ent.LabelQuery) {
	// Apply ordering
	orderOption := h.getOrderOption(c)
	if orderOption == nil {
		return // Error response already sent
	}
	query = query.Order(orderOption)

	// Pagination
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

	// Get total count for pagination info with same filters
	countQuery := h.Client.Label.Query()

	// Apply the same filters as the main query
	if name := c.Query("name"); name != "" {
		countQuery = countQuery.Where(label.NameContains(name))
	}
	if description := c.Query("description"); description != "" {
		countQuery = countQuery.Where(label.DescriptionContains(description))
	}
	if color := c.Query("color"); color != "" {
		countQuery = countQuery.Where(label.ColorEQ(color))
	}
	if orgIDStr := c.Query("org_id"); orgIDStr != "" {
		if orgID, err := strconv.Atoi(orgIDStr); err == nil {
			countQuery = countQuery.Where(label.OrgIDEQ(orgID))
		}
	}

	total, err := countQuery.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := gin.H{
		"data": labels,
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

// listWithCursorPagination handles cursor-based pagination
func (h *LabelHandler) listWithCursorPagination(c *gin.Context, query *ent.LabelQuery) {
	// Apply ordering
	orderOption := h.getOrderOption(c)
	if orderOption == nil {
		return // Error response already sent
	}
	query = query.Order(orderOption)

	// Cursor pagination setup
	limit := 10
	if limitStr := c.Query("cursor_limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse cursor if provided
	if cursor := c.Query("cursor"); cursor != "" {
		cursorData, err := h.decodeCursor(cursor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid cursor format",
			})
			return
		}

		// Apply cursor filter (assuming ID-based cursor)
		if id, ok := cursorData["id"].(float64); ok {
			query = query.Where(label.IDGT(int(id)))
		}
	}

	// Get one extra item to determine if there's a next page
	labels, err := query.Limit(limit + 1).All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	hasNext := len(labels) > limit
	if hasNext {
		labels = labels[:limit] // Remove the extra item
	}

	var nextCursor *string
	if hasNext && len(labels) > 0 {
		lastItem := labels[len(labels)-1]
		cursorStr := h.encodeCursor(map[string]interface{}{
			"id": lastItem.ID,
		})
		nextCursor = &cursorStr
	}

	response := gin.H{
		"data": labels,
		"pagination": gin.H{
			"type":        "cursor",
			"per_page":    limit,
			"has_next":    hasNext,
			"next_cursor": nextCursor,
		},
	}

	c.JSON(http.StatusOK, response)
}

// getOrderOption returns the appropriate order option based on query parameters
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
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order_by field. Valid fields: id, name, description, color, org_id, created_at, updated_at",
		})
		return nil
	}

	return orderOption
}

// encodeCursor encodes cursor data to base64
func (h *LabelHandler) encodeCursor(data map[string]interface{}) string {
	jsonData, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(jsonData)
}

// decodeCursor decodes base64 cursor to data map
func (h *LabelHandler) decodeCursor(cursor string) (map[string]interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}
