package service

import (
	"fmt"
	"testing"
	"time"

	"gin-fast/app/global/app"
	"gin-fast/app/models"
)

func validCompletionRuleItems() []models.EduLessonCompletionRuleItemRequest {
	items := make([]models.EduLessonCompletionRuleItemRequest, 0, len(models.EduLessonCompletionResultTypes))
	for _, resultType := range models.EduLessonCompletionResultTypes {
		items = append(items, models.EduLessonCompletionRuleItemRequest{
			ResultType:          resultType,
			DeductEnabled:       0,
			DeductMode:          models.EduLessonCompletionDeductModeFixedCount,
			FixedCount:          1,
			DurationUnitMinutes: 60,
			InsufficientPolicy:  models.EduLessonCompletionInsufficientBlock,
			Status:              1,
		})
	}
	return items
}

func TestEduLessonCompletionServiceRules(t *testing.T) {
	setupEduTestDB(t)
	ctx := contextWithTenantAndUser(1, 99)
	svc := NewEduLessonCompletionService()

	t.Run("requires complete rule set", func(t *testing.T) {
		items := validCompletionRuleItems()
		items = items[:len(items)-1]
		err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items})
		if err == nil {
			t.Fatalf("expected incomplete rule set to fail")
		}
	})

	t.Run("saves and lists six tenant rules", func(t *testing.T) {
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: validCompletionRuleItems()}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		rules, err := svc.ListRules(ctx, 1)
		if err != nil {
			t.Fatalf("list rules: %v", err)
		}
		if len(rules) != len(models.EduLessonCompletionResultTypes) {
			t.Fatalf("expected %d rules, got %d", len(models.EduLessonCompletionResultTypes), len(rules))
		}
	})

	t.Run("upserts existing tenant rules", func(t *testing.T) {
		items := validCompletionRuleItems()
		items[0].DeductEnabled = 1
		items[0].FixedCount = 2
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save updated rules: %v", err)
		}
		rules, err := svc.ListRules(ctx, 1)
		if err != nil {
			t.Fatalf("list updated rules: %v", err)
		}
		if len(rules) != len(models.EduLessonCompletionResultTypes) {
			t.Fatalf("expected upsert to keep six rules, got %d", len(rules))
		}
		var attended *models.EduLessonCompletionRule
		for _, rule := range rules {
			if rule.ResultType == models.EduLessonCompletionResultAttended {
				attended = &rule
				break
			}
		}
		if attended == nil || attended.DeductEnabled != 1 || attended.FixedCount != 2 {
			t.Fatalf("expected attended rule to update, got %+v", attended)
		}
	})

	t.Run("returns default rules when tenant has no configured rules", func(t *testing.T) {
		rules, err := svc.ListRules(contextWithTenantAndUser(3, 99), 3)
		if err != nil {
			t.Fatalf("list default rules: %v", err)
		}
		if len(rules) != len(models.EduLessonCompletionResultTypes) {
			t.Fatalf("expected %d rules, got %d", len(models.EduLessonCompletionResultTypes), len(rules))
		}
		var attended *models.EduLessonCompletionRule
		for i := range rules {
			if rules[i].ResultType == models.EduLessonCompletionResultAttended {
				attended = &rules[i]
				break
			}
		}
		if attended == nil {
			t.Fatalf("expected attended default rule")
		}
		if attended.DeductEnabled != 1 || attended.DeductMode != models.EduLessonCompletionDeductModeFixedCount || attended.FixedCount != 1 {
			t.Fatalf("unexpected attended default rule: %+v", attended)
		}
	})
}

func TestEduLessonCompletionServiceStatus(t *testing.T) {
	setupEduTestDB(t)
	ctx := contextWithTenantAndUser(1, 99)
	db := app.DB()
	svc := NewEduLessonCompletionService()

	course := &models.EduCourse{Name: "数学", Code: "MATH", TenantID: 1}
	if err := db.Create(course).Error; err != nil {
		t.Fatalf("create course: %v", err)
	}
	class := &models.EduClass{Name: "一班", Code: "C1", ClassType: "group", CourseID: course.ID, TeacherID: 10, Capacity: 10, Status: 1, TenantID: 1}
	if err := db.Create(class).Error; err != nil {
		t.Fatalf("create class: %v", err)
	}
	students := []models.EduStudent{
		{Name: "学生A", Phone: "13000000001", Status: 1, TenantID: 1},
		{Name: "学生B", Phone: "13000000002", Status: 1, TenantID: 1},
	}
	if err := db.Create(&students).Error; err != nil {
		t.Fatalf("create students: %v", err)
	}
	for _, student := range students {
		member := &models.EduClassMember{ClassID: class.ID, StudentID: student.ID, Status: "studying", TenantID: 1}
		if err := db.Create(member).Error; err != nil {
			t.Fatalf("create member: %v", err)
		}
	}
	lessonDate := models.NewJSONTime(time.Now())
	lesson := &models.EduLesson{LessonType: "class", LessonDate: &lessonDate, StartTime: "09:00", EndTime: "10:00", ClassID: class.ID, CourseID: course.ID, TeacherID: 10, Status: "scheduled", TenantID: 1}
	if err := db.Create(lesson).Error; err != nil {
		t.Fatalf("create lesson: %v", err)
	}
	completion := &models.EduLessonStudentCompletion{LessonID: lesson.ID, StudentID: students[0].ID, ClassID: class.ID, CourseID: course.ID, TeacherID: 10, ResultType: models.EduLessonCompletionResultAttended, Status: models.EduLessonCompletionStatusCompleted, TenantID: 1}
	if err := db.Create(completion).Error; err != nil {
		t.Fatalf("create completion: %v", err)
	}

	got, err := svc.GetLessonCompletionStatus(ctx, 1, lesson.ID)
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if got.Lesson == nil || got.Lesson.ID != lesson.ID {
		t.Fatalf("expected lesson in response, got %+v", got.Lesson)
	}
	if len(got.Students) != 2 {
		t.Fatalf("expected two students, got %d", len(got.Students))
	}
	if got.UnprocessedCount != 1 {
		t.Fatalf("expected one unprocessed student, got %d", got.UnprocessedCount)
	}
}

func seedCompletionLessonWithBenefit(t *testing.T, tenantID uint, remaining int) (*models.EduLesson, *models.EduStudentBenefit, uint) {
	t.Helper()
	db := app.DB()
	course := &models.EduCourse{Name: fmt.Sprintf("课程-%d", tenantID), Code: fmt.Sprintf("COMP-%d", tenantID), TenantID: tenantID}
	if err := db.Create(course).Error; err != nil {
		t.Fatalf("create course: %v", err)
	}
	student := &models.EduStudent{Name: fmt.Sprintf("学生-%d", tenantID), Phone: fmt.Sprintf("139%08d", tenantID), Status: 1, TenantID: tenantID}
	if err := db.Create(student).Error; err != nil {
		t.Fatalf("create student: %v", err)
	}
	benefit := &models.EduStudentBenefit{
		StudentID:       student.ID,
		BenefitType:     "course",
		CalculationMode: "count_limited",
		CourseID:        course.ID,
		TotalCount:      remaining,
		RemainingCount:  remaining,
		Status:          1,
		TenantID:        tenantID,
	}
	if err := db.Create(benefit).Error; err != nil {
		t.Fatalf("create benefit: %v", err)
	}
	lessonDate := models.NewJSONTime(time.Now())
	lesson := &models.EduLesson{
		LessonType: "one_to_one",
		LessonDate: &lessonDate,
		StartTime:  "09:00",
		EndTime:    "10:30",
		StudentID:  student.ID,
		CourseID:   course.ID,
		TeacherID:  10,
		Status:     "scheduled",
		TenantID:   tenantID,
	}
	if err := db.Create(lesson).Error; err != nil {
		t.Fatalf("create lesson: %v", err)
	}
	return lesson, benefit, student.ID
}

func completionSubmitRequest(lessonID, studentID uint, resultType string) *models.EduLessonCompletionSubmitRequest {
	return &models.EduLessonCompletionSubmitRequest{
		LessonID: models.FlexString(fmt.Sprintf("%d", lessonID)),
		Items: []models.EduLessonCompletionSubmitItemRequest{{
			StudentID:  models.FlexString(fmt.Sprintf("%d", studentID)),
			ResultType: resultType,
			Reason:     "课堂处理",
		}},
	}
}

func TestEduLessonCompletionServiceSubmit(t *testing.T) {
	t.Run("fixed count deducts benefit and writes ledger", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, benefit, studentID := seedCompletionLessonWithBenefit(t, 1, 5)
		items := validCompletionRuleItems()
		for i := range items {
			if items[i].ResultType == models.EduLessonCompletionResultAttended {
				items[i].DeductEnabled = 1
				items[i].FixedCount = 2
			}
		}
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		rows, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended))
		if err != nil {
			t.Fatalf("submit: %v", err)
		}
		if len(rows) != 1 || rows[0].Status != models.EduLessonCompletionStatusCompleted {
			t.Fatalf("unexpected completion rows: %+v", rows)
		}
		var updated models.EduStudentBenefit
		if err := app.DB().First(&updated, benefit.ID).Error; err != nil {
			t.Fatalf("load benefit: %v", err)
		}
		if updated.RemainingCount != 3 || updated.UsedCount != 2 {
			t.Fatalf("expected benefit remaining=3 used=2, got remaining=%d used=%d", updated.RemainingCount, updated.UsedCount)
		}
		var ledger models.EduBenefitLedger
		if err := app.DB().Where("student_id = ? AND action_type = ?", studentID, "lesson_complete").First(&ledger).Error; err != nil {
			t.Fatalf("load ledger: %v", err)
		}
		if ledger.ChangeCount != -2 || ledger.BizType != "lesson_completion" || ledger.BizID != rows[0].ID {
			t.Fatalf("unexpected ledger: %+v, completion id=%d", ledger, rows[0].ID)
		}
	})

	t.Run("duration mode rounds up by unit minutes", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, benefit, studentID := seedCompletionLessonWithBenefit(t, 1, 5)
		items := validCompletionRuleItems()
		for i := range items {
			if items[i].ResultType == models.EduLessonCompletionResultAttended {
				items[i].DeductEnabled = 1
				items[i].DeductMode = models.EduLessonCompletionDeductModeDuration
				items[i].DurationUnitMinutes = 60
			}
		}
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err != nil {
			t.Fatalf("submit: %v", err)
		}
		var updated models.EduStudentBenefit
		if err := app.DB().First(&updated, benefit.ID).Error; err != nil {
			t.Fatalf("load benefit: %v", err)
		}
		if updated.RemainingCount != 3 || updated.UsedCount != 2 {
			t.Fatalf("expected 90 minutes to deduct 2 units, got remaining=%d used=%d", updated.RemainingCount, updated.UsedCount)
		}
	})

	t.Run("insufficient block rejects submit", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, _, studentID := seedCompletionLessonWithBenefit(t, 1, 0)
		items := validCompletionRuleItems()
		for i := range items {
			if items[i].ResultType == models.EduLessonCompletionResultAttended {
				items[i].DeductEnabled = 1
				items[i].FixedCount = 1
				items[i].InsufficientPolicy = models.EduLessonCompletionInsufficientBlock
			}
		}
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err == nil {
			t.Fatalf("expected insufficient benefit to fail")
		}
	})

	t.Run("expired benefit is not deductible", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, benefit, studentID := seedCompletionLessonWithBenefit(t, 1, 5)
		expiredAt := models.NewJSONTime(time.Now().Add(-24 * time.Hour))
		if err := app.DB().Model(&models.EduStudentBenefit{}).Where("id = ?", benefit.ID).Update("valid_to", &expiredAt).Error; err != nil {
			t.Fatalf("expire benefit: %v", err)
		}
		items := validCompletionRuleItems()
		for i := range items {
			if items[i].ResultType == models.EduLessonCompletionResultAttended {
				items[i].DeductEnabled = 1
				items[i].FixedCount = 1
				items[i].InsufficientPolicy = models.EduLessonCompletionInsufficientBlock
			}
		}
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err == nil {
			t.Fatalf("expected expired benefit to fail")
		}
	})

	t.Run("non course benefit is not deductible", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, benefit, studentID := seedCompletionLessonWithBenefit(t, 1, 5)
		if err := app.DB().Model(&models.EduStudentBenefit{}).Where("id = ?", benefit.ID).Update("benefit_type", "welfare").Error; err != nil {
			t.Fatalf("update benefit type: %v", err)
		}
		items := validCompletionRuleItems()
		for i := range items {
			if items[i].ResultType == models.EduLessonCompletionResultAttended {
				items[i].DeductEnabled = 1
				items[i].FixedCount = 1
				items[i].InsufficientPolicy = models.EduLessonCompletionInsufficientBlock
			}
		}
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err == nil {
			t.Fatalf("expected non course benefit to fail")
		}
	})

	t.Run("benefit without matching course is not deductible", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, benefit, studentID := seedCompletionLessonWithBenefit(t, 1, 5)
		if err := app.DB().Model(&models.EduStudentBenefit{}).Where("id = ?", benefit.ID).Update("course_id", 0).Error; err != nil {
			t.Fatalf("clear benefit course: %v", err)
		}
		items := validCompletionRuleItems()
		for i := range items {
			if items[i].ResultType == models.EduLessonCompletionResultAttended {
				items[i].DeductEnabled = 1
				items[i].FixedCount = 1
				items[i].InsufficientPolicy = models.EduLessonCompletionInsufficientBlock
			}
		}
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err == nil {
			t.Fatalf("expected benefit without matching course to fail")
		}
	})

	t.Run("allow arrears creates arrears completion without ledger", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, _, studentID := seedCompletionLessonWithBenefit(t, 1, 0)
		items := validCompletionRuleItems()
		for i := range items {
			if items[i].ResultType == models.EduLessonCompletionResultAttended {
				items[i].DeductEnabled = 1
				items[i].FixedCount = 1
				items[i].InsufficientPolicy = models.EduLessonCompletionInsufficientAllowArrears
			}
		}
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err != nil {
			t.Fatalf("submit arrears: %v", err)
		}
		var completion models.EduLessonStudentCompletion
		if err := app.DB().Where("lesson_id = ? AND student_id = ?", lesson.ID, studentID).First(&completion).Error; err != nil {
			t.Fatalf("load completion: %v", err)
		}
		if completion.Status != models.EduLessonCompletionStatusArrears || completion.LedgerID != 0 || completion.StudentBenefitID != 0 {
			t.Fatalf("expected arrears without ledger, got %+v", completion)
		}
	})

	t.Run("duplicate active submit is blocked", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, _, studentID := seedCompletionLessonWithBenefit(t, 1, 5)
		items := validCompletionRuleItems()
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err != nil {
			t.Fatalf("first submit: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultStudentLeave)); err == nil {
			t.Fatalf("expected duplicate active submit to fail")
		}
	})

	t.Run("submits with default rules when tenant has no configured rules", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(3, 99)
		svc := NewEduLessonCompletionService()
		lesson, benefit, studentID := seedCompletionLessonWithBenefit(t, 3, 3)
		rows, err := svc.SubmitCompletions(ctx, 3, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended))
		if err != nil {
			t.Fatalf("submit with default rules: %v", err)
		}
		if len(rows) != 1 || rows[0].Status != models.EduLessonCompletionStatusCompleted {
			t.Fatalf("unexpected completion rows: %+v", rows)
		}
		var updated models.EduStudentBenefit
		if err := app.DB().First(&updated, benefit.ID).Error; err != nil {
			t.Fatalf("load benefit: %v", err)
		}
		if updated.RemainingCount != 2 || updated.UsedCount != 1 {
			t.Fatalf("expected default attended rule to deduct 1, got remaining=%d used=%d", updated.RemainingCount, updated.UsedCount)
		}
	})
}

func TestEduLessonCompletionServiceRevoke(t *testing.T) {
	t.Run("restores deducted benefit and allows resubmit", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, benefit, studentID := seedCompletionLessonWithBenefit(t, 1, 5)
		items := validCompletionRuleItems()
		for i := range items {
			if items[i].ResultType == models.EduLessonCompletionResultAttended {
				items[i].DeductEnabled = 1
				items[i].FixedCount = 2
			}
		}
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err != nil {
			t.Fatalf("submit: %v", err)
		}
		if err := svc.RevokeCompletion(ctx, 1, &models.EduLessonCompletionRevokeRequest{
			LessonID:  models.FlexString(fmt.Sprintf("%d", lesson.ID)),
			StudentID: models.FlexString(fmt.Sprintf("%d", studentID)),
			Reason:    "误操作",
		}); err != nil {
			t.Fatalf("revoke: %v", err)
		}
		var updated models.EduStudentBenefit
		if err := app.DB().First(&updated, benefit.ID).Error; err != nil {
			t.Fatalf("load benefit: %v", err)
		}
		if updated.RemainingCount != 5 || updated.UsedCount != 0 {
			t.Fatalf("expected benefit restored, got remaining=%d used=%d", updated.RemainingCount, updated.UsedCount)
		}
		var completion models.EduLessonStudentCompletion
		if err := app.DB().Where("lesson_id = ? AND student_id = ?", lesson.ID, studentID).First(&completion).Error; err != nil {
			t.Fatalf("load completion: %v", err)
		}
		if completion.Status != models.EduLessonCompletionStatusRevoked || completion.RevokeLedgerID == 0 {
			t.Fatalf("expected revoked completion with revoke ledger, got %+v", completion)
		}
		var revokeLedger models.EduBenefitLedger
		if err := app.DB().Where("id = ?", completion.RevokeLedgerID).First(&revokeLedger).Error; err != nil {
			t.Fatalf("load revoke ledger: %v", err)
		}
		if revokeLedger.ChangeCount != 2 || revokeLedger.ActionType != "lesson_revoke" || revokeLedger.BizID != completion.ID {
			t.Fatalf("unexpected revoke ledger: %+v", revokeLedger)
		}

		rows, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultStudentLeave))
		if err != nil {
			t.Fatalf("resubmit after revoke: %v", err)
		}
		if len(rows) != 1 || rows[0].Version != 2 {
			t.Fatalf("expected resubmitted completion version 2, got %+v", rows)
		}
	})

	t.Run("revokes arrears without restore ledger", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, _, studentID := seedCompletionLessonWithBenefit(t, 1, 0)
		items := validCompletionRuleItems()
		for i := range items {
			if items[i].ResultType == models.EduLessonCompletionResultAttended {
				items[i].DeductEnabled = 1
				items[i].InsufficientPolicy = models.EduLessonCompletionInsufficientAllowArrears
			}
		}
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err != nil {
			t.Fatalf("submit arrears: %v", err)
		}
		if err := svc.RevokeCompletion(ctx, 1, &models.EduLessonCompletionRevokeRequest{
			LessonID:  models.FlexString(fmt.Sprintf("%d", lesson.ID)),
			StudentID: models.FlexString(fmt.Sprintf("%d", studentID)),
			Reason:    "欠费记录误操作",
		}); err != nil {
			t.Fatalf("revoke arrears: %v", err)
		}
		var completion models.EduLessonStudentCompletion
		if err := app.DB().Where("lesson_id = ? AND student_id = ?", lesson.ID, studentID).First(&completion).Error; err != nil {
			t.Fatalf("load completion: %v", err)
		}
		if completion.Status != models.EduLessonCompletionStatusRevoked || completion.RevokeLedgerID != 0 {
			t.Fatalf("expected revoked arrears without ledger, got %+v", completion)
		}
		var count int64
		if err := app.DB().Model(&models.EduBenefitLedger{}).Where("action_type = ?", "lesson_revoke").Count(&count).Error; err != nil {
			t.Fatalf("count revoke ledgers: %v", err)
		}
		if count != 0 {
			t.Fatalf("expected no revoke ledger for arrears, got %d", count)
		}
	})
}

func TestEduLessonCompletionServiceCompleteLesson(t *testing.T) {
	t.Run("blocks when students are unprocessed", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, _, _ := seedCompletionLessonWithBenefit(t, 1, 5)
		if err := svc.CompleteLesson(ctx, 1, &models.EduLessonCompleteRequest{LessonID: models.FlexString(fmt.Sprintf("%d", lesson.ID))}); err == nil {
			t.Fatalf("expected unprocessed lesson to fail")
		}
	})

	t.Run("updates lesson to completed after all students processed", func(t *testing.T) {
		setupEduTestDB(t)
		ctx := contextWithTenantAndUser(1, 99)
		svc := NewEduLessonCompletionService()
		lesson, _, studentID := seedCompletionLessonWithBenefit(t, 1, 5)
		items := validCompletionRuleItems()
		if err := svc.SaveRules(ctx, 1, &models.EduLessonCompletionRuleSaveRequest{Items: items}); err != nil {
			t.Fatalf("save rules: %v", err)
		}
		if _, err := svc.SubmitCompletions(ctx, 1, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err != nil {
			t.Fatalf("submit: %v", err)
		}
		if err := svc.CompleteLesson(ctx, 1, &models.EduLessonCompleteRequest{LessonID: models.FlexString(fmt.Sprintf("%d", lesson.ID)), Reason: "结课"}); err != nil {
			t.Fatalf("complete lesson: %v", err)
		}
		var updated models.EduLesson
		if err := app.DB().First(&updated, lesson.ID).Error; err != nil {
			t.Fatalf("load lesson: %v", err)
		}
		if updated.Status != "completed" {
			t.Fatalf("expected completed lesson, got %s", updated.Status)
		}
		var log models.EduLessonChangeLog
		if err := app.DB().Where("lesson_id = ? AND action_type = ?", lesson.ID, "complete").First(&log).Error; err != nil {
			t.Fatalf("load change log: %v", err)
		}
	})
}

func TestEduLessonCompletionServiceTenantIsolation(t *testing.T) {
	setupEduTestDB(t)
	ctxTenant2 := contextWithTenantAndUser(2, 99)
	svc := NewEduLessonCompletionService()
	lesson, _, studentID := seedCompletionLessonWithBenefit(t, 1, 5)
	if err := svc.SaveRules(ctxTenant2, 2, &models.EduLessonCompletionRuleSaveRequest{Items: validCompletionRuleItems()}); err != nil {
		t.Fatalf("save tenant 2 rules: %v", err)
	}
	if _, err := svc.GetLessonCompletionStatus(ctxTenant2, 2, lesson.ID); err == nil {
		t.Fatalf("expected tenant 2 to be unable to read tenant 1 lesson")
	}
	if _, err := svc.SubmitCompletions(ctxTenant2, 2, completionSubmitRequest(lesson.ID, studentID, models.EduLessonCompletionResultAttended)); err == nil {
		t.Fatalf("expected tenant 2 to be unable to submit tenant 1 lesson")
	}
}
