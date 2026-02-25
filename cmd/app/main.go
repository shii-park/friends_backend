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

	// serviceセットアップ
	registerSvc := service.NewRegisterService(queries)

	r := gin.Default()

	//TODO: クッキーの秘密鍵や名前の変更
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 認証なしエンドポイント
	// 新規登録
	r.POST("/register", handler.RegisterHandler(registerSvc))

	// 認証が必要なエンドポイント
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		// TODO:以下は動作検証用エンドポイントなので後で削除
		auth.GET("/ping", testHandler)
		// バトルスタート
		auth.POST("/battle", testHandler)
		// ガチャのラインナップを取得
		auth.GET("/gacha/lineup", testHandler)
		// ガチャを引く
		auth.POST("/gacha/draw", testHandler)
		// ストレージの内容を取得
		auth.GET("/storage", testHandler)
		// ユーザー情報を取得
		auth.GET("/user/:userID/get", testHandler)
		// ユーザー情報を削除
		auth.DELETE("/user/:userID/delete", testHandler)
		// ユーザー名を更新
		auth.PUT("/user/:userID/update", testHandler)
		// カードを強化する
		auth.POST("/card/:cardID/upgrade", testHandler)
		// マッチング部屋に参加
		auth.POST("/matching", testHandler)
	}

	//ポート8080番でリッスン
	r.Run()
}

// TODO: すべてのハンドラーが完成したら以下を削除
func testHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
