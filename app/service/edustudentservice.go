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

type EduStudentService struct{}

func NewEduStudentService() *EduStudentService { return &EduStudentService{} }

func (s *EduStudentService) Create(ctx context.Context, student *models.EduStudent) error {
	if student == nil {
		return errors.New("学生不能为空")
	}
	if err := ensureTenantID(student.TenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Contacts").Create(student).Error; err != nil {
			return err
		}
		return replaceStudentContacts(tx, student)
	})
}

func (s *EduStudentService) Update(ctx context.Context, student *models.EduStudent) error {
	if student == nil {
		return errors.New("学生不能为空")
	}
	if err := ensureTenantID(student.TenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduStudent{}), student.TenantID).Where("id = ?", student.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("学生不存在或不属于当前租户")
		}
		if err := tx.Omit("Contacts").Save(student).Error; err != nil {
			return err
		}
		return replaceStudentContacts(tx, student)
	})
}

func (s *EduStudentService) Delete(ctx context.Context, tenantID, id uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduClassMember{}), tenantID).Where("student_id = ? AND status IN ?", id, []string{"studying", "paused"}).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("存在活跃班级成员，禁止删除")
		}
		if err := requireTenant(tx.Where("student_id = ?", id), tenantID).Delete(&models.EduStudentContact{}).Error; err != nil {
			return err
		}
		return requireTenant(tx.Where("id = ?", id), tenantID).Delete(&models.EduStudent{}).Error
	})
}

func replaceStudentContacts(tx *gorm.DB, student *models.EduStudent) error {
	if student == nil {
		return nil
	}
	if err := tx.Where("student_id = ?", student.ID).Delete(&models.EduStudentContact{}).Error; err != nil {
		return err
	}
	for i := range student.Contacts {
		contact := student.Contacts[i]
		if contact == nil {
			continue
		}
		contact.ID = 0
		contact.StudentID = student.ID
		if contact.TenantID == 0 {
			contact.TenantID = student.TenantID
		}
		if err := tx.Create(contact).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *EduStudentService) ImportRows(ctx context.Context, tenantID uint, rows []models.EduStudentImportRow) (*EduImportResult, error) {
	result := &EduImportResult{}
	if err := ensureTenantID(tenantID); err != nil {
		return result, err
	}
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, row := range rows {
			rowNum := i + 1
			phone := strings.TrimSpace(row.Phone)
			if phone == "" {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "phone", Reason: "不能为空"})
				continue
			}
			student := models.EduStudent{}
			err := requireTenant(tx.Model(&models.EduStudent{}), tenantID).Where("phone = ?", phone).First(&student).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				student = models.EduStudent{Name: row.Name, Phone: phone, Gender: row.Gender, TenantID: tenantID}
				if err := tx.Create(&student).Error; err != nil {
					result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "student", Reason: err.Error()})
					continue
				}
				result.Created++
				continue
			}
			student.Name = row.Name
			student.Gender = row.Gender
			student.Phone = phone
			if err := tx.Save(&student).Error; err != nil {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "student", Reason: err.Error()})
				continue
			}
			result.Updated++
		}
		if len(result.Errors) > 0 {
			return errors.New("导入数据存在错误")
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	result.Success = true
	return result, nil
}

func (s *EduStudentService) ExportRows(ctx context.Context, tenantID uint, ids []uint) ([]models.EduStudentExportRow, error) {
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	var list []models.EduStudent
	db := requireTenant(app.DB().WithContext(ctx).Model(&models.EduStudent{}), tenantID)
	if len(ids) > 0 {
		db = db.Where("id IN ?", ids)
	}
	if err := db.Order("id asc").Find(&list).Error; err != nil {
		return nil, err
	}
	rows := make([]models.EduStudentExportRow, 0, len(list))
	for _, item := range list {
		rows = append(rows, models.EduStudentExportRow{
			ID:               item.ID,
			Name:             item.Name,
			Gender:           item.Gender,
			Birthday:         item.Birthday,
			Phone:            item.Phone,
			Status:           item.Status,
			Avatar:           item.Avatar,
			School:           item.School,
			Grade:            item.Grade,
			SchoolClass:      item.SchoolClass,
			StudentType:      item.StudentType,
			SourceChannel:    item.SourceChannel,
			EnrollDate:       item.EnrollDate,
			HealthNote:       item.HealthNote,
			AllergyNote:      item.AllergyNote,
			EmergencyContact: item.EmergencyContact,
			PickupNote:       item.PickupNote,
			Remark:           item.Remark,
			CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return rows, nil
}
