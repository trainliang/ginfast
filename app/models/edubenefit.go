package models

import (
	"context"
	"gin-fast/app/global/app"

	"gorm.io/gorm"
)

// EduBenefitProduct 权益产品模型
type EduBenefitProduct struct {
	BaseModel
	Name            string `gorm:"column:name;size:100;not null;comment:权益产品名称" json:"name"`
	Code            string `gorm:"column:code;size:64;not null;comment:权益产品编码" json:"code"`
	BenefitType     string `gorm:"column:benefit_type;size:32;not null;comment:权益类型 course/welfare/external" json:"benefitType"`
	CalculationMode string `gorm:"column:calculation_mode;size:32;not null;comment:核算方式 count_limited/period_unlimited/none" json:"calculationMode"`
	TotalCount      int    `gorm:"column:total_count;default:0;comment:总次数" json:"totalCount"`
	ValidDays       int    `gorm:"column:valid_days;default:0;comment:有效天数" json:"validDays"`
	Status          int8   `gorm:"column:status;default:1;comment:状态" json:"status"`
	Remark          string `gorm:"column:remark;size:500;comment:备注" json:"remark"`
	CreatedBy       uint   `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
	TenantID        uint   `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduBenefitProductList []*EduBenefitProduct

func NewEduBenefitProduct() *EduBenefitProduct { return &EduBenefitProduct{} }

func NewEduBenefitProductList() EduBenefitProductList { return EduBenefitProductList{} }

func (EduBenefitProduct) TableName() string { return "edu_benefit_product" }

func (m *EduBenefitProduct) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduBenefitProduct) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

func (m *EduBenefitProduct) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduBenefitProduct) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (m *EduBenefitProduct) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

func (l *EduBenefitProductList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduBenefitProduct{}).Scopes(funcs...).Find(l).Error
}

func (l EduBenefitProductList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduBenefitProduct{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduBenefitProductCourse 权益产品适用课程
type EduBenefitProductCourse struct {
	BaseModel
	ProductID uint `gorm:"column:product_id;not null;comment:权益产品ID" json:"productId"`
	CourseID  uint `gorm:"column:course_id;not null;comment:课程/项目ID" json:"courseId"`
	TenantID  uint `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduBenefitProductCourseList []*EduBenefitProductCourse

func NewEduBenefitProductCourse() *EduBenefitProductCourse { return &EduBenefitProductCourse{} }

func NewEduBenefitProductCourseList() EduBenefitProductCourseList {
	return EduBenefitProductCourseList{}
}

func (EduBenefitProductCourse) TableName() string { return "edu_benefit_product_course" }

func (m *EduBenefitProductCourse) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduBenefitProductCourse) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

func (m *EduBenefitProductCourse) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduBenefitProductCourse) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (m *EduBenefitProductCourse) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

func (l *EduBenefitProductCourseList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduBenefitProductCourse{}).Scopes(funcs...).Find(l).Error
}

func (l EduBenefitProductCourseList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduBenefitProductCourse{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduStudentBenefit 学生权益模型
type EduStudentBenefit struct {
	BaseModel
	StudentID       uint      `gorm:"column:student_id;not null;comment:学生ID" json:"studentId"`
	ProductID       uint      `gorm:"column:product_id;not null;comment:权益产品ID" json:"productId"`
	BenefitType     string    `gorm:"column:benefit_type;size:32;not null;comment:权益类型" json:"benefitType"`
	CalculationMode string    `gorm:"column:calculation_mode;size:32;not null;comment:核算方式" json:"calculationMode"`
	CourseID        uint      `gorm:"column:course_id;default:0;comment:课程/项目ID" json:"courseId"`
	ClassID         uint      `gorm:"column:class_id;default:0;comment:班级ID" json:"classId"`
	TeacherID       uint      `gorm:"column:teacher_id;default:0;comment:教师ID" json:"teacherId"`
	ValidFrom       *JSONTime `gorm:"column:valid_from;comment:开始有效期" json:"validFrom"`
	ValidTo         *JSONTime `gorm:"column:valid_to;comment:结束有效期" json:"validTo"`
	TotalCount      int       `gorm:"column:total_count;default:0;comment:总次数" json:"totalCount"`
	UsedCount       int       `gorm:"column:used_count;default:0;comment:已用次数" json:"usedCount"`
	RemainingCount  int       `gorm:"column:remaining_count;default:0;comment:剩余次数" json:"remainingCount"`
	Status          int8      `gorm:"column:status;default:1;comment:状态" json:"status"`
	SourceType      string    `gorm:"column:source_type;size:32;comment:来源类型" json:"sourceType"`
	TenantID        uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
	CreatedBy       uint      `gorm:"column:created_by;default:0;comment:创建人" json:"createdBy"`
}

type EduStudentBenefitList []*EduStudentBenefit

func NewEduStudentBenefit() *EduStudentBenefit { return &EduStudentBenefit{} }

func NewEduStudentBenefitList() EduStudentBenefitList { return EduStudentBenefitList{} }

func (EduStudentBenefit) TableName() string { return "edu_student_benefit" }

func (m *EduStudentBenefit) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduStudentBenefit) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

func (m *EduStudentBenefit) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduStudentBenefit) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (m *EduStudentBenefit) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

func (l *EduStudentBenefitList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduStudentBenefit{}).Scopes(funcs...).Find(l).Error
}

func (l EduStudentBenefitList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduStudentBenefit{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduBenefitLedger 权益流水模型
type EduBenefitLedger struct {
	BaseModel
	StudentBenefitID uint      `gorm:"column:student_benefit_id;not null;comment:学生权益ID" json:"studentBenefitId"`
	StudentID        uint      `gorm:"column:student_id;not null;comment:学生ID" json:"studentId"`
	ActionType       string    `gorm:"column:action_type;size:32;not null;comment:动作类型 reserve/use/release/adjust/grant_external" json:"actionType"`
	BizType          string    `gorm:"column:biz_type;size:32;not null;comment:业务类型 schedule/lesson/external/manual" json:"bizType"`
	BizID            uint      `gorm:"column:biz_id;default:0;comment:业务ID" json:"bizId"`
	ChangeCount      int       `gorm:"column:change_count;default:0;comment:变更数量" json:"changeCount"`
	BeforeCount      int       `gorm:"column:before_count;default:0;comment:变更前数量" json:"beforeCount"`
	AfterCount       int       `gorm:"column:after_count;default:0;comment:变更后数量" json:"afterCount"`
	OccurredAt       *JSONTime `gorm:"column:occurred_at;comment:发生时间" json:"occurredAt"`
	OperatorID       uint      `gorm:"column:operator_id;default:0;comment:操作人ID" json:"operatorId"`
	Remark           string    `gorm:"column:remark;size:500;comment:备注" json:"remark"`
	TenantID         uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduBenefitLedgerList []*EduBenefitLedger

func NewEduBenefitLedger() *EduBenefitLedger { return &EduBenefitLedger{} }

func NewEduBenefitLedgerList() EduBenefitLedgerList { return EduBenefitLedgerList{} }

func (EduBenefitLedger) TableName() string { return "edu_benefit_ledger" }

func (m *EduBenefitLedger) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduBenefitLedger) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

func (m *EduBenefitLedger) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduBenefitLedger) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (m *EduBenefitLedger) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

func (l *EduBenefitLedgerList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduBenefitLedger{}).Scopes(funcs...).Find(l).Error
}

func (l EduBenefitLedgerList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduBenefitLedger{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduBenefitEvent 权益事件模型
type EduBenefitEvent struct {
	BaseModel
	StudentID        uint      `gorm:"column:student_id;not null;comment:学生ID" json:"studentId"`
	StudentBenefitID uint      `gorm:"column:student_benefit_id;default:0;comment:学生权益ID" json:"studentBenefitId"`
	EventType        string    `gorm:"column:event_type;size:32;not null;comment:事件类型 grant/renew/expire/adjust/cancel" json:"eventType"`
	CourseID         uint      `gorm:"column:course_id;default:0;comment:课程/项目ID" json:"courseId"`
	ClassID          uint      `gorm:"column:class_id;default:0;comment:班级ID" json:"classId"`
	EffectiveFrom    *JSONTime `gorm:"column:effective_from;comment:生效开始时间" json:"effectiveFrom"`
	EffectiveTo      *JSONTime `gorm:"column:effective_to;comment:生效结束时间" json:"effectiveTo"`
	Payload          string    `gorm:"column:payload;type:text;comment:载荷" json:"payload"`
	ProcessedAt      *JSONTime `gorm:"column:processed_at;comment:处理时间" json:"processedAt"`
	TenantID         uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduBenefitEventList []*EduBenefitEvent

func NewEduBenefitEvent() *EduBenefitEvent { return &EduBenefitEvent{} }

func NewEduBenefitEventList() EduBenefitEventList { return EduBenefitEventList{} }

func (EduBenefitEvent) TableName() string { return "edu_benefit_event" }

func (m *EduBenefitEvent) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduBenefitEvent) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

func (m *EduBenefitEvent) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduBenefitEvent) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (m *EduBenefitEvent) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

func (l *EduBenefitEventList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduBenefitEvent{}).Scopes(funcs...).Find(l).Error
}

func (l EduBenefitEventList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduBenefitEvent{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduBenefitExternalSync 外部同步任务模型
type EduBenefitExternalSync struct {
	BaseModel
	StudentBenefitID uint      `gorm:"column:student_benefit_id;not null;comment:学生权益ID" json:"studentBenefitId"`
	StudentID        uint      `gorm:"column:student_id;not null;comment:学生ID" json:"studentId"`
	ProviderCode     string    `gorm:"column:provider_code;size:64;not null;comment:供应商编码" json:"providerCode"`
	IdempotencyKey   string    `gorm:"column:idempotency_key;size:128;not null;comment:幂等键" json:"idempotencyKey"`
	Payload          string    `gorm:"column:payload;type:text;comment:载荷" json:"payload"`
	Status           string    `gorm:"column:status;size:32;not null;comment:状态 pending/success/failed/retrying/canceled" json:"status"`
	ExternalOrderNo  string    `gorm:"column:external_order_no;size:128;comment:外部单号" json:"externalOrderNo"`
	LastError        string    `gorm:"column:last_error;size:500;comment:最后错误" json:"lastError"`
	RetryCount       int       `gorm:"column:retry_count;default:0;comment:重试次数" json:"retryCount"`
	NextRetryAt      *JSONTime `gorm:"column:next_retry_at;comment:下次重试时间" json:"nextRetryAt"`
	TenantID         uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduBenefitExternalSyncList []*EduBenefitExternalSync

func NewEduBenefitExternalSync() *EduBenefitExternalSync { return &EduBenefitExternalSync{} }

func NewEduBenefitExternalSyncList() EduBenefitExternalSyncList { return EduBenefitExternalSyncList{} }

func (EduBenefitExternalSync) TableName() string { return "edu_benefit_external_sync" }

func (m *EduBenefitExternalSync) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduBenefitExternalSync) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

func (m *EduBenefitExternalSync) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduBenefitExternalSync) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (m *EduBenefitExternalSync) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

func (l *EduBenefitExternalSyncList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduBenefitExternalSync{}).Scopes(funcs...).Find(l).Error
}

func (l EduBenefitExternalSyncList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduBenefitExternalSync{}).Scopes(query...).Count(&total).Error
	return total, err
}

// EduLessonStudentEligibility 课次学生资格模型
type EduLessonStudentEligibility struct {
	BaseModel
	LessonID          uint      `gorm:"column:lesson_id;not null;comment:课次ID" json:"lessonId"`
	StudentID         uint      `gorm:"column:student_id;not null;comment:学生ID" json:"studentId"`
	CourseID          uint      `gorm:"column:course_id;not null;comment:课程/项目ID" json:"courseId"`
	ClassID           uint      `gorm:"column:class_id;default:0;comment:班级ID" json:"classId"`
	StudentBenefitID  uint      `gorm:"column:student_benefit_id;default:0;comment:学生权益ID" json:"studentBenefitId"`
	EligibilityStatus string    `gorm:"column:eligibility_status;size:32;not null;comment:资格状态 eligible/ineligible/warn/override" json:"eligibilityStatus"`
	ReasonCode        string    `gorm:"column:reason_code;size:64;comment:原因代码 no_benefit/expired/insufficient_count/class_policy/manual_override" json:"reasonCode"`
	CheckedAt         *JSONTime `gorm:"column:checked_at;comment:检查时间" json:"checkedAt"`
	ResolvedAt        *JSONTime `gorm:"column:resolved_at;comment:解决时间" json:"resolvedAt"`
	TenantID          uint      `gorm:"type:int(11);column:tenant_id;comment:租户ID" json:"tenantID"`
}

type EduLessonStudentEligibilityList []*EduLessonStudentEligibility

func NewEduLessonStudentEligibility() *EduLessonStudentEligibility {
	return &EduLessonStudentEligibility{}
}

func NewEduLessonStudentEligibilityList() EduLessonStudentEligibilityList {
	return EduLessonStudentEligibilityList{}
}

func (EduLessonStudentEligibility) TableName() string { return "edu_lesson_student_eligibility" }

func (m *EduLessonStudentEligibility) IsEmpty() bool { return m == nil || m.ID == 0 }

func (m *EduLessonStudentEligibility) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Find(m).Error
}

func (m *EduLessonStudentEligibility) Create(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Scopes(funcs...).Create(m).Error
}

func (m *EduLessonStudentEligibility) Update(ctx context.Context) error {
	return app.DB().WithContext(ctx).Save(m).Error
}

func (m *EduLessonStudentEligibility) Delete(ctx context.Context) error {
	return app.DB().WithContext(ctx).Delete(m).Error
}

func (l *EduLessonStudentEligibilityList) Find(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) error {
	return app.DB().WithContext(ctx).Model(&EduLessonStudentEligibility{}).Scopes(funcs...).Find(l).Error
}

func (l EduLessonStudentEligibilityList) GetTotal(ctx context.Context, query ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := app.DB().WithContext(ctx).Model(&EduLessonStudentEligibility{}).Scopes(query...).Count(&total).Error
	return total, err
}
