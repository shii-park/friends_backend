package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shii-park/friends/internal/service"
)

type CollectionGinHandler struct {
	svc service.CollectionService
}

func NewCollectionGinHandler(svc service.CollectionService) *CollectionGinHandler {
	return &CollectionGinHandler{svc: svc}
}

// GET /collection
func (h *CollectionGinHandler) List(c *gin.Context) {
	userID, err := mustUserIDFromSession(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	list, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		writeGinDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entries": list,
		"count":   len(list),
	})
}
