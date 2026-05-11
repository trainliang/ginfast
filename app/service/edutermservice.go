package service

import (
	"context"
	"errors"
	"fmt"

	"gin-fast/app/global/app"
	"gin-fast/app/models"

	"gorm.io/gorm"
)

type EduTermService struct{}

func NewEduTermService() *EduTermService { return &EduTermService{} }

func (s *EduTermService) validateTerm(term *models.EduTerm) error {
	if term == nil {
		return errors.New("学期不能为空")
	}
	if err := ensureTenantID(term.TenantID); err != nil {
		return err
	}
	if term.StartDate == nil || term.EndDate == nil {
		return errors.New("开始日期和结束日期不能为空")
	}
	if term.EndDate.Time.Before(term.StartDate.Time) {
		return errors.New("结束日期不能早于开始日期")
	}
	return nil
}

func (s *EduTermService) Create(ctx context.Context, term *models.EduTerm) error {
	if err := s.validateTerm(term); err != nil {
		return err
	}
	if term.Status == 0 {
		term.Status = 1
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduTerm{}), term.TenantID).Where("name = ?", term.Name).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("学期名称已存在")
		}
		return tx.Create(term).Error
	})
}

func (s *EduTermService) Update(ctx context.Context, term *models.EduTerm) error {
	if err := s.validateTerm(term); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduTerm{}), term.TenantID).Where("id = ?", term.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("学期不存在或不属于当前租户")
		}
		if err := requireTenant(tx.Model(&models.EduTerm{}), term.TenantID).Where("name = ? AND id <> ?", term.Name, term.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("学期名称已存在")
		}
		result := requireTenant(tx.Model(&models.EduTerm{}), term.TenantID).
			Where("id = ?", term.ID).
			Updates(map[string]interface{}{
				"name":       term.Name,
				"start_date": term.StartDate,
				"end_date":   term.EndDate,
				"status":     term.Status,
			})
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func (s *EduTermService) Delete(ctx context.Context, tenantID, id uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduTerm{}), tenantID).Where("id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("学期不存在或不属于当前租户")
		}
		if err := requireTenant(tx.Model(&models.EduTermClosedDay{}), tenantID).Where("term_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("学期已被停课日引用，禁止删除")
		}
		result := requireTenant(tx.Where("id = ?", id), tenantID).Delete(&models.EduTerm{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("学期不存在或不属于当前租户")
		}
		return nil
	})
}

func (s *EduTermService) ClosedDays(ctx context.Context, tenantID, termID uint) ([]models.EduTermClosedDay, error) {
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	var term models.EduTerm
	if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduTerm{}), tenantID).Where("id = ?", termID).First(&term).Error; err != nil {
		return nil, err
	}
	var rows []models.EduTermClosedDay
	if err := requireTenant(app.DB().WithContext(ctx).Model(&models.EduTermClosedDay{}), tenantID).Where("term_id = ?", termID).Order("closed_date asc, id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *EduTermService) SaveClosedDays(ctx context.Context, tenantID, termID uint, rows []models.EduTermClosedDay) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var term models.EduTerm
		if err := requireTenant(tx.Model(&models.EduTerm{}), tenantID).Where("id = ?", termID).First(&term).Error; err != nil {
			return err
		}
		if term.IsEmpty() {
			return errors.New("学期不存在或不属于当前租户")
		}
		if err := requireTenant(tx.Where("term_id = ?", termID), tenantID).Delete(&models.EduTermClosedDay{}).Error; err != nil {
			return err
		}
		for i := range rows {
			if rows[i].ClosedDate == nil || rows[i].ClosedDate.Time.IsZero() {
				return errors.New("停课日期不能为空")
			}
			rows[i].ID = 0
			rows[i].TermID = termID
			rows[i].TenantID = tenantID
			if err := tx.Create(&rows[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
