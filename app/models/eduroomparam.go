package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduRoomListRequest 场地列表请求
type EduRoomListRequest struct {
	BasePaging
	Validator
	ID       *FlexUint `form:"id" json:"id"`
	Name     string `form:"name" json:"name"`
	Code     string `form:"code" json:"code"`
	Type     string `form:"type" json:"type"`
	Status   *FlexInt8 `form:"status" json:"status"`
	Location string `form:"location" json:"location"`
}

func (r *EduRoomListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduRoomListRequest) Handler() func(db *gorm.DB) *gorm.DB {
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
		if r.Location != "" {
			db = db.Where("location LIKE ?", "%"+r.Location+"%")
		}
		return db
	}
}

// EduRoomAddRequest 新增场地请求
type EduRoomAddRequest struct {
	Validator
	Name     string `form:"name" json:"name" validate:"required" message:"场地名称不能为空"`
	Code     string `form:"code" json:"code" validate:"required" message:"场地编码不能为空"`
	Type     string `form:"type" json:"type"`
	Capacity FlexInt `form:"capacity" json:"capacity" validate:"required" message:"容量不能为空"`
	Location string `form:"location" json:"location"`
	Status   *FlexInt8 `form:"status" json:"status"`
	Remark   string `form:"remark" json:"remark"`
}

func (r *EduRoomAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomUpdateRequest 更新场地请求
type EduRoomUpdateRequest struct {
	Validator
	ID       FlexUint `form:"id" json:"id" validate:"required" message:"场地ID不能为空"`
	Name     string `form:"name" json:"name" validate:"required" message:"场地名称不能为空"`
	Code     string `form:"code" json:"code" validate:"required" message:"场地编码不能为空"`
	Type     string `form:"type" json:"type"`
	Capacity FlexInt `form:"capacity" json:"capacity" validate:"required" message:"容量不能为空"`
	Location string `form:"location" json:"location"`
	Status   *FlexInt8 `form:"status" json:"status"`
	Remark   string `form:"remark" json:"remark"`
}

func (r *EduRoomUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomDeleteRequest 删除场地请求
type EduRoomDeleteRequest struct {
	Validator
	ID FlexUint `form:"id" json:"id" validate:"required" message:"场地ID不能为空"`
}

func (r *EduRoomDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomGetRequest 获取场地请求
type EduRoomGetRequest struct {
	Validator
	ID FlexUint `form:"id" uri:"id" json:"id" validate:"required" message:"场地ID不能为空"`
}

func (r *EduRoomGetRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomWeeklyRuleRequest 场地周规则请求
type EduRoomWeeklyRuleRequest struct {
	ID        FlexUint `json:"id" form:"id"`
	RoomID    FlexUint `json:"roomId" form:"roomId"`
	Weekday   FlexInt8 `json:"weekday" form:"weekday"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
	Available FlexInt8 `json:"available" form:"available"`
	Remark    string `json:"remark" form:"remark"`
}

// EduRoomWeeklyRuleListRequest 场地周规则列表请求
type EduRoomWeeklyRuleListRequest struct {
	BasePaging
	Validator
	ID      *FlexUint `form:"id" json:"id"`
	RoomID  *FlexUint `form:"roomId" json:"roomId"`
	Weekday *FlexInt8 `form:"weekday" json:"weekday"`
}

func (r *EduRoomWeeklyRuleListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduRoomWeeklyRuleListRequest) Handler() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID != nil {
			db = db.Where("id = ?", *r.ID)
		}
		if r.RoomID != nil {
			db = db.Where("room_id = ?", *r.RoomID)
		}
		if r.Weekday != nil {
			db = db.Where("weekday = ?", *r.Weekday)
		}
		return db
	}
}

// EduRoomWeeklyRuleAddRequest 新增场地周规则请求
type EduRoomWeeklyRuleAddRequest struct {
	Validator
	RoomID    FlexUint `json:"roomId" form:"roomId" validate:"required" message:"场地ID不能为空"`
	Weekday   FlexInt8 `json:"weekday" form:"weekday" validate:"required" message:"星期不能为空"`
	StartTime string `json:"startTime" form:"startTime" validate:"required" message:"开始时间不能为空"`
	EndTime   string `json:"endTime" form:"endTime" validate:"required" message:"结束时间不能为空"`
	Available FlexInt8 `json:"available" form:"available"`
	Remark    string `json:"remark" form:"remark"`
}

func (r *EduRoomWeeklyRuleAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomWeeklyRuleUpdateRequest 更新场地周规则请求
type EduRoomWeeklyRuleUpdateRequest struct {
	Validator
	ID        FlexUint `json:"id" form:"id" validate:"required" message:"规则ID不能为空"`
	RoomID    FlexUint `json:"roomId" form:"roomId" validate:"required" message:"场地ID不能为空"`
	Weekday   FlexInt8 `json:"weekday" form:"weekday" validate:"required" message:"星期不能为空"`
	StartTime string `json:"startTime" form:"startTime" validate:"required" message:"开始时间不能为空"`
	EndTime   string `json:"endTime" form:"endTime" validate:"required" message:"结束时间不能为空"`
	Available FlexInt8 `json:"available" form:"available"`
	Remark    string `json:"remark" form:"remark"`
}

func (r *EduRoomWeeklyRuleUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomWeeklyRuleDeleteRequest 删除场地周规则请求
type EduRoomWeeklyRuleDeleteRequest struct {
	Validator
	ID FlexUint `json:"id" form:"id" validate:"required" message:"规则ID不能为空"`
}

func (r *EduRoomWeeklyRuleDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomExceptionRequest 场地例外请求
type EduRoomExceptionRequest struct {
	ID            FlexUint  `json:"id" form:"id"`
	RoomID        FlexUint  `json:"roomId" form:"roomId"`
	ExceptionDate *JSONTime `json:"exceptionDate" form:"exceptionDate"`
	Type          string    `json:"type" form:"type"`
	StartTime     string    `json:"startTime" form:"startTime"`
	EndTime       string    `json:"endTime" form:"endTime"`
	Reason        string    `json:"reason" form:"reason"`
}

// EduRoomExceptionListRequest 场地例外列表请求
type EduRoomExceptionListRequest struct {
	BasePaging
	Validator
	ID     *FlexUint `form:"id" json:"id"`
	RoomID *FlexUint `form:"roomId" json:"roomId"`
	Type   string `form:"type" json:"type"`
}

func (r *EduRoomExceptionListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduRoomExceptionListRequest) Handler() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ID != nil {
			db = db.Where("id = ?", *r.ID)
		}
		if r.RoomID != nil {
			db = db.Where("room_id = ?", *r.RoomID)
		}
		if r.Type != "" {
			db = db.Where("type = ?", r.Type)
		}
		return db
	}
}

// EduRoomExceptionAddRequest 新增场地例外请求
type EduRoomExceptionAddRequest struct {
	Validator
	RoomID        FlexUint  `json:"roomId" form:"roomId" validate:"required" message:"场地ID不能为空"`
	ExceptionDate *JSONTime `json:"exceptionDate" form:"exceptionDate" validate:"required" message:"例外日期不能为空"`
	Type          string    `json:"type" form:"type" validate:"required" message:"例外类型不能为空"`
	StartTime     string    `json:"startTime" form:"startTime"`
	EndTime       string    `json:"endTime" form:"endTime"`
	Reason        string    `json:"reason" form:"reason"`
}

func (r *EduRoomExceptionAddRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomExceptionUpdateRequest 更新场地例外请求
type EduRoomExceptionUpdateRequest struct {
	Validator
	ID            FlexUint  `json:"id" form:"id" validate:"required" message:"例外ID不能为空"`
	RoomID        FlexUint  `json:"roomId" form:"roomId" validate:"required" message:"场地ID不能为空"`
	ExceptionDate *JSONTime `json:"exceptionDate" form:"exceptionDate" validate:"required" message:"例外日期不能为空"`
	Type          string    `json:"type" form:"type" validate:"required" message:"例外类型不能为空"`
	StartTime     string    `json:"startTime" form:"startTime"`
	EndTime       string    `json:"endTime" form:"endTime"`
	Reason        string    `json:"reason" form:"reason"`
}

func (r *EduRoomExceptionUpdateRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomExceptionDeleteRequest 删除场地例外请求
type EduRoomExceptionDeleteRequest struct {
	Validator
	ID FlexUint `json:"id" form:"id" validate:"required" message:"例外ID不能为空"`
}

func (r *EduRoomExceptionDeleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomImportRow 场地导入行
type EduRoomImportRow struct {
	Name     string `json:"name" form:"name"`
	Code     string `json:"code" form:"code"`
	Type     string `json:"type" form:"type"`
	Capacity FlexInt `json:"capacity" form:"capacity"`
	Location string `json:"location" form:"location"`
	Status   *FlexInt8 `json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
}

// EduRoomWeeklyRuleImportRow 场地周规则导入行
type EduRoomWeeklyRuleImportRow struct {
	RoomID    FlexUint `json:"roomId" form:"roomId"`
	Weekday   FlexInt8 `json:"weekday" form:"weekday"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
	Available FlexInt8 `json:"available" form:"available"`
	Remark    string `json:"remark" form:"remark"`
}

// EduRoomExceptionImportRow 场地例外导入行
type EduRoomExceptionImportRow struct {
	RoomID        FlexUint  `json:"roomId" form:"roomId"`
	ExceptionDate *JSONTime `json:"exceptionDate" form:"exceptionDate"`
	Type          string    `json:"type" form:"type"`
	StartTime     string    `json:"startTime" form:"startTime"`
	EndTime       string    `json:"endTime" form:"endTime"`
	Reason        string    `json:"reason" form:"reason"`
}

// EduRoomExportRow 场地导出行
type EduRoomExportRow struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Type      string `json:"type"`
	Capacity  int    `json:"capacity"`
	Location  string `json:"location"`
	Status    int8   `json:"status"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// EduRoomWeeklyRuleExportRow 场地周规则导出行
type EduRoomWeeklyRuleExportRow struct {
	ID        uint   `json:"id"`
	RoomID    uint   `json:"roomId"`
	Weekday   int8   `json:"weekday"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Available int8   `json:"available"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// EduRoomExceptionExportRow 场地例外导出行
type EduRoomExceptionExportRow struct {
	ID            uint      `json:"id"`
	RoomID        uint      `json:"roomId"`
	ExceptionDate *JSONTime `json:"exceptionDate"`
	Type          string    `json:"type"`
	StartTime     string    `json:"startTime"`
	EndTime       string    `json:"endTime"`
	Reason        string    `json:"reason"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
}

// EduRoomImportRequest 场地导入请求
type EduRoomImportRequest struct {
	Validator
	Rows []EduRoomImportRow `json:"rows" form:"rows" validate:"required" message:"导入数据不能为空"`
}

func (r *EduRoomImportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomWeeklyRuleImportRequest 场地周规则导入请求
type EduRoomWeeklyRuleImportRequest struct {
	Validator
	Rows []EduRoomWeeklyRuleImportRow `json:"rows" form:"rows" validate:"required" message:"导入数据不能为空"`
}

func (r *EduRoomWeeklyRuleImportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomExceptionImportRequest 场地例外导入请求
type EduRoomExceptionImportRequest struct {
	Validator
	Rows []EduRoomExceptionImportRow `json:"rows" form:"rows" validate:"required" message:"导入数据不能为空"`
}

func (r *EduRoomExceptionImportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduRoomExportRequest 场地导出请求
type EduRoomExportRequest struct {
	Validator
	IDs []FlexUint `form:"ids" json:"ids"`
}

func (r *EduRoomExportRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

// EduTeacherOption 教师下拉选项
type EduTeacherOption struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	NickName string `json:"nickName"`
	Name     string `json:"name"`
}
