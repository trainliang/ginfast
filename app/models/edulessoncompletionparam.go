package models

import "github.com/gin-gonic/gin"

type EduLessonCompletionRuleItemRequest struct {
	ResultType          string `form:"resultType" json:"resultType" validate:"required" message:"结课结果不能为空"`
	DeductEnabled       int8   `form:"deductEnabled" json:"deductEnabled"`
	DeductMode          string `form:"deductMode" json:"deductMode"`
	FixedCount          int    `form:"fixedCount" json:"fixedCount"`
	DurationUnitMinutes int    `form:"durationUnitMinutes" json:"durationUnitMinutes"`
	InsufficientPolicy  string `form:"insufficientPolicy" json:"insufficientPolicy"`
	Status              int8   `form:"status" json:"status"`
	Remark              string `form:"remark" json:"remark"`
}

type EduLessonCompletionRuleSaveRequest struct {
	Validator
	Items []EduLessonCompletionRuleItemRequest `form:"items" json:"items" validate:"required" message:"消课规则不能为空"`
}

func (r *EduLessonCompletionRuleSaveRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

type EduLessonCompletionListRequest struct {
	Validator
	LessonID FlexString `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
}

func (r *EduLessonCompletionListRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduLessonCompletionListRequest) GetLessonIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.LessonID))
}

type EduLessonCompletionSubmitItemRequest struct {
	StudentID  FlexString `form:"studentId" json:"studentId" validate:"required" message:"学生ID不能为空"`
	ResultType string     `form:"resultType" json:"resultType" validate:"required" message:"结课结果不能为空"`
	Reason     string     `form:"reason" json:"reason"`
}

func (r *EduLessonCompletionSubmitItemRequest) GetStudentIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.StudentID))
}

type EduLessonCompletionSubmitRequest struct {
	Validator
	LessonID FlexString                             `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
	Items    []EduLessonCompletionSubmitItemRequest `form:"items" json:"items" validate:"required" message:"消课学生不能为空"`
}

func (r *EduLessonCompletionSubmitRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduLessonCompletionSubmitRequest) GetLessonIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.LessonID))
}

type EduLessonCompletionRevokeRequest struct {
	Validator
	LessonID  FlexString `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
	StudentID FlexString `form:"studentId" json:"studentId" validate:"required" message:"学生ID不能为空"`
	Reason    string     `form:"reason" json:"reason" validate:"required" message:"撤销原因不能为空"`
}

func (r *EduLessonCompletionRevokeRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduLessonCompletionRevokeRequest) GetLessonIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.LessonID))
}

func (r *EduLessonCompletionRevokeRequest) GetStudentIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.StudentID))
}

type EduLessonCompleteRequest struct {
	Validator
	LessonID FlexString `form:"lessonId" json:"lessonId" validate:"required" message:"课次ID不能为空"`
	Reason   string     `form:"reason" json:"reason"`
}

func (r *EduLessonCompleteRequest) Validate(c *gin.Context) error {
	return r.Validator.Check(c, r)
}

func (r *EduLessonCompleteRequest) GetLessonIDUint() uint {
	if r == nil {
		return 0
	}
	return StringToUint(string(r.LessonID))
}

type EduLessonCompletionStudentStatus struct {
	StudentID  uint                        `json:"studentId"`
	Completion *EduLessonStudentCompletion `json:"completion,omitempty"`
	Processed  bool                        `json:"processed"`
}

type EduLessonCompletionStatusResponse struct {
	Lesson           *EduLesson                         `json:"lesson"`
	Students         []EduLessonCompletionStudentStatus `json:"students"`
	UnprocessedCount int                                `json:"unprocessedCount"`
}
