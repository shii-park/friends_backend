package handler

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/service"
)

const initialCardID = 4
const initialEquipID = 38

// RegisterHandler はユーザー新規登録を処理するハンドラーです
// ユーザー情報をデータベースに保存し，セッションを保存します．
func RegisterHandler(svc *service.RegisterService, storageSvc service.StorageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json struct {
			Username       string `json:"userName" binding:"required"`
			Icon           string `json:"icon,omitempty"`
			ProfileMessage string `json:"profileMsg,omitempty"`
			Birthmonth     int    `json:"birthmonth,omitempty" `
			Birthday       int    `json:"birthday,omitempty" `
		}

		if err := c.ShouldBindBodyWithJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.Birthmonth == 0 {
			json.Birthmonth = 1
		}
		if json.Birthday == 0 {
			json.Birthday = 1
		}

		// DBにユーザーデータを保存
		user, err := svc.RegisterUser(c.Request.Context(), json.Username, json.Icon, json.ProfileMessage, json.Birthmonth, json.Birthday)
		if err != nil {
			if errors.Is(err, errs.ErrUserNameRequired) || errors.Is(err, errs.ErrInvalidBirthday) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
			}
			return
		}

		// 初期カードをDBに保存
		userUUID, err := uuid.Parse(user.ID)
		if err != nil {
			_ = svc.DeleteUser(c.Request.Context(), user.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
			return
		}
		if _, err := storageSvc.AddCard(c.Request.Context(), userUUID, initialCardID); err != nil {
			_ = svc.DeleteUser(c.Request.Context(), user.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
			return
		}
		if _, err := storageSvc.AddCard(c.Request.Context(), userUUID, initialEquipID); err != nil {
			_ = svc.DeleteUser(c.Request.Context(), user.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
			return
		}

		// セッション保存
		session := sessions.Default(c)
		session.Set("userID", user.ID)
		if err := session.Save(); err != nil {
			_ = svc.DeleteUser(c.Request.Context(), user.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"userId":          user.ID,
			"userName":        user.Name,
			"icon":            user.Icon,
			"profileMsg":      user.ProfileMessage,
			"latestLoginDate": user.LatestLoginAt,
			"streakLogin":     user.StreakLogin,
			"registeredDate":  user.RegisteredAt,
			"rp":              user.RankPoint,
			"coin":            user.Coin,
		})
	}
}
