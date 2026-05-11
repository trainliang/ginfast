package service

import (
	"testing"
	"time"

	"gin-fast/app/models"
)

func TestEduTermServiceCompile(t *testing.T) {
	setupEduTestDB(t)
	_ = NewEduTermService()
}

func TestEduTermServiceRejectsInvalidDateRange(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduTermService()
	err := svc.Create(contextWithTenant(1), &models.EduTerm{
		Name:      "2026春季",
		StartDate: &models.JSONTime{Time: time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 5, 9, 0, 0, 0, 0, time.UTC)},
		Status:    1,
		TenantID:  1,
	})
	if err == nil {
		t.Fatalf("expected invalid date range to be rejected")
	}
}

func TestEduTermServiceRejectsDuplicateNameWithinTenantAndAllowsCrossTenantSameName(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduTermService()
	ctx1 := contextWithTenant(1)

	if err := db.Create(&models.EduTerm{
		Name:      "2026春季",
		StartDate: &models.JSONTime{Time: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)},
		Status:    1,
		TenantID:  1,
	}).Error; err != nil {
		t.Fatalf("seed tenant 1 term: %v", err)
	}

	err := svc.Create(ctx1, &models.EduTerm{
		Name:      "2026春季",
		StartDate: &models.JSONTime{Time: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)},
		Status:    1,
		TenantID:  1,
	})
	if err == nil {
		t.Fatalf("expected duplicate name within tenant to be rejected")
	}

	err = svc.Create(contextWithTenant(2), &models.EduTerm{
		Name:      "2026春季",
		StartDate: &models.JSONTime{Time: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)},
		Status:    1,
		TenantID:  2,
	})
	if err != nil {
		t.Fatalf("expected cross-tenant same name to be allowed, got %v", err)
	}
}

func TestEduTermServiceSaveClosedDaysReplacesTenantScopedRows(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduTermService()
	ctx := contextWithTenant(1)

	if err := db.Create(&models.EduTerm{Name: "2026春季", Status: 1, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed term: %v", err)
	}
	if err := db.Create(&models.EduTermClosedDay{TermID: 1, ClosedDate: &models.JSONTime{Time: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)}, Reason: "旧数据", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed tenant closed day: %v", err)
	}
	if err := db.Create(&models.EduTermClosedDay{TermID: 1, ClosedDate: &models.JSONTime{Time: time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)}, Reason: "其他租户旧数据", TenantID: 2}).Error; err != nil {
		t.Fatalf("seed other tenant closed day: %v", err)
	}

	err := svc.SaveClosedDays(ctx, 1, 1, []models.EduTermClosedDay{
		{ClosedDate: &models.JSONTime{Time: time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC)}, Reason: "新停课日", TenantID: 1},
	})
	if err != nil {
		t.Fatalf("save closed days: %v", err)
	}

	var rows []models.EduTermClosedDay
	if err := db.Where("tenant_id = ?", 1).Order("id asc").Find(&rows).Error; err != nil {
		t.Fatalf("load tenant closed days: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 tenant-scoped row, got %d", len(rows))
	}
	if rows[0].Reason != "新停课日" {
		t.Fatalf("expected tenant rows to be replaced, got %+v", rows[0])
	}
	var otherCount int64
	if err := db.Model(&models.EduTermClosedDay{}).Where("tenant_id = ?", 2).Count(&otherCount).Error; err != nil {
		t.Fatalf("count other tenant rows: %v", err)
	}
	if otherCount != 1 {
		t.Fatalf("expected other tenant rows to remain, got %d", otherCount)
	}
}

func TestEduTermServiceUpdatePreservesCreatedAtAndRejectsMissingOrCrossTenantTerm(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduTermService()
	ctx1 := contextWithTenant(1)
	createdAt := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

	if err := db.Create(&models.EduTerm{
		BaseModel: models.BaseModel{
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
		Name:      "2026春季",
		StartDate: &models.JSONTime{Time: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)},
		Status:    1,
		TenantID:  1,
	}).Error; err != nil {
		t.Fatalf("seed term: %v", err)
	}

	if err := svc.Update(ctx1, &models.EduTerm{
		BaseModel: models.BaseModel{ID: 1},
		Name:      "2026春季（修订）",
		StartDate: &models.JSONTime{Time: time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC)},
		Status:    2,
		TenantID:  1,
	}); err != nil {
		t.Fatalf("update term: %v", err)
	}

	var updated models.EduTerm
	if err := db.Where("id = ?", 1).First(&updated).Error; err != nil {
		t.Fatalf("load updated term: %v", err)
	}
	if !updated.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected created_at to stay %v, got %v", createdAt, updated.CreatedAt)
	}
	if updated.Name != "2026春季（修订）" || updated.Status != 2 {
		t.Fatalf("expected updated fields to be persisted, got %+v", updated)
	}

	err := svc.Update(ctx1, &models.EduTerm{
		BaseModel: models.BaseModel{ID: 999},
		Name:      "不存在",
		StartDate: &models.JSONTime{Time: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		Status:    1,
		TenantID:  1,
	})
	if err == nil {
		t.Fatalf("expected missing term update to fail")
	}

	if err := db.Create(&models.EduTerm{
		Name:      "2026秋季",
		StartDate: &models.JSONTime{Time: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)},
		Status:    1,
		TenantID:  2,
	}).Error; err != nil {
		t.Fatalf("seed tenant 2 term: %v", err)
	}
	err = svc.Update(ctx1, &models.EduTerm{
		BaseModel: models.BaseModel{ID: 2},
		Name:      "2026秋季（跨租户）",
		StartDate: &models.JSONTime{Time: time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 12, 2, 0, 0, 0, 0, time.UTC)},
		Status:    2,
		TenantID:  1,
	})
	if err == nil {
		t.Fatalf("expected cross-tenant term update to fail")
	}
}

func TestEduTermServiceUpdateAllowsNoopUpdate(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduTermService()
	ctx1 := contextWithTenant(1)
	createdAt := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	startDate := &models.JSONTime{Time: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)}
	endDate := &models.JSONTime{Time: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)}

	if err := db.Create(&models.EduTerm{
		BaseModel: models.BaseModel{
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
		Name:      "2026春季",
		StartDate: startDate,
		EndDate:   endDate,
		Status:    1,
		TenantID:  1,
	}).Error; err != nil {
		t.Fatalf("seed term: %v", err)
	}

	if err := svc.Update(ctx1, &models.EduTerm{
		BaseModel: models.BaseModel{ID: 1},
		Name:      "2026春季",
		StartDate: &models.JSONTime{Time: startDate.Time},
		EndDate:   &models.JSONTime{Time: endDate.Time},
		Status:    1,
		TenantID:  1,
	}); err != nil {
		t.Fatalf("noop update should succeed, got %v", err)
	}

	var updated models.EduTerm
	if err := db.Where("id = ?", 1).First(&updated).Error; err != nil {
		t.Fatalf("load updated term: %v", err)
	}
	if !updated.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected created_at to stay %v, got %v", createdAt, updated.CreatedAt)
	}
	if updated.Name != "2026春季" || updated.Status != 1 {
		t.Fatalf("expected original values to remain, got %+v", updated)
	}
}

func TestEduTermServiceDeleteRejectsMissingCrossTenantAndReferencedTerms(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduTermService()
	ctx1 := contextWithTenant(1)
	ctx2 := contextWithTenant(2)

	if err := db.Create(&models.EduTerm{
		Name:      "2026春季",
		StartDate: &models.JSONTime{Time: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)},
		Status:    1,
		TenantID:  1,
	}).Error; err != nil {
		t.Fatalf("seed tenant 1 term: %v", err)
	}
	if err := db.Create(&models.EduTerm{
		Name:      "2026秋季",
		StartDate: &models.JSONTime{Time: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:   &models.JSONTime{Time: time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)},
		Status:    1,
		TenantID:  2,
	}).Error; err != nil {
		t.Fatalf("seed tenant 2 term: %v", err)
	}
	if err := db.Create(&models.EduTermClosedDay{
		TermID:     1,
		ClosedDate: &models.JSONTime{Time: time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC)},
		Reason:     "新停课日",
		TenantID:   1,
	}).Error; err != nil {
		t.Fatalf("seed closed day: %v", err)
	}

	if err := svc.Delete(ctx1, 1, 999); err == nil {
		t.Fatalf("expected deleting missing term to fail")
	}
	if err := svc.Delete(ctx1, 1, 2); err == nil {
		t.Fatalf("expected deleting cross-tenant term to fail")
	}
	if err := svc.Delete(ctx2, 2, 1); err == nil {
		t.Fatalf("expected deleting cross-tenant term from tenant 2 to fail")
	}
	if err := svc.Delete(ctx1, 1, 1); err == nil {
		t.Fatalf("expected deleting referenced term to fail")
	}
}
