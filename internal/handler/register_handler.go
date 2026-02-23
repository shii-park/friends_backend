package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shii-park/friends/internal/domain"
)

// RegisterHandler はユーザー新規登録を処理するハンドラーです。
func RegisterHandler(c *gin.Context) {
	var json struct {
		Username       string    `json:"userName" binding:"required"`
		Icon           string    `json:"icon,omitempty"`
		ProfileMessage string    `json:"profileMsg,omitempty"`
		Birthday       time.Time `json:"birthday" binding:"required"`
	}

	if err := c.ShouldBindBodyWithJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//日付の検証(未来の日付であればここでエラー)
	now := time.Now()
	if json.Birthday.After(now) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "誕生日が正しくありません"})
		return
	}

	user, err := domain.NewUser(json.Username, json.Icon, json.ProfileMessage, json.Birthday)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//TODO: ユーザー情報をDBに保存する処理を追加
	//TODO: 下のPrintlnを削除(未使用の変数があるとエラーが出るので置いています)
	fmt.Println(user)

	//TODO: フロントへのレスポンスは話し合って調整する
	c.JSON(http.StatusOK, gin.H{"message": "ユーザ登録が完了しました"})

}
