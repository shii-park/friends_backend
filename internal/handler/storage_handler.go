package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/service"
)

type StorageGinHandler struct {
	svc service.StorageService
}

func NewStorageGinHandler(svc service.StorageService) *StorageGinHandler {
	return &StorageGinHandler{svc: svc}
}

type addCardRequest struct {
	CardID int `json:"cardID"`
}

type cardInstanceResponse struct {
	InstanceID string `json:"instanceID"`
	CardID     int    `json:"cardID"`
}

type listCardsResponse struct {
	Cards []cardInstanceResponse `json:"cards"`
	Count int                    `json:"count"`
}

// GET /storage
func (h *StorageGinHandler) ListCards(c *gin.Context) {
	userID, err := mustUserIDFromSession(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	list, err := h.svc.ListCards(c.Request.Context(), userID)
	if err != nil {
		writeGinDomainError(c, err)
		return
	}

	resp := listCardsResponse{
		Cards: make([]cardInstanceResponse, 0, len(list)),
		Count: len(list),
	}
	for _, ci := range list {
		resp.Cards = append(resp.Cards, cardInstanceResponse{
			InstanceID: ci.InstanceID.String(),
			CardID:     ci.CardID,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// POST /storage/cards  body: { "cardID": 1 }
func (h *StorageGinHandler) AddCard(c *gin.Context) {
	userID, err := mustUserIDFromSession(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req addCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	ci, err := h.svc.AddCard(c.Request.Context(), userID, req.CardID)
	if err != nil {
		writeGinDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, cardInstanceResponse{
		InstanceID: ci.InstanceID.String(),
		CardID:     ci.CardID,
	})
}

// DELETE /storage/cards/:instanceID
func (h *StorageGinHandler) RemoveCard(c *gin.Context) {
	userID, err := mustUserIDFromSession(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	instanceID, err := uuid.Parse(c.Param("instanceID"))
	if err != nil || instanceID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instanceID"})
		return
	}

	if err := h.svc.RemoveCard(c.Request.Context(), userID, instanceID); err != nil {
		writeGinDomainError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// --- helpers ---

func mustUserIDFromSession(c *gin.Context) (uuid.UUID, error) {
	v, ok := c.Get("userID")
	if !ok || v == nil {
		return uuid.Nil, errs.ErrInvalidUserID
	}

	switch t := v.(type) {
	case uuid.UUID:
		if t == uuid.Nil {
			return uuid.Nil, errs.ErrInvalidUserID
		}
		return t, nil
	case string:
		id, err := uuid.Parse(t)
		if err != nil || id == uuid.Nil {
			return uuid.Nil, errs.ErrInvalidUserID
		}
		return id, nil
	default:
		return uuid.Nil, errs.ErrInvalidUserID
	}
}

func writeGinDomainError(c *gin.Context, err error) {
	switch err {
	case errs.ErrStorageNil,
		errs.ErrInvalidCardID,
		errs.ErrInvalidInstanceID,
		errs.ErrInvalidUserID:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errs.ErrCardNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
