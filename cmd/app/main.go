package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shii-park/friends/internal/db"
	"github.com/shii-park/friends/internal/sqlc"
)

func main() {
	// DBのセットアップ
	dbConn, err := db.Setup("postgre", "db", 2556, "user", "password", "dbname", "disable")
	if err != nil {
		fmt.Errorf("データベースのセットアップに失敗しました: %w", err)
	}

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
