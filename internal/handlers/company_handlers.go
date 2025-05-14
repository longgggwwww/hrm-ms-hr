package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

type PermGroupHandler struct {
	Client *ent.Client
}

func (h *PermGroupHandler) CreateCompanyGroups(c *gin.Context) {

}

func (h *PermGroupHandler) respondWithError(c *gin.Context, statusCode int, err error) {
	utils.RespondWithError(c, statusCode, err)
}
