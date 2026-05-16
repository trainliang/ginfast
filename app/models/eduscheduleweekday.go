package models

import (
	"context"
	"gin-fast/app/global/app"

	"gorm.io/gorm"
)

// EduScheduleRuleWeekday 排课规则周几子表
type EduScheduleRuleWeekday struct {
	BaseModel
	RuleID   uint `gorm:"column:rule_id;not null;comment:排课规则ID" json:"ruleId"`
	Weekday  int8 `gorm:"column:weekday;not null;comment:星期 1-7" json:"weekday"`
	TenantID uint `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

func (EduScheduleRuleWeekday) TableName() string { return "edu_schedule_rule_weekday" }

func (m *EduScheduleRuleWeekday) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}
