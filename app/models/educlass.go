package models

import (
	"context"
	"gin-fast/app/global/app"

	"gorm.io/gorm"
)

// EduClass 班级模型
type EduClass struct {
	BaseModel
	Name               string    `gorm:"column:name;size:100;not null;comment:班级名称" json:"name"`
	Code               string    `gorm:"column:code;size:64;not null;comment:班级编码" json:"code"`
	ClassType          string    `gorm:"column:class_type;size:32;not null;comment:班级类型 daycare/group" json:"classType"`
	CourseID           uint      `gorm:"column:course_id;not null;comment:课程/项目ID" json:"courseId"`
	TeacherID          uint      `gorm:"column:teacher_id;not null;comment:负责教师用户ID" json:"teacherId"`
	RoomID             uint      `gorm:"column:room_id;default:0;comment:默认场地ID" json:"roomId"`
	Capacity           int       `gorm:"column:capacity;not null;comment:容量" json:"capacity"`
	BenefitCheckPolicy string    `gorm:"column:benefit_check_policy;size:32;default:'required';comment:权益校验策略 required/warn/none" json:"benefitCheckPolicy"`
	Status             int8      `gorm:"column:status;default:1;comment:状态" json:"status"`
	StartDate          *JSONTime `gorm:"column:start_date;comment:开班日期" json:"startDate"`
	EndDate            *JSONTime `gorm:"column:end_date;comment:结班日期" json:"endDate"`
	Remark             string    `gorm:"column:remark;size:500;comment:备注" json:"remark"`
	CreatedBy          uint      `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID           uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

// EduClassList 班级列表
type EduClassList []*EduClass

// NewEduClass 创建班级实例
func NewEduClass() *EduClass {
	return &EduClass{}
}

// NewEduClassList 创建班级列表实例
func NewEduClassList() EduClassList {
	return EduClassList{}
}

// TableName 指定表名
func (EduClass) TableName() string {
	return "edu_class"
}

// IsEmpty 判断是否为空
func (m *EduClass) IsEmpty() bool {
	return m == nil || m.ID == 0
}

// Find 查询班级
func (m *EduClass) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建班级
func (m *EduClass) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新班级
func (m *EduClass) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除班级
func (m *EduClass) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询班级列表
func (l *EduClassList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduClass{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取班级总数
func (l EduClassList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduClass{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduClassMember 班级成员模型
type EduClassMember struct {
	BaseModel
	ClassID   uint      `gorm:"column:class_id;not null;comment:班级ID" json:"classId"`
	StudentID uint      `gorm:"column:student_id;not null;comment:学生ID" json:"studentId"`
	JoinDate  *JSONTime `gorm:"column:join_date;comment:入班日期" json:"joinDate"`
	LeaveDate *JSONTime `gorm:"column:leave_date;comment:退班日期" json:"leaveDate"`
	Status    string    `gorm:"column:status;size:32;not null;comment:状态 studying/paused/left" json:"status"`
	Remark    string    `gorm:"column:remark;size:500;comment:备注" json:"remark"`
	CreatedBy uint      `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID  uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

// EduClassMemberList 班级成员列表
type EduClassMemberList []*EduClassMember

// NewEduClassMember 创建班级成员实例
func NewEduClassMember() *EduClassMember {
	return &EduClassMember{}
}

// NewEduClassMemberList 创建班级成员列表实例
func NewEduClassMemberList() EduClassMemberList {
	return EduClassMemberList{}
}

// TableName 指定表名
func (EduClassMember) TableName() string {
	return "edu_class_member"
}

// IsEmpty 判断是否为空
func (m *EduClassMember) IsEmpty() bool {
	return m == nil || m.ID == 0
}

// Find 查询班级成员
func (m *EduClassMember) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建班级成员
func (m *EduClassMember) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新班级成员
func (m *EduClassMember) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除班级成员
func (m *EduClassMember) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询班级成员列表
func (l *EduClassMemberList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduClassMember{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取班级成员总数
func (l EduClassMemberList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduClassMember{}).Scopes(query...).Count(&total).Error
	return total, err
}
