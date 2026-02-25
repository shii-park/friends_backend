package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shii-park/friends/internal/db"
	"github.com/shii-park/friends/internal/sqlc"

	"github.com/shii-park/friends/internal/handler"
	"github.com/shii-park/friends/internal/middleware"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// DBのセットアップ
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbConn, err := db.Setup(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSLMODE"),
	)
	if err != nil {
		log.Fatalf("データベースのセットアップに失敗しました: %v", err)
	}
	log.Println("データベースのセットアップが正常に完了しました")
	defer dbConn.Close()

	// sqlcセットアップ
	queries := sqlc.New(dbConn)
	_ = queries // 未使用によるコンパイルエラー回避のために一時的に代入しています．クエリが必要になった時に，queriesからもらってください．その際にこの行は削除してください．

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
