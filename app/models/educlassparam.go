package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduClassListRequest 班级列表请求
type EduClassListRequest struct {
	BasePaging
	Validator
	ID        *uint  `form:"id" json:"id"`
	Name      string `form:"name" json:"name"`
	Code      string `form:"code" json:"code"`
	ClassType string `form:"classType" json:"classType"`
	CourseID  *uint  `form:"courseId" json:"courseId"`
	TeacherID *uint  `form:"teacherId" json:"teacherId"`
	RoomID    *uint  `form:"roomId" json:"roomId"`
	Status    *int8  `form:"status" json:"status"`
}

func (r *EduClassListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduClassListRequest) Handler() func(db *gorm.DB) *gorm.DB {
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
		if r.ClassType != "" {
			db = db.Where("class_type = ?", r.ClassType)
		}
		if r.CourseID != nil {
			db = db.Where("course_id = ?", *r.CourseID)
		}
		if r.TeacherID != nil {
			db = db.Where("teacher_id = ?", *r.TeacherID)
		}
		if r.RoomID != nil {
			db = db.Where("room_id = ?", *r.RoomID)
		}
		if r.Status != nil {
			db = db.Where("status = ?", *r.Status)
		}
		return db
	}
}

// EduClassAddRequest 新增班级请求
type EduClassAddRequest struct {
	Validator
	Name               string    `form:"name" json:"name" validate:"required" message:"班级名称不能为空"`
	Code               string    `form:"code" json:"code" validate:"required" message:"班级编码不能为空"`
	ClassType          string    `form:"classType" json:"classType" validate:"required" message:"班级类型不能为空"`
	CourseID           uint      `form:"courseId" json:"courseId" validate:"required" message:"课程/项目ID不能为空"`
	TeacherID          uint      `form:"teacherId" json:"teacherId" validate:"required" message:"负责教师不能为空"`
	RoomID             *uint     `form:"roomId" json:"roomId"`
	Capacity           int       `form:"capacity" json:"capacity" validate:"required" message:"容量不能为空"`
	BenefitCheckPolicy string    `form:"benefitCheckPolicy" json:"benefitCheckPolicy"`
	Status             *int8     `form:"status" json:"status"`
	StartDate          *JSONTime `form:"startDate" json:"startDate"`
	EndDate            *JSONTime `form:"endDate" json:"endDate"`
	Remark             string    `form:"remark" json:"remark"`
}

func (r *EduClassAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduClassUpdateRequest 更新班级请求
type EduClassUpdateRequest struct {
	Validator
	ID                 uint      `form:"id" json:"id" validate:"required" message:"班级ID不能为空"`
	Name               string    `form:"name" json:"name" validate:"required" message:"班级名称不能为空"`
	Code               string    `form:"code" json:"code" validate:"required" message:"班级编码不能为空"`
	ClassType          string    `form:"classType" json:"classType" validate:"required" message:"班级类型不能为空"`
	CourseID           uint      `form:"courseId" json:"courseId" validate:"required" message:"课程/项目ID不能为空"`
	TeacherID          uint      `form:"teacherId" json:"teacherId" validate:"required" message:"负责教师不能为空"`
	RoomID             *uint     `form:"roomId" json:"roomId"`
	Capacity           int       `form:"capacity" json:"capacity" validate:"required" message:"容量不能为空"`
	BenefitCheckPolicy string    `form:"benefitCheckPolicy" json:"benefitCheckPolicy"`
	Status             *int8     `form:"status" json:"status"`
	StartDate          *JSONTime `form:"startDate" json:"startDate"`
	EndDate            *JSONTime `form:"endDate" json:"endDate"`
	Remark             string    `form:"remark" json:"remark"`
}

func (r *EduClassUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduClassDeleteRequest 删除班级请求
type EduClassDeleteRequest struct {
	Validator
	ID uint `form:"id" json:"id" validate:"required" message:"班级ID不能为空"`
}

func (r *EduClassDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduClassGetRequest 获取班级请求
type EduClassGetRequest struct {
	Validator
	ID uint `form:"id" uri:"id" json:"id" validate:"required" message:"班级ID不能为空"`
}

func (r *EduClassGetRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduClassMemberListRequest 班级成员列表请求
type EduClassMemberListRequest struct {
	BasePaging
	Validator
	ID        *uint  `form:"id" json:"id"`
	ClassID   *uint  `form:"classId" json:"classId"`
	StudentID *uint  `form:"studentId" json:"studentId"`
	Status    string `form:"status" json:"status"`
}

func (r *EduClassMemberListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduClassMemberListRequest) Handler() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID != nil {
			db = db.Where("id = ?", *r.ID)
		}
		if r.ClassID != nil {
			db = db.Where("class_id = ?", *r.ClassID)
		}
		if r.StudentID != nil {
			db = db.Where("student_id = ?", *r.StudentID)
		}
		if r.Status != "" {
			db = db.Where("status = ?", r.Status)
		}
		return db
	}
}

// EduClassMemberAddRequest 新增班级成员请求
type EduClassMemberAddRequest struct {
	Validator
	ClassID   uint      `form:"classId" json:"classId" validate:"required" message:"班级ID不能为空"`
	StudentID uint      `form:"studentId" json:"studentId" validate:"required" message:"学生ID不能为空"`
	JoinDate  *JSONTime `form:"joinDate" json:"joinDate"`
	LeaveDate *JSONTime `form:"leaveDate" json:"leaveDate"`
	Status    string    `form:"status" json:"status" validate:"required" message:"状态不能为空"`
	Remark    string    `form:"remark" json:"remark"`
}

func (r *EduClassMemberAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduClassMemberUpdateRequest 更新班级成员请求
type EduClassMemberUpdateRequest struct {
	Validator
	ID        uint      `form:"id" json:"id" validate:"required" message:"班级成员ID不能为空"`
	ClassID   uint      `form:"classId" json:"classId" validate:"required" message:"班级ID不能为空"`
	StudentID uint      `form:"studentId" json:"studentId" validate:"required" message:"学生ID不能为空"`
	JoinDate  *JSONTime `form:"joinDate" json:"joinDate"`
	LeaveDate *JSONTime `form:"leaveDate" json:"leaveDate"`
	Status    string    `form:"status" json:"status" validate:"required" message:"状态不能为空"`
	Remark    string    `form:"remark" json:"remark"`
}

func (r *EduClassMemberUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduClassMemberDeleteRequest 删除班级成员请求
type EduClassMemberDeleteRequest struct {
	Validator
	ID uint `form:"id" json:"id" validate:"required" message:"班级成员ID不能为空"`
}

func (r *EduClassMemberDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduClassImportRow 班级导入行
type EduClassImportRow struct {
	Name               string    `json:"name" form:"name"`
	Code               string    `json:"code" form:"code"`
	ClassType          string    `json:"classType" form:"classType"`
	CourseID           uint      `json:"courseId" form:"courseId"`
	TeacherID          uint      `json:"teacherId" form:"teacherId"`
	RoomID             *uint     `json:"roomId" form:"roomId"`
	Capacity           int       `json:"capacity" form:"capacity"`
	BenefitCheckPolicy string    `json:"benefitCheckPolicy" form:"benefitCheckPolicy"`
	Status             *int8     `json:"status" form:"status"`
	StartDate          *JSONTime `json:"startDate" form:"startDate"`
	EndDate            *JSONTime `json:"endDate" form:"endDate"`
	Remark             string    `json:"remark" form:"remark"`
}

// EduClassMemberImportRow 班级成员导入行
type EduClassMemberImportRow struct {
	ClassID   uint      `json:"classId" form:"classId"`
	StudentID uint      `json:"studentId" form:"studentId"`
	JoinDate  *JSONTime `json:"joinDate" form:"joinDate"`
	LeaveDate *JSONTime `json:"leaveDate" form:"leaveDate"`
	Status    string    `json:"status" form:"status"`
	Remark    string    `json:"remark" form:"remark"`
}

// EduClassExportRow 班级导出行
type EduClassExportRow struct {
	ID                 uint      `json:"id"`
	Name               string    `json:"name"`
	Code               string    `json:"code"`
	ClassType          string    `json:"classType"`
	CourseID           uint      `json:"courseId"`
	TeacherID          uint      `json:"teacherId"`
	RoomID             uint      `json:"roomId"`
	Capacity           int       `json:"capacity"`
	BenefitCheckPolicy string    `json:"benefitCheckPolicy"`
	Status             int8      `json:"status"`
	StartDate          *JSONTime `json:"startDate"`
	EndDate            *JSONTime `json:"endDate"`
	Remark             string    `json:"remark"`
	CreatedAt          string    `json:"createdAt"`
	UpdatedAt          string    `json:"updatedAt"`
}

// EduClassMemberExportRow 班级成员导出行
type EduClassMemberExportRow struct {
	ID        uint      `json:"id"`
	ClassID   uint      `json:"classId"`
	StudentID uint      `json:"studentId"`
	JoinDate  *JSONTime `json:"joinDate"`
	LeaveDate *JSONTime `json:"leaveDate"`
	Status    string    `json:"status"`
	Remark    string    `json:"remark"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

// EduClassImportRequest 班级导入请求
type EduClassImportRequest struct {
	Validator
	Rows []EduClassImportRow `json:"rows" form:"rows" validate:"required" message:"导入数据不能为空"`
}

func (r *EduClassImportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduClassMemberImportRequest 班级成员导入请求
type EduClassMemberImportRequest struct {
	Validator
	Rows []EduClassMemberImportRow `json:"rows" form:"rows" validate:"required" message:"导入数据不能为空"`
}

func (r *EduClassMemberImportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduClassExportRequest 班级导出请求
type EduClassExportRequest struct {
	Validator
	IDs []uint `form:"ids" json:"ids"`
}

func (r *EduClassExportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduTeacherRoleConfigRequest 保存教师角色配置请求
type EduTeacherRoleConfigRequest struct {
	Validator
	RoleIDs []uint `form:"roleIds" json:"roleIds"`
}

func (r *EduTeacherRoleConfigRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduTeacherRoleConfigRoleItem 教师角色配置中的角色选项
type EduTeacherRoleConfigRoleItem struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Status   int8   `json:"status"`
	TenantID uint   `json:"tenantID"`
}

// EduTeacherRoleConfigResponse 教师角色配置响应
type EduTeacherRoleConfigResponse struct {
	Roles           []EduTeacherRoleConfigRoleItem `json:"roles"`
	SelectedRoleIDs []uint                         `json:"selectedRoleIds"`
}
