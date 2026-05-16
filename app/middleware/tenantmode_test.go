package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-fast/app/global/app"
	"gin-fast/app/global/consts"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequireTenantBusinessModeRejectsPlatformMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(consts.BindContextKeyName, &app.Claims{ClaimsUser: app.ClaimsUser{UserID: 1, TenantID: 0, Mode: "platform"}})
	})
	r.Use(RequireTenantBusinessMode())
	r.GET("/api/edu/students/list", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/edu/students/list", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireTenantBusinessModeAllowsTenantMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(consts.BindContextKeyName, &app.Claims{ClaimsUser: app.ClaimsUser{UserID: 1, TenantID: 2, Mode: "tenant"}})
	})
	r.Use(RequireTenantBusinessMode())
	r.GET("/api/edu/students/list", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/edu/students/list", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
