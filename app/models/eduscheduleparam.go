package models

import "github.com/gin-gonic/gin"

type EduScheduleRuleListRequest struct {
	BasePaging
	Validator
	ID         FlexString `form:"id" json:"id"`
	Name       string     `form:"name" json:"name"`
	RuleType   string     `form:"ruleType" json:"ruleType"`
	RepeatType string     `form:"repeatType" json:"repeatType"`
	ClassID    FlexString `form:"classId" json:"classId"`
	StudentID  FlexString `form:"studentId" json:"studentId"`
	CourseID   FlexString `form:"courseId" json:"courseId"`
	TeacherID  FlexString `form:"teacherId" json:"teacherId"`
	Status     *FlexInt8  `form:"status" json:"status"`
}

func (r *EduScheduleRuleListRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduScheduleRuleAddRequest struct {
	Validator
	Name          string     `form:"name" json:"name" validate:"required" message:"规则名称不能为空"`
	RuleType      string     `form:"ruleType" json:"ruleType" validate:"required" message:"规则类型不能为空"`
	RepeatType    string     `form:"repeatType" json:"repeatType" validate:"required" message:"重复类型不能为空"`
	TermID        FlexString `form:"termId" json:"termId"`
	StartDate     *JSONTime  `form:"startDate" json:"startDate"`
	EndDate       *JSONTime  `form:"endDate" json:"endDate"`
	ClassID       FlexString `form:"classId" json:"classId"`
	StudentID     FlexString `form:"studentId" json:"studentId"`
	CourseID      FlexString `form:"courseId" json:"courseId" validate:"required" message:"课程ID不能为空"`
	TeacherID     FlexString `form:"teacherId" json:"teacherId" validate:"required" message:"教师ID不能为空"`
	TeachingMode  string     `form:"teachingMode" json:"teachingMode"`
	RequiresRoom  int8       `form:"requiresRoom" json:"requiresRoom"`
	RoomID        FlexString `form:"roomId" json:"roomId"`
	Weekdays      []int8     `form:"weekdays" json:"weekdays"`
	StartTime     string     `form:"startTime" json:"startTime" validate:"required" message:"开始时间不能为空"`
	EndTime       string     `form:"endTime" json:"endTime" validate:"required" message:"结束时间不能为空"`
	Status        *FlexInt8  `form:"status" json:"status"`
	Version       *FlexInt   `form:"version" json:"version"`
	EffectiveFrom *JSONTime  `form:"effectiveFrom" json:"effectiveFrom"`
	Remark        string     `form:"remark" json:"remark"`
}

func (r *EduScheduleRuleAddRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduScheduleRuleUpdateRequest struct {
	EduScheduleRuleAddRequest
	ID FlexString `form:"id" json:"id" validate:"required" message:"规则ID不能为空"`
}

type EduScheduleRuleDeleteRequest struct {
	Validator
	ID FlexString `form:"id" json:"id" validate:"required" message:"规则ID不能为空"`
}

func (r *EduScheduleRuleDeleteRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduScheduleRulePreviewChangeRequest struct {
	Validator
	ID FlexString `form:"id" json:"id" validate:"required" message:"规则ID不能为空"`
}

func (r *EduScheduleRulePreviewChangeRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduLessonListRequest struct {
	BasePaging
	Validator
	ID        FlexString `form:"id" json:"id"`
	RuleID    FlexString `form:"ruleId" json:"ruleId"`
	ClassID   FlexString `form:"classId" json:"classId"`
	StudentID FlexString `form:"studentId" json:"studentId"`
	CourseID  FlexString `form:"courseId" json:"courseId"`
	TeacherID FlexString `form:"teacherId" json:"teacherId"`
	Status    string     `form:"status" json:"status"`
}

func (r *EduLessonListRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduLessonCalendarRequest struct {
	Validator
	StartDate *JSONTime  `form:"startDate" json:"startDate"`
	EndDate   *JSONTime  `form:"endDate" json:"endDate"`
	TeacherID FlexString `form:"teacherId" json:"teacherId"`
	ClassID   FlexString `form:"classId" json:"classId"`
}

func (r *EduLessonCalendarRequest) Validate(c *gin.Context) error { return r.Validator.Check(c, r) }

type EduLessonRescheduleRequest struct {
	Validator
	LessonID              FlexString `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
	LessonDate            *JSONTime  `form:"lessonDate" json:"lessonDate" validate:"required" message:"新日期不能为空"`
	StartTime             string     `form:"startTime" json:"startTime" validate:"required" message:"新开始时间不能为空"`
	EndTime               string     `form:"endTime" json:"endTime" validate:"required" message:"新结束时间不能为空"`
	TeacherID             FlexString `form:"teacherId" json:"teacherId" validate:"required" message:"教师ID不能为空"`
	TeachingMode          string     `form:"teachingMode" json:"teachingMode" validate:"required" message:"授课方式不能为空"`
	RoomID                FlexString `form:"roomId" json:"roomId"`
	AllowConflictOverride bool       `form:"allowConflictOverride" json:"allowConflictOverride"`
	OverrideReason        string     `form:"overrideReason" json:"overrideReason"`
}

func (r *EduLessonRescheduleRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduLessonStopRequest struct {
	Validator
	LessonID FlexString `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
	Reason   string     `form:"reason" json:"reason" validate:"required" message:"停课原因不能为空"`
}

func (r *EduLessonStopRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduLessonCancelRequest struct {
	Validator
	LessonID FlexString `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
	Reason   string     `form:"reason" json:"reason" validate:"required" message:"取消原因不能为空"`
}

func (r *EduLessonCancelRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduLessonRestoreRequest struct {
	Validator
	LessonID FlexString `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
	Reason   string     `form:"reason" json:"reason" validate:"required" message:"恢复原因不能为空"`
}

func (r *EduLessonRestoreRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduLessonMakeupRequest struct {
	Validator
	LessonID              FlexString `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
	LessonDate            *JSONTime  `form:"lessonDate" json:"lessonDate" validate:"required" message:"新日期不能为空"`
	StartTime             string     `form:"startTime" json:"startTime" validate:"required" message:"新开始时间不能为空"`
	EndTime               string     `form:"endTime" json:"endTime" validate:"required" message:"新结束时间不能为空"`
	TeacherID             FlexString `form:"teacherId" json:"teacherId" validate:"required" message:"教师ID不能为空"`
	TeachingMode          string     `form:"teachingMode" json:"teachingMode" validate:"required" message:"授课方式不能为空"`
	RoomID                FlexString `form:"roomId" json:"roomId"`
	Reason                string     `form:"reason" json:"reason" validate:"required" message:"补课原因不能为空"`
	AllowConflictOverride bool       `form:"allowConflictOverride" json:"allowConflictOverride"`
	OverrideReason        string     `form:"overrideReason" json:"overrideReason"`
}

func (r *EduLessonMakeupRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduLessonChangeLogListRequest struct {
	BasePaging
	Validator
	LessonID FlexString `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
}

func (r *EduLessonChangeLogListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduScheduleConflictCheckRequest struct {
	Validator
	LessonID              string     `form:"lessonId" json:"lessonId"`
	RuleID                FlexString `form:"ruleId" json:"ruleId"`
	LessonDate            *JSONTime  `form:"lessonDate" json:"lessonDate"`
	StartTime             string     `form:"startTime" json:"startTime"`
	EndTime               string     `form:"endTime" json:"endTime"`
	ClassID               FlexString `form:"classId" json:"classId"`
	StudentID             FlexString `form:"studentId" json:"studentId"`
	CourseID              FlexString `form:"courseId" json:"courseId"`
	TeacherID             FlexString `form:"teacherId" json:"teacherId"`
	TeachingMode          string     `form:"teachingMode" json:"teachingMode"`
	RequiresRoom          int8       `form:"requiresRoom" json:"requiresRoom"`
	RoomID                FlexString `form:"roomId" json:"roomId"`
	AllowConflictOverride bool       `form:"allowConflictOverride" json:"allowConflictOverride"`
	OverrideReason        string     `form:"overrideReason" json:"overrideReason"`
	OperatorID            FlexString `form:"operatorID" json:"operatorID"`
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

func (r *EduScheduleRuleListRequest) GetIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.ID))
}

func (r *EduScheduleRuleListRequest) GetClassIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.ClassID))
}

func (r *EduScheduleRuleListRequest) GetStudentIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.StudentID))
}

func (r *EduScheduleRuleListRequest) GetCourseIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.CourseID))
}

func (r *EduScheduleRuleListRequest) GetTeacherIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.TeacherID))
}

func (r *EduScheduleRuleAddRequest) GetTermIDUint() uint        { return StringToUint(string(r.TermID)) }
func (r *EduScheduleRuleAddRequest) GetClassIDUint() uint       { return StringToUint(string(r.ClassID)) }
func (r *EduScheduleRuleAddRequest) GetStudentIDUint() uint     { return StringToUint(string(r.StudentID)) }
func (r *EduScheduleRuleAddRequest) GetCourseIDUint() uint      { return StringToUint(string(r.CourseID)) }
func (r *EduScheduleRuleAddRequest) GetTeacherIDUint() uint     { return StringToUint(string(r.TeacherID)) }
func (r *EduScheduleRuleAddRequest) GetRoomIDUint() uint        { return StringToUint(string(r.RoomID)) }
func (r *EduScheduleRuleAddRequest) GetStatusInt8() int8        { if r == nil || r.Status == nil { return 0 }; return int8(*r.Status) }
func (r *EduScheduleRuleAddRequest) GetVersionInt() int         { if r == nil || r.Version == nil { return 0 }; return int(*r.Version) }
func (r *EduScheduleRuleUpdateRequest) GetIDUint() uint         { return StringToUint(string(r.ID)) }
func (r *EduScheduleRuleDeleteRequest) GetIDUint() uint         { return StringToUint(string(r.ID)) }
func (r *EduScheduleRulePreviewChangeRequest) GetIDUint() uint  { return StringToUint(string(r.ID)) }

func (r *EduLessonListRequest) GetIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.ID))
}

func (r *EduLessonListRequest) GetRuleIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.RuleID))
}

func (r *EduLessonListRequest) GetClassIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.ClassID))
}

func (r *EduLessonListRequest) GetStudentIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.StudentID))
}

func (r *EduLessonListRequest) GetCourseIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.CourseID))
}

func (r *EduLessonListRequest) GetTeacherIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.TeacherID))
}

func (r *EduLessonCalendarRequest) GetTeacherID() string {
	if r == nil {
		return ""
	}
	return string(r.TeacherID)
}

func (r *EduLessonCalendarRequest) GetTeacherIDUint() uint { return StringToUint(r.GetTeacherID()) }

func (r *EduLessonCalendarRequest) GetClassID() string {
	if r == nil {
		return ""
	}
	return string(r.ClassID)
}

func (r *EduLessonCalendarRequest) GetClassIDUint() uint { return StringToUint(r.GetClassID()) }

func (r *EduLessonRescheduleRequest) GetLessonIDUint() uint { return StringToUint(string(r.LessonID)) }
func (r *EduLessonRescheduleRequest) GetTeacherIDUint() uint { return StringToUint(string(r.TeacherID)) }
func (r *EduLessonRescheduleRequest) GetRoomID() string {
	if r == nil {
		return ""
	}
	return string(r.RoomID)
}

func (r *EduLessonRescheduleRequest) GetRoomIDUint() uint { return StringToUint(r.GetRoomID()) }

func (r *EduLessonStopRequest) GetLessonIDUint() uint    { return StringToUint(string(r.LessonID)) }
func (r *EduLessonCancelRequest) GetLessonIDUint() uint  { return StringToUint(string(r.LessonID)) }
func (r *EduLessonRestoreRequest) GetLessonIDUint() uint { return StringToUint(string(r.LessonID)) }

func (r *EduLessonMakeupRequest) GetLessonIDUint() uint  { return StringToUint(string(r.LessonID)) }
func (r *EduLessonMakeupRequest) GetTeacherIDUint() uint { return StringToUint(string(r.TeacherID)) }
func (r *EduLessonMakeupRequest) GetRoomID() string {
	if r == nil {
		return ""
	}
	return string(r.RoomID)
}

func (r *EduLessonMakeupRequest) GetRoomIDUint() uint { return StringToUint(r.GetRoomID()) }

func (r *EduLessonChangeLogListRequest) GetLessonIDUint() uint { return StringToUint(string(r.LessonID)) }

func (r *EduScheduleConflictCheckRequest) GetLessonID() string      { return r.LessonID }
func (r *EduScheduleConflictCheckRequest) GetLessonIDUint() uint    { return StringToUint(r.LessonID) }
func (r *EduScheduleConflictCheckRequest) GetRuleID() string        { return string(r.RuleID) }
func (r *EduScheduleConflictCheckRequest) GetRuleIDUint() uint      { return StringToUint(string(r.RuleID)) }
func (r *EduScheduleConflictCheckRequest) GetClassID() string       { return string(r.ClassID) }
func (r *EduScheduleConflictCheckRequest) GetClassIDUint() uint     { return StringToUint(string(r.ClassID)) }
func (r *EduScheduleConflictCheckRequest) GetStudentID() string     { return string(r.StudentID) }
func (r *EduScheduleConflictCheckRequest) GetStudentIDUint() uint   { return StringToUint(string(r.StudentID)) }
func (r *EduScheduleConflictCheckRequest) GetCourseID() string      { return string(r.CourseID) }
func (r *EduScheduleConflictCheckRequest) GetCourseIDUint() uint    { return StringToUint(string(r.CourseID)) }
func (r *EduScheduleConflictCheckRequest) GetTeacherID() string     { return string(r.TeacherID) }
func (r *EduScheduleConflictCheckRequest) GetTeacherIDUint() uint   { return StringToUint(string(r.TeacherID)) }
func (r *EduScheduleConflictCheckRequest) GetRoomID() string        { return string(r.RoomID) }
func (r *EduScheduleConflictCheckRequest) GetRoomIDUint() uint      { return StringToUint(string(r.RoomID)) }
func (r *EduScheduleConflictCheckRequest) GetOperatorIDUint() uint  { return StringToUint(string(r.OperatorID)) }
