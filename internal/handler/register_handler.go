package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterHandler はユーザー新規登録を処理するハンドラーです。
// リクエストボディからユーザー名とパスワードを受け取り、
// パスワードをbcryptでハッシュ化したうえでユーザー情報を作成します。
func RegisterHandler(c *gin.Context) {
	var json struct {
		Username string `json:"Username" binding:"required"`
	}

	if err := c.ShouldBindBodyWithJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//TODO: Userタイプは仮実装．別ファイルに定義する
	type User struct {
		ID        uint
		Username  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	user := User{
		Username:  json.Username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	//TODO: ユーザ情報をDBに保存する処理を追加
	//TODO: 下のPrintlnを削除(未使用の変数があるとエラーが出るので置いています)
	fmt.Println(user)

	//TODO: フロントへのレスポンスは話し合って調整する
	c.JSON(http.StatusOK, gin.H{"message": "ユーザ登録が完了しました"})

}
