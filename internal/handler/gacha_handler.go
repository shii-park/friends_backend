package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/service"
)

type GachaGinHandler struct {
	svc *service.GachaService
}

func NewGachaGinHandler(svc *service.GachaService) *GachaGinHandler {
	return &GachaGinHandler{svc: svc}
}

type drawRequest struct {
	Count int `json:"count"` // 1 or 10
}

func (h *GachaGinHandler) Draw(c *gin.Context) {
	sess := sessions.Default(c)
	raw := sess.Get("userID")
	userIDStr, _ := raw.(string)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session userID"})
		return
	}

	var req drawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	newStone, results, err := h.svc.Draw(c.Request.Context(), userID, req.Count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userID":     userID.String(),
		"gachaStone": newStone,
		"results":    results,
	})
}

func (h *GachaGinHandler) ListCharacters(c *gin.Context) {
	items, err := h.svc.ListCharacters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
