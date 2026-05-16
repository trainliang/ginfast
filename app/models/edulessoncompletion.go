package models

import (
	"context"
	"gin-fast/app/global/app"

	"gorm.io/gorm"
)

const (
	EduLessonCompletionResultAttended            = "attended"
	EduLessonCompletionResultStudentLeave        = "student_leave"
	EduLessonCompletionResultStudentAbsent       = "student_absent"
	EduLessonCompletionResultTeacherStopped      = "teacher_stopped"
	EduLessonCompletionResultInstitutionCanceled = "institution_canceled"
	EduLessonCompletionResultMakeupCompleted     = "makeup_completed"

	EduLessonCompletionDeductModeFixedCount = "fixed_count"
	EduLessonCompletionDeductModeDuration   = "duration"

	EduLessonCompletionInsufficientBlock        = "block"
	EduLessonCompletionInsufficientAllowArrears = "allow_arrears"

	EduLessonCompletionStatusCompleted = "completed"
	EduLessonCompletionStatusArrears   = "arrears"
	EduLessonCompletionStatusRevoked   = "revoked"
)

var EduLessonCompletionResultTypes = []string{
	EduLessonCompletionResultAttended,
	EduLessonCompletionResultStudentLeave,
	EduLessonCompletionResultStudentAbsent,
	EduLessonCompletionResultTeacherStopped,
	EduLessonCompletionResultInstitutionCanceled,
	EduLessonCompletionResultMakeupCompleted,
}

type EduLessonCompletionRule struct {
	BaseModel
	ResultType          string `gorm:"column:result_type;size:64;not null;comment:结课结果" json:"resultType"`
	DeductEnabled       int8   `gorm:"column:deduct_enabled;default:0;comment:是否扣减" json:"deductEnabled"`
	DeductMode          string `gorm:"column:deduct_mode;size:32;default:'fixed_count';comment:扣减方式 fixed_count/duration" json:"deductMode"`
	FixedCount          int    `gorm:"column:fixed_count;default:1;comment:固定扣减次数" json:"fixedCount"`
	DurationUnitMinutes int    `gorm:"column:duration_unit_minutes;default:60;comment:多少分钟折算1课时" json:"durationUnitMinutes"`
	InsufficientPolicy  string `gorm:"column:insufficient_policy;size:32;default:'block';comment:权益不足策略 block/allow_arrears" json:"insufficientPolicy"`
	Status              int8   `gorm:"column:status;default:1;comment:状态" json:"status"`
	Remark              string `gorm:"column:remark;size:500;comment:备注" json:"remark"`
	CreatedBy           uint   `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID            uint   `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduLessonCompletionRuleList []*EduLessonCompletionRule

func NewEduLessonCompletionRule() *EduLessonCompletionRule { return &EduLessonCompletionRule{} }

func NewEduLessonCompletionRuleList() EduLessonCompletionRuleList {
	return EduLessonCompletionRuleList{}
}

func (EduLessonCompletionRule) TableName() string { return "edu_lesson_completion_rule" }

func (m *EduLessonCompletionRule) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduLessonCompletionRule) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduLessonCompletionRule) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (l *EduLessonCompletionRuleList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduLessonCompletionRule{}).Scopes(funcs...).Find(l).Error
}

type EduLessonStudentCompletion struct {
	BaseModel
	LessonID         uint      `gorm:"column:lesson_id;not null;comment:课次ID" json:"lessonId"`
	StudentID        uint      `gorm:"column:student_id;not null;comment:学生ID" json:"studentId"`
	ClassID          uint      `gorm:"column:class_id;default:0;comment:班级ID" json:"classId"`
	CourseID         uint      `gorm:"column:course_id;default:0;comment:课程ID" json:"courseId"`
	TeacherID        uint      `gorm:"column:teacher_id;default:0;comment:教师ID" json:"teacherId"`
	ResultType       string    `gorm:"column:result_type;size:64;not null;comment:结课结果" json:"resultType"`
	DeductEnabled    int8      `gorm:"column:deduct_enabled;default:0;comment:是否扣减" json:"deductEnabled"`
	DeductCount      int       `gorm:"column:deduct_count;default:0;comment:实际扣减数量" json:"deductCount"`
	StudentBenefitID uint      `gorm:"column:student_benefit_id;default:0;comment:学生权益ID" json:"studentBenefitId"`
	Status           string    `gorm:"column:status;size:32;not null;comment:状态 completed/arrears/revoked" json:"status"`
	LedgerID         uint      `gorm:"column:ledger_id;default:0;comment:扣减流水ID" json:"ledgerId"`
	RevokeLedgerID   uint      `gorm:"column:revoke_ledger_id;default:0;comment:撤销流水ID" json:"revokeLedgerId"`
	Reason           string    `gorm:"column:reason;size:500;comment:原因" json:"reason"`
	OperatorID       uint      `gorm:"column:operator_id;default:0;comment:操作人ID" json:"operatorId"`
	CompletedAt      *JSONTime `gorm:"column:completed_at;comment:消课时间" json:"completedAt"`
	RevokedAt        *JSONTime `gorm:"column:revoked_at;comment:撤销时间" json:"revokedAt"`
	RevokedBy        uint      `gorm:"column:revoked_by;default:0;comment:撤销人ID" json:"revokedBy"`
	RevokeReason     string    `gorm:"column:revoke_reason;size:500;comment:撤销原因" json:"revokeReason"`
	Version          int       `gorm:"column:version;default:1;comment:版本" json:"version"`
	CreatedBy        uint      `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID         uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduLessonStudentCompletionList []*EduLessonStudentCompletion

func NewEduLessonStudentCompletion() *EduLessonStudentCompletion {
	return &EduLessonStudentCompletion{}
}

func NewEduLessonStudentCompletionList() EduLessonStudentCompletionList {
	return EduLessonStudentCompletionList{}
}

func (EduLessonStudentCompletion) TableName() string { return "edu_lesson_student_completion" }

func (m *EduLessonStudentCompletion) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduLessonStudentCompletion) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduLessonStudentCompletion) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (l *EduLessonStudentCompletionList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduLessonStudentCompletion{}).Scopes(funcs...).Find(l).Error
}
