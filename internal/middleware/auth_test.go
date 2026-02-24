package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shii-park/friends/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(userID any) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	store := cookie.NewStore([]byte("test-secret"))
	r.Use(sessions.Sessions("mysession", store))

	// セッションにuserIDをセットするエンドポイント（テスト用）
	r.GET("/set-session", func(c *gin.Context) {
		session := sessions.Default(c)
		if userID != nil {
			session.Set("userID", userID)
			err := session.Save()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			}
		}
		c.Status(http.StatusOK)
	})

	r.GET("/protected", middleware.AuthRequired(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r
}

func TestAuthRequired_正常系(t *testing.T) {
	r := setupRouter("user-123")

	// セッションCookieを取得
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/set-session", nil)
	if err != nil {
		t.Logf("正常系テストの1つ目のhttpリクエストの作成失敗: %v", err)
	}
	r.ServeHTTP(w, req)
	cookie := w.Result().Header.Get("Set-Cookie")

	// 保護されたエンドポイントにCookieつきでアクセス
	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Logf("正常系テストの2つ目のhttpリクエストの作成失敗: %v", err)
	}
	req2.Header.Set("Cookie", cookie)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestAuthRequired_セッション情報なし(t *testing.T) {
	r := setupRouter(nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Logf("セッション情報なしテストのhttpリクエストの作成失敗: %v", err)
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
