package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gin-fast/app/global/app"
	"gin-fast/app/global/consts"
	"gin-fast/app/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupEduTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.EduCourse{},
		&models.EduTerm{},
		&models.EduTermClosedDay{},
		&models.EduBenefitProduct{},
		&models.EduBenefitProductCourse{},
		&models.EduStudentBenefit{},
		&models.EduBenefitLedger{},
		&models.EduBenefitEvent{},
		&models.EduBenefitExternalSync{},
		&models.EduLessonStudentEligibility{},
		&models.EduStudent{},
		&models.EduStudentContact{},
		&models.EduClass{},
		&models.EduClassMember{},
		&models.EduRoom{},
		&models.EduRoomWeeklyRule{},
		&models.EduRoomException{},
		&models.EduTeacherRoleConfig{},
		&models.User{},
		&models.SysRole{},
		&models.SysUserRole{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	app.GormDbMysql = db
	app.ConfigYml = fakeEduConfig{}
	return db
}

type fakeEduConfig struct{}

func (fakeEduConfig) ConfigFileChangeListen(...func()) {}
func (fakeEduConfig) Get(string) interface{}           { return nil }
func (fakeEduConfig) GetString(string) string          { return "mysql" }
func (fakeEduConfig) GetBool(string) bool              { return false }
func (fakeEduConfig) GetInt(string) int                { return 0 }
func (fakeEduConfig) GetInt32(string) int32            { return 0 }
func (fakeEduConfig) GetInt64(string) int64            { return 0 }
func (fakeEduConfig) GetFloat64(string) float64        { return 0 }
func (fakeEduConfig) GetDuration(string) time.Duration { return 0 }
func (fakeEduConfig) GetStringSlice(string) []string   { return nil }
func (fakeEduConfig) GetUintSlice(string) []uint       { return nil }
func (fakeEduConfig) Set(string, interface{})          {}
func (fakeEduConfig) SaveConfig() error                { return nil }

func contextWithTenant(tenantID uint) context.Context {
	return context.WithValue(context.Background(), consts.BindContextKeyName, &app.Claims{ClaimsUser: app.ClaimsUser{TenantID: tenantID}})
}

func seedEduTenantData(t *testing.T, tenantID uint) {
	t.Helper()
	db := app.DB()
	mustCreate := func(value interface{}) {
		t.Helper()
		if err := db.Create(value).Error; err != nil {
			t.Fatalf("seed %T: %v", value, err)
		}
	}

	mustCreate(&models.EduCourse{Name: fmt.Sprintf("课程-%d", tenantID), Code: fmt.Sprintf("C%03d", tenantID), TenantID: tenantID})
	mustCreate(&models.EduRoom{Name: fmt.Sprintf("场地-%d", tenantID), Code: fmt.Sprintf("R%03d", tenantID), Capacity: 10, TenantID: tenantID})
	mustCreate(&models.User{Username: fmt.Sprintf("teacher-%d", tenantID), Password: "pwd", Description: "teacher", TenantID: tenantID})
	mustCreate(&models.SysRole{Name: "教师", Status: 1, TenantID: tenantID})

	var user models.User
	if err := db.Where("username = ?", fmt.Sprintf("teacher-%d", tenantID)).First(&user).Error; err != nil {
		t.Fatalf("load seeded teacher: %v", err)
	}
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", tenantID).First(&role).Error; err != nil {
		t.Fatalf("load seeded role: %v", err)
	}
	mustCreate(&models.SysUserRole{UserID: user.ID, RoleID: role.ID})
}
