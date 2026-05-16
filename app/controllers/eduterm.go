package controllers

import (
	"gin-fast/app/models"
	"gin-fast/app/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduTermController 学期控制器
type EduTermController struct {
	Common
	EduTermService *service.EduTermService
}

func NewEduTermController() *EduTermController {
	return &EduTermController{Common: Common{}, EduTermService: service.NewEduTermService()}
}

func (ctl *EduTermController) List(c *gin.Context) {
	var req models.EduTermListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	total, err := models.NewEduTermList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计学期数量失败", err)
	}
	list := models.NewEduTermList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取学期列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduTermController) Add(c *gin.Context) {
	var req models.EduTermAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	term := &models.EduTerm{
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    1,
		CreatedBy: ctl.GetCurrentUserID(c),
		TenantID:  ctl.RequireTenant(c),
	}
	if req.Status != nil {
		term.Status = int8(*req.Status)
	}
	if err := ctl.EduTermService.Create(c, term); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学期创建成功", term)
}

func (ctl *EduTermController) Update(c *gin.Context) {
	var req models.EduTermUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	term := &models.EduTerm{
		BaseModel: models.BaseModel{ID: uint(req.ID)},
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    1,
		CreatedBy: ctl.GetCurrentUserID(c),
		TenantID:  ctl.RequireTenant(c),
	}
	if req.Status != nil {
		term.Status = int8(*req.Status)
	}
	if err := ctl.EduTermService.Update(c, term); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学期更新成功", term)
}

func (ctl *EduTermController) Delete(c *gin.Context) {
	var req models.EduTermDeleteRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduTermService.Delete(c, ctl.RequireTenant(c), uint(req.ID)); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学期删除成功", nil)
}

func (ctl *EduTermController) ClosedDays(c *gin.Context) {
	termID, err := parseOptionalPathUint(c, "id")
	if err != nil {
		ctl.FailAndAbort(c, "学期ID格式错误", err)
	}
	rows, err := ctl.EduTermService.ClosedDays(c, ctl.RequireTenant(c), termID)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, gin.H{"list": rows})
}

func (ctl *EduTermController) SaveClosedDays(c *gin.Context) {
	var req models.EduTermClosedDaysSaveRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	termID, err := parseOptionalPathUint(c, "id")
	if err != nil {
		ctl.FailAndAbort(c, "学期ID格式错误", err)
	}
	rows := make([]models.EduTermClosedDay, 0, len(req.Rows))
	for _, row := range req.Rows {
		rows = append(rows, models.EduTermClosedDay{
			BaseModel:  models.BaseModel{ID: uint(row.ID)},
			TermID:     uint(row.TermID),
			ClosedDate: row.ClosedDate,
			Reason:     row.Reason,
		})
	}
	if err := ctl.EduTermService.SaveClosedDays(c, ctl.RequireTenant(c), termID, rows); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学期停课日保存成功", nil)
}
