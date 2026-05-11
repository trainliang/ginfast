package models

import (
	"context"
	"gin-fast/app/global/app"

	"gorm.io/gorm"
)

// EduTerm 学期模型
type EduTerm struct {
	BaseModel
	Name      string    `gorm:"column:name;size:100;not null;comment:学期名称" json:"name"`
	StartDate *JSONTime `gorm:"column:start_date;comment:开始日期" json:"startDate"`
	EndDate   *JSONTime `gorm:"column:end_date;comment:结束日期" json:"endDate"`
	Status    int8      `gorm:"column:status;default:1;comment:状态" json:"status"`
	CreatedBy uint      `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID  uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

// EduTermList 学期列表
type EduTermList []*EduTerm

// NewEduTerm 创建学期实例
func NewEduTerm() *EduTerm { return &EduTerm{} }

// NewEduTermList 创建学期列表实例
func NewEduTermList() EduTermList { return EduTermList{} }

// TableName 指定表名
func (EduTerm) TableName() string { return "edu_term" }

// IsEmpty 判断是否为空
func (m *EduTerm) IsEmpty() bool { return m == nil || m.ID == 0 }

// Find 查询学期
func (m *EduTerm) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建学期
func (m *EduTerm) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新学期
func (m *EduTerm) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除学期
func (m *EduTerm) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询学期列表
func (l *EduTermList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduTerm{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取学期总数
func (l EduTermList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduTerm{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduTermClosedDay 学期停课日模型
type EduTermClosedDay struct {
	BaseModel
	TermID     uint      `gorm:"column:term_id;not null;comment:学期ID" json:"termId"`
	ClosedDate *JSONTime `gorm:"column:closed_date;not null;comment:停课日期" json:"closedDate"`
	Reason     string    `gorm:"column:reason;size:255;comment:原因" json:"reason"`
	CreatedBy  uint      `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID   uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

// EduTermClosedDayList 学期停课日列表
type EduTermClosedDayList []*EduTermClosedDay

// NewEduTermClosedDay 创建学期停课日实例
func NewEduTermClosedDay() *EduTermClosedDay { return &EduTermClosedDay{} }

// NewEduTermClosedDayList 创建学期停课日列表实例
func NewEduTermClosedDayList() EduTermClosedDayList { return EduTermClosedDayList{} }

// TableName 指定表名
func (EduTermClosedDay) TableName() string { return "edu_term_closed_day" }

// IsEmpty 判断是否为空
func (m *EduTermClosedDay) IsEmpty() bool { return m == nil || m.ID == 0 }

// Find 查询学期停课日
func (m *EduTermClosedDay) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建学期停课日
func (m *EduTermClosedDay) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新学期停课日
func (m *EduTermClosedDay) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除学期停课日
func (m *EduTermClosedDay) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询学期停课日列表
func (l *EduTermClosedDayList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduTermClosedDay{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取学期停课日总数
func (l EduTermClosedDayList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduTermClosedDay{}).Scopes(query...).Count(&total).Error
	return total, err
}
