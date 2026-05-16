package middleware

import (
	"net/http"

	"gin-fast/app/utils/tenanthelper"

	"github.com/gin-gonic/gin"
)

func RequireTenantBusinessMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantCtx, err := tenanthelper.FromGinContext(c)
		if err != nil || tenantCtx.EffectiveTenantID == 0 || tenantCtx.Mode == tenanthelper.ModePlatform {
			c.JSON(http.StatusForbidden, gin.H{"message": "当前处于平台态，禁止访问租户业务数据，请先切换到具体租户"})
			c.Abort()
			return
		}
		c.Next()
	}
}
