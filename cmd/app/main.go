package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/shii-park/friends/internal/handler"
	"github.com/shii-park/friends/internal/middleware"
)

func main() {
	r := gin.Default()

	//TODO: クッキーの秘密鍵や名前の変更
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/register", handler.RegisterHandler)

	//保護されたエンドポイントの動作検証(後々削除)
	r.GET("/ping", middleware.AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	//ポート8080番でリッスン
	r.Run()
}
