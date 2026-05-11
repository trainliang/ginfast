package service

import (
	"testing"
	"time"

	"gin-fast/app/models"
)

func TestEduRoomServiceCompile(t *testing.T) {
	setupEduTestDB(t)
	_ = NewEduRoomService()
}

func TestEduRoomServiceRejectsOverlappingWeeklyRules(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduRoomService()
	err := svc.ValidateEduRoomWeeklyRules([]models.EduRoomWeeklyRuleImportRow{
		{RoomID: 1, Weekday: 1, StartTime: "09:00", EndTime: "10:00", Available: 1},
		{RoomID: 1, Weekday: 1, StartTime: "09:30", EndTime: "10:30", Available: 1},
	})
	if err == nil {
		t.Fatalf("expected overlapping weekly rules error")
	}
}

func TestEduRoomServiceAllowsNonOverlappingWeeklyRules(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduRoomService()
	err := svc.ValidateEduRoomWeeklyRules([]models.EduRoomWeeklyRuleImportRow{
		{RoomID: 1, Weekday: 1, StartTime: "09:00", EndTime: "10:00", Available: 1},
		{RoomID: 1, Weekday: 1, StartTime: "10:00", EndTime: "11:00", Available: 1},
	})
	if err != nil {
		t.Fatalf("expected non-overlapping weekly rules to pass, got %v", err)
	}
}

func TestEduRoomServiceRejectsInvalidTimeRange(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduRoomService()
	err := svc.ValidateEduRoomWeeklyRules([]models.EduRoomWeeklyRuleImportRow{
		{RoomID: 1, Weekday: 1, StartTime: "11:00", EndTime: "10:00", Available: 1},
	})
	if err == nil {
		t.Fatalf("expected invalid time range error")
	}
}

func TestEduRoomServiceImportRowsPersistsData(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduRoomService()
	ctx := contextWithTenant(1)

	result, err := svc.ImportRows(ctx, 1, []models.EduRoomImportRow{
		{Name: "场地A", Code: "R001", Capacity: 10},
	})
	if err != nil {
		t.Fatalf("import room rows: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success result")
	}

	var count int64
	if err := db.Model(&models.EduRoom{}).Where("tenant_id = ?", 1).Count(&count).Error; err != nil {
		t.Fatalf("count rooms: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected room persisted, got %d", count)
	}
}

func TestEduRoomServiceRejectsZeroTenantCreate(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduRoomService()
	err := svc.Create(contextWithTenant(0), &models.EduRoom{Name: "场地", Code: "R001", Capacity: 10, TenantID: 0})
	if err == nil {
		t.Fatalf("expected zero tenant create to be rejected")
	}
}

func TestEduRoomServiceRejectsCrossTenantUpdate(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduRoomService()
	if err := db.Create(&models.EduRoom{Name: "场地", Code: "R001", Capacity: 10, TenantID: 0}).Error; err != nil {
		t.Fatalf("seed room: %v", err)
	}
	err := svc.Update(contextWithTenant(1), &models.EduRoom{BaseModel: models.BaseModel{ID: 1}, Name: "场地", Code: "R001", Capacity: 10, TenantID: 1})
	if err == nil {
		t.Fatalf("expected cross tenant update to be rejected")
	}
}

func TestEduRoomServiceRejectsZeroTenantWeeklyRules(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduRoomService()
	_, err := svc.SaveWeeklyRules(contextWithTenant(0), 0, []models.EduRoomWeeklyRuleImportRow{
		{RoomID: 1, Weekday: 1, StartTime: "09:00", EndTime: "10:00", Available: 1},
	})
	if err == nil {
		t.Fatalf("expected zero tenant weekly rules to be rejected")
	}
}

func TestEduRoomServiceImportWeeklyRowsPersistsDataAndRejectsInvalidTenantRoom(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduRoomService()
	ctx := contextWithTenant(1)
	if err := db.Create(&models.EduRoom{Name: "场地A", Code: "R001", Capacity: 10, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed room: %v", err)
	}

	result, err := svc.ImportWeeklyRows(ctx, 1, []models.EduRoomWeeklyRuleImportRow{
		{RoomID: 1, Weekday: 1, StartTime: "09:00", EndTime: "10:00", Available: 1},
	})
	if err != nil {
		t.Fatalf("import weekly rows: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success result")
	}

	if err := db.Create(&models.EduRoom{Name: "别租户", Code: "R002", Capacity: 10, TenantID: 2}).Error; err != nil {
		t.Fatalf("seed other room: %v", err)
	}
	if result2, err := svc.ImportWeeklyRows(ctx, 1, []models.EduRoomWeeklyRuleImportRow{
		{RoomID: 2, Weekday: 1, StartTime: "11:00", EndTime: "12:00", Available: 1},
	}); err == nil {
		t.Fatalf("expected tenant isolation error")
	} else if result2 == nil {
		t.Fatalf("expected result on error")
	}

	var count int64
	if err := db.Model(&models.EduRoomWeeklyRule{}).Where("tenant_id = ?", 1).Count(&count).Error; err != nil {
		t.Fatalf("count weekly rules: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected only first weekly rule persisted, got %d", count)
	}
}

func TestEduRoomServiceImportExceptionRowsPersistsData(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduRoomService()
	ctx := contextWithTenant(1)
	if err := db.Create(&models.EduRoom{Name: "场地A", Code: "R001", Capacity: 10, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed room: %v", err)
	}

	result, err := svc.ImportExceptionRows(ctx, 1, []models.EduRoomExceptionImportRow{
		{RoomID: 1, ExceptionDate: &models.JSONTime{Time: time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)}, Type: "open"},
	})
	if err != nil {
		t.Fatalf("import exception rows: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success result")
	}

	var count int64
	if err := db.Model(&models.EduRoomException{}).Where("tenant_id = ?", 1).Count(&count).Error; err != nil {
		t.Fatalf("count exceptions: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exception persisted, got %d", count)
	}
}

func TestEduRoomServiceImportRowsRollsBackOnError(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduRoomService()
	ctx := contextWithTenant(1)

	result, err := svc.ImportRows(ctx, 1, []models.EduRoomImportRow{
		{Name: "场地A", Code: "R001", Capacity: 10},
		{Name: "错误场地", Code: "R002", Capacity: 0},
	})
	if err == nil {
		t.Fatalf("expected import error")
	}
	if result == nil {
		t.Fatalf("expected result on error")
	}

	var count int64
	if err := db.Model(&models.EduRoom{}).Where("tenant_id = ?", 1).Count(&count).Error; err != nil {
		t.Fatalf("count rooms: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected rollback without partial insert, got %d", count)
	}
}
