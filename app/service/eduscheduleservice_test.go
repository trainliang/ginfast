package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"gin-fast/app/models"

	"gorm.io/gorm"
)

func TestEduScheduleServiceCompile(t *testing.T) {
	setupEduTestDB(t)
	_ = NewEduScheduleService()
}

func TestEduScheduleServiceGenerateOneToOneSingleLesson(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedStudent(t, db, 1, 1)
	seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC), time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC))

	rule := &models.EduScheduleRule{
		RuleType:     "one_to_one",
		RepeatType:   "single",
		StudentID:    1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		StartDate:    &models.JSONTime{Time: time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)},
		StartTime:    "09:00",
		EndTime:      "10:00",
		Status:       1,
		TenantID:     1,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("seed rule: %v", err)
	}

	lessons, err := svc.GenerateLessonsForRule(ctx, rule.ID, 1)
	if err != nil {
		t.Fatalf("generate lessons: %v", err)
	}
	if len(lessons) != 1 {
		t.Fatalf("expected 1 lesson, got %d", len(lessons))
	}
	if lessons[0].StudentID != 1 || lessons[0].LessonType != "one_to_one" || lessons[0].LessonDate == nil {
		t.Fatalf("unexpected lesson: %+v", lessons[0])
	}
}

func TestEduScheduleServiceGenerateWeeklyClassLessonsByTerm(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedStudents(t, db, 1, 2)
	seedClassWithMembers(t, db, 1, 1, []uint{1, 2}, "required")
	seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))
	seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 2, 2, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))
	if err := db.Create(&models.EduTerm{Name: "2026春季", StartDate: &models.JSONTime{Time: time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)}, EndDate: &models.JSONTime{Time: time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)}, Status: 1, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed term: %v", err)
	}

	rule := &models.EduScheduleRule{
		RuleType:     "class",
		RepeatType:   "weekly",
		TermID:       1,
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		Weekday:      1,
		StartTime:    "09:00",
		EndTime:      "10:00",
		Status:       1,
		TenantID:     1,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("seed rule: %v", err)
	}

	lessons, err := svc.GenerateLessonsForRule(ctx, rule.ID, 1)
	if err != nil {
		t.Fatalf("generate lessons: %v", err)
	}
	if len(lessons) != 3 {
		t.Fatalf("expected 3 lessons within term, got %d", len(lessons))
	}
}

func TestEduScheduleServiceSkipsClosedDays(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedStudents(t, db, 1, 1)
	seedClassWithMembers(t, db, 1, 1, []uint{1}, "required")
	seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))
	if err := db.Create(&models.EduTerm{Name: "2026春季", StartDate: &models.JSONTime{Time: time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)}, EndDate: &models.JSONTime{Time: time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)}, Status: 1, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed term: %v", err)
	}
	if err := db.Create(&models.EduTermClosedDay{TermID: 1, ClosedDate: &models.JSONTime{Time: time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)}, Reason: "假期", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed closed day: %v", err)
	}

	rule := &models.EduScheduleRule{
		RuleType:     "class",
		RepeatType:   "weekly",
		TermID:       1,
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		Weekday:      1,
		StartTime:    "09:00",
		EndTime:      "10:00",
		Status:       1,
		TenantID:     1,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("seed rule: %v", err)
	}

	lessons, err := svc.GenerateLessonsForRule(ctx, rule.ID, 1)
	if err != nil {
		t.Fatalf("generate lessons: %v", err)
	}
	if len(lessons) != 2 {
		t.Fatalf("expected 2 lessons after skipping closed day, got %d", len(lessons))
	}
	for _, lesson := range lessons {
		if lesson.LessonDate != nil && lesson.LessonDate.Time.Equal(time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("closed day lesson should not be generated")
		}
	}
}

func TestEduScheduleServiceRequiresRoomForOfflineAndAllowsOnline(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedStudent(t, db, 1, 1)
	seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC), time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC))

	offlineRule := &models.EduScheduleRule{
		RuleType:     "one_to_one",
		RepeatType:   "single",
		StudentID:    1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		StartDate:    &models.JSONTime{Time: time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)},
		StartTime:    "09:00",
		EndTime:      "10:00",
		Status:       1,
		TenantID:     1,
	}
	if err := db.Create(offlineRule).Error; err != nil {
		t.Fatalf("seed offline rule: %v", err)
	}
	if _, err := svc.GenerateLessonsForRule(ctx, offlineRule.ID, 1); err == nil {
		t.Fatalf("expected offline lesson without room to fail")
	}

	onlineRule := &models.EduScheduleRule{
		RuleType:     "one_to_one",
		RepeatType:   "single",
		StudentID:    1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "online",
		RequiresRoom: 0,
		StartDate:    &models.JSONTime{Time: time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)},
		StartTime:    "09:00",
		EndTime:      "10:00",
		Status:       1,
		TenantID:     1,
	}
	if err := db.Create(onlineRule).Error; err != nil {
		t.Fatalf("seed online rule: %v", err)
	}
	lessons, err := svc.GenerateLessonsForRule(ctx, onlineRule.ID, 1)
	if err != nil {
		t.Fatalf("generate online lessons: %v", err)
	}
	if len(lessons) != 1 {
		t.Fatalf("expected online lesson to generate, got %d", len(lessons))
	}
}

func TestEduScheduleServiceWritesEligibilityForActiveClassMembers(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedStudents(t, db, 1, 3)
	seedClassWithMembersWithDates(t, db, 1, 1, []classMemberSeed{
		{StudentID: 1, JoinDate: datePtr(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)), Status: "studying"},
		{StudentID: 2, JoinDate: datePtr(time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)), Status: "studying"},
		{StudentID: 3, JoinDate: datePtr(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)), LeaveDate: datePtr(time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)), Status: "left"},
	}, "required")
	seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))
	seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 2, 2, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))
	if err := db.Create(&models.EduTerm{Name: "2026春季", StartDate: &models.JSONTime{Time: time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)}, EndDate: &models.JSONTime{Time: time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)}, Status: 1, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed term: %v", err)
	}

	rule := &models.EduScheduleRule{
		RuleType:     "class",
		RepeatType:   "weekly",
		TermID:       1,
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		Weekday:      1,
		StartTime:    "09:00",
		EndTime:      "10:00",
		Status:       1,
		TenantID:     1,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("seed rule: %v", err)
	}

	lessons, err := svc.GenerateLessonsForRule(ctx, rule.ID, 1)
	if err != nil {
		t.Fatalf("generate lessons: %v", err)
	}
	if len(lessons) != 3 {
		t.Fatalf("expected 3 lessons, got %d", len(lessons))
	}

	var eligibilities []models.EduLessonStudentEligibility
	if err := db.Order("student_id asc").Find(&eligibilities).Error; err != nil {
		t.Fatalf("load eligibilities: %v", err)
	}
	if len(eligibilities) != 4 {
		t.Fatalf("expected 4 eligibility rows, got %d", len(eligibilities))
	}
	if eligibilities[0].StudentID != 1 || eligibilities[0].EligibilityStatus != "eligible" {
		t.Fatalf("unexpected eligibility rows: %+v", eligibilities)
	}
	if eligibilities[3].StudentID != 2 || eligibilities[3].EligibilityStatus != "eligible" {
		t.Fatalf("unexpected final eligibility row: %+v", eligibilities[3])
	}
}

func TestEduScheduleServiceGenerateRequiredClassPolicyRejectsIneligibleOnlineLesson(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedStudents(t, db, 1, 1)
	seedClassWithMembers(t, db, 1, 1, []uint{1}, "required")
	rule := &models.EduScheduleRule{
		RuleType:     "class",
		RepeatType:   "single",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "online",
		RequiresRoom: 0,
		StartDate:    &models.JSONTime{Time: time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)},
		StartTime:    "09:00",
		EndTime:      "10:00",
		Status:       1,
		TenantID:     1,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("seed rule: %v", err)
	}

	if _, err := svc.GenerateLessonsForRule(ctx, rule.ID, 1); err == nil {
		t.Fatalf("expected required class policy to reject ineligible online lesson")
	}
	var lessonCount int64
	if err := db.Model(&models.EduLesson{}).Count(&lessonCount).Error; err != nil {
		t.Fatalf("count lessons: %v", err)
	}
	if lessonCount != 0 {
		t.Fatalf("expected transaction rollback, got %d lessons", lessonCount)
	}
}

func TestEduScheduleServiceGenerateNoneClassPolicySkipsEligibilityRows(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedStudents(t, db, 1, 1)
	seedClassWithMembers(t, db, 1, 1, []uint{1}, "none")
	rule := &models.EduScheduleRule{
		RuleType:     "class",
		RepeatType:   "single",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    1,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       1,
		StartDate:    &models.JSONTime{Time: time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)},
		StartTime:    "09:00",
		EndTime:      "10:00",
		Status:       1,
		TenantID:     1,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("seed rule: %v", err)
	}

	lessons, err := svc.GenerateLessonsForRule(ctx, rule.ID, 1)
	if err != nil {
		t.Fatalf("generate lessons: %v", err)
	}
	if len(lessons) != 1 {
		t.Fatalf("expected lesson to generate, got %d", len(lessons))
	}
	var eligibilityCount int64
	if err := db.Model(&models.EduLessonStudentEligibility{}).Count(&eligibilityCount).Error; err != nil {
		t.Fatalf("count eligibility rows: %v", err)
	}
	if eligibilityCount != 0 {
		t.Fatalf("expected no eligibility rows for none policy, got %d", eligibilityCount)
	}
}

func TestEduScheduleServiceConflictRejectsTeacherOverlap(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "one_to_one",
		LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:00",
		EndTime:      "10:00",
		StudentID:    1,
		CourseID:     1,
		TeacherID:    7,
		TeachingMode: "online",
		RequiresRoom: 0,
		Status:       "scheduled",
	})

	result, err := svc.CheckConflicts(ctx, &models.EduScheduleConflictCheckRequest{
		LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:30",
		EndTime:      "10:30",
		StudentID:    2,
		CourseID:     1,
		TeacherID:    7,
		TeachingMode: "online",
		RequiresRoom: 0,
	})
	if err == nil {
		t.Fatalf("expected teacher conflict error")
	}
	if result == nil || !result.HasConflict || len(result.Items) != 1 || result.Items[0].ConflictType != "teacher" {
		t.Fatalf("expected teacher conflict result, got %+v", result)
	}
}

func TestEduScheduleServiceConflictRejectsRoomOverlap(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "class",
		LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:00",
		EndTime:      "10:00",
		ClassID:      1,
		CourseID:     1,
		TeacherID:    7,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       3,
		Status:       "scheduled",
	})

	result, err := svc.CheckConflicts(ctx, &models.EduScheduleConflictCheckRequest{
		LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:30",
		EndTime:      "10:30",
		ClassID:      2,
		CourseID:     1,
		TeacherID:    8,
		TeachingMode: "offline",
		RequiresRoom: 1,
		RoomID:       3,
	})
	if err == nil {
		t.Fatalf("expected room conflict error")
	}
	if result == nil || !result.HasConflict || len(result.Items) != 1 || result.Items[0].ConflictType != "room" {
		t.Fatalf("expected room conflict result, got %+v", result)
	}
}

func TestEduScheduleServiceConflictRejectsStudentOverlap(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "one_to_one",
		LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:00",
		EndTime:      "10:00",
		StudentID:    5,
		CourseID:     1,
		TeacherID:    7,
		TeachingMode: "online",
		RequiresRoom: 0,
		Status:       "scheduled",
	})

	result, err := svc.CheckConflicts(ctx, &models.EduScheduleConflictCheckRequest{
		LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:30",
		EndTime:      "10:30",
		StudentID:    5,
		CourseID:     1,
		TeacherID:    8,
		TeachingMode: "online",
		RequiresRoom: 0,
	})
	if err == nil {
		t.Fatalf("expected student conflict error")
	}
	if result == nil || !result.HasConflict || len(result.Items) != 1 || result.Items[0].ConflictType != "student" {
		t.Fatalf("expected student conflict result, got %+v", result)
	}
}

func TestEduScheduleServiceConflictRejectsClassOverlap(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "class",
		LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:00",
		EndTime:      "10:00",
		ClassID:      9,
		CourseID:     1,
		TeacherID:    7,
		TeachingMode: "online",
		RequiresRoom: 0,
		Status:       "scheduled",
	})

	result, err := svc.CheckConflicts(ctx, &models.EduScheduleConflictCheckRequest{
		LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:30",
		EndTime:      "10:30",
		ClassID:      9,
		CourseID:     1,
		TeacherID:    8,
		TeachingMode: "online",
		RequiresRoom: 0,
	})
	if err == nil {
		t.Fatalf("expected class conflict error")
	}
	if result == nil || !result.HasConflict || len(result.Items) != 1 || result.Items[0].ConflictType != "class" {
		t.Fatalf("expected class conflict result, got %+v", result)
	}
}

func TestEduScheduleServiceConflictOverrideRequiresReasonAndRecordsOverride(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedExistingLesson(t, db, 1, &models.EduLesson{
		LessonType:   "one_to_one",
		LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:    "09:00",
		EndTime:      "10:00",
		StudentID:    1,
		CourseID:     1,
		TeacherID:    7,
		TeachingMode: "online",
		RequiresRoom: 0,
		Status:       "scheduled",
	})
	req := &models.EduScheduleConflictCheckRequest{
		LessonDate:            datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:             "09:30",
		EndTime:               "10:30",
		StudentID:             2,
		CourseID:              1,
		TeacherID:             7,
		TeachingMode:          "online",
		RequiresRoom:          0,
		AllowConflictOverride: true,
		OperatorID:            99,
	}
	if _, err := svc.CheckConflicts(ctx, req); err == nil {
		t.Fatalf("expected override without reason to fail")
	}

	req.OverrideReason = "管理员确认可覆盖"
	result, err := svc.CheckConflicts(ctx, req)
	if err != nil {
		t.Fatalf("override conflict: %v", err)
	}
	if result == nil || !result.HasConflict || len(result.Items) != 1 {
		t.Fatalf("expected conflict result with override, got %+v", result)
	}
	var overrides []models.EduScheduleConflictOverride
	if err := db.Find(&overrides).Error; err != nil {
		t.Fatalf("load overrides: %v", err)
	}
	if len(overrides) != 1 || overrides[0].ConflictType != "teacher" || overrides[0].Reason != "管理员确认可覆盖" || overrides[0].OperatorID != 99 {
		t.Fatalf("unexpected override row: %+v", overrides)
	}
}

func TestEduScheduleServiceRescheduleLesson(t *testing.T) {
	t.Run("updates date time teacher room and logs change", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		seedStudent(t, db, 1, 1)
		seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))

		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    7,
			TeachingMode: "offline",
			RequiresRoom: 1,
			RoomID:       3,
			Status:       "scheduled",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)

		err := svc.RescheduleLesson(ctx, 1, &models.EduLessonRescheduleRequest{
			LessonID:   lesson.ID,
			LessonDate: datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
			StartTime:  "10:00",
			EndTime:    "11:00",
			TeacherID:  8,
			TeachingMode: "offline",
			RoomID:    uintPtr(5),
		})
		if err != nil {
			t.Fatalf("reschedule lesson: %v", err)
		}

		var updated models.EduLesson
		if err := db.First(&updated, lesson.ID).Error; err != nil {
			t.Fatalf("load lesson: %v", err)
		}
		if updated.LessonDate == nil || !updated.LessonDate.Time.Equal(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("unexpected lesson date: %+v", updated)
		}
		if updated.StartTime != "10:00" || updated.EndTime != "11:00" || updated.TeacherID != 8 || updated.RoomID != 5 {
			t.Fatalf("unexpected updated lesson: %+v", updated)
		}
		if updated.IsManualAdjusted != 1 {
			t.Fatalf("expected manual adjusted lesson, got %+v", updated)
		}

		var logs []models.EduLessonChangeLog
		if err := db.Order("id asc").Find(&logs).Error; err != nil {
			t.Fatalf("load change logs: %v", err)
		}
		if len(logs) != 1 {
			t.Fatalf("expected 1 change log, got %d", len(logs))
		}
		if logs[0].ActionType != "reschedule" || logs[0].LessonID != lesson.ID {
			t.Fatalf("unexpected change log: %+v", logs[0])
		}
		if logs[0].BeforeData == "" || logs[0].AfterData == "" {
			t.Fatalf("expected before/after data in log: %+v", logs[0])
		}
		if !strings.Contains(logs[0].BeforeData, "2026-05-11") || !strings.Contains(logs[0].AfterData, "2026-05-12") {
			t.Fatalf("unexpected change log payload: %+v", logs[0])
		}
	})

	t.Run("conflict rejects normal user", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		seedStudent(t, db, 1, 1)
		seedStudent(t, db, 1, 2)
		seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))
		seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 2, 2, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))

		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    7,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "scheduled",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)
		seedExistingLesson(t, db, 1, &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
			StartTime:    "10:30",
			EndTime:      "11:30",
			StudentID:    2,
			CourseID:     1,
			TeacherID:    8,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "scheduled",
			TenantID:     1,
		})

		err := svc.RescheduleLesson(ctx, 1, &models.EduLessonRescheduleRequest{
			LessonID:    lesson.ID,
			LessonDate:  datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
			StartTime:   "10:00",
			EndTime:     "11:00",
			TeacherID:   8,
			TeachingMode: "online",
		})
		if err == nil {
			t.Fatalf("expected conflict to be rejected")
		}
		var overrides []models.EduScheduleConflictOverride
		if err := db.Find(&overrides).Error; err != nil {
			t.Fatalf("load overrides: %v", err)
		}
		if len(overrides) != 0 {
			t.Fatalf("expected no conflict override rows, got %+v", overrides)
		}
	})

	t.Run("admin override writes conflict override", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		seedStudent(t, db, 1, 1)
		seedStudent(t, db, 1, 2)
		seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))
		seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 2, 2, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))

		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    7,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "scheduled",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)
		seedExistingLesson(t, db, 1, &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
			StartTime:    "10:30",
			EndTime:      "11:30",
			StudentID:    2,
			CourseID:     1,
			TeacherID:    8,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "scheduled",
			TenantID:     1,
		})

		err := svc.RescheduleLesson(ctx, 1, &models.EduLessonRescheduleRequest{
			LessonID:              lesson.ID,
			LessonDate:            datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
			StartTime:             "10:00",
			EndTime:               "11:00",
			TeacherID:             8,
			TeachingMode:          "online",
			AllowConflictOverride: true,
			OverrideReason:        "管理员确认可覆盖",
		})
		if err != nil {
			t.Fatalf("reschedule with override: %v", err)
		}
		var overrides []models.EduScheduleConflictOverride
		if err := db.Order("id asc").Find(&overrides).Error; err != nil {
			t.Fatalf("load overrides: %v", err)
		}
		if len(overrides) != 1 {
			t.Fatalf("expected 1 override row, got %+v", overrides)
		}
		if overrides[0].Reason != "管理员确认可覆盖" || overrides[0].ConflictType == "" || overrides[0].ConflictKey == "" {
			t.Fatalf("unexpected override row: %+v", overrides[0])
		}
	})

	t.Run("completed lesson is blocked", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		seedStudent(t, db, 1, 1)
		seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))

		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    7,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "completed",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)

		err := svc.RescheduleLesson(ctx, 1, &models.EduLessonRescheduleRequest{
			LessonID:    lesson.ID,
			LessonDate:  datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
			StartTime:   "10:00",
			EndTime:     "11:00",
			TeacherID:   8,
			TeachingMode: "online",
		})
		if err == nil {
			t.Fatalf("expected completed lesson to be blocked")
		}
	})

	t.Run("offline lesson requires room and invalid time is rejected", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		seedStudent(t, db, 1, 1)
		seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))

		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    7,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "scheduled",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)

		if err := svc.RescheduleLesson(ctx, 1, &models.EduLessonRescheduleRequest{
			LessonID:    lesson.ID,
			LessonDate:  datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
			StartTime:   "11:00",
			EndTime:     "10:00",
			TeacherID:   8,
			TeachingMode: "online",
		}); err == nil {
			t.Fatalf("expected invalid time to fail")
		}

		if err := svc.RescheduleLesson(ctx, 1, &models.EduLessonRescheduleRequest{
			LessonID:    lesson.ID,
			LessonDate:  datePtr(time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)),
			StartTime:   "10:00",
			EndTime:     "11:00",
			TeacherID:   8,
			TeachingMode: "offline",
		}); err == nil {
			t.Fatalf("expected offline lesson without room to fail")
		}
	})
}

func TestEduScheduleServiceStopCancelRestoreLesson(t *testing.T) {
	t.Run("stop lesson updates status and writes change log", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    1,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "scheduled",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)

		if err := svc.StopLesson(ctx, 1, &models.EduLessonStopRequest{LessonID: lesson.ID, Reason: "临时停课"}); err != nil {
			t.Fatalf("stop lesson: %v", err)
		}

		var updated models.EduLesson
		if err := db.First(&updated, lesson.ID).Error; err != nil {
			t.Fatalf("load lesson: %v", err)
		}
		if updated.Status != "stopped" {
			t.Fatalf("expected stopped status, got %+v", updated)
		}
		var logs []models.EduLessonChangeLog
		if err := db.Order("id asc").Find(&logs).Error; err != nil {
			t.Fatalf("load change logs: %v", err)
		}
		if len(logs) != 1 {
			t.Fatalf("expected 1 change log, got %d", len(logs))
		}
		if logs[0].ActionType != "stop" || logs[0].LessonID != lesson.ID || logs[0].Reason != "临时停课" {
			t.Fatalf("unexpected stop log: %+v", logs[0])
		}
	})

	t.Run("cancel lesson updates status and writes change log", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    1,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "scheduled",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)

		if err := svc.CancelLesson(ctx, 1, &models.EduLessonCancelRequest{LessonID: lesson.ID, Reason: "客户取消"}); err != nil {
			t.Fatalf("cancel lesson: %v", err)
		}

		var updated models.EduLesson
		if err := db.First(&updated, lesson.ID).Error; err != nil {
			t.Fatalf("load lesson: %v", err)
		}
		if updated.Status != "canceled" {
			t.Fatalf("expected canceled status, got %+v", updated)
		}
		var logs []models.EduLessonChangeLog
		if err := db.Order("id asc").Find(&logs).Error; err != nil {
			t.Fatalf("load change logs: %v", err)
		}
		if len(logs) != 1 {
			t.Fatalf("expected 1 change log, got %d", len(logs))
		}
		if logs[0].ActionType != "cancel" || logs[0].LessonID != lesson.ID || logs[0].Reason != "客户取消" {
			t.Fatalf("unexpected cancel log: %+v", logs[0])
		}
	})

	t.Run("restore lesson returns to scheduled and rechecks conflict and eligibility", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		seedStudent(t, db, 1, 1)
		seedStudent(t, db, 1, 2)
		seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 1, 1, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))
		seedBenefitProductAndStudentBenefitForSchedule(t, db, 1, 2, 2, 1, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC))

		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    1,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "stopped",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)
		seedExistingLesson(t, db, 1, &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:30",
			EndTime:      "10:30",
			StudentID:    2,
			CourseID:     1,
			TeacherID:    1,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "scheduled",
			TenantID:     1,
		})

		if err := svc.RestoreLesson(ctx, 1, &models.EduLessonRestoreRequest{LessonID: lesson.ID, Reason: "恢复排课"}); err == nil {
			t.Fatalf("expected conflict to block normal restore")
		}

		var updated models.EduLesson
		if err := db.First(&updated, lesson.ID).Error; err != nil {
			t.Fatalf("load lesson: %v", err)
		}
		if updated.Status != "stopped" {
			t.Fatalf("restore should rollback on conflict, got %+v", updated)
		}
	})

	t.Run("restore lesson rejects when eligibility is invalid", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		seedStudent(t, db, 1, 1)

		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    1,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "canceled",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)

		if err := svc.RestoreLesson(ctx, 1, &models.EduLessonRestoreRequest{LessonID: lesson.ID, Reason: "恢复排课"}); err == nil {
			t.Fatalf("expected restore to fail when eligibility is invalid")
		}
		var updated models.EduLesson
		if err := db.First(&updated, lesson.ID).Error; err != nil {
			t.Fatalf("load lesson: %v", err)
		}
		if updated.Status != "canceled" {
			t.Fatalf("restore should rollback on eligibility failure, got %+v", updated)
		}
	})

	t.Run("completed lesson is blocked for stop cancel restore", func(t *testing.T) {
		db := setupEduTestDB(t)
		svc := NewEduScheduleService()
		ctx := contextWithTenant(1)
		seedEduTenantData(t, 1)
		lesson := &models.EduLesson{
			LessonType:   "one_to_one",
			LessonDate:   datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
			StartTime:    "09:00",
			EndTime:      "10:00",
			StudentID:    1,
			CourseID:     1,
			TeacherID:    1,
			TeachingMode: "online",
			RequiresRoom: 0,
			Status:       "completed",
			TenantID:     1,
		}
		seedExistingLesson(t, db, 1, lesson)

		if err := svc.StopLesson(ctx, 1, &models.EduLessonStopRequest{LessonID: lesson.ID, Reason: "临时停课"}); err == nil {
			t.Fatalf("expected completed stop to fail")
		}
		if err := svc.CancelLesson(ctx, 1, &models.EduLessonCancelRequest{LessonID: lesson.ID, Reason: "客户取消"}); err == nil {
			t.Fatalf("expected completed cancel to fail")
		}
		if err := svc.RestoreLesson(ctx, 1, &models.EduLessonRestoreRequest{LessonID: lesson.ID, Reason: "恢复"}); err == nil {
			t.Fatalf("expected completed restore to fail")
		}
	})
}

func TestEduScheduleServiceRuleChangePreviewAndRegenerate(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduScheduleService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedStudent(t, db, 1, 1)
	seedStudent(t, db, 1, 2)
	seedStudent(t, db, 1, 3)
	seedStudent(t, db, 1, 4)

	rule := &models.EduScheduleRule{
		RuleType:      "one_to_one",
		RepeatType:    "weekly",
		StudentID:     1,
		CourseID:      1,
		TeacherID:     1,
		TeachingMode:  "offline",
		RequiresRoom:  1,
		RoomID:        1,
		Weekday:       1,
		StartTime:     "09:00",
		EndTime:       "10:00",
		StartDate:     datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		EndDate:       datePtr(time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)),
		Status:        1,
		Version:       1,
		EffectiveFrom: datePtr(time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)),
		TenantID:      1,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("seed rule: %v", err)
	}

	seedExistingLesson(t, db, 1, &models.EduLesson{
		RuleID:      rule.ID,
		RuleVersion: 1,
		LessonType:  "one_to_one",
		LessonDate:  datePtr(time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)),
		StartTime:   "09:00",
		EndTime:     "10:00",
		StudentID:   1,
		CourseID:    1,
		TeacherID:   1,
		Status:      "scheduled",
	})
	seedExistingLesson(t, db, 1, &models.EduLesson{
		RuleID:      rule.ID,
		RuleVersion: 1,
		LessonType:  "one_to_one",
		LessonDate:  datePtr(time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)),
		StartTime:   "09:00",
		EndTime:     "10:00",
		StudentID:   2,
		CourseID:    1,
		TeacherID:   1,
		Status:      "scheduled",
	})
	seedExistingLesson(t, db, 1, &models.EduLesson{
		RuleID:      rule.ID,
		RuleVersion: 1,
		LessonType:  "one_to_one",
		LessonDate:  datePtr(time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)),
		StartTime:   "09:00",
		EndTime:     "10:00",
		StudentID:   3,
		CourseID:    1,
		TeacherID:   1,
		Status:      "completed",
	})
	seedExistingLesson(t, db, 1, &models.EduLesson{
		RuleID:           rule.ID,
		RuleVersion:      1,
		LessonType:       "one_to_one",
		LessonDate:       datePtr(time.Date(2026, 5, 26, 0, 0, 0, 0, time.UTC)),
		StartTime:        "09:00",
		EndTime:          "10:00",
		StudentID:        4,
		CourseID:         1,
		TeacherID:        1,
		Status:           "scheduled",
		IsManualAdjusted: 1,
	})
	seedExistingLesson(t, db, 1, &models.EduLesson{
		RuleID:      rule.ID,
		RuleVersion: 1,
		LessonType:  "one_to_one",
		LessonDate:  datePtr(time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)),
		StartTime:   "09:00",
		EndTime:     "10:00",
		StudentID:   1,
		CourseID:    1,
		TeacherID:   1,
		Status:      "scheduled",
	})

	preview, err := svc.PreviewRuleChange(ctx, 1, rule.ID)
	if err != nil {
		t.Fatalf("preview rule change: %v", err)
	}
	if len(preview) != 1 {
		t.Fatalf("expected 1 replaceable future lesson, got %d", len(preview))
	}
	if preview[0].LessonDate == nil || !preview[0].LessonDate.Time.Equal(time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected preview first lesson: %+v", preview[0])
	}

	rule.Version = 2
	if err := db.Model(&models.EduScheduleRule{}).Where("id = ?", rule.ID).Update("version", 2).Error; err != nil {
		t.Fatalf("bump rule version: %v", err)
	}
	regenerated, err := svc.RegenerateFutureLessons(ctx, 1, rule.ID, 99)
	if err != nil {
		t.Fatalf("regenerate future lessons: %v", err)
	}
	if len(regenerated) != 1 {
		t.Fatalf("expected 1 regenerated lesson, got %d", len(regenerated))
	}
	for _, lesson := range regenerated {
		if lesson.RuleVersion != 2 {
			t.Fatalf("expected regenerated lesson version 2, got %+v", lesson)
		}
		if lesson.LessonDate != nil && lesson.LessonDate.Time.Before(time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("regenerated lesson should not include past lessons: %+v", lesson)
		}
	}

	var lessons []models.EduLesson
	if err := db.Order("lesson_date asc, id asc").Find(&lessons).Error; err != nil {
		t.Fatalf("load lessons: %v", err)
	}
	if len(lessons) != 5 {
		t.Fatalf("expected 5 lessons after regeneration, got %d", len(lessons))
	}
	for _, lesson := range lessons {
		if lesson.Status == "completed" || lesson.IsManualAdjusted == 1 {
			continue
		}
		if lesson.LessonDate == nil {
			t.Fatalf("lesson date missing: %+v", lesson)
		}
		if !lesson.LessonDate.Time.Before(time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)) && lesson.RuleVersion != 2 {
			t.Fatalf("expected future replaceable lessons to use new version: %+v", lesson)
		}
	}

	var logs []models.EduLessonChangeLog
	if err := db.Order("id asc").Find(&logs).Error; err != nil {
		t.Fatalf("load change logs: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 change log, got %d", len(logs))
	}
	for _, log := range logs {
		if log.ActionType != "rule_regenerate" {
			t.Fatalf("unexpected action type: %+v", log)
		}
		if log.OperatorID != 99 {
			t.Fatalf("unexpected operator id: %+v", log)
		}
	}
}

func seedStudent(t *testing.T, db *gorm.DB, tenantID, id uint) {
	t.Helper()
	if err := db.Create(&models.EduStudent{BaseModel: models.BaseModel{ID: id}, Name: fmt.Sprintf("学生-%d", id), Phone: fmt.Sprintf("1380000%04d", id), TenantID: tenantID}).Error; err != nil {
		t.Fatalf("seed student: %v", err)
	}
}

func seedStudents(t *testing.T, db *gorm.DB, tenantID uint, count int) {
	t.Helper()
	for i := 1; i <= count; i++ {
		seedStudent(t, db, tenantID, uint(i))
	}
}

type classMemberSeed struct {
	StudentID uint
	JoinDate  *models.JSONTime
	LeaveDate *models.JSONTime
	Status    string
}

func seedClassWithMembers(t *testing.T, db *gorm.DB, tenantID, classID uint, studentIDs []uint, policy string) {
	t.Helper()
	members := make([]classMemberSeed, 0, len(studentIDs))
	for _, studentID := range studentIDs {
		members = append(members, classMemberSeed{StudentID: studentID, JoinDate: datePtr(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)), Status: "studying"})
	}
	seedClassWithMembersWithDates(t, db, tenantID, classID, members, policy)
}

func seedClassWithMembersWithDates(t *testing.T, db *gorm.DB, tenantID, classID uint, members []classMemberSeed, policy string) {
	t.Helper()
	if err := db.Create(&models.EduClass{Name: "班级", Code: fmt.Sprintf("CL-%d", classID), ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 10, BenefitCheckPolicy: policy, TenantID: tenantID, RoomID: 1}).Error; err != nil {
		t.Fatalf("seed class: %v", err)
	}
	for _, member := range members {
		if err := db.Create(&models.EduClassMember{ClassID: classID, StudentID: member.StudentID, JoinDate: member.JoinDate, LeaveDate: member.LeaveDate, Status: member.Status, TenantID: tenantID}).Error; err != nil {
			t.Fatalf("seed class member: %v", err)
		}
	}
}

func seedBenefitProductAndStudentBenefitForSchedule(t *testing.T, db *gorm.DB, tenantID, studentID, productID, courseID uint, validFrom, validTo time.Time) {
	t.Helper()
	if err := db.Create(&models.EduBenefitProduct{BaseModel: models.BaseModel{ID: productID}, Name: fmt.Sprintf("权益-%d", productID), Code: fmt.Sprintf("B%03d", productID), BenefitType: "course", CalculationMode: "period_unlimited", Status: 1, TenantID: tenantID}).Error; err != nil {
		t.Fatalf("seed product: %v", err)
	}
	if err := db.Create(&models.EduStudentBenefit{StudentID: studentID, ProductID: productID, BenefitType: "course", CalculationMode: "period_unlimited", CourseID: courseID, TotalCount: 10, RemainingCount: 10, Status: 1, ValidFrom: &models.JSONTime{Time: validFrom}, ValidTo: &models.JSONTime{Time: validTo}, TenantID: tenantID}).Error; err != nil {
		t.Fatalf("seed student benefit: %v", err)
	}
}

func datePtr(t time.Time) *models.JSONTime {
	return &models.JSONTime{Time: t}
}

func seedExistingLesson(t *testing.T, db *gorm.DB, tenantID uint, lesson *models.EduLesson) {
	t.Helper()
	lesson.TenantID = tenantID
	if lesson.Status == "" {
		lesson.Status = "scheduled"
	}
	if err := db.Create(lesson).Error; err != nil {
		t.Fatalf("seed existing lesson: %v", err)
	}
}
