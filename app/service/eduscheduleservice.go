package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gin-fast/app/global/app"
	"gin-fast/app/models"

	"gorm.io/gorm"
)

type EduScheduleService struct{}

func NewEduScheduleService() *EduScheduleService { return &EduScheduleService{} }

func (s *EduScheduleService) CreateRule(ctx context.Context, rule *models.EduScheduleRule) error {
	return app.DB().WithContext(ctx).Create(rule).Error
}

func (s *EduScheduleService) UpdateRule(ctx context.Context, rule *models.EduScheduleRule) error {
	return app.DB().WithContext(ctx).Save(rule).Error
}

func (s *EduScheduleService) DeleteRule(ctx context.Context, tenantID, id uint) error {
	return requireTenant(app.DB().WithContext(ctx).Where("id = ?", id), tenantID).Delete(&models.EduScheduleRule{}).Error
}

func (s *EduScheduleService) PreviewRuleChange(ctx context.Context, tenantID, id uint) ([]models.EduLesson, error) {
	var rows []models.EduLesson
	err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduLesson{}), tenantID).Where("rule_id = ?", id).Order("lesson_date asc, id asc").Find(&rows).Error
	return rows, err
}

func (s *EduScheduleService) GenerateLessonsForRule(ctx context.Context, ruleID, operatorID uint) ([]models.EduLesson, error) {
	tenantID := tenantIDFromContext(ctx)
	if tenantID == 0 {
		return nil, errors.New("缺少租户信息")
	}
	var rule models.EduScheduleRule
	if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduScheduleRule{}), tenantID).Where("id = ?", ruleID).First(&rule).Error; err != nil {
		return nil, err
	}
	return s.generateLessonsForRuleTx(ctx, &rule, operatorID)
}

func (s *EduScheduleService) RegenerateFutureLessons(ctx context.Context, tenantID, ruleID, operatorID uint) ([]models.EduLesson, error) {
	var rule models.EduScheduleRule
	if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduScheduleRule{}), tenantID).Where("id = ?", ruleID).First(&rule).Error; err != nil {
		return nil, err
	}
	return s.generateLessonsForRuleTx(ctx, &rule, operatorID)
}

func (s *EduScheduleService) ListLessons(ctx context.Context, tenantID uint, req *models.EduLessonListRequest) ([]models.EduLesson, error) {
	db := requireTenant(app.DB().WithContext(ctx).Model(&models.EduLesson{}), tenantID)
	if req != nil {
		if req.ID != nil {
			db = db.Where("id = ?", *req.ID)
		}
		if req.RuleID != nil {
			db = db.Where("rule_id = ?", *req.RuleID)
		}
		if req.ClassID != nil {
			db = db.Where("class_id = ?", *req.ClassID)
		}
		if req.StudentID != nil {
			db = db.Where("student_id = ?", *req.StudentID)
		}
		if req.CourseID != nil {
			db = db.Where("course_id = ?", *req.CourseID)
		}
		if req.TeacherID != nil {
			db = db.Where("teacher_id = ?", *req.TeacherID)
		}
		if strings.TrimSpace(req.Status) != "" {
			db = db.Where("status = ?", req.Status)
		}
	}
	var rows []models.EduLesson
	return rows, db.Order("lesson_date asc, id asc").Find(&rows).Error
}

func (s *EduScheduleService) CalendarLessons(ctx context.Context, tenantID uint, req *models.EduLessonCalendarRequest) ([]models.EduLesson, error) {
	db := requireTenant(app.DB().WithContext(ctx).Model(&models.EduLesson{}), tenantID)
	if req != nil {
		if req.StartDate != nil && !req.StartDate.Time.IsZero() {
			db = db.Where("lesson_date >= ?", req.StartDate.Time)
		}
		if req.EndDate != nil && !req.EndDate.Time.IsZero() {
			db = db.Where("lesson_date <= ?", req.EndDate.Time)
		}
		if req.TeacherID != nil {
			db = db.Where("teacher_id = ?", *req.TeacherID)
		}
		if req.ClassID != nil {
			db = db.Where("class_id = ?", *req.ClassID)
		}
	}
	var rows []models.EduLesson
	return rows, db.Order("lesson_date asc, start_time asc, id asc").Find(&rows).Error
}

func (s *EduScheduleService) generateLessonsForRuleTx(ctx context.Context, rule *models.EduScheduleRule, operatorID uint) ([]models.EduLesson, error) {
	if rule == nil {
		return nil, errors.New("规则不能为空")
	}
	tenantID := rule.TenantID
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
	}
	if tenantID == 0 {
		return nil, errors.New("缺少租户信息")
	}
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	if rule.CourseID == 0 || rule.TeacherID == 0 {
		return nil, errors.New("课程ID和教师ID不能为空")
	}
	if strings.TrimSpace(rule.StartTime) == "" || strings.TrimSpace(rule.EndTime) == "" {
		return nil, errors.New("开始时间和结束时间不能为空")
	}
	if err := validateLessonRoom(rule); err != nil {
		return nil, err
	}
	startDate, endDate, err := s.resolveScheduleWindow(ctx, rule)
	if err != nil {
		return nil, err
	}
	lessonDates, err := s.buildLessonDates(ctx, rule, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var lessons []models.EduLesson
	err = app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, lessonDate := range lessonDates {
			lesson := models.EduLesson{
				RuleID:       rule.ID,
				RuleVersion:  rule.Version,
				LessonType:   rule.RuleType,
				LessonDate:   &models.JSONTime{Time: lessonDate},
				StartTime:    rule.StartTime,
				EndTime:      rule.EndTime,
				ClassID:      rule.ClassID,
				StudentID:    rule.StudentID,
				CourseID:     rule.CourseID,
				TeacherID:    rule.TeacherID,
				TeachingMode: rule.TeachingMode,
				RequiresRoom: rule.RequiresRoom,
				RoomID:       rule.RoomID,
				Status:       "scheduled",
				CreatedBy:    operatorID,
				TenantID:     tenantID,
			}
			if err := tx.Create(&lesson).Error; err != nil {
				return err
			}
			if err := s.writeEligibilityRowsTx(tx, ctx, &lesson, rule); err != nil {
				return err
			}
			lessons = append(lessons, lesson)
		}
		return nil
	})
	return lessons, err
}

func (s *EduScheduleService) resolveScheduleWindow(ctx context.Context, rule *models.EduScheduleRule) (time.Time, time.Time, error) {
	var startDate, endDate time.Time
	if rule.StartDate != nil && !rule.StartDate.Time.IsZero() {
		startDate = rule.StartDate.Time
	}
	if rule.EndDate != nil && !rule.EndDate.Time.IsZero() {
		endDate = rule.EndDate.Time
	}
	if rule.TermID != 0 {
		var term models.EduTerm
		if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduTerm{}), rule.TenantID).Where("id = ?", rule.TermID).First(&term).Error; err != nil {
			return time.Time{}, time.Time{}, err
		}
		if startDate.IsZero() && term.StartDate != nil {
			startDate = term.StartDate.Time
		}
		if endDate.IsZero() && term.EndDate != nil {
			endDate = term.EndDate.Time
		}
	}
	if startDate.IsZero() {
		return time.Time{}, time.Time{}, errors.New("开始日期不能为空")
	}
	if endDate.IsZero() {
		endDate = startDate
	}
	if endDate.Before(startDate) {
		return time.Time{}, time.Time{}, errors.New("结束日期不能早于开始日期")
	}
	return startDate, endDate, nil
}

func (s *EduScheduleService) buildLessonDates(ctx context.Context, rule *models.EduScheduleRule, startDate, endDate time.Time) ([]time.Time, error) {
	switch strings.TrimSpace(rule.RepeatType) {
	case "single":
		return []time.Time{truncateDate(startDate)}, nil
	case "weekly":
		closedDates, err := s.loadClosedDates(ctx, rule.TenantID, rule.TermID)
		if err != nil {
			return nil, err
		}
		var dates []time.Time
		for current := truncateDate(startDate); !current.After(endDate); current = current.AddDate(0, 0, 1) {
			if int(current.Weekday()) == int(rule.Weekday)%7 && rule.Weekday > 0 {
				if _, closed := closedDates[current.Format("2006-01-02")]; closed {
					continue
				}
				dates = append(dates, current)
			}
		}
		return dates, nil
	default:
		return nil, fmt.Errorf("不支持的重复类型")
	}
}

func (s *EduScheduleService) loadClosedDates(ctx context.Context, tenantID, termID uint) (map[string]struct{}, error) {
	result := map[string]struct{}{}
	if termID == 0 {
		return result, nil
	}
	var rows []models.EduTermClosedDay
	if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduTermClosedDay{}), tenantID).Where("term_id = ?", termID).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row.ClosedDate == nil || row.ClosedDate.Time.IsZero() {
			continue
		}
		result[row.ClosedDate.Time.Format("2006-01-02")] = struct{}{}
	}
	return result, nil
}

func (s *EduScheduleService) writeEligibilityRowsTx(tx *gorm.DB, ctx context.Context, lesson *models.EduLesson, rule *models.EduScheduleRule) error {
	if lesson == nil || rule == nil {
		return nil
	}
	lessonDate := lesson.LessonDate
	if lessonDate == nil || lessonDate.Time.IsZero() {
		return errors.New("课次日期不能为空")
	}
	studentIDs, err := s.resolveLessonStudentIDsTx(tx, ctx, rule, lessonDate.Time)
	if err != nil {
		return err
	}
	classPolicy := ""
	if rule.ClassID != 0 && rule.RuleType == "class" {
		classPolicy, err = s.classBenefitPolicyTx(tx, rule.TenantID, rule.ClassID)
		if err != nil {
			return err
		}
		if classPolicy == "none" {
			return nil
		}
	}
	for _, studentID := range studentIDs {
		check, err := checkEligibilityTx(tx, lessonDate, rule.TenantID, studentID, rule.CourseID, rule.ClassID, rule.TeacherID)
		if err != nil {
			return err
		}
		status := check.Status
		if status == "" {
			status = "ineligible"
		}
		eligibility := &models.EduLessonStudentEligibility{
			LessonID:          lesson.ID,
			StudentID:         studentID,
			CourseID:          rule.CourseID,
			ClassID:           rule.ClassID,
			StudentBenefitID:  check.StudentBenefitID,
			EligibilityStatus: status,
			ReasonCode:        check.ReasonCode,
			CheckedAt:         lessonDate,
			TenantID:          rule.TenantID,
		}
		if err := tx.Create(eligibility).Error; err != nil {
			return err
		}
		if rule.ClassID != 0 && rule.RuleType == "class" && !check.Eligible && classPolicy == "required" {
			return fmt.Errorf("学生权益不足")
		}
	}
	return nil
}

func checkEligibilityTx(tx *gorm.DB, checkedAt *models.JSONTime, tenantID, studentID, courseID, classID, teacherID uint) (*models.EduBenefitCheckResult, error) {
	if tx == nil {
		return nil, errors.New("数据库不能为空")
	}
	if checkedAt == nil || checkedAt.Time.IsZero() {
		now := models.NewJSONTime(time.Now().UTC())
		checkedAt = &now
	}
	var benefits []models.EduStudentBenefit
	err := requireTenant(tx.Model(&models.EduStudentBenefit{}), tenantID).
		Where("student_id = ? AND course_id = ? AND status = 1 AND benefit_type = ?", studentID, courseID, "course").
		Order("id asc").
		Find(&benefits).Error
	if err != nil {
		return nil, err
	}
	if len(benefits) == 0 {
		return &models.EduBenefitCheckResult{Eligible: false, Status: "ineligible", ReasonCode: "no_benefit"}, nil
	}
	var expiredResult *models.EduBenefitCheckResult
	var insufficientResult *models.EduBenefitCheckResult
	for i := range benefits {
		benefit := benefits[i]
		if !studentBenefitMatchesDimension(&benefit, classID, teacherID) {
			continue
		}
		if benefit.ValidFrom != nil && checkedAt.Time.Before(benefit.ValidFrom.Time) {
			if expiredResult == nil {
				expiredResult = &models.EduBenefitCheckResult{Eligible: false, Status: "ineligible", ReasonCode: "expired", StudentBenefitID: benefit.ID}
			}
			continue
		}
		if benefit.ValidTo != nil && checkedAt.Time.After(benefit.ValidTo.Time) {
			if expiredResult == nil {
				expiredResult = &models.EduBenefitCheckResult{Eligible: false, Status: "ineligible", ReasonCode: "expired", StudentBenefitID: benefit.ID}
			}
			continue
		}
		if benefit.CalculationMode == "count_limited" && benefit.RemainingCount <= 0 {
			if insufficientResult == nil {
				insufficientResult = &models.EduBenefitCheckResult{Eligible: false, Status: "ineligible", ReasonCode: "insufficient_count", StudentBenefitID: benefit.ID}
			}
			continue
		}
		return &models.EduBenefitCheckResult{Eligible: true, Status: "eligible", ReasonCode: "", StudentBenefitID: benefit.ID}, nil
	}
	if expiredResult != nil {
		return expiredResult, nil
	}
	if insufficientResult != nil {
		return insufficientResult, nil
	}
	return &models.EduBenefitCheckResult{Eligible: false, Status: "ineligible", ReasonCode: "no_benefit"}, nil
}

func (s *EduScheduleService) resolveLessonStudentIDsTx(tx *gorm.DB, ctx context.Context, rule *models.EduScheduleRule, lessonDate time.Time) ([]uint, error) {
	if rule.RuleType == "one_to_one" || rule.StudentID != 0 {
		return []uint{rule.StudentID}, nil
	}
	if rule.ClassID == 0 {
		return nil, errors.New("班级ID不能为空")
	}
	var members []models.EduClassMember
	if err := requireTenant(tx.Model(&models.EduClassMember{}), rule.TenantID).Where("class_id = ?", rule.ClassID).Find(&members).Error; err != nil {
		return nil, err
	}
	var studentIDs []uint
	for _, member := range members {
		if !isActiveMemberStatus(member.Status) {
			continue
		}
		if member.JoinDate != nil && lessonDate.Before(member.JoinDate.Time) {
			continue
		}
		if member.LeaveDate != nil && !member.LeaveDate.Time.IsZero() && lessonDate.After(member.LeaveDate.Time) {
			continue
		}
		studentIDs = append(studentIDs, member.StudentID)
	}
	sort.Slice(studentIDs, func(i, j int) bool { return studentIDs[i] < studentIDs[j] })
	return studentIDs, nil
}

func (s *EduScheduleService) classBenefitPolicyTx(tx *gorm.DB, tenantID, classID uint) (string, error) {
	var class models.EduClass
	if err := requireTenant(tx.Model(&models.EduClass{}), tenantID).Where("id = ?", classID).First(&class).Error; err != nil {
		return "", err
	}
	if strings.TrimSpace(class.BenefitCheckPolicy) == "" {
		return "required", nil
	}
	return strings.ToLower(strings.TrimSpace(class.BenefitCheckPolicy)), nil
}

func validateLessonRoom(rule *models.EduScheduleRule) error {
	if rule == nil {
		return errors.New("规则不能为空")
	}
	if strings.EqualFold(strings.TrimSpace(rule.TeachingMode), "online") {
		return nil
	}
	if rule.RequiresRoom == 1 && rule.RoomID == 0 {
		return errors.New("线下课需要场地")
	}
	return nil
}

func truncateDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
