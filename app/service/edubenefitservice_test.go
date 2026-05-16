package service

import (
	"fmt"
	"testing"
	"time"

	"gin-fast/app/models"

	"gorm.io/gorm"
)

func TestEduBenefitServiceCheckEligibilityPeriodUnlimited(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduBenefitService()
	ctx := contextWithTenant(1)
	now := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)

	seedBenefitProduct(t, db, 1, 1, 1, "course", "period_unlimited", 0, 30)
	seedStudentBenefit(t, db, 1, 1, 1, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 0, now.Add(-24*time.Hour), now.Add(24*time.Hour))

	result, err := svc.CheckEligibility(ctx, &models.EduBenefitCheckRequest{
		TenantID:  uintPtr(1),
		StudentID: 1,
		CourseID:  1,
		ClassID:   0,
		TeacherID: 0,
		CheckedAt: &models.JSONTime{Time: now},
	})
	if err != nil {
		t.Fatalf("check eligibility: %v", err)
	}
	if !result.Eligible || result.Status != "eligible" || result.StudentBenefitID == 0 {
		t.Fatalf("expected eligible result, got %+v", result)
	}
}

func TestEduBenefitServiceCheckEligibilityExpired(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduBenefitService()
	ctx := contextWithTenant(1)
	now := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)

	seedBenefitProduct(t, db, 1, 1, 1, "course", "period_unlimited", 0, 30)
	seedStudentBenefit(t, db, 1, 1, 1, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 0, now.Add(-48*time.Hour), now.Add(-24*time.Hour))

	result, err := svc.CheckEligibility(ctx, &models.EduBenefitCheckRequest{
		TenantID:  uintPtr(1),
		StudentID: 1,
		CourseID:  1,
		ClassID:   0,
		TeacherID: 0,
		CheckedAt: &models.JSONTime{Time: now},
	})
	if err != nil {
		t.Fatalf("check eligibility: %v", err)
	}
	if result.Eligible || result.Status != "ineligible" || result.ReasonCode != "expired" {
		t.Fatalf("expected expired result, got %+v", result)
	}
}

func TestEduBenefitServiceCheckEligibilityInsufficientCount(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduBenefitService()
	ctx := contextWithTenant(1)
	now := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)

	seedBenefitProduct(t, db, 1, 1, 1, "course", "count_limited", 10, 30)
	seedStudentBenefit(t, db, 1, 1, 1, 1, 1, 0, 0, "course", "count_limited", 2, 2, 0, now.Add(-24*time.Hour), now.Add(24*time.Hour))

	result, err := svc.CheckEligibility(ctx, &models.EduBenefitCheckRequest{
		TenantID:  uintPtr(1),
		StudentID: 1,
		CourseID:  1,
		ClassID:   0,
		TeacherID: 0,
		CheckedAt: &models.JSONTime{Time: now},
	})
	if err != nil {
		t.Fatalf("check eligibility: %v", err)
	}
	if result.Eligible || result.Status != "ineligible" || result.ReasonCode != "insufficient_count" {
		t.Fatalf("expected insufficient_count result, got %+v", result)
	}
}

func TestEduBenefitServiceCheckEligibilityUsesLaterEligibleBenefit(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduBenefitService()
	ctx := contextWithTenant(1)
	now := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)

	seedBenefitProduct(t, db, 1, 1, 1, "course", "period_unlimited", 0, 30)
	seedStudentBenefit(t, db, 1, 1, 1, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 0, now.Add(-72*time.Hour), now.Add(-24*time.Hour))
	seedStudentBenefit(t, db, 1, 2, 1, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 0, now.Add(-24*time.Hour), now.Add(72*time.Hour))

	result, err := svc.CheckEligibility(ctx, &models.EduBenefitCheckRequest{
		TenantID:  uintPtr(1),
		StudentID: 1,
		CourseID:  1,
		CheckedAt: &models.JSONTime{Time: now},
	})
	if err != nil {
		t.Fatalf("check eligibility: %v", err)
	}
	if !result.Eligible || result.StudentBenefitID != 2 {
		t.Fatalf("expected later eligible benefit to be selected, got %+v", result)
	}
}

func TestEduBenefitServiceRepairScheduleEligibilityAfterRenewal(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduBenefitService()
	ctx := contextWithTenant(1)
	now := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	checkedAt := models.JSONTime{Time: now}

	seedBenefitProduct(t, db, 1, 1, 1, "course", "period_unlimited", 0, 30)
	seedStudentBenefit(t, db, 1, 1, 1, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 0, now.Add(-72*time.Hour), now.Add(-24*time.Hour))
	seedEligibility(t, db, 1, 1, 1, 1, 1, 0, "ineligible", "expired", &checkedAt, nil)

	if err := db.Model(&models.EduStudentBenefit{}).Where("id = ? AND tenant_id = ?", 1, 1).Updates(map[string]interface{}{
		"valid_to": &models.JSONTime{Time: now.Add(72 * time.Hour)},
	}).Error; err != nil {
		t.Fatalf("renew student benefit: %v", err)
	}

	result, err := svc.RepairScheduleEligibility(ctx, &models.EduBenefitRepairRequest{
		TenantID:      uintPtr(1),
		StudentID:     flexUintPtr(1),
		CourseID:      flexUintPtr(1),
		EffectiveFrom: &models.JSONTime{Time: now.Add(-time.Hour)},
	})
	if err != nil {
		t.Fatalf("repair schedule eligibility: %v", err)
	}
	if result.CheckedCount != 1 || result.RepairedCount != 1 {
		t.Fatalf("expected one repaired record, got %+v", result)
	}

	var updated models.EduLessonStudentEligibility
	if err := db.Where("id = ?", 1).First(&updated).Error; err != nil {
		t.Fatalf("load repaired eligibility: %v", err)
	}
	if updated.EligibilityStatus != "eligible" || updated.StudentBenefitID == 0 || updated.ResolvedAt == nil {
		t.Fatalf("expected repaired eligibility, got %+v", updated)
	}
}

func TestEduBenefitServiceRepairDefaultsToFutureEligibilityRows(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduBenefitService()
	ctx := contextWithTenant(1)
	past := time.Now().UTC().Add(-24 * time.Hour)
	future := time.Now().UTC().Add(24 * time.Hour)

	seedBenefitProduct(t, db, 1, 1, 1, "course", "period_unlimited", 0, 30)
	seedStudentBenefit(t, db, 1, 1, 1, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 0, past.Add(-24*time.Hour), future.Add(24*time.Hour))
	pastCheckedAt := models.JSONTime{Time: past}
	futureCheckedAt := models.JSONTime{Time: future}
	seedEligibility(t, db, 1, 1, 1, 1, 1, 0, "ineligible", "expired", &pastCheckedAt, nil)
	seedEligibility(t, db, 1, 2, 2, 1, 1, 0, "ineligible", "expired", &futureCheckedAt, nil)

	result, err := svc.RepairScheduleEligibility(ctx, &models.EduBenefitRepairRequest{
		TenantID:  uintPtr(1),
		StudentID: flexUintPtr(1),
		CourseID:  flexUintPtr(1),
	})
	if err != nil {
		t.Fatalf("repair schedule eligibility: %v", err)
	}
	if result.CheckedCount != 1 || result.RepairedCount != 1 {
		t.Fatalf("expected only future row to be repaired, got %+v", result)
	}

	var rows []models.EduLessonStudentEligibility
	if err := db.Order("id asc").Find(&rows).Error; err != nil {
		t.Fatalf("load eligibility rows: %v", err)
	}
	if rows[0].EligibilityStatus != "ineligible" || rows[0].ResolvedAt != nil {
		t.Fatalf("expected past row to remain unresolved, got %+v", rows[0])
	}
	if rows[1].EligibilityStatus != "eligible" || rows[1].ResolvedAt == nil {
		t.Fatalf("expected future row to be repaired, got %+v", rows[1])
	}
}

func uintPtr(v uint) *uint { return &v }

func flexUintPtr(v uint) *models.FlexUint {
	value := models.FlexUint(v)
	return &value
}

func seedBenefitProduct(t *testing.T, db *gorm.DB, tenantID, id, codeID uint, benefitType, calculationMode string, totalCount, validDays int) {
	t.Helper()
	product := &models.EduBenefitProduct{
		BaseModel:       models.BaseModel{ID: id},
		Name:            "权益产品",
		Code:            fmt.Sprintf("B%03d", codeID),
		BenefitType:     benefitType,
		CalculationMode: calculationMode,
		TotalCount:      totalCount,
		ValidDays:       validDays,
		Status:          1,
		TenantID:        tenantID,
	}
	if err := db.Create(product).Error; err != nil {
		t.Fatalf("seed benefit product: %v", err)
	}
}

func seedStudentBenefit(t *testing.T, db *gorm.DB, tenantID, id, studentID, productID, courseID, classID, teacherID uint, benefitType, calculationMode string, totalCount, usedCount, remainingCount int, validFrom, validTo time.Time) {
	t.Helper()
	row := &models.EduStudentBenefit{
		BaseModel:       models.BaseModel{ID: id},
		StudentID:       studentID,
		ProductID:       productID,
		BenefitType:     benefitType,
		CalculationMode: calculationMode,
		CourseID:        courseID,
		ClassID:         classID,
		TeacherID:       teacherID,
		TotalCount:      totalCount,
		UsedCount:       usedCount,
		RemainingCount:  remainingCount,
		Status:          1,
		TenantID:        tenantID,
		ValidFrom:       &models.JSONTime{Time: validFrom},
		ValidTo:         &models.JSONTime{Time: validTo},
	}
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("seed student benefit: %v", err)
	}
}

func seedEligibility(t *testing.T, db *gorm.DB, tenantID, id, lessonID, studentID, courseID, classID uint, status, reason string, checkedAt, resolvedAt *models.JSONTime) {
	t.Helper()
	row := &models.EduLessonStudentEligibility{
		BaseModel:         models.BaseModel{ID: id},
		LessonID:          lessonID,
		StudentID:         studentID,
		CourseID:          courseID,
		ClassID:           classID,
		EligibilityStatus: status,
		ReasonCode:        reason,
		CheckedAt:         checkedAt,
		ResolvedAt:        resolvedAt,
		TenantID:          tenantID,
	}
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("seed eligibility: %v", err)
	}
}
