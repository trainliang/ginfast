package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduBenefitProductListRequest 权益产品列表请求
type EduBenefitProductListRequest struct {
	BasePaging
	Validator
	ID              *uint  `form:"id" json:"id"`
	Name            string `form:"name" json:"name"`
	Code            string `form:"code" json:"code"`
	BenefitType     string `form:"benefitType" json:"benefitType"`
	CalculationMode string `form:"calculationMode" json:"calculationMode"`
	Status          *int8  `form:"status" json:"status"`
}

func (r *EduBenefitProductListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduBenefitProductAddRequest 新增权益产品请求
type EduBenefitProductAddRequest struct {
	Validator
	Name            string `form:"name" json:"name" validate:"required" message:"权益产品名称不能为空"`
	Code            string `form:"code" json:"code" validate:"required" message:"权益产品编码不能为空"`
	BenefitType     string `form:"benefitType" json:"benefitType"`
	CalculationMode string `form:"calculationMode" json:"calculationMode"`
	TotalCount      int    `form:"totalCount" json:"totalCount"`
	ValidDays       int    `form:"validDays" json:"validDays"`
	Status          *int8  `form:"status" json:"status"`
	Remark          string `form:"remark" json:"remark"`
}

func (r *EduBenefitProductAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduBenefitProductUpdateRequest 更新权益产品请求
type EduBenefitProductUpdateRequest struct {
	Validator
	ID              uint   `form:"id" json:"id" validate:"required" message:"权益产品ID不能为空"`
	Name            string `form:"name" json:"name" validate:"required" message:"权益产品名称不能为空"`
	Code            string `form:"code" json:"code" validate:"required" message:"权益产品编码不能为空"`
	BenefitType     string `form:"benefitType" json:"benefitType"`
	CalculationMode string `form:"calculationMode" json:"calculationMode"`
	TotalCount      int    `form:"totalCount" json:"totalCount"`
	ValidDays       int    `form:"validDays" json:"validDays"`
	Status          *int8  `form:"status" json:"status"`
	Remark          string `form:"remark" json:"remark"`
}

func (r *EduBenefitProductUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduBenefitProductDeleteRequest 删除权益产品请求
type EduBenefitProductDeleteRequest struct {
	Validator
	ID uint `form:"id" json:"id" validate:"required" message:"权益产品ID不能为空"`
}

func (r *EduBenefitProductDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduBenefitProductListRequest) Handler() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID != nil {
			db = db.Where("id = ?", *r.ID)
		}
		if r.Name != "" {
			db = db.Where("name LIKE ?", "%"+r.Name+"%")
		}
		if r.Code != "" {
			db = db.Where("code LIKE ?", "%"+r.Code+"%")
		}
		if r.BenefitType != "" {
			db = db.Where("benefit_type = ?", r.BenefitType)
		}
		if r.CalculationMode != "" {
			db = db.Where("calculation_mode = ?", r.CalculationMode)
		}
		if r.Status != nil {
			db = db.Where("status = ?", *r.Status)
		}
		return db
	}
}

// EduStudentBenefitListRequest 学生权益列表请求
type EduStudentBenefitListRequest struct {
	BasePaging
	Validator
	ID              *uint  `form:"id" json:"id"`
	StudentID       *uint  `form:"studentId" json:"studentId"`
	ProductID       *uint  `form:"productId" json:"productId"`
	BenefitType     string `form:"benefitType" json:"benefitType"`
	CalculationMode string `form:"calculationMode" json:"calculationMode"`
	CourseID        *uint  `form:"courseId" json:"courseId"`
	ClassID         *uint  `form:"classId" json:"classId"`
	TeacherID       *uint  `form:"teacherId" json:"teacherId"`
	Status          *int8  `form:"status" json:"status"`
}

func (r *EduStudentBenefitListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduStudentBenefitAddRequest 新增学生权益请求
type EduStudentBenefitAddRequest struct {
	Validator
	StudentID       uint      `form:"studentId" json:"studentId" validate:"required" message:"学生ID不能为空"`
	ProductID       uint      `form:"productId" json:"productId"`
	BenefitType     string    `form:"benefitType" json:"benefitType"`
	CalculationMode string    `form:"calculationMode" json:"calculationMode"`
	CourseID        uint      `form:"courseId" json:"courseId" validate:"required" message:"课程ID不能为空"`
	ClassID         uint      `form:"classId" json:"classId"`
	TeacherID       uint      `form:"teacherId" json:"teacherId"`
	ValidFrom       *JSONTime `form:"validFrom" json:"validFrom"`
	ValidTo         *JSONTime `form:"validTo" json:"validTo"`
	TotalCount      int       `form:"totalCount" json:"totalCount"`
	UsedCount       int       `form:"usedCount" json:"usedCount"`
	RemainingCount  int       `form:"remainingCount" json:"remainingCount"`
	Status          *int8     `form:"status" json:"status"`
	SourceType      string    `form:"sourceType" json:"sourceType"`
}

func (r *EduStudentBenefitAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduStudentBenefitUpdateRequest 更新学生权益请求
type EduStudentBenefitUpdateRequest struct {
	Validator
	ID              uint      `form:"id" json:"id" validate:"required" message:"学生权益ID不能为空"`
	StudentID       uint      `form:"studentId" json:"studentId" validate:"required" message:"学生ID不能为空"`
	ProductID       uint      `form:"productId" json:"productId"`
	BenefitType     string    `form:"benefitType" json:"benefitType"`
	CalculationMode string    `form:"calculationMode" json:"calculationMode"`
	CourseID        uint      `form:"courseId" json:"courseId" validate:"required" message:"课程ID不能为空"`
	ClassID         uint      `form:"classId" json:"classId"`
	TeacherID       uint      `form:"teacherId" json:"teacherId"`
	ValidFrom       *JSONTime `form:"validFrom" json:"validFrom"`
	ValidTo         *JSONTime `form:"validTo" json:"validTo"`
	TotalCount      int       `form:"totalCount" json:"totalCount"`
	UsedCount       int       `form:"usedCount" json:"usedCount"`
	RemainingCount  int       `form:"remainingCount" json:"remainingCount"`
	Status          *int8     `form:"status" json:"status"`
	SourceType      string    `form:"sourceType" json:"sourceType"`
}

func (r *EduStudentBenefitUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduBenefitLedgerListRequest 权益流水列表请求
type EduBenefitLedgerListRequest struct {
	BasePaging
	Validator
	ID               *uint  `form:"id" json:"id"`
	StudentBenefitID *uint  `form:"studentBenefitId" json:"studentBenefitId"`
	StudentID        *uint  `form:"studentId" json:"studentId"`
	ActionType       string `form:"actionType" json:"actionType"`
	BizType          string `form:"bizType" json:"bizType"`
}

func (r *EduBenefitLedgerListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduBenefitLedgerListRequest) Handler() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID != nil {
			db = db.Where("id = ?", *r.ID)
		}
		if r.StudentBenefitID != nil {
			db = db.Where("student_benefit_id = ?", *r.StudentBenefitID)
		}
		if r.StudentID != nil {
			db = db.Where("student_id = ?", *r.StudentID)
		}
		if r.ActionType != "" {
			db = db.Where("action_type = ?", r.ActionType)
		}
		if r.BizType != "" {
			db = db.Where("biz_type = ?", r.BizType)
		}
		return db
	}
}

// EduBenefitExternalSyncListRequest 外部同步任务列表请求
type EduBenefitExternalSyncListRequest struct {
	BasePaging
	Validator
	ID               *uint  `form:"id" json:"id"`
	StudentBenefitID *uint  `form:"studentBenefitId" json:"studentBenefitId"`
	StudentID        *uint  `form:"studentId" json:"studentId"`
	ProviderCode     string `form:"providerCode" json:"providerCode"`
	IdempotencyKey   string `form:"idempotencyKey" json:"idempotencyKey"`
	Status           string `form:"status" json:"status"`
}

func (r *EduBenefitExternalSyncListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduBenefitExternalSyncListRequest) Handler() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID != nil {
			db = db.Where("id = ?", *r.ID)
		}
		if r.StudentBenefitID != nil {
			db = db.Where("student_benefit_id = ?", *r.StudentBenefitID)
		}
		if r.StudentID != nil {
			db = db.Where("student_id = ?", *r.StudentID)
		}
		if r.ProviderCode != "" {
			db = db.Where("provider_code = ?", r.ProviderCode)
		}
		if r.IdempotencyKey != "" {
			db = db.Where("idempotency_key = ?", r.IdempotencyKey)
		}
		if r.Status != "" {
			db = db.Where("status = ?", r.Status)
		}
		return db
	}
}

// EduBenefitExternalSyncRetryRequest 外部同步重试请求
type EduBenefitExternalSyncRetryRequest struct {
	Validator
	ID uint `form:"id" json:"id" validate:"required" message:"外部同步任务ID不能为空"`
}

func (r *EduBenefitExternalSyncRetryRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduStudentBenefitListRequest) Handler() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID != nil {
			db = db.Where("id = ?", *r.ID)
		}
		if r.StudentID != nil {
			db = db.Where("student_id = ?", *r.StudentID)
		}
		if r.ProductID != nil {
			db = db.Where("product_id = ?", *r.ProductID)
		}
		if r.BenefitType != "" {
			db = db.Where("benefit_type = ?", r.BenefitType)
		}
		if r.CalculationMode != "" {
			db = db.Where("calculation_mode = ?", r.CalculationMode)
		}
		if r.CourseID != nil {
			db = db.Where("course_id = ?", *r.CourseID)
		}
		if r.ClassID != nil {
			db = db.Where("class_id = ?", *r.ClassID)
		}
		if r.TeacherID != nil {
			db = db.Where("teacher_id = ?", *r.TeacherID)
		}
		if r.Status != nil {
			db = db.Where("status = ?", *r.Status)
		}
		return db
	}
}

// EduBenefitCheckRequest 权益资格校验请求
type EduBenefitCheckRequest struct {
	Validator
	TenantID  *uint     `json:"tenantID" form:"tenantID"`
	StudentID uint      `json:"studentId" form:"studentId" validate:"required" message:"学生ID不能为空"`
	CourseID  uint      `json:"courseId" form:"courseId" validate:"required" message:"课程ID不能为空"`
	ClassID   uint      `json:"classId" form:"classId"`
	TeacherID uint      `json:"teacherId" form:"teacherId"`
	CheckedAt *JSONTime `json:"checkedAt" form:"checkedAt"`
}

func (r *EduBenefitCheckRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduBenefitCheckResult 权益资格校验结果
type EduBenefitCheckResult struct {
	Eligible         bool   `json:"eligible"`
	Status           string `json:"status"`
	ReasonCode       string `json:"reasonCode"`
	StudentBenefitID uint   `json:"studentBenefitId"`
}

// EduBenefitRepairRequest 课次资格修复请求
type EduBenefitRepairRequest struct {
	Validator
	TenantID      *uint     `json:"tenantID" form:"tenantID"`
	StudentID     *uint     `json:"studentId" form:"studentId"`
	CourseID      *uint     `json:"courseId" form:"courseId"`
	EffectiveFrom *JSONTime `json:"effectiveFrom" form:"effectiveFrom"`
}

func (r *EduBenefitRepairRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduBenefitRepairResult 修复结果
type EduBenefitRepairResult struct {
	CheckedCount  int `json:"checkedCount"`
	RepairedCount int `json:"repairedCount"`
}
