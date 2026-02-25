package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shii-park/friends/internal/db"
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

	queries := sqlc.New(dbConn)
	r := gin.Default()

	_ = queries // 未使用によるコンパイルエラー回避のために一時的に代入しています．クエリが必要になった時に，queriesからもらってください．その際にこの行は削除してください．

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	//ポート8080番でリッスン
	r.Run()
}
