package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gin-fast/app/global/app"
	"gin-fast/app/utils/response"

	"github.com/gin-gonic/gin"
)

func TestGetVerifyImgStringReturnsDisabledWhenCaptchaClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalConfig := app.ConfigYml
	originalResponse := app.Response
	t.Cleanup(func() {
		app.ConfigYml = originalConfig
		app.Response = originalResponse
	})

	app.ConfigYml = testYmlConfig{bools: map[string]bool{"captcha.open": false}}
	app.Response = response.NewResponseHandler()

	router := gin.New()
	router.GET("/captcha/verify", NewAuthController().GetVerifyImgString)

	req := httptest.NewRequest(http.MethodGet, "/captcha/verify", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	want := `"enabled":false`
	if !containsJSONField(w.Body.String(), want) {
		t.Fatalf("expected response to include %s, got %s", want, w.Body.String())
	}
}

type testYmlConfig struct {
	bools map[string]bool
}

func (c testYmlConfig) ConfigFileChangeListen(fns ...func())     {}
func (c testYmlConfig) Get(keyName string) interface{}           { return nil }
func (c testYmlConfig) GetString(keyName string) string          { return "" }
func (c testYmlConfig) GetBool(keyName string) bool              { return c.bools[keyName] }
func (c testYmlConfig) GetInt(keyName string) int                { return 0 }
func (c testYmlConfig) GetInt32(keyName string) int32            { return 0 }
func (c testYmlConfig) GetInt64(keyName string) int64            { return 0 }
func (c testYmlConfig) GetFloat64(keyName string) float64        { return 0 }
func (c testYmlConfig) GetDuration(keyName string) time.Duration { return 0 }
func (c testYmlConfig) GetStringSlice(keyName string) []string   { return nil }
func (c testYmlConfig) GetUintSlice(keyName string) []uint       { return nil }
func (c testYmlConfig) Set(keyName string, value interface{})    {}
func (c testYmlConfig) SaveConfig() error                        { return nil }

func containsJSONField(body string, field string) bool {
	for i := 0; i+len(field) <= len(body); i++ {
		if body[i:i+len(field)] == field {
			return true
		}
	}
	return false
}
