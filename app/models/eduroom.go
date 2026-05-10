package models

import (
	"context"
	"gin-fast/app/global/app"

	"gorm.io/gorm"
)

// EduRoom 场地模型
type EduRoom struct {
	BaseModel
	Name      string `gorm:"column:name;size:100;not null;comment:场地名称" json:"name"`
	Code      string `gorm:"column:code;size:64;not null;comment:场地编码" json:"code"`
	Type      string `gorm:"column:type;size:64;comment:场地类型" json:"type"`
	Capacity  int    `gorm:"column:capacity;not null;comment:容量" json:"capacity"`
	Location  string `gorm:"column:location;size:255;comment:位置说明" json:"location"`
	Status    int8   `gorm:"column:status;default:1;comment:状态" json:"status"`
	Remark    string `gorm:"column:remark;size:500;comment:备注" json:"remark"`
	CreatedBy uint   `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID  uint   `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

// EduRoomList 场地列表
type EduRoomList []*EduRoom

// NewEduRoom 创建场地实例
func NewEduRoom() *EduRoom {
	return &EduRoom{}
}

// NewEduRoomList 创建场地列表实例
func NewEduRoomList() EduRoomList {
	return EduRoomList{}
}

// TableName 指定表名
func (EduRoom) TableName() string {
	return "edu_room"
}

// IsEmpty 判断是否为空
func (m *EduRoom) IsEmpty() bool {
	return m == nil || m.ID == 0
}

// Find 查询场地
func (m *EduRoom) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建场地
func (m *EduRoom) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新场地
func (m *EduRoom) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除场地
func (m *EduRoom) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询场地列表
func (l *EduRoomList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduRoom{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取场地总数
func (l EduRoomList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduRoom{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduRoomWeeklyRule 场地周规则模型
type EduRoomWeeklyRule struct {
	BaseModel
	RoomID    uint   `gorm:"column:room_id;not null;comment:场地ID" json:"roomId"`
	Weekday   int8   `gorm:"column:weekday;not null;comment:星期 1-7" json:"weekday"`
	StartTime string `gorm:"column:start_time;size:5;not null;comment:开始时间 HH:mm" json:"startTime"`
	EndTime   string `gorm:"column:end_time;size:5;not null;comment:结束时间 HH:mm" json:"endTime"`
	Available int8   `gorm:"column:available;default:1;comment:是否可用" json:"available"`
	Remark    string `gorm:"column:remark;size:255;comment:备注" json:"remark"`
	CreatedBy uint   `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID  uint   `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

// EduRoomWeeklyRuleList 场地周规则列表
type EduRoomWeeklyRuleList []*EduRoomWeeklyRule

// NewEduRoomWeeklyRule 创建场地周规则实例
func NewEduRoomWeeklyRule() *EduRoomWeeklyRule {
	return &EduRoomWeeklyRule{}
}

// NewEduRoomWeeklyRuleList 创建场地周规则列表实例
func NewEduRoomWeeklyRuleList() EduRoomWeeklyRuleList {
	return EduRoomWeeklyRuleList{}
}

// TableName 指定表名
func (EduRoomWeeklyRule) TableName() string {
	return "edu_room_weekly_rule"
}

// IsEmpty 判断是否为空
func (m *EduRoomWeeklyRule) IsEmpty() bool {
	return m == nil || m.ID == 0
}

// Find 查询场地周规则
func (m *EduRoomWeeklyRule) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建场地周规则
func (m *EduRoomWeeklyRule) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新场地周规则
func (m *EduRoomWeeklyRule) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除场地周规则
func (m *EduRoomWeeklyRule) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询场地周规则列表
func (l *EduRoomWeeklyRuleList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduRoomWeeklyRule{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取场地周规则总数
func (l EduRoomWeeklyRuleList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduRoomWeeklyRule{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduRoomException 场地例外模型
type EduRoomException struct {
	BaseModel
	RoomID        uint      `gorm:"column:room_id;not null;comment:场地ID" json:"roomId"`
	ExceptionDate *JSONTime `gorm:"column:exception_date;not null;comment:例外日期" json:"exceptionDate"`
	Type          string    `gorm:"column:type;size:32;not null;comment:例外类型 open/closed" json:"type"`
	StartTime     string    `gorm:"column:start_time;size:5;comment:开始时间 HH:mm" json:"startTime"`
	EndTime       string    `gorm:"column:end_time;size:5;comment:结束时间 HH:mm" json:"endTime"`
	Reason        string    `gorm:"column:reason;size:255;comment:原因" json:"reason"`
	CreatedBy     uint      `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID      uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

// EduRoomExceptionList 场地例外列表
type EduRoomExceptionList []*EduRoomException

// NewEduRoomException 创建场地例外实例
func NewEduRoomException() *EduRoomException {
	return &EduRoomException{}
}

// NewEduRoomExceptionList 创建场地例外列表实例
func NewEduRoomExceptionList() EduRoomExceptionList {
	return EduRoomExceptionList{}
}

// TableName 指定表名
func (EduRoomException) TableName() string {
	return "edu_room_exception"
}

// IsEmpty 判断是否为空
func (m *EduRoomException) IsEmpty() bool {
	return m == nil || m.ID == 0
}

// Find 查询场地例外
func (m *EduRoomException) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

// Create 创建场地例外
func (m *EduRoomException) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

// Update 更新场地例外
func (m *EduRoomException) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

// Delete 删除场地例外
func (m *EduRoomException) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

// Find 查询场地例外列表
func (l *EduRoomExceptionList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduRoomException{}).Scopes(funcs...).Find(l).Error
}

// GetTotal 获取场地例外总数
func (l EduRoomExceptionList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduRoomException{}).Scopes(query...).Count(&total).Error
	return total, err
}
