package tenanthelper

import (
	"context"
	"testing"

	"gin-fast/app/global/app"
	"gin-fast/app/global/consts"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromContextRequiresTenantForBusiness(t *testing.T) {
	ctx := context.WithValue(context.Background(), consts.BindContextKeyName, &app.Claims{
		ClaimsUser: app.ClaimsUser{UserID: 1, TenantID: 0, Mode: "platform", IsPlatformAdmin: true},
	})

	tenantCtx, err := FromContext(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint(0), tenantCtx.EffectiveTenantID)
	assert.Equal(t, "platform", tenantCtx.Mode)
	assert.True(t, tenantCtx.IsPlatformAdmin)

	_, err = RequireBusinessTenant(ctx)
	require.ErrorContains(t, err, "当前处于平台态")
}

func TestRequireBusinessTenantAcceptsTenantMode(t *testing.T) {
	ctx := context.WithValue(context.Background(), consts.BindContextKeyName, &app.Claims{
		ClaimsUser: app.ClaimsUser{UserID: 2, TenantID: 9, TenantCode: "t9", Mode: "tenant", AuthSource: "password"},
	})

	tenantCtx, err := RequireBusinessTenant(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint(9), tenantCtx.EffectiveTenantID)
	assert.Equal(t, "t9", tenantCtx.TenantCode)
	assert.Equal(t, "tenant", tenantCtx.Mode)
}

func TestFromGinContextUsesClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set(consts.BindContextKeyName, &app.Claims{
		ClaimsUser: app.ClaimsUser{UserID: 3, TenantID: 11, Mode: "impersonation", ActorUserID: 3, IsPlatformAdmin: true},
	})

	tenantCtx, err := FromGinContext(c)
	require.NoError(t, err)
	assert.Equal(t, uint(11), tenantCtx.EffectiveTenantID)
	assert.Equal(t, "impersonation", tenantCtx.Mode)
	assert.True(t, tenantCtx.IsPlatformAdmin)
}
