package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gin-fast/app/global/app"
	"gin-fast/app/models"

	"gorm.io/gorm"
)

type EduCourseService struct{}

func NewEduCourseService() *EduCourseService { return &EduCourseService{} }

func (s *EduCourseService) Create(ctx context.Context, course *models.EduCourse) error {
	if course == nil {
		return errors.New("课程不能为空")
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduCourse{}), course.TenantID).Where("code = ?", course.Code).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("课程编码已存在")
		}
		return tx.Create(course).Error
	})
}

func (s *EduCourseService) Update(ctx context.Context, course *models.EduCourse) error {
	if course == nil {
		return errors.New("课程不能为空")
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduCourse{}), course.TenantID).Where("code = ? AND id <> ?", course.Code, course.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("课程编码已存在")
		}
		return tx.Save(course).Error
	})
}

func (s *EduCourseService) Delete(ctx context.Context, tenantID, id uint) error {
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduClass{}), tenantID).Where("course_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("课程被班级引用，禁止删除")
		}
		return requireTenant(tx.Where("id = ?", id), tenantID).Delete(&models.EduCourse{}).Error
	})
}

func (s *EduCourseService) ImportRows(ctx context.Context, tenantID uint, rows []models.EduCourseImportRow) (*EduImportResult, error) {
	result := &EduImportResult{Success: false}
	if len(rows) == 0 {
		result.Success = true
		return result, nil
	}
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		seen := map[string]struct{}{}
		for i, row := range rows {
			rowNum := i + 1
			code := strings.TrimSpace(row.Code)
			if code == "" {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "code", Reason: "不能为空"})
				continue
			}
			if _, ok := seen[code]; ok {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "code", Reason: "重复"})
				continue
			}
			seen[code] = struct{}{}
		}
		if len(result.Errors) > 0 {
			return errors.New("导入数据存在错误")
		}
		for _, row := range rows {
			course := models.EduCourse{}
			err := requireTenant(tx.Model(&models.EduCourse{}), tenantID).Where("code = ?", strings.TrimSpace(row.Code)).First(&course).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				course = models.EduCourse{
					Name:        row.Name,
					Code:        strings.TrimSpace(row.Code),
					Type:        row.Type,
					GradeRange:  row.GradeRange,
					Status:      1,
					Sort:        0,
					Description: row.Description,
					TenantID:    tenantID,
				}
				if row.Status != nil {
					course.Status = *row.Status
				}
				if row.Sort != nil {
					course.Sort = *row.Sort
				}
				if err := tx.Create(&course).Error; err != nil {
					return err
				}
				result.Created++
				continue
			}
			course.Name = row.Name
			course.Type = row.Type
			course.GradeRange = row.GradeRange
			course.Description = row.Description
			if row.Status != nil {
				course.Status = *row.Status
			}
			if row.Sort != nil {
				course.Sort = *row.Sort
			}
			if err := tx.Save(&course).Error; err != nil {
				return err
			}
			result.Updated++
		}
		result.Success = true
		return nil
	})
	if err != nil && len(result.Errors) == 0 {
		return result, err
	}
	return result, err
}

func (s *EduCourseService) ExportRows(ctx context.Context, tenantID uint, ids []uint) ([]models.EduCourseExportRow, error) {
	var courses []models.EduCourse
	db := requireTenant(app.DB().WithContext(ctx).Model(&models.EduCourse{}), tenantID)
	if len(ids) > 0 {
		db = db.Where("id IN ?", ids)
	}
	if err := db.Order("id asc").Find(&courses).Error; err != nil {
		return nil, err
	}
	rows := make([]models.EduCourseExportRow, 0, len(courses))
	for _, course := range courses {
		rows = append(rows, models.EduCourseExportRow{
			ID:          course.ID,
			Name:        course.Name,
			Code:        course.Code,
			Type:        course.Type,
			GradeRange:  course.GradeRange,
			Status:      course.Status,
			Sort:        course.Sort,
			Description: course.Description,
			CreatedAt:   course.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   course.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return rows, nil
}
