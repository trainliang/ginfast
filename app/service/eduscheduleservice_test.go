package service

import (
	"fmt"
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
