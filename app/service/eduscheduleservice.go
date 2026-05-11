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
	db := requireTenant(app.DB().WithContext(ctx).Model(&models.EduLesson{}), tenantID).Where("rule_id = ?", id)
	var rule models.EduScheduleRule
	if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduScheduleRule{}), tenantID).Where("id = ?", id).First(&rule).Error; err != nil {
		return nil, err
	}
	if rule.EffectiveFrom != nil && !rule.EffectiveFrom.Time.IsZero() {
		db = db.Where("lesson_date >= ?", truncateDate(rule.EffectiveFrom.Time))
	}
	db = db.Where("status NOT IN ?", []string{"completed", "canceled", "stopped"}).
		Where("is_manual_adjusted = 0")
	err := db.Order("lesson_date asc, id asc").Find(&rows).Error
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
	return s.generateLessonsForRuleTxWithDB(ctx, app.DB().WithContext(ctx), &rule, operatorID)
}

func (s *EduScheduleService) RegenerateFutureLessons(ctx context.Context, tenantID, ruleID, operatorID uint) ([]models.EduLesson, error) {
	var rule models.EduScheduleRule
	if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduScheduleRule{}), tenantID).Where("id = ?", ruleID).First(&rule).Error; err != nil {
		return nil, err
	}
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	if rule.EffectiveFrom == nil || rule.EffectiveFrom.Time.IsZero() {
		return s.generateLessonsForRuleTxWithDB(ctx, app.DB().WithContext(ctx), &rule, operatorID)
	}
	effectiveFrom := truncateDate(rule.EffectiveFrom.Time)
	newRuleVersion := rule.Version
	if newRuleVersion < 1 {
		newRuleVersion = 1
	}

	var generated []models.EduLesson
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var replaceable []models.EduLesson
		query := requireTenant(tx.Model(&models.EduLesson{}), tenantID).
			Where("rule_id = ?", ruleID).
			Where("lesson_date >= ?", effectiveFrom).
			Where("status NOT IN ?", []string{"completed", "canceled", "stopped"}).
			Where("is_manual_adjusted = 0").
			Order("lesson_date asc, id asc")
		if err := query.Find(&replaceable).Error; err != nil {
			return err
		}
		if len(replaceable) > 0 {
			ids := make([]uint, 0, len(replaceable))
			for _, lesson := range replaceable {
				ids = append(ids, lesson.ID)
			}
			if err := requireTenant(tx.Where("id IN ?", ids), tenantID).Delete(&models.EduLesson{}).Error; err != nil {
				return err
			}
			for _, lesson := range replaceable {
				beforeData := fmt.Sprintf("lesson_id=%d,version=%d,status=%s,date=%s", lesson.ID, lesson.RuleVersion, lesson.Status, lessonDateString(lesson.LessonDate))
				afterData := fmt.Sprintf("deleted_for_regeneration=true,rule_version=%d", newRuleVersion)
				log := &models.EduLessonChangeLog{
					LessonID:   lesson.ID,
					RuleID:     ruleID,
					ActionType: "rule_regenerate",
					BeforeData: beforeData,
					AfterData:  afterData,
					Reason:     "规则变更重算",
					OperatorID: operatorID,
					OccurredAt: &models.JSONTime{Time: time.Now().UTC()},
					TenantID:   tenantID,
				}
				if err := tx.Create(log).Error; err != nil {
					return err
				}
			}
		}
		dates := make([]time.Time, 0, len(replaceable))
		for _, lesson := range replaceable {
			if lesson.LessonDate == nil || lesson.LessonDate.Time.IsZero() {
				continue
			}
			dates = append(dates, truncateDate(lesson.LessonDate.Time))
		}
		rule.Version = newRuleVersion
		created, err := s.generateLessonsForDatesTx(ctx, tx, &rule, operatorID, dates)
		if err != nil {
			return err
		}
		generated = created
		return nil
	})
	return generated, err
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

func (s *EduScheduleService) CheckConflicts(ctx context.Context, req *models.EduScheduleConflictCheckRequest) (*models.EduScheduleConflictCheckResult, error) {
	return s.checkConflictsWithDB(ctx, app.DB().WithContext(ctx), req)
}

func (s *EduScheduleService) checkConflictsWithDB(ctx context.Context, db *gorm.DB, req *models.EduScheduleConflictCheckRequest) (*models.EduScheduleConflictCheckResult, error) {
	if req == nil {
		return nil, errors.New("请求不能为空")
	}
	tenantID := tenantIDFromContext(ctx)
	if tenantID == 0 {
		return nil, errors.New("缺少租户信息")
	}
	if req.LessonDate == nil || req.LessonDate.Time.IsZero() {
		return nil, errors.New("课次日期不能为空")
	}
	newStart, err := parseHHMM(req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("开始时间格式错误: %w", err)
	}
	newEnd, err := parseHHMM(req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("结束时间格式错误: %w", err)
	}
	if newEnd <= newStart {
		return nil, errors.New("结束时间必须晚于开始时间")
	}

	var lessons []models.EduLesson
	query := requireTenant(db.Model(&models.EduLesson{}), tenantID).
		Where("lesson_date = ? AND status = ?", truncateDate(req.LessonDate.Time), "scheduled")
	if req.LessonID != 0 {
		query = query.Where("id <> ?", req.LessonID)
	}
	if err := query.Order("id asc").Find(&lessons).Error; err != nil {
		return nil, err
	}
	result := &models.EduScheduleConflictCheckResult{}
	seen := map[string]struct{}{}
	for _, lesson := range lessons {
		existingStart, err := parseHHMM(lesson.StartTime)
		if err != nil {
			return nil, fmt.Errorf("既有课次开始时间格式错误: %w", err)
		}
		existingEnd, err := parseHHMM(lesson.EndTime)
		if err != nil {
			return nil, fmt.Errorf("既有课次结束时间格式错误: %w", err)
		}
		if !(newStart < existingEnd && existingStart < newEnd) {
			continue
		}
		appendConflict(result, seen, "teacher", fmt.Sprintf("%d", req.TeacherID), req.TeacherID != 0 && lesson.TeacherID == req.TeacherID)
		appendConflict(result, seen, "room", fmt.Sprintf("%d", req.RoomID), req.RequiresRoom == 1 && req.RoomID != 0 && lesson.RequiresRoom == 1 && lesson.RoomID == req.RoomID)
		appendConflict(result, seen, "student", fmt.Sprintf("%d", req.StudentID), req.StudentID != 0 && lesson.StudentID == req.StudentID)
		appendConflict(result, seen, "class", fmt.Sprintf("%d", req.ClassID), req.ClassID != 0 && lesson.ClassID == req.ClassID)
	}
	if !result.HasConflict {
		return result, nil
	}
	if !req.AllowConflictOverride {
		return result, errors.New("存在排课冲突")
	}
	if strings.TrimSpace(req.OverrideReason) == "" {
		return result, errors.New("覆盖冲突必须填写原因")
	}
	now := models.NewJSONTime(time.Now().UTC())
	for _, item := range result.Items {
		override := &models.EduScheduleConflictOverride{
			RuleID:       req.RuleID,
			LessonID:     req.LessonID,
			ConflictType: item.ConflictType,
			ConflictKey:  item.ConflictKey,
			Reason:       strings.TrimSpace(req.OverrideReason),
			OperatorID:   req.OperatorID,
			OccurredAt:   &now,
			TenantID:     tenantID,
		}
		if err := db.Create(override).Error; err != nil {
			return result, err
		}
	}
	return result, nil
}

func (s *EduScheduleService) RescheduleLesson(ctx context.Context, tenantID uint, req *models.EduLessonRescheduleRequest) error {
	if req == nil {
		return errors.New("请求不能为空")
	}
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
	}
	if tenantID == 0 {
		return errors.New("缺少租户信息")
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	newDate := req.LessonDate
	if newDate == nil || newDate.Time.IsZero() {
		return errors.New("新日期不能为空")
	}
	newStart, err := parseHHMM(req.StartTime)
	if err != nil {
		return fmt.Errorf("开始时间格式错误: %w", err)
	}
	newEnd, err := parseHHMM(req.EndTime)
	if err != nil {
		return fmt.Errorf("结束时间格式错误: %w", err)
	}
	if newEnd <= newStart {
		return errors.New("结束时间必须晚于开始时间")
	}
	if strings.EqualFold(strings.TrimSpace(req.TeachingMode), "offline") {
		if req.RoomID == nil || *req.RoomID == 0 {
			return errors.New("线下课需要场地")
		}
	}

	var conflictResult *models.EduScheduleConflictCheckResult
	err = app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lesson models.EduLesson
		if err := requireTenant(tx.Model(&models.EduLesson{}), tenantID).Where("id = ?", req.LessonID).First(&lesson).Error; err != nil {
			return err
		}
		if strings.EqualFold(strings.TrimSpace(lesson.Status), "completed") {
			return errors.New("已完成课次禁止调课")
		}
		if strings.EqualFold(strings.TrimSpace(lesson.Status), "canceled") || strings.EqualFold(strings.TrimSpace(lesson.Status), "stopped") {
			return errors.New("当前课次状态不允许调课")
		}

		roomID := lesson.RoomID
		requiresRoom := lesson.RequiresRoom
		if strings.EqualFold(strings.TrimSpace(req.TeachingMode), "offline") {
			requiresRoom = 1
			if req.RoomID != nil {
				roomID = *req.RoomID
			}
		} else {
			requiresRoom = 0
			roomID = 0
		}

		conflictResult, err = s.checkConflictsWithDB(ctx, tx, &models.EduScheduleConflictCheckRequest{
			LessonID:              lesson.ID,
			RuleID:                lesson.RuleID,
			LessonDate:            newDate,
			StartTime:             req.StartTime,
			EndTime:               req.EndTime,
			ClassID:               lesson.ClassID,
			StudentID:             lesson.StudentID,
			CourseID:              lesson.CourseID,
			TeacherID:             req.TeacherID,
			TeachingMode:          req.TeachingMode,
			RequiresRoom:          requiresRoom,
			RoomID:                roomID,
			AllowConflictOverride: req.AllowConflictOverride,
			OverrideReason:        req.OverrideReason,
			OperatorID:            0,
		})
		if err != nil {
			return err
		}
		if conflictResult != nil && conflictResult.HasConflict && !req.AllowConflictOverride {
			return errors.New("存在排课冲突")
		}

		beforeData := fmt.Sprintf("lesson_date=%s,start_time=%s,end_time=%s,teacher_id=%d,room_id=%d,teaching_mode=%s,requires_room=%d,status=%s,is_manual_adjusted=%d",
			lessonDateString(lesson.LessonDate), lesson.StartTime, lesson.EndTime, lesson.TeacherID, lesson.RoomID, lesson.TeachingMode, lesson.RequiresRoom, lesson.Status, lesson.IsManualAdjusted)
		lesson.LessonDate = &models.JSONTime{Time: truncateDate(newDate.Time)}
		lesson.StartTime = req.StartTime
		lesson.EndTime = req.EndTime
		lesson.TeacherID = req.TeacherID
		lesson.TeachingMode = strings.TrimSpace(req.TeachingMode)
		lesson.RequiresRoom = requiresRoom
		lesson.RoomID = roomID
		lesson.IsManualAdjusted = 1
		if err := tx.Save(&lesson).Error; err != nil {
			return err
		}

		afterData := fmt.Sprintf("lesson_date=%s,start_time=%s,end_time=%s,teacher_id=%d,room_id=%d,teaching_mode=%s,requires_room=%d,status=%s,is_manual_adjusted=%d",
			lessonDateString(lesson.LessonDate), lesson.StartTime, lesson.EndTime, lesson.TeacherID, lesson.RoomID, lesson.TeachingMode, lesson.RequiresRoom, lesson.Status, lesson.IsManualAdjusted)
		changeLog := &models.EduLessonChangeLog{
			LessonID:   lesson.ID,
			RuleID:     lesson.RuleID,
			ActionType: "reschedule",
			BeforeData: beforeData,
			AfterData:  afterData,
			Reason:     strings.TrimSpace(req.OverrideReason),
			OccurredAt: &models.JSONTime{Time: time.Now().UTC()},
			TenantID:   tenantID,
		}
		if err := tx.Create(changeLog).Error; err != nil {
			return err
		}

		if lesson.StudentID != 0 {
			if err := s.recheckLessonEligibilityTx(tx, ctx, &lesson); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *EduScheduleService) StopLesson(ctx context.Context, tenantID uint, req *models.EduLessonStopRequest) error {
	return s.changeLessonStatus(ctx, tenantID, req.LessonID, "stopped", "stop", req.Reason, false)
}

func (s *EduScheduleService) CancelLesson(ctx context.Context, tenantID uint, req *models.EduLessonCancelRequest) error {
	return s.changeLessonStatus(ctx, tenantID, req.LessonID, "canceled", "cancel", req.Reason, false)
}

func (s *EduScheduleService) RestoreLesson(ctx context.Context, tenantID uint, req *models.EduLessonRestoreRequest) error {
	return s.changeLessonStatus(ctx, tenantID, req.LessonID, "scheduled", "restore", req.Reason, true)
}

func (s *EduScheduleService) changeLessonStatus(ctx context.Context, tenantID, lessonID uint, targetStatus, actionType, reason string, recheckConflict bool) error {
	if lessonID == 0 {
		return errors.New("课次ID不能为空")
	}
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
	}
	if tenantID == 0 {
		return errors.New("缺少租户信息")
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return errors.New("原因不能为空")
	}

	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lesson models.EduLesson
		if err := requireTenant(tx.Model(&models.EduLesson{}), tenantID).Where("id = ?", lessonID).First(&lesson).Error; err != nil {
			return err
		}
		if strings.EqualFold(strings.TrimSpace(lesson.Status), "completed") {
			return errors.New("已完成课次禁止变更状态")
		}
		if !recheckConflict && strings.EqualFold(strings.TrimSpace(lesson.Status), targetStatus) {
			return nil
		}
		if recheckConflict {
			if strings.TrimSpace(lesson.Status) != "stopped" && strings.TrimSpace(lesson.Status) != "canceled" {
				return errors.New("当前课次状态不允许恢复")
			}
			conflictResult, err := s.checkConflictsWithDB(ctx, tx, &models.EduScheduleConflictCheckRequest{
				LessonID:     lesson.ID,
				RuleID:       lesson.RuleID,
				LessonDate:   lesson.LessonDate,
				StartTime:    lesson.StartTime,
				EndTime:      lesson.EndTime,
				ClassID:      lesson.ClassID,
				StudentID:    lesson.StudentID,
				CourseID:     lesson.CourseID,
				TeacherID:    lesson.TeacherID,
				TeachingMode: lesson.TeachingMode,
				RequiresRoom: lesson.RequiresRoom,
				RoomID:       lesson.RoomID,
			})
			if err != nil {
				return err
			}
			if conflictResult != nil && conflictResult.HasConflict {
				return errors.New("存在排课冲突")
			}
		} else {
			if strings.EqualFold(strings.TrimSpace(lesson.Status), "stopped") || strings.EqualFold(strings.TrimSpace(lesson.Status), "canceled") {
				// 允许重复停课/取消返回成功，保持幂等性。
			} else if strings.TrimSpace(lesson.Status) != "scheduled" {
				return errors.New("当前课次状态不允许变更")
			}
		}

		beforeData := fmt.Sprintf("lesson_date=%s,start_time=%s,end_time=%s,teacher_id=%d,room_id=%d,teaching_mode=%s,requires_room=%d,status=%s,is_manual_adjusted=%d",
			lessonDateString(lesson.LessonDate), lesson.StartTime, lesson.EndTime, lesson.TeacherID, lesson.RoomID, lesson.TeachingMode, lesson.RequiresRoom, lesson.Status, lesson.IsManualAdjusted)
		lesson.Status = targetStatus
		if err := tx.Save(&lesson).Error; err != nil {
			return err
		}
		afterData := fmt.Sprintf("lesson_date=%s,start_time=%s,end_time=%s,teacher_id=%d,room_id=%d,teaching_mode=%s,requires_room=%d,status=%s,is_manual_adjusted=%d",
			lessonDateString(lesson.LessonDate), lesson.StartTime, lesson.EndTime, lesson.TeacherID, lesson.RoomID, lesson.TeachingMode, lesson.RequiresRoom, lesson.Status, lesson.IsManualAdjusted)
		changeLog := &models.EduLessonChangeLog{
			LessonID:   lesson.ID,
			RuleID:     lesson.RuleID,
			ActionType: actionType,
			BeforeData: beforeData,
			AfterData:  afterData,
			Reason:     reason,
			OccurredAt: &models.JSONTime{Time: time.Now().UTC()},
			TenantID:   tenantID,
		}
		if err := tx.Create(changeLog).Error; err != nil {
			return err
		}
		if recheckConflict {
			if lesson.StudentID != 0 {
				if err := s.recheckLessonEligibilityTx(tx, ctx, &lesson); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *EduScheduleService) generateLessonsForRuleTx(ctx context.Context, rule *models.EduScheduleRule, operatorID uint) ([]models.EduLesson, error) {
	return s.generateLessonsForRuleTxWithDB(ctx, app.DB().WithContext(ctx), rule, operatorID)
}

func (s *EduScheduleService) generateLessonsForRuleTxWithDB(ctx context.Context, db *gorm.DB, rule *models.EduScheduleRule, operatorID uint) ([]models.EduLesson, error) {
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
	err = db.Transaction(func(tx *gorm.DB) error {
		created, err := s.generateLessonsForDatesTx(ctx, tx, rule, operatorID, lessonDates)
		if err != nil {
			return err
		}
		lessons = created
		return nil
	})
	return lessons, err
}

func (s *EduScheduleService) generateLessonsForDatesTx(ctx context.Context, tx *gorm.DB, rule *models.EduScheduleRule, operatorID uint, lessonDates []time.Time) ([]models.EduLesson, error) {
	var lessons []models.EduLesson
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
			TenantID:     rule.TenantID,
		}
		if err := tx.Create(&lesson).Error; err != nil {
			return nil, err
		}
		if err := s.writeEligibilityRowsTx(tx, ctx, &lesson, rule); err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
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

func lessonDateString(value *models.JSONTime) string {
	if value == nil || value.Time.IsZero() {
		return ""
	}
	return value.Time.Format("2006-01-02")
}

func appendConflict(result *models.EduScheduleConflictCheckResult, seen map[string]struct{}, conflictType, conflictKey string, matched bool) {
	if !matched {
		return
	}
	key := conflictType + ":" + conflictKey
	if _, ok := seen[key]; ok {
		return
	}
	seen[key] = struct{}{}
	result.HasConflict = true
	result.Items = append(result.Items, models.EduScheduleConflictItem{
		ConflictType: conflictType,
		ConflictKey:  conflictKey,
		Reason:       "时间段冲突",
	})
}

func (s *EduScheduleService) recheckLessonEligibilityTx(tx *gorm.DB, ctx context.Context, lesson *models.EduLesson) error {
	if tx == nil || lesson == nil {
		return nil
	}
	if lesson.LessonDate == nil || lesson.LessonDate.Time.IsZero() {
		return errors.New("课次日期不能为空")
	}
	var rows []models.EduLessonStudentEligibility
	if err := requireTenant(tx.Model(&models.EduLessonStudentEligibility{}), lesson.TenantID).Where("lesson_id = ?", lesson.ID).Find(&rows).Error; err != nil {
		return err
	}
	if len(rows) == 0 {
		if lesson.StudentID == 0 {
			return nil
		}
		check, err := checkEligibilityTx(tx, lesson.LessonDate, lesson.TenantID, lesson.StudentID, lesson.CourseID, lesson.ClassID, lesson.TeacherID)
		if err != nil {
			return err
		}
		if check == nil || !check.Eligible {
			return errors.New("学生权益不足")
		}
		return nil
	}
	for _, row := range rows {
		check, err := checkEligibilityTx(tx, lesson.LessonDate, lesson.TenantID, row.StudentID, lesson.CourseID, lesson.ClassID, lesson.TeacherID)
		if err != nil {
			return err
		}
		if check == nil || !check.Eligible {
			return errors.New("学生权益不足")
		}
		status := check.Status
		if status == "" {
			status = "ineligible"
		}
		updates := map[string]interface{}{
			"eligibility_status": status,
			"reason_code":        check.ReasonCode,
			"student_benefit_id": check.StudentBenefitID,
			"checked_at":         lesson.LessonDate,
		}
		if err := requireTenant(tx.Model(&models.EduLessonStudentEligibility{}), lesson.TenantID).
			Where("id = ?", row.ID).
			Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}
