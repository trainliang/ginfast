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

type EduRoomService struct{}

func NewEduRoomService() *EduRoomService { return &EduRoomService{} }

func (s *EduRoomService) ValidateEduRoomWeeklyRules(rows []models.EduRoomWeeklyRuleImportRow) error {
	type slot struct {
		roomID  uint
		weekday int8
		start   int
		end     int
	}
	var slots []slot
	for _, row := range rows {
		if row.Available != 1 {
			continue
		}
		start, err := parseHHMM(row.StartTime)
		if err != nil {
			return err
		}
		end, err := parseHHMM(row.EndTime)
		if err != nil {
			return err
		}
		if end <= start {
			return errors.New("结束时间必须晚于开始时间")
		}
		slots = append(slots, slot{roomID: row.RoomID, weekday: row.Weekday, start: start, end: end})
	}
	for i := range slots {
		for j := i + 1; j < len(slots); j++ {
			if slots[i].roomID == slots[j].roomID && slots[i].weekday == slots[j].weekday {
				if slots[i].start < slots[j].end && slots[j].start < slots[i].end {
					return errors.New("同一场地同一星期 available=1 的时间段不能重叠")
				}
			}
		}
	}
	return nil
}

func (s *EduRoomService) Create(ctx context.Context, room *models.EduRoom) error {
	if room == nil {
		return errors.New("场地不能为空")
	}
	if err := ensureTenantID(room.TenantID); err != nil {
		return err
	}
	if room.Capacity <= 0 {
		return errors.New("capacity 必须大于 0")
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduRoom{}), room.TenantID).Where("code = ?", room.Code).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("场地编码已存在")
		}
		return tx.Create(room).Error
	})
}

func (s *EduRoomService) Update(ctx context.Context, room *models.EduRoom) error {
	if room == nil {
		return errors.New("场地不能为空")
	}
	if err := ensureTenantID(room.TenantID); err != nil {
		return err
	}
	if room.Capacity <= 0 {
		return errors.New("capacity 必须大于 0")
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduRoom{}), room.TenantID).Where("id = ?", room.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("场地不存在或不属于当前租户")
		}
		if err := requireTenant(tx.Model(&models.EduRoom{}), room.TenantID).Where("code = ? AND id <> ?", room.Code, room.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("场地编码已存在")
		}
		return tx.Save(room).Error
	})
}

func (s *EduRoomService) Delete(ctx context.Context, tenantID, id uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduClass{}), tenantID).Where("room_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("场地被班级引用，禁止删除")
		}
		if err := requireTenant(tx.Model(&models.EduRoomWeeklyRule{}), tenantID).Where("room_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("场地被周规则引用，禁止删除")
		}
		if err := requireTenant(tx.Model(&models.EduRoomException{}), tenantID).Where("room_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("场地被例外引用，禁止删除")
		}
		return requireTenant(tx.Where("id = ?", id), tenantID).Delete(&models.EduRoom{}).Error
	})
}

func (s *EduRoomService) SaveWeeklyRules(ctx context.Context, tenantID uint, rows []models.EduRoomWeeklyRuleImportRow) (*EduImportResult, error) {
	if err := ensureTenantID(tenantID); err != nil {
		return &EduImportResult{Success: false}, err
	}
	if err := s.ValidateEduRoomWeeklyRules(rows); err != nil {
		return &EduImportResult{Success: false}, err
	}
	result := &EduImportResult{}
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		roomIDs := make(map[uint]bool)
		for _, row := range rows {
			if row.RoomID > 0 {
				roomIDs[row.RoomID] = true
			}
		}
		for roomID := range roomIDs {
			if err := ensureRoomBelongsToTenant(tx, tenantID, roomID); err != nil {
				return err
			}
			if err := requireTenant(tx.Where("room_id = ?", roomID), tenantID).Delete(&models.EduRoomWeeklyRule{}).Error; err != nil {
				return err
			}
		}
		for _, row := range rows {
			rule := models.EduRoomWeeklyRule{RoomID: row.RoomID, Weekday: row.Weekday, StartTime: row.StartTime, EndTime: row.EndTime, Available: row.Available, Remark: row.Remark, TenantID: tenantID}
			if row.Available == 0 {
				rule.Available = 0
			}
			if err := tx.Create(&rule).Error; err != nil {
				return err
			}
			result.Created++
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	result.Success = true
	return result, nil
}

func (s *EduRoomService) AddException(ctx context.Context, exception *models.EduRoomException) error {
	if exception == nil {
		return errors.New("场地例外不能为空")
	}
	tenantID := exception.TenantID
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
		exception.TenantID = tenantID
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensureRoomBelongsToTenant(tx, tenantID, exception.RoomID); err != nil {
			return err
		}
		return tx.Create(exception).Error
	})
}

func (s *EduRoomService) UpdateException(ctx context.Context, exception *models.EduRoomException) error {
	if exception == nil {
		return errors.New("场地例外不能为空")
	}
	tenantID := exception.TenantID
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
		exception.TenantID = tenantID
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensureRoomBelongsToTenant(tx, tenantID, exception.RoomID); err != nil {
			return err
		}
		return requireTenant(tx.Where("id = ?", exception.ID), tenantID).Save(exception).Error
	})
}

func (s *EduRoomService) DeleteException(ctx context.Context, tenantID, id uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return requireTenant(app.DB().WithContext(ctx).Where("id = ?", id), tenantID).Delete(&models.EduRoomException{}).Error
}

func (s *EduRoomService) ImportRows(ctx context.Context, tenantID uint, rows []models.EduRoomImportRow) (*EduImportResult, error) {
	result := &EduImportResult{}
	if err := ensureTenantID(tenantID); err != nil {
		return result, err
	}
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, row := range rows {
			rowNum := i + 1
			if strings.TrimSpace(row.Code) == "" {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "code", Reason: "不能为空"})
				continue
			}
			if row.Capacity <= 0 {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "capacity", Reason: "必须大于 0"})
				continue
			}
			room := models.EduRoom{}
			err := requireTenant(tx.Model(&models.EduRoom{}), tenantID).Where("code = ?", strings.TrimSpace(row.Code)).First(&room).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				room = models.EduRoom{
					Name:     row.Name,
					Code:     strings.TrimSpace(row.Code),
					Type:     row.Type,
					Capacity: row.Capacity,
					Location: row.Location,
					Status:   1,
					Remark:   row.Remark,
					TenantID: tenantID,
				}
				if row.Status != nil {
					room.Status = *row.Status
				}
				if err := tx.Create(&room).Error; err != nil {
					return err
				}
				result.Created++
				continue
			}
			room.Name = row.Name
			room.Type = row.Type
			room.Capacity = row.Capacity
			room.Location = row.Location
			room.Remark = row.Remark
			if row.Status != nil {
				room.Status = *row.Status
			}
			if err := tx.Save(&room).Error; err != nil {
				return err
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

func (s *EduRoomService) ImportWeeklyRows(ctx context.Context, tenantID uint, rows []models.EduRoomWeeklyRuleImportRow) (*EduImportResult, error) {
	result := &EduImportResult{}
	if err := ensureTenantID(tenantID); err != nil {
		return result, err
	}
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.ValidateEduRoomWeeklyRules(rows); err != nil {
			return err
		}
		for i, row := range rows {
			rowNum := i + 1
			if err := ensureRoomBelongsToTenant(tx, tenantID, row.RoomID); err != nil {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "roomId", Reason: err.Error()})
				continue
			}
			rule := models.EduRoomWeeklyRule{RoomID: row.RoomID, Weekday: row.Weekday, StartTime: row.StartTime, EndTime: row.EndTime, Available: row.Available, Remark: row.Remark, TenantID: tenantID}
			if err := tx.Create(&rule).Error; err != nil {
				return err
			}
			result.Created++
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

func (s *EduRoomService) ImportExceptionRows(ctx context.Context, tenantID uint, rows []models.EduRoomExceptionImportRow) (*EduImportResult, error) {
	result := &EduImportResult{}
	if err := ensureTenantID(tenantID); err != nil {
		return result, err
	}
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, row := range rows {
			rowNum := i + 1
			if row.ExceptionDate == nil {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "exceptionDate", Reason: "不能为空"})
				continue
			}
			if strings.TrimSpace(row.Type) == "" {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "type", Reason: "不能为空"})
				continue
			}
			if err := ensureRoomBelongsToTenant(tx, tenantID, row.RoomID); err != nil {
				result.Errors = append(result.Errors, EduImportRowError{Row: rowNum, Field: "roomId", Reason: err.Error()})
				continue
			}
			exception := models.EduRoomException{RoomID: row.RoomID, ExceptionDate: row.ExceptionDate, Type: row.Type, StartTime: row.StartTime, EndTime: row.EndTime, Reason: row.Reason, TenantID: tenantID}
			if err := tx.Create(&exception).Error; err != nil {
				return err
			}
			result.Created++
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

func (s *EduRoomService) ExportRows(ctx context.Context, tenantID uint, ids []uint) ([]models.EduRoomExportRow, error) {
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	var list []models.EduRoom
	db := requireTenant(app.DB().WithContext(ctx).Model(&models.EduRoom{}), tenantID)
	if len(ids) > 0 {
		db = db.Where("id IN ?", ids)
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	rows := make([]models.EduRoomExportRow, 0, len(list))
	for _, item := range list {
		rows = append(rows, models.EduRoomExportRow{ID: item.ID, Name: item.Name, Code: item.Code, Type: item.Type, Capacity: item.Capacity, Location: item.Location, Status: item.Status, Remark: item.Remark, CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"), UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05")})
	}
	return rows, nil
}

func ensureRoomBelongsToTenant(tx *gorm.DB, tenantID, roomID uint) error {
	if roomID == 0 {
		return errors.New("roomId 不能为空")
	}
	var count int64
	if err := requireTenant(tx.Model(&models.EduRoom{}), tenantID).Where("id = ?", roomID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("场地不存在或不属于当前租户")
	}
	return nil
}
