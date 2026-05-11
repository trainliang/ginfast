package service

import (
	"testing"

	"gin-fast/app/models"
)

func TestEduCourseServiceCompile(t *testing.T) {
	setupEduTestDB(t)
	_ = NewEduCourseService()
}

func TestEduCourseServiceRejectsDuplicateCodeWithinTenant(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduCourseService()
	ctx := contextWithTenant(1)
	if err := db.Create(&models.EduCourse{Name: "A", Code: "C001", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed course: %v", err)
	}
	err := svc.Create(ctx, &models.EduCourse{Name: "B", Code: "C001", TenantID: 1})
	if err == nil {
		t.Fatalf("expected duplicate code error")
	}
}

func TestEduCourseServiceRejectsZeroTenantCreate(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduCourseService()
	err := svc.Create(contextWithTenant(0), &models.EduCourse{Name: "A", Code: "C001", TenantID: 0})
	if err == nil {
		t.Fatalf("expected zero tenant create to be rejected")
	}
}

func TestEduCourseServiceRejectsCrossTenantUpdate(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduCourseService()
	if err := db.Create(&models.EduCourse{Name: "A", Code: "C001", TenantID: 0}).Error; err != nil {
		t.Fatalf("seed course: %v", err)
	}
	err := svc.Update(contextWithTenant(1), &models.EduCourse{BaseModel: models.BaseModel{ID: 1}, Name: "A", Code: "C001", TenantID: 1})
	if err == nil {
		t.Fatalf("expected cross tenant update to be rejected")
	}
}

func TestEduCourseServicePersistsTeachingMode(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduCourseService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)

	course := &models.EduCourse{Name: "A", Code: "C002", DefaultTeachingMode: "online", TenantID: 1}
	if err := svc.Create(ctx, course); err != nil {
		t.Fatalf("create course: %v", err)
	}

	var saved models.EduCourse
	if err := db.Where("code = ? AND tenant_id = ?", "C002", 1).First(&saved).Error; err != nil {
		t.Fatalf("load course: %v", err)
	}
	if saved.DefaultTeachingMode != "online" {
		t.Fatalf("expected defaultTeachingMode online, got %q", saved.DefaultTeachingMode)
	}
	if saved.RequiresRoom != 0 {
		t.Fatalf("expected requiresRoom 0, got %d", saved.RequiresRoom)
	}
}

func TestEduCourseServiceDefaultsEmptyTeachingModeToOffline(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduCourseService()
	ctx := contextWithTenant(1)

	course := &models.EduCourse{Name: "A", Code: "C002", RequiresRoom: -1, TenantID: 1}
	if err := svc.Create(ctx, course); err != nil {
		t.Fatalf("create course: %v", err)
	}

	var saved models.EduCourse
	if err := db.Where("code = ? AND tenant_id = ?", "C002", 1).First(&saved).Error; err != nil {
		t.Fatalf("load course: %v", err)
	}
	if saved.DefaultTeachingMode != "offline" {
		t.Fatalf("expected defaultTeachingMode offline, got %q", saved.DefaultTeachingMode)
	}
	if saved.RequiresRoom != 1 {
		t.Fatalf("expected requiresRoom 1, got %d", saved.RequiresRoom)
	}
}

func TestEduCourseServiceHonorsExplicitRequiresRoom(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduCourseService()
	ctx := contextWithTenant(1)

	course := &models.EduCourse{Name: "A", Code: "C002", DefaultTeachingMode: "online", RequiresRoom: 1, TenantID: 1}
	if err := svc.Create(ctx, course); err != nil {
		t.Fatalf("create course: %v", err)
	}

	var saved models.EduCourse
	if err := db.Where("code = ? AND tenant_id = ?", "C002", 1).First(&saved).Error; err != nil {
		t.Fatalf("load course: %v", err)
	}
	if saved.DefaultTeachingMode != "online" {
		t.Fatalf("expected defaultTeachingMode online, got %q", saved.DefaultTeachingMode)
	}
	if saved.RequiresRoom != 1 {
		t.Fatalf("expected explicit requiresRoom 1, got %d", saved.RequiresRoom)
	}
}

func TestEduCourseServiceUpdatePreservesOmittedTeachingFields(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduCourseService()
	ctx := contextWithTenant(1)
	if err := svc.Create(ctx, &models.EduCourse{Name: "A", Code: "C002", DefaultTeachingMode: "online", RequiresRoom: -1, TenantID: 1}); err != nil {
		t.Fatalf("seed course: %v", err)
	}

	if err := svc.Update(ctx, &models.EduCourse{
		BaseModel:    models.BaseModel{ID: 1},
		Name:         "A2",
		Code:         "C002",
		RequiresRoom: -1,
		TenantID:     1,
	}); err != nil {
		t.Fatalf("update course: %v", err)
	}

	var saved models.EduCourse
	if err := db.Where("id = ?", 1).First(&saved).Error; err != nil {
		t.Fatalf("load course: %v", err)
	}
	if saved.DefaultTeachingMode != "online" {
		t.Fatalf("expected omitted defaultTeachingMode to preserve online, got %q", saved.DefaultTeachingMode)
	}
	if saved.RequiresRoom != 0 {
		t.Fatalf("expected omitted requiresRoom to preserve 0, got %d", saved.RequiresRoom)
	}
}

func TestEduCourseServiceImportExistingPreservesOmittedTeachingFields(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduCourseService()
	ctx := contextWithTenant(1)
	if err := svc.Create(ctx, &models.EduCourse{Name: "A", Code: "C002", DefaultTeachingMode: "online", RequiresRoom: -1, TenantID: 1}); err != nil {
		t.Fatalf("seed course: %v", err)
	}

	result, err := svc.ImportRows(ctx, 1, []models.EduCourseImportRow{{Name: "A2", Code: "C002"}})
	if err != nil {
		t.Fatalf("import course: %v", err)
	}
	if !result.Success || result.Updated != 1 {
		t.Fatalf("expected one updated row, got %+v", result)
	}

	var saved models.EduCourse
	if err := db.Where("code = ? AND tenant_id = ?", "C002", 1).First(&saved).Error; err != nil {
		t.Fatalf("load course: %v", err)
	}
	if saved.DefaultTeachingMode != "online" {
		t.Fatalf("expected omitted import defaultTeachingMode to preserve online, got %q", saved.DefaultTeachingMode)
	}
	if saved.RequiresRoom != 0 {
		t.Fatalf("expected omitted import requiresRoom to preserve 0, got %d", saved.RequiresRoom)
	}
}
