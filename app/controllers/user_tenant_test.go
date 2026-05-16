package controllers

import (
	"testing"

	"gin-fast/app/global/app"
	"gin-fast/app/global/consts"
	"gin-fast/app/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestEnsureUserInCurrentTenantRejectsCrossTenantUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	app.GormDbMysql = db
	app.ConfigYml = sysDepartmentTestConfig{}
	require.NoError(t, db.Create(&models.User{BaseModel: models.BaseModel{ID: 2}, Username: "u2", TenantID: 2}).Error)

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set(consts.BindContextKeyName, &app.Claims{ClaimsUser: app.ClaimsUser{UserID: 1, TenantID: 1, Mode: "tenant"}})

	err = ensureUserInCurrentTenant(c, 2)
	require.Error(t, err)
}

func TestEnsureUserInCurrentTenantAllowsPlatformAdmin(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	app.GormDbMysql = db
	app.ConfigYml = sysDepartmentTestConfig{}
	require.NoError(t, db.Create(&models.User{BaseModel: models.BaseModel{ID: 2}, Username: "u2", TenantID: 2}).Error)

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set(consts.BindContextKeyName, &app.Claims{ClaimsUser: app.ClaimsUser{UserID: 1, TenantID: 0, Mode: "platform", IsPlatformAdmin: true}})

	err = ensureUserInCurrentTenant(c, 2)
	require.NoError(t, err)
}
