package service

import (
	"testing"

	"gin-fast/app/models"
)

func TestEduClassServiceCompile(t *testing.T) {
	setupEduTestDB(t)
	_ = NewEduClassService()
}

func TestEduClassServiceRejectsCapacityOverflow(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	if err := db.Create(&models.EduStudent{Name: "学生A", Phone: "13800000000", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed student: %v", err)
	}

	if err := db.Create(&models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed class: %v", err)
	}
	if err := db.Create(&models.EduClassMember{ClassID: 1, StudentID: 1, Status: "studying", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	err := svc.AddMember(ctx, &models.EduClassMember{ClassID: 1, StudentID: 2, Status: "studying", TenantID: 1})
	if err == nil {
		t.Fatalf("expected capacity overflow error")
	}
}

func TestEduClassServiceRejectsDuplicateActiveMember(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	if err := db.Create(&models.EduStudent{Name: "学生A", Phone: "13800000000", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed student: %v", err)
	}

	if err := db.Create(&models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 3, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed class: %v", err)
	}
	if err := svc.AddMember(ctx, &models.EduClassMember{ClassID: 1, StudentID: 1, Status: "studying", TenantID: 1}); err != nil {
		t.Fatalf("seed member: %v", err)
	}

	err := svc.AddMember(ctx, &models.EduClassMember{ClassID: 1, StudentID: 1, Status: "paused", TenantID: 1})
	if err == nil {
		t.Fatalf("expected duplicate active member error")
	}
}

func TestEduClassServiceRejectsOneOnOneClassType(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduClassService()
	err := svc.Create(contextWithTenant(1), &models.EduClass{Name: "班级", Code: "CL001", ClassType: "one_to_one", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1})
	if err == nil {
		t.Fatalf("expected one-on-one class type rejection")
	}
}

func TestEduClassServiceRejectsZeroTenantCreate(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduClassService()
	err := svc.Create(contextWithTenant(0), &models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 0})
	if err == nil {
		t.Fatalf("expected zero tenant create to be rejected")
	}
}

func TestEduClassServiceRejectsCrossTenantUpdate(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	if err := db.Create(&models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 0}).Error; err != nil {
		t.Fatalf("seed class: %v", err)
	}
	err := svc.Update(contextWithTenant(1), &models.EduClass{BaseModel: models.BaseModel{ID: 1}, Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1})
	if err == nil {
		t.Fatalf("expected cross tenant update to be rejected")
	}
}

func TestEduClassServiceDefaultsBenefitCheckPolicy(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 1).First(&role).Error; err != nil {
		t.Fatalf("load role: %v", err)
	}
	if err := db.Create(&models.EduTeacherRoleConfig{RoleID: role.ID, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed teacher role config: %v", err)
	}

	class := &models.EduClass{Name: "班级", Code: "CL002", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1}
	if err := svc.Create(ctx, class); err != nil {
		t.Fatalf("create class: %v", err)
	}

	var saved models.EduClass
	if err := db.Where("code = ? AND tenant_id = ?", "CL002", 1).First(&saved).Error; err != nil {
		t.Fatalf("load class: %v", err)
	}
	if saved.BenefitCheckPolicy != "required" {
		t.Fatalf("expected benefitCheckPolicy required, got %q", saved.BenefitCheckPolicy)
	}
}

func TestEduClassServiceNormalizesBenefitCheckPolicy(t *testing.T) {
	setupEduTestDB(t)
	class := &models.EduClass{BenefitCheckPolicy: " WARN "}
	normalizeBenefitCheckPolicy(class)
	if class.BenefitCheckPolicy != "warn" {
		t.Fatalf("expected warn policy, got %q", class.BenefitCheckPolicy)
	}

	class.BenefitCheckPolicy = "invalid"
	normalizeBenefitCheckPolicy(class)
	if class.BenefitCheckPolicy != "required" {
		t.Fatalf("expected invalid policy to default required, got %q", class.BenefitCheckPolicy)
	}
}

func TestEduClassServiceUpdatePreservesOmittedBenefitCheckPolicy(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 1).First(&role).Error; err != nil {
		t.Fatalf("load role: %v", err)
	}
	if err := db.Create(&models.EduTeacherRoleConfig{RoleID: role.ID, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed teacher role config: %v", err)
	}
	if err := db.Create(&models.EduClass{Name: "班级", Code: "CL002", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, BenefitCheckPolicy: "warn", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed class: %v", err)
	}

	if err := svc.Update(ctx, &models.EduClass{
		BaseModel: models.BaseModel{ID: 1},
		Name:      "班级2",
		Code:      "CL002",
		ClassType: "group",
		CourseID:  1,
		TeacherID: 1,
		Capacity:  1,
		TenantID:  1,
	}); err != nil {
		t.Fatalf("update class: %v", err)
	}

	var saved models.EduClass
	if err := db.Where("id = ?", 1).First(&saved).Error; err != nil {
		t.Fatalf("load class: %v", err)
	}
	if saved.BenefitCheckPolicy != "warn" {
		t.Fatalf("expected omitted benefitCheckPolicy to preserve warn, got %q", saved.BenefitCheckPolicy)
	}
}

func TestEduClassServiceImportExistingPreservesOmittedBenefitCheckPolicy(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 1).First(&role).Error; err != nil {
		t.Fatalf("load role: %v", err)
	}
	if err := db.Create(&models.EduTeacherRoleConfig{RoleID: role.ID, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed teacher role config: %v", err)
	}
	if err := db.Create(&models.EduClass{Name: "班级", Code: "CL002", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, BenefitCheckPolicy: "none", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed class: %v", err)
	}

	result, err := svc.ImportRows(ctx, 1, []models.EduClassImportRow{{Name: "班级2", Code: "CL002", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1}})
	if err != nil {
		t.Fatalf("import class: %v", err)
	}
	if !result.Success || result.Updated != 1 {
		t.Fatalf("expected one updated row, got %+v", result)
	}

	var saved models.EduClass
	if err := db.Where("code = ? AND tenant_id = ?", "CL002", 1).First(&saved).Error; err != nil {
		t.Fatalf("load class: %v", err)
	}
	if saved.BenefitCheckPolicy != "none" {
		t.Fatalf("expected omitted import benefitCheckPolicy to preserve none, got %q", saved.BenefitCheckPolicy)
	}
}

func TestEduClassServiceRejectsUnconfiguredTeacherRole(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)

	err := svc.Create(ctx, &models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1})
	if err == nil {
		t.Fatalf("expected unconfigured teacher role to be rejected")
	}
}

func TestEduClassServiceAllowsConfiguredTeacherRole(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 1).First(&role).Error; err != nil {
		t.Fatalf("load role: %v", err)
	}
	if err := db.Create(&models.EduTeacherRoleConfig{RoleID: role.ID, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed teacher role config: %v", err)
	}

	err := svc.Create(ctx, &models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1})
	if err != nil {
		t.Fatalf("expected configured teacher role to pass, got %v", err)
	}
}

func TestEduClassServiceRejectsCrossTenantTeacherRoleConfig(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	seedEduTenantData(t, 2)
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 2).First(&role).Error; err != nil {
		t.Fatalf("load tenant 2 role: %v", err)
	}
	if err := db.Create(&models.EduTeacherRoleConfig{RoleID: role.ID, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed cross tenant teacher role config: %v", err)
	}

	err := svc.Create(ctx, &models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1})
	if err == nil {
		t.Fatalf("expected cross tenant teacher role config to be rejected")
	}
}

func TestEduClassServiceRejectsDisabledConfiguredTeacherRole(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 1).First(&role).Error; err != nil {
		t.Fatalf("load role: %v", err)
	}
	if err := db.Model(&models.SysRole{}).Where("id = ?", role.ID).Update("status", 0).Error; err != nil {
		t.Fatalf("disable role: %v", err)
	}
	if err := db.Create(&models.EduTeacherRoleConfig{RoleID: role.ID, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed teacher role config: %v", err)
	}

	err := svc.Create(ctx, &models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1})
	if err == nil {
		t.Fatalf("expected disabled configured teacher role to be rejected")
	}
}

func TestEduClassServiceUpdateMemberAllowsSameRecordWithoutCapacityDrop(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	if err := db.Create(&models.EduStudent{Name: "学生A", Phone: "13800000000", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed student: %v", err)
	}

	if err := db.Create(&models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed class: %v", err)
	}
	member := &models.EduClassMember{ClassID: 1, StudentID: 1, Status: "studying", TenantID: 1}
	if err := db.Create(member).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	updated := &models.EduClassMember{BaseModel: models.BaseModel{ID: member.ID}, ClassID: 1, StudentID: 1, Status: "paused", TenantID: 1}
	if err := svc.UpdateMember(ctx, updated); err != nil {
		t.Fatalf("update member: %v", err)
	}

	var count int64
	if err := db.Model(&models.EduClassMember{}).Where("class_id = ? AND status IN ?", 1, []string{"studying", "paused"}).Count(&count).Error; err != nil {
		t.Fatalf("count active members: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected no false capacity drop, got %d", count)
	}
}

func TestEduClassServiceImportMemberRowsRollsBackOnError(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)
	if err := db.Create(&models.EduStudent{Name: "学生A", Phone: "13800000000", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed student: %v", err)
	}

	if err := db.Create(&models.EduClass{Name: "班级", Code: "CL001", ClassType: "group", CourseID: 1, TeacherID: 1, Capacity: 1, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed class: %v", err)
	}
	if err := db.Create(&models.EduStudent{Name: "学生A", Phone: "13800000000", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed student: %v", err)
	}

	result, err := svc.ImportMemberRows(ctx, 1, []models.EduClassMemberImportRow{
		{ClassID: 1, StudentID: 1, Status: "studying"},
		{ClassID: 1, StudentID: 1, Status: "bad"},
	})
	if err == nil {
		t.Fatalf("expected import error")
	}
	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("expected row errors")
	}

	var count int64
	if err := db.Model(&models.EduClassMember{}).Where("tenant_id = ?", 1).Count(&count).Error; err != nil {
		t.Fatalf("count members: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected rollback without partial insert, got %d", count)
	}
}

func TestEduClassServiceTeacherOptionsDeduplicatesUsers(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	ctx := contextWithTenant(1)
	seedEduTenantData(t, 1)

	var user models.User
	if err := db.Where("username = ?", "teacher-1").First(&user).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if err := db.Create(&models.SysRole{Name: "教师", Status: 1, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed duplicate role: %v", err)
	}

	var roles []models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 1).Find(&roles).Error; err != nil {
		t.Fatalf("load roles: %v", err)
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 teacher roles, got %d", len(roles))
	}
	for _, role := range roles {
		if err := db.Create(&models.EduTeacherRoleConfig{RoleID: role.ID, TenantID: 1}).Error; err != nil {
			t.Fatalf("seed teacher role config: %v", err)
		}
	}
	if err := db.Create(&models.SysUserRole{UserID: user.ID, RoleID: roles[1].ID}).Error; err != nil {
		t.Fatalf("seed user-role: %v", err)
	}

	opts, err := svc.TeacherOptions(ctx, 1)
	if err != nil {
		t.Fatalf("teacher options: %v", err)
	}
	if len(opts) != 1 {
		t.Fatalf("expected distinct teacher options, got %d", len(opts))
	}
}

func TestEduClassServiceSaveTeacherRoleConfigRejectsCrossTenantRole(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	seedEduTenantData(t, 1)
	seedEduTenantData(t, 2)
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 2).First(&role).Error; err != nil {
		t.Fatalf("load tenant 2 role: %v", err)
	}

	err := svc.SaveTeacherRoleConfig(contextWithTenant(1), 1, []uint{role.ID}, 99)
	if err == nil {
		t.Fatalf("expected cross tenant role config to be rejected")
	}
}

func TestEduClassServiceSaveTeacherRoleConfigRejectsDisabledRole(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	seedEduTenantData(t, 1)
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 1).First(&role).Error; err != nil {
		t.Fatalf("load role: %v", err)
	}
	if err := db.Model(&models.SysRole{}).Where("id = ?", role.ID).Update("status", 0).Error; err != nil {
		t.Fatalf("disable role: %v", err)
	}

	err := svc.SaveTeacherRoleConfig(contextWithTenant(1), 1, []uint{role.ID}, 99)
	if err == nil {
		t.Fatalf("expected disabled role config to be rejected")
	}
}

func TestEduClassServiceTeacherRoleConfigReturnsTenantRolesAndSelectedIDs(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduClassService()
	seedEduTenantData(t, 1)
	seedEduTenantData(t, 2)
	if err := db.Create(&models.SysRole{Name: "助教", Status: 1, TenantID: 1}).Error; err != nil {
		t.Fatalf("seed assistant role: %v", err)
	}
	var role models.SysRole
	if err := db.Where("name = ? AND tenant_id = ?", "教师", 1).First(&role).Error; err != nil {
		t.Fatalf("load tenant 1 role: %v", err)
	}
	if err := svc.SaveTeacherRoleConfig(contextWithTenant(1), 1, []uint{role.ID}, 99); err != nil {
		t.Fatalf("save teacher role config: %v", err)
	}

	config, err := svc.TeacherRoleConfig(contextWithTenant(1), 1)
	if err != nil {
		t.Fatalf("teacher role config: %v", err)
	}
	if len(config.SelectedRoleIDs) != 1 || config.SelectedRoleIDs[0] != role.ID {
		t.Fatalf("unexpected selected role ids: %#v", config.SelectedRoleIDs)
	}
	for _, item := range config.Roles {
		if item.TenantID != 1 {
			t.Fatalf("expected only tenant 1 roles, got role %#v", item)
		}
	}
	if len(config.Roles) != 2 {
		t.Fatalf("expected 2 tenant roles, got %d", len(config.Roles))
	}
}
