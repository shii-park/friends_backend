package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthRequiredはユーザーが認証済であるかを検証するミドルウェアです
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userID") //keyはRegisterHandler()と統一する
		if userID == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
