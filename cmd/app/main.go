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
	"github.com/shii-park/friends/internal/handler"
	"github.com/shii-park/friends/internal/middleware"
	"github.com/shii-park/friends/internal/repository"
	"github.com/shii-park/friends/internal/service"
	"github.com/shii-park/friends/internal/sqlc"
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

	// Storage
	storageRepo := repository.NewStorageRepository(queries)
	storageSvc := service.NewStorageService(storageRepo)
	storageHandler := handler.NewStorageGinHandler(storageSvc)

	// Enhancement
	enhRepo := repository.NewEnhancementRepository(dbConn, queries)
	enhSvc := service.NewEnhancementService(enhRepo)
	enhHandler := handler.NewEnhancementGinHandler(enhSvc)

	r := gin.Default()

	// TODO: クッキーの秘密鍵や名前の変更
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 認証なしエンドポイント
	r.POST("/register", handler.RegisterHandler)

	// 認証が必要なエンドポイント
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		// TODO:以下は動作検証用エンドポイントなので後で削除
		auth.GET("/ping", testHandler)
		auth.POST("/battle", testHandler)
		auth.GET("/gacha/lineup", testHandler)
		auth.POST("/gacha/draw", testHandler)

		auth.GET("/storage", storageHandler.ListCards)
		auth.POST("/storage/cards", storageHandler.AddCard)
		auth.DELETE("/storage/cards/:instanceID", storageHandler.RemoveCard)

		auth.GET("/user/:userID/get", testHandler)
		auth.DELETE("/user/:userID/delete", testHandler)
		auth.PUT("/user/:userID/update", testHandler)
		auth.POST("/card/:instanceID/enhacement", enhHandler.Enhance)
		auth.POST("/matching", testHandler)
	}

	// ポート8080番でリッスン
	r.Run()
}

// TODO: すべてのハンドラーが完成したら以下を削除
func testHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
