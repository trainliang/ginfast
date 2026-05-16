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

type EduLessonCompletionService struct{}

func NewEduLessonCompletionService() *EduLessonCompletionService {
	return &EduLessonCompletionService{}
}

func validCompletionResultType(value string) bool {
	for _, item := range models.EduLessonCompletionResultTypes {
		if value == item {
			return true
		}
	}
	return false
}

func validateCompletionRuleItem(item models.EduLessonCompletionRuleItemRequest) error {
	if !validCompletionResultType(strings.TrimSpace(item.ResultType)) {
		return fmt.Errorf("无效结课结果: %s", item.ResultType)
	}
	if item.DeductEnabled != 0 && item.DeductEnabled != 1 {
		return errors.New("是否扣减只能为0或1")
	}
	if item.DeductMode == "" {
		item.DeductMode = models.EduLessonCompletionDeductModeFixedCount
	}
	if item.DeductMode != models.EduLessonCompletionDeductModeFixedCount && item.DeductMode != models.EduLessonCompletionDeductModeDuration {
		return errors.New("扣减方式无效")
	}
	if item.DeductEnabled == 1 && item.DeductMode == models.EduLessonCompletionDeductModeFixedCount && item.FixedCount <= 0 {
		return errors.New("固定扣减次数必须大于0")
	}
	if item.DeductEnabled == 1 && item.DeductMode == models.EduLessonCompletionDeductModeDuration && item.DurationUnitMinutes <= 0 {
		return errors.New("课时时长换算单位必须大于0")
	}
	if item.InsufficientPolicy == "" {
		item.InsufficientPolicy = models.EduLessonCompletionInsufficientBlock
	}
	if item.InsufficientPolicy != models.EduLessonCompletionInsufficientBlock && item.InsufficientPolicy != models.EduLessonCompletionInsufficientAllowArrears {
		return errors.New("权益不足策略无效")
	}
	return nil
}

func validateCompletionRuleSet(items []models.EduLessonCompletionRuleItemRequest) error {
	if len(items) != len(models.EduLessonCompletionResultTypes) {
		return errors.New("必须配置全部结课结果")
	}
	seen := map[string]bool{}
	for _, item := range items {
		resultType := strings.TrimSpace(item.ResultType)
		if seen[resultType] {
			return fmt.Errorf("结课结果重复: %s", resultType)
		}
		if err := validateCompletionRuleItem(item); err != nil {
			return err
		}
		seen[resultType] = true
	}
	for _, resultType := range models.EduLessonCompletionResultTypes {
		if !seen[resultType] {
			return fmt.Errorf("缺少结课结果配置: %s", resultType)
		}
	}
	return nil
}

func defaultLessonCompletionRules() []models.EduLessonCompletionRule {
	items := make([]models.EduLessonCompletionRule, 0, len(models.EduLessonCompletionResultTypes))
	for _, resultType := range models.EduLessonCompletionResultTypes {
		rule := models.EduLessonCompletionRule{
			ResultType:          resultType,
			DeductEnabled:       0,
			DeductMode:          models.EduLessonCompletionDeductModeFixedCount,
			FixedCount:          1,
			DurationUnitMinutes: 60,
			InsufficientPolicy:  models.EduLessonCompletionInsufficientBlock,
			Status:              1,
		}
		if resultType == models.EduLessonCompletionResultAttended {
			rule.DeductEnabled = 1
		}
		items = append(items, rule)
	}
	return items
}

func mergeLessonCompletionRules(overrides []models.EduLessonCompletionRule, tenantID uint) []models.EduLessonCompletionRule {
	defaults := defaultLessonCompletionRules()
	byResultType := make(map[string]models.EduLessonCompletionRule, len(overrides))
	for _, item := range overrides {
		byResultType[item.ResultType] = item
	}
	for i := range defaults {
		defaults[i].TenantID = tenantID
		if override, ok := byResultType[defaults[i].ResultType]; ok {
			defaults[i] = override
		}
		if defaults[i].DeductMode == "" {
			defaults[i].DeductMode = models.EduLessonCompletionDeductModeFixedCount
		}
		if defaults[i].DurationUnitMinutes <= 0 {
			defaults[i].DurationUnitMinutes = 60
		}
		if defaults[i].FixedCount <= 0 {
			defaults[i].FixedCount = 1
		}
		if defaults[i].InsufficientPolicy == "" {
			defaults[i].InsufficientPolicy = models.EduLessonCompletionInsufficientBlock
		}
		if defaults[i].Status == 0 {
			defaults[i].Status = 1
		}
	}
	return defaults
}

func (s *EduLessonCompletionService) SaveRules(ctx context.Context, tenantID uint, req *models.EduLessonCompletionRuleSaveRequest) error {
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	if req == nil {
		return errors.New("消课规则不能为空")
	}
	if err := validateCompletionRuleSet(req.Items); err != nil {
		return err
	}
	operatorID := userIDFromContext(ctx)
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Items {
			deductMode := item.DeductMode
			if deductMode == "" {
				deductMode = models.EduLessonCompletionDeductModeFixedCount
			}
			insufficientPolicy := item.InsufficientPolicy
			if insufficientPolicy == "" {
				insufficientPolicy = models.EduLessonCompletionInsufficientBlock
			}
			status := item.Status
			if status == 0 {
				status = 1
			}
			updates := map[string]interface{}{
				"deduct_enabled":        item.DeductEnabled,
				"deduct_mode":           deductMode,
				"fixed_count":           item.FixedCount,
				"duration_unit_minutes": item.DurationUnitMinutes,
				"insufficient_policy":   insufficientPolicy,
				"status":                status,
				"remark":                item.Remark,
			}
			var existing models.EduLessonCompletionRule
			err := requireTenant(tx.Model(&models.EduLessonCompletionRule{}), tenantID).
				Where("result_type = ?", item.ResultType).
				First(&existing).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rule := &models.EduLessonCompletionRule{
					ResultType:          item.ResultType,
					DeductEnabled:       item.DeductEnabled,
					DeductMode:          deductMode,
					FixedCount:          item.FixedCount,
					DurationUnitMinutes: item.DurationUnitMinutes,
					InsufficientPolicy:  insufficientPolicy,
					Status:              status,
					Remark:              item.Remark,
					CreatedBy:           operatorID,
					TenantID:            tenantID,
				}
				if err := tx.Create(rule).Error; err != nil {
					return err
				}
				continue
			}
			if err != nil {
				return err
			}
			if err := requireTenant(tx.Model(&models.EduLessonCompletionRule{}), tenantID).
				Where("id = ?", existing.ID).
				Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *EduLessonCompletionService) ListRules(ctx context.Context, tenantID uint) ([]models.EduLessonCompletionRule, error) {
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
	}
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	var rows []models.EduLessonCompletionRule
	if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduLessonCompletionRule{}), tenantID).
		Order("result_type asc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	rows = mergeLessonCompletionRules(rows, tenantID)
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].ResultType < rows[j].ResultType })
	return rows, nil
}

func (s *EduLessonCompletionService) GetLessonCompletionStatus(ctx context.Context, tenantID, lessonID uint) (*models.EduLessonCompletionStatusResponse, error) {
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
	}
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	var lesson models.EduLesson
	if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduLesson{}), tenantID).
		Where("id = ?", lessonID).
		First(&lesson).Error; err != nil {
		return nil, err
	}
	studentIDs, err := s.lessonStudentIDs(ctx, app.DB(), tenantID, &lesson)
	if err != nil {
		return nil, err
	}
	var completions []models.EduLessonStudentCompletion
	if len(studentIDs) > 0 {
		if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduLessonStudentCompletion{}), tenantID).
			Where("lesson_id = ? AND student_id IN ? AND status IN ?", lesson.ID, studentIDs, []string{models.EduLessonCompletionStatusCompleted, models.EduLessonCompletionStatusArrears}).
			Find(&completions).Error; err != nil {
			return nil, err
		}
	}
	byStudent := map[uint]*models.EduLessonStudentCompletion{}
	for i := range completions {
		completion := completions[i]
		byStudent[completion.StudentID] = &completion
	}
	resp := &models.EduLessonCompletionStatusResponse{Lesson: &lesson}
	for _, studentID := range studentIDs {
		item := models.EduLessonCompletionStudentStatus{StudentID: studentID}
		if completion, ok := byStudent[studentID]; ok {
			item.Completion = completion
			item.Processed = true
		} else {
			resp.UnprocessedCount++
		}
		resp.Students = append(resp.Students, item)
	}
	return resp, nil
}

func (s *EduLessonCompletionService) lessonStudentIDs(ctx context.Context, db *gorm.DB, tenantID uint, lesson *models.EduLesson) ([]uint, error) {
	if lesson == nil || lesson.ID == 0 {
		return nil, errors.New("课次不能为空")
	}
	if lesson.StudentID != 0 {
		return []uint{lesson.StudentID}, nil
	}
	if lesson.ClassID == 0 {
		return nil, errors.New("课次缺少学生或班级信息")
	}
	var rows []models.EduClassMember
	if err := requireTenant(db.WithContext(ctx).Model(&models.EduClassMember{}), tenantID).
		Where("class_id = ? AND status = ?", lesson.ClassID, "studying").
		Order("student_id asc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.StudentID)
	}
	return ids, nil
}

func (s *EduLessonCompletionService) calculateDeductCount(lesson *models.EduLesson, rule *models.EduLessonCompletionRule) (int, error) {
	if rule == nil || rule.DeductEnabled == 0 {
		return 0, nil
	}
	switch rule.DeductMode {
	case "", models.EduLessonCompletionDeductModeFixedCount:
		if rule.FixedCount <= 0 {
			return 0, errors.New("固定扣减次数必须大于0")
		}
		return rule.FixedCount, nil
	case models.EduLessonCompletionDeductModeDuration:
		if rule.DurationUnitMinutes <= 0 {
			return 0, errors.New("课时时长换算单位必须大于0")
		}
		start, err := time.Parse("15:04", lesson.StartTime)
		if err != nil {
			return 0, errors.New("课次开始时间无效")
		}
		end, err := time.Parse("15:04", lesson.EndTime)
		if err != nil {
			return 0, errors.New("课次结束时间无效")
		}
		minutes := int(end.Sub(start).Minutes())
		if minutes <= 0 {
			return 0, errors.New("课次时长无效")
		}
		count := minutes / rule.DurationUnitMinutes
		if minutes%rule.DurationUnitMinutes != 0 {
			count++
		}
		if count <= 0 {
			count = 1
		}
		return count, nil
	default:
		return 0, errors.New("扣减方式无效")
	}
}

func (s *EduLessonCompletionService) loadRuleMap(ctx context.Context, tx *gorm.DB, tenantID uint) (map[string]models.EduLessonCompletionRule, error) {
	var rows []models.EduLessonCompletionRule
	if err := requireTenant(tx.WithContext(ctx).Model(&models.EduLessonCompletionRule{}), tenantID).
		Where("status = ?", 1).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	rows = mergeLessonCompletionRules(rows, tenantID)
	result := map[string]models.EduLessonCompletionRule{}
	for _, row := range rows {
		if row.Status != 1 {
			continue
		}
		result[row.ResultType] = row
	}
	for _, resultType := range models.EduLessonCompletionResultTypes {
		if _, ok := result[resultType]; !ok {
			return nil, errors.New("消课规则不完整")
		}
	}
	return result, nil
}

func (s *EduLessonCompletionService) nextCompletionVersion(ctx context.Context, tx *gorm.DB, tenantID, lessonID, studentID uint) int {
	var count int64
	_ = requireTenant(tx.WithContext(ctx).Model(&models.EduLessonStudentCompletion{}), tenantID).
		Where("lesson_id = ? AND student_id = ?", lessonID, studentID).
		Count(&count).Error
	return int(count) + 1
}

func (s *EduLessonCompletionService) findDeductibleBenefit(ctx context.Context, tx *gorm.DB, tenantID, studentID, courseID, classID, teacherID uint, deductCount int) (*models.EduStudentBenefit, error) {
	var benefits []models.EduStudentBenefit
	if err := requireTenant(tx.WithContext(ctx).Model(&models.EduStudentBenefit{}), tenantID).
		Where("student_id = ? AND status = 1 AND benefit_type = ? AND calculation_mode = ?", studentID, "course", "count_limited").
		Order("id asc").
		Find(&benefits).Error; err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	for i := range benefits {
		benefit := &benefits[i]
		if benefit.RemainingCount < deductCount {
			continue
		}
		if benefit.ValidFrom != nil && now.Before(benefit.ValidFrom.Time) {
			continue
		}
		if benefit.ValidTo != nil && now.After(benefit.ValidTo.Time) {
			continue
		}
		if benefit.CourseID != courseID {
			continue
		}
		if benefit.ClassID != 0 && benefit.ClassID != classID {
			continue
		}
		if benefit.TeacherID != 0 && benefit.TeacherID != teacherID {
			continue
		}
		return benefit, nil
	}
	return nil, errors.New("权益不足")
}

func (s *EduLessonCompletionService) applyLessonDeduction(ctx context.Context, tx *gorm.DB, tenantID uint, completion *models.EduLessonStudentCompletion, benefit *models.EduStudentBenefit, deductCount int, operatorID uint) (uint, error) {
	if benefit == nil || completion == nil {
		return 0, errors.New("权益或消课记录不能为空")
	}
	before := benefit.RemainingCount
	after := before - deductCount
	if after < 0 {
		return 0, errors.New("权益不足")
	}
	benefit.UsedCount += deductCount
	benefit.RemainingCount = after
	if err := requireTenant(tx.Model(&models.EduStudentBenefit{}), tenantID).
		Where("id = ?", benefit.ID).
		Updates(map[string]interface{}{"used_count": benefit.UsedCount, "remaining_count": benefit.RemainingCount}).Error; err != nil {
		return 0, err
	}
	now := models.NewJSONTime(time.Now().UTC())
	ledger := &models.EduBenefitLedger{
		StudentBenefitID: benefit.ID,
		StudentID:        benefit.StudentID,
		ActionType:       "lesson_complete",
		BizType:          "lesson_completion",
		BizID:            completion.LessonID,
		ChangeCount:      -deductCount,
		BeforeCount:      before,
		AfterCount:       after,
		OccurredAt:       &now,
		OperatorID:       operatorID,
		TenantID:         tenantID,
	}
	if err := tx.Create(ledger).Error; err != nil {
		return 0, err
	}
	return ledger.ID, nil
}

func (s *EduLessonCompletionService) SubmitCompletions(ctx context.Context, tenantID uint, req *models.EduLessonCompletionSubmitRequest) ([]models.EduLessonStudentCompletion, error) {
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
	}
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	if req == nil || len(req.Items) == 0 {
		return nil, errors.New("消课学生不能为空")
	}
	operatorID := userIDFromContext(ctx)
	created := make([]models.EduLessonStudentCompletion, 0, len(req.Items))
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lesson models.EduLesson
		if err := requireTenant(tx.Model(&models.EduLesson{}), tenantID).Where("id = ?", req.GetLessonIDUint()).First(&lesson).Error; err != nil {
			return err
		}
		if lesson.Status != "scheduled" {
			return errors.New("只有待上课课次允许消课")
		}
		studentIDs, err := s.lessonStudentIDs(ctx, tx, tenantID, &lesson)
		if err != nil {
			return err
		}
		allowed := map[uint]bool{}
		for _, id := range studentIDs {
			allowed[id] = true
		}
		rules, err := s.loadRuleMap(ctx, tx, tenantID)
		if err != nil {
			return err
		}
		for _, item := range req.Items {
			studentID := item.GetStudentIDUint()
			if !allowed[studentID] {
				return errors.New("学生不属于该课次")
			}
			if !validCompletionResultType(item.ResultType) {
				return errors.New("结课结果无效")
			}
			var activeCount int64
			if err := requireTenant(tx.Model(&models.EduLessonStudentCompletion{}), tenantID).
				Where("lesson_id = ? AND student_id = ? AND status IN ?", lesson.ID, studentID, []string{models.EduLessonCompletionStatusCompleted, models.EduLessonCompletionStatusArrears}).
				Count(&activeCount).Error; err != nil {
				return err
			}
			if activeCount > 0 {
				return errors.New("学生已消课，不能重复提交")
			}
			rule := rules[item.ResultType]
			deductCount, err := s.calculateDeductCount(&lesson, &rule)
			if err != nil {
				return err
			}
			completedAt := models.NewJSONTime(time.Now().UTC())
			completion := models.EduLessonStudentCompletion{
				LessonID:      lesson.ID,
				StudentID:     studentID,
				ClassID:       lesson.ClassID,
				CourseID:      lesson.CourseID,
				TeacherID:     lesson.TeacherID,
				ResultType:    item.ResultType,
				DeductEnabled: rule.DeductEnabled,
				DeductCount:   deductCount,
				Status:        models.EduLessonCompletionStatusCompleted,
				Reason:        item.Reason,
				OperatorID:    operatorID,
				CompletedAt:   &completedAt,
				Version:       s.nextCompletionVersion(ctx, tx, tenantID, lesson.ID, studentID),
				CreatedBy:     operatorID,
				TenantID:      tenantID,
			}
			var ledgerID uint
			if rule.DeductEnabled == 1 && deductCount > 0 {
				benefit, err := s.findDeductibleBenefit(ctx, tx, tenantID, studentID, lesson.CourseID, lesson.ClassID, lesson.TeacherID, deductCount)
				if err != nil {
					if rule.InsufficientPolicy != models.EduLessonCompletionInsufficientAllowArrears {
						return err
					}
					completion.Status = models.EduLessonCompletionStatusArrears
					completion.StudentBenefitID = 0
				} else {
					ledgerID, err = s.applyLessonDeduction(ctx, tx, tenantID, &completion, benefit, deductCount, operatorID)
					if err != nil {
						return err
					}
					completion.StudentBenefitID = benefit.ID
					completion.LedgerID = ledgerID
				}
			}
			if err := tx.Create(&completion).Error; err != nil {
				return err
			}
			if ledgerID != 0 {
				if err := tx.Model(&models.EduBenefitLedger{}).Where("id = ?", ledgerID).Update("biz_id", completion.ID).Error; err != nil {
					return err
				}
			}
			created = append(created, completion)
		}
		return nil
	})
	return created, err
}

func (s *EduLessonCompletionService) RevokeCompletion(ctx context.Context, tenantID uint, req *models.EduLessonCompletionRevokeRequest) error {
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	if req == nil {
		return errors.New("撤销请求不能为空")
	}
	operatorID := userIDFromContext(ctx)
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var completion models.EduLessonStudentCompletion
		if err := requireTenant(tx.Model(&models.EduLessonStudentCompletion{}), tenantID).
			Where("lesson_id = ? AND student_id = ? AND status IN ?", req.GetLessonIDUint(), req.GetStudentIDUint(), []string{models.EduLessonCompletionStatusCompleted, models.EduLessonCompletionStatusArrears}).
			First(&completion).Error; err != nil {
			return err
		}
		var revokeLedgerID uint
		if completion.Status == models.EduLessonCompletionStatusCompleted && completion.LedgerID != 0 && completion.StudentBenefitID != 0 && completion.DeductCount > 0 {
			var benefit models.EduStudentBenefit
			if err := requireTenant(tx.Model(&models.EduStudentBenefit{}), tenantID).Where("id = ?", completion.StudentBenefitID).First(&benefit).Error; err != nil {
				return err
			}
			before := benefit.RemainingCount
			after := before + completion.DeductCount
			used := benefit.UsedCount - completion.DeductCount
			if used < 0 {
				used = 0
			}
			if err := requireTenant(tx.Model(&models.EduStudentBenefit{}), tenantID).
				Where("id = ?", benefit.ID).
				Updates(map[string]interface{}{"used_count": used, "remaining_count": after}).Error; err != nil {
				return err
			}
			occurredAt := models.NewJSONTime(time.Now().UTC())
			ledger := &models.EduBenefitLedger{
				StudentBenefitID: benefit.ID,
				StudentID:        completion.StudentID,
				ActionType:       "lesson_revoke",
				BizType:          "lesson_completion",
				BizID:            completion.ID,
				ChangeCount:      completion.DeductCount,
				BeforeCount:      before,
				AfterCount:       after,
				OccurredAt:       &occurredAt,
				OperatorID:       operatorID,
				Remark:           req.Reason,
				TenantID:         tenantID,
			}
			if err := tx.Create(ledger).Error; err != nil {
				return err
			}
			revokeLedgerID = ledger.ID
		}
		revokedAt := models.NewJSONTime(time.Now().UTC())
		return requireTenant(tx.Model(&models.EduLessonStudentCompletion{}), tenantID).
			Where("id = ?", completion.ID).
			Updates(map[string]interface{}{
				"status":           models.EduLessonCompletionStatusRevoked,
				"revoke_ledger_id": revokeLedgerID,
				"revoked_at":       &revokedAt,
				"revoked_by":       operatorID,
				"revoke_reason":    req.Reason,
			}).Error
	})
}

func (s *EduLessonCompletionService) CompleteLesson(ctx context.Context, tenantID uint, req *models.EduLessonCompleteRequest) error {
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	if req == nil {
		return errors.New("结课请求不能为空")
	}
	operatorID := userIDFromContext(ctx)
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lesson models.EduLesson
		if err := requireTenant(tx.Model(&models.EduLesson{}), tenantID).Where("id = ?", req.GetLessonIDUint()).First(&lesson).Error; err != nil {
			return err
		}
		if lesson.Status != "scheduled" {
			return errors.New("只有待上课课次允许结课完成")
		}
		studentIDs, err := s.lessonStudentIDs(ctx, tx, tenantID, &lesson)
		if err != nil {
			return err
		}
		if len(studentIDs) == 0 {
			return errors.New("课次没有应处理学生")
		}
		var count int64
		if err := requireTenant(tx.Model(&models.EduLessonStudentCompletion{}), tenantID).
			Where("lesson_id = ? AND student_id IN ? AND status IN ?", lesson.ID, studentIDs, []string{models.EduLessonCompletionStatusCompleted, models.EduLessonCompletionStatusArrears}).
			Count(&count).Error; err != nil {
			return err
		}
		if int(count) != len(studentIDs) {
			return errors.New("课次还有未处理学生")
		}
		beforeData := fmt.Sprintf("status=%s", lesson.Status)
		if err := requireTenant(tx.Model(&models.EduLesson{}), tenantID).Where("id = ?", lesson.ID).Update("status", "completed").Error; err != nil {
			return err
		}
		occurredAt := models.NewJSONTime(time.Now().UTC())
		changeLog := &models.EduLessonChangeLog{
			LessonID:   lesson.ID,
			RuleID:     lesson.RuleID,
			ActionType: "complete",
			BeforeData: beforeData,
			AfterData:  "status=completed",
			Reason:     req.Reason,
			OperatorID: operatorID,
			OccurredAt: &occurredAt,
			TenantID:   tenantID,
		}
		return tx.Create(changeLog).Error
	})
}
