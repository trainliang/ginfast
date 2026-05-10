package controllers

import (
	"gin-fast/app/models"
	"gin-fast/app/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduRoomController 教培教室/场地控制器
type EduRoomController struct {
	Common
	EduRoomService *service.EduRoomService
}

func NewEduRoomController() *EduRoomController {
	return &EduRoomController{Common: Common{}, EduRoomService: service.NewEduRoomService()}
}

func (ctl *EduRoomController) List(c *gin.Context) {
	var req models.EduRoomListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	total, err := models.NewEduRoomList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计场地数量失败", err)
	}
	list := models.NewEduRoomList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取场地列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduRoomController) Options(c *gin.Context) {
	tenantID := ctl.GetCurrentTenantID(c)
	list := models.NewEduRoomList()
	if err := list.Find(c, func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ? AND status = ?", tenantID, 1).Order("id asc")
	}); err != nil {
		ctl.FailAndAbort(c, "获取场地选项失败", err)
	}
	ctl.Success(c, gin.H{"list": list})
}

func (ctl *EduRoomController) GetByID(c *gin.Context) {
	var req models.EduRoomGetRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	room := models.NewEduRoom()
	if err := room.Find(c, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ? AND tenant_id = ?", req.ID, ctl.GetCurrentTenantID(c))
	}); err != nil {
		ctl.FailAndAbort(c, "查询场地失败", err)
	}
	if room.IsEmpty() {
		ctl.FailAndAbort(c, "场地不存在", nil)
	}
	ctl.Success(c, room)
}

func (ctl *EduRoomController) Add(c *gin.Context) {
	var req models.EduRoomAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	room := buildRoomFromAddRequest(&req, ctl.GetCurrentTenantID(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduRoomService.Create(c, room); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "场地创建成功", room)
}

func (ctl *EduRoomController) Update(c *gin.Context) {
	var req models.EduRoomUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	room := buildRoomFromUpdateRequest(&req, ctl.GetCurrentTenantID(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduRoomService.Update(c, room); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "场地更新成功", room)
}

func (ctl *EduRoomController) Delete(c *gin.Context) {
	var req models.EduRoomDeleteRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduRoomService.Delete(c, ctl.GetCurrentTenantID(c), req.ID); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "场地删除成功", nil)
}

func (ctl *EduRoomController) WeeklyRules(c *gin.Context) {
	var req models.EduRoomWeeklyRuleListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if roomID, err := parseOptionalPathUint(c, "id"); err != nil {
		ctl.FailAndAbort(c, "场地ID格式错误", err)
	} else if roomID > 0 {
		req.RoomID = &roomID
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	list := models.NewEduRoomWeeklyRuleList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取场地周规则失败", err)
	}
	ctl.Success(c, gin.H{"list": list})
}

func (ctl *EduRoomController) SaveWeeklyRules(c *gin.Context) {
	var req models.EduRoomWeeklyRuleImportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if roomID, err := parseOptionalPathUint(c, "id"); err != nil {
		ctl.FailAndAbort(c, "场地ID格式错误", err)
	} else if roomID > 0 {
		for i := range req.Rows {
			req.Rows[i].RoomID = roomID
		}
	}
	result, err := ctl.EduRoomService.SaveWeeklyRules(c, ctl.GetCurrentTenantID(c), req.Rows)
	if err != nil && len(result.Errors) == 0 {
		ctl.FailAndAbort(c, "保存场地周规则失败", err)
	}
	ctl.Success(c, result)
}

func (ctl *EduRoomController) Exceptions(c *gin.Context) {
	var req models.EduRoomExceptionListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if roomID, err := parseOptionalPathUint(c, "id"); err != nil {
		ctl.FailAndAbort(c, "场地ID格式错误", err)
	} else if roomID > 0 {
		req.RoomID = &roomID
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	list := models.NewEduRoomExceptionList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取场地例外日期失败", err)
	}
	ctl.Success(c, gin.H{"list": list})
}

func (ctl *EduRoomController) AddException(c *gin.Context) {
	var req models.EduRoomExceptionAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if roomID, err := parseOptionalPathUint(c, "id"); err != nil {
		ctl.FailAndAbort(c, "场地ID格式错误", err)
	} else if roomID > 0 {
		req.RoomID = roomID
	}
	exception := buildRoomExceptionFromAddRequest(&req, ctl.GetCurrentTenantID(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduRoomService.AddException(c, exception); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "场地例外日期创建成功", exception)
}

func (ctl *EduRoomController) UpdateException(c *gin.Context) {
	var req models.EduRoomExceptionUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if roomID, err := parseOptionalPathUint(c, "id"); err != nil {
		ctl.FailAndAbort(c, "场地ID格式错误", err)
	} else if roomID > 0 {
		req.RoomID = roomID
	}
	exception := buildRoomExceptionFromUpdateRequest(&req, ctl.GetCurrentTenantID(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduRoomService.UpdateException(c, exception); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "场地例外日期更新成功", exception)
}

func (ctl *EduRoomController) DeleteException(c *gin.Context) {
	var req models.EduRoomExceptionDeleteRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduRoomService.DeleteException(c, ctl.GetCurrentTenantID(c), req.ID); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "场地例外日期删除成功", nil)
}

func (ctl *EduRoomController) Import(c *gin.Context) {
	var req models.EduRoomImportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	result, err := ctl.EduRoomService.ImportRows(c, ctl.GetCurrentTenantID(c), req.Rows)
	if err != nil && len(result.Errors) == 0 {
		ctl.FailAndAbort(c, "导入场地失败", err)
	}
	ctl.Success(c, result)
}

func (ctl *EduRoomController) ImportWeeklyRules(c *gin.Context) {
	var req models.EduRoomWeeklyRuleImportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	result, err := ctl.EduRoomService.ImportWeeklyRows(c, ctl.GetCurrentTenantID(c), req.Rows)
	if err != nil && len(result.Errors) == 0 {
		ctl.FailAndAbort(c, "导入场地周规则失败", err)
	}
	ctl.Success(c, result)
}

func (ctl *EduRoomController) ImportExceptions(c *gin.Context) {
	var req models.EduRoomExceptionImportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	result, err := ctl.EduRoomService.ImportExceptionRows(c, ctl.GetCurrentTenantID(c), req.Rows)
	if err != nil && len(result.Errors) == 0 {
		ctl.FailAndAbort(c, "导入场地例外日期失败", err)
	}
	ctl.Success(c, result)
}

func (ctl *EduRoomController) Export(c *gin.Context) {
	var req models.EduRoomExportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	rows, err := ctl.EduRoomService.ExportRows(c, ctl.GetCurrentTenantID(c), req.IDs)
	if err != nil {
		ctl.FailAndAbort(c, "导出场地失败", err)
	}
	ctl.Success(c, gin.H{"list": rows})
}

func buildRoomFromAddRequest(req *models.EduRoomAddRequest, tenantID, userID uint) *models.EduRoom {
	room := &models.EduRoom{Name: req.Name, Code: req.Code, Type: req.Type, Capacity: req.Capacity, Location: req.Location, Status: 1, Remark: req.Remark, CreatedBy: userID, TenantID: tenantID}
	if req.Status != nil {
		room.Status = *req.Status
	}
	return room
}

func buildRoomFromUpdateRequest(req *models.EduRoomUpdateRequest, tenantID, userID uint) *models.EduRoom {
	room := &models.EduRoom{BaseModel: models.BaseModel{ID: req.ID}, Name: req.Name, Code: req.Code, Type: req.Type, Capacity: req.Capacity, Location: req.Location, Status: 1, Remark: req.Remark, CreatedBy: userID, TenantID: tenantID}
	if req.Status != nil {
		room.Status = *req.Status
	}
	return room
}

func buildRoomExceptionFromAddRequest(req *models.EduRoomExceptionAddRequest, tenantID, userID uint) *models.EduRoomException {
	return &models.EduRoomException{RoomID: req.RoomID, ExceptionDate: req.ExceptionDate, Type: req.Type, StartTime: req.StartTime, EndTime: req.EndTime, Reason: req.Reason, CreatedBy: userID, TenantID: tenantID}
}

func buildRoomExceptionFromUpdateRequest(req *models.EduRoomExceptionUpdateRequest, tenantID, userID uint) *models.EduRoomException {
	return &models.EduRoomException{BaseModel: models.BaseModel{ID: req.ID}, RoomID: req.RoomID, ExceptionDate: req.ExceptionDate, Type: req.Type, StartTime: req.StartTime, EndTime: req.EndTime, Reason: req.Reason, CreatedBy: userID, TenantID: tenantID}
}
