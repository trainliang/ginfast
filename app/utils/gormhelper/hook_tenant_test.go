package gormhelper

import (
	"context"
	"testing"

	"gin-fast/app/global/app"
	"gin-fast/app/global/consts"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type hookTenantRecord struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	TenantID  uint
	CreatedBy uint
}

func setupHookTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if err := db.AutoMigrate(&hookTenantRecord{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// Register the same hook as production.
	_ = db.Callback().Create().Before("gorm:before_create").Register("CreateBeforeHook", CreateBeforeHook)
	return db
}

func contextWithTenant(tenantID uint) context.Context {
	return context.WithValue(context.Background(), consts.BindContextKeyName, &app.Claims{
		ClaimsUser: app.ClaimsUser{TenantID: tenantID},
	})
}

func TestCreateBeforeHook_BatchCreateSetsTenantID(t *testing.T) {
	db := setupHookTestDB(t)
	ctx := contextWithTenant(123)

	rows := []hookTenantRecord{
		{Name: "A"},
		{Name: "B"},
	}
	if err := db.WithContext(ctx).Create(&rows).Error; err != nil {
		t.Fatalf("batch create: %v", err)
	}
	for i := range rows {
		if rows[i].TenantID != 123 {
			t.Fatalf("expected TenantID=123, got %d (idx=%d)", rows[i].TenantID, i)
		}
	}
}

func TestCreateBeforeHook_SingleCreateRejectsCrossTenantID(t *testing.T) {
	db := setupHookTestDB(t)
	ctx := contextWithTenant(123)

	row := &hookTenantRecord{Name: "A", TenantID: 999}
	err := db.WithContext(ctx).Create(row).Error
	if err == nil {
		t.Fatalf("expected cross tenant create to fail")
	}
}

func TestCreateBeforeHook_SingleCreateAllowsMatchingTenantID(t *testing.T) {
	db := setupHookTestDB(t)
	ctx := contextWithTenant(123)

	row := &hookTenantRecord{Name: "A", TenantID: 123}
	if err := db.WithContext(ctx).Create(row).Error; err != nil {
		t.Fatalf("create with matching tenant: %v", err)
	}
	if row.TenantID != 123 {
		t.Fatalf("expected TenantID=123, got %d", row.TenantID)
	}
}

func TestCreateBeforeHook_SingleCreateRejectsPlatformMode(t *testing.T) {
	db := setupHookTestDB(t)
	ctx := contextWithTenant(0)

	row := &hookTenantRecord{Name: "A"}
	err := db.WithContext(ctx).Create(row).Error
	if err == nil {
		t.Fatalf("expected platform mode create to fail")
	}
}
