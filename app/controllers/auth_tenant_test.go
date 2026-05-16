package controllers

import (
	"testing"

	"gin-fast/app/global/app"

	"github.com/stretchr/testify/assert"
)

func TestBuildRefreshClaimsUserPreservesTenantContext(t *testing.T) {
	refreshClaims := &app.RefreshTokenClaims{
		ClaimsUser: app.ClaimsUser{
			UserID:          1,
			Username:        "old",
			TenantID:        8,
			TenantCode:      "tenant8",
			Mode:            "impersonation",
			AuthSource:      "password",
			IsPlatformAdmin: true,
			ActorUserID:     1,
		},
	}

	next := buildRefreshClaimsUser(refreshClaims, "admin")
	assert.Equal(t, uint(1), next.UserID)
	assert.Equal(t, "admin", next.Username)
	assert.Equal(t, uint(8), next.TenantID)
	assert.Equal(t, "tenant8", next.TenantCode)
	assert.Equal(t, "impersonation", next.Mode)
	assert.Equal(t, "password", next.AuthSource)
	assert.True(t, next.IsPlatformAdmin)
	assert.Equal(t, uint(1), next.ActorUserID)
}
