package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/service"
)

type CoinGinHandler struct {
	svc *service.CoinService
}

func NewCoinGinHandler(svc *service.CoinService) *CoinGinHandler {
	return &CoinGinHandler{svc: svc}
}

type addCoinRequest struct {
	Amount int `json:"amount"`
}

func (h *CoinGinHandler) Add(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return
	}

	var req addCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	v, err := h.svc.Add(c.Request.Context(), userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userID": userID.String(),
		"coin":   v,
	})
}
