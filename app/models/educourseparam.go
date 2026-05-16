package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduCourseListRequest 课程/项目列表请求
type EduCourseListRequest struct {
	BasePaging
	Validator
	ID          *FlexUint `form:"id" json:"id"`
	Name        string `form:"name" json:"name"`
	Code        string `form:"code" json:"code"`
	Type        string `form:"type" json:"type"`
	Status      *FlexInt8 `form:"status" json:"status"`
	GradeRange  string `form:"gradeRange" json:"gradeRange"`
	Description string `form:"description" json:"description"`
}

func (r *EduCourseListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduCourseListRequest) Handler() func(db *gorm.DB) *gorm.DB {
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
		if r.Type != "" {
			db = db.Where("type = ?", r.Type)
		}
		if r.Status != nil {
			db = db.Where("status = ?", int8(*r.Status))
		}
		if r.GradeRange != "" {
			db = db.Where("grade_range LIKE ?", "%"+r.GradeRange+"%")
		}
		if r.Description != "" {
			db = db.Where("description LIKE ?", "%"+r.Description+"%")
		}
		return db
	}
}

// EduCourseAddRequest 新增课程/项目请求
type EduCourseAddRequest struct {
	Validator
	Name                string `form:"name" json:"name" validate:"required" message:"课程/项目名称不能为空"`
	Code                string `form:"code" json:"code" validate:"required" message:"课程/项目编码不能为空"`
	Type                string `form:"type" json:"type"`
	GradeRange          string `form:"gradeRange" json:"gradeRange"`
	DefaultTeachingMode string `form:"defaultTeachingMode" json:"defaultTeachingMode"`
	RequiresRoom        *FlexInt8 `form:"requiresRoom" json:"requiresRoom"`
	Status              *FlexInt8 `form:"status" json:"status"`
	Sort                *FlexInt `form:"sort" json:"sort"`
	Description         string `form:"description" json:"description"`
}

func (r *EduCourseAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduCourseUpdateRequest 更新课程/项目请求
type EduCourseUpdateRequest struct {
	Validator
	ID                  FlexUint `form:"id" json:"id" validate:"required" message:"课程/项目ID不能为空"`
	Name                string `form:"name" json:"name" validate:"required" message:"课程/项目名称不能为空"`
	Code                string `form:"code" json:"code" validate:"required" message:"课程/项目编码不能为空"`
	Type                string `form:"type" json:"type"`
	GradeRange          string `form:"gradeRange" json:"gradeRange"`
	DefaultTeachingMode string `form:"defaultTeachingMode" json:"defaultTeachingMode"`
	RequiresRoom        *FlexInt8 `form:"requiresRoom" json:"requiresRoom"`
	Status              *FlexInt8 `form:"status" json:"status"`
	Sort                *FlexInt `form:"sort" json:"sort"`
	Description         string `form:"description" json:"description"`
}

func (r *EduCourseUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduCourseDeleteRequest 删除课程/项目请求
type EduCourseDeleteRequest struct {
	Validator
	ID FlexUint `form:"id" json:"id" validate:"required" message:"课程/项目ID不能为空"`
}

func (r *EduCourseDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduCourseGetRequest 获取课程/项目请求
type EduCourseGetRequest struct {
	Validator
	ID FlexUint `form:"id" uri:"id" json:"id" validate:"required" message:"课程/项目ID不能为空"`
}

func (r *EduCourseGetRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduCourseImportRow 课程/项目导入行
type EduCourseImportRow struct {
	Name                string `json:"name" form:"name"`
	Code                string `json:"code" form:"code"`
	Type                string `json:"type" form:"type"`
	GradeRange          string `json:"gradeRange" form:"gradeRange"`
	DefaultTeachingMode string `json:"defaultTeachingMode" form:"defaultTeachingMode"`
	RequiresRoom        *FlexInt8 `json:"requiresRoom" form:"requiresRoom"`
	Status              *FlexInt8 `json:"status" form:"status"`
	Sort                *FlexInt `json:"sort" form:"sort"`
	Description         string `json:"description" form:"description"`
}

// EduCourseExportRow 课程/项目导出行
type EduCourseExportRow struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	Code                string `json:"code"`
	Type                string `json:"type"`
	GradeRange          string `json:"gradeRange"`
	DefaultTeachingMode string `json:"defaultTeachingMode"`
	RequiresRoom        int8   `json:"requiresRoom"`
	Status              int8   `json:"status"`
	Sort                int    `json:"sort"`
	Description         string `json:"description"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

// EduCourseImportRequest 课程/项目导入请求
type EduCourseImportRequest struct {
	Validator
	Rows []EduCourseImportRow `json:"rows" form:"rows" validate:"required" message:"导入数据不能为空"`
}

func (r *EduCourseImportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduCourseExportRequest 课程/项目导出请求
type EduCourseExportRequest struct {
	Validator
	IDs []FlexUint `form:"ids" json:"ids"`
}

func (r *EduCourseExportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}
