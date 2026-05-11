package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"gin-fast/app/global/app"
	"gin-fast/app/global/consts"

	"gorm.io/gorm"
)

type EduImportRowError struct {
	Row    int    `json:"row"`
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type EduImportResult struct {
	Success bool                `json:"success"`
	Errors  []EduImportRowError `json:"errors"`
	Created int                 `json:"created"`
	Updated int                 `json:"updated"`
}

func tenantIDFromContext(ctx context.Context) uint {
	if ctx == nil {
		return 0
	}
	if claims, ok := ctx.Value(consts.BindContextKeyName).(*app.Claims); ok && claims != nil {
		return claims.TenantID
	}
	return 0
}

func requireTenant(db *gorm.DB, tenantID uint) *gorm.DB {
	if db == nil {
		return nil
	}
	if tenantID == 0 {
		return db.Where("1 = 0")
	}
	return db.Where("tenant_id = ?", tenantID)
}

func ensureTenantID(tenantID uint) error {
	if tenantID == 0 {
		return errors.New("当前处于全局租户，禁止维护租户业务数据，请先切换到具体租户")
	}
	return nil
}

func isActiveMemberStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "studying", "paused":
		return true
	default:
		return false
	}
}

func parseHHMM(value string) (int, error) {
	value = strings.TrimSpace(value)
	t, err := time.Parse("15:04", value)
	if err != nil {
		return 0, err
	}
	return t.Hour()*60 + t.Minute(), nil
}
