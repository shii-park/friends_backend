package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/sqlc"
)

// RegisterHandler はユーザー新規登録を処理するハンドラーです
// ユーザー情報をデータベースに保存し，セッションを保存します．
func RegisterHandler(queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json struct {
			Username       string `json:"userName" binding:"required"`
			Icon           string `json:"icon,omitempty"`
			ProfileMessage string `json:"profileMsg,omitempty"`
			Birthmonth     int    `json:"birthmonth" binding:"required"`
			Birthday       int    `json:"birthday" binding:"required"`
		}

		if err := c.ShouldBindBodyWithJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := domain.NewUser(json.Username, json.Icon, json.ProfileMessage, json.Birthmonth, json.Birthday)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ユーザー情報をDBに保存する
		err = queries.CreateUser(c.Request.Context(), sqlc.CreateUserParams{
			UserID:   user.ID,
			UserName: user.Name,
			IconUrl: sql.NullString{
				String: user.Icon,
				Valid:  user.Icon != "",
			},
			ProfileMessage: sql.NullString{
				String: user.ProfileMessage,
				Valid:  user.ProfileMessage != "",
			},
			BirthMonth: sql.NullInt16{
				Int16: int16(user.Birthmonth),
				Valid: true,
			},
			BirthDay: sql.NullInt16{
				Int16: int16(user.Birthday),
				Valid: true,
			},
			LatestLoginDate: sql.NullTime{
				Time:  user.LatestLoginAt,
				Valid: true,
			},
			StreakLoginDays: int32(user.StreakLogin),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// セッション保存
		session := sessions.Default(c)
		session.Set("userID", user.ID)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//TODO: フロントへのレスポンスは話し合って調整する
		c.JSON(http.StatusOK, gin.H{"message": "ユーザ登録が完了しました"})
	}
}
