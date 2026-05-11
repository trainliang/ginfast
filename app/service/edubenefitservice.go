package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gin-fast/app/global/app"
	"gin-fast/app/models"

	"gorm.io/gorm"
)

type EduBenefitService struct{}

func NewEduBenefitService() *EduBenefitService { return &EduBenefitService{} }

func (s *EduBenefitService) resolveTenantID(ctx context.Context, requested *uint) (uint, error) {
	if requested != nil && *requested != 0 {
		return *requested, nil
	}
	tenantID := tenantIDFromContext(ctx)
	if tenantID == 0 {
		return 0, errors.New("缺少租户信息")
	}
	return tenantID, nil
}

func (s *EduBenefitService) CheckEligibility(ctx context.Context, req *models.EduBenefitCheckRequest) (*models.EduBenefitCheckResult, error) {
	if req == nil {
		return nil, errors.New("请求不能为空")
	}
	tenantID, err := s.resolveTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	if req.StudentID == 0 || req.CourseID == 0 {
		return nil, errors.New("学生ID和课程ID不能为空")
	}
	checkedAt := time.Now().UTC()
	if req.CheckedAt != nil && !req.CheckedAt.Time.IsZero() {
		checkedAt = req.CheckedAt.Time
	}

	var benefits []models.EduStudentBenefit
	err = app.DB().WithContext(ctx).
		Where("tenant_id = ? AND student_id = ? AND course_id = ? AND status = 1 AND benefit_type = ?", tenantID, req.StudentID, req.CourseID, "course").
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
		if !studentBenefitMatchesDimension(&benefit, req.ClassID, req.TeacherID) {
			continue
		}
		if benefit.ValidFrom != nil && checkedAt.Before(benefit.ValidFrom.Time) {
			if expiredResult == nil {
				expiredResult = &models.EduBenefitCheckResult{Eligible: false, Status: "ineligible", ReasonCode: "expired", StudentBenefitID: benefit.ID}
			}
			continue
		}
		if benefit.ValidTo != nil && checkedAt.After(benefit.ValidTo.Time) {
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

func studentBenefitMatchesDimension(benefit *models.EduStudentBenefit, classID, teacherID uint) bool {
	if benefit == nil {
		return false
	}
	if classID == 0 {
		if benefit.ClassID != 0 {
			return false
		}
	} else if benefit.ClassID != 0 && benefit.ClassID != classID {
		return false
	}
	if teacherID == 0 {
		if benefit.TeacherID != 0 {
			return false
		}
	} else if benefit.TeacherID != 0 && benefit.TeacherID != teacherID {
		return false
	}
	return true
}

func (s *EduBenefitService) RepairScheduleEligibility(ctx context.Context, req *models.EduBenefitRepairRequest) (*models.EduBenefitRepairResult, error) {
	if req == nil {
		return nil, errors.New("请求不能为空")
	}
	tenantID, err := s.resolveTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	query := app.DB().WithContext(ctx).Model(&models.EduLessonStudentEligibility{}).Where("tenant_id = ?", tenantID).
		Where("eligibility_status IN ?", []string{"ineligible", "warn"}).
		Where("resolved_at IS NULL")
	if req.StudentID != nil {
		query = query.Where("student_id = ?", *req.StudentID)
	}
	if req.CourseID != nil {
		query = query.Where("course_id = ?", *req.CourseID)
	}
	effectiveFrom := time.Now().UTC()
	if req.EffectiveFrom != nil && !req.EffectiveFrom.Time.IsZero() {
		effectiveFrom = req.EffectiveFrom.Time
	}
	query = query.Where("checked_at >= ?", effectiveFrom)
	var rows []models.EduLessonStudentEligibility
	if err := query.Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := &models.EduBenefitRepairResult{}
	for i := range rows {
		row := &rows[i]
		if row.CheckedAt == nil || row.CheckedAt.Time.IsZero() {
			continue
		}
		result.CheckedCount++
		check, err := s.CheckEligibility(ctx, &models.EduBenefitCheckRequest{
			TenantID:  &tenantID,
			StudentID: row.StudentID,
			CourseID:  row.CourseID,
			ClassID:   row.ClassID,
			TeacherID: 0,
			CheckedAt: row.CheckedAt,
		})
		if err != nil {
			return nil, err
		}
		if check.Eligible {
			now := models.NewJSONTime(time.Now().UTC())
			updateResult := app.DB().WithContext(ctx).Model(&models.EduLessonStudentEligibility{}).
				Where("id = ? AND tenant_id = ? AND resolved_at IS NULL", row.ID, tenantID).
				Updates(map[string]interface{}{
					"eligibility_status": "eligible",
					"reason_code":        "",
					"student_benefit_id": check.StudentBenefitID,
					"resolved_at":        &now,
				})
			if updateResult.Error != nil {
				return nil, updateResult.Error
			}
			if updateResult.RowsAffected > 0 {
				result.RepairedCount++
			}
		}
	}
	return result, nil
}

func (s *EduBenefitService) CreateProduct(ctx context.Context, product *models.EduBenefitProduct) error {
	if product == nil {
		return errors.New("权益产品不能为空")
	}
	if err := ensureTenantID(product.TenantID); err != nil {
		return err
	}
	if product.BenefitType == "" {
		product.BenefitType = "course"
	}
	if product.CalculationMode == "" {
		product.CalculationMode = "count_limited"
	}
	if product.Status == 0 {
		product.Status = 1
	}
	var count int64
	if err := app.DB().WithContext(ctx).Model(&models.EduBenefitProduct{}).
		Where("tenant_id = ? AND code = ?", product.TenantID, product.Code).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("权益产品编码已存在")
	}
	return app.DB().WithContext(ctx).Create(product).Error
}

func (s *EduBenefitService) UpdateProduct(ctx context.Context, product *models.EduBenefitProduct) error {
	if product == nil {
		return errors.New("权益产品不能为空")
	}
	if err := ensureTenantID(product.TenantID); err != nil {
		return err
	}
	var existing models.EduBenefitProduct
	if err := app.DB().WithContext(ctx).Where("id = ? AND tenant_id = ?", product.ID, product.TenantID).First(&existing).Error; err != nil {
		return err
	}
	var count int64
	if err := app.DB().WithContext(ctx).Model(&models.EduBenefitProduct{}).
		Where("tenant_id = ? AND code = ? AND id <> ?", product.TenantID, product.Code, product.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("权益产品编码已存在")
	}
	updates := map[string]interface{}{
		"name":             product.Name,
		"code":             product.Code,
		"benefit_type":     product.BenefitType,
		"calculation_mode": product.CalculationMode,
		"total_count":      product.TotalCount,
		"valid_days":       product.ValidDays,
		"status":           product.Status,
		"remark":           product.Remark,
	}
	return app.DB().WithContext(ctx).Model(&models.EduBenefitProduct{}).
		Where("id = ? AND tenant_id = ?", product.ID, product.TenantID).
		Updates(updates).Error
}

func (s *EduBenefitService) DeleteProduct(ctx context.Context, tenantID, id uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	var count int64
	if err := app.DB().WithContext(ctx).Model(&models.EduStudentBenefit{}).
		Where("tenant_id = ? AND product_id = ?", tenantID, id).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("权益产品已被学生权益引用，禁止删除")
	}
	result := app.DB().WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.EduBenefitProduct{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("权益产品不存在或不属于当前租户")
	}
	return nil
}

func (s *EduBenefitService) CreateStudentBenefit(ctx context.Context, benefit *models.EduStudentBenefit) error {
	if benefit == nil {
		return errors.New("学生权益不能为空")
	}
	if err := ensureTenantID(benefit.TenantID); err != nil {
		return err
	}
	if benefit.BenefitType == "" {
		benefit.BenefitType = "course"
	}
	if benefit.CalculationMode == "" {
		benefit.CalculationMode = "count_limited"
	}
	if benefit.Status == 0 {
		benefit.Status = 1
	}
	if benefit.CalculationMode == "count_limited" && benefit.RemainingCount == 0 && benefit.TotalCount > benefit.UsedCount {
		benefit.RemainingCount = benefit.TotalCount - benefit.UsedCount
	}
	return app.DB().WithContext(ctx).Create(benefit).Error
}

func (s *EduBenefitService) UpdateStudentBenefit(ctx context.Context, benefit *models.EduStudentBenefit) error {
	if benefit == nil {
		return errors.New("学生权益不能为空")
	}
	if err := ensureTenantID(benefit.TenantID); err != nil {
		return err
	}
	var existing models.EduStudentBenefit
	if err := app.DB().WithContext(ctx).Where("id = ? AND tenant_id = ?", benefit.ID, benefit.TenantID).First(&existing).Error; err != nil {
		return err
	}
	if benefit.BenefitType == "" {
		benefit.BenefitType = "course"
	}
	if benefit.CalculationMode == "" {
		benefit.CalculationMode = "count_limited"
	}
	if benefit.Status == 0 {
		benefit.Status = 1
	}
	if benefit.CalculationMode == "count_limited" && benefit.RemainingCount == 0 && benefit.TotalCount > benefit.UsedCount {
		benefit.RemainingCount = benefit.TotalCount - benefit.UsedCount
	}
	updates := map[string]interface{}{
		"student_id":       benefit.StudentID,
		"product_id":       benefit.ProductID,
		"benefit_type":     benefit.BenefitType,
		"calculation_mode": benefit.CalculationMode,
		"course_id":        benefit.CourseID,
		"class_id":         benefit.ClassID,
		"teacher_id":       benefit.TeacherID,
		"valid_from":       benefit.ValidFrom,
		"valid_to":         benefit.ValidTo,
		"total_count":      benefit.TotalCount,
		"used_count":       benefit.UsedCount,
		"remaining_count":  benefit.RemainingCount,
		"status":           benefit.Status,
		"source_type":      benefit.SourceType,
	}
	return app.DB().WithContext(ctx).Model(&models.EduStudentBenefit{}).
		Where("id = ? AND tenant_id = ?", benefit.ID, benefit.TenantID).
		Updates(updates).Error
}

func (s *EduBenefitService) DeleteStudentBenefit(ctx context.Context, tenantID, id uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	result := app.DB().WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.EduStudentBenefit{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("学生权益不存在或不属于当前租户")
	}
	return nil
}

func (s *EduBenefitService) RetryExternalSync(ctx context.Context, tenantID, id uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	var row models.EduBenefitExternalSync
	if err := app.DB().WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).First(&row).Error; err != nil {
		return err
	}
	switch row.Status {
	case "failed", "retrying", "pending":
		return app.DB().WithContext(ctx).Model(&models.EduBenefitExternalSync{}).
			Where("id = ? AND tenant_id = ?", id, tenantID).
			Updates(map[string]interface{}{
				"status":        "retrying",
				"retry_count":   gorm.Expr("retry_count + 1"),
				"last_error":    row.LastError,
				"next_retry_at": row.NextRetryAt,
			}).Error
	case "success", "canceled":
		return errors.New("当前状态不允许重试")
	default:
		return errors.New("当前状态不允许重试")
	}
}
