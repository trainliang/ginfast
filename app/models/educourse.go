package models

import (
	"context"
	"gin-fast/app/global/app"

	"gorm.io/gorm"
)

// EduCourse 课程/项目模型
type EduCourse struct {
	BaseModel
	Name                string `gorm:"column:name;size:100;not null;comment:课程/项目名称" json:"name"`
	Code                string `gorm:"column:code;size:64;not null;comment:课程/项目编码" json:"code"`
	Type                string `gorm:"column:type;size:64;comment:课程/项目类型" json:"type"`
	GradeRange          string `gorm:"column:grade_range;size:128;comment:适用年龄/年级" json:"gradeRange"`
	DefaultTeachingMode string `gorm:"column:default_teaching_mode;size:32;default:'offline';comment:默认授课方式 offline/online/home" json:"defaultTeachingMode"`
	RequiresRoom        int8   `gorm:"column:requires_room;default:1;comment:是否需要场地" json:"requiresRoom"`
	Status              int8   `gorm:"column:status;default:1;comment:状态 0停用 1启用" json:"status"`
	Sort                int    `gorm:"column:sort;default:0;comment:排序" json:"sort"`
	Description         string `gorm:"column:description;size:500;comment:描述" json:"description"`
	CreatedBy           uint   `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID            uint   `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

// EduCourseList 课程列表
type EduCourseList []*EduCourse

// NewEduCourse 创建课程实例
func NewEduCourse() *EduCourse {
	return &EduCourse{}
}

// NewEduCourseList 创建课程列表实例
func NewEduCourseList() EduCourseList {
	return EduCourseList{}
}

// TableName 指定表名
func (EduCourse) TableName() string {
	return "edu_course"
}

// IsEmpty 判断是否为空
func (m *EduCourse) IsEmpty() bool {
	return m == nil || m.ID == 0
}

// Find 查询课程
func (m *EduCourse) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建课程
func (m *EduCourse) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新课程
func (m *EduCourse) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除课程
func (m *EduCourse) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询课程列表
func (l *EduCourseList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduCourse{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取课程总数
func (l EduCourseList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduCourse{}).Scopes(query...).Count(&total).Error
	return total, err
}
