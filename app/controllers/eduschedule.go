package controllers

import (
	"gin-fast/app/models"
	"gin-fast/app/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduScheduleController 排课控制器
type EduScheduleController struct {
	Common
	EduScheduleService *service.EduScheduleService
}

func NewEduScheduleController() *EduScheduleController {
	return &EduScheduleController{Common: Common{}, EduScheduleService: service.NewEduScheduleService()}
}

func (ctl *EduScheduleController) RuleList(c *gin.Context) {
	var req models.EduScheduleRuleListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := buildEduScheduleRuleListScope(&req, tenantID)
	total, err := models.NewEduScheduleRuleList().GetTotal(c, scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计排课规则数量失败", err)
	}
	list := models.NewEduScheduleRuleList()
	if err := list.Find(c, req.Paginate(), scope); err != nil {
		ctl.FailAndAbort(c, "获取排课规则列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduScheduleController) RuleAdd(c *gin.Context) {
	var req models.EduScheduleRuleAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	rule := buildEduScheduleRuleFromAddRequest(&req, ctl.RequireTenant(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduScheduleService.CreateRule(c, rule); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "排课规则创建成功", rule)
}

func (ctl *EduScheduleController) RuleUpdate(c *gin.Context) {
	var req models.EduScheduleRuleUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	rule := buildEduScheduleRuleFromUpdateRequest(&req, ctl.RequireTenant(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduScheduleService.UpdateRule(c, rule); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "排课规则更新成功", rule)
}

func (ctl *EduScheduleController) RulePreviewChange(c *gin.Context) {
	var req models.EduScheduleRulePreviewChangeRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	rows, err := ctl.EduScheduleService.PreviewRuleChange(c, ctl.RequireTenant(c), req.ID)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, gin.H{"list": rows})
}

func (ctl *EduScheduleController) RuleDelete(c *gin.Context) {
	var req models.EduScheduleRuleDeleteRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduScheduleService.DeleteRule(c, ctl.RequireTenant(c), req.ID); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "排课规则删除成功", nil)
}

func (ctl *EduScheduleController) LessonCalendar(c *gin.Context) {
	var req models.EduLessonCalendarRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	rows, err := ctl.EduScheduleService.CalendarLessons(c, ctl.RequireTenant(c), &req)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, gin.H{"list": rows})
}

func (ctl *EduScheduleController) LessonList(c *gin.Context) {
	var req models.EduLessonListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := buildEduScheduleLessonListScope(&req, tenantID)
	total, err := models.NewEduLessonList().GetTotal(c, scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计课次数量失败", err)
	}
	list := models.NewEduLessonList()
	if err := list.Find(c, req.Paginate(), scope); err != nil {
		ctl.FailAndAbort(c, "获取课次列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduScheduleController) CheckConflicts(c *gin.Context) {
	var req models.EduScheduleConflictCheckRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	result, err := ctl.EduScheduleService.CheckConflicts(c, &req)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, result)
}

func buildEduScheduleRuleFromAddRequest(req *models.EduScheduleRuleAddRequest, tenantID, userID uint) *models.EduScheduleRule {
	rule := &models.EduScheduleRule{
		RuleType:      req.RuleType,
		RepeatType:    req.RepeatType,
		TermID:        req.TermID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		ClassID:       req.ClassID,
		StudentID:     req.StudentID,
		CourseID:      req.CourseID,
		TeacherID:     req.TeacherID,
		TeachingMode:  req.TeachingMode,
		RequiresRoom:  req.RequiresRoom,
		RoomID:        req.RoomID,
		Weekday:       req.Weekday,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Status:        1,
		Version:       1,
		EffectiveFrom: req.EffectiveFrom,
		CreatedBy:     userID,
		TenantID:      tenantID,
	}
	if req.Status != nil {
		rule.Status = *req.Status
	}
	if req.Version != nil {
		rule.Version = *req.Version
	}
	return rule
}

func buildEduScheduleRuleFromUpdateRequest(req *models.EduScheduleRuleUpdateRequest, tenantID, userID uint) *models.EduScheduleRule {
	rule := buildEduScheduleRuleFromAddRequest(&req.EduScheduleRuleAddRequest, tenantID, userID)
	rule.BaseModel = models.BaseModel{ID: req.ID}
	return rule
}

func buildEduScheduleRuleListScope(req *models.EduScheduleRuleListRequest, tenantID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("tenant_id = ?", tenantID)
		if req == nil {
			return db
		}
		if req.ID != nil {
			db = db.Where("id = ?", *req.ID)
		}
		if req.RuleType != "" {
			db = db.Where("rule_type = ?", req.RuleType)
		}
		if req.RepeatType != "" {
			db = db.Where("repeat_type = ?", req.RepeatType)
		}
		if req.ClassID != nil {
			db = db.Where("class_id = ?", *req.ClassID)
		}
		if req.StudentID != nil {
			db = db.Where("student_id = ?", *req.StudentID)
		}
		if req.CourseID != nil {
			db = db.Where("course_id = ?", *req.CourseID)
		}
		if req.TeacherID != nil {
			db = db.Where("teacher_id = ?", *req.TeacherID)
		}
		if req.Status != nil {
			db = db.Where("status = ?", *req.Status)
		}
		return db
	}
}

func buildEduScheduleLessonListScope(req *models.EduLessonListRequest, tenantID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("tenant_id = ?", tenantID)
		if req == nil {
			return db
		}
		if req.ID != nil {
			db = db.Where("id = ?", *req.ID)
		}
		if req.RuleID != nil {
			db = db.Where("rule_id = ?", *req.RuleID)
		}
		if req.ClassID != nil {
			db = db.Where("class_id = ?", *req.ClassID)
		}
		if req.StudentID != nil {
			db = db.Where("student_id = ?", *req.StudentID)
		}
		if req.CourseID != nil {
			db = db.Where("course_id = ?", *req.CourseID)
		}
		if req.TeacherID != nil {
			db = db.Where("teacher_id = ?", *req.TeacherID)
		}
		if req.Status != "" {
			db = db.Where("status = ?", req.Status)
		}
		return db
	}
}
