package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/shii-park/friends/internal/handler"
)

func main() {
	r := gin.Default()

	//TODO: クッキーの秘密鍵や名前の変更
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/register", handler.RegisterHandler)

	//ポート8080番でリッスン
	r.Run()
}
