package tenanthelper

import (
	"context"
	"errors"

	"gin-fast/app/global/app"
	"gin-fast/app/global/consts"
	"gin-fast/app/utils/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	ModePlatform      = "platform"
	ModeTenant        = "tenant"
	ModeImpersonation = "impersonation"
	ModeExternal      = "external"
)

type TenantContext struct {
	ActorUserID       uint
	UserID            uint
	EffectiveTenantID uint
	TenantCode        string
	IsPlatformAdmin   bool
	Mode              string
	AuthSource        string
	ExternalClientID  string
}

func FromGinContext(c *gin.Context) (TenantContext, error) {
	if c == nil {
		return TenantContext{}, errors.New("请求上下文不存在")
	}
	return fromClaims(common.GetClaims(c))
}

func FromContext(ctx context.Context) (TenantContext, error) {
	if ctx == nil {
		return TenantContext{}, errors.New("请求上下文不存在")
	}
	claims, _ := ctx.Value(consts.BindContextKeyName).(*app.Claims)
	return fromClaims(claims)
}

func RequireBusinessTenant(ctx context.Context) (TenantContext, error) {
	tenantCtx, err := FromContext(ctx)
	if err != nil {
		return TenantContext{}, err
	}
	if tenantCtx.EffectiveTenantID == 0 || tenantCtx.Mode == ModePlatform {
		return TenantContext{}, errors.New("当前处于平台态，禁止访问租户业务数据，请先切换到具体租户")
	}
	return tenantCtx, nil
}

func TenantScopeFromContext(ctx context.Context) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tenantCtx, err := RequireBusinessTenant(ctx)
		if err != nil {
			return db.Where("1 = 0")
		}
		return db.Where("tenant_id = ?", tenantCtx.EffectiveTenantID)
	}
}

func fromClaims(claims *app.Claims) (TenantContext, error) {
	if claims == nil {
		return TenantContext{}, errors.New("登录信息不存在")
	}
	mode := claims.Mode
	if mode == "" {
		if claims.TenantID > 0 {
			mode = ModeTenant
		} else {
			mode = ModePlatform
		}
	}
	actorUserID := claims.ActorUserID
	if actorUserID == 0 {
		actorUserID = claims.UserID
	}
	return TenantContext{
		ActorUserID:       actorUserID,
		UserID:            claims.UserID,
		EffectiveTenantID: claims.TenantID,
		TenantCode:        claims.TenantCode,
		IsPlatformAdmin:   claims.IsPlatformAdmin,
		Mode:              mode,
		AuthSource:        claims.AuthSource,
		ExternalClientID:  claims.ExternalClientID,
	}, nil
}
