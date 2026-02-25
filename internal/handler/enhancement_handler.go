package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/service"
)

type EnhancementGinHandler struct {
	svc service.EnhancementService
}

func NewEnhancementGinHandler(svc service.EnhancementService) *EnhancementGinHandler {
	return &EnhancementGinHandler{svc: svc}
}

type enhancementReq struct {
	Times int `json:"times"`
}

func (h *EnhancementGinHandler) Enhance(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID"})
		return
	}

	instanceID, err := uuid.Parse(c.Param("instanceID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instanceID"})
		return
	}

	var req enhancementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// body無しなら 1回強化に寄せたい場合は req.Times=1 でもOK
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if req.Times <= 0 {
		req.Times = 1
	}

	newLevel, cost, err := h.svc.EnhanceCard(c.Request.Context(), userID, instanceID, req.Times)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"instanceID": instanceID.String(),
		"newLevel":   newLevel,
		"cost":       cost,
	})
}
