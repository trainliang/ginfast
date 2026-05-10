package models

import (
	"context"
	"gin-fast/app/global/app"

	"gorm.io/gorm"
)

// EduStudent 学生模型
type EduStudent struct {
	BaseModel
	Name             string                `gorm:"column:name;size:100;not null;comment:学生姓名" json:"name"`
	Gender           string                `gorm:"column:gender;size:32;comment:性别" json:"gender"`
	Birthday         *JSONTime             `gorm:"column:birthday;comment:生日" json:"birthday"`
	Phone            string                `gorm:"column:phone;size:64;comment:联系电话" json:"phone"`
	Status           int8                  `gorm:"column:status;default:1;comment:状态" json:"status"`
	Avatar           string                `gorm:"column:avatar;size:255;comment:头像/照片" json:"avatar"`
	School           string                `gorm:"column:school;size:128;comment:就读学校" json:"school"`
	Grade            string                `gorm:"column:grade;size:64;comment:年级" json:"grade"`
	SchoolClass      string                `gorm:"column:school_class;size:64;comment:原班级" json:"schoolClass"`
	SourceChannel    string                `gorm:"column:source_channel;size:64;comment:来源渠道" json:"sourceChannel"`
	EnrollDate       *JSONTime             `gorm:"column:enroll_date;comment:报名日期" json:"enrollDate"`
	HealthNote       string                `gorm:"column:health_note;size:500;comment:健康说明" json:"healthNote"`
	AllergyNote      string                `gorm:"column:allergy_note;size:500;comment:过敏史" json:"allergyNote"`
	EmergencyContact string                `gorm:"column:emergency_contact;size:128;comment:紧急联系人" json:"emergencyContact"`
	PickupNote       string                `gorm:"column:pickup_note;size:500;comment:接送备注" json:"pickupNote"`
	Remark           string                `gorm:"column:remark;size:500;comment:备注" json:"remark"`
	CreatedBy        uint                  `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID         uint                  `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
	Contacts         EduStudentContactList `gorm:"foreignKey:student_id;references:id" json:"contacts"`
}

// EduStudentList 学生列表
type EduStudentList []*EduStudent

// NewEduStudent 创建学生实例
func NewEduStudent() *EduStudent {
	return &EduStudent{}
}

// NewEduStudentList 创建学生列表实例
func NewEduStudentList() EduStudentList {
	return EduStudentList{}
}

// TableName 指定表名
func (EduStudent) TableName() string {
	return "edu_student"
}

// IsEmpty 判断是否为空
func (m *EduStudent) IsEmpty() bool {
	return m == nil || m.ID == 0
}

// Find 查询学生
func (m *EduStudent) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建学生
func (m *EduStudent) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新学生
func (m *EduStudent) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除学生
func (m *EduStudent) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询学生列表
func (l *EduStudentList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduStudent{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取学生总数
func (l EduStudentList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduStudent{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduStudentContact 学生联系人模型
type EduStudentContact struct {
	BaseModel
	StudentID uint   `gorm:"column:student_id;not null;comment:学生ID" json:"studentId"`
	Relation  string `gorm:"column:relation;size:64;comment:关系" json:"relation"`
	Name      string `gorm:"column:name;size:100;not null;comment:联系人姓名" json:"name"`
	Phone     string `gorm:"column:phone;size:64;not null;comment:联系人电话" json:"phone"`
	IsPrimary int8   `gorm:"column:is_primary;default:0;comment:是否主联系人" json:"isPrimary"`
	CanPickup int8   `gorm:"column:can_pickup;default:0;comment:是否允许接送" json:"canPickup"`
	Remark    string `gorm:"column:remark;size:255;comment:备注" json:"remark"`
	CreatedBy uint   `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID  uint   `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

// EduStudentContactList 学生联系人列表
type EduStudentContactList []*EduStudentContact

// NewEduStudentContact 创建学生联系人实例
func NewEduStudentContact() *EduStudentContact {
	return &EduStudentContact{}
}

// NewEduStudentContactList 创建学生联系人列表实例
func NewEduStudentContactList() EduStudentContactList {
	return EduStudentContactList{}
}

// TableName 指定表名
func (EduStudentContact) TableName() string {
	return "edu_student_contact"
}

// IsEmpty 判断是否为空
func (m *EduStudentContact) IsEmpty() bool {
	return m == nil || m.ID == 0
}

// Find 查询学生联系人
func (m *EduStudentContact) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建学生联系人
func (m *EduStudentContact) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新学生联系人
func (m *EduStudentContact) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除学生联系人
func (m *EduStudentContact) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询学生联系人列表
func (l *EduStudentContactList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduStudentContact{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取学生联系人总数
func (l EduStudentContactList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduStudentContact{}).Scopes(query...).Count(&total).Error
	return total, err
}
