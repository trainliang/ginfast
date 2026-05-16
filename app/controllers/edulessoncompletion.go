package controllers

import (
	"gin-fast/app/models"
	"gin-fast/app/service"

	"github.com/gin-gonic/gin"
)

// EduLessonCompletionController 消课控制器
type EduLessonCompletionController struct {
	Common
	EduLessonCompletionService *service.EduLessonCompletionService
}

func NewEduLessonCompletionController() *EduLessonCompletionController {
	return &EduLessonCompletionController{
		Common:                     Common{},
		EduLessonCompletionService: service.NewEduLessonCompletionService(),
	}
}

func (ctl *EduLessonCompletionController) Rules(c *gin.Context) {
	rows, err := ctl.EduLessonCompletionService.ListRules(c, ctl.RequireTenant(c))
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, gin.H{"list": rows})
}

func (ctl *EduLessonCompletionController) SaveRules(c *gin.Context) {
	var req models.EduLessonCompletionRuleSaveRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduLessonCompletionService.SaveRules(c, ctl.RequireTenant(c), &req); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "消课规则保存成功", nil)
}

func (ctl *EduLessonCompletionController) Completions(c *gin.Context) {
	req := models.EduLessonCompletionListRequest{LessonID: models.FlexString(c.Param("id"))}
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	req.LessonID = models.FlexString(c.Param("id"))
	result, err := ctl.EduLessonCompletionService.GetLessonCompletionStatus(c, ctl.RequireTenant(c), req.GetLessonIDUint())
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, result)
}

func (ctl *EduLessonCompletionController) SubmitCompletions(c *gin.Context) {
	var req models.EduLessonCompletionSubmitRequest
	req.LessonID = models.FlexString(c.Param("id"))
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	req.LessonID = models.FlexString(c.Param("id"))
	rows, err := ctl.EduLessonCompletionService.SubmitCompletions(c, ctl.RequireTenant(c), &req)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "消课提交成功", gin.H{"list": rows})
}

func (ctl *EduLessonCompletionController) RevokeCompletion(c *gin.Context) {
	var req models.EduLessonCompletionRevokeRequest
	req.LessonID = models.FlexString(c.Param("id"))
	req.StudentID = models.FlexString(c.Param("studentId"))
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	req.LessonID = models.FlexString(c.Param("id"))
	req.StudentID = models.FlexString(c.Param("studentId"))
	if err := ctl.EduLessonCompletionService.RevokeCompletion(c, ctl.RequireTenant(c), &req); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学生消课已撤销", nil)
}

func (ctl *EduLessonCompletionController) CompleteLesson(c *gin.Context) {
	var req models.EduLessonCompleteRequest
	req.LessonID = models.FlexString(c.Param("id"))
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	req.LessonID = models.FlexString(c.Param("id"))
	if err := ctl.EduLessonCompletionService.CompleteLesson(c, ctl.RequireTenant(c), &req); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "课次已完成", nil)
}
