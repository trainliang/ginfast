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
