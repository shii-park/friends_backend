package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shii-park/friends/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler_正常系(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	//正常なリクエストボディ
	body := `{"Username":"testuser","Password":"testpass123"}`
	c.Request = httptest.NewRequest("POST", "/register", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.RegisterHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ユーザ登録が完了しました")
}
