package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/service"
)

type GachaStoneGinHandler struct {
	svc *service.GachaStoneService
}

func NewGachaStoneGinHandler(svc *service.GachaStoneService) *GachaStoneGinHandler {
	return &GachaStoneGinHandler{svc: svc}
}

type addGachaStoneRequest struct {
	Amount int `json:"amount"`
}

func (h *GachaStoneGinHandler) Add(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return
	}

	var req addGachaStoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	v, err := h.svc.Add(c.Request.Context(), userID, req.Amount)
	if err != nil {
		// 既存のエラー→HTTP変換があるならそれに寄せると綺麗
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userID":     userID.String(),
		"gachaStone": v,
	})
}
