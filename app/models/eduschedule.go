package models

import (
	"context"
	"gin-fast/app/global/app"

	"gorm.io/gorm"
)

// EduScheduleRule 排课规则模型
type EduScheduleRule struct {
	BaseModel
	Name          string    `gorm:"column:name;size:100;not null;comment:规则名称" json:"name"`
	RuleType      string    `gorm:"column:rule_type;size:32;not null;comment:规则类型 class/one_to_one" json:"ruleType"`
	RepeatType    string    `gorm:"column:repeat_type;size:32;not null;comment:重复类型 single/weekly" json:"repeatType"`
	TermID        uint      `gorm:"column:term_id;default:0;comment:学期ID" json:"termId"`
	StartDate     *JSONTime `gorm:"column:start_date;comment:开始日期" json:"startDate"`
	EndDate       *JSONTime `gorm:"column:end_date;comment:结束日期" json:"endDate"`
	ClassID       uint      `gorm:"column:class_id;default:0;comment:班级ID" json:"classId"`
	StudentID     uint      `gorm:"column:student_id;default:0;comment:学生ID" json:"studentId"`
	CourseID      uint      `gorm:"column:course_id;not null;comment:课程/项目ID" json:"courseId"`
	TeacherID     uint      `gorm:"column:teacher_id;not null;comment:教师ID" json:"teacherId"`
	TeachingMode  string    `gorm:"column:teaching_mode;size:32;default:'offline';comment:授课方式 offline/online/home" json:"teachingMode"`
	RequiresRoom  int8      `gorm:"column:requires_room;default:1;comment:是否需要场地" json:"requiresRoom"`
	RoomID        uint      `gorm:"column:room_id;default:0;comment:场地ID" json:"roomId"`
	Weekdays      []int8    `gorm:"-" json:"weekdays"`
	StartTime     string    `gorm:"column:start_time;size:16;not null;comment:开始时间" json:"startTime"`
	EndTime       string    `gorm:"column:end_time;size:16;not null;comment:结束时间" json:"endTime"`
	Status        int8      `gorm:"column:status;default:1;comment:状态" json:"status"`
	Version       int       `gorm:"column:version;default:1;comment:版本" json:"version"`
	EffectiveFrom *JSONTime `gorm:"column:effective_from;comment:生效开始时间" json:"effectiveFrom"`
	Remark        string    `gorm:"column:remark;size:500;comment:备注" json:"remark"`
	CreatedBy     uint      `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID      uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduScheduleRuleList []*EduScheduleRule

func NewEduScheduleRule() *EduScheduleRule { return &EduScheduleRule{} }

func NewEduScheduleRuleList() EduScheduleRuleList { return EduScheduleRuleList{} }

func (EduScheduleRule) TableName() string { return "edu_schedule_rule" }

func (m *EduScheduleRule) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduScheduleRule) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

func (m *EduScheduleRule) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduScheduleRule) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (m *EduScheduleRule) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

func (l *EduScheduleRuleList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduScheduleRule{}).Scopes(funcs...).Find(l).Error
}

func (l EduScheduleRuleList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduScheduleRule{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduLesson 课次模型
type EduLesson struct {
	BaseModel
	RuleID           uint      `gorm:"column:rule_id;default:0;comment:排课规则ID" json:"ruleId"`
	RuleVersion      int       `gorm:"column:rule_version;default:1;comment:规则版本" json:"ruleVersion"`
	LessonType       string    `gorm:"column:lesson_type;size:32;not null;comment:课次类型 class/one_to_one" json:"lessonType"`
	LessonDate       *JSONTime `gorm:"column:lesson_date;not null;comment:上课日期" json:"lessonDate"`
	StartTime        string    `gorm:"column:start_time;size:16;not null;comment:开始时间" json:"startTime"`
	EndTime          string    `gorm:"column:end_time;size:16;not null;comment:结束时间" json:"endTime"`
	ClassID          uint      `gorm:"column:class_id;default:0;comment:班级ID" json:"classId"`
	StudentID        uint      `gorm:"column:student_id;default:0;comment:学生ID" json:"studentId"`
	CourseID         uint      `gorm:"column:course_id;not null;comment:课程/项目ID" json:"courseId"`
	TeacherID        uint      `gorm:"column:teacher_id;not null;comment:教师ID" json:"teacherId"`
	TeachingMode     string    `gorm:"column:teaching_mode;size:32;default:'offline';comment:授课方式" json:"teachingMode"`
	RequiresRoom     int8      `gorm:"column:requires_room;default:1;comment:是否需要场地" json:"requiresRoom"`
	RoomID           uint      `gorm:"column:room_id;default:0;comment:场地ID" json:"roomId"`
	Status           string    `gorm:"column:status;size:32;default:'scheduled';comment:状态" json:"status"`
	IsManualAdjusted int8      `gorm:"column:is_manual_adjusted;default:0;comment:是否人工调整" json:"isManualAdjusted"`
	SourceLessonID   uint      `gorm:"column:source_lesson_id;default:0;comment:来源课次ID" json:"sourceLessonId"`
	CreatedBy        uint      `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID         uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduLessonList []*EduLesson

func NewEduLesson() *EduLesson { return &EduLesson{} }

func NewEduLessonList() EduLessonList { return EduLessonList{} }

func (EduLesson) TableName() string { return "edu_lesson" }

func (m *EduLesson) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduLesson) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

func (m *EduLesson) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduLesson) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (m *EduLesson) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

func (l *EduLessonList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduLesson{}).Scopes(funcs...).Find(l).Error
}

func (l EduLessonList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduLesson{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduLessonChangeLog 课次变更日志模型
type EduLessonChangeLog struct {
	BaseModel
	LessonID   uint      `gorm:"column:lesson_id;not null;comment:课次ID" json:"lessonId"`
	RuleID     uint      `gorm:"column:rule_id;default:0;comment:排课规则ID" json:"ruleId"`
	ActionType string    `gorm:"column:action_type;size:32;not null;comment:动作类型 generate/regenerate" json:"actionType"`
	BeforeData string    `gorm:"column:before_data;type:text;comment:变更前数据" json:"beforeData"`
	AfterData  string    `gorm:"column:after_data;type:text;comment:变更后数据" json:"afterData"`
	Reason     string    `gorm:"column:reason;size:500;comment:原因" json:"reason"`
	OperatorID uint      `gorm:"column:operator_id;default:0;comment:操作人ID" json:"operatorId"`
	OccurredAt *JSONTime `gorm:"column:occurred_at;comment:发生时间" json:"occurredAt"`
	TenantID   uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduLessonChangeLogList []*EduLessonChangeLog

func NewEduLessonChangeLog() *EduLessonChangeLog { return &EduLessonChangeLog{} }

func NewEduLessonChangeLogList() EduLessonChangeLogList { return EduLessonChangeLogList{} }

func (EduLessonChangeLog) TableName() string { return "edu_lesson_change_log" }

func (m *EduLessonChangeLog) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduLessonChangeLog) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (l *EduLessonChangeLogList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduLessonChangeLog{}).Scopes(funcs...).Find(l).Error
}

// EduScheduleConflictOverride 冲突覆盖记录模型
type EduScheduleConflictOverride struct {
	BaseModel
	RuleID       uint      `gorm:"column:rule_id;default:0;comment:排课规则ID" json:"ruleId"`
	LessonID     uint      `gorm:"column:lesson_id;default:0;comment:课次ID" json:"lessonId"`
	ConflictType string    `gorm:"column:conflict_type;size:32;not null;comment:冲突类型 teacher/room/student/class" json:"conflictType"`
	ConflictKey  string    `gorm:"column:conflict_key;size:128;not null;comment:冲突键" json:"conflictKey"`
	Reason       string    `gorm:"column:reason;size:500;not null;comment:覆盖原因" json:"reason"`
	OperatorID   uint      `gorm:"column:operator_id;default:0;comment:操作人ID" json:"operatorId"`
	OccurredAt   *JSONTime `gorm:"column:occurred_at;comment:发生时间" json:"occurredAt"`
	TenantID     uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduScheduleConflictOverrideList []*EduScheduleConflictOverride

func NewEduScheduleConflictOverride() *EduScheduleConflictOverride {
	return &EduScheduleConflictOverride{}
}

func NewEduScheduleConflictOverrideList() EduScheduleConflictOverrideList {
	return EduScheduleConflictOverrideList{}
}

func (EduScheduleConflictOverride) TableName() string { return "edu_schedule_conflict_override" }

func (m *EduScheduleConflictOverride) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduScheduleConflictOverride) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (l *EduScheduleConflictOverrideList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduScheduleConflictOverride{}).Scopes(funcs...).Find(l).Error
}
