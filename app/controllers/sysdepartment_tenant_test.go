package controllers

import (
	"context"
	"testing"
	"time"

	"gin-fast/app/global/app"
	"gin-fast/app/global/consts"
	"gin-fast/app/models"
	"gin-fast/app/service"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSysDepartmentServiceRejectsCrossTenantUpdate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.SysDepartment{}))
	app.GormDbMysql = db
	app.ConfigYml = sysDepartmentTestConfig{}

	require.NoError(t, db.Create(&models.SysDepartment{BaseModel: models.BaseModel{ID: 1}, Name: "T2", TenantID: 2}).Error)

	ctx := context.WithValue(context.Background(), consts.BindContextKeyName, &app.Claims{
		ClaimsUser: app.ClaimsUser{UserID: 1, TenantID: 1, Mode: "tenant"},
	})

	_, err = service.NewSysDepartmentService().Update(ctx, &models.SysDepartmentUpdateRequest{ID: 1, Name: "hacked"})
	require.Error(t, err)
}

type sysDepartmentTestConfig struct{}

func (sysDepartmentTestConfig) ConfigFileChangeListen(...func()) {}
func (sysDepartmentTestConfig) Get(string) interface{}           { return nil }
func (sysDepartmentTestConfig) GetString(string) string          { return "mysql" }
func (sysDepartmentTestConfig) GetBool(string) bool              { return false }
func (sysDepartmentTestConfig) GetInt(string) int                { return 0 }
func (sysDepartmentTestConfig) GetInt32(string) int32            { return 0 }
func (sysDepartmentTestConfig) GetInt64(string) int64            { return 0 }
func (sysDepartmentTestConfig) GetFloat64(string) float64        { return 0 }
func (sysDepartmentTestConfig) GetDuration(string) time.Duration { return 0 }
func (sysDepartmentTestConfig) GetStringSlice(string) []string   { return nil }
func (sysDepartmentTestConfig) GetUintSlice(string) []uint       { return nil }
func (sysDepartmentTestConfig) Set(string, interface{})          {}
func (sysDepartmentTestConfig) SaveConfig() error                { return nil }
