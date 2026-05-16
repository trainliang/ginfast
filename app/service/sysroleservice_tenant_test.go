package service

import (
	"testing"
	"time"

	"gin-fast/app/global/app"
	"gin-fast/app/global/consts"
	"gin-fast/app/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSysRoleServiceRejectsCrossTenantUpdate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.SysRole{}))
	app.GormDbMysql = db
	app.ConfigYml = sysRoleTestConfig{}

	require.NoError(t, db.Create(&models.SysRole{BaseModel: models.BaseModel{ID: 1}, Name: "T2", TenantID: 2}).Error)

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set(consts.BindContextKeyName, &app.Claims{ClaimsUser: app.ClaimsUser{UserID: 1, TenantID: 1, Mode: "tenant"}})

	_, err = NewSysRoleService().Update(c, models.SysRoleUpdateRequest{ID: 1, Name: "hacked"})
	require.Error(t, err)
}

type sysRoleTestConfig struct{}

func (sysRoleTestConfig) ConfigFileChangeListen(...func()) {}
func (sysRoleTestConfig) Get(string) interface{}           { return nil }
func (sysRoleTestConfig) GetString(string) string          { return "mysql" }
func (sysRoleTestConfig) GetBool(string) bool              { return false }
func (sysRoleTestConfig) GetInt(string) int                { return 0 }
func (sysRoleTestConfig) GetInt32(string) int32            { return 0 }
func (sysRoleTestConfig) GetInt64(string) int64            { return 0 }
func (sysRoleTestConfig) GetFloat64(string) float64        { return 0 }
func (sysRoleTestConfig) GetDuration(string) time.Duration { return 0 }
func (sysRoleTestConfig) GetStringSlice(string) []string   { return nil }
func (sysRoleTestConfig) GetUintSlice(string) []uint       { return nil }
func (sysRoleTestConfig) Set(string, interface{})          {}
func (sysRoleTestConfig) SaveConfig() error                { return nil }
