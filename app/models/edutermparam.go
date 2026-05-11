package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduTermListRequest 学期列表请求
type EduTermListRequest struct {
	BasePaging
	Validator
	ID     *uint  `form:"id" json:"id"`
	Name   string `form:"name" json:"name"`
	Status *int8  `form:"status" json:"status"`
}

func (r *EduTermListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduTermListRequest) Handler() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID != nil {
			db = db.Where("id = ?", *r.ID)
		}
		if r.Name != "" {
			db = db.Where("name LIKE ?", "%"+r.Name+"%")
		}
		if r.Status != nil {
			db = db.Where("status = ?", *r.Status)
		}
		return db
	}
}

// EduTermAddRequest 新增学期请求
type EduTermAddRequest struct {
	Validator
	Name      string    `form:"name" json:"name" validate:"required" message:"学期名称不能为空"`
	StartDate *JSONTime `form:"startDate" json:"startDate" validate:"required" message:"开始日期不能为空"`
	EndDate   *JSONTime `form:"endDate" json:"endDate" validate:"required" message:"结束日期不能为空"`
	Status    *int8     `form:"status" json:"status"`
}

func (r *EduTermAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduTermUpdateRequest 更新学期请求
type EduTermUpdateRequest struct {
	Validator
	ID        uint      `form:"id" json:"id" validate:"required" message:"学期ID不能为空"`
	Name      string    `form:"name" json:"name" validate:"required" message:"学期名称不能为空"`
	StartDate *JSONTime `form:"startDate" json:"startDate" validate:"required" message:"开始日期不能为空"`
	EndDate   *JSONTime `form:"endDate" json:"endDate" validate:"required" message:"结束日期不能为空"`
	Status    *int8     `form:"status" json:"status"`
}

func (r *EduTermUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduTermDeleteRequest 删除学期请求
type EduTermDeleteRequest struct {
	Validator
	ID uint `form:"id" json:"id" validate:"required" message:"学期ID不能为空"`
}

func (r *EduTermDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduTermClosedDaysSaveRequest 保存学期停课日请求
type EduTermClosedDaysSaveRequest struct {
	Validator
	Rows []EduTermClosedDay `form:"rows" json:"rows" validate:"required" message:"停课日不能为空"`
}

func (r *EduTermClosedDaysSaveRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}
