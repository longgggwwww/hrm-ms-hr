package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/huynhthanhthao/hrm_user_service/proto/user"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/position"
)

type PositionHandler struct {
	Client     *ent.Client
	UserClient *pb.UserServiceClient
}

func NewPositionHandler(client *ent.Client, userClient *pb.UserServiceClient) *PositionHandler {
	return &PositionHandler{
		Client:     client,
		UserClient: userClient,
	}
}

// GetPositions trả về danh sách tất cả chức vụ
func (h *PositionHandler) GetPositions(c *gin.Context) {
	positions, err := h.Client.Position.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch positions"})
		return
	}
	c.JSON(http.StatusOK, positions)
}

// GetPositionByID trả về thông tin chức vụ theo ID
func (h *PositionHandler) GetPositionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position ID"})
		return
	}
	positionObj, err := h.Client.Position.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		return
	}
	c.JSON(http.StatusOK, positionObj)
}

// UpdatePosition cập nhật thông tin chức vụ (chưa implement)
func (h *PositionHandler) UpdatePosition(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// DeletePosition xóa chức vụ theo ID
func (h *PositionHandler) DeletePosition(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position ID"})
		return
	}
	_, err = h.Client.Position.Delete().Where(position.ID(id)).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *PositionHandler) RegisterRoutes(r *gin.Engine) {
	positions := r.Group("/positions")
	{
		positions.GET("", h.GetPositions)
		positions.GET(":id", h.GetPositionByID)
		positions.PUT(":id", h.UpdatePosition)
		positions.DELETE(":id", h.DeletePosition)
	}
}
