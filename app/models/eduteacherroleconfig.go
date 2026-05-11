package models

// EduTeacherRoleConfig 教培教师角色配置
type EduTeacherRoleConfig struct {
	BaseModel
	RoleID    uint `gorm:"column:role_id;not null;comment:角色ID;uniqueIndex:idx_edu_teacher_role_tenant_role" json:"roleId"`
	CreatedBy uint `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID  uint `gorm:"type:int(11);column:tenant_id;not null;comment:租户ID;uniqueIndex:idx_edu_teacher_role_tenant_role" json:"tenantID"`
}

// TableName 指定表名
func (EduTeacherRoleConfig) TableName() string {
	return "edu_teacher_role_config"
}
