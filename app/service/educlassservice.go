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

type EduClassService struct{}

func NewEduClassService() *EduClassService { return &EduClassService{} }

func (s *EduClassService) Create(ctx context.Context, class *models.EduClass) error {
	if class == nil {
		return errors.New("班级不能为空")
	}
	if strings.EqualFold(class.ClassType, "one_to_one") || strings.Contains(class.ClassType, "一对一") {
		return errors.New("一对一预约不属于班级管理范围")
	}
	if class.ClassType != "daycare" && class.ClassType != "group" {
		return errors.New("classType 只允许 daycare/group")
	}
	if class.Capacity <= 0 {
		return errors.New("capacity 必须大于 0")
	}
	tenantID := class.TenantID
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
		class.TenantID = tenantID
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduClass{}), tenantID).Where("code = ?", class.Code).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("班级编码已存在")
		}
		if err := validateClassReferences(tx, tenantID, class.CourseID, class.TeacherID, class.RoomID); err != nil {
			return err
		}
		return tx.Create(class).Error
	})
}

func (s *EduClassService) Update(ctx context.Context, class *models.EduClass) error {
	if class == nil {
		return errors.New("班级不能为空")
	}
	if strings.EqualFold(class.ClassType, "one_to_one") || strings.Contains(class.ClassType, "一对一") {
		return errors.New("一对一预约不属于班级管理范围")
	}
	if class.ClassType != "daycare" && class.ClassType != "group" {
		return errors.New("classType 只允许 daycare/group")
	}
	if class.Capacity <= 0 {
		return errors.New("capacity 必须大于 0")
	}
	tenantID := class.TenantID
	if tenantID == 0 {
		tenantID = tenantIDFromContext(ctx)
		class.TenantID = tenantID
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduClass{}), tenantID).Where("id = ?", class.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("班级不存在或不属于当前租户")
		}
		if err := requireTenant(tx.Model(&models.EduClass{}), tenantID).Where("code = ? AND id <> ?", class.Code, class.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("班级编码已存在")
		}
		if err := validateClassReferences(tx, tenantID, class.CourseID, class.TeacherID, class.RoomID); err != nil {
			return err
		}
		return tx.Save(class).Error
	})
}

func (s *EduClassService) Delete(ctx context.Context, tenantID, id uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := requireTenant(tx.Model(&models.EduClassMember{}), tenantID).Where("class_id = ? AND status IN ?", id, []string{"studying", "paused"}).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("班级存在活跃成员，禁止删除")
		}
		return requireTenant(tx.Where("id = ?", id), tenantID).Delete(&models.EduClass{}).Error
	})
}

func (s *EduClassService) AddMember(ctx context.Context, member *models.EduClassMember) error {
	return s.saveMember(ctx, member, false)
}

func (s *EduClassService) UpdateMember(ctx context.Context, member *models.EduClassMember) error {
	return s.saveMember(ctx, member, true)
}

func (s *EduClassService) DeleteMember(ctx context.Context, tenantID, id uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	return requireTenant(app.DB().WithContext(ctx).Where("id = ?", id), tenantID).Delete(&models.EduClassMember{}).Error
}

func (s *EduClassService) saveMember(ctx context.Context, member *models.EduClassMember, allowUpdate bool) error {
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.saveMemberTx(tx, member, allowUpdate)
	})
}

func (s *EduClassService) saveMemberTx(tx *gorm.DB, member *models.EduClassMember, allowUpdate bool) error {
	if member == nil {
		return errors.New("班级成员不能为空")
	}
	if !isActiveMemberStatus(member.Status) && strings.TrimSpace(member.Status) != "left" {
		return errors.New("成员状态不合法")
	}
	tenantID := member.TenantID
	if tenantID == 0 {
		tenantID = tenantIDFromContext(tx.Statement.Context)
		member.TenantID = tenantID
	}
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	var class models.EduClass
	if err := requireTenant(tx.Model(&models.EduClass{}), tenantID).Where("id = ?", member.ClassID).First(&class).Error; err != nil {
		return err
	}
	if class.Capacity <= 0 {
		return errors.New("capacity 必须大于 0")
	}
	if err := ensureStudentBelongsToTenant(tx, tenantID, member.StudentID); err != nil {
		return err
	}
	var activeCount int64
	memberScope := requireTenant(tx.Model(&models.EduClassMember{}), tenantID).Where("class_id = ? AND status IN ?", member.ClassID, []string{"studying", "paused"})
	if err := memberScope.Count(&activeCount).Error; err != nil {
		return err
	}
	if allowUpdate && member.ID > 0 {
		var current models.EduClassMember
		if err := requireTenant(tx.Model(&models.EduClassMember{}), tenantID).Where("id = ?", member.ID).First(&current).Error; err != nil {
			return err
		}
		if current.ClassID == member.ClassID && isActiveMemberStatus(current.Status) {
			activeCount--
		}
	}
	if isActiveMemberStatus(member.Status) && int(activeCount) >= class.Capacity {
		return fmt.Errorf("班级容量已满")
	}
	var dupCount int64
	dupScope := requireTenant(tx.Model(&models.EduClassMember{}), tenantID).Where("class_id = ? AND student_id = ? AND status IN ?", member.ClassID, member.StudentID, []string{"studying", "paused"})
	if allowUpdate && member.ID > 0 {
		dupScope = dupScope.Where("id <> ?", member.ID)
	}
	if err := dupScope.Count(&dupCount).Error; err != nil {
		return err
	}
	if dupCount > 0 && isActiveMemberStatus(member.Status) {
		return fmt.Errorf("同一学生同一班级只能有一条有效记录")
	}
	if member.ID == 0 {
		return tx.Create(member).Error
	}
	return tx.Save(member).Error
}

func (s *EduClassService) ImportRows(ctx context.Context, tenantID uint, rows []models.EduClassImportRow) (*EduImportResult, error) {
	result := &EduImportResult{}
	if err := ensureTenantID(tenantID); err != nil {
		return result, err
	}
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, row := range rows {
			if strings.EqualFold(row.ClassType, "one_to_one") || strings.Contains(row.ClassType, "一对一") {
				result.Errors = append(result.Errors, EduImportRowError{Row: i + 1, Field: "classType", Reason: "一对一预约不属于班级管理范围"})
				continue
			}
			if row.ClassType != "daycare" && row.ClassType != "group" {
				result.Errors = append(result.Errors, EduImportRowError{Row: i + 1, Field: "classType", Reason: "classType 只允许 daycare/group"})
				continue
			}
			if row.Capacity <= 0 {
				result.Errors = append(result.Errors, EduImportRowError{Row: i + 1, Field: "capacity", Reason: "必须大于 0"})
				continue
			}
			if err := validateClassReferences(tx, tenantID, row.CourseID, row.TeacherID, roomIDOrZero(row.RoomID)); err != nil {
				result.Errors = append(result.Errors, EduImportRowError{Row: i + 1, Field: "reference", Reason: err.Error()})
				continue
			}
		}
		if len(result.Errors) > 0 {
			return errors.New("导入数据存在错误")
		}
		for _, row := range rows {
			var class models.EduClass
			err := requireTenant(tx.Model(&models.EduClass{}), tenantID).Where("code = ?", row.Code).First(&class).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				class = models.EduClass{Name: row.Name, Code: row.Code, ClassType: row.ClassType, CourseID: row.CourseID, TeacherID: row.TeacherID, Capacity: row.Capacity, TenantID: tenantID}
				if row.RoomID != nil {
					class.RoomID = *row.RoomID
				}
				if row.Status != nil {
					class.Status = *row.Status
				}
				if err := tx.Create(&class).Error; err != nil {
					return err
				}
				result.Created++
				continue
			}
			class.Name = row.Name
			class.ClassType = row.ClassType
			class.CourseID = row.CourseID
			class.TeacherID = row.TeacherID
			class.Capacity = row.Capacity
			if row.RoomID != nil {
				class.RoomID = *row.RoomID
			}
			if row.Status != nil {
				class.Status = *row.Status
			}
			if err := tx.Save(&class).Error; err != nil {
				return err
			}
			result.Updated++
		}
		result.Success = true
		return nil
	})
	return result, err
}

func (s *EduClassService) ImportMemberRows(ctx context.Context, tenantID uint, rows []models.EduClassMemberImportRow) (*EduImportResult, error) {
	result := &EduImportResult{}
	if err := ensureTenantID(tenantID); err != nil {
		return result, err
	}
	err := app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, row := range rows {
			member := &models.EduClassMember{ClassID: row.ClassID, StudentID: row.StudentID, Status: row.Status, TenantID: tenantID}
			if err := s.saveMemberTx(tx, member, false); err != nil {
				result.Errors = append(result.Errors, EduImportRowError{Row: i + 1, Field: "member", Reason: err.Error()})
				continue
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

func (s *EduClassService) ExportRows(ctx context.Context, tenantID uint, ids []uint) ([]models.EduClassExportRow, error) {
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	var list []models.EduClass
	db := requireTenant(app.DB().WithContext(ctx).Model(&models.EduClass{}), tenantID)
	if len(ids) > 0 {
		db = db.Where("id IN ?", ids)
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	rows := make([]models.EduClassExportRow, 0, len(list))
	for _, item := range list {
		rows = append(rows, models.EduClassExportRow{ID: item.ID, Name: item.Name, Code: item.Code, ClassType: item.ClassType, CourseID: item.CourseID, TeacherID: item.TeacherID, RoomID: item.RoomID, Capacity: item.Capacity, Status: item.Status, StartDate: item.StartDate, EndDate: item.EndDate, Remark: item.Remark, CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"), UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05")})
	}
	return rows, nil
}

func (s *EduClassService) ExportMemberRows(ctx context.Context, tenantID uint, ids []uint) ([]models.EduClassMemberExportRow, error) {
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	var list []models.EduClassMember
	db := requireTenant(app.DB().WithContext(ctx).Model(&models.EduClassMember{}), tenantID)
	if len(ids) > 0 {
		db = db.Where("id IN ?", ids)
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	rows := make([]models.EduClassMemberExportRow, 0, len(list))
	for _, item := range list {
		rows = append(rows, models.EduClassMemberExportRow{ID: item.ID, ClassID: item.ClassID, StudentID: item.StudentID, JoinDate: item.JoinDate, LeaveDate: item.LeaveDate, Status: item.Status, Remark: item.Remark, CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"), UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05")})
	}
	return rows, nil
}

func (s *EduClassService) TeacherOptions(ctx context.Context, tenantID uint) ([]models.EduTeacherOption, error) {
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	var users []models.User
	err := app.DB().WithContext(ctx).Model(&models.User{}).
		Where("sys_users.tenant_id = ?", tenantID).
		Select("DISTINCT sys_users.id, sys_users.username, sys_users.nick_name").
		Joins("JOIN sys_user_role sur ON sur.user_id = sys_users.id").
		Joins("JOIN edu_teacher_role_config etrc ON etrc.role_id = sur.role_id AND etrc.tenant_id = ?", tenantID).
		Joins("JOIN sys_role r ON r.id = etrc.role_id AND r.tenant_id = ? AND r.status = 1", tenantID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	opts := make([]models.EduTeacherOption, 0, len(users))
	for _, user := range users {
		opts = append(opts, models.EduTeacherOption{ID: user.ID, Username: user.Username, NickName: user.NickName, Name: user.NickName})
	}
	return opts, nil
}

func (s *EduClassService) TeacherRoleConfig(ctx context.Context, tenantID uint) (*models.EduTeacherRoleConfigResponse, error) {
	if err := ensureTenantID(tenantID); err != nil {
		return nil, err
	}
	var roles []models.SysRole
	if err := app.DB().WithContext(ctx).Where("tenant_id = ?", tenantID).Order("sort ASC, id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	var configs []models.EduTeacherRoleConfig
	if err := app.DB().WithContext(ctx).Where("tenant_id = ?", tenantID).Order("role_id ASC").Find(&configs).Error; err != nil {
		return nil, err
	}
	resp := &models.EduTeacherRoleConfigResponse{
		Roles:           make([]models.EduTeacherRoleConfigRoleItem, 0, len(roles)),
		SelectedRoleIDs: make([]uint, 0, len(configs)),
	}
	for _, role := range roles {
		resp.Roles = append(resp.Roles, models.EduTeacherRoleConfigRoleItem{ID: role.ID, Name: role.Name, Status: role.Status, TenantID: role.TenantID})
	}
	for _, config := range configs {
		resp.SelectedRoleIDs = append(resp.SelectedRoleIDs, config.RoleID)
	}
	return resp, nil
}

func (s *EduClassService) SaveTeacherRoleConfig(ctx context.Context, tenantID uint, roleIDs []uint, userID uint) error {
	if err := ensureTenantID(tenantID); err != nil {
		return err
	}
	uniqueRoleIDs := dedupeUint(roleIDs)
	return app.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(uniqueRoleIDs) > 0 {
			var count int64
			if err := tx.Model(&models.SysRole{}).Where("tenant_id = ? AND status = 1 AND id IN ?", tenantID, uniqueRoleIDs).Count(&count).Error; err != nil {
				return err
			}
			if count != int64(len(uniqueRoleIDs)) {
				return errors.New("教师角色只能选择当前租户启用角色")
			}
		}
		if err := tx.Unscoped().Where("tenant_id = ?", tenantID).Delete(&models.EduTeacherRoleConfig{}).Error; err != nil {
			return err
		}
		for _, roleID := range uniqueRoleIDs {
			if err := tx.Create(&models.EduTeacherRoleConfig{RoleID: roleID, CreatedBy: userID, TenantID: tenantID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func validateClassReferences(tx *gorm.DB, tenantID, courseID, teacherID, roomID uint) error {
	if courseID > 0 {
		var count int64
		if err := requireTenant(tx.Model(&models.EduCourse{}), tenantID).Where("id = ?", courseID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("课程不存在或不属于当前租户")
		}
	}
	if teacherID > 0 {
		var count int64
		if err := requireTenant(tx.Model(&models.User{}), tenantID).Where("id = ?", teacherID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("教师不存在或不属于当前租户")
		}
		if err := tx.Model(&models.SysUserRole{}).
			Joins("JOIN edu_teacher_role_config etrc ON etrc.role_id = sys_user_role.role_id AND etrc.tenant_id = ?", tenantID).
			Joins("JOIN sys_role r ON r.id = etrc.role_id AND r.tenant_id = ? AND r.status = 1", tenantID).
			Where("sys_user_role.user_id = ?", teacherID).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("教师未关联当前租户配置的教师角色")
		}
	}
	if roomID > 0 {
		var count int64
		if err := requireTenant(tx.Model(&models.EduRoom{}), tenantID).Where("id = ?", roomID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("场地不存在或不属于当前租户")
		}
	}
	return nil
}

func dedupeUint(values []uint) []uint {
	seen := make(map[uint]struct{}, len(values))
	result := make([]uint, 0, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func roomIDOrZero(roomID *uint) uint {
	if roomID == nil {
		return 0
	}
	return *roomID
}

func ensureStudentBelongsToTenant(tx *gorm.DB, tenantID, studentID uint) error {
	if studentID == 0 {
		return errors.New("studentId 不能为空")
	}
	var count int64
	if err := requireTenant(tx.Model(&models.EduStudent{}), tenantID).Where("id = ?", studentID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("学生不存在或不属于当前租户")
	}
	return nil
}
