package service

import (
	"testing"
	"time"

	"gin-fast/app/models"

	"gorm.io/gorm"
)

func TestEduStatisticsServiceCompile(t *testing.T) {
	setupEduTestDB(t)
	_ = NewEduStatisticsService()
}

func TestEduStatisticsServiceTeacherMinutesCountsScheduledAndCompletedOnly(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "class",
		LessonDate:   datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:00",
		EndTime:      "10:30",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		Status:       "scheduled",
	})
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "class",
		LessonDate:   datePtr(time.Date(2026, 5, 13, 0, 0, 0, 0, time.UTC)),
		StartTime:    "10:00",
		EndTime:      "11:00",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		Status:       "completed",
	})
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "class",
		LessonDate:   datePtr(time.Date(2026, 5, 14, 0, 0, 0, 0, time.UTC)),
		StartTime:    "11:00",
		EndTime:      "12:00",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		Status:       "canceled",
	})

	resp, err := svc.ScheduleStatistics(ctx, &models.EduStatisticsRangeRequest{TeacherID: uintPtr(1)})
	if err != nil {
		t.Fatalf("schedule statistics: %v", err)
	}
	if len(resp.TeacherLessons) != 1 {
		t.Fatalf("expected 1 teacher summary, got %+v", resp.TeacherLessons)
	}
	got := resp.TeacherLessons[0]
	if got.LessonCount != 2 || got.TotalMinutes != 150 {
		t.Fatalf("expected scheduled/completed only with 150 minutes, got %+v", got)
	}
}

func TestEduStatisticsServiceRoomUsageRequiresRoomAndRoomID(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "class",
		LessonDate:   datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:00",
		EndTime:      "10:00",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		Status:       "scheduled",
	})
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "class",
		LessonDate:   datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
		StartTime:    "10:00",
		EndTime:      "11:00",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 0,
		RoomID:       1,
		Status:       "scheduled",
	})
	if err := db.Model(&models.EduLesson{}).Where("start_time = ? AND teacher_id = ? AND tenant_id = ?", "10:00", 1, 1).
		Update("requires_room", 0).Error; err != nil {
		t.Fatalf("force requires_room=0: %v", err)
	}
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "class",
		LessonDate:   datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
		StartTime:    "11:00",
		EndTime:      "12:00",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       0,
		Status:       "scheduled",
	})

	resp, err := svc.ScheduleStatistics(ctx, &models.EduStatisticsRangeRequest{RoomID: uintPtr(1)})
	if err != nil {
		t.Fatalf("schedule statistics: %v", err)
	}
	if len(resp.RoomUsages) != 1 {
		t.Fatalf("expected 1 room summary, got %+v", resp.RoomUsages)
	}
	got := resp.RoomUsages[0]
	if got.LessonCount != 1 || got.TotalMinutes != 60 {
		t.Fatalf("expected only requires_room=1 and room_id!=0, got %+v", got)
	}
}

func TestEduStatisticsServiceClassAndOneToOneSeparated(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "class",
		LessonDate:   datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:00",
		EndTime:      "10:00",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		Status:       "scheduled",
	})
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "one_to_one",
		LessonDate:   datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
		StartTime:    "10:00",
		EndTime:      "11:00",
		StudentID:    1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "online",
		RequiresRoom: 0,
		RoomID:       0,
		Status:       "scheduled",
	})

	resp, err := svc.ScheduleStatistics(ctx, &models.EduStatisticsRangeRequest{})
	if err != nil {
		t.Fatalf("schedule statistics: %v", err)
	}
	if len(resp.ClassLessons) != 1 || resp.ClassLessons[0].LessonCount != 1 {
		t.Fatalf("expected one class lesson summary, got %+v", resp.ClassLessons)
	}
	if len(resp.OneToOneLessons) != 1 || resp.OneToOneLessons[0].LessonCount != 1 {
		t.Fatalf("expected one one-to-one summary, got %+v", resp.OneToOneLessons)
	}
}

func TestEduStatisticsServiceBenefitExpiringSoonUsesNext30Days(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	now := time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)
	seedBenefitProduct(t, db, 1, 1, 1, "course", "period_unlimited", 0, 30)
	seedStudentBenefit(t, db, 1, 1, 1, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 10, now.Add(-24*time.Hour), now.Add(15*24*time.Hour))
	seedStudentBenefit(t, db, 1, 2, 2, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 10, now.Add(-24*time.Hour), now.Add(31*24*time.Hour))
	seedStudentBenefit(t, db, 1, 3, 3, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 10, now.Add(-48*time.Hour), now.Add(-1*time.Hour))

	resp, err := svc.BenefitStatistics(ctx, &models.EduStatisticsRangeRequest{
		StartDate: datePtr(now),
		EndDate:   datePtr(now.Add(30 * 24 * time.Hour)),
	})
	if err != nil {
		t.Fatalf("benefit statistics: %v", err)
	}
	if resp.ExpiringSoonCount != 1 {
		t.Fatalf("expected only next 30 days benefit, got %+v", resp)
	}
}

func TestEduStatisticsServiceExternalSyncGroupsByStatus(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedExternalSync(t, db, 1, 1, 1, "pending")
	seedExternalSync(t, db, 1, 2, 1, "success")
	seedExternalSync(t, db, 1, 3, 1, "failed")
	seedExternalSync(t, db, 1, 4, 1, "retrying")
	seedExternalSync(t, db, 1, 5, 1, "canceled")

	resp, err := svc.ExternalSyncStatistics(ctx, &models.EduStatisticsRangeRequest{})
	if err != nil {
		t.Fatalf("external sync statistics: %v", err)
	}
	if len(resp.StatusCounts) != 5 {
		t.Fatalf("expected 5 status groups, got %+v", resp.StatusCounts)
	}
}

func TestEduStatisticsServiceConflictOverridesGroupedByType(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedScheduleConflictOverride(t, db, 1, 1, "teacher")
	seedScheduleConflictOverride(t, db, 1, 2, "room")
	seedScheduleConflictOverride(t, db, 1, 3, "room")

	resp, err := svc.ScheduleStatistics(ctx, &models.EduStatisticsRangeRequest{})
	if err != nil {
		t.Fatalf("schedule statistics: %v", err)
	}
	if len(resp.ConflictOverrides) != 2 {
		t.Fatalf("expected grouped conflict overrides, got %+v", resp.ConflictOverrides)
	}
}

func TestEduStatisticsServiceScheduleIncludesEligibilityWarningsAndIneligible(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedLessonEligibility(t, db, 1, 1, 1, 1, 1, "warn", "about_to_expire")
	seedLessonEligibility(t, db, 1, 2, 2, 1, 1, "ineligible", "no_benefit")
	seedLessonEligibility(t, db, 1, 3, 3, 1, 1, "eligible", "ok")

	resp, err := svc.ScheduleStatistics(ctx, &models.EduStatisticsRangeRequest{})
	if err != nil {
		t.Fatalf("schedule statistics: %v", err)
	}
	if len(resp.EligibilityDetails) != 2 {
		t.Fatalf("expected 2 eligibility anomalies, got %+v", resp.EligibilityDetails)
	}
}

func TestEduStatisticsServiceBenefitAggregatesCourseCountsAndZeroRemaining(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	now := time.Now().UTC()
	seedBenefitProduct(t, db, 1, 1, 1, "course", "count_limited", 10, 30)
	seedBenefitProduct(t, db, 1, 2, 2, "course", "period_unlimited", 0, 30)
	seedStudentBenefit(t, db, 1, 1, 1, 1, 1, 0, 0, "course", "count_limited", 10, 10, 0, now.Add(-24*time.Hour), now.Add(10*24*time.Hour))
	seedStudentBenefit(t, db, 1, 2, 2, 2, 2, 0, 0, "course", "period_unlimited", 0, 0, 0, now.Add(-24*time.Hour), now.Add(40*24*time.Hour))
	seedStudentBenefit(t, db, 1, 3, 3, 1, 1, 0, 0, "course", "count_limited", 10, 1, 9, now.Add(-24*time.Hour), now.Add(5*24*time.Hour))

	resp, err := svc.BenefitStatistics(ctx, &models.EduStatisticsRangeRequest{})
	if err != nil {
		t.Fatalf("benefit statistics: %v", err)
	}
	if resp.ZeroRemainingCount != 1 {
		t.Fatalf("expected 1 zero remaining count benefit, got %+v", resp)
	}
	if len(resp.CourseBenefitCounts) != 2 {
		t.Fatalf("expected 2 course aggregates, got %+v", resp.CourseBenefitCounts)
	}
}

func TestEduStatisticsServiceExternalSyncIncludesDueRetryCount(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedExternalSyncWithRetry(t, db, 1, 1, 1, "retrying", time.Now().UTC().Add(-time.Hour))
	seedExternalSyncWithRetry(t, db, 1, 2, 1, "retrying", time.Now().UTC().Add(time.Hour))
	seedExternalSyncWithRetry(t, db, 1, 3, 1, "failed", time.Time{})

	resp, err := svc.ExternalSyncStatistics(ctx, &models.EduStatisticsRangeRequest{})
	if err != nil {
		t.Fatalf("external sync statistics: %v", err)
	}
	if resp.RetryingCount != 2 {
		t.Fatalf("expected retrying count 2, got %+v", resp)
	}
	if resp.FailedCount != 1 {
		t.Fatalf("expected failed count 1, got %+v", resp)
	}
	if resp.DueRetryCount != 1 {
		t.Fatalf("expected due retry count 1, got %+v", resp)
	}
}

func TestEduStatisticsServiceScheduleIncludesEligibilityAndConflictDetails(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedEligibility(t, db, 1, 1, 11, 1, 1, 1, "warn", "class_policy", &models.JSONTime{Time: time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)}, nil)
	seedScheduleConflictOverride(t, db, 1, 1, "teacher")
	seedScheduleConflictOverride(t, db, 1, 2, "room")

	resp, err := svc.ScheduleStatistics(ctx, &models.EduStatisticsRangeRequest{})
	if err != nil {
		t.Fatalf("schedule statistics: %v", err)
	}
	if len(resp.EligibilityDetails) != 1 || resp.EligibilityDetails[0].EligibilityStatus != "warn" {
		t.Fatalf("expected eligibility details, got %+v", resp.EligibilityDetails)
	}
	if len(resp.ConflictDetails) != 2 {
		t.Fatalf("expected conflict details, got %+v", resp.ConflictDetails)
	}
}

func TestEduStatisticsServiceBenefitCountsByCourse(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	now := time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)
	seedBenefitProduct(t, db, 1, 1, 1, "course", "period_unlimited", 0, 30)
	seedStudentBenefit(t, db, 1, 1, 1, 1, 1, 0, 0, "course", "period_unlimited", 10, 0, 10, now.Add(-24*time.Hour), now.Add(15*24*time.Hour))
	seedStudentBenefit(t, db, 1, 2, 2, 1, 2, 0, 0, "course", "period_unlimited", 10, 0, 10, now.Add(-24*time.Hour), now.Add(15*24*time.Hour))

	resp, err := svc.BenefitStatistics(ctx, &models.EduStatisticsRangeRequest{})
	if err != nil {
		t.Fatalf("benefit statistics: %v", err)
	}
	if len(resp.CourseBenefitCounts) != 2 {
		t.Fatalf("expected course benefit counts, got %+v", resp.CourseBenefitCounts)
	}
}

func TestEduStatisticsServiceExternalSyncCountsFailedRetryingAndDueRetry(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStatisticsService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	now := time.Now().UTC()
	seedExternalSync(t, db, 1, 1, 1, "pending")
	seedExternalSync(t, db, 1, 2, 1, "success")
	seedExternalSync(t, db, 1, 3, 1, "failed")
	seedExternalSyncWithRetry(t, db, 1, 4, 1, "retrying", now.Add(-time.Hour))
	seedExternalSyncWithRetry(t, db, 1, 5, 1, "retrying", now.Add(time.Hour))

	resp, err := svc.ExternalSyncStatistics(ctx, &models.EduStatisticsRangeRequest{})
	if err != nil {
		t.Fatalf("external sync statistics: %v", err)
	}
	if resp.FailedCount != 1 || resp.RetryingCount != 2 || resp.DueRetryCount != 1 {
		t.Fatalf("expected failed/retrying/dueRetry counts, got %+v", resp)
	}
}

func seedExternalSync(t *testing.T, db *gorm.DB, tenantID, id, studentBenefitID uint, status string) {
	t.Helper()
	row := &models.EduBenefitExternalSync{
		BaseModel:        models.BaseModel{ID: id},
		StudentBenefitID: studentBenefitID,
		StudentID:        studentBenefitID,
		ProviderCode:     "provider",
		IdempotencyKey:   "sync-" + status,
		Status:           status,
		TenantID:         tenantID,
	}
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("seed external sync: %v", err)
	}
}

func seedExternalSyncWithRetry(t *testing.T, db *gorm.DB, tenantID, id, studentBenefitID uint, status string, nextRetryAt time.Time) {
	t.Helper()
	row := &models.EduBenefitExternalSync{
		BaseModel:        models.BaseModel{ID: id},
		StudentBenefitID: studentBenefitID,
		StudentID:        studentBenefitID,
		ProviderCode:     "provider",
		IdempotencyKey:   "sync-" + status,
		Status:           status,
		NextRetryAt:      &models.JSONTime{Time: nextRetryAt},
		TenantID:         tenantID,
	}
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("seed external sync retry: %v", err)
	}
}

func seedScheduleConflictOverride(t *testing.T, db *gorm.DB, tenantID, id uint, conflictType string) {
	t.Helper()
	row := &models.EduScheduleConflictOverride{
		BaseModel:    models.BaseModel{ID: id},
		ConflictType: conflictType,
		ConflictKey:  conflictType + "-key",
		Reason:       "override",
		TenantID:     tenantID,
	}
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("seed conflict override: %v", err)
	}
}

func seedLessonEligibility(t *testing.T, db *gorm.DB, tenantID, lessonID, studentID, courseID, benefitID uint, status, reasonCode string) {
	t.Helper()
	row := &models.EduLessonStudentEligibility{
		LessonID:           lessonID,
		StudentID:          studentID,
		CourseID:           courseID,
		StudentBenefitID:   benefitID,
		EligibilityStatus:  status,
		ReasonCode:         reasonCode,
		CheckedAt:          datePtr(time.Now().UTC()),
		TenantID:           tenantID,
	}
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("seed lesson eligibility: %v", err)
	}
}
