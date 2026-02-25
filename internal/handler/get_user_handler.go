package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shii-park/friends/internal/service"
)

// GetUserHandler はクエリパラメータで受け取ったユーザーIDのユーザー情報を返すハンドラーです
func GetUserHandler(svc *service.GetUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userIDが指定されていません"})
			return
		}

		user, err := svc.GetUser(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"userId":          user.ID,
			"userName":        user.Name,
			"icon":            user.Icon,
			"profileMsg":      user.ProfileMessage,
			"birthmonth":      user.Birthmonth,
			"birthday":        user.Birthday,
			"registeredDate":  user.RegisteredAt,
			"latestLoginDate": user.LatestLoginAt,
			"streakLogin":     user.StreakLogin,
			"rp":              user.RankPoint,
			"coin":            user.Coin,
			"gachaStone":      user.GachaStone,
		})
	}
}
