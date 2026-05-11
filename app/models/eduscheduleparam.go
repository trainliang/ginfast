package models

import "github.com/gin-gonic/gin"

type EduScheduleRuleListRequest struct {
	BasePaging
	Validator
	ID         *uint  `form:"id" json:"id"`
	RuleType   string `form:"ruleType" json:"ruleType"`
	RepeatType string `form:"repeatType" json:"repeatType"`
	ClassID    *uint  `form:"classId" json:"classId"`
	StudentID  *uint  `form:"studentId" json:"studentId"`
	CourseID   *uint  `form:"courseId" json:"courseId"`
	TeacherID  *uint  `form:"teacherId" json:"teacherId"`
	Status     *int8  `form:"status" json:"status"`
}

func (r *EduScheduleRuleListRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduScheduleRuleAddRequest struct {
	Validator
	RuleType      string    `form:"ruleType" json:"ruleType" validate:"required" message:"规则类型不能为空"`
	RepeatType    string    `form:"repeatType" json:"repeatType" validate:"required" message:"重复类型不能为空"`
	TermID        uint      `form:"termId" json:"termId"`
	StartDate     *JSONTime `form:"startDate" json:"startDate"`
	EndDate       *JSONTime `form:"endDate" json:"endDate"`
	ClassID       uint      `form:"classId" json:"classId"`
	StudentID     uint      `form:"studentId" json:"studentId"`
	CourseID      uint      `form:"courseId" json:"courseId" validate:"required" message:"课程ID不能为空"`
	TeacherID     uint      `form:"teacherId" json:"teacherId" validate:"required" message:"教师ID不能为空"`
	TeachingMode  string    `form:"teachingMode" json:"teachingMode"`
	RequiresRoom  int8      `form:"requiresRoom" json:"requiresRoom"`
	RoomID        uint      `form:"roomId" json:"roomId"`
	Weekday       int8      `form:"weekday" json:"weekday"`
	StartTime     string    `form:"startTime" json:"startTime" validate:"required" message:"开始时间不能为空"`
	EndTime       string    `form:"endTime" json:"endTime" validate:"required" message:"结束时间不能为空"`
	Status        *int8     `form:"status" json:"status"`
	Version       *int      `form:"version" json:"version"`
	EffectiveFrom *JSONTime `form:"effectiveFrom" json:"effectiveFrom"`
}

func (r *EduScheduleRuleAddRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduScheduleRuleUpdateRequest struct {
	EduScheduleRuleAddRequest
	ID uint `form:"id" json:"id" validate:"required" message:"规则ID不能为空"`
}

type EduScheduleRuleDeleteRequest struct {
	Validator
	ID uint `form:"id" json:"id" validate:"required" message:"规则ID不能为空"`
}

func (r *EduScheduleRuleDeleteRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduScheduleRulePreviewChangeRequest struct {
	Validator
	ID uint `form:"id" json:"id" validate:"required" message:"规则ID不能为空"`
}

func (r *EduScheduleRulePreviewChangeRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduLessonListRequest struct {
	BasePaging
	Validator
	ID        *uint  `form:"id" json:"id"`
	RuleID    *uint  `form:"ruleId" json:"ruleId"`
	ClassID   *uint  `form:"classId" json:"classId"`
	StudentID *uint  `form:"studentId" json:"studentId"`
	CourseID  *uint  `form:"courseId" json:"courseId"`
	TeacherID *uint  `form:"teacherId" json:"teacherId"`
	Status    string `form:"status" json:"status"`
}

func (r *EduLessonListRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduLessonCalendarRequest struct {
	Validator
	StartDate *JSONTime `form:"startDate" json:"startDate"`
	EndDate   *JSONTime `form:"endDate" json:"endDate"`
	TeacherID *uint     `form:"teacherId" json:"teacherId"`
	ClassID   *uint     `form:"classId" json:"classId"`
}

func (r *EduLessonCalendarRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduScheduleConflictCheckRequest struct {
	Validator
	LessonID uint `form:"lessonId" json:"lessonId"`
	RuleID   uint `form:"ruleId" json:"ruleId"`
}

func (r *EduScheduleConflictCheckRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduScheduleConflictItem struct {
	ConflictType string `json:"conflictType"`
	ConflictKey  string `json:"conflictKey"`
	Reason       string `json:"reason"`
}

type EduScheduleConflictCheckResult struct {
	HasConflict bool                      `json:"hasConflict"`
	Items       []EduScheduleConflictItem `json:"items"`
}
