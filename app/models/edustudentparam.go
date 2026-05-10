package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduStudentListRequest 学生列表请求
type EduStudentListRequest struct {
	BasePaging
	Validator
	ID               *uint  `form:"id" json:"id"`
	Name             string `form:"name" json:"name"`
	Gender           string `form:"gender" json:"gender"`
	Phone            string `form:"phone" json:"phone"`
	Status           *int8  `form:"status" json:"status"`
	School           string `form:"school" json:"school"`
	Grade            string `form:"grade" json:"grade"`
	SchoolClass      string `form:"schoolClass" json:"schoolClass"`
	SourceChannel    string `form:"sourceChannel" json:"sourceChannel"`
	EmergencyContact string `form:"emergencyContact" json:"emergencyContact"`
}

func (r *EduStudentListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduStudentListRequest) Handler() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID != nil {
			db = db.Where("id = ?", *r.ID)
		}
		if r.Name != "" {
			db = db.Where("name LIKE ?", "%"+r.Name+"%")
		}
		if r.Gender != "" {
			db = db.Where("gender = ?", r.Gender)
		}
		if r.Phone != "" {
			db = db.Where("phone LIKE ?", "%"+r.Phone+"%")
		}
		if r.Status != nil {
			db = db.Where("status = ?", *r.Status)
		}
		if r.School != "" {
			db = db.Where("school LIKE ?", "%"+r.School+"%")
		}
		if r.Grade != "" {
			db = db.Where("grade LIKE ?", "%"+r.Grade+"%")
		}
		if r.SchoolClass != "" {
			db = db.Where("school_class LIKE ?", "%"+r.SchoolClass+"%")
		}
		if r.SourceChannel != "" {
			db = db.Where("source_channel = ?", r.SourceChannel)
		}
		if r.EmergencyContact != "" {
			db = db.Where("emergency_contact LIKE ?", "%"+r.EmergencyContact+"%")
		}
		return db
	}
}

// EduStudentAddRequest 新增学生请求
type EduStudentAddRequest struct {
	Validator
	Name             string                        `form:"name" json:"name" validate:"required" message:"学生姓名不能为空"`
	Gender           string                        `form:"gender" json:"gender"`
	Birthday         *JSONTime                     `form:"birthday" json:"birthday"`
	Phone            string                        `form:"phone" json:"phone"`
	Status           *int8                         `form:"status" json:"status"`
	Avatar           string                        `form:"avatar" json:"avatar"`
	School           string                        `form:"school" json:"school"`
	Grade            string                        `form:"grade" json:"grade"`
	SchoolClass      string                        `form:"schoolClass" json:"schoolClass"`
	SourceChannel    string                        `form:"sourceChannel" json:"sourceChannel"`
	EnrollDate       *JSONTime                     `form:"enrollDate" json:"enrollDate"`
	HealthNote       string                        `form:"healthNote" json:"healthNote"`
	AllergyNote      string                        `form:"allergyNote" json:"allergyNote"`
	EmergencyContact string                        `form:"emergencyContact" json:"emergencyContact"`
	PickupNote       string                        `form:"pickupNote" json:"pickupNote"`
	Remark           string                        `form:"remark" json:"remark"`
	Contacts         []EduStudentContactAddRequest `form:"contacts" json:"contacts"`
}

func (r *EduStudentAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduStudentUpdateRequest 更新学生请求
type EduStudentUpdateRequest struct {
	Validator
	ID               uint                             `form:"id" json:"id" validate:"required" message:"学生ID不能为空"`
	Name             string                           `form:"name" json:"name" validate:"required" message:"学生姓名不能为空"`
	Gender           string                           `form:"gender" json:"gender"`
	Birthday         *JSONTime                        `form:"birthday" json:"birthday"`
	Phone            string                           `form:"phone" json:"phone"`
	Status           *int8                            `form:"status" json:"status"`
	Avatar           string                           `form:"avatar" json:"avatar"`
	School           string                           `form:"school" json:"school"`
	Grade            string                           `form:"grade" json:"grade"`
	SchoolClass      string                           `form:"schoolClass" json:"schoolClass"`
	SourceChannel    string                           `form:"sourceChannel" json:"sourceChannel"`
	EnrollDate       *JSONTime                        `form:"enrollDate" json:"enrollDate"`
	HealthNote       string                           `form:"healthNote" json:"healthNote"`
	AllergyNote      string                           `form:"allergyNote" json:"allergyNote"`
	EmergencyContact string                           `form:"emergencyContact" json:"emergencyContact"`
	PickupNote       string                           `form:"pickupNote" json:"pickupNote"`
	Remark           string                           `form:"remark" json:"remark"`
	Contacts         []EduStudentContactUpdateRequest `form:"contacts" json:"contacts"`
}

func (r *EduStudentUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduStudentDeleteRequest 删除学生请求
type EduStudentDeleteRequest struct {
	Validator
	ID uint `form:"id" json:"id" validate:"required" message:"学生ID不能为空"`
}

func (r *EduStudentDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduStudentGetRequest 获取学生请求
type EduStudentGetRequest struct {
	Validator
	ID uint `form:"id" uri:"id" json:"id" validate:"required" message:"学生ID不能为空"`
}

func (r *EduStudentGetRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduStudentContactAddRequest 学生联系人新增请求
type EduStudentContactAddRequest struct {
	ID        uint   `json:"id" form:"id"`
	StudentID uint   `json:"studentId" form:"studentId"`
	Relation  string `json:"relation" form:"relation"`
	Name      string `json:"name" form:"name" validate:"required" message:"联系人姓名不能为空"`
	Phone     string `json:"phone" form:"phone" validate:"required" message:"联系人电话不能为空"`
	IsPrimary *int8  `json:"isPrimary" form:"isPrimary"`
	CanPickup *int8  `json:"canPickup" form:"canPickup"`
	Remark    string `json:"remark" form:"remark"`
}

// EduStudentContactUpdateRequest 学生联系人更新请求
type EduStudentContactUpdateRequest struct {
	ID        uint   `json:"id" form:"id" validate:"required" message:"联系人ID不能为空"`
	StudentID uint   `json:"studentId" form:"studentId"`
	Relation  string `json:"relation" form:"relation"`
	Name      string `json:"name" form:"name" validate:"required" message:"联系人姓名不能为空"`
	Phone     string `json:"phone" form:"phone" validate:"required" message:"联系人电话不能为空"`
	IsPrimary *int8  `json:"isPrimary" form:"isPrimary"`
	CanPickup *int8  `json:"canPickup" form:"canPickup"`
	Remark    string `json:"remark" form:"remark"`
}

// EduStudentContactDeleteRequest 学生联系人删除请求
type EduStudentContactDeleteRequest struct {
	Validator
	ID uint `form:"id" json:"id" validate:"required" message:"联系人ID不能为空"`
}

func (r *EduStudentContactDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduStudentImportRow 学生导入行
type EduStudentImportRow struct {
	Name             string                       `json:"name" form:"name"`
	Gender           string                       `json:"gender" form:"gender"`
	Birthday         *JSONTime                    `json:"birthday" form:"birthday"`
	Phone            string                       `json:"phone" form:"phone"`
	Status           *int8                        `json:"status" form:"status"`
	Avatar           string                       `json:"avatar" form:"avatar"`
	School           string                       `json:"school" form:"school"`
	Grade            string                       `json:"grade" form:"grade"`
	SchoolClass      string                       `json:"schoolClass" form:"schoolClass"`
	SourceChannel    string                       `json:"sourceChannel" form:"sourceChannel"`
	EnrollDate       *JSONTime                    `json:"enrollDate" form:"enrollDate"`
	HealthNote       string                       `json:"healthNote" form:"healthNote"`
	AllergyNote      string                       `json:"allergyNote" form:"allergyNote"`
	EmergencyContact string                       `json:"emergencyContact" form:"emergencyContact"`
	PickupNote       string                       `json:"pickupNote" form:"pickupNote"`
	Remark           string                       `json:"remark" form:"remark"`
	Contacts         []EduStudentContactImportRow `json:"contacts" form:"contacts"`
}

// EduStudentContactImportRow 学生联系人导入行
type EduStudentContactImportRow struct {
	Relation  string `json:"relation" form:"relation"`
	Name      string `json:"name" form:"name"`
	Phone     string `json:"phone" form:"phone"`
	IsPrimary *int8  `json:"isPrimary" form:"isPrimary"`
	CanPickup *int8  `json:"canPickup" form:"canPickup"`
	Remark    string `json:"remark" form:"remark"`
}

// EduStudentExportRow 学生导出行
type EduStudentExportRow struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Gender           string    `json:"gender"`
	Birthday         *JSONTime `json:"birthday"`
	Phone            string    `json:"phone"`
	Status           int8      `json:"status"`
	Avatar           string    `json:"avatar"`
	School           string    `json:"school"`
	Grade            string    `json:"grade"`
	SchoolClass      string    `json:"schoolClass"`
	SourceChannel    string    `json:"sourceChannel"`
	EnrollDate       *JSONTime `json:"enrollDate"`
	HealthNote       string    `json:"healthNote"`
	AllergyNote      string    `json:"allergyNote"`
	EmergencyContact string    `json:"emergencyContact"`
	PickupNote       string    `json:"pickupNote"`
	Remark           string    `json:"remark"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
}

// EduStudentContactExportRow 学生联系人导出行
type EduStudentContactExportRow struct {
	ID        uint   `json:"id"`
	StudentID uint   `json:"studentId"`
	Relation  string `json:"relation"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	IsPrimary int8   `json:"isPrimary"`
	CanPickup int8   `json:"canPickup"`
	Remark    string `json:"remark"`
}

// EduStudentImportRequest 学生导入请求
type EduStudentImportRequest struct {
	Validator
	Rows []EduStudentImportRow `json:"rows" form:"rows" validate:"required" message:"导入数据不能为空"`
}

func (r *EduStudentImportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduStudentExportRequest 学生导出请求
type EduStudentExportRequest struct {
	Validator
	IDs []uint `form:"ids" json:"ids"`
}

func (r *EduStudentExportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}
