package handler

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/service"
)

// RegisterHandler はユーザー新規登録を処理するハンドラーです
// ユーザー情報をデータベースに保存し，セッションを保存します．
func RegisterHandler(svc *service.RegisterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json struct {
			Username       string `json:"userName" binding:"required"`
			Icon           string `json:"icon,omitempty"`
			ProfileMessage string `json:"profileMsg,omitempty"`
			Birthmonth     int    `json:"birthmonth,omitempty"`
			Birthday       int    `json:"birthday,omitempty"`
		}

		if err := c.ShouldBindBodyWithJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
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

		// セッション保存
		session := sessions.Default(c)
		session.Set("userID", user.ID)
		if err := session.Save(); err != nil {
			_ = svc.DeleteUser(c.Request.Context(), user.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// TODO: フロントへのレスポンスは話し合って調整する
		c.JSON(http.StatusOK, gin.H{"message": "ユーザー登録が完了しました"})
	}
}
